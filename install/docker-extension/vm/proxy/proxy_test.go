package proxy_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
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
	expectMsg(proxy.UnauthenticatedMsg) // current state on connect

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

// dialWS opens a websocket to the proxy's /ws endpoint and checks the two
// messages every connection starts with: the handshake, then the current auth
// state (wantState).
func dialWS(t *testing.T, srvURL, wantState string) *websocket.Conn {
	t.Helper()
	ws, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srvURL, "http")+"/ws", nil)
	if err != nil {
		t.Fatalf("dial /ws: %v", err)
	}
	t.Cleanup(func() { ws.Close() })
	_ = ws.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, msg, err := ws.ReadMessage(); err != nil || string(msg) != proxy.InitiationSuccessful {
		t.Fatalf("expected %q, got %q (err %v)", proxy.InitiationSuccessful, msg, err)
	}
	if _, msg, err := ws.ReadMessage(); err != nil || string(msg) != wantState {
		t.Fatalf("expected current state %q on connect, got %q (err %v)", wantState, msg, err)
	}
	return ws
}

// When the UI reconnects while its old websocket is still open, every auth
// change must reach the newest connection instead of being split between the
// two listeners.
func TestServeHTTP_NewestWSClientGetsAllAuthChanges(t *testing.T) {
	srv := httptest.NewServer(&proxy.Proxy{})
	defer srv.Close()

	older := dialWS(t, srv.URL, proxy.UnauthenticatedMsg)
	newer := dialWS(t, srv.URL, proxy.UnauthenticatedMsg)

	for i := 0; i < 4; i++ {
		req, _ := http.NewRequest(http.MethodGet, srv.URL+"/token/store?token=abc", nil)
		want := proxy.AuthenticatedMsg
		if i%2 == 1 {
			req, _ = http.NewRequest(http.MethodDelete, srv.URL+"/token", nil)
			want = proxy.UnauthenticatedMsg
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("%s %s: %v", req.Method, req.URL.Path, err)
		}
		resp.Body.Close()

		_ = newer.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, msg, err := newer.ReadMessage()
		if err != nil {
			t.Fatalf("change %d: newest client got no notification: %v", i, err)
		}
		if string(msg) != want {
			t.Fatalf("change %d: expected %q, got %q", i, want, msg)
		}
	}

	_ = older.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	if _, msg, err := older.ReadMessage(); err == nil {
		t.Fatalf("replaced client should get no auth notifications, got %q", msg)
	}
}

// The auth token is a credential, so neither /token/store nor a proxied
// request carrying the token cookie may write it to the log.
func TestServeHTTP_TokenNotLogged(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer backend.Close()
	old := proxy.MesheryServerHost
	proxy.MesheryServerHost = strings.TrimPrefix(backend.URL, "http://")
	t.Cleanup(func() { proxy.MesheryServerHost = old })

	buf := captureLog(t)

	const secret = "s3cr3t-token-value"
	p := &proxy.Proxy{}
	serveWithTimeout(t, p, httptest.NewRequest(http.MethodGet, "http://localhost/token/store?token="+secret, nil))
	serveWithTimeout(t, p, httptest.NewRequest(http.MethodGet, "http://localhost/api/system/version", nil))

	if strings.Contains(buf.String(), secret) {
		t.Fatalf("auth token leaked into logs:\n%s", buf.String())
	}
	if !strings.Contains(buf.String(), "token=REDACTED") {
		t.Fatalf("expected the /token/store request to be logged with a redacted token:\n%s", buf.String())
	}
}

