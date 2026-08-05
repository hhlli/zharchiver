package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"zharchiver/models"
)

func sendMessageToTelegram(token, chatID, text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	
	payload := map[string]string{
		"chat_id": chatID,
		"text":    text,
	}
	body, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	client.Do(req)
}

func sendMessageToTelegramWithMarkup(token, chatID, text string, markup interface{}) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	if markup != nil {
		payload["reply_markup"] = markup
	}
	body, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	client.Do(req)
}

func editMessageText(token string, chatID int64, messageID int, text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText", token)
	
	payload := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       text,
	}
	body, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	client.Do(req)
}

func answerCallbackQuery(token, callbackQueryID string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/answerCallbackQuery", token)
	
	payload := map[string]interface{}{
		"callback_query_id": callbackQueryID,
	}
	body, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	client.Do(req)
}

func generateTagMarkup(db *sql.DB, answerID string) InlineKeyboardMarkup {
	tags := models.GetAllTags(db)
	var rows [][]InlineKeyboardButton
	var currentRow []InlineKeyboardButton
	
	for _, t := range tags {
		currentRow = append(currentRow, InlineKeyboardButton{
			Text:         t,
			CallbackData: fmt.Sprintf("set_tag:%s:%s", answerID, t),
		})
		if len(currentRow) == 3 {
			rows = append(rows, currentRow)
			currentRow = nil
		}
	}
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}
	
	rows = append(rows, []InlineKeyboardButton{
		{Text: "➕ 新增标签", CallbackData: fmt.Sprintf("new_tag:%s", answerID)},
	})
	
	return InlineKeyboardMarkup{InlineKeyboard: rows}
}

