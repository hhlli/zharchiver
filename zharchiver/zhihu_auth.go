package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type WsQRMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Image   string `json:"image,omitempty"`
}

func (s *Server) WsZhihuQRHandler(w http.ResponseWriter, r *http.Request) {
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
	l.Set("disable-web-security") // 解决知乎跨域图片污染 Canvas 导致 toDataURL 失败的问题
	u, err := l.Launch()
	if err != nil {
		sendMsg("error", "启动浏览器失败: "+err.Error())
		return
	}

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose() // 保证连接断开时销毁进程

	// 监听前端断开
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

		// 等待 canvas 出现 (知乎现在使用 canvas 渲染二维码)
		var qrSrc string
		for i := 0; i < 30; i++ {
			time.Sleep(1 * time.Second)
			canvas, err := page.Element("canvas")
			if err == nil {
				res, err := canvas.Eval(`() => this.toDataURL()`)
				if err == nil && res != nil {
					val := res.Value.String()
					// 去掉首尾的双引号
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

		// 轮询检查是否跳转到了首页
		for i := 0; i < 60; i++ {
			time.Sleep(1 * time.Second)
			url := page.MustInfo().URL
			if url == "https://www.zhihu.com/" || url == "https://www.zhihu.com" {
				sendMsg("success", "扫码成功，正在保存凭证...")
				
				// 提取 Cookie
				cookies, err := browser.GetCookies()
				if err == nil {
					cookieJSON, _ := json.Marshal(cookies)
					s.setSetting("zhihu_cookies", string(cookieJSON))
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
