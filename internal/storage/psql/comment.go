package psql

import (
	"blog/internal/models"
	"context"
	"time"
)

func (pgdb *PostgreSQL) CreateComment(comment *models.Comment) (int, error) {
	var commentID int
	query := `INSERT INTO comments(post_id, user_id, content, created) VALUES($1, $2, $3, $4) RETURNING id`
	err := pgdb.pool.QueryRow(
		context.Background(),
		query,
		comment.PostID, comment.UserID, comment.Content, time.Now(),
	).Scan(&commentID)
	if err != nil {
		return 0, err
	}
	return commentID, nil
}

func (pgdb *PostgreSQL) GetCommentsByPostID(postID int) ([]*models.Comment, error) {
	query := `SELECT * FROM comments WHERE post_id = $1`
	rows, err := pgdb.pool.Query(context.Background(), query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]*models.Comment, 0)

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Created)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (pgdb *PostgreSQL) GetCommentByID(id int) (*models.Comment, error) {
	var comment models.Comment
	err := pgdb.pool.QueryRow(
		context.Background(),
		"SELECT * FROM comments WHERE id = $1",
		id,
	).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Created)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (pgdb *PostgreSQL) DeleteComment(id int) error {
	_, err := pgdb.pool.Exec(
		context.Background(),
		"DELETE FROM comments WHERE id = $1",
		id,
	)
	return err
}

func (pgdb *PostgreSQL) UpdateComment(comment *models.Comment) error {
	_, err := pgdb.pool.Exec(
		context.Background(),
		"UPDATE comments SET content = $1 WHERE id = $2",
		comment.Content, comment.ID,
	)
	return err
}
