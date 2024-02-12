package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) routes() {
	s.router.Use(middleware.Logger)

	s.router.Route("/users", func(r chi.Router) {
		r.Post("/register", s.Register())
		r.Post("/login", s.Login())
	})

	s.router.Route("/posts", func(r chi.Router) {
		r.With(s.CheckAuthentication).Get("/", s.GetPosts())
		r.With(s.CheckAuthentication).Get("/{id}", s.GetPostByID())
		r.With(s.CheckAuthentication).Post("/", s.CreatePost())
		r.With(s.CheckAuthentication).With(s.CheckOwnerPost).Put("/{id}", s.UpdatePost())
		r.With(s.CheckAuthentication).With(s.CheckOwnerPost).Delete("/{id}", s.DeletePost())
		r.With(s.CheckAuthentication).Post("/{id}/tags", s.CreateTags())
		r.With(s.CheckAuthentication).Get("/posts/{id}/comments", s.GetCommentsByPostID())
	})

	s.router.Route("/tags", func(r chi.Router) {
		r.With(s.CheckAuthentication).Get("/", s.GetTags())
	})

	s.router.Route("/comments", func(r chi.Router) {
		r.With(s.CheckAuthentication).Post("/", s.CreateComment())
		r.With(s.CheckAuthentication).With(s.CheckOwnerComment).Put("/{id}", s.UpdateComment())
		r.With(s.CheckAuthentication).With(s.CheckOwnerComment).Delete("/{id}", s.DeleteComment())
	})
}
