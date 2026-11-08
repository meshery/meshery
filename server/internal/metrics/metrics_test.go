package metrics

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
)

type mockFlusherHijacker struct {
	http.ResponseWriter
	flushed  bool
	hijacked bool
}

func (m *mockFlusherHijacker) Flush() {
	m.flushed = true
}

func (m *mockFlusherHijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	m.hijacked = true
	serverConn, clientConn := net.Pipe()
	_ = clientConn.Close()
	return serverConn, bufio.NewReadWriter(bufio.NewReader(serverConn), bufio.NewWriter(serverConn)), nil
}

func TestNewMetricsInitialization(t *testing.T) {
	version := "v0.7.50"
	commitSHA := "abcdef123456"
	channel := "edge"

	m := New(version, commitSHA, channel)
	if m == nil {
		t.Fatal("expected non-nil Metrics instance")
	}

	gauge, err := m.buildInfo.GetMetricWithLabelValues(version, commitSHA, channel)
	if err != nil {
		t.Fatalf("unexpected error getting build_info metric: %v", err)
	}

	val := testutil.ToFloat64(gauge)
	if val != 1.0 {
		t.Fatalf("expected build_info value 1.0, got %f", val)
	}
}

func TestHTTPMiddlewareRoutingAndLatency(t *testing.T) {
	m := New("test", "test", "test")
	r := mux.NewRouter()
	r.Use(m.HTTPMiddleware())

	r.HandleFunc("/api/user/{id}", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("user-response"))
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/user/123", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// Verify counter incremented for route template "/api/user/{id}"
	counter, err := m.httpRequestsTotal.GetMetricWithLabelValues("GET", "/api/user/{id}", "200")
	if err != nil {
		t.Fatalf("unexpected error fetching counter: %v", err)
	}
	if count := testutil.ToFloat64(counter); count != 1.0 {
		t.Fatalf("expected request count 1.0, got %f", count)
	}

	// Verify duration histogram was observed
	observer, err := m.httpRequestDuration.GetMetricWithLabelValues("GET", "/api/user/{id}")
	if err != nil {
		t.Fatalf("unexpected error fetching duration observer: %v", err)
	}
	promMetric, ok := observer.(prometheus.Metric)
	if !ok {
		t.Fatal("observer does not implement prometheus.Metric")
	}
	pbMetric := &dto.Metric{}
	if err := promMetric.Write(pbMetric); err != nil {
		t.Fatalf("failed to write metric proto: %v", err)
	}
	if pbMetric.GetHistogram().GetSampleCount() != 1 {
		t.Fatalf("expected 1 sample in histogram, got %d", pbMetric.GetHistogram().GetSampleCount())
	}
}

func TestHTTPMiddlewareStreamingExclusion(t *testing.T) {
	m := New("test", "test", "test")
	r := mux.NewRouter()
	r.Use(m.HTTPMiddleware())

	r.HandleFunc("/api/providers/stream", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("data: event\n\n"))
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/providers/stream", nil)
	req.Header.Set("Accept", "text/event-stream")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	// Counter MUST be incremented
	counter, err := m.httpRequestsTotal.GetMetricWithLabelValues("GET", "/api/providers/stream", "200")
	if err != nil {
		t.Fatalf("unexpected error fetching counter: %v", err)
	}
	if count := testutil.ToFloat64(counter); count != 1.0 {
		t.Fatalf("expected request count 1.0 for stream, got %f", count)
	}

	// Latency histogram MUST NOT be observed for streaming connections
	observer, err := m.httpRequestDuration.GetMetricWithLabelValues("GET", "/api/providers/stream")
	if err != nil {
		t.Fatalf("unexpected error fetching duration observer: %v", err)
	}
	promMetric, ok := observer.(prometheus.Metric)
	if !ok {
		t.Fatal("observer does not implement prometheus.Metric")
	}
	pbMetric := &dto.Metric{}
	if err := promMetric.Write(pbMetric); err != nil {
		t.Fatalf("failed to write metric proto: %v", err)
	}
	if pbMetric.GetHistogram().GetSampleCount() != 0 {
		t.Fatalf("expected 0 samples in streaming duration histogram, got %d", pbMetric.GetHistogram().GetSampleCount())
	}
}

func TestHTTPMiddlewareUnmatchedRoute(t *testing.T) {
	m := New("test", "test", "test")
	r := mux.NewRouter()
	r.Use(m.HTTPMiddleware())
	r.NotFoundHandler = m.HTTPMiddleware()(http.NotFoundHandler())

	req := httptest.NewRequest("GET", "/non-existent-endpoint", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	// Route should be recorded as "unmatched"
	counter, err := m.httpRequestsTotal.GetMetricWithLabelValues("GET", "unmatched", "404")
	if err != nil {
		t.Fatalf("unexpected error fetching counter: %v", err)
	}
	if count := testutil.ToFloat64(counter); count != 1.0 {
		t.Fatalf("expected request count 1.0 for unmatched, got %f", count)
	}
}

func TestResponseWriterInterfaces(t *testing.T) {
	rec := httptest.NewRecorder()
	mock := &mockFlusherHijacker{ResponseWriter: rec}
	rw := newResponseWriter(mock)

	// Test Flush passthrough
	rw.Flush()
	if !mock.flushed {
		t.Fatal("expected Flush() to be propagated to underlying writer")
	}

	// Test Hijack passthrough
	conn, _, err := rw.Hijack()
	if err != nil {
		t.Fatalf("unexpected hijack error: %v", err)
	}
	_ = conn.Close()
	if !mock.hijacked {
		t.Fatal("expected Hijack() to be propagated to underlying writer")
	}

	// Test Hijack error on standard ResponseRecorder (non-hijacker)
	stdRW := newResponseWriter(httptest.NewRecorder())
	_, _, err = stdRW.Hijack()
	if err == nil {
		t.Fatal("expected error when underlying writer does not implement Hijacker")
	}
}

func TestStartServerAndEndpoints(t *testing.T) {
	m := New("v0.7.0", "123456", "stable")

	// Trigger one request through middleware so http_requests_total is populated
	r := mux.NewRouter()
	r.Use(m.HTTPMiddleware())
	r.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/healthz", nil))

	// Pick free port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on ephemeral port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()

	srv, err := StartServer(context.Background(), port, m, nil)
	if err != nil {
		t.Fatalf("unexpected error starting metrics server: %v", err)
	}
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	var resp *http.Response
	for i := 0; i < 20; i++ {
		time.Sleep(25 * time.Millisecond)
		resp, err = http.Get(baseURL + "/healthz")
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("health check failed on metrics server: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on /healthz, got %d", resp.StatusCode)
	}

	// Test /metrics endpoint
	metricsResp, err := http.Get(baseURL + "/metrics")
	if err != nil {
		t.Fatalf("failed to query /metrics: %v", err)
	}
	defer func() { _ = metricsResp.Body.Close() }()

	if metricsResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on /metrics, got %d", metricsResp.StatusCode)
	}

	body, err := io.ReadAll(metricsResp.Body)
	if err != nil {
		t.Fatalf("failed to read /metrics body: %v", err)
	}

	bodyStr := string(body)
	if !strings.Contains(bodyStr, "meshery_build_info") {
		t.Fatal("expected /metrics to expose meshery_build_info")
	}
	if !strings.Contains(bodyStr, "meshery_http_requests_total") {
		t.Fatal("expected /metrics to expose meshery_http_requests_total")
	}
}
