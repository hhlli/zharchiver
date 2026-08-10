package models

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// InitDB 初始化数据库连接，并创建相关表结构
func InitDB(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %v", err)
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	
	// 启用 WAL 模式提高并发读写性能
	_, _ = db.Exec("PRAGMA journal_mode=WAL")
	// 限制单个 SQLite 文件的最大打开连接数为 1，防止 database is locked 错误
	db.SetMaxOpenConns(1)

	schema := `
	CREATE TABLE IF NOT EXISTS answers (
		answer_id TEXT PRIMARY KEY,
		question_id TEXT NOT NULL,
		title TEXT NOT NULL,
		author_name TEXT,
		author_avatar TEXT,
		content_html TEXT,
		created_time INTEGER,
		updated_time INTEGER,
		saved_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	// 数据库迁移: 增加新字段 (忽略错误，因为如果列已存在会报错)
	_, _ = db.Exec(`ALTER TABLE answers ADD COLUMN tag TEXT DEFAULT ''`)
	_, _ = db.Exec(`ALTER TABLE answers ADD COLUMN tag_color TEXT DEFAULT 'blue'`)
	_, _ = db.Exec(`ALTER TABLE answers ADD COLUMN author_avatar TEXT DEFAULT ''`)
	_, _ = db.Exec(`ALTER TABLE answers ADD COLUMN is_favorite INTEGER DEFAULT 0`)

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS settings (
            key TEXT PRIMARY KEY,
            value TEXT NOT NULL
        );
    `)
	if err != nil {
		return nil, fmt.Errorf("创建 settings 表失败: %v", err)
	}

	// 为提高性能，创建相关索引（如果不存在）
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_answers_title_qid ON answers (title, question_id);`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_answers_saved_at ON answers (saved_at);`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_answers_tag ON answers (tag);`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			answer_id TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("创建 comments 表失败: %v", err)
	}

	return db, nil
}
