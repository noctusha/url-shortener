package shortenerhandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
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
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Warn("invalid json", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp.Error("invalid JSON"))
			return
		}

		// логика правил (struct tags)
		if err := h.v.Struct(req); err != nil {
			valErr := err.(validator.ValidationErrors)
			logger.Warn("validation failed", slog.Any("errors", valErr))
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(resp.ValidationError(valErr))
			return
		}

		logger.Info("request body decoded", slog.Any("req", req))

		id, alias, err := h.svc.SaveURL(r.Context(), req.Url, req.Alias)
		if err != nil {
			if errors.Is(err, short.ErrAliasAlreadyExists) {
				logger.Info("alias already exists", slog.String("alias", req.Alias))
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(resp.Error("alias already exists"))
				return
			}
			logger.Error("unexpected error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp.Error("failed to save url"))
			return
		}

		logger.Info("url saved", slog.Int("id", int(id)))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Alias:    alias,
			Response: resp.OK(),
		})
	}
}
