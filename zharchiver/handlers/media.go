package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"zharchiver/models"
)

type MediaItem struct {
	Type        string `json:"type"` // "image" or "video"
	URL         string `json:"url"`
	AnswerID    string `json:"answer_id"`
	SavedAt     string `json:"saved_at"`
}

func (env *HandlerEnv) handleGetMedia(w http.ResponseWriter, r *http.Request) {
	answers, err := models.GetAnswers(env.db)
	if err != nil {
		http.Error(w, "获取归档记录失败", http.StatusInternalServerError)
		return
	}

	var mediaList []MediaItem
	
	imageExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	videoExts := map[string]bool{".mp4": true, ".webm": true, ".mov": true, ".ogg": true}

	for _, ans := range answers {
		dir := filepath.Join("storage", "images", ans.AnswerID)
		entries, err := os.ReadDir(dir)
		if err != nil {
			// 目录不存在或者读取失败则跳过
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			
			var mediaType string
			if imageExts[ext] {
				mediaType = "image"
			} else if videoExts[ext] {
				mediaType = "video"
			} else {
				continue // 不支持的媒体类型
			}

			mediaList = append(mediaList, MediaItem{
				Type:     mediaType,
				URL:      "/storage/images/" + ans.AnswerID + "/" + entry.Name(),
				AnswerID: ans.AnswerID,
				SavedAt:  ans.SavedAt,
			})
		}
	}

	// 如果没有数据，返回空数组而不是 null
	if mediaList == nil {
		mediaList = []MediaItem{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mediaList)
}
