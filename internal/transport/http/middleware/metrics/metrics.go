package metrics

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	appmetrics "github.com/noctusha/url-shortener/internal/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		route := chi.RouteContext(r.Context()).RoutePattern()
		if route == "" {
			route = "unknown"
		}

		appmetrics.HTTPRequestsTotal.With(prometheus.Labels{
			"method": r.Method,
			"route":  route,
			"status": http.StatusText(rw.status),
		}).Inc()

		appmetrics.HTTPRequestDuration.With(prometheus.Labels{
			"method": r.Method,
			"route":  route,
		}).Observe(time.Since(start).Seconds())
	})
}
