package shortenerhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/noctusha/url-shortener/internal/logger"
	"github.com/noctusha/url-shortener/internal/service/shortener"
	"github.com/noctusha/url-shortener/internal/transport/http/shortenerhandler/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type redirectErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func TestRedirectHandler(t *testing.T) {
	type tc struct {
		name         string
		alias        string
		mockURL      string
		mockErr      error
		wantStatus   int
		wantLocation string
		wantRespErr  string
	}

	tests := []tc{
		{
			name:         "Success",
			alias:        "bbc",
			mockURL:      "https://bbc.com",
			mockErr:      nil,
			wantStatus:   http.StatusFound,
			wantLocation: "https://bbc.com",
			wantRespErr:  "",
		},
		{
			name:        "Empty alias",
			alias:       "",
			mockURL:     "",
			mockErr:     nil,
			wantStatus:  http.StatusNotFound,
			wantRespErr: "invalid request",
		},
		{
			name:        "URL not found",
			alias:       "something complicated and definitely not used in DB",
			mockURL:     "",
			mockErr:     shortener.ErrURLNotFound,
			wantStatus:  http.StatusNotFound,
			wantRespErr: "url not found",
		},
		{
			name:        "Internal error",
			alias:       "bbc",
			mockURL:     "",
			mockErr:     errors.New("db is down"),
			wantStatus:  http.StatusInternalServerError,
			wantRespErr: "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewShortener(t)
			v := validator.New()
			h := New(logger.NewEmptyLogger(), v, m)

			if tt.alias != "" {
				m.On("GetURL",
					mock.Anything,
					tt.alias,
				).Return(
					tt.mockURL,
					tt.mockErr,
				).Once()
			}

			url := "/"
			if tt.alias != "" {
				url = "/" + tt.alias
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			r := chi.NewRouter()
			r.Get("/{alias}", h.Redirect())

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)

			switch tt.wantStatus {
			case http.StatusFound:
				location := rr.Header().Get("Location")
				require.Equal(t, tt.wantLocation, location)
			case http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError:
				if tt.alias == "" {
					return
				}
				var resp redirectErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, tt.wantRespErr, resp.Error)
				require.Equal(t, "Error", resp.Status)
			default:
				t.Fatalf("Unexpected status code: %d", rr.Code)
			}
		})
	}
}
