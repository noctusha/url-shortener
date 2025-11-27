package shortenerhandler

import (
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

func (h *Handler) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.shortenerhandler.url.save"

		logger := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		// render.Bind = DecodeJSON (json.Unmarshal) + Bind(); логика DTO
		if err := render.Bind(r, &req); err != nil {
			logger.Error("failed to bind request", slog.String("error", err.Error()))
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		// логика правил (struct tags)
		if err := h.v.Struct(req); err != nil {
			valErr := err.(validator.ValidationErrors)
			logger.Error("failed to validate request", slog.String("error", err.Error()))
			render.JSON(w, r, resp.ValidationError(valErr))
			return
		}

		logger.Info("request body decoded", slog.Any("req", req))

		id, alias, err := h.svc.URLSave(r.Context(), req.Url, req.Alias)
		if err != nil {
			if errors.Is(err, short.ErrAliasAlreadyExists) {
				logger.Info("url already exists", slog.String("url", req.Url))
				render.JSON(w, r, resp.Error("url already exists"))
				return
			}
			logger.Error("failed to save url", slog.String("error", err.Error()))
			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		logger.Info("url saved", slog.Int("id", int(id)))

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
