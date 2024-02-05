package models

import (
	"time"
)

// User представляет собой структуру данных для пользователя
type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
}

// Post представляет собой структуру данных для поста
type Post struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	UserID  string    `json:"user_id"`
	Created time.Time `json:"created"`
	*Comment
	*Tags
}

// Comment представляет собой структуру данных для комментария
type Comment struct {
	ID      int       `json:"id"`
	PostID  int       `json:"post_id"`
	UserID  string    `json:"user_id"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}

// Tag представляет собой структуру данных для тега
type Tags struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	PostID int    `json:"post_id"`
}
