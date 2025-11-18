package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

// обертка над стандартным middleware chi
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	// сам middleware, соответствующий сигнатуре chi
	return func(next http.Handler) http.Handler {
		// перезапись логгера с доп. полем
		log = log.With(
			slog.String("component", "middleware/logger"),
		)
		// однократный инфо-лог о включении middleware
		log.Info("logger middleware enabled")

		// инициализация http-handler, которым станет middleware при каждом запросе
		fn := func(w http.ResponseWriter, r *http.Request) {
			// инициализация логгера для этого хэндлера с доп. полями
			// выполняется до обработки запроса
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			// обертка стандартного ResponseWriter в chi.ResponseWriter для возможности трекинга Status() + BytesWritten()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			// выполняется после обработки запроса
			defer func() {
				// использование логгера с инфой от chi.ResponseWriter
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.Duration("duration", time.Since(t1)),
				)
			}()
			// Передаём управление следующему обработчику в цепочке middleware.
			// next - это либо следующий middleware, либо конечный handler маршрута.
			// Logger → ...(e.g. Auth → RateLimit) → Business Logic → Response
			next.ServeHTTP(ww, r)
		}
		// превращение func(w http.ResponseWriter, r *http.Request) в интерфейс http.Handler
		return http.HandlerFunc(fn)
	}
}
