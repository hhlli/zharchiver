package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// LogEntry 定义发送给前端的日志格式
type LogEntry struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

var (
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	wsClients   = make(map[*websocket.Conn]bool)
	logHistory  []LogEntry
	wsClientsMu sync.Mutex
)

const maxLogHistory = 500

// WsLogHandler 处理前端的 WebSocket 连接请求
func WsLogHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	wsClientsMu.Lock()
	wsClients[conn] = true
	// 发送历史日志给新连接的客户端
	history := make([]LogEntry, len(logHistory))
	copy(history, logHistory)
	wsClientsMu.Unlock()

	for _, entry := range history {
		msg, err := json.Marshal(entry)
		if err == nil {
			_ = conn.WriteMessage(websocket.TextMessage, msg)
		}
	}

	// 阻塞监听，处理前端主动断开连接
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			wsClientsMu.Lock()
			delete(wsClients, conn)
			wsClientsMu.Unlock()
			break
		}
	}
}

// BroadcastLog 接收日志内容并向所有已连接的客户端广播
func BroadcastLog(level, message string) {
	entry := LogEntry{
		Time:    time.Now().Format("15:04:05"),
		Level:   level,
		Message: message,
	}

	msg, err := json.Marshal(entry)
	if err != nil {
		return
	}

	wsClientsMu.Lock()
	logHistory = append(logHistory, entry)
	if len(logHistory) > maxLogHistory {
		logHistory = logHistory[len(logHistory)-maxLogHistory:]
	}
	
	// Copy the current clients into a local slice to release lock before I/O
	var clientsToNotify []*websocket.Conn
	for client := range wsClients {
		clientsToNotify = append(clientsToNotify, client)
	}
	wsClientsMu.Unlock()

	for _, client := range clientsToNotify {
		if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
			client.Close()
			wsClientsMu.Lock()
			delete(wsClients, client)
			wsClientsMu.Unlock()
		}
	}
}