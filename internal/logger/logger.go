package logger

import (
	"log/slog"
	"os"
)

const (
	JSON = iota + 1
	TEXT
)

func Logger(model int) *slog.Logger {
	switch model {
	case JSON:
		logHandler := slog.NewJSONHandler(os.Stdout, nil)
		return slog.New(logHandler)

	default:
		logHandler := slog.NewTextHandler(os.Stdout, nil)
		return slog.New(logHandler)
	}
}
