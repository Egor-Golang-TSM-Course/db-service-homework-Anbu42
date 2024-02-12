package models

import (
	"time"
)

type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
}

type Post struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	UserID  string    `json:"user_id"`
	Created time.Time `json:"created"`
	*Comment
	*Tags
}

type Comment struct {
	ID      int       `json:"id"`
	PostID  int       `json:"post_id"`
	UserID  string    `json:"user_id"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}

type Tags struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	PostID int    `json:"post_id"`
}
