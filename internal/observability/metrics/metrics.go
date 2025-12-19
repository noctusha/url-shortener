package metrics

import "github.com/prometheus/client_golang/prometheus"

var HTTPRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "url_shortener",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests",
	},
	[]string{"method", "route", "status"},
)

var HTTPRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "url_shortener",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "HTTP request latency",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"method", "route"},
)

var RateLimitBlockedTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace:   "url_shortener",
		Subsystem:   "ratelimit",
		Name:        "blocked_total",
		Help:        "Number of blocked requests by rate limiter",
		ConstLabels: nil,
	},
	[]string{"scope"}, // ip / ipAlias
)

var RateLimitErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "url_shortener",
		Subsystem: "ratelimit",
		Name:      "errors_total",
		Help:      "Rate limiter backend errors",
	},
	[]string{"scope"},
)

func Init() {
	prometheus.MustRegister(HTTPRequestsTotal, HTTPRequestDuration, RateLimitBlockedTotal, RateLimitErrorsTotal)
}
