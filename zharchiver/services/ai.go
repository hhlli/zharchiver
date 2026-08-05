package services

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"zharchiver/models"
)

type aiVisionRequest struct {
	Model    string        `json:"model"`
	Messages []aiMessage   `json:"messages"`
}

type aiMessage struct {
	Role    string      `json:"role"`
	Content []aiContent `json:"content"`
}

type aiContent struct {
	Type     string      `json:"type"`
	Text     string      `json:"text,omitempty"`
	ImageURL *aiImageURL `json:"image_url,omitempty"`
}

type aiImageURL struct {
	URL string `json:"url"`
}

type aiVisionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type aiExtractedData struct {
	Title       string `json:"title"`
	AuthorName  string `json:"author_name"`
	ContentHTML string `json:"content_html"`
}

func ProcessImageArchiveTask(db *sql.DB, imgBytes []byte) (*models.AnswerData, error) {
	baseURL := models.GetSetting(db, "ai_base_url")
	apiKey := models.GetSetting(db, "ai_api_key")
	modelName := models.GetSetting(db, "ai_model_name")

	if baseURL == "" || apiKey == "" || modelName == "" {
		return nil, errors.New("AI 助手未完全配置，请在设置中完善信息")
	}

	base64Img := base64.StdEncoding.EncodeToString(imgBytes)
	dataURL := fmt.Sprintf("data:image/jpeg;base64,%s", base64Img)

	prompt := `你是一个专门提取截图信息的内容提取器。请识别截图中的问答内容。
提取出：
1. 标题(如有)
2. 作者名称(如有)
3. 回答正文(转化为基本的HTML格式，如 <p>, <b> 等，不要使用 Markdown)

请严格输出纯 JSON 对象，格式必须完全符合以下结构，不要输出任何其他解释文字或 Markdown 代码块：
{
  "title": "...",
  "author_name": "...",
  "content_html": "..."
}`

	reqBody := aiVisionRequest{
		Model: modelName,
		Messages: []aiMessage{
			{
				Role: "user",
				Content: []aiContent{
					{Type: "text", Text: prompt},
					{Type: "image_url", ImageURL: &aiImageURL{URL: dataURL}},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用大模型 API 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("大模型返回错误状态码 %d: %s", resp.StatusCode, string(respBody))
	}

	var aiResp aiVisionResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, fmt.Errorf("解析大模型响应失败: %v", err)
	}

	if len(aiResp.Choices) == 0 {
		return nil, errors.New("大模型未返回内容")
	}

	contentStr := aiResp.Choices[0].Message.Content

	// Clean up markdown block if present
	contentStr = strings.TrimSpace(contentStr)
	if strings.HasPrefix(contentStr, "```json") {
		contentStr = strings.TrimPrefix(contentStr, "```json")
		contentStr = strings.TrimSuffix(contentStr, "```")
	} else if strings.HasPrefix(contentStr, "```") {
		contentStr = strings.TrimPrefix(contentStr, "```")
		contentStr = strings.TrimSuffix(contentStr, "```")
	}
	contentStr = strings.TrimSpace(contentStr)

	var extracted aiExtractedData
	if err := json.Unmarshal([]byte(contentStr), &extracted); err != nil {
		return nil, fmt.Errorf("解析提取出的 JSON 失败: %v\n返回内容为: %s", err, contentStr)
	}

	// 存入数据库
	uuid := fmt.Sprintf("img_%d", time.Now().UnixNano())
	now := time.Now().Format("2006-01-02 15:04:05")
	
	_, err = db.Exec(`
		INSERT INTO answers (answer_id, question_id, title, author_name, author_avatar, content_html, created_time, updated_time, saved_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, uuid, 0, extracted.Title, extracted.AuthorName, "", extracted.ContentHTML, 0, 0, now)

	if err != nil {
		return nil, fmt.Errorf("存入数据库失败: %v", err)
	}

	return &models.AnswerData{
		AnswerID:    uuid,
		Title:       extracted.Title,
		AuthorName:  extracted.AuthorName,
		ContentHTML: extracted.ContentHTML,
	}, nil
}
