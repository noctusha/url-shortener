package shortenerhandler

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
)

//go:generate mockery --config=../../../../.mockery.yml
type Shortener interface {
	SaveURL(ctx context.Context, url, alias string, expireAt *time.Time) (int32, string, error)
	GetURL(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, alias string) error
}
type Handler struct {
	log *slog.Logger
	v   *validator.Validate
	svc Shortener
}

func New(log *slog.Logger, v *validator.Validate, svc Shortener) *Handler {
	return &Handler{log: log, v: v, svc: svc}
}
