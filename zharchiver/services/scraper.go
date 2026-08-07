package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"zharchiver/models"
	"zharchiver/utils"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/tidwall/gjson"
)

var zhihuRegex = regexp.MustCompile(`zhihu\.com/question/(\d+)/answer/(\d+)`)

const dC0Value = "kdOYbtmYmhyPTuJmX8VHlc5ERZxJtktiX88=|1784191170" // 仅作为备用

type ZhihuTarget struct {
	QuestionID string
	AnswerID   string
	CleanURL   string
}

func parseZhihuLink(input string) (*ZhihuTarget, error) {
	matches := zhihuRegex.FindStringSubmatch(input)
	if len(matches) < 3 {
		return nil, errors.New("未识别到有效的知乎回答链接")
	}
	return &ZhihuTarget{
		QuestionID: matches[1],
		AnswerID:   matches[2],
		CleanURL:   fmt.Sprintf("https://www.zhihu.com/question/%s/answer/%s", matches[1], matches[2]),
	}, nil
}

func fetchWithRod(targetURL, zhihuCookies string) (string, error) {
	utils.BroadcastLog("INFO", "正在初始化无头浏览器环境...")
	l := launcher.New().Headless(true)
	l.Set("disable-blink-features", "AutomationControlled")
	l.Set("disable-web-security")

	u, err := l.Launch()
	if err != nil {
		utils.BroadcastLog("ERROR", fmt.Sprintf("浏览器启动失败, 错误详情: %v", err))
		return "", fmt.Errorf("启动浏览器失败: %w", err)
	}

	utils.BroadcastLog("INFO", "成功启动浏览器进程，正在连接控制端口...")
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage()
	page = page.Timeout(20 * time.Second)

	_ = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})

	var cookies []*proto.NetworkCookieParam
	if zhihuCookies != "" {
		var rodCookies []*proto.NetworkCookie
		if err := json.Unmarshal([]byte(zhihuCookies), &rodCookies); err == nil {
			for _, c := range rodCookies {
				cookies = append(cookies, &proto.NetworkCookieParam{
					Name:     c.Name,
					Value:    c.Value,
					Domain:   c.Domain,
					Path:     c.Path,
					Secure:   c.Secure,
					HTTPOnly: c.HTTPOnly,
					SameSite: c.SameSite,
					Expires:  c.Expires,
				})
			}
		}
	}
	
	if len(cookies) == 0 {
		cookies = []*proto.NetworkCookieParam{
			{
				Name:   "d_c0",
				Value:  dC0Value,
				Domain: ".zhihu.com",
				Path:   "/",
			},
		}
	}

	_ = page.SetCookies(cookies)

	utils.BroadcastLog("INFO", fmt.Sprintf("开始导航至目标页面: %s", targetURL))
	if err := page.Navigate(targetURL); err != nil {
		utils.BroadcastLog("ERROR", fmt.Sprintf("页面导航失败，可能由于网络问题或超时: %v", err))
		return "", fmt.Errorf("页面导航失败: %w", err)
	}

	utils.BroadcastLog("INFO", "页面加载中，正在等待 #js-initialData 核心节点出现...")
	el, err := page.Element("#js-initialData")
	if err != nil {
		utils.BroadcastLog("ERROR", fmt.Sprintf("未能找到核心数据节点 #js-initialData，可能被反爬或链接已失效: %v", err))
		return "", fmt.Errorf("等待 #js-initialData 节点超时或失败: %w", err)
	}

	utils.BroadcastLog("INFO", "成功提取到核心 JSON 字符串数据")
	return el.Text()
}

