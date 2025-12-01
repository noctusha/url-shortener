package logger

import (
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func New(env string) *slog.Logger {
	var log *slog.Logger

	timestampReplace := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.String("timestamp", a.Value.Time().Format("2006-01-02 15:04:05"))
		}
		return a
	}

	switch env {
	case envLocal:
		// использую stderr (standard error), поскольку это канал для ошибок и диагностических сообщений
		log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
			// кастомизация логов - например, времени
			ReplaceAttr: timestampReplace,
		}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			AddSource:   true,
			ReplaceAttr: timestampReplace,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level:     slog.LevelError,
			AddSource: true,
			// кастомизация логов, например, времени
			ReplaceAttr: timestampReplace,
		}))
	default:
		// fallback
		log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:       slog.LevelDebug,
			AddSource:   true,
			ReplaceAttr: timestampReplace,
		}))

		log.Warn("unknown env, fallback to local", slog.String("env", env))
	}
	return log
}
