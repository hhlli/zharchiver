package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"zharchiver/models"
	"zharchiver/utils"
)

func (env *HandlerEnv) handleGetAISettings(w http.ResponseWriter, r *http.Request) {
	settings := map[string]string{
		"ai_base_url":        models.GetSetting(env.db, "ai_base_url"),
		"ai_api_key":         models.GetSetting(env.db, "ai_api_key"),
		"ai_model_name":      models.GetSetting(env.db, "ai_model_name"),
		"telegram_bot_token": models.GetSetting(env.db, "telegram_bot_token"),
		"telegram_chat_id":   models.GetSetting(env.db, "telegram_chat_id"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (env *HandlerEnv) handleSaveAISettings(w http.ResponseWriter, r *http.Request) {
	var settings map[string]string
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	models.SetSetting(env.db, "ai_base_url", settings["ai_base_url"])
	models.SetSetting(env.db, "ai_api_key", settings["ai_api_key"])
	models.SetSetting(env.db, "ai_model_name", settings["ai_model_name"])
	models.SetSetting(env.db, "telegram_bot_token", settings["telegram_bot_token"])
	models.SetSetting(env.db, "telegram_chat_id", settings["telegram_chat_id"])
	
	utils.BroadcastLog("INFO", "已更新工具(AI与Telegram)配置")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (env *HandlerEnv) handleTestAIConnection(w http.ResponseWriter, r *http.Request) {
	var settings map[string]string
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	
	baseURL := settings["ai_base_url"]
	apiKey := settings["ai_api_key"]
	modelName := settings["ai_model_name"]
	
	if baseURL == "" || apiKey == "" || modelName == "" {
		utils.BroadcastLog("WARN", "测试 AI 连通性失败：配置不完整")
		http.Error(w, "请填写完整配置", http.StatusBadRequest)
		return
	}

	utils.BroadcastLog("INFO", fmt.Sprintf("开始测试 AI 连通性 (模型: %s)...", modelName))

	reqBody := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": "Hi, this is a test connection. Reply exactly with 'OK'.",
			},
		},
	}
	
	jsonBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		utils.BroadcastLog("ERROR", "测试 AI 连通性失败：网络请求出错 ("+err.Error()+")")
		http.Error(w, "网络请求失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		utils.BroadcastLog("ERROR", fmt.Sprintf("测试 AI 连通性失败：状态码 %d", resp.StatusCode))
		http.Error(w, fmt.Sprintf("API 返回错误状态码 %d: %s", resp.StatusCode, string(respBody)), http.StatusBadRequest)
		return
	}
	
	utils.BroadcastLog("INFO", "AI 连通性测试成功")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
