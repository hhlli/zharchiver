package models

import (
	"database/sql"
)

func GetSetting(db *sql.DB, key string) string {
	var value string
	err := db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return ""
	}
	return value
}

func SetSetting(db *sql.DB, key string, value string) error {
	query := `
		INSERT INTO settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`
	_, err := db.Exec(query, key, value)
	return err
}

func GetTelegramAPIEndpoint(db *sql.DB) string {
	ep := GetSetting(db, "telegram_api_endpoint")
	if ep == "" {
		return "https://api.telegram.org"
	}
	// 移除末尾可能的斜杠
	if ep[len(ep)-1] == '/' {
		ep = ep[:len(ep)-1]
	}
	return ep
}
