package services

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"zharchiver/models"
	"zharchiver/utils"
)

// stripHTML 去除 HTML 标签，保留换行
func stripHTML(htmlStr string) string {
	// 替换块级标签和换行为 \n
	htmlStr = strings.ReplaceAll(htmlStr, "<br>", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "<br/>", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "</p>", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "</div>", "\n")
	
	// 移除所有剩余标签
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(htmlStr, "")
	
	// 去除多余的连续空行
	reEmptyLines := regexp.MustCompile(`\n{3,}`)
	text = reEmptyLines.ReplaceAllString(text, "\n\n")
	
	return strings.TrimSpace(text)
}

// splitText 按长度切分文本
func splitText(text string, limit int) []string {
	var chunks []string
	runes := []rune(text)
	for i := 0; i < len(runes); i += limit {
		end := i + limit
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}

type MediaItem struct {
	Type              string `json:"type"`
	Media             string `json:"media"`
	Caption           string `json:"caption,omitempty"`
	SupportsStreaming bool   `json:"supports_streaming,omitempty"`
	Width             int    `json:"width,omitempty"`
	Height            int    `json:"height,omitempty"`
}

// sendMediaGroup 发送媒体组
func sendMediaGroup(apiEndpoint, token, chatID string, mediaPaths []string, caption string, replyToMessageID int) (int, error) {
	url := fmt.Sprintf("%s/bot%s/sendMediaGroup", apiEndpoint, token)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("chat_id", chatID)
	if replyToMessageID > 0 {
		writer.WriteField("reply_to_message_id", fmt.Sprintf("%d", replyToMessageID))
	}

	var mediaArr []MediaItem
	for i, path := range mediaPaths {
		ext := strings.ToLower(filepath.Ext(path))
		mediaType := "photo"
		if ext == ".mp4" || ext == ".mov" || ext == ".gif" { // TG treats GIF as animation/video
			mediaType = "video"
		}

		attachKey := fmt.Sprintf("media%d", i)
		
		item := MediaItem{
			Type:  mediaType,
			Media: fmt.Sprintf("attach://media%d", i),
		}
		
		if mediaType == "video" {
			item.SupportsStreaming = true
			
		localPath := filepath.Join("storage", path)
			w, h, err := utils.GetVideoDimensions(localPath)
			if err == nil && w > 0 && h > 0 {
				item.Width = w
				item.Height = h
			}
		}
		
		// 仅把标题/文字放在第一张媒体上
		if i == 0 && caption != "" {
			item.Caption = caption
		}
		mediaArr = append(mediaArr, item)

		// 注意: 本地路径是相对 storage 的
		localPath := filepath.Join("storage", path)

		file, err := os.Open(localPath)
		if err != nil {
			return 0, fmt.Errorf("读取文件失败 %s: %v", localPath, err)
		}

		part, err := writer.CreateFormFile(attachKey, filepath.Base(localPath))
		if err != nil {
			file.Close()
			return 0, err
		}
		io.Copy(part, file)
		file.Close()
	}

	mediaJSON, _ := json.Marshal(mediaArr)
	writer.WriteField("media", string(mediaJSON))

	err := writer.Close()
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Timeout: 15 * time.Minute, Transport: tr} // 上传大文件可能需要较长时间
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var tgResp struct {
		Ok     bool `json:"ok"`
		Result []TelegramMessage `json:"result"`
	}
	json.Unmarshal(respBody, &tgResp)

	if !tgResp.Ok {
		return 0, fmt.Errorf("sendMediaGroup 失败: %s", string(respBody))
	}

	if len(tgResp.Result) > 0 {
		return tgResp.Result[0].MessageID, nil
	}
	return 0, nil
}

