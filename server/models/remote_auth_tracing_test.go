package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/meshery/meshkit/tracing"
)

// The remote provider calls for token refresh, keys, revoke and introspect should
// carry the same trace context as the other outbound provider calls, so a trace
// that reaches the provider is not broken at these requests.
func TestRemoteAuthClient_PropagatesTraceContext(t *testing.T) {
	// InitTracer is what the server uses at startup: it installs the global tracer
	// provider and the W3C trace context propagator. The exporter connects lazily,
	// so the collector endpoint is never contacted by this test.
	tp, err := tracing.InitTracer(context.Background(), tracing.Config{
		ServiceName: "remote-auth-test",
		Endpoint:    "127.0.0.1:1",
		Insecure:    true,
	})
	if err != nil {
		t.Fatalf("init tracer: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = tp.Shutdown(ctx)
	})

	gotTraceparent := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTraceparent <- r.Header.Get("traceparent")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	resp, err := remoteAuthClient.Get(srv.URL)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	_ = resp.Body.Close()

	if got := <-gotTraceparent; got == "" {
		t.Fatal("expected the request to carry a traceparent header, got none")
	}
}
