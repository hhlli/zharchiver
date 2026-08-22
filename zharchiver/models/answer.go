package models

import (
	"database/sql"
)

type AnswerData struct {
	QuestionID  string   `json:"question_id"`
	AnswerID    string   `json:"answer_id"`
	Title       string   `json:"title"`
	AuthorName  string   `json:"author_name"`
	ContentHTML string   `json:"content_html"`
	CreatedTime int64    `json:"created_time"`
	UpdatedTime int64    `json:"updated_time"`
	SavedAt     string   `json:"saved_at"`
	ImageURLs   []string `json:"image_urls"`
	Tag         string   `json:"tag"`
	TagColor    string   `json:"tag_color"`
	IsFavorite  int      `json:"is_favorite"`
}

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
	GroupCount  int    `json:"group_count"`
	IsFavorite  int    `json:"is_favorite"`
}

func SaveAnswer(db *sql.DB, data *AnswerData) error {
	if data.Tag != "" {
		var existingColor string
		err := db.QueryRow("SELECT tag_color FROM answers WHERE tag = ? AND tag_color != '' LIMIT 1", data.Tag).Scan(&existingColor)
		if err == nil && existingColor != "" {
			data.TagColor = existingColor
		}
	}
	if data.TagColor == "" {
		data.TagColor = "blue"
	}
	query := `
	INSERT INTO answers (answer_id, question_id, title, author_name, content_html, created_time, updated_time, tag, tag_color)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(answer_id) DO UPDATE SET
		title = excluded.title,
		author_name = excluded.author_name,
		content_html = excluded.content_html,
		updated_time = excluded.updated_time,
		tag = CASE WHEN excluded.tag != '' THEN excluded.tag ELSE answers.tag END,
		tag_color = CASE WHEN excluded.tag != '' THEN excluded.tag_color ELSE answers.tag_color END,
		saved_at = CURRENT_TIMESTAMP;
	`

	_, err := db.Exec(query,
		data.AnswerID,
		data.QuestionID,
		data.Title,
		data.AuthorName,
		data.ContentHTML,
		data.CreatedTime,
		data.UpdatedTime,
		data.Tag,
		data.TagColor,
	)

	return err
}