// sendText 发送普通文本消息，超长自动切分
func sendText(apiEndpoint, token, chatID string, text string, replyToMessageID int) (int, error) {
	chunks := splitText(text, 4000)
	var lastMsgID int

	for i, chunk := range chunks {
		url := fmt.Sprintf("%s/bot%s/sendMessage", apiEndpoint, token)
		payload := map[string]interface{}{
			"chat_id": chatID,
			"text":    chunk,
		}
		
		// 如果有 reply_to，或者是切片的后续部分，进行回复串联
		if replyToMessageID > 0 {
			payload["reply_to_message_id"] = replyToMessageID
		} else if i > 0 && lastMsgID > 0 {
			payload["reply_to_message_id"] = lastMsgID
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client := &http.Client{Timeout: 10 * time.Second, Transport: tr}
		resp, err := client.Do(req)
		if err != nil {
			return 0, err
		}
		
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var tgResp struct {
			Ok     bool `json:"ok"`
			Result TelegramMessage `json:"result"`
		}
		json.Unmarshal(respBody, &tgResp)

		if !tgResp.Ok {
			return 0, fmt.Errorf("sendMessage 失败: %s", string(respBody))
		}
		lastMsgID = tgResp.Result.MessageID
	}
	
	return lastMsgID, nil
}

// ShareAnswerToTelegram 发送完整的归档内容到 TG
func ShareAnswerToTelegram(db *sql.DB, token, chatID string, data *models.AnswerData) error {
	apiEndpoint := models.GetTelegramAPIEndpoint(db)
	// 1. 准备文字内容
	header := fmt.Sprintf("📄 【%s】\n👤 作者：%s\n\n", data.Title, data.AuthorName)
	rawText := stripHTML(data.ContentHTML)
	fullText := header + rawText

	// 2. 切分媒体 (TG 的 MediaGroup 限制最多 10 个文件)
	var mediaChunks [][]string
	for i := 0; i < len(data.ImageURLs); i += 10 {
		end := i + 10
		if end > len(data.ImageURLs) {
			end = len(data.ImageURLs)
		}
		mediaChunks = append(mediaChunks, data.ImageURLs[i:end])
	}

	var lastMessageID int
	var err error

	// 场景 A: 没有媒体文件，直接发纯文本
	if len(mediaChunks) == 0 {
		_, err = sendText(apiEndpoint, token, chatID, fullText, 0)
		return err
	}

	// 场景 B: 有媒体文件
	// 处理附言长度限制 (TG Caption 限制 1024)
	caption := fullText
	remainingText := ""
	if len([]rune(fullText)) > 1000 {
		runes := []rune(fullText)
		caption = string(runes[:1000]) + "...\n(全文较长，见回复)"
		remainingText = string(runes[1000:])
	}

	// 循环发送媒体块
	for i, chunk := range mediaChunks {
		chunkCaption := ""
		if i == 0 {
			chunkCaption = caption
		}

		if len(chunk) == 1 {
			// TG 如果只有一个媒体，不能使用 MediaGroup，需要降级到 sendPhoto / sendVideo。
			// 为保持简单且复用逻辑，我们可以强行用 sendPhoto/sendVideo 接口，或在 sendMediaGroup 里处理。
			// 我们写一个辅助函数或使用现成的 API。这里直接处理：
			msgID, err := sendSingleMedia(apiEndpoint, token, chatID, chunk[0], chunkCaption, lastMessageID)
			if err != nil {
				return err
			}
			lastMessageID = msgID
		} else {
			// 正常 2~10 张图的专辑
			msgID, err := sendMediaGroup(apiEndpoint, token, chatID, chunk, chunkCaption, lastMessageID)
			if err != nil {
				return err
			}
			lastMessageID = msgID
		}
	}

	// 如果有未发完的纯文本（因为 caption 超长），则用 sendMessage 回复发送
	if remainingText != "" {
		_, err = sendText(apiEndpoint, token, chatID, "...接上文：\n"+remainingText, lastMessageID)
		if err != nil {
			return err
		}
	}

	return nil
}

// AutoPushToTelegram checks user settings and auto pushes an answer to TG if enabled
func AutoPushToTelegram(db *sql.DB, answerID string) {
	if models.GetSetting(db, "auto_push_enabled") != "true" {
		return
	}

	token := models.GetSetting(db, "telegram_push_bot_token")
	chatID := models.GetSetting(db, "telegram_push_chat_id")
	if token == "" || chatID == "" {
		// 如果推送机器人未单独配置，退化使用默认归档机器人
		token = models.GetSetting(db, "telegram_bot_token")
		chatID = models.GetSetting(db, "telegram_chat_id")
	}

	if token == "" || chatID == "" {
		utils.BroadcastLog("INFO", "自动推送未配置完整的 Push Token 或 Chat ID，且无默认机器人配置，已跳过")
		return
	}

	utils.BroadcastLog("INFO", "开启自动推送：准备将归档文章发送至 Telegram...")
	data, err := models.GetAnswerByID(db, answerID)
	if err != nil {
		return
	}

	// Populate ImageURLs from local disk
	imageDir := filepath.Join("storage", "images", answerID)
	entries, err := os.ReadDir(imageDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				data.ImageURLs = append(data.ImageURLs, filepath.Join("images", answerID, entry.Name()))
			}
		}
	}

	err = ShareAnswerToTelegram(db, token, chatID, data)
	if err != nil {
		utils.BroadcastLog("ERROR", "自动推送失败: " + err.Error())
	} else {
		utils.BroadcastLog("INFO", "✅ 自动推送成功！已将归档文章发送至指定的频道/群组")
	}
}

// sendSingleMedia 单图/单视频发送
func sendSingleMedia(apiEndpoint, token, chatID string, mediaPath string, caption string, replyToMessageID int) (int, error) {
	ext := strings.ToLower(filepath.Ext(mediaPath))
	method := "sendPhoto"
	fileKey := "photo"
	if ext == ".mp4" || ext == ".mov" || ext == ".gif" {
		method = "sendVideo"
		fileKey = "video"
	}

	url := fmt.Sprintf("%s/bot%s/%s", apiEndpoint, token, method)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("chat_id", chatID)
	
	localPath := filepath.Join("storage", mediaPath)
	
	if method == "sendVideo" {
		writer.WriteField("supports_streaming", "true")
		w, h, err := utils.GetVideoDimensions(localPath)
		if err == nil && w > 0 && h > 0 {
			writer.WriteField("width", fmt.Sprintf("%d", w))
			writer.WriteField("height", fmt.Sprintf("%d", h))
		}
	}
	if replyToMessageID > 0 {
		writer.WriteField("reply_to_message_id", fmt.Sprintf("%d", replyToMessageID))
	}
	if caption != "" {
		writer.WriteField("caption", caption)
	}

	file, err := os.Open(localPath)
	if err != nil {
		return 0, fmt.Errorf("读取单媒体失败 %s: %v", localPath, err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fileKey, filepath.Base(localPath))
	if err != nil {
		return 0, err
	}
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Timeout: 15 * time.Minute, Transport: tr}
	
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	respBody, _ := io.ReadAll(resp.Body)
	var tgResp struct {
		Ok     bool `json:"ok"`
		Result TelegramMessage `json:"result"`
	}
	json.Unmarshal(respBody, &tgResp)

	if !tgResp.Ok {
		return 0, fmt.Errorf("%s 失败: %s", method, string(respBody))
	}

	return tgResp.Result.MessageID, nil
}
