package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"stellar_journal/internal/apod_worker"
	"stellar_journal/internal/config"
	mwLg "stellar_journal/internal/http-server/middleware/logger"
	"stellar_journal/internal/lib/logger/sl"
	"stellar_journal/internal/stellar_api/nasa_api"
	"stellar_journal/internal/storage/postgresql"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting stellar journal",
		slog.String("env", cfg.Env),
		slog.String("version", "1.0.0"),
	)

	log.Debug("debug messages are enabled")

	storage, err := postgresql.NewStorage(cfg.Storage.DbUri)
	if err != nil {
		log.Error("failed to create storage", sl.Err(err))
		os.Exit(1)
	}

	apiConn := nasa_api.NewNasaApiConnect(cfg.NasaApi.Host, cfg.NasaApi.Token)

	apodWorker := apod_worker.NewAPODWorker(apiConn, storage, log)
	go apodWorker.Run()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLg.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/journal", func(r chi.Router) {

	})

	log.Info("starting server", slog.String("address", cfg.HttpServer.Host))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HttpServer.Host,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.CtxTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	if err := storage.Close(); err != nil {
		log.Error("failed to close storage", sl.Err(err))
	}

	log.Info("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
