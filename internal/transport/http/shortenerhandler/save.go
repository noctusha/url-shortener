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

type SaveRequest struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type SaveResponse struct {
	Alias string `json:"alias,omitempty"`
	resp.Response
}

func (h *Handler) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		const op = "http.shortenerhandler.save"

		logger := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SaveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Warn("invalid json", slog.String("error", err.Error()))
			resp.WriteJSON(w, http.StatusBadRequest, resp.Error("invalid JSON"))
			return
		}

		// логика правил (struct tags)
		if err := h.v.Struct(req); err != nil {
			valErr, ok := err.(validator.ValidationErrors)
			if !ok {
				resp.WriteJSON(w, http.StatusBadRequest, resp.Error("invalid request"))
				return
			}
			logger.Warn("validation failed", slog.Any("errors", valErr))
			resp.WriteJSON(w, http.StatusUnprocessableEntity, resp.ValidationError(valErr))
			return
		}

		logger.Info("request body decoded", slog.Any("req", req))

		id, alias, err := h.svc.SaveURL(r.Context(), req.URL, req.Alias)
		if err != nil {
			if errors.Is(err, short.ErrAliasAlreadyExists) {
				logger.Info("alias already exists",
					slog.String("url", req.URL),
					slog.String("alias", req.Alias),
				)
				resp.WriteJSON(w, http.StatusConflict, resp.Error("alias already exists"))
				return
			}
			logger.Error("failed to save url",
				slog.String("url", req.URL),
				slog.String("alias", alias),
				slog.String("error", err.Error()),
			)
			resp.WriteJSON(w, http.StatusInternalServerError, resp.Error("internal error"))
			return
		}

		logger.Info("url saved", slog.Int("id", int(id)))

		resp.WriteJSON(w, http.StatusOK, SaveResponse{
			Alias:    alias,
			Response: resp.OK(),
		})
	}
}
