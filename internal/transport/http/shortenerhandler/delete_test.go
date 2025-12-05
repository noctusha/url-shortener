package shortenerhandler

import (
	"encoding/json"
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

type deleteErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func TestDeleteHandler(t *testing.T) {
	type tc struct {
		name        string
		alias       string
		mockErr     error
		wantStatus  int
		wantRespErr string
	}
	tests := []tc{
		{
			name:        "Success",
			alias:       "bbc",
			mockErr:     nil,
			wantStatus:  http.StatusOK,
			wantRespErr: "",
		},
		{
			name:        "Empty alias",
			alias:       "",
			mockErr:     nil,
			wantStatus:  http.StatusBadRequest,
			wantRespErr: "invalid request",
		},
		{
			name:        "URL not found",
			alias:       "something complicated and definitely not used in DB",
			mockErr:     shortener.ErrURLNotFound,
			wantStatus:  http.StatusNotFound,
			wantRespErr: "url not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewShortener(t)
			v := validator.New()
			h := New(logger.NewEmptyLogger(), v, m)

			if tt.alias != "" {
				m.On("DeleteURL",
					mock.Anything,
					tt.alias,
				).Return(
					tt.mockErr,
				).Once()
			}

			url := "/url/"
			if tt.alias != "" {
				url = "/url/" + tt.alias
			}

			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			var handler http.Handler
			if tt.alias == "" {
				handler = h.Delete()
			} else {
				r := chi.NewRouter()
				r.Delete("/url/{alias}", h.Delete())
				handler = r
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)

			if tt.wantStatus == http.StatusOK {
				var resp DeleteResponse
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
				require.Equal(t, tt.alias, resp.Alias)
				require.Equal(t, "OK", resp.Status)
				return
			}

			var resp deleteErrorResponse
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			require.Equal(t, "Error", resp.Status)
			require.Equal(t, tt.wantRespErr, resp.Error)
		})
	}
}
