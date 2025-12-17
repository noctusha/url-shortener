package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	appmetrics "github.com/noctusha/url-shortener/internal/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (w *responseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.status = code
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
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
			"status": strconv.Itoa(rw.status),
		}).Inc()

		appmetrics.HTTPRequestDuration.With(prometheus.Labels{
			"method": r.Method,
			"route":  route,
		}).Observe(time.Since(start).Seconds())
	})
}
