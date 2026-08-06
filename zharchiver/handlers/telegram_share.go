package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"zharchiver/models"
	"zharchiver/services"
	"zharchiver/utils"
)

func (env *HandlerEnv) handleShareToTelegram(w http.ResponseWriter, r *http.Request) {
	// 路径解析 /api/answers/{id}/share/telegram
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "无效的请求路径", http.StatusBadRequest)
		return
	}
	answerID := parts[3]

	// 1. 获取答卷详情
	data, err := models.GetAnswerByID(env.db, answerID)
	if err != nil {
		http.Error(w, "未找到该归档", http.StatusNotFound)
		return
	}

	// 从本地磁盘自动读取所有该归档关联的图片和视频 (因为数据库没有存储 ImageURLs)
	imageDir := filepath.Join("storage", "images", answerID)
	entries, err := os.ReadDir(imageDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				data.ImageURLs = append(data.ImageURLs, filepath.Join("images", answerID, entry.Name()))
			}
		}
	}

	// 2. 获取 Telegram 推送配置
	token := models.GetSetting(env.db, "telegram_push_bot_token")
	chatID := models.GetSetting(env.db, "telegram_push_chat_id")
	if token == "" || chatID == "" {
		http.Error(w, "请先在设置中配置 Telegram 推送机器人的 Bot Token 和 Chat ID", http.StatusBadRequest)
		return
	}

	// 3. 执行推送
	utils.BroadcastLog("INFO", "=== 开始手动推送归档到 Telegram ===")
	err = services.ShareAnswerToTelegram(token, chatID, data)
	if err != nil {
		utils.BroadcastLog("ERROR", "推送失败: "+err.Error())
		http.Error(w, "推送失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.BroadcastLog("INFO", "✅ 手动推送成功！")
	// 4. 返回成功
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
	})
}
