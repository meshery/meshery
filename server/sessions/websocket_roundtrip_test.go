package sessions

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/meshery/meshkit/logger"
)

// wsURL rewrites an httptest http:// URL into a ws:// one.
func wsURL(t *testing.T, serverURL string) string {
	t.Helper()
	parsed, err := url.Parse(serverURL)
	if err != nil {
		t.Fatalf("url.Parse(%q): %v", serverURL, err)
	}
	parsed.Scheme = "ws"
	return parsed.String()
}

func testLogger(t *testing.T) logger.Handler {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("logger.New: %v", err)
	}
	return log
}

// TestSocketRoundTrip drives the transport end to end over a real WebSocket:
// the ready frame, stdin as binary frames, stdout back as binary frames, a
// resize control frame reaching the driver, and the terminal exit frame
// followed by the close code.
func TestSocketRoundTrip(t *testing.T) {
	log := testLogger(t)

	// Reports what the "driver" side observed, read only after the socket closes.
	stdinSeen := make(chan string, 1)
	sizeSeen := make(chan TerminalSize, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := Upgrade(w, r, log)
		if err != nil {
			t.Errorf("Upgrade: %v", err)
			return
		}
		ctx := socket.Start(r.Context())
		defer socket.Close(CloseSessionEnded, "done")

		if err := socket.SendControl(ControlMessage{
			Type:         ControlReady,
			Capabilities: &Capabilities{Terminal: true, DefaultContainer: "app"},
		}); err != nil {
			t.Errorf("SendControl(ready): %v", err)
			return
		}

		buf := make([]byte, len("hello"))
		if _, err := io.ReadFull(socket.Stdin(), buf); err != nil {
			t.Errorf("read stdin: %v", err)
			return
		}
		stdinSeen <- string(buf)

		if _, err := socket.Stdout().Write([]byte("echo:" + string(buf))); err != nil {
			t.Errorf("write stdout: %v", err)
			return
		}

		select {
		case size := <-socket.Resize():
			sizeSeen <- size
		case <-ctx.Done():
			t.Error("context cancelled before a resize arrived")
			return
		case <-time.After(5 * time.Second):
			t.Error("timed out waiting for a resize")
			return
		}

		if err := socket.SendControl(ControlMessage{Type: ControlExit, ExitCode: 7}); err != nil {
			t.Errorf("SendControl(exit): %v", err)
		}
	}))
	defer server.Close()

	client, _, err := websocket.DefaultDialer.Dial(wsURL(t, server.URL), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer client.Close()

	// 1. The ready frame arrives first, carrying resolved capabilities.
	messageType, payload, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("read ready: %v", err)
	}
	if messageType != websocket.TextMessage {
		t.Fatalf("ready frame type = %d, want TextMessage", messageType)
	}
	if !strings.Contains(string(payload), `"type":"ready"`) ||
		!strings.Contains(string(payload), `"defaultContainer":"app"`) {
		t.Fatalf("ready frame = %s, want a ready frame carrying capabilities", payload)
	}

	// 2. stdin travels as a binary frame.
	if err := client.WriteMessage(websocket.BinaryMessage, []byte("hello")); err != nil {
		t.Fatalf("write stdin: %v", err)
	}

	// 3. stdout comes back as a binary frame.
	messageType, payload, err = client.ReadMessage()
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	if messageType != websocket.BinaryMessage {
		t.Fatalf("stdout frame type = %d, want BinaryMessage", messageType)
	}
	if string(payload) != "echo:hello" {
		t.Fatalf("stdout = %q, want %q", payload, "echo:hello")
	}

	// 4. A resize control frame reaches the driver.
	if err := client.WriteMessage(websocket.TextMessage, []byte(`{"type":"resize","cols":120,"rows":40}`)); err != nil {
		t.Fatalf("write resize: %v", err)
	}

	// 5. The session ends with an exit frame, then a close frame carrying the
	//    session-ended code.
	messageType, payload, err = client.ReadMessage()
	if err != nil {
		t.Fatalf("read exit: %v", err)
	}
	if messageType != websocket.TextMessage || !strings.Contains(string(payload), `"exitCode":7`) {
		t.Fatalf("final frame = %s, want an exit frame with code 7", payload)
	}

	if _, _, err = client.ReadMessage(); err == nil {
		t.Fatal("expected the socket to close after the exit frame")
	} else if !websocket.IsCloseError(err, CloseSessionEnded) {
		t.Fatalf("close error = %v, want close code %d", err, CloseSessionEnded)
	}

	if got := <-stdinSeen; got != "hello" {
		t.Errorf("driver saw stdin %q, want %q", got, "hello")
	}
	if got := <-sizeSeen; got.Width != 120 || got.Height != 40 {
		t.Errorf("driver saw geometry %+v, want 120x40", got)
	}
}

