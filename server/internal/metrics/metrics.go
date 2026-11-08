package metrics

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/meshery/meshkit/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricNamespace = "meshery"
	unmatchedRoute  = "unmatched"
)

// Metrics manages isolated Prometheus telemetry for Meshery Server.
type Metrics struct {
	Registry             *prometheus.Registry
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	httpRequestsInFlight prometheus.Gauge
	buildInfo            *prometheus.GaugeVec
}

// New initializes an isolated Prometheus Registry and registers Go runtime,
// process, and Meshery-specific HTTP & build telemetry metrics.
func New(version, commitSHA, releaseChannel string) *Metrics {
	reg := prometheus.NewRegistry()

	// Register standard Go runtime and process metrics on isolated registry
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
		Namespace: metricNamespace,
	}))

	buildInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricNamespace,
			Name:      "build_info",
			Help:      "Build information for Meshery Server.",
		},
		[]string{"version", "commit_sha", "release_channel"},
	)

	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricNamespace,
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests processed by Meshery Server.",
		},
		[]string{"method", "route", "code"},
	)

	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricNamespace,
			Name:      "http_request_duration_seconds",
			Help:      "Histogram of HTTP request durations for non-streaming endpoints.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	httpRequestsInFlight := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: metricNamespace,
			Name:      "http_requests_in_flight",
			Help:      "Current number of in-flight HTTP requests being processed by Meshery Server.",
		},
	)

	reg.MustRegister(buildInfo, httpRequestsTotal, httpRequestDuration, httpRequestsInFlight)

	// Set constant build_info gauge to 1 with provided build metadata
	buildInfo.WithLabelValues(version, commitSHA, releaseChannel).Set(1)

	return &Metrics{
		Registry:             reg,
		httpRequestsTotal:    httpRequestsTotal,
		httpRequestDuration:  httpRequestDuration,
		httpRequestsInFlight: httpRequestsInFlight,
		buildInfo:            buildInfo,
	}
}

// HTTPMiddleware returns a gorilla/mux MiddlewareFunc that records HTTP
// telemetry. It tracks in-flight requests, request counts partitioned by
// method, route template, and status code, and latency distributions for
// non-streaming endpoints.
func (m *Metrics) HTTPMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.httpRequestsInFlight.Inc()
			defer m.httpRequestsInFlight.Dec()

			start := time.Now()
			route := extractRoutePattern(r)
			isStream := isStreamingRequest(r, route)

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			codeStr := strconv.Itoa(rw.statusCode)

			m.httpRequestsTotal.WithLabelValues(r.Method, route, codeStr).Inc()

			// Exclude indefinite-lived streaming connections (SSE, WebSockets) from
			// duration histogram to prevent histogram distortion.
			if !isStream {
				m.httpRequestDuration.WithLabelValues(r.Method, route).Observe(duration.Seconds())
			}
		})
	}
}

// StartServer spins up a dedicated opt-in HTTP listener serving the Prometheus
// /metrics endpoint on the specified port.
func StartServer(ctx context.Context, port int, m *Metrics, log logger.Handler) (*http.Server, error) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if log != nil {
			log.Info(fmt.Sprintf("Prometheus metrics server listening on :%d", port))
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if log != nil {
				log.Error(fmt.Errorf("metrics server listener error: %w", err))
			}
		}
	}()

	return srv, nil
}

// extractRoutePattern safely retrieves the gorilla/mux template pattern
// to maintain low label cardinality. Falls back to "unmatched".
func extractRoutePattern(r *http.Request) string {
	if currentRoute := mux.CurrentRoute(r); currentRoute != nil {
		if pathTpl, err := currentRoute.GetPathTemplate(); err == nil && pathTpl != "" {
			return pathTpl
		}
		if pathRegex, err := currentRoute.GetPathRegexp(); err == nil && pathRegex != "" {
			return pathRegex
		}
	}
	return unmatchedRoute
}

// isStreamingRequest checks if the request represents a persistent stream
// such as Server-Sent Events or WebSockets.
func isStreamingRequest(r *http.Request, route string) bool {
	if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return true
	}
	if strings.Contains(strings.ToLower(r.Header.Get("Accept")), "text/event-stream") {
		return true
	}
	if strings.HasSuffix(route, "/stream") {
		return true
	}
	return false
}
