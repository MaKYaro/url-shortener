package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/MaKYaro/url-shortener/internal/config"
	"github.com/MaKYaro/url-shortener/internal/http-server/router"
	"github.com/MaKYaro/url-shortener/internal/http-server/server"
	"github.com/MaKYaro/url-shortener/internal/lib/random"
	urlshortener "github.com/MaKYaro/url-shortener/internal/services/url-shortener"
	"github.com/MaKYaro/url-shortener/internal/storage/postgres"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("logger is working", slog.String("env", cfg.Env))
	log.Debug("debug masseges are enabled")

	storage, err := postgres.New(cfg.DBConn)
	defer storage.Close()

	if err != nil {
		log.Error("can't init storage", slog.String("error", err.Error()))
	}
	log.Info("storage is working", slog.Any("storage", cfg.DBConn))

	aliasGenerator := random.NewGenerator(cfg.Alias.Length)
	shortener := urlshortener.New(
		log,
		storage,
		storage,
		storage,
		aliasGenerator,
		cfg.Alias.LifeLength,
	)

	router := router.NewRouter(log, http.NewServeMux(), shortener)

	log.Info("starting server", slog.String("address", cfg.Server.Address))

	srv := server.NewServer(
		cfg.Server.Address,
		router.InitRoutes(),
		cfg.Server.Timeout,
		cfg.Server.IdleTimeout,
	)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

		sign := <-stop

		log.Info(
			"stopping application",
			slog.String("signal", sign.String()),
		)

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Error(
				"failed to shutdown server",
				slog.String("error", err.Error()),
			)
		}
		close(stop)
	}()

	if err := srv.Run(); err != nil {
		log.Error("failed to start server")
	}

	log.Info("server stopped")
	log.Info("application stopped")
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
