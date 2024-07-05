package main

import (
	"fmt"
	"os"
	"time"

	"log/slog"

	"github.com/MaKYaro/url-shortener/internal/config"
	"github.com/MaKYaro/url-shortener/internal/domain"
	"github.com/MaKYaro/url-shortener/internal/storage/postgres"
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

	storage, err := postgres.New(cfg.DBConn)
	if err != nil {
		log.Error("can't init storage", slog.String("error", err.Error()))
	}
	err = storage.SaveURL(
		&domain.Alias{
			Value:  "goose",
			URL:    "https://habr.com/ru/articles/780280/",
			Expire: time.Now(),
		},
	)
	if err != nil {
		log.Error("can't save url", slog.String("error", err.Error()))
	}
	err = storage.SaveURL(
		&domain.Alias{
			Value:  "goose",
			URL:    "https://habr.com/ru/articles/780280/",
			Expire: time.Now(),
		},
	)
	if err != nil {
		log.Error("can't save url", slog.String("error", err.Error()))
	}
	alias, err := storage.GetURL("goose")
	fmt.Println(alias)
	if err != nil {
		log.Error("can't get url", slog.String("error", err.Error()))
	}
	_, err = storage.GetURL("goo")
	if err != nil {
		log.Error("can't get url", slog.String("error", err.Error()))
	}

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