type AnswerListResult struct {
	Items []AnswerSummary `json:"items"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

func GetAnswersPaginated(db *sql.DB, page, limit int, tag, search string, isFavorite bool) (*AnswerListResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}
	offset := (page - 1) * limit

	// 构建 WHERE 条件
	where := "WHERE 1=1"
	args := []interface{}{}
	if tag != "" {
		where += " AND tag = ?"
		args = append(args, tag)
	}
	if search != "" {
		where += " AND (title LIKE ? OR author_name LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	if isFavorite {
		where += " AND is_favorite = 1"
	}

	// 查总数
	var total int
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	db.QueryRow("SELECT COUNT(*) FROM answers "+where, countArgs...).Scan(&total)

	// 分页查询
	paginateArgs := append(args, limit, offset)
	rows, err := db.Query(`
		SELECT answer_id, question_id, title, author_name, created_time, updated_time, saved_at, tag, tag_color, is_favorite,
			(SELECT COUNT(*) FROM answers sub 
			 WHERE (sub.title = answers.title AND answers.title != '') 
			    OR (sub.question_id = answers.question_id AND answers.question_id != '0' AND answers.title == '')) as group_count
		FROM answers `+where+`
		ORDER BY saved_at DESC
		LIMIT ? OFFSET ?
	`, paginateArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AnswerSummary
	for rows.Next() {
		var item AnswerSummary
		if err := rows.Scan(&item.AnswerID, &item.QuestionID, &item.Title, &item.AuthorName, &item.CreatedTime, &item.UpdatedTime, &item.SavedAt, &item.Tag, &item.TagColor, &item.IsFavorite, &item.GroupCount); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if list == nil {
		list = []AnswerSummary{}
	}
	return &AnswerListResult{Items: list, Total: total, Page: page, Limit: limit}, nil
}

func GetAnswerByID(db *sql.DB, id string) (*AnswerData, error) {
	var data AnswerData
	err := db.QueryRow(`
		SELECT answer_id, question_id, title, author_name, content_html, created_time, updated_time, saved_at, tag, tag_color, is_favorite
		FROM answers
		WHERE answer_id = ?
	`, id).Scan(&data.AnswerID, &data.QuestionID, &data.Title, &data.AuthorName, &data.ContentHTML, &data.CreatedTime, &data.UpdatedTime, &data.SavedAt, &data.Tag, &data.TagColor, &data.IsFavorite)

	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetGroupAnswers 获取同一问题下的所有回答（用于渲染详情页的作者标签），按保存时间倒序
func GetGroupAnswers(db *sql.DB, title, questionID string) ([]AnswerSummary, error) {
	query := `
		SELECT answer_id, question_id, title, author_name, created_time, updated_time, saved_at, tag, tag_color, is_favorite
		FROM answers
		WHERE 1=1
	`
	var args []interface{}
	var condition string

	// 与前端分组逻辑保持一致：有非空 title 按 title，否则按 question_id
	if title != "" {
		condition = " AND title = ?"
		args = append(args, title)
		// 如果此时还有有效的 question_id 作为补充条件
		if questionID != "" && questionID != "0" {
			condition = " AND (title = ? OR question_id = ?)"
			args = []interface{}{title, questionID}
		}
	} else if questionID != "" && questionID != "0" {
		condition = " AND question_id = ?"
		args = append(args, questionID)
	} else {
		// 都为空，无法成组，返回空
		return []AnswerSummary{}, nil
	}

	query += condition + ` ORDER BY saved_at DESC`
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []AnswerSummary
	for rows.Next() {
		var s AnswerSummary
		if err := rows.Scan(
			&s.AnswerID, &s.QuestionID, &s.Title, &s.AuthorName,
			&s.CreatedTime, &s.UpdatedTime, &s.SavedAt,
			&s.Tag, &s.TagColor, &s.IsFavorite,
		); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	return results, nil
}

func UpdateTag(db *sql.DB, id string, tag string, color string) error {
	_, err := db.Exec("UPDATE answers SET tag = ?, tag_color = ? WHERE answer_id = ?", tag, color, id)
	if err != nil {
		return err
	}
	// 同步更新所有拥有该标签的记录颜色
	if tag != "" {
		_, _ = db.Exec("UPDATE answers SET tag_color = ? WHERE tag = ?", color, tag)
	}
	return nil
}

func RenameGlobalTag(db *sql.DB, oldTag string, newTag string, color string) error {
	_, err := db.Exec("UPDATE answers SET tag = ?, tag_color = ? WHERE tag = ?", newTag, color, oldTag)
	return err
}

func DeleteGlobalTag(db *sql.DB, tag string) error {
	_, err := db.Exec("UPDATE answers SET tag = '', tag_color = '' WHERE tag = ?", tag)
	return err
}

func UpdateAnswerContent(db *sql.DB, id string, title string, contentHTML string) error {
	_, err := db.Exec("UPDATE answers SET title = ?, content_html = ? WHERE answer_id = ?", title, contentHTML, id)
	return err
}

func DeleteAnswer(db *sql.DB, id string) (int64, error) {
	_, err := db.Exec("DELETE FROM comments WHERE answer_id = ?", id)
	if err != nil {
		return 0, err
	}

	result, err := db.Exec("DELETE FROM answers WHERE answer_id = ?", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func DeleteGroup(db *sql.DB, title string, questionID string) (int64, error) {
	var args []interface{}
	var condition string

	if title != "" {
		condition = "title = ?"
		args = append(args, title)
		if questionID != "" && questionID != "0" {
			condition = "(title = ? OR question_id = ?)"
			args = append(args, questionID)
		}
	} else {
		condition = "question_id = ?"
		args = append(args, questionID)
	}

	// First find all answer_ids to delete their comments
	rows, err := db.Query("SELECT answer_id FROM answers WHERE " + condition, args...)
	if err == nil {
		defer rows.Close()
		var ids []string
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err == nil {
				ids = append(ids, id)
			}
		}
		for _, id := range ids {
			db.Exec("DELETE FROM comments WHERE answer_id = ?", id)
		}
	}

	result, err := db.Exec("DELETE FROM answers WHERE " + condition, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func ToggleFavorite(db *sql.DB, id string) (int, error) {
	var current int
	err := db.QueryRow("SELECT is_favorite FROM answers WHERE answer_id = ?", id).Scan(&current)
	if err != nil {
		return 0, err
	}
	
	newFav := 0
	if current == 0 {
		newFav = 1
	}
	
	_, err = db.Exec("UPDATE answers SET is_favorite = ? WHERE answer_id = ?", newFav, id)
	return newFav, err
}

func ToggleGroupFavorite(db *sql.DB, title string, questionID string) (int, error) {
	var args []interface{}
	var condition string

	if title != "" {
		condition = "title = ?"
		args = append(args, title)
		if questionID != "" && questionID != "0" {
			condition = "(title = ? OR question_id = ?)"
			args = append(args, questionID)
		}
	} else {
		condition = "question_id = ?"
		args = append(args, questionID)
	}

	// Read current status from the first one
	var current int
	err := db.QueryRow("SELECT is_favorite FROM answers WHERE " + condition + " LIMIT 1", args...).Scan(&current)
	if err != nil {
		return 0, err
	}

	newFav := 0
	if current == 0 {
		newFav = 1
	}

	updateArgs := append([]interface{}{newFav}, args...)
	_, err = db.Exec("UPDATE answers SET is_favorite = ? WHERE " + condition, updateArgs...)
	return newFav, err
}

func GetAllTags(db *sql.DB) []string {
	rows, err := db.Query("SELECT DISTINCT tag FROM answers WHERE tag != ''")
	var tags []string
	if err != nil {
		return tags
	}
	defer rows.Close()
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err == nil {
			tags = append(tags, tag)
		}
	}
	return tags
}
