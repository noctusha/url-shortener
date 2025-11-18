package main

import (
	"log/slog"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/noctusha/url-shortener/internal/config"
	"github.com/noctusha/url-shortener/internal/logger"
	"github.com/noctusha/url-shortener/internal/service/shortener"
	"github.com/noctusha/url-shortener/internal/storage/postgres"
	mw "github.com/noctusha/url-shortener/internal/transport/http/middleware/logger"
)

func main() {
	// init config: cleanenv
	cfg := config.New()

	// init logger: slog
	log := logger.New(cfg.Env)
	log.Info("starting main", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// init storage: postgres
	storage, err := postgres.New(cfg, log)
	if err != nil {
		log.Error("failed to initialize storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer storage.Close()

	// repository (postgres)
	urlRepo := postgres.NewURLRepository(storage.Conn())

	// service
	service := shortener.NewService(urlRepo, log)
	_ = service

	// init router: chi
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID) // tracing requester's ID
	router.Use(middleware.Logger)    // additional logs
	router.Use(mw.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// run server
}
