package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"zharchiver/models"
	"zharchiver/utils"
)

type FxTwitterResponse struct {
	Code    int `json:"code"`
	Message string `json:"message"`
	Tweet   *FxTweet `json:"tweet"`
}

type FxTweet struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Author    *FxAuthor `json:"author"`
	Media     *FxMedia  `json:"media"`
	CreatedAt string    `json:"created_at"`
	Timestamp int64     `json:"created_timestamp"`
}

type FxAuthor struct {
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	AvatarURL  string `json:"avatar_url"`
}

type FxMedia struct {
	All []FxMediaItem `json:"all"`
}

type FxMediaItem struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// 提取 Tweet ID
func extractTweetID(urlStr string) (string, error) {
	// 匹配 twitter.com/xxx/status/12345 或 x.com/xxx/status/12345
	re := regexp.MustCompile(`(?:twitter\.com|x\.com)/(?:[^/]+)/status/(\d+)`)
	matches := re.FindStringSubmatch(urlStr)
	if len(matches) < 2 {
		return "", fmt.Errorf("无法从链接中提取推文 ID: %s", urlStr)
	}
	return matches[1], nil
}

// ProcessTwitterTask 处理 X 平台的解析归档
func ProcessTwitterTask(db *sql.DB, urlStr string, tag string) (*models.AnswerData, error) {
	tweetID, err := extractTweetID(urlStr)
	if err != nil {
		utils.BroadcastLog("ERROR", err.Error())
		return nil, err
	}

	utils.BroadcastLog("INFO", fmt.Sprintf("已识别推文 ID: %s, 正在请求解析接口...", tweetID))

	// 使用 fxtwitter 的公开 API 进行无感极速解析
	apiURL := fmt.Sprintf("https://api.fxtwitter.com/i/status/%s", tweetID)
	resp, err := http.Get(apiURL)
	if err != nil {
		utils.BroadcastLog("ERROR", fmt.Sprintf("请求解析接口失败: %v", err))
		return nil, fmt.Errorf("解析请求失败: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var fxResp FxTwitterResponse
	if err := json.Unmarshal(bodyBytes, &fxResp); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}

	if fxResp.Code != 200 || fxResp.Tweet == nil {
		return nil, fmt.Errorf("推文不存在、已被删除或账号私密 (Code: %d, Msg: %s)", fxResp.Code, fxResp.Message)
	}

	tweet := fxResp.Tweet
	authorName := "Unknown"
	if tweet.Author != nil {
		authorName = fmt.Sprintf("%s (@%s)", tweet.Author.Name, tweet.Author.ScreenName)
	}

	title := tweet.Text
	// 将换行替换为空格用于标题
	title = strings.ReplaceAll(title, "\n", " ")
	if title == "" {
		title = "无文字内容的推文"
	}

	// 内容转为简单的 HTML
	contentHTML := "<p>" + strings.ReplaceAll(tweet.Text, "\n", "<br>") + "</p>"

	// 提取媒体链接
	var mediaURLs []string
	if tweet.Media != nil && len(tweet.Media.All) > 0 {
		for _, item := range tweet.Media.All {
			mediaURLs = append(mediaURLs, item.URL)
			// 如果需要将图片嵌入到正文
			if item.Type == "video" || item.Type == "gif" {
				contentHTML += fmt.Sprintf(`<div style="margin-top: 10px;"><video controls src="%s" style="max-width: 100%%; border-radius: 8px;"></video></div>`, item.URL)
			} else {
				contentHTML += fmt.Sprintf(`<div style="margin-top: 10px;"><img src="%s" style="max-width: 100%%; border-radius: 8px;" /></div>`, item.URL)
			}
		}
		utils.BroadcastLog("INFO", fmt.Sprintf("解析到 %d 个媒体文件", len(mediaURLs)))
	}

	data := &models.AnswerData{
		AnswerID:    tweet.ID,
		QuestionID:  "twitter", // X平台数据统一放在 twitter 分类下
		Title:       title,
		AuthorName:  authorName,
		ContentHTML: contentHTML,
		CreatedTime: tweet.Timestamp,
		UpdatedTime: time.Now().Unix(),
		ImageURLs:   mediaURLs,
		Tag:         tag,
		TagColor:    "blue",
	}

	utils.BroadcastLog("INFO", "推文内容提取成功！开始本地化存储...")

	// 1. 下载媒体文件
	if len(data.ImageURLs) > 0 {
		for i, remoteURL := range data.ImageURLs {
			localPath, err := downloadImage(remoteURL, data.AnswerID, i+1)
			if err != nil {
				utils.BroadcastLog("ERROR", fmt.Sprintf("媒体文件下载失败 [%s]: %v", remoteURL, err))
				continue
			}
			data.ContentHTML = strings.Replace(data.ContentHTML, remoteURL, localPath, -1)
			data.ImageURLs[i] = localPath
		}
	}

	// 2. 存入数据库
	err = models.SaveAnswer(db, data)
	if err != nil {
		utils.BroadcastLog("ERROR", fmt.Sprintf("X 平台内容存入数据库失败: %v", err))
		return nil, fmt.Errorf("保存失败: %v", err)
	}

	utils.BroadcastLog("INFO", "X 平台内容归档完毕并成功入库")

	// 触发自动推送
	go AutoPushToTelegram(db, data.AnswerID)

	return data, nil
}
