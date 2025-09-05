package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Current number of HTTP requests being processed",
	})
)

// PrometheusMiddleware returns a Fiber middleware that collects metrics
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		httpRequestsInFlight.Inc()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get endpoint path (normalize it to avoid high cardinality)
		endpoint := normalizeEndpoint(c.Path())
		method := c.Method()
		statusCode := strconv.Itoa(c.Response().StatusCode())

		// Record metrics
		httpRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint, statusCode).Observe(duration)
		httpRequestsInFlight.Dec()

		return err
	}
}

// normalizeEndpoint normalizes endpoint paths to avoid high cardinality
func normalizeEndpoint(path string) string {
	switch {
	case path == "/v1/course/getcourse":
		return "/v1/course/getcourse"
	case path == "/v1/user/getuser":
		return "/v1/user/getuser"
	case path == "/v1/report/get-top-courses-report":
		return "/v1/report/get-top-courses-report"
	case path == "/metrics":
		return "/metrics"
	default:
		return "other"
	}
}
