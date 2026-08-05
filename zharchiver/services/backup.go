package services

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"zharchiver/models"
	"zharchiver/utils"
)

func CreateBackupZip(w io.Writer, db *sql.DB) error {
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

func StartTelegramAutoBackup(db *sql.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			token := models.GetSetting(db, "telegram_bot_token")
			chatID := models.GetSetting(db, "telegram_chat_id")
			backupEnabled := models.GetSetting(db, "telegram_backup_enabled")
			backupTime := models.GetSetting(db, "telegram_backup_time") // Format: HH:MM

			if backupEnabled != "true" || token == "" || chatID == "" || backupTime == "" {
				continue
			}

			now := time.Now()
			if fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()) == backupTime {
				lastBackup := models.GetSetting(db, "telegram_last_backup")
				todayStr := now.Format("2006-01-02")
				if lastBackup == todayStr {
					continue
				}

				utils.BroadcastLog("INFO", "触发定时任务：开始生成自动备份...")
				var buf bytes.Buffer
				if err := CreateBackupZip(&buf, db); err == nil {
					errSend := SendDocumentToTelegramWithErr(token, chatID, buf.Bytes(), "zharchiver_backup_"+todayStr+".zip")
					if errSend == nil {
						models.SetSetting(db, "telegram_last_backup", todayStr)
						utils.BroadcastLog("INFO", "自动备份已成功发送至 Telegram")
					} else {
						utils.BroadcastLog("ERROR", "自动备份发送失败："+errSend.Error())
					}
				} else {
					utils.BroadcastLog("ERROR", "自动备份生成失败："+err.Error())
				}
			}
		}
	}()
}

func SendDocumentToTelegramWithErr(token, chatID string, data []byte, filename string) error {
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
