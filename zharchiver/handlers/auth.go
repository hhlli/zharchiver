package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"

	"zharchiver/models"
	"zharchiver/utils"
)

func getJWTSecret(db *sql.DB) ([]byte, error) {
	secretHex := models.GetSetting(db, "jwt_secret")
	if secretHex == "" {
		b := make([]byte, 32)
		rand.Read(b)
		secretHex = hex.EncodeToString(b)
		models.SetSetting(db, "jwt_secret", secretHex)
	}
	return hex.DecodeString(secretHex)
}

func refreshJWTSecret(db *sql.DB) error {
	b := make([]byte, 32)
	rand.Read(b)
	secretHex := hex.EncodeToString(b)
	return models.SetSetting(db, "jwt_secret", secretHex)
}

type UpdatePasswordReq struct {
	Enabled     bool   `json:"enabled"`
	Password    string `json:"password"`
	OldPassword string `json:"oldPassword"`
}

func (env *HandlerEnv) handleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req UpdatePasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败", http.StatusBadRequest)
		return
	}

	enabledStr := models.GetSetting(env.db, "is_password_enabled")
	if enabledStr == "true" {
		oldHash := models.GetSetting(env.db, "password_hash")
		if bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(req.OldPassword)) != nil {
			utils.BroadcastLog("WARN", "修改密码失败：原密码验证错误")
			http.Error(w, "原密码错误", http.StatusUnauthorized)
			return
		}
	}

	if req.Enabled {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.BroadcastLog("ERROR", "更新密码失败：密码加密出错")
			http.Error(w, "密码处理失败", http.StatusInternalServerError)
			return
		}
		models.SetSetting(env.db, "password_hash", string(hash))
		models.SetSetting(env.db, "is_password_enabled", "true")
		_ = refreshJWTSecret(env.db)
		utils.BroadcastLog("INFO", "管理员密码及访问保护已启用/更新")
	} else {
		models.SetSetting(env.db, "is_password_enabled", "false")
		utils.BroadcastLog("WARN", "管理员访问保护已关闭 (系统处于无密码状态)")
	}

	w.WriteHeader(http.StatusOK)
}

type LoginReq struct {
	Password string `json:"password"`
}

func (env *HandlerEnv) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败", http.StatusBadRequest)
		return
	}

	hash := models.GetSetting(env.db, "password_hash")
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		utils.BroadcastLog("WARN", "管理员登录失败：密码错误")
		http.Error(w, "密码错误", http.StatusUnauthorized)
		return
	}
	
	utils.BroadcastLog("INFO", "管理员登录成功")

	secret, err := getJWTSecret(env.db)
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

func (env *HandlerEnv) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	enabled := models.GetSetting(env.db, "is_password_enabled")
	isEnabled := false
	if enabled == "true" {
		isEnabled = true
	}

	zhihuCookies := models.GetSetting(env.db, "zhihu_cookies")
	zhihuConfigured := false
	if zhihuCookies != "" && zhihuCookies != "[]" {
		zhihuConfigured = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"enabled": isEnabled,
		"zhihu_configured": zhihuConfigured,
	})
}

func (env *HandlerEnv) handleGetApiKey(w http.ResponseWriter, r *http.Request) {
	apiKey := models.GetSetting(env.db, "api_key")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"api_key": apiKey,
	})
}

func (env *HandlerEnv) handleGenerateApiKey(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	rand.Read(b)
	newKey := fmt.Sprintf("zk_%x", b)
	models.SetSetting(env.db, "api_key", newKey)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"api_key": newKey,
	})
}

type WsQRMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Image   string `json:"image,omitempty"`
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (env *HandlerEnv) WsZhihuQRHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	sendMsg := func(status, msg string, image ...string) {
		m := WsQRMessage{Status: status, Message: msg}
		if len(image) > 0 {
			m.Image = image[0]
		}
		_ = conn.WriteJSON(m)
	}

	sendMsg("loading", "正在启动安全环境...")

	l := launcher.New().Headless(true)
	l.Set("disable-blink-features", "AutomationControlled")
	l.Set("disable-web-security")
	u, err := l.Launch()
	if err != nil {
		sendMsg("error", "启动浏览器失败: "+err.Error())
		return
	}

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	clientDisconnected := make(chan bool)
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				clientDisconnected <- true
				return
			}
		}
	}()

	page := browser.MustPage()
	sendMsg("loading", "正在加载知乎登录页...")

	go func() {
		_ = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		})

		err := page.Navigate("https://www.zhihu.com/signin?next=%2F")
		if err != nil {
			sendMsg("error", "页面加载超时或网络失败")
			return
		}

		sendMsg("loading", "正在获取二维码...")

		var qrSrc string
		for i := 0; i < 30; i++ {
			time.Sleep(1 * time.Second)
			canvas, err := page.Element("canvas")
			if err == nil {
				res, err := canvas.Eval(`() => this.toDataURL()`)
				if err == nil && res != nil {
					val := res.Value.String()
					if len(val) > 2 && val[0] == '"' {
						qrSrc = val[1 : len(val)-1]
					} else {
						qrSrc = val
					}
					if len(qrSrc) > 20 && qrSrc[:4] == "data" {
						break
					}
				}
			}
		}

		if qrSrc == "" {
			sendMsg("error", "未能抓取到二维码，可能被知乎风控")
			return
		}

		sendMsg("qrcode", "请打开知乎 App 扫码", qrSrc)

		for i := 0; i < 60; i++ {
			time.Sleep(1 * time.Second)
			url := page.MustInfo().URL
			if url == "https://www.zhihu.com/" || url == "https://www.zhihu.com" {
				sendMsg("success", "扫码成功，正在保存凭证...")
				
				cookies, err := browser.GetCookies()
				if err == nil {
					cookieJSON, _ := json.Marshal(cookies)
					models.SetSetting(env.db, "zhihu_cookies", string(cookieJSON))
					sendMsg("done", "凭证保存完成")
				} else {
					sendMsg("error", "获取 Cookie 失败")
				}
				return
			}
		}
		sendMsg("error", "登录超时，请重试")
	}()

	<-clientDisconnected
}
