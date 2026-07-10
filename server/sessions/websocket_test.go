package sessions

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/meshery/meshkit/logger"
)

func TestOriginAllowed(t *testing.T) {
	allowed := map[string]struct{}{
		"https://cloud.example.com": {},
	}

	tests := []struct {
		name   string
		host   string
		origin string
		want   bool
	}{
		{
			name:   "no origin is a non-browser client",
			host:   "localhost:9081",
			origin: "",
			want:   true,
		},
		{
			name:   "same origin",
			host:   "localhost:9081",
			origin: "http://localhost:9081",
			want:   true,
		},
		{
			name:   "same origin over https",
			host:   "meshery.example.com",
			origin: "https://meshery.example.com",
			want:   true,
		},
		{
			name:   "same origin through the UI dev server proxy",
			host:   "localhost:3000",
			origin: "http://localhost:3000",
			want:   true,
		},
		{
			name:   "host comparison is case insensitive",
			host:   "Meshery.Example.com",
			origin: "https://meshery.example.com",
			want:   true,
		},
		{
			name:   "explicitly allowlisted cross origin",
			host:   "meshery.example.com",
			origin: "https://cloud.example.com",
			want:   true,
		},
		{
			// The whole point of the check: an attacker's page must not be able
			// to open a shell using the victim's ambient session cookie.
			name:   "cross site origin is rejected",
			host:   "localhost:9081",
			origin: "https://evil.example.com",
			want:   false,
		},
		{
			name:   "port mismatch is a different origin",
			host:   "localhost:9081",
			origin: "http://localhost:3000",
			want:   false,
		},
		{
			name:   "subdomain of the host is a different origin",
			host:   "example.com",
			origin: "https://evil.example.com",
			want:   false,
		},
		{
			name:   "host as a prefix of the origin is a different origin",
			host:   "example.com",
			origin: "https://example.com.evil.net",
			want:   false,
		},
		{
			name:   "unparseable origin is rejected",
			host:   "localhost:9081",
			origin: "://not a url",
			want:   false,
		},
		{
			name:   "origin without a host is rejected",
			host:   "localhost:9081",
			origin: "null",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{Host: tt.host, Header: http.Header{}}
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if got := originAllowed(req, allowed); got != tt.want {
				t.Errorf("originAllowed(host=%q, origin=%q) = %v, want %v", tt.host, tt.origin, got, tt.want)
			}
		})
	}
}

// newTestSocket builds a Socket with no underlying connection. The queueing
// paths — enqueue, Stdout, Stdin — touch only the channels, so they can be
// exercised without a live WebSocket.
func newTestSocket(t *testing.T) *Socket {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("logger.New: %v", err)
	}
	return &Socket{
		log:           log,
		outbound:      make(chan frame, outboundBufferFrames),
		stdin:         make(chan []byte, stdinBufferFrames),
		resize:        make(chan TerminalSize, resizeBufferFrames),
		done:          make(chan struct{}),
		writePumpDone: make(chan struct{}),
	}
}

// TestBinaryWriterCopiesPayload guards the aliasing hazard that makes terminal
// output non-deterministic: client-go reuses the buffer it hands to Write, and
// frames sit queued until the write pump drains them.
func TestBinaryWriterCopiesPayload(t *testing.T) {
	socket := newTestSocket(t)
	writer := socket.Stdout()

	scratch := []byte("first")
	if _, err := writer.Write(scratch); err != nil {
		t.Fatalf("Write: %v", err)
	}

	// Simulate client-go reusing the buffer before the pump drains the queue.
	copy(scratch, []byte("SECON"))

	select {
	case f := <-socket.outbound:
		if f.messageType != websocket.BinaryMessage {
			t.Errorf("messageType = %d, want BinaryMessage", f.messageType)
		}
		if !bytes.Equal(f.payload, []byte("first")) {
			t.Errorf("queued payload = %q, want %q (writer aliased the caller's buffer)", f.payload, "first")
		}
	default:
		t.Fatal("no frame was queued")
	}
}

