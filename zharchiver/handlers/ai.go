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
	profilesJSON := models.GetSetting(env.db, "ai_profiles")
	activeProfileID := models.GetSetting(env.db, "ai_active_profile_id")
	
	// 自动向下兼容迁移
	if profilesJSON == "" {
		oldBaseURL := models.GetSetting(env.db, "ai_base_url")
		oldAPIKey := models.GetSetting(env.db, "ai_api_key")
		oldModelName := models.GetSetting(env.db, "ai_model_name")
		
		if oldBaseURL != "" || oldAPIKey != "" || oldModelName != "" {
			defaultProfile := map[string]string{
				"id": "default_1",
				"name": "默认配置",
				"base_url": oldBaseURL,
				"api_key": oldAPIKey,
				"model_name": oldModelName,
			}
			profilesList := []map[string]string{defaultProfile}
			bytes, _ := json.Marshal(profilesList)
			profilesJSON = string(bytes)
			activeProfileID = "default_1"
			
			models.SetSetting(env.db, "ai_profiles", profilesJSON)
			models.SetSetting(env.db, "ai_active_profile_id", activeProfileID)
		} else {
			profilesJSON = "[]"
		}
	}

	var profiles []map[string]string
	json.Unmarshal([]byte(profilesJSON), &profiles)
	if profiles == nil {
		profiles = []map[string]string{}
	}

	settings := map[string]interface{}{
		"ai_profiles":          profiles,
		"ai_active_profile_id": activeProfileID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (env *HandlerEnv) handleSaveAISettings(w http.ResponseWriter, r *http.Request) {
	var settings struct {
		Profiles        []map[string]string `json:"ai_profiles"`
		ActiveProfileID string              `json:"ai_active_profile_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	
	profilesBytes, _ := json.Marshal(settings.Profiles)
	models.SetSetting(env.db, "ai_profiles", string(profilesBytes))
	models.SetSetting(env.db, "ai_active_profile_id", settings.ActiveProfileID)
	
	// 同步活跃配置到全局旧配置键
	for _, p := range settings.Profiles {
		if p["id"] == settings.ActiveProfileID {
			models.SetSetting(env.db, "ai_base_url", p["base_url"])
			models.SetSetting(env.db, "ai_api_key", p["api_key"])
			models.SetSetting(env.db, "ai_model_name", p["model_name"])
			break
		}
	}
	
	utils.BroadcastLog("INFO", "已更新 AI 视觉模型配置")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (env *HandlerEnv) handleTestAIConnection(w http.ResponseWriter, r *http.Request) {
	var reqPayload struct {
		BaseURL   string `json:"base_url"`
		APIKey    string `json:"api_key"`
		ModelName string `json:"model_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	
	baseURL := reqPayload.BaseURL
	apiKey := reqPayload.APIKey
	modelName := reqPayload.ModelName
	
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
