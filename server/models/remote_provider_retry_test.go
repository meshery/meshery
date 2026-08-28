package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestLoadCapabilities_SucceedsOnFirstAttempt(t *testing.T) {
	var callCount atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"providerURL":"http://` + r.Host + `"}`))
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	viper.Set("BUILD", "v1.2.3")
	viper.Set("OS", "linux")
	viper.Set("PLAYGROUND", "")
	t.Cleanup(func() {
		viper.Set("BUILD", "")
		viper.Set("OS", "")
		viper.Set("PLAYGROUND", "")
	})

	props, err := provider.loadCapabilities(context.Background(), "")
	if err != nil {
		t.Fatalf("loadCapabilities failed: %v", err)
	}
	if callCount.Load() != 1 {
		t.Fatalf("expected 1 request, got %d", callCount.Load())
	}
	if props.ProviderURL != server.URL {
		t.Errorf("expected ProviderURL %q, got %q", server.URL, props.ProviderURL)
	}
}

func TestLoadCapabilities_RetriesOnConnectionErrors(t *testing.T) {
	var callCount atomic.Int64

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := callCount.Add(1)
		if count < 3 {
			hj, ok := w.(http.Hijacker)
			if !ok {
				t.Fatal("server does not support hijacking")
			}
			conn, _, err := hj.Hijack()
			if err != nil {
				t.Fatalf("hijack failed: %v", err)
			}
			_ = conn.Close()
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"providerURL":"http://` + r.Host + `"}`))
	}))
	server.Start()
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	viper.Set("BUILD", "v1.2.3")
	viper.Set("OS", "linux")
	viper.Set("PLAYGROUND", "")
	t.Cleanup(func() {
		viper.Set("BUILD", "")
		viper.Set("OS", "")
		viper.Set("PLAYGROUND", "")
	})

	start := time.Now()
	props, err := provider.loadCapabilities(context.Background(), "")
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("loadCapabilities failed after retries: %v", err)
	}
	if callCount.Load() != 3 {
		t.Fatalf("expected 3 requests (2 connection failures + 1 success), got %d", callCount.Load())
	}
	if elapsed > 10*time.Second {
		t.Fatalf("retries took too long: %v (expected ~3-4s with exponential backoff)", elapsed)
	}
	if props.ProviderURL != server.URL {
		t.Errorf("expected ProviderURL %q, got %q", server.URL, props.ProviderURL)
	}
}

func TestLoadCapabilities_ExhaustsRetries(t *testing.T) {
	var callCount atomic.Int64

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("server does not support hijacking")
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			t.Fatalf("hijack failed: %v", err)
		}
		_ = conn.Close()
	}))
	server.Start()
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	viper.Set("BUILD", "v1.2.3")
	viper.Set("OS", "linux")
	viper.Set("PLAYGROUND", "")
	t.Cleanup(func() {
		viper.Set("BUILD", "")
		viper.Set("OS", "")
		viper.Set("PLAYGROUND", "")
	})

	start := time.Now()
	_, err := provider.loadCapabilities(context.Background(), "")
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected error when all retries exhausted, got nil")
	}
	attempts := callCount.Load()
	if attempts == 0 {
		t.Fatal("expected at least 1 request, got 0")
	}
	t.Logf("retries exhausted after %d attempts in %v (backoff: 1s/2s/4s/8s, max elapsed 15s)", attempts, elapsed)
}

func TestLoadCapabilities_ContextCancelled(t *testing.T) {
	var callCount atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	viper.Set("BUILD", "v1.2.3")
	viper.Set("OS", "linux")
	viper.Set("PLAYGROUND", "")
	t.Cleanup(func() {
		viper.Set("BUILD", "")
		viper.Set("OS", "")
		viper.Set("PLAYGROUND", "")
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := provider.loadCapabilities(ctx, "")
	if err == nil {
		t.Fatal("expected error when context is cancelled, got nil")
	}
}

func TestLoadCapabilities_WithToken(t *testing.T) {
	var callCount atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"providerURL":"http://` + r.Host + `"}`))
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	viper.Set("BUILD", "v1.2.3")
	viper.Set("OS", "linux")
	viper.Set("PLAYGROUND", "")
	t.Cleanup(func() {
		viper.Set("BUILD", "")
		viper.Set("OS", "")
		viper.Set("PLAYGROUND", "")
	})

	props, err := provider.loadCapabilities(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("loadCapabilities with token failed: %v", err)
	}
	if callCount.Load() != 1 {
		t.Fatalf("expected 1 request, got %d", callCount.Load())
	}
	if props.ProviderURL != server.URL {
		t.Errorf("expected ProviderURL %q, got %q", server.URL, props.ProviderURL)
	}
}

func TestLoadCapabilities_Non200StatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	viper.Set("BUILD", "v1.2.3")
	viper.Set("OS", "linux")
	viper.Set("PLAYGROUND", "")
	t.Cleanup(func() {
		viper.Set("BUILD", "")
		viper.Set("OS", "")
		viper.Set("PLAYGROUND", "")
	})

	_, err := provider.loadCapabilities(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for non-200 status code, got nil")
	}
}
