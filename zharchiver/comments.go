package main

import (
	"encoding/json"
	"net/http"
)

type Comment struct {
	ID        int    `json:"id"`
	AnswerID  string `json:"answer_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type AddCommentReq struct {
	Content string `json:"content"`
}

func (s *Server) handleGetComments(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "缺少 id 参数", http.StatusBadRequest)
		return
	}

	rows, err := s.db.Query(`
		SELECT id, answer_id, content, created_at
		FROM comments
		WHERE answer_id = ?
		ORDER BY created_at ASC
	`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.AnswerID, &c.Content, &c.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		list = append(list, c)
	}
	if list == nil {
		list = []Comment{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (s *Server) handleAddComment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "缺少 id 参数", http.StatusBadRequest)
		return
	}

	var req AddCommentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		http.Error(w, "无效的评论内容", http.StatusBadRequest)
		return
	}

	_, err := s.db.Exec(`
		INSERT INTO comments (answer_id, content) VALUES (?, ?)
	`, id, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "评论添加成功"})
}
