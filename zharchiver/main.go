package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/tidwall/gjson"
	_ "modernc.org/sqlite"
)

var zhihuRegex = regexp.MustCompile(`zhihu\.com/question/(\d+)/answer/(\d+)`)

const dC0Value = "kdOYbtmYmhyPTuJmX8VHlc5ERZxJtktiX88=|1784191170" // 仅作为备用

type ZhihuTarget struct {
	QuestionID string
	AnswerID   string
	CleanURL   string
}

type AnswerData struct {
	QuestionID  string   `json:"question_id"`
	AnswerID    string   `json:"answer_id"`
	Title       string   `json:"title"`
	AuthorName  string   `json:"author_name"`
	ContentHTML string   `json:"content_html"`
	CreatedTime int64    `json:"created_time"`
	UpdatedTime int64    `json:"updated_time"`
	ImageURLs   []string `json:"image_urls"`
	Tag         string   `json:"tag"`
	TagColor    string   `json:"tag_color"`
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
	BroadcastLog("INFO", "正在初始化无头浏览器环境...")
	l := launcher.New().Headless(true)
	l.Set("disable-blink-features", "AutomationControlled")
	l.Set("disable-web-security")

	u, err := l.Launch()
	if err != nil {
		BroadcastLog("ERROR", fmt.Sprintf("浏览器启动失败, 错误详情: %v", err))
		return "", fmt.Errorf("启动浏览器失败: %w", err)
	}

	BroadcastLog("INFO", "成功启动浏览器进程，正在连接控制端口...")
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

	BroadcastLog("INFO", fmt.Sprintf("开始导航至目标页面: %s", targetURL))
	if err := page.Navigate(targetURL); err != nil {
		BroadcastLog("ERROR", fmt.Sprintf("页面导航失败，可能由于网络问题或超时: %v", err))
		return "", fmt.Errorf("页面导航失败: %w", err)
	}

	BroadcastLog("INFO", "页面加载中，正在等待 #js-initialData 核心节点出现...")
	el, err := page.Element("#js-initialData")
	if err != nil {
		BroadcastLog("ERROR", fmt.Sprintf("未能找到核心数据节点 #js-initialData，可能被反爬或链接已失效: %v", err))
		return "", fmt.Errorf("等待 #js-initialData 节点超时或失败: %w", err)
	}

	BroadcastLog("INFO", "成功提取到核心 JSON 字符串数据")
	return el.Text()
}

func parseInitialJSON(jsonData string, target *ZhihuTarget) (*AnswerData, error) {
	BroadcastLog("INFO", fmt.Sprintf("开始解析 JSON 数据，目标问题 ID: %s, 回答 ID: %s", target.QuestionID, target.AnswerID))
	answerPath := fmt.Sprintf("initialState.entities.answers.%s", target.AnswerID)
	questionPath := fmt.Sprintf("initialState.entities.questions.%s", target.QuestionID)

	answerNode := gjson.Get(jsonData, answerPath)
	if !answerNode.Exists() {
		BroadcastLog("ERROR", fmt.Sprintf("在 JSON 中未找到对应回答实体 [%s]，可能该回答不存在、被屏蔽或需要登录", target.AnswerID))
		return nil, fmt.Errorf("JSON 中未找到回答 ID [%s] 的实体数据", target.AnswerID)
	}

	questionNode := gjson.Get(jsonData, questionPath)
	contentHTML := answerNode.Get("content").String()

	imgRegex := regexp.MustCompile(`(?:data-actualsrc|data-original)="([^"]+)"`)
	imgMatches := imgRegex.FindAllStringSubmatch(contentHTML, -1)

	var imageURLs []string
	seen := make(map[string]bool)

	for _, m := range imgMatches {
		if len(m) > 1 {
			rawURL := m[1]
			if !seen[rawURL] {
				seen[rawURL] = true
				imageURLs = append(imageURLs, rawURL)
			}
		}
	}

	// 清除占位图 src，防止双 src 冲突
	placeholderRegex := regexp.MustCompile(`src="data:image[^"]+"`)
	contentHTML = placeholderRegex.ReplaceAllString(contentHTML, "")

	data := &AnswerData{
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

func initDB(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %v", err)
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	
	// 启用 WAL 模式提高并发读写性能
	_, _ = db.Exec("PRAGMA journal_mode=WAL")
	// 限制单个 SQLite 文件的最大打开连接数为 1，防止 database is locked 错误
	db.SetMaxOpenConns(1)

	schema := `
	CREATE TABLE IF NOT EXISTS answers (
		answer_id TEXT PRIMARY KEY,
		question_id TEXT NOT NULL,
		title TEXT NOT NULL,
		author_name TEXT,
		author_avatar TEXT,
		content_html TEXT,
		created_time INTEGER,
		updated_time INTEGER,
		saved_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	// 数据库迁移: 增加新字段 (忽略错误，因为如果列已存在会报错)
	_, _ = db.Exec(`ALTER TABLE answers ADD COLUMN tag TEXT DEFAULT ''`)
	_, _ = db.Exec(`ALTER TABLE answers ADD COLUMN tag_color TEXT DEFAULT 'blue'`)
	_, _ = db.Exec(`ALTER TABLE answers ADD COLUMN author_avatar TEXT DEFAULT ''`)

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS settings (
            key TEXT PRIMARY KEY,
            value TEXT NOT NULL
        );
    `)
    if err != nil {
        return nil, fmt.Errorf("创建 settings 表失败: %v", err)
    }

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			answer_id TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("创建 comments 表失败: %v", err)
	}

	return db, nil
}

func saveAnswer(db *sql.DB, data *AnswerData) error {
	if data.TagColor == "" {
		data.TagColor = "blue"
	}
	query := `
	INSERT INTO answers (answer_id, question_id, title, author_name, content_html, created_time, updated_time, tag, tag_color)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(answer_id) DO UPDATE SET
		title = excluded.title,
		author_name = excluded.author_name,
		content_html = excluded.content_html,
		updated_time = excluded.updated_time,
		tag = CASE WHEN excluded.tag != '' THEN excluded.tag ELSE answers.tag END,
		tag_color = CASE WHEN excluded.tag != '' THEN excluded.tag_color ELSE answers.tag_color END,
		saved_at = CURRENT_TIMESTAMP;
	`

	_, err := db.Exec(query,
		data.AnswerID,
		data.QuestionID,
		data.Title,
		data.AuthorName,
		data.ContentHTML,
		data.CreatedTime,
		data.UpdatedTime,
		data.Tag,
		data.TagColor,
	)

	return err
}

func main() {
	db, err := initDB("db/zharchiver.db")
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	startServer(db, 8080)
}