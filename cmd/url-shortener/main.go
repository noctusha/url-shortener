package main

import (
	"log/slog"

	"github.com/noctusha/url-shortener/internal/config"
	"github.com/noctusha/url-shortener/internal/logger"
)

func main() {
	// init config: cleanenv
	cfg := config.New()

	// init logger: slog
	log := logger.New(cfg.Env)
	log.Info("starting main", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// init storage: postgres

	// init router: chi

	// run server
}
