package logger

import (
	"context"
	"log/slog"
)

func NewEmptyLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type EmptyLogger struct{}

func NewDiscardHandler() *EmptyLogger {
	return &EmptyLogger{}
}

func (h *EmptyLogger) Handle(_ context.Context, _ slog.Record) error {
	// Просто игнорируем запись журнала
	return nil
}

func (h *EmptyLogger) WithAttrs(_ []slog.Attr) slog.Handler {
	// Возвращает тот же обработчик, так как нет атрибутов для сохранения
	return h
}

func (h *EmptyLogger) WithGroup(_ string) slog.Handler {
	// Возвращает тот же обработчик, так как нет группы для сохранения
	return h
}

func (h *EmptyLogger) Enabled(_ context.Context, _ slog.Level) bool {
	// Всегда возвращает false, так как запись журнала игнорируется
	return false
}
