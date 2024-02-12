package psql

import (
	"blog/internal/models"
	"context"
)

func (pgdb *PostgreSQL) CreateTags(tags *models.Tags) error {
	tx, err := pgdb.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background()) // Откатываем транзакцию в случае ошибки

	_, err = tx.Exec(context.Background(), "INSERT INTO tags(name, post_id) VALUES($1, $2)", tags.Name, tags.PostID)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (pgdb *PostgreSQL) GetTags() ([]string, error) {
	query := `SELECT DISTINCT name FROM tags`

	rows, err := pgdb.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string

	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
