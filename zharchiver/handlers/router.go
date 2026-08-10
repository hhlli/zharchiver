package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"zharchiver/models"
	"zharchiver/utils"
)

type HandlerEnv struct {
	db *sql.DB
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/api/") && !strings.HasPrefix(r.URL.Path, "/storage/") {
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/api/auth/login" || r.URL.Path == "/api/auth/status" {
			next.ServeHTTP(w, r)
			return
		}

		enabled := models.GetSetting(db, "is_password_enabled")
		if enabled != "true" {
			next.ServeHTTP(w, r)
			return
		}

		tokenStr := r.Header.Get("Authorization")
		if tokenStr != "" && strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		} else {
			tokenStr = r.URL.Query().Get("token")
		}

		if tokenStr == "" {
			http.Error(w, "未授权访问", http.StatusUnauthorized)
			return
		}

		apiKey := models.GetSetting(db, "api_key")
		if apiKey != "" && tokenStr == apiKey {
			next.ServeHTTP(w, r)
			return
		}

		secret, err := getJWTSecret(db)
		if err != nil {
			http.Error(w, "系统错误", http.StatusInternalServerError)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token 无效或已过期", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" || strings.HasSuffix(r.URL.Path, ".json") {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		} else if strings.HasPrefix(r.URL.Path, "/assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		next.ServeHTTP(w, r)
	})
}

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) http.Handler {
	env := &HandlerEnv{db: db}

	// 静态文件服务
	mux.Handle("/", cacheControlMiddleware(http.FileServer(http.Dir("dist"))))
	mux.Handle("/storage/", http.StripPrefix("/storage/", http.FileServer(http.Dir("storage"))))

	mux.HandleFunc("GET /api/answers", env.handleGetAnswers)
	mux.HandleFunc("GET /api/group/answers", env.handleGetGroupAnswers)
	mux.HandleFunc("GET /api/answers/{id}", env.handleGetAnswerByID)
	mux.HandleFunc("DELETE /api/answers/{id}", env.handleDeleteAnswer)
	mux.HandleFunc("POST /api/answers/{id}/favorite", env.handleToggleFavorite)
	mux.HandleFunc("DELETE /api/group/answers", env.handleDeleteGroup)
	mux.HandleFunc("POST /api/group/answers/favorite", env.handleToggleGroupFavorite)
	mux.HandleFunc("PATCH /api/answers/{id}/tag", env.handleUpdateTag)
	mux.HandleFunc("PATCH /api/answers/{id}/content", env.handleUpdateAnswer)
	mux.HandleFunc("GET /api/tags", env.handleGetAllTags)
	mux.HandleFunc("PUT /api/tags", env.handleRenameGlobalTag)
	mux.HandleFunc("DELETE /api/tags", env.handleDeleteGlobalTag)
	mux.HandleFunc("GET /api/answers/{id}/comments", env.handleGetComments)
	mux.HandleFunc("POST /api/answers/{id}/comments", env.handleAddComment)
	mux.HandleFunc("PATCH /api/answers/{id}/comments/{comment_id}", env.handleUpdateComment)
	mux.HandleFunc("DELETE /api/answers/{id}/comments/{comment_id}", env.handleDeleteComment)
	mux.HandleFunc("POST /api/answers/{id}/share/telegram", env.handleShareToTelegram)
	mux.HandleFunc("POST /api/archive", env.handleArchive)
	mux.HandleFunc("GET /api/media", env.handleGetMedia)
	mux.HandleFunc("GET /api/logs/ws", utils.WsLogHandler)
	
	mux.HandleFunc("GET /api/auth/status", env.handleAuthStatus)
	mux.HandleFunc("POST /api/auth/login", env.handleLogin)
	mux.HandleFunc("POST /api/auth/update", env.handleUpdatePassword)
	mux.HandleFunc("GET /api/auth/zhihu/qrcode/ws", env.WsZhihuQRHandler)
	
	mux.HandleFunc("GET /api/settings/api_key", env.handleGetApiKey)
	mux.HandleFunc("POST /api/settings/api_key/generate", env.handleGenerateApiKey)

	mux.HandleFunc("GET /api/backup/download", env.handleDownloadBackup)
	mux.HandleFunc("POST /api/backup/restore", env.handleRestoreBackup)
	mux.HandleFunc("POST /api/backup/telegram/send", env.handleSendTelegramBackup)
	mux.HandleFunc("GET /api/settings/backup", env.handleGetBackupSettings)
	mux.HandleFunc("POST /api/settings/backup", env.handleSaveBackupSettings)

	mux.HandleFunc("GET /api/settings/ai", env.handleGetAISettings)
	mux.HandleFunc("POST /api/settings/ai", env.handleSaveAISettings)
	mux.HandleFunc("POST /api/settings/ai/test", env.handleTestAIConnection)
	mux.HandleFunc("GET /api/settings/telegram", env.handleGetTelegramSettings)
	mux.HandleFunc("POST /api/settings/telegram", env.handleSaveTelegramSettings)
	mux.HandleFunc("POST /api/settings/telegram/test", env.handleTestTelegramConnection)
	mux.HandleFunc("GET /api/settings/preferences", env.handleGetPreferences)
	mux.HandleFunc("POST /api/settings/preferences", env.handleSavePreferences)

	return corsMiddleware(authMiddleware(db, mux))
}
