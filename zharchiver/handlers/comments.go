package handlers

import (
	"encoding/json"
	"net/http"

	"zharchiver/models"
)

type AddCommentReq struct {
	Content string `json:"content"`
}

func (env *HandlerEnv) handleGetComments(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "缺少 id 参数", http.StatusBadRequest)
		return
	}

	list, err := models.GetComments(env.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (env *HandlerEnv) handleAddComment(w http.ResponseWriter, r *http.Request) {
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

	if err := models.AddComment(env.db, id, req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "评论添加成功"})
}

func (env *HandlerEnv) handleDeleteComment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	commentID := r.PathValue("comment_id")
	if id == "" || commentID == "" {
		http.Error(w, "缺少参数", http.StatusBadRequest)
		return
	}

	if err := models.DeleteComment(env.db, id, commentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "已删除"})
}

func (env *HandlerEnv) handleUpdateComment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	commentID := r.PathValue("comment_id")
	if id == "" || commentID == "" {
		http.Error(w, "缺少参数", http.StatusBadRequest)
		return
	}

	var req AddCommentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		http.Error(w, "无效的评论内容", http.StatusBadRequest)
		return
	}

	if err := models.UpdateComment(env.db, id, commentID, req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "已更新"})
}
