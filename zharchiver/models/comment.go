package models

import "database/sql"

type Comment struct {
	ID        int    `json:"id"`
	AnswerID  string `json:"answer_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func GetComments(db *sql.DB, answerID string) ([]Comment, error) {
	rows, err := db.Query(`
		SELECT id, answer_id, content, created_at
		FROM comments
		WHERE answer_id = ?
		ORDER BY created_at ASC
	`, answerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.AnswerID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	if list == nil {
		list = []Comment{}
	}
	return list, nil
}

func AddComment(db *sql.DB, answerID string, content string) error {
	_, err := db.Exec(`
		INSERT INTO comments (answer_id, content) VALUES (?, ?)
	`, answerID, content)
	return err
}
