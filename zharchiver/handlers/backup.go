package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"os"
	"path/filepath"
	"strings"
	"archive/zip"
	"io"

	"zharchiver/models"
	"zharchiver/services"
	"zharchiver/utils"
)

func (env *HandlerEnv) handleDownloadBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=zharchiver_backup.zip")
	err := services.CreateBackupZip(w, env.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (env *HandlerEnv) handleRestoreBackup(w http.ResponseWriter, r *http.Request) {
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
	env.db.Close()

	// 2. 清理旧的图片文件夹
	os.RemoveAll("storage/images")

	// 3. 解压并覆盖文件
	for _, f := range zr.File {
		path := filepath.FromSlash(f.Name)
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
	db, err := models.InitDB("db/zharchiver.db")
	if err != nil {
		http.Error(w, "数据库重启失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	env.db = db

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (env *HandlerEnv) handleSendTelegramBackup(w http.ResponseWriter, r *http.Request) {
	token := models.GetSetting(env.db, "telegram_bot_token")
	chatID := models.GetSetting(env.db, "telegram_chat_id")
	
	if token == "" || chatID == "" {
		utils.BroadcastLog("WARN", "手动发送 Telegram 备份失败：未配置 Bot Token 或 Chat ID")
		http.Error(w, "未配置 Telegram Bot Token 或 Chat ID", http.StatusBadRequest)
		return
	}

	utils.BroadcastLog("INFO", "开始手动打包并发送 Telegram 备份...")

	var buf bytes.Buffer
	if err := services.CreateBackupZip(&buf, env.db); err != nil {
		utils.BroadcastLog("ERROR", "手动发送 Telegram 备份失败：创建备份出错 ("+err.Error()+")")
		http.Error(w, "创建备份失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	todayStr := time.Now().Format("2006-01-02_150405")
	err := services.SendDocumentToTelegramWithErr(token, chatID, buf.Bytes(), "zharchiver_manual_backup_"+todayStr+".zip")
	if err != nil {
		utils.BroadcastLog("ERROR", "手动发送 Telegram 备份失败："+err.Error())
		http.Error(w, "发送到 Telegram 失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.BroadcastLog("INFO", "手动 Telegram 备份发送成功！")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (env *HandlerEnv) handleGetBackupSettings(w http.ResponseWriter, r *http.Request) {
	settings := map[string]string{
		"telegram_backup_enabled": models.GetSetting(env.db, "telegram_backup_enabled"),
		"telegram_backup_time":    models.GetSetting(env.db, "telegram_backup_time"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (env *HandlerEnv) handleSaveBackupSettings(w http.ResponseWriter, r *http.Request) {
	var settings map[string]string
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	models.SetSetting(env.db, "telegram_backup_enabled", settings["telegram_backup_enabled"])
	models.SetSetting(env.db, "telegram_backup_time", settings["telegram_backup_time"])
	
	utils.BroadcastLog("INFO", "已更新数据备份配置")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
