package logger

import (
	"blog/internal/config"
	"strings"

	"github.com/sirupsen/logrus"
)

func Logger(cfg *config.Config) *logrus.Logger {

	logger := logrus.New()

	logLevel, err := logrus.ParseLevel(strings.ToLower(cfg.Log.Level))
	if err != nil {
		logger.Fatalf("Invalid log level: %s", cfg.Log.Level)
	}
	logger.SetLevel(logLevel)

	if cfg.Log.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	return logger
}
