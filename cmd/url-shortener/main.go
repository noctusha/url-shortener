package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/noctusha/url-shortener/internal/config"
	"github.com/noctusha/url-shortener/internal/logger"
	"github.com/noctusha/url-shortener/internal/service/shortener"
	"github.com/noctusha/url-shortener/internal/storage/postgres"
	mw "github.com/noctusha/url-shortener/internal/transport/http/middleware/logger"
	handler "github.com/noctusha/url-shortener/internal/transport/http/shortenerhandler"
)

func main() {
	// каналы для graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

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

	// middleware
	router.Use(middleware.RequestID) // tracing requester's ID
	router.Use(middleware.Logger)    // additional logs
	router.Use(mw.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	v := validator.New()
	// http layer
	hand := handler.New(log, v, service)
	router.Post("/url", hand.Save())

	log.Info("starting server", slog.String("address", cfg.HTTPAddr))
	// init server
	srv := http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,      // ограничение чтения запроса
		WriteTimeout: cfg.Timeout,      // ограничение записи ответа
		IdleTimeout:  30 * time.Second, // время простоя keep-alive соединения
	}

	// run server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server run error", slog.String("error", err.Error()))
		}
	}()

	<-signalChan
	log.Info("termination signal received, shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// плавная остановка сервера
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced shutdown", slog.String("error", err.Error()))
		return
	}
	log.Info("server gracefully stopped")
}
