package main

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (s *Server) getSetting(key string) string {
	var val string
	err := s.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&val)
	if err != nil {
		return ""
	}
	return val
}

func (s *Server) setSetting(key, value string) error {
	_, err := s.db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)", key, value)
	return err
}

func createBackupZip(w io.Writer, db *sql.DB) error {
	// 在打包之前强制执行 WAL Checkpoint，将缓存写入主文件，保证备份的完整性
	_, _ = db.Exec("PRAGMA wal_checkpoint(TRUNCATE)")

	zw := zip.NewWriter(w)
	defer zw.Close()

	if err := addFileToZip(zw, "db/zharchiver.db", "zharchiver.db"); err != nil {
		return err
	}

	err := filepath.Walk("storage/images", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		return addFileToZip(zw, path, path)
	})
	if err != nil {
		return err
	}
	return nil
}

func addFileToZip(zw *zip.Writer, srcPath, zipPath string) error {
	fileToZip, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(zipPath)
	header.Method = zip.Deflate

	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func (s *Server) handleDownloadBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=zharchiver_backup.zip")
	err := createBackupZip(w, s.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleRestoreBackup(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("backup")
	if err != nil {
		http.Error(w, "无效的文件上传: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	zr, err := zip.NewReader(file, handler.Size)
	if err != nil {
		http.Error(w, "无法读取 ZIP 文件: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 1. 关闭现有数据库连接
	s.db.Close()

	// 2. 清理旧的图片文件夹
	os.RemoveAll("storage/images")

	// 3. 解压并覆盖文件
	for _, f := range zr.File {
		path := filepath.FromSlash(f.Name)
		// 安全检查，防止路径穿越
		if path != "zharchiver.db" && !strings.HasPrefix(path, filepath.Join("storage", "images")) {
			continue
		}

		if path == "zharchiver.db" {
			path = filepath.Join("db", "zharchiver.db")
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			continue
		}
		srcFile, err := f.Open()
		if err != nil {
			dstFile.Close()
			continue
		}
		io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()
	}

	// 4. 重新初始化数据库连接
	db, err := initDB("db/zharchiver.db")
	if err != nil {
		http.Error(w, "数据库重启失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	s.db = db

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) handleGetBackupSettings(w http.ResponseWriter, r *http.Request) {
	settings := map[string]string{
		"telegram_backup_enabled": s.getSetting("telegram_backup_enabled"),
		"telegram_backup_time":    s.getSetting("telegram_backup_time"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (s *Server) handleSaveBackupSettings(w http.ResponseWriter, r *http.Request) {
	var settings map[string]string
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	s.setSetting("telegram_backup_enabled", settings["telegram_backup_enabled"])
	s.setSetting("telegram_backup_time", settings["telegram_backup_time"])
	
	BroadcastLog("INFO", "已更新数据备份配置")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) handleGetAISettings(w http.ResponseWriter, r *http.Request) {
	settings := map[string]string{
		"ai_base_url":   s.getSetting("ai_base_url"),
		"ai_api_key":    s.getSetting("ai_api_key"),
		"ai_model_name":      s.getSetting("ai_model_name"),
		"telegram_bot_token": s.getSetting("telegram_bot_token"),
		"telegram_chat_id":   s.getSetting("telegram_chat_id"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (s *Server) handleSaveAISettings(w http.ResponseWriter, r *http.Request) {
	var settings map[string]string
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	s.setSetting("ai_base_url", settings["ai_base_url"])
	s.setSetting("ai_api_key", settings["ai_api_key"])
	s.setSetting("ai_model_name", settings["ai_model_name"])
	s.setSetting("telegram_bot_token", settings["telegram_bot_token"])
	s.setSetting("telegram_chat_id", settings["telegram_chat_id"])
	
	BroadcastLog("INFO", "已更新工具(AI与Telegram)配置")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) handleTestAIConnection(w http.ResponseWriter, r *http.Request) {
	var settings map[string]string
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	
	baseURL := settings["ai_base_url"]
	apiKey := settings["ai_api_key"]
	modelName := settings["ai_model_name"]
	
	if baseURL == "" || apiKey == "" || modelName == "" {
		BroadcastLog("WARN", "测试 AI 连通性失败：配置不完整")
		http.Error(w, "请填写完整配置", http.StatusBadRequest)
		return
	}

	BroadcastLog("INFO", fmt.Sprintf("开始测试 AI 连通性 (模型: %s)...", modelName))

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
		BroadcastLog("ERROR", "测试 AI 连通性失败：网络请求出错 ("+err.Error()+")")
		http.Error(w, "网络请求失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		BroadcastLog("ERROR", fmt.Sprintf("测试 AI 连通性失败：状态码 %d", resp.StatusCode))
		http.Error(w, fmt.Sprintf("API 返回错误状态码 %d: %s", resp.StatusCode, string(respBody)), http.StatusBadRequest)
		return
	}
	
	BroadcastLog("INFO", "AI 连通性测试成功")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) handleSendTelegramBackup(w http.ResponseWriter, r *http.Request) {
	token := s.getSetting("telegram_bot_token")
	chatID := s.getSetting("telegram_chat_id")
	
	if token == "" || chatID == "" {
		BroadcastLog("WARN", "手动发送 Telegram 备份失败：未配置 Bot Token 或 Chat ID")
		http.Error(w, "未配置 Telegram Bot Token 或 Chat ID", http.StatusBadRequest)
		return
	}

	BroadcastLog("INFO", "开始手动打包并发送 Telegram 备份...")

	var buf bytes.Buffer
	if err := createBackupZip(&buf, s.db); err != nil {
		BroadcastLog("ERROR", "手动发送 Telegram 备份失败：创建备份出错 ("+err.Error()+")")
		http.Error(w, "创建备份失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	todayStr := time.Now().Format("2006-01-02_150405")
	err := sendDocumentToTelegramWithErr(token, chatID, buf.Bytes(), "zharchiver_manual_backup_"+todayStr+".zip")
	if err != nil {
		BroadcastLog("ERROR", "手动发送 Telegram 备份失败："+err.Error())
		http.Error(w, "发送到 Telegram 失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	BroadcastLog("INFO", "手动 Telegram 备份发送成功！")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func startTelegramAutoBackup(s *Server) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			token := s.getSetting("telegram_bot_token")
			chatID := s.getSetting("telegram_chat_id")
			backupEnabled := s.getSetting("telegram_backup_enabled")
			backupTime := s.getSetting("telegram_backup_time") // Format: HH:MM

			if backupEnabled != "true" || token == "" || chatID == "" || backupTime == "" {
				continue
			}

			now := time.Now()
			if fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()) == backupTime {
				lastBackup := s.getSetting("telegram_last_backup")
				todayStr := now.Format("2006-01-02")
				if lastBackup == todayStr {
					continue
				}

				BroadcastLog("INFO", "触发定时任务：开始生成自动备份...")
				var buf bytes.Buffer
				if err := createBackupZip(&buf, s.db); err == nil {
					errSend := sendDocumentToTelegramWithErr(token, chatID, buf.Bytes(), "zharchiver_backup_"+todayStr+".zip")
					if errSend == nil {
						s.setSetting("telegram_last_backup", todayStr)
						BroadcastLog("INFO", "自动备份已成功发送至 Telegram")
					} else {
						BroadcastLog("ERROR", "自动备份发送失败："+errSend.Error())
					}
				} else {
					BroadcastLog("ERROR", "自动备份生成失败："+err.Error())
				}
			}
		}
	}()
}

func sendDocumentToTelegram(token, chatID string, data []byte, filename string) {
	_ = sendDocumentToTelegramWithErr(token, chatID, data, filename)
}

func sendDocumentToTelegramWithErr(token, chatID string, data []byte, filename string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", token)
	
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	w.WriteField("chat_id", chatID)
	fw, err := w.CreateFormFile("document", filename)
	if err != nil {
		return err
	}
	fw.Write(data)
	w.Close()
	
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
