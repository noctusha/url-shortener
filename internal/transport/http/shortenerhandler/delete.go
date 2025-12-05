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

type DeleteRequest struct {
	Alias string `json:"alias"`
}

type DeleteResponse struct {
	Alias string `json:"alias"`
	resp.Response
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.shortenerhandler.Delete"

		logger := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Warn("alias is empty")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp.Error("invalid request"))
			return
		}

		err := h.svc.DeleteURL(r.Context(), alias)
		if err != nil {
			if errors.Is(err, shortener.ErrURLNotFound) {
				logger.Info("url not found", slog.String("alias", alias))
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(resp.Error("url not found"))
				return
			}
			logger.Error("failed to delete url by alias",
				slog.String("alias", alias),
				slog.String("error", err.Error()),
			)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp.Error("internal error"))
			return
		}

		logger.Info("url for alias deleted", slog.String("alias", alias))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteResponse{
			Alias:    alias,
			Response: resp.OK(),
		})
	}
}
