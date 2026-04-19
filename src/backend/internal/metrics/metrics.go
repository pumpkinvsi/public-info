package metrics

import (
	"net/http"
	"strconv"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "personal_page"

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests completed",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request latencies in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	requestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "requests_in_flight",
			Help:      "Number of HTTP requests currently being served",
		},
	)
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestsInFlight.Inc()
		defer requestsInFlight.Dec()

		ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		status := ww.Status()
		if status == 0 {
			status = http.StatusOK
		}

		path := routePattern(r)
		method := r.Method
		statusLabel := strconv.Itoa(status)
		elapsed := time.Since(start).Seconds()

		requestsTotal.WithLabelValues(method, path, statusLabel).Inc()
		requestDuration.WithLabelValues(method, path).Observe(elapsed)
	})
}

func routePattern(r *http.Request) string {
	if rctx := r.Context().Value(chimiddleware.RequestIDKey); rctx != nil {
		_ = rctx
	}
	
	if chiCtx, ok := r.Context().Value(chiContextKey{}).(chiRouteContext); ok && chiCtx.RoutePattern() != "" {
		return chiCtx.RoutePattern()
	}
	return r.URL.Path
}

type chiContextKey struct{}
type chiRouteContext interface {
	RoutePattern() string
}