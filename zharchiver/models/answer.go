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
}

func SaveAnswer(db *sql.DB, data *AnswerData) error {
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

func GetAnswers(db *sql.DB) ([]AnswerSummary, error) {
	rows, err := db.Query(`
		SELECT answer_id, question_id, title, author_name, created_time, updated_time, saved_at, tag, tag_color
		FROM answers
		ORDER BY saved_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AnswerSummary
	for rows.Next() {
		var item AnswerSummary
		if err := rows.Scan(&item.AnswerID, &item.QuestionID, &item.Title, &item.AuthorName, &item.CreatedTime, &item.UpdatedTime, &item.SavedAt, &item.Tag, &item.TagColor); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	if list == nil {
		list = []AnswerSummary{}
	}
	return list, nil
}

func GetAnswerByID(db *sql.DB, id string) (*AnswerData, error) {
	var data AnswerData
	err := db.QueryRow(`
		SELECT answer_id, question_id, title, author_name, content_html, created_time, updated_time, saved_at, tag, tag_color
		FROM answers
		WHERE answer_id = ?
	`, id).Scan(&data.AnswerID, &data.QuestionID, &data.Title, &data.AuthorName, &data.ContentHTML, &data.CreatedTime, &data.UpdatedTime, &data.SavedAt, &data.Tag, &data.TagColor)

	if err != nil {
		return nil, err
	}
	return &data, nil
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
