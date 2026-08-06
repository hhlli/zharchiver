package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"zharchiver/models"
	"zharchiver/utils"
)

func (env *HandlerEnv) handleGetTelegramSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"telegram_bot_token":      models.GetSetting(env.db, "telegram_bot_token"),
		"telegram_chat_id":        models.GetSetting(env.db, "telegram_chat_id"),
		"telegram_push_bot_token": models.GetSetting(env.db, "telegram_push_bot_token"),
		"telegram_push_chat_id":   models.GetSetting(env.db, "telegram_push_chat_id"),
		"telegram_api_endpoint":   models.GetSetting(env.db, "telegram_api_endpoint"),
	})
}

func (env *HandlerEnv) handleSaveTelegramSettings(w http.ResponseWriter, r *http.Request) {
	var settings map[string]string
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	
	models.SetSetting(env.db, "telegram_bot_token", settings["telegram_bot_token"])
	models.SetSetting(env.db, "telegram_chat_id", settings["telegram_chat_id"])
	models.SetSetting(env.db, "telegram_push_bot_token", settings["telegram_push_bot_token"])
	models.SetSetting(env.db, "telegram_push_chat_id", settings["telegram_push_chat_id"])
	models.SetSetting(env.db, "telegram_api_endpoint", settings["telegram_api_endpoint"])
	
	utils.BroadcastLog("INFO", "已更新 Telegram 机器人配置")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (env *HandlerEnv) handleTestTelegramConnection(w http.ResponseWriter, r *http.Request) {
	type TelegramTestReq struct {
		BotType      string `json:"bot_type"` // "archive" or "push"
		BotToken     string `json:"telegram_bot_token"`
		ChatID       string `json:"telegram_chat_id"`
		PushBotToken string `json:"telegram_push_bot_token"`
		PushChatID   string `json:"telegram_push_chat_id"`
		ApiEndpoint  string `json:"telegram_api_endpoint"`
	}

	var req TelegramTestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	token := req.BotToken
	chatID := req.ChatID
	msgText := "你好！这是来自 ZHArchiver 的归档机器人连通性测试消息 🎉"
	if req.BotType == "push" {
		token = req.PushBotToken
		chatID = req.PushChatID
		msgText = "你好！这是来自 ZHArchiver 的推送机器人连通性测试消息 🚀"
	}

	if token == "" || chatID == "" {
		http.Error(w, "请提供完整的 Bot Token 和 Chat ID", http.StatusBadRequest)
		return
	}

	apiEndpoint := req.ApiEndpoint
	if apiEndpoint == "" {
		apiEndpoint = "https://api.telegram.org"
	}
	if apiEndpoint[len(apiEndpoint)-1] == '/' {
		apiEndpoint = apiEndpoint[:len(apiEndpoint)-1]
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", apiEndpoint, token)
	reqBody := map[string]interface{}{
		"chat_id": chatID,
		"text":    msgText,
	}
	jsonBytes, _ := json.Marshal(reqBody)
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		http.Error(w, "请求 Telegram API 失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Telegram 返回错误状态码 (请检查 Token 或 Chat ID 是否正确)", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