func TestBinaryWriterEmptyWriteQueuesNothing(t *testing.T) {
	socket := newTestSocket(t)

	n, err := socket.Stdout().Write(nil)
	if err != nil || n != 0 {
		t.Fatalf("Write(nil) = (%d, %v), want (0, nil)", n, err)
	}
	if len(socket.outbound) != 0 {
		t.Errorf("queued %d frames for an empty write, want 0", len(socket.outbound))
	}
}

func TestStdinReaderSpansFrames(t *testing.T) {
	socket := newTestSocket(t)
	socket.stdin <- []byte("hello")
	socket.stdin <- []byte(" world")
	close(socket.stdin)

	got, err := io.ReadAll(socket.Stdin())
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != "hello world" {
		t.Errorf("read %q, want %q", got, "hello world")
	}
}

// TestStdinReaderCarriesRemainder covers the case where the consumer's buffer
// is smaller than an inbound frame, e.g. a paste larger than the exec stream's
// read buffer.
func TestStdinReaderCarriesRemainder(t *testing.T) {
	socket := newTestSocket(t)
	socket.stdin <- []byte("abcdef")

	reader := socket.Stdin()
	buf := make([]byte, 4)

	n, err := reader.Read(buf)
	if err != nil || n != 4 || string(buf[:n]) != "abcd" {
		t.Fatalf("first Read = (%d, %q, %v), want (4, \"abcd\", nil)", n, buf[:n], err)
	}

	n, err = reader.Read(buf)
	if err != nil || n != 2 || string(buf[:n]) != "ef" {
		t.Fatalf("second Read = (%d, %q, %v), want (2, \"ef\", nil)", n, buf[:n], err)
	}
}

func TestStdinReaderReportsEOFOnShutdown(t *testing.T) {
	socket := newTestSocket(t)
	socket.shutdown()

	if _, err := socket.Stdin().Read(make([]byte, 8)); err != io.EOF {
		t.Errorf("Read after shutdown = %v, want io.EOF", err)
	}
}

func TestEnqueueAfterShutdownFails(t *testing.T) {
	socket := newTestSocket(t)
	socket.shutdown()

	if _, err := socket.Stdout().Write([]byte("data")); err == nil {
		t.Error("Write after shutdown succeeded, want an error")
	}
}

func TestShutdownIsIdempotent(t *testing.T) {
	socket := newTestSocket(t)
	socket.shutdown()
	socket.shutdown() // must not panic on a double close of done
}

// TestHandleControlResizeDropsOldest asserts that a burst of resize events
// cannot block the read pump, and that the newest geometry is the one that
// survives.
func TestHandleControlResizeDropsOldest(t *testing.T) {
	socket := newTestSocket(t)

	total := resizeBufferFrames + 3
	for i := 1; i <= total; i++ {
		socket.handleControl([]byte(`{"type":"resize","cols":` + strconv.Itoa(i) + `,"rows":24}`))
	}

	if len(socket.resize) != resizeBufferFrames {
		t.Fatalf("resize buffer holds %d, want it saturated at %d", len(socket.resize), resizeBufferFrames)
	}

	// Drain: the newest geometry must be present, the oldest must be gone.
	var last TerminalSize
	for len(socket.resize) > 0 {
		last = <-socket.resize
	}
	if last.Width != uint16(total) {
		t.Errorf("newest geometry = %d cols, want %d", last.Width, total)
	}
}

func TestHandleControlIgnoresJunk(t *testing.T) {
	socket := newTestSocket(t)

	// Malformed JSON, an unknown type, a server-to-client type, and a
	// degenerate geometry must all be ignored rather than fatal.
	socket.handleControl([]byte(`{`))
	socket.handleControl([]byte(`{"type":"teleport"}`))
	socket.handleControl([]byte(`{"type":"ready"}`))
	socket.handleControl([]byte(`{"type":"resize","cols":0,"rows":0}`))

	if len(socket.resize) != 0 {
		t.Errorf("queued %d resize events, want 0", len(socket.resize))
	}
}
