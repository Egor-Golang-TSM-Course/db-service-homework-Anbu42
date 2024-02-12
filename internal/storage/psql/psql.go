package psql

import (
	"blog/internal/config"
	"blog/internal/models"
	"context"
	"time"

	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

type PSQL interface {
	User
	Post
	Comment
	Tags
}

type PostgreSQL struct {
	pool *pgxpool.Pool
}

func NewPostgreSQL(cfg *config.Config) (*PostgreSQL, error) {
	connStr := "user=" + cfg.Components.Database.Username +
		" password=" + cfg.Components.Database.Password +
		" dbname=" + cfg.Components.Database.Name +
		" host=" + cfg.Components.Database.Host +
		" port=" + cfg.Components.Database.Port +
		" sslmode=disable"

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.Components.Database.ConnectionsLimit)
	poolConfig.ConnConfig.ConnectTimeout, _ = time.ParseDuration(cfg.Components.Database.ConnectionTimeout)
	poolConfig.MaxConnLifetime, _ = time.ParseDuration(cfg.Components.Database.ConnectionLifetime)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return &PostgreSQL{pool: pool}, nil
}

type User interface {
	Register(JwtSecretKey string, user *models.User) (string, error)
	Login(JwtSecretKey string, user *models.User) (string, error)
}

type Post interface {
	CreatePost(post *models.Post) (int, error)
	GetPostByID(id int) (*models.Post, error)
	GetPosts(pageSize, page int, tags []string, date *time.Time) ([]*models.Post, error)
	UpdatePost(post *models.Post) error
	DeletePost(id int) error
}

type Comment interface {
	CreateComment(comment *models.Comment) (int, error)
	GetCommentsByPostID(postID int) ([]*models.Comment, error)
	GetCommentByID(id int) (*models.Comment, error)
	DeleteComment(id int) error
	UpdateComment(comment *models.Comment) error
}

type Tags interface {
	CreateTags(*models.Tags) error
	GetTags() ([]string, error)
}
