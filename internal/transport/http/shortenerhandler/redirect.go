package shortenerhandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/noctusha/url-shortener/internal/service/shortener"
	resp "github.com/noctusha/url-shortener/internal/transport/http/response"
)

func (h *Handler) Redirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.shortenerhandler.Redirect"

		logger := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Warn("alias is empty")
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp.Error("invalid request"))
			return
		}

		url, err := h.svc.GetURL(r.Context(), alias)
		if err != nil {
			if errors.Is(err, shortener.ErrURLNotFound) {
				logger.Info("url not found", slog.String("alias", alias))
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp.Error("url not found"))
				return
			}
			logger.Error("failed to get url",
				slog.String("alias", alias),
				slog.String("error", err.Error()),
			)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp.Error("internal error"))
			return
		}

		logger.Info("redirecting to url", slog.String("url", url))
		http.Redirect(w, r, url, http.StatusFound)
	}
}
