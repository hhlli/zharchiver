package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 获取或生成 JWT Secret
func getJWTSecret(db *sql.DB) ([]byte, error) {
	var secretHex string
	err := db.QueryRow("SELECT value FROM settings WHERE key = 'jwt_secret'").Scan(&secretHex)
	if err == sql.ErrNoRows {
		b := make([]byte, 32)
		rand.Read(b)
		secretHex = hex.EncodeToString(b)
		_, err = db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('jwt_secret', ?)", secretHex)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return hex.DecodeString(secretHex)
}

// 刷新 JWT Secret 使所有 Token 失效
func refreshJWTSecret(db *sql.DB) error {
	b := make([]byte, 32)
	rand.Read(b)
	secretHex := hex.EncodeToString(b)
	_, err := db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('jwt_secret', ?)", secretHex)
	return err
}

// 鉴权中间件
func authMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 放行 OPTIONS 预检请求
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// 放行前端静态文件 (以 /api/ 或 /storage/ 开头的才需要鉴权)
		if !strings.HasPrefix(r.URL.Path, "/api/") && !strings.HasPrefix(r.URL.Path, "/storage/") {
			next.ServeHTTP(w, r)
			return
		}

		// 放行登录和状态接口
		if r.URL.Path == "/api/auth/login" || r.URL.Path == "/api/auth/status" {
			next.ServeHTTP(w, r)
			return
		}

		// 检查是否开启了密码保护
		var enabled string
		err := db.QueryRow("SELECT value FROM settings WHERE key = 'is_password_enabled'").Scan(&enabled)
		if err == sql.ErrNoRows || enabled != "true" {
			next.ServeHTTP(w, r)
			return
		}

		// 校验 Token
		tokenStr := r.Header.Get("Authorization")
		if tokenStr != "" && strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		} else {
			// 支持从 Query 获取，兼容 WebSocket
			tokenStr = r.URL.Query().Get("token")
		}

		if tokenStr == "" {
			http.Error(w, "未授权访问", http.StatusUnauthorized)
			return
		}

		// 优先检查是否是永久 API Key
		var apiKey string
		_ = db.QueryRow("SELECT value FROM settings WHERE key = 'api_key'").Scan(&apiKey)
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

type UpdatePasswordReq struct {
	Enabled     bool   `json:"enabled"`
	Password    string `json:"password"`
	OldPassword string `json:"oldPassword"`
}

func (s *Server) handleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req UpdatePasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败", http.StatusBadRequest)
		return
	}

	var enabledStr string
	_ = s.db.QueryRow("SELECT value FROM settings WHERE key = 'is_password_enabled'").Scan(&enabledStr)
	
	if enabledStr == "true" {
		var oldHash string
		err := s.db.QueryRow("SELECT value FROM settings WHERE key = 'password_hash'").Scan(&oldHash)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(req.OldPassword)) != nil {
			BroadcastLog("WARN", "修改密码失败：原密码验证错误")
			http.Error(w, "原密码错误", http.StatusUnauthorized)
			return
		}
	}

	if req.Enabled {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			BroadcastLog("ERROR", "更新密码失败：密码加密出错")
			http.Error(w, "密码处理失败", http.StatusInternalServerError)
			return
		}
		s.db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('password_hash', ?)", string(hash))
		s.db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('is_password_enabled', 'true')")
		_ = refreshJWTSecret(s.db)
		BroadcastLog("INFO", "管理员密码及访问保护已启用/更新")
	} else {
		s.db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('is_password_enabled', 'false')")
		BroadcastLog("WARN", "管理员访问保护已关闭 (系统处于无密码状态)")
	}

	w.WriteHeader(http.StatusOK)
}

type LoginReq struct {
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败", http.StatusBadRequest)
		return
	}

	var hash string
	err := s.db.QueryRow("SELECT value FROM settings WHERE key = 'password_hash'").Scan(&hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		BroadcastLog("WARN", "管理员登录失败：密码错误")
		http.Error(w, "密码错误", http.StatusUnauthorized)
		return
	}
	
	BroadcastLog("INFO", "管理员登录成功")

	secret, err := getJWTSecret(s.db)
	if err != nil {
		http.Error(w, "系统错误", http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, "Token 生成失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func (s *Server) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	var enabled string
	err := s.db.QueryRow("SELECT value FROM settings WHERE key = 'is_password_enabled'").Scan(&enabled)
	
	isEnabled := false
	if err == nil && enabled == "true" {
		isEnabled = true
	}

	var zhihuCookies string
	err = s.db.QueryRow("SELECT value FROM settings WHERE key = 'zhihu_cookies'").Scan(&zhihuCookies)
	zhihuConfigured := false
	if err == nil && zhihuCookies != "" && zhihuCookies != "[]" {
		zhihuConfigured = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"enabled": isEnabled,
		"zhihu_configured": zhihuConfigured,
	})
}

func (s *Server) handleGetApiKey(w http.ResponseWriter, r *http.Request) {
	apiKey := s.getSetting("api_key")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"api_key": apiKey,
	})
}

func (s *Server) handleGenerateApiKey(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	rand.Read(b)
	newKey := fmt.Sprintf("zk_%x", b)
	s.setSetting("api_key", newKey)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"api_key": newKey,
	})
}