package handlers

import (
	"encoding/json"
	"net/http"

	"zharchiver/models"
)

func (env *HandlerEnv) handleGetPreferences(w http.ResponseWriter, r *http.Request) {
	settings := map[string]string{
		"auto_push_enabled": models.GetSetting(env.db, "auto_push_enabled"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (env *HandlerEnv) handleSavePreferences(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求格式", http.StatusBadRequest)
		return
	}
	
	if val, ok := req["auto_push_enabled"]; ok {
		models.SetSetting(env.db, "auto_push_enabled", val)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
	})
}
