package main

import (
	"os"

	"log/slog"

	"github.com/MaKYaro/url-shortener/internal/config"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	// TODO: parse config
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("logger is working", slog.String("env", cfg.Env))
	log.Debug("debug masseges are enabled")

	// TODO: init storage

	// TODO: init router

	// TODO: run app
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
		)
	case envProd:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		),
		)
	}
	return log
}
