package shortenerhandler

import (
	"context"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

type Shortener interface {
	URLSave(ctx context.Context, url, alias string) (int32, string, error)
}
type Handler struct {
	log *slog.Logger
	v   *validator.Validate
	svc Shortener
}

func New(log *slog.Logger, v *validator.Validate, svc Shortener) *Handler {
	return &Handler{log: log, v: v, svc: svc}
}
