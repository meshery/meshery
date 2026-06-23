package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/meshery/meshkit/logger"
)

func TestNewPrometheusClient_UsesHTTPTimeout(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("logger.New: %v", err)
	}

	p := NewPrometheusClient(&log)
	if p == nil {
		t.Fatalf("NewPrometheusClient returned nil")
	}
	if p.grafanaClient == nil {
		t.Fatalf("grafanaClient is nil")
	}
	if p.grafanaClient.httpClient == nil {
		t.Fatalf("grafanaClient.httpClient is nil")
	}

	if got := p.grafanaClient.httpClient.Timeout; got != defaultPrometheusClientTimeout {
		t.Fatalf("http client Timeout = %s; want %s", got, defaultPrometheusClientTimeout)
	}
}

func TestPrometheusClient_getAllNodes_UsesConfiguredHTTPClientTimeout(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("logger.New: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay longer than the client's Timeout to prove we actually honor it.
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":[]}`))
	}))
	defer srv.Close()

	p := NewPrometheusClientWithHTTPClient(&http.Client{Timeout: 50 * time.Millisecond}, &log)

	_, err = p.getAllNodes(context.Background(), srv.URL)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
}

func TestPrometheusClient_QueryRangeUsingClient_UsesConfiguredHTTPClientTimeout(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("logger.New: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[]}}`))
	}))
	defer srv.Close()

	p := NewPrometheusClientWithHTTPClient(&http.Client{Timeout: 50 * time.Millisecond}, &log)

	start := time.Unix(0, 0)
	end := start.Add(1 * time.Minute)
	_, err = p.QueryRangeUsingClient(context.Background(), srv.URL, "up", start, end, 15*time.Second)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
}
