package proxy_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/meshery-extension/vm/proxy"
)

// runRecover runs f and reports whether it panicked.
func runRecover(f func()) (panicked bool, val interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			val = r
		}
	}()
	f()
	return
}

// When the Meshery server is unreachable, client.Do returns a non-nil error and
// a nil response. ServeHTTP must write a clean 500 and return, not dereference
// the nil response in `defer resp.Body.Close()`.
func TestServeHTTP_BackendUnreachable(t *testing.T) {
	old := proxy.MesheryServerHost
	proxy.MesheryServerHost = "127.0.0.1:1" // connection refused
	t.Cleanup(func() { proxy.MesheryServerHost = old })
	p := &proxy.Proxy{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/system/version", nil)
	rec := httptest.NewRecorder()

	panicked, val := runRecover(func() { p.ServeHTTP(rec, req) })
	if panicked {
		t.Fatalf("ServeHTTP panicked on unreachable backend: %v", val)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

// A plain (non-websocket) GET to /ws cannot be upgraded, so upgrader.Upgrade
// returns a nil connection and an error. ServeHTTP must check the error and
// return before calling ws.SetCloseHandler / handleWsMessage on the nil conn.
func TestServeHTTP_FailedWSUpgrade(t *testing.T) {
	p := &proxy.Proxy{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost/ws", nil) // no upgrade headers
	rec := httptest.NewRecorder()

	panicked, val := runRecover(func() { p.ServeHTTP(rec, req) })
	if panicked {
		t.Fatalf("ServeHTTP panicked on failed /ws upgrade: %v", val)
	}
}
