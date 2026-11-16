package proxy_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/docker/meshery-extension/vm/proxy"
	"github.com/gorilla/websocket"
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
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

// serveWithTimeout runs ServeHTTP and fails the test if it does not return
// within a second, which is how a blocked channel send shows up.
func serveWithTimeout(t *testing.T, p *proxy.Proxy, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	done := make(chan struct{})
	go func() {
		defer close(done)
		p.ServeHTTP(rec, req)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("ServeHTTP(%s %s) did not return; it is blocked", req.Method, req.URL.Path)
	}
	return rec
}

// Storing a token must not hang when no /ws client has connected yet, which is
// the case when the browser completes login before the extension UI opens the
// websocket.
func TestServeHTTP_TokenStoreWithoutWSClient(t *testing.T) {
	p := &proxy.Proxy{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost/token/store?token=abc", nil)

	rec := serveWithTimeout(t, p, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	rec = serveWithTimeout(t, p, httptest.NewRequest(http.MethodGet, "http://localhost/token", nil))
	if got := rec.Body.String(); got != "abc" {
		t.Fatalf("expected stored token %q, got %q", "abc", got)
	}
}

// Logging out must not hang when no /ws client is connected.
func TestServeHTTP_TokenDeleteWithoutWSClient(t *testing.T) {
	p := &proxy.Proxy{}
	serveWithTimeout(t, p, httptest.NewRequest(http.MethodGet, "http://localhost/token/store?token=abc", nil))

	serveWithTimeout(t, p, httptest.NewRequest(http.MethodDelete, "http://localhost/token", nil))

	rec := serveWithTimeout(t, p, httptest.NewRequest(http.MethodGet, "http://localhost/token", nil))
	if got := rec.Body.String(); got != "null" {
		t.Fatalf("expected token to be cleared, got %q", got)
	}
}

// A connected /ws client is still told when the auth state changes.
func TestServeHTTP_WSClientNotifiedOfAuthChanges(t *testing.T) {
	srv := httptest.NewServer(&proxy.Proxy{})
	defer srv.Close()

	ws, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http")+"/ws", nil)
	if err != nil {
		t.Fatalf("dial /ws: %v", err)
	}
	defer ws.Close()

	expectMsg := func(want string) {
		t.Helper()
		_ = ws.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, msg, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("waiting for %q: %v", want, err)
		}
		if string(msg) != want {
			t.Fatalf("expected %q, got %q", want, msg)
		}
	}
	expectMsg(proxy.InitiationSuccessful)

	resp, err := http.Get(srv.URL + "/token/store?token=abc")
	if err != nil {
		t.Fatalf("store token: %v", err)
	}
	resp.Body.Close()
	expectMsg(proxy.AuthenticatedMsg)

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/token", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("delete token: %v", err)
	}
	resp.Body.Close()
	expectMsg(proxy.UnauthenticatedMsg)
}
