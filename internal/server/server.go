package server

import (
	"blog/internal/config"
	"blog/internal/storage/psql"
	"blog/internal/storage/redis"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg       *config.Config
	psql      psql.PSQL
	cache     redis.RedisCache
	router    chi.Router
	log       *logrus.Logger
	validator *validator.Validate
}

func New(cfg *config.Config,
	psql psql.PSQL,
	cache *redis.RedisCache,
	router chi.Router,
	log *logrus.Logger,
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
