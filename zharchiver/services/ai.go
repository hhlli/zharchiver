package services

import (
	"bytes"
	"crypto/tls"
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
	Tag         string `json:"tag,omitempty"`
}

func getExistingTagsStr(db *sql.DB) string {
	rows, err := db.Query("SELECT DISTINCT tag FROM answers WHERE tag != ''")
	if err != nil {
		return ""
	}
	defer rows.Close()
	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err == nil {
			tags = append(tags, tag)
		}
	}
	if len(tags) == 0 {
		return ""
	}
	return strings.Join(tags, ", ")
}

// SuggestTagByAI 用于纯文本的自动打标
func SuggestTagByAI(db *sql.DB, title, cleanContent string) string {
	baseURL := models.GetSetting(db, "ai_base_url")
	apiKey := models.GetSetting(db, "ai_api_key")
	modelName := models.GetSetting(db, "ai_model_name")

	if baseURL == "" || apiKey == "" || modelName == "" {
		return ""
	}

	tagsStr := getExistingTagsStr(db)
	
	// 截取最多 1500 字符防止内容过长
	if len(cleanContent) > 1500 {
		runes := []rune(cleanContent)
		if len(runes) > 1500 {
			cleanContent = string(runes[:1500])
		}
	}

	prompt := fmt.Sprintf(`你是一个归档整理助手。请阅读以下文章标题和正文摘要，并选择最匹配的标签。
候选标签列表：[%s]。
如果都不合适，请结合内容自己提炼一个不超过4个字的新标签。
请严格输出纯 JSON 对象，不要带Markdown格式，结构如下：
{"recommended_tag": "..."}

文章标题：%s
正文摘要：%s`, tagsStr, title, cleanContent)

	reqBody := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		return ""
	}
	defer resp.Body.Close()

	var aiResp aiVisionResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil || len(aiResp.Choices) == 0 {
		return ""
	}

	contentStr := strings.TrimSpace(aiResp.Choices[0].Message.Content)
	if strings.HasPrefix(contentStr, "```json") {
		contentStr = strings.TrimPrefix(contentStr, "```json")
		contentStr = strings.TrimSuffix(contentStr, "```")
	} else if strings.HasPrefix(contentStr, "```") {
		contentStr = strings.TrimPrefix(contentStr, "```")
		contentStr = strings.TrimSuffix(contentStr, "```")
	}

	var extracted struct {
		RecommendedTag string `json:"recommended_tag"`
	}
	if err := json.Unmarshal([]byte(contentStr), &extracted); err == nil {
		return extracted.RecommendedTag
	}
	return ""
}

func ProcessImageArchiveTask(db *sql.DB, imgBytes []byte) (*models.AnswerData, error) {
	baseURL := models.GetSetting(db, "ai_base_url")
	apiKey := models.GetSetting(db, "ai_api_key")
	modelName := models.GetSetting(db, "ai_model_name")
	autoCategorize := models.GetSetting(db, "auto_categorization_enabled") == "true"

	if baseURL == "" || apiKey == "" || modelName == "" {
		return nil, errors.New("AI 助手未完全配置，请在设置中完善信息")
	}

	mimeType := http.DetectContentType(imgBytes)
	base64Img := base64.StdEncoding.EncodeToString(imgBytes)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Img)

	prompt := `你是一个专门提取截图信息的内容提取器。请识别截图中的问答内容。
提取出：
1. 标题(如有)
2. 作者名称(如有)
3. 回答正文(转化为基本的HTML格式，如 <p>, <b> 等，不要使用 Markdown)`

	jsonStructure := `{
  "title": "...",
  "author_name": "...",
  "content_html": "..."`

	if autoCategorize {
		tagsStr := getExistingTagsStr(db)
		prompt += fmt.Sprintf("\n4. 标签：请结合内容，从已有标签列表 [%s] 中选择最匹配的一个。如果都不合适，请提炼一个不超过4个字的新标签。", tagsStr)
		jsonStructure += ",\n  \"tag\": \"...\""
	}
	jsonStructure += "\n}"

	prompt += fmt.Sprintf("\n\n请严格输出纯 JSON 对象，格式必须完全符合以下结构，不要输出任何其他解释文字或 Markdown 代码块：\n%s", jsonStructure)

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

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}
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
	
	tagToSave := extracted.Tag
	if tagToSave == "" && autoCategorize {
		tagToSave = "未分类"
	}
	
	data := &models.AnswerData{
		AnswerID:    uuid,
		QuestionID:  "0",
		Title:       extracted.Title,
		AuthorName:  extracted.AuthorName,
		ContentHTML: extracted.ContentHTML,
		CreatedTime: time.Now().Unix(),
		UpdatedTime: time.Now().Unix(),
		Tag:         tagToSave,
		TagColor:    "blue",
	}

	err = models.SaveAnswer(db, data)
	if err != nil {
		return nil, fmt.Errorf("存入数据库失败: %v", err)
	}

	return data, nil
}