// syncBuffer is a bytes.Buffer safe for concurrent use. The global logger can
// be written by websocket goroutines left over from earlier tests while a test
// reads what it captured.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// captureLog redirects the global logger to a buffer for the rest of the test.
func captureLog(t *testing.T) *syncBuffer {
	t.Helper()
	buf := &syncBuffer{}
	log.SetOutput(buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	return buf
}

// When the backend is unreachable, the client error embeds the request URL,
// so a token in the query must not reach the log through it.
func TestServeHTTP_TokenNotLoggedOnBackendError(t *testing.T) {
	old := proxy.MesheryServerHost
	proxy.MesheryServerHost = "127.0.0.1:1" // connection refused
	t.Cleanup(func() { proxy.MesheryServerHost = old })

	buf := captureLog(t)

	const secret = "s3cr3t-token-value"
	rec := serveWithTimeout(t, &proxy.Proxy{}, httptest.NewRequest(http.MethodGet, "http://localhost/api/system/version?token="+secret, nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if strings.Contains(buf.String(), secret) {
		t.Fatalf("auth token leaked into logs:\n%s", buf.String())
	}
	if !strings.Contains(buf.String(), "backend request failed") {
		t.Fatalf("expected the backend failure to be logged:\n%s", buf.String())
	}
}

// A client is told the current auth state when it connects, so a login or
// logout that lands between the UI reading /token and opening /ws is not lost.
func TestServeHTTP_WSClientGetsCurrentStateOnConnect(t *testing.T) {
	srv := httptest.NewServer(&proxy.Proxy{})
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/token/store?token=abc")
	if err != nil {
		t.Fatalf("store token: %v", err)
	}
	resp.Body.Close()
	dialWS(t, srv.URL, proxy.AuthenticatedMsg)

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/token", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("delete token: %v", err)
	}
	resp.Body.Close()
	dialWS(t, srv.URL, proxy.UnauthenticatedMsg)
}

// A client that sends messages gets them echoed back, and those echoes must
// not race the handshake or auth-state writes on the same connection: gorilla
// websocket panics on concurrent writes, which would take down the proxy.
func TestServeHTTP_EchoAndAuthWritesDoNotRace(t *testing.T) {
	srv := httptest.NewServer(&proxy.Proxy{})
	defer srv.Close()

	ws := dialWS(t, srv.URL, proxy.UnauthenticatedMsg)
	// dialWS leaves a short read deadline on the connection; clear it so a
	// slow run doesn't stop the drainer and stall the server's echo writes.
	_ = ws.SetReadDeadline(time.Time{})

	// Drain everything the server sends so its writes never block.
	go func() {
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				return
			}
		}
	}()

	stop := make(chan struct{})
	sendDone := make(chan struct{})
	go func() {
		defer close(sendDone)
		for {
			select {
			case <-stop:
				return
			default:
			}
			if err := ws.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				return
			}
		}
	}()

	for i := 0; i < 200; i++ {
		req, _ := http.NewRequest(http.MethodGet, srv.URL+"/token/store?token=abc", nil)
		if i%2 == 1 {
			req, _ = http.NewRequest(http.MethodDelete, srv.URL+"/token", nil)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("%s %s: %v", req.Method, req.URL.Path, err)
		}
		resp.Body.Close()
	}
	close(stop)
	ws.Close() // unblock a sender stuck in WriteMessage
	<-sendDone
}

// The client error for a malformed redirect quotes the Location header, which
// can carry a token, so the error text itself must be redacted too.
func TestServeHTTP_TokenNotLoggedOnMalformedRedirect(t *testing.T) {
	const secret = "s3cr3t-token-value"
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", "/next%zz?token="+secret)
		w.WriteHeader(http.StatusFound)
	}))
	defer backend.Close()
	old := proxy.MesheryServerHost
	proxy.MesheryServerHost = strings.TrimPrefix(backend.URL, "http://")
	t.Cleanup(func() { proxy.MesheryServerHost = old })

	buf := captureLog(t)
	rec := serveWithTimeout(t, &proxy.Proxy{}, httptest.NewRequest(http.MethodGet, "http://localhost/api/system/version", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if strings.Contains(buf.String(), secret) {
		t.Fatalf("auth token leaked into logs:\n%s", buf.String())
	}
	if !strings.Contains(buf.String(), "?REDACTED") {
		t.Fatalf("expected the redirect error to be logged with its query redacted:\n%s", buf.String())
	}
}
