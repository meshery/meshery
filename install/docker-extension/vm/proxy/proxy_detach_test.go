package proxy

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// A client that disappears without sending a close frame must still have its
// auth listener stopped, instead of leaving it to fail a write on every later
// auth change.
func TestWSListenerDetachesWhenClientVanishes(t *testing.T) {
	p := &Proxy{}
	srv := httptest.NewServer(p)
	defer srv.Close()

	ws, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http")+"/ws", nil)
	if err != nil {
		t.Fatalf("dial /ws: %v", err)
	}
	// Read the handshake and the initial auth state.
	_ = ws.SetReadDeadline(time.Now().Add(2 * time.Second))
	for i := 0; i < 2; i++ {
		if _, _, err := ws.ReadMessage(); err != nil {
			t.Fatalf("reading initial message %d: %v", i, err)
		}
	}

	p.mu.Lock()
	attached := p.authWsChan != nil
	p.mu.Unlock()
	if !attached {
		t.Fatal("expected the connection's auth channel to be installed")
	}

	// Drop the TCP connection without a close frame.
	_ = ws.UnderlyingConn().Close()

	deadline := time.Now().Add(2 * time.Second)
	for {
		p.mu.Lock()
		detached := p.authWsChan == nil
		p.mu.Unlock()
		if detached {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("auth channel still installed after the client vanished; its listener keeps running")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
