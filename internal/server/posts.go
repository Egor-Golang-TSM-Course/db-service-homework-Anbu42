package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"blog/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func (s *Server) GetPostByID() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		cacheKey := fmt.Sprintf("post:%d", postID)
		cachedData, err := s.cache.Get(cacheKey)
		if err == nil {
			var post models.Post
			if err := json.Unmarshal([]byte(cachedData), &post); err != nil {
				s.log.WithFields(logrus.Fields{
					"method": r.Method,
					"URL":    r.URL.Path,
				}).Error(err)

				http.Error(w, "Error decoding cached data", http.StatusInternalServerError)
				return
			}
			s.respond(w, http.StatusOK, post)
			return
		}

		post, err := s.psql.GetPostByID(postID)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Error fetching post", http.StatusInternalServerError)
			return
		}

		postJSON, err := json.Marshal(post)
		if err == nil {
			s.cache.Set(cacheKey, string(postJSON))
		}

		s.respond(w, http.StatusOK, post)
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
				s.log.WithFields(logrus.Fields{
					"method": r.Method,
					"URL":    r.URL.Path,
				}).Error(err)
				http.Error(w, "Invalid date format", http.StatusBadRequest)
				return
			}
			date = &parsedDate
		}

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
			posts, err = s.psql.GetPosts(pageSize, page, tags, date)
			if err != nil {
				s.log.WithFields(logrus.Fields{
					"method": r.Method,
					"URL":    r.URL.Path,
				}).Error(err)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if err := s.cache.SetPosts(pageSize, page, tags, date, posts); err != nil {
				s.log.WithFields(logrus.Fields{
					"method": r.Method,
					"URL":    r.URL.Path,
				}).Error(err)
			}
		}

		s.respond(w, http.StatusOK, posts)
	}
}

func (s *Server) CreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post models.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		post.Created = time.Now()

		postID, err := s.psql.CreatePost(&post)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

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
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)
			return
		}

		updatedPost.ID = postID
		err = s.psql.UpdatePost(&updatedPost)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		cacheKey := fmt.Sprintf("post:%d", postID)
		if err := s.cache.UpdatePost(cacheKey, &updatedPost); err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)
			http.Error(w, "Error updating cache", http.StatusInternalServerError)
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
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		err = s.psql.DeletePost(postID)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		cacheKey := fmt.Sprintf("posts:%d", postID)
		err = s.cache.Delete(cacheKey)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)
			fmt.Printf("Error deleting cache for post ID %d: %v\n", postID, err)
		}

		s.respond(w, http.StatusOK, map[string]string{"message": "Post deleted successfully"})
	}
}
