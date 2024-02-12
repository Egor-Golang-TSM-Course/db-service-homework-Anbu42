package server

import (
	"blog/internal/models"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (s *Server) CreateTags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tags models.Tags
		if err := json.NewDecoder(r.Body).Decode(&tags); err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := s.psql.CreateTags(&tags)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Failed to create tags", http.StatusInternalServerError)
			return
		}

		s.respond(w, http.StatusOK, map[string]string{"message": "Tags created successfully"})
	}
}

func (s *Server) GetTags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tags, err := s.psql.GetTags()
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"method": r.Method,
				"URL":    r.URL.Path,
			}).Error(err)

			http.Error(w, "Failed to get tags", http.StatusInternalServerError)
			return
		}

		s.respond(w, http.StatusOK, tags)
	}
}