// TestSocketStdinEOFOnClientDisconnect: a driver blocked reading stdin must be
// released when the browser goes away, or the remote exec would hang forever.
func TestSocketStdinEOFOnClientDisconnect(t *testing.T) {
	log := testLogger(t)
	stdinClosed := make(chan error, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := Upgrade(w, r, log)
		if err != nil {
			t.Errorf("Upgrade: %v", err)
			return
		}
		socket.Start(r.Context())
		_, err = io.ReadAll(socket.Stdin())
		stdinClosed <- err
	}))
	defer server.Close()

	client, _, err := websocket.DefaultDialer.Dial(wsURL(t, server.URL), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	client.Close()

	select {
	case err := <-stdinClosed:
		if err != nil {
			t.Errorf("reading stdin after disconnect = %v, want a clean EOF", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("stdin never reported EOF after the client disconnected")
	}
}

// TestSocketContextCancelledOnDisconnect: the context handed to a driver must
// fire when the client vanishes, so the upstream stream is torn down promptly.
func TestSocketContextCancelledOnDisconnect(t *testing.T) {
	log := testLogger(t)
	cancelled := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := Upgrade(w, r, log)
		if err != nil {
			t.Errorf("Upgrade: %v", err)
			return
		}
		ctx := socket.Start(context.Background())
		<-ctx.Done()
		close(cancelled)
	}))
	defer server.Close()

	client, _, err := websocket.DefaultDialer.Dial(wsURL(t, server.URL), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	client.Close()

	select {
	case <-cancelled:
	case <-time.After(5 * time.Second):
		t.Fatal("the driver's context was never cancelled after the client disconnected")
	}
}

// TestUpgradeRejectsCrossSiteOrigin is the handshake-level counterpart to
// TestOriginAllowed: a page on another origin must not be able to open a
// session using the victim's ambient cookie.
func TestUpgradeRejectsCrossSiteOrigin(t *testing.T) {
	log := testLogger(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := Upgrade(w, r, log); err == nil {
			t.Error("Upgrade accepted a cross-site origin")
		}
	}))
	defer server.Close()

	_, response, err := websocket.DefaultDialer.Dial(wsURL(t, server.URL), http.Header{
		"Origin": []string{"https://evil.example.com"},
	})
	if err == nil {
		t.Fatal("dial from a cross-site origin succeeded")
	}
	if response == nil || response.StatusCode != http.StatusForbidden {
		t.Fatalf("cross-site handshake status = %v, want 403", response)
	}
}

// A same-origin handshake, which is what the browser sends, must succeed.
func TestUpgradeAcceptsSameOrigin(t *testing.T) {
	log := testLogger(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := Upgrade(w, r, log)
		if err != nil {
			t.Errorf("Upgrade rejected a same-origin handshake: %v", err)
			return
		}
		socket.Start(r.Context())
		socket.Close(CloseSessionEnded, "done")
	}))
	defer server.Close()

	client, _, err := websocket.DefaultDialer.Dial(wsURL(t, server.URL), http.Header{
		"Origin": []string{server.URL},
	})
	if err != nil {
		t.Fatalf("same-origin dial failed: %v", err)
	}
	client.Close()
}
