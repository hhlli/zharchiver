package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"zharchiver/models"
	"zharchiver/services"
	"zharchiver/utils"
)

type ArchiveRequest struct {
	URL string `json:"url"`
	Tag string `json:"tag"`
}

func (env *HandlerEnv) handleGetAnswers(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 50
	tag := r.URL.Query().Get("tag")
	search := r.URL.Query().Get("search")

	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			page = v
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}

	result, err := models.GetAnswersPaginated(env.db, page, limit, tag, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (env *HandlerEnv) handleGetGroupAnswers(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	questionID := r.URL.Query().Get("question_id")

	if title == "" && questionID == "" {
		http.Error(w, "必须提供 title 或 question_id", http.StatusBadRequest)
		return
	}

	answers, err := models.GetGroupAnswers(env.db, title, questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answers)
}

func (env *HandlerEnv) handleGetAnswerByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "缺少 id 参数", http.StatusBadRequest)
		return
	}

	data, err := models.GetAnswerByID(env.db, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "未找到该回答", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (env *HandlerEnv) handleUpdateTag(w http.ResponseWriter, r *http.Request) {
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

	if err := models.UpdateTag(env.db, id, req.Tag, req.Color); err != nil {
		utils.BroadcastLog("ERROR", fmt.Sprintf("更新标签失败 (ID: %s)：%v", id, err))
		http.Error(w, fmt.Sprintf("更新标签失败: %v", err), http.StatusInternalServerError)
		return
	}

	utils.BroadcastLog("INFO", fmt.Sprintf("更新了归档 (ID: %s) 的标签为: %s", id, req.Tag))
	w.WriteHeader(http.StatusOK)
}

func (env *HandlerEnv) handleRenameGlobalTag(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OldTag string `json:"old_tag"`
		NewTag string `json:"new_tag"`
		Color  string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "解析请求体失败", http.StatusBadRequest)
		return
	}
	if req.OldTag == "" || req.NewTag == "" {
		http.Error(w, "标签名称不能为空", http.StatusBadRequest)
		return
	}
	if req.Color == "" {
		req.Color = "blue"
	}

	if err := models.RenameGlobalTag(env.db, req.OldTag, req.NewTag, req.Color); err != nil {
		utils.BroadcastLog("ERROR", fmt.Sprintf("重命名全局标签失败：%v", err))
		http.Error(w, fmt.Sprintf("重命名标签失败: %v", err), http.StatusInternalServerError)
		return
	}

	utils.BroadcastLog("INFO", fmt.Sprintf("成功将标签 '%s' 重命名为 '%s'", req.OldTag, req.NewTag))
	w.WriteHeader(http.StatusOK)
}

func (env *HandlerEnv) handleDeleteGlobalTag(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Tag string `json:"tag"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "解析请求体失败", http.StatusBadRequest)
		return
	}
	if req.Tag == "" {
		http.Error(w, "标签名称不能为空", http.StatusBadRequest)
		return
	}

	if err := models.DeleteGlobalTag(env.db, req.Tag); err != nil {
		utils.BroadcastLog("ERROR", fmt.Sprintf("删除全局标签失败：%v", err))
		http.Error(w, fmt.Sprintf("删除标签失败: %v", err), http.StatusInternalServerError)
		return
	}

	utils.BroadcastLog("INFO", fmt.Sprintf("成功删除了标签 '%s'", req.Tag))
	w.WriteHeader(http.StatusOK)
}

func (env *HandlerEnv) handleGetAllTags(w http.ResponseWriter, r *http.Request) {
	type TagItem struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	rows, err := env.db.Query("SELECT DISTINCT tag, tag_color FROM answers WHERE tag != '' ORDER BY tag")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var tags []TagItem
	for rows.Next() {
		var t TagItem
		rows.Scan(&t.Name, &t.Color)
		tags = append(tags, t)
	}
	if tags == nil {
		tags = []TagItem{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func (env *HandlerEnv) handleUpdateAnswer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "缺少 id 参数", http.StatusBadRequest)
		return
	}

	var req struct {
		Title       string `json:"title"`
		ContentHTML string `json:"content_html"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "解析请求体失败", http.StatusBadRequest)
		return
	}

	if err := models.UpdateAnswerContent(env.db, id, req.Title, req.ContentHTML); err != nil {
		utils.BroadcastLog("ERROR", fmt.Sprintf("更新文章内容失败 (ID: %s)：%v", id, err))
		http.Error(w, fmt.Sprintf("更新文章内容失败: %v", err), http.StatusInternalServerError)
		return
	}

	utils.BroadcastLog("INFO", fmt.Sprintf("修改了归档内容 (ID: %s)", id))
	w.WriteHeader(http.StatusOK)
}

func (env *HandlerEnv) handleDeleteAnswer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "缺少 id 参数", http.StatusBadRequest)
		return
	}

	rowsAffected, err := models.DeleteAnswer(env.db, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("删除归档失败: %v", err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		utils.BroadcastLog("WARN", fmt.Sprintf("尝试删除归档 (ID: %s)，但未找到该记录", id))
		http.Error(w, "未找到该归档记录", http.StatusNotFound)
		return
	}

	utils.BroadcastLog("INFO", fmt.Sprintf("删除了归档及其评论 (ID: %s)", id))

	imgDir := filepath.Join("storage", "images", id)
	os.RemoveAll(imgDir)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	})
}

func (env *HandlerEnv) handleArchive(w http.ResponseWriter, r *http.Request) {
	var req ArchiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.URL) == "" {
		utils.BroadcastLog("ERROR", "收到无效的归档请求: 参数错误")
		http.Error(w, "请求参数错误", http.StatusBadRequest)
		return
	}

	go func(url, tag string) {
		utils.BroadcastLog("TASK_START", "开始后台归档任务")
		_, err := services.ProcessArchiveTask(env.db, url, tag)
		if err != nil {
			utils.BroadcastLog("ERROR", fmt.Sprintf("后台任务失败: %v", err))
			utils.BroadcastLog("TASK_FAILED", "归档失败")
		} else {
			utils.BroadcastLog("TASK_SUCCESS", "归档成功")
		}
	}(req.URL, req.Tag)

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    0,
		"message": "已加入后台解析任务",
	})
}
