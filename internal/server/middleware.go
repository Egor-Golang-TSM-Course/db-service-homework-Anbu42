package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	userIDContextKey contextKey = "userID"
)

// Authentication проверяет JWT токен и добавляет UserID в контекст запроса
func (s *Server) CheckAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				http.Error(w, fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]), http.StatusBadRequest)
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(s.cfg.JwtSecretKey), nil
		})

		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Unable to extract claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "Unauthorized: Unable to extract user ID", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) CheckOwnerPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		post, err := s.psql.GetPostByID(postID)
		if err != nil {
			http.Error(w, "Unauthorized: User not owner this post", http.StatusInternalServerError)
			return
		}

		userID, ok := r.Context().Value(userIDContextKey).(string)
		if !ok {
			http.Error(w, "Unauthorized: User ID not found in context", http.StatusUnauthorized)
			return
		}

		if post.UserID != userID {
			http.Error(w, "Unauthorized: Unable to extract claims", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) CheckOwnerComment(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		comIDStr := chi.URLParam(r, "id")
		comID, err := strconv.Atoi(comIDStr)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		com, err := s.psql.GetCommentByID(comID)
		if err != nil {
			http.Error(w, "Unauthorized: User not owner this comment", http.StatusInternalServerError)
			return
		}

		userID, ok := r.Context().Value(userIDContextKey).(string)
		if !ok {
			http.Error(w, "Unauthorized: User ID not found in context", http.StatusUnauthorized)
			return
		}

		if com.UserID != userID {
			http.Error(w, "Unauthorized: Unable to extract claims", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
