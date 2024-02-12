package server

import (
	"blog/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func (s *Server) CreateComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var comment models.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		commentID, err := s.psql.CreateComment(&comment)
		if err != nil {
			http.Error(w, "Failed to create comment", http.StatusInternalServerError)
			return
		}

		response := map[string]int{"comment_id": commentID}
		s.respond(w, http.StatusCreated, response)
	}
}

func (s *Server) GetCommentsByPostID() http.HandlerFunc {
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

		comments, err := s.psql.GetCommentsByPostID(postID)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Failed to get comment", http.StatusInternalServerError)
			return
		}

		s.respond(w, http.StatusOK, comments)
	}
}

func (s *Server) DeleteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commentIDStr := chi.URLParam(r, "id")
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		err = s.psql.DeleteComment(commentID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.respond(w, http.StatusOK, map[string]string{"message": "Comment deleted successfully"})
	}
}

func (s *Server) UpdateComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updatedComment models.Comment
		err := json.NewDecoder(r.Body).Decode(&updatedComment)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		commentIDStr := chi.URLParam(r, "id")
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		updatedComment.ID = commentID

		err = s.psql.UpdateComment(&updatedComment)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.respond(w, http.StatusOK, map[string]string{"message": "Comment updated successfully"})
	}
}
