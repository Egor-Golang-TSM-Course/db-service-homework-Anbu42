package psql

import (
	"blog/internal/models"
	"context"
	"time"
)

func (pgdb *PostgreSQL) GetPostByID(id int) (*models.Post, error) {
	query := `SELECT * FROM posts WHERE posts.id=$1`

	row, err := pgdb.pool.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var post models.Post
	err = row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Created)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (pgdb *PostgreSQL) GetPosts(pageSize, page int, tags []string, date *time.Time) ([]*models.Post, error) {
	offset := pageSize * (page - 1)
	query := `
	SELECT
		posts.id,
		posts.title,
		posts.content,
		posts.user_id,
		posts.created
	FROM
		posts
	LEFT JOIN
		tags ON tags.post_id = posts.id
	WHERE
		tags.name = ANY($1::text[]) OR $2::timestamp IS NULL OR date_trunc('day', posts.created) = date_trunc('day', $2)
	ORDER BY
		posts.id
	OFFSET $3
	LIMIT $4
`
	rows, err := pgdb.pool.Query(context.Background(), query, tags, date, offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (pgdb *PostgreSQL) CreatePost(post *models.Post) (int, error) {
	var postID int
	err := pgdb.pool.QueryRow(
		context.Background(),
		"INSERT INTO posts(title, content, user_id, created) VALUES($1, $2, $3, $4) RETURNING id",
		post.Title, post.Content, post.UserID, time.Now(),
	).Scan(&postID)
	if err != nil {
		return 0, err
	}
	return postID, nil
}

func (pgdb *PostgreSQL) UpdatePost(post *models.Post) error {
	_, err := pgdb.pool.Exec(
		context.Background(),
		"UPDATE posts SET title = $1, content = $2 WHERE id = $3",
		post.Title, post.Content, post.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (pgdb *PostgreSQL) DeletePost(id int) error {
	_, err := pgdb.pool.Exec(
		context.Background(),
		"DELETE FROM posts WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