func parseInitialJSON(jsonData string, target *ZhihuTarget) (*models.AnswerData, error) {
	utils.BroadcastLog("INFO", fmt.Sprintf("开始解析 JSON 数据，目标问题 ID: %s, 回答 ID: %s", target.QuestionID, target.AnswerID))
	answerPath := fmt.Sprintf("initialState.entities.answers.%s", target.AnswerID)
	questionPath := fmt.Sprintf("initialState.entities.questions.%s", target.QuestionID)

	answerNode := gjson.Get(jsonData, answerPath)
	if !answerNode.Exists() {
		utils.BroadcastLog("ERROR", fmt.Sprintf("在 JSON 中未找到对应回答实体 [%s]，可能该回答不存在、被屏蔽或需要登录", target.AnswerID))
		return nil, fmt.Errorf("JSON 中未找到回答 ID [%s] 的实体数据", target.AnswerID)
	}

	questionNode := gjson.Get(jsonData, questionPath)
	contentHTML := answerNode.Get("content").String()

	imgTagRegex := regexp.MustCompile(`<img[^>]+>`)
	originalRegex := regexp.MustCompile(`data-original="([^"]+)"`)
	actualRegex := regexp.MustCompile(`data-actualsrc="([^"]+)"`)

	var imageURLs []string
	seen := make(map[string]bool)

	for _, imgTag := range imgTagRegex.FindAllString(contentHTML, -1) {
		var rawURL string
		if m := originalRegex.FindStringSubmatch(imgTag); len(m) > 1 {
			rawURL = m[1]
		} else if m := actualRegex.FindStringSubmatch(imgTag); len(m) > 1 {
			rawURL = m[1]
		}
		
		if rawURL != "" && !seen[rawURL] {
			seen[rawURL] = true
			imageURLs = append(imageURLs, rawURL)
		}
	}

	// 清除占位图 src，防止双 src 冲突
	placeholderRegex := regexp.MustCompile(`src="data:image[^"]+"`)
	contentHTML = placeholderRegex.ReplaceAllString(contentHTML, "")

	data := &models.AnswerData{
		QuestionID:  target.QuestionID,
		AnswerID:    target.AnswerID,
		Title:       questionNode.Get("title").String(),
		AuthorName:  answerNode.Get("author.name").String(),
		ContentHTML: contentHTML,
		CreatedTime: answerNode.Get("createdTime").Int(),
		UpdatedTime: answerNode.Get("updatedTime").Int(),
		ImageURLs:   imageURLs,
	}

	return data, nil
}

func ProcessArchiveTask(db *sql.DB, url string, tag string) (*models.AnswerData, error) {
	utils.BroadcastLog("INFO", "=== 开始新的归档任务 ===")
	utils.BroadcastLog("INFO", fmt.Sprintf("目标链接: %s", url))
	utils.BroadcastLog("INFO", fmt.Sprintf("标签: %s", tag))

	if strings.Contains(url, "twitter.com") || strings.Contains(url, "x.com") {
		return ProcessTwitterTask(db, url, tag)
	}

    target, err := parseZhihuLink(url)
    if err != nil {
        utils.BroadcastLog("ERROR", fmt.Sprintf("链接解析失败: %v", err))
        return nil, err
    }
    utils.BroadcastLog("INFO", fmt.Sprintf("链接解析成功，识别到回答 ID: %s", target.AnswerID))

    utils.BroadcastLog("INFO", "启动浏览器引擎(Rod)抓取页面数据...")
    cookiesStr := models.GetSetting(db, "zhihu_cookies")
    jsonData, err := fetchWithRod(target.CleanURL, cookiesStr)
    if err != nil {
        utils.BroadcastLog("ERROR", fmt.Sprintf("页面抓取失败: %v", err))
        return nil, fmt.Errorf("抓取失败: %v", err)
    }

    utils.BroadcastLog("INFO", "页面抓取完成，开始提取核心 JSON 数据...")
    data, err := parseInitialJSON(jsonData, target)
    if err != nil {
        utils.BroadcastLog("ERROR", fmt.Sprintf("数据提取失败: %v", err))
        return nil, fmt.Errorf("解析失败: %v", err)
    }
    utils.BroadcastLog("INFO", fmt.Sprintf("数据提取成功，标题: [%s]，作者: [%s]", data.Title, data.AuthorName))

    if len(data.ImageURLs) > 0 {
        utils.BroadcastLog("INFO", fmt.Sprintf("发现 %d 张相关图片，启动本地化下载流程...", len(data.ImageURLs)))
        _ = ProcessImages(data)
        utils.BroadcastLog("INFO", "图片本地化处理完成")
    } else {
        utils.BroadcastLog("INFO", "当前回答无内嵌图片，跳过下载流程")
    }

    utils.BroadcastLog("INFO", "准备将归档数据写入 SQLite 数据库...")
   	data.Tag = tag

	// 保存到数据库
	if err := models.SaveAnswer(db, data); err != nil {
        utils.BroadcastLog("ERROR", fmt.Sprintf("数据库写入失败: %v", err))
        return nil, fmt.Errorf("保存失败: %v", err)
    }
    
    utils.BroadcastLog("INFO", "数据库写入成功，归档流程全部完成")

    // 触发自动推送
    go AutoPushToTelegram(db, data.AnswerID)

    return data, nil
}
