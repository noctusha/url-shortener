package shortener_handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	short "github.com/noctusha/url-shortener/internal/service/shortener"
	resp "github.com/noctusha/url-shortener/internal/transport/http/response"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Alias string `json:"alias,omitempty"`
	resp.Response
}

type Shortener interface {
	URLSave(ctx context.Context, url, alias string) (int32, string, error)
}

func New(log *slog.Logger, v *validator.Validate, shortener Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.shortener_handler.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		// render.Bind = DecodeJSON (json.Unmarshal) + Bind(); логика DTO
		if err := render.Bind(r, &req); err != nil {
			log.Error("failed to bind request", slog.String("error", err.Error()))
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		// логика правил (struct tags)
		if err := v.Struct(req); err != nil {
			valErr := err.(validator.ValidationErrors)
			log.Error("failed to validate request", slog.String("error", err.Error()))
			render.JSON(w, r, resp.ValidationError(valErr))
			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		id, alias, err := shortener.URLSave(r.Context(), req.Url, req.Alias)
		if err != nil {
			if errors.Is(err, short.ErrAliasAlreadyExists) {
				log.Info("url already exists", slog.String("url", req.Url))
				render.JSON(w, r, resp.Error("url already exists"))
				return
			}
			log.Error("failed to save url", slog.String("error", err.Error()))
			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		log.Info("url saved", slog.Int("id", int(id)))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}

func (r *Request) Bind(_ *http.Request) error {
	if r.Url == "" {
		return errors.New("url is required")
	}
	return nil
}
