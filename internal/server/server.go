package server

import (
	"blog/internal/config"
	"blog/internal/storage/psql"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	cfg       *config.Config
	psql      psql.PSQL
	router    chi.Router
	log       *slog.Logger
	validator *validator.Validate
}

func New(cfg *config.Config,
	psql psql.PSQL,
	router chi.Router,
	log *slog.Logger,
	validator *validator.Validate) *Server {
	return &Server{
		cfg:       cfg,
		psql:      psql,
		router:    router,
		log:       log,
		validator: validator,
	}
}

func (s *Server) Start() error {
	s.log.Info("server started...")
	s.routes()

	return http.ListenAndServe(":8080", s.router)
}