type TelegramPhotoSize struct {
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

func downloadTelegramFile(token, fileID string) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", token, fileID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var fileResp struct {
		Ok     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil || !fileResp.Ok {
		return nil, fmt.Errorf("getFile failed")
	}
	
	dlUrl := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", token, fileResp.Result.FilePath)
	dlResp, err := http.Get(dlUrl)
	if err != nil {
		return nil, err
	}
	defer dlResp.Body.Close()
	
	return io.ReadAll(dlResp.Body)
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type ForceReply struct {
	ForceReply bool `json:"force_reply"`
}

type TelegramMessage struct {
	MessageID int `json:"message_id"`
	Chat      struct {
		ID int64 `json:"id"`
	} `json:"chat"`
	Text           string              `json:"text"`
	Photo          []TelegramPhotoSize `json:"photo"`
	ReplyToMessage *TelegramMessage    `json:"reply_to_message,omitempty"`
}

type TelegramCallbackQuery struct {
	ID      string           `json:"id"`
	Data    string           `json:"data"`
	Message *TelegramMessage `json:"message"`
}

type TelegramUpdate struct {
	UpdateID      int                    `json:"update_id"`
	Message       *TelegramMessage       `json:"message"`
	CallbackQuery *TelegramCallbackQuery `json:"callback_query"`
}

type TelegramUpdateResponse struct {
	Ok     bool             `json:"ok"`
	Result []TelegramUpdate `json:"result"`
}

func StartTelegramBotListener(db *sql.DB) {
	go func() {
		lastUpdateID := 0
		client := &http.Client{Timeout: 35 * time.Second} // getUpdates has timeout=30

		for {
			token := models.GetSetting(db, "telegram_bot_token")
			authorizedChatIDStr := models.GetSetting(db, "telegram_chat_id")

			if token == "" || authorizedChatIDStr == "" {
				// 未配置 Telegram 机器人，休眠后重试
				time.Sleep(30 * time.Second)
				continue
			}

			authorizedChatID, err := strconv.ParseInt(authorizedChatIDStr, 10, 64)
			if err != nil {
				time.Sleep(30 * time.Second)
				continue
			}

			url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=30", token, lastUpdateID+1)
			
			resp, err := client.Get(url)
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			var updateResp TelegramUpdateResponse
			if err := json.Unmarshal(body, &updateResp); err != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			if !updateResp.Ok {
				time.Sleep(5 * time.Second)
				continue
			}

			for _, update := range updateResp.Result {
				if update.UpdateID > lastUpdateID {
					lastUpdateID = update.UpdateID
				}

				if update.CallbackQuery != nil {
					cq := update.CallbackQuery
					if cq.Message.Chat.ID != authorizedChatID {
						continue
					}
					
					data := cq.Data
					if strings.HasPrefix(data, "set_tag:") {
						parts := strings.SplitN(data, ":", 3)
						if len(parts) == 3 {
							answerID := parts[1]
							tag := parts[2]
							
							// 查找该标签已有的颜色
							var tagColor string
							err := db.QueryRow("SELECT tag_color FROM answers WHERE tag = ? AND tag_color != '' AND tag_color IS NOT NULL LIMIT 1", tag).Scan(&tagColor)
							if err != nil || tagColor == "" {
								tagColor = "blue" // 默认蓝色
							}
							
							db.Exec("UPDATE answers SET tag = ?, tag_color = ? WHERE answer_id = ?", tag, tagColor, answerID)
							editMessageText(token, cq.Message.Chat.ID, cq.Message.MessageID, fmt.Sprintf("✅ 归档成功并已添加标签：%s", tag))
						}
					} else if strings.HasPrefix(data, "new_tag:") {
						parts := strings.SplitN(data, ":", 2)
						if len(parts) == 2 {
							answerID := parts[1]
							markup := ForceReply{ForceReply: true}
							sendMessageToTelegramWithMarkup(token, strconv.FormatInt(cq.Message.Chat.ID, 10), fmt.Sprintf("为归档 %s 输入新标签：", answerID), markup)
						}
					}
					answerCallbackQuery(token, cq.ID)
					continue
				}

				if update.Message == nil {
					continue
				}

				if update.Message.Chat.ID != authorizedChatID {
					// 忽略非授权用户的消息
					continue
				}

				if update.Message.ReplyToMessage != nil && strings.HasPrefix(update.Message.ReplyToMessage.Text, "为归档 ") {
					replyText := update.Message.ReplyToMessage.Text
					parts := strings.Split(replyText, " ")
					if len(parts) >= 2 {
						answerID := parts[1]
						newTag := strings.TrimSpace(update.Message.Text)
						if newTag != "" {
							var tagColor string
							err := db.QueryRow("SELECT tag_color FROM answers WHERE tag = ? AND tag_color != '' AND tag_color IS NOT NULL LIMIT 1", newTag).Scan(&tagColor)
							if err != nil || tagColor == "" {
								tagColor = "blue"
							}
							db.Exec("UPDATE answers SET tag = ?, tag_color = ? WHERE answer_id = ?", newTag, tagColor, answerID)
							sendMessageToTelegram(token, authorizedChatIDStr, fmt.Sprintf("✅ 成功为归档 %s 添加新标签：%s", answerID, newTag))
						}
					}
					continue
				}

				if len(update.Message.Photo) > 0 {
					photo := update.Message.Photo[len(update.Message.Photo)-1]
					go func(fileID string) {
						sendMessageToTelegram(token, authorizedChatIDStr, "收到图片，正在召唤 AI 视觉提取...")
						imgBytes, err := downloadTelegramFile(token, fileID)
						if err != nil {
							sendMessageToTelegram(token, authorizedChatIDStr, "❌ 图片下载失败")
							return
						}
						
						data, err := ProcessImageArchiveTask(db, imgBytes)
						if err != nil {
							sendMessageToTelegram(token, authorizedChatIDStr, fmt.Sprintf("❌ 提取归档失败：\n%v", err))
						} else {
							markup := generateTagMarkup(db, data.AnswerID)
							sendMessageToTelegramWithMarkup(token, authorizedChatIDStr, fmt.Sprintf("✅ 视觉提取并归档成功！\n标题：%s\n作者：%s", data.Title, data.AuthorName), markup)
						}
					}(photo.FileID)
					continue
				}

				text := strings.TrimSpace(update.Message.Text)
				if text == "" {
					continue
				}

				if strings.Contains(text, "zhihu.com") {
					// 发送收到反馈
					go func(targetUrl string) {
						sendMessageToTelegram(token, authorizedChatIDStr, "收到链接，正在为您归档...")
						
						data, err := ProcessArchiveTask(db, targetUrl, "")
						if err != nil {
							sendMessageToTelegram(token, authorizedChatIDStr, fmt.Sprintf("❌ 归档失败：\n%v", err))
						} else {
							markup := generateTagMarkup(db, data.AnswerID)
							sendMessageToTelegramWithMarkup(token, authorizedChatIDStr, fmt.Sprintf("✅ 归档成功！\n标题：%s\n作者：%s", data.Title, data.AuthorName), markup)
						}
					}(text)
				}
			}
		}
	}()
}
