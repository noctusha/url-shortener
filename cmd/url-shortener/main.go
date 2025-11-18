package main

import (
	"log/slog"
	"os"

	"github.com/noctusha/url-shortener/internal/config"
	"github.com/noctusha/url-shortener/internal/logger"
	"github.com/noctusha/url-shortener/internal/storage/postgres"
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

	// init router: chi

	// middleware

	// run server
}
