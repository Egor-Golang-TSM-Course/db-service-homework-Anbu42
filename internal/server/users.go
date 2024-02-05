package server

import (
	"blog/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RegisterHandler обрабатывает запрос на регистрацию пользователя
func (s *Server) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			Name     string `json:"name"`
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		user := models.User{
			Name:     requestData.Name,
			Login:    requestData.Login,
			Password: requestData.Password,
			Created:  time.Now(),
		}

		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			fmt.Println("err")
			return
		}

		token, err := s.psql.Register(&user)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			fmt.Println(err, "err2")
			return
		}

		s.respond(w, http.StatusCreated, map[string]string{"token": token})
	}
}

// LoginHandler обрабатывает запрос на вход пользователя
func (s *Server) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		user := models.User{
			Login:    requestData.Login,
			Password: requestData.Password,
		}

		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		token, err := s.psql.Login(&user)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		s.respond(w, http.StatusCreated, map[string]string{"token": token})
	}
}
