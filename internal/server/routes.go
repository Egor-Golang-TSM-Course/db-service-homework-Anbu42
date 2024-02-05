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
		r.Get("/", s.GetPosts())
		r.Get("/{id}", s.GetPostByID())
		r.Post("/", s.CreatePost())
		r.Put("/{id}", s.UpdatePost())
		r.Delete("/{id}", s.DeletePost())
	})

}
