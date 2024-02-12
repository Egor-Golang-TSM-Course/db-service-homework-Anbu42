package cmd

import (
	"blog/internal/config"
	"blog/internal/logger"
	"blog/internal/server"
	"blog/internal/storage/psql"
	"blog/internal/storage/redis"
	"flag"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func Start() error {
	router := chi.NewRouter()

	configPath := flag.String("cfg", "configs/config.yaml", "path to file config")
	flag.Parse()
	if *configPath == "" {
		panic("nil config file")
	}

	validator := validator.New()
	cfg := config.GetConfig(*configPath, validator)

	log := logger.Logger(cfg)

	database, err := psql.NewPostgreSQL(cfg)
	if err != nil {
		slog.Error("blogservice.NewDB", slog.String("err", err.Error()))
		os.Exit(1)
	}

	cache := redis.NewRedisCache(cfg)

	serve := server.New(cfg, database, cache, router, log, validator)

	return serve.Start()
}
