package handlers

import (
	"encoding/json"
	"net/http"

	"zharchiver/models"
)

func (env *HandlerEnv) handleGetPreferences(w http.ResponseWriter, r *http.Request) {
	settings := map[string]string{
		"auto_push_enabled":           models.GetSetting(env.db, "auto_push_enabled"),
		"auto_categorization_enabled": models.GetSetting(env.db, "auto_categorization_enabled"),
		"tag_sort_order":              models.GetSetting(env.db, "tag_sort_order"),
		"theme_mode":                  models.GetSetting(env.db, "theme_mode"),
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
	if val, ok := req["auto_categorization_enabled"]; ok {
		models.SetSetting(env.db, "auto_categorization_enabled", val)
	}
	if val, ok := req["tag_sort_order"]; ok {
		models.SetSetting(env.db, "tag_sort_order", val)
	}
	if val, ok := req["theme_mode"]; ok {
		models.SetSetting(env.db, "theme_mode", val)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
	})
}
