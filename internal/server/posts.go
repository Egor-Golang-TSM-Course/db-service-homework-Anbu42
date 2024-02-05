package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"blog/internal/models"

	"github.com/go-chi/chi/v5"
)

func (s *Server) GetPostByID() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		posts, err := s.psql.GetPostByID(postID)
		if err != nil {
			http.Error(w, "Error fetching post", http.StatusInternalServerError)
			return
		}

		s.respond(w, http.StatusOK, posts)
	}
}

func (s *Server) GetPosts() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tags := r.URL.Query()["tags"]
		dateStr := r.URL.Query().Get("date")

		var date *time.Time
		if dateStr != "" {
			parsedDate, err := time.Parse(time.RFC3339, dateStr)
			if err != nil {
				// Обработка ошибки парсинга даты, если необходимо
				http.Error(w, "Invalid date format", http.StatusBadRequest)
				return
			}
			date = &parsedDate
		}

		// Обработка параметров пагинации
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		posts, err := s.psql.GetPosts(pageSize, page, tags, date)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.respond(w, http.StatusOK, posts)
	}
}

func (s *Server) CreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post models.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		post.Created = time.Now()

		postID, err := s.psql.CreatePost(&post)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.respond(w, http.StatusCreated, map[string]int{"postID": postID})
	}
}

func (s *Server) UpdatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updatedPost models.Post
		err := json.NewDecoder(r.Body).Decode(&updatedPost)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		updatedPost.ID = postID
		err = s.psql.UpdatePost(&updatedPost)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.respond(w, http.StatusOK, map[string]string{"message": "Post updated successfully"})
	}
}

func (s *Server) DeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		err = s.psql.DeletePost(postID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.respond(w, http.StatusOK, map[string]string{"message": "Post deleted successfully"})
	}
}
