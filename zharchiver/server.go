package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	db *sql.DB
}

type ArchiveRequest struct {
	URL string `json:"url"`
	Tag string `json:"tag"`
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

func startServer(db *sql.DB, port int) {
	s := &Server{db: db}

	startTelegramAutoBackup(s)
	startTelegramBotListener(s)

	mux := http.NewServeMux()

	// 静态文件服务
	mux.Handle("/", http.FileServer(http.Dir("dist")))
	mux.Handle("/storage/", http.StripPrefix("/storage/", http.FileServer(http.Dir("storage"))))

	mux.HandleFunc("GET /api/answers", s.handleGetAnswers)
	mux.HandleFunc("GET /api/answers/{id}", s.handleGetAnswerByID)
	mux.HandleFunc("DELETE /api/answers/{id}", s.handleDeleteAnswer)
	mux.HandleFunc("PATCH /api/answers/{id}/tag", s.handleUpdateTag)
	mux.HandleFunc("GET /api/answers/{id}/comments", s.handleGetComments)
	mux.HandleFunc("POST /api/answers/{id}/comments", s.handleAddComment)
	mux.HandleFunc("POST /api/archive", s.handleArchive)
	mux.HandleFunc("GET /api/logs/ws", WsLogHandler)
	mux.HandleFunc("GET /api/auth/status", s.handleAuthStatus)
    mux.HandleFunc("POST /api/auth/login", s.handleLogin)
    mux.HandleFunc("POST /api/auth/update", s.handleUpdatePassword)
	mux.HandleFunc("GET /api/auth/zhihu/qrcode/ws", s.WsZhihuQRHandler)
	
	mux.HandleFunc("GET /api/settings/api_key", s.handleGetApiKey)
	mux.HandleFunc("POST /api/settings/api_key/generate", s.handleGenerateApiKey)

	mux.HandleFunc("GET /api/backup/download", s.handleDownloadBackup)
	mux.HandleFunc("POST /api/backup/restore", s.handleRestoreBackup)
	mux.HandleFunc("POST /api/backup/telegram/send", s.handleSendTelegramBackup)
	mux.HandleFunc("GET /api/settings/backup", s.handleGetBackupSettings)
	mux.HandleFunc("POST /api/settings/backup", s.handleSaveBackupSettings)

	mux.HandleFunc("GET /api/settings/ai", s.handleGetAISettings)
	mux.HandleFunc("POST /api/settings/ai", s.handleSaveAISettings)
	mux.HandleFunc("POST /api/settings/ai/test", s.handleTestAIConnection)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("服务已启动: http://localhost%s\n", addr)

	handler := corsMiddleware(authMiddleware(db, mux))

if err := http.ListenAndServe(addr, handler); err != nil {
    log.Fatalf("服务启动失败: %v", err)
}
}

func (s *Server) handleGetAnswers(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(`
		SELECT answer_id, question_id, title, author_name, created_time, updated_time, saved_at, tag, tag_color
		FROM answers
		ORDER BY saved_at DESC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type AnswerSummary struct {
		AnswerID    string `json:"answer_id"`
		QuestionID  string `json:"question_id"`
		Title       string `json:"title"`
		AuthorName  string `json:"author_name"`
		CreatedTime int64  `json:"created_time"`
		UpdatedTime int64  `json:"updated_time"`
		SavedAt     string `json:"saved_at"`
		Tag         string `json:"tag"`
		TagColor    string `json:"tag_color"`
	}

	var list []AnswerSummary
	for rows.Next() {
		var item AnswerSummary
		if err := rows.Scan(&item.AnswerID, &item.QuestionID, &item.Title, &item.AuthorName, &item.CreatedTime, &item.UpdatedTime, &item.SavedAt, &item.Tag, &item.TagColor); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		list = append(list, item)
	}

	if list == nil {
		list = []AnswerSummary{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (s *Server) handleGetAnswerByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "缺少 id 参数", http.StatusBadRequest)
		return
	}

	var data AnswerData
	var savedAt string
	err := s.db.QueryRow(`
		SELECT answer_id, question_id, title, author_name, content_html, created_time, updated_time, saved_at, tag, tag_color
		FROM answers
		WHERE answer_id = ?
	`, id).Scan(&data.AnswerID, &data.QuestionID, &data.Title, &data.AuthorName, &data.ContentHTML, &data.CreatedTime, &data.UpdatedTime, &savedAt, &data.Tag, &data.TagColor)

	if err == sql.ErrNoRows {
		http.Error(w, "未找到该回答", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *Server) handleUpdateTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "缺少 id 参数", http.StatusBadRequest)
		return
	}

	var req struct {
		Tag   string `json:"tag"`
		Color string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "解析请求体失败", http.StatusBadRequest)
		return
	}
	if req.Color == "" {
		req.Color = "blue"
	}

	_, err := s.db.Exec("UPDATE answers SET tag = ?, tag_color = ? WHERE answer_id = ?", req.Tag, req.Color, id)
	if err != nil {
		BroadcastLog("ERROR", fmt.Sprintf("更新标签失败 (ID: %s)：%v", id, err))
		http.Error(w, fmt.Sprintf("更新标签失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 同步更新所有拥有该标签的记录颜色
	if req.Tag != "" {
		_, _ = s.db.Exec("UPDATE answers SET tag_color = ? WHERE tag = ?", req.Color, req.Tag)
	}

	BroadcastLog("INFO", fmt.Sprintf("更新了归档 (ID: %s) 的标签为: %s", id, req.Tag))
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleDeleteAnswer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "缺少 id 参数", http.StatusBadRequest)
		return
	}

	_, err := s.db.Exec("DELETE FROM comments WHERE answer_id = ?", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("删除关联评论失败: %v", err), http.StatusInternalServerError)
		return
	}

	result, err := s.db.Exec("DELETE FROM answers WHERE answer_id = ?", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("删除归档失败: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		BroadcastLog("WARN", fmt.Sprintf("尝试删除归档 (ID: %s)，但未找到该记录", id))
		http.Error(w, "未找到该归档记录", http.StatusNotFound)
		return
	}

	BroadcastLog("INFO", fmt.Sprintf("删除了归档及其评论 (ID: %s)", id))

	imgDir := filepath.Join("storage", "images", id)
	os.RemoveAll(imgDir)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	})
}

func (s *Server) processArchiveTask(url string, tag string) (*AnswerData, error) {
	BroadcastLog("INFO", "=== 开始新的归档任务 ===")
	BroadcastLog("INFO", fmt.Sprintf("目标链接: %s", url))
	BroadcastLog("INFO", fmt.Sprintf("标签: %s", tag))

    target, err := parseZhihuLink(url)
    if err != nil {
        BroadcastLog("ERROR", fmt.Sprintf("链接解析失败: %v", err))
        return nil, err
    }
    BroadcastLog("INFO", fmt.Sprintf("链接解析成功，识别到回答 ID: %s", target.AnswerID))

    BroadcastLog("INFO", "启动浏览器引擎(Rod)抓取页面数据...")
    cookiesStr := s.getSetting("zhihu_cookies")
    jsonData, err := fetchWithRod(target.CleanURL, cookiesStr)
    if err != nil {
        BroadcastLog("ERROR", fmt.Sprintf("页面抓取失败: %v", err))
        return nil, fmt.Errorf("抓取失败: %v", err)
    }

    BroadcastLog("INFO", "页面抓取完成，开始提取核心 JSON 数据...")
    data, err := parseInitialJSON(jsonData, target)
    if err != nil {
        BroadcastLog("ERROR", fmt.Sprintf("数据提取失败: %v", err))
        return nil, fmt.Errorf("解析失败: %v", err)
    }
    BroadcastLog("INFO", fmt.Sprintf("数据提取成功，标题: [%s]，作者: [%s]", data.Title, data.AuthorName))

    if len(data.ImageURLs) > 0 {
        BroadcastLog("INFO", fmt.Sprintf("发现 %d 张相关图片，启动本地化下载流程...", len(data.ImageURLs)))
        _ = processImages(data)
        BroadcastLog("INFO", "图片本地化处理完成")
    } else {
        BroadcastLog("INFO", "当前回答无内嵌图片，跳过下载流程")
    }

    BroadcastLog("INFO", "准备将归档数据写入 SQLite 数据库...")
   	data.Tag = tag

	// 保存到数据库
	if err := saveAnswer(s.db, data); err != nil {
        BroadcastLog("ERROR", fmt.Sprintf("数据库写入失败: %v", err))
        return nil, fmt.Errorf("保存失败: %v", err)
    }
    
    BroadcastLog("INFO", "数据库写入成功，归档流程全部完成")
    return data, nil
}

func (s *Server) handleArchive(w http.ResponseWriter, r *http.Request) {
    var req ArchiveRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.URL) == "" {
        BroadcastLog("ERROR", "收到无效的归档请求: 参数错误")
        http.Error(w, "请求参数错误", http.StatusBadRequest)
        return
    }

    data, err := s.processArchiveTask(req.URL, req.Tag)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code":    0,
        "message": "归档成功",
        "data":    data,
    })
}