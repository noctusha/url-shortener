package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/noctusha/url-shortener/internal/config"
	"github.com/noctusha/url-shortener/internal/logger"
	"github.com/noctusha/url-shortener/internal/service/shortener"
	"github.com/noctusha/url-shortener/internal/storage/postgres"
	mw "github.com/noctusha/url-shortener/internal/transport/http/middleware/logger"
	handler "github.com/noctusha/url-shortener/internal/transport/http/shortener_handler"
)

func main() {
	// init config: cleanenv
	cfg := config.MustLoad()

	// init logger: slog
	log := logger.New(cfg.Env)
	log.Info("starting main", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// init storage: postgres (infrastructure layer)
	storage, err := postgres.New(cfg, log)
	if err != nil {
		log.Error("failed to initialize storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer storage.Close()

	// repository, storage repo
	urlRepo := postgres.NewURLRepository(storage.Conn())

	// service (business logic)
	service := shortener.NewService(urlRepo, log)
	_ = service

	// init router: chi
	router := chi.NewRouter()

	// http layer?

	// middleware
	router.Use(middleware.RequestID) // tracing requester's ID
	router.Use(middleware.Logger)    // additional logs
	router.Use(mw.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	v := validator.New()
	router.Post("/url", handler.New(log, v, service))

	log.Info("starting server", slog.String("address", cfg.Host))
	// init server
	srv := http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,      // ограничение чтения запроса
		WriteTimeout: cfg.Timeout,      // ограничение записи ответа
		IdleTimeout:  30 * time.Second, // время простоя keep-alive соединения
	}

	// run server
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", slog.String("error", err.Error()))
	}
	log.Info("server stopped", slog.String("address", cfg.Host))
}
