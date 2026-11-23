package console

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/meshery/meshkit/logger"
)

const (
	// writeWait bounds a single frame write, and also how long a stalled client
	// may block the producer before the console is torn down.
	writeWait = 10 * time.Second
	// pongWait is how long the server tolerates silence from the client. The
	// read deadline is pushed out by this much on every pong.
	pongWait = 60 * time.Second
	// pingPeriod must be meaningfully shorter than pongWait so a ping and its
	// pong can complete before the read deadline elapses.
	pingPeriod = (pongWait * 9) / 10
	// maxInboundMessageSize caps a single client frame. Terminal stdin is
	// keystrokes and pastes; 64 KiB is generous and bounds a hostile client.
	maxInboundMessageSize = 64 << 10
	// outboundBufferFrames absorbs bursts of output (a `cat` of a large file)
	// without blocking the copy from the remote target.
	outboundBufferFrames = 256
	// stdinBufferFrames absorbs a paste burst ahead of the remote's reads.
	stdinBufferFrames = 32
	// resizeBufferFrames is small because only the most recent geometry
	// matters; older ones are superseded and dropped.
	resizeBufferFrames = 4
)

// allowedOriginsEnv names a comma-separated list of extra origins permitted to
// open console sockets, for deployments that serve the UI from a different
// origin than the API. Same-origin requests are always allowed and need no
// configuration.
const allowedOriginsEnv = "MESHERY_ALLOWED_WS_ORIGINS"

// upgrader is shared: it holds no per-connection state.
var upgrader = newUpgrader(originsFromEnv())

func originsFromEnv() []string {
	raw := os.Getenv(allowedOriginsEnv)
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			origins = append(origins, strings.ToLower(p))
		}
	}
	return origins
}

func newUpgrader(allowedOrigins []string) *websocket.Upgrader {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = struct{}{}
	}
	return &websocket.Upgrader{
		HandshakeTimeout: 10 * time.Second,
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
		CheckOrigin:      func(r *http.Request) bool { return originAllowed(r, allowed) },
	}
}

// originAllowed enforces same-origin on the handshake, with an opt-in allowlist.
//
// This deliberately does not mirror the `return true` used by the GraphQL
// WebSocket transport. A console socket hands its holder an interactive shell
// inside the user's cluster, and the handshake is authenticated by an ambient
// session cookie, so accepting any Origin would be a cross-site WebSocket
// hijacking hole: any page the user visits could open a socket, authenticated
// as them, and exec into their workloads. The same-origin check is what stops
// that, because a browser cannot forge the Origin header.
func originAllowed(r *http.Request, allowed map[string]struct{}) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		// No Origin means a non-browser client (mesheryctl, tests, scripts).
		// A browser always sends Origin on a WebSocket handshake, so its
		// absence cannot be forged from a page. Such clients are authenticated
		// by the usual middleware; there is no ambient-credential risk to guard.
		return true
	}
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return false
	}
	if strings.EqualFold(u.Host, r.Host) {
		return true
	}
	_, ok := allowed[strings.ToLower(origin)]
	return ok
}

// frame is one queued outbound WebSocket message.
type frame struct {
	messageType int
	payload     []byte
}

// Socket is the server side of a console's WebSocket connection.
//
// It owns the underlying connection for the console's lifetime and enforces
// gorilla/websocket's contract that at most one goroutine writes at a time: all
// writes are funnelled through a single write pump, fed by a buffered channel.
// Callers therefore treat Stdout and SendControl as safe to use concurrently
// with the read side.
type Socket struct {
	ws  *websocket.Conn
	log logger.Handler

	outbound chan frame
	stdin    chan []byte
	resize   chan TerminalSize

	// done is closed exactly once, when either pump fails or Close is called.
	done      chan struct{}
	closeOnce sync.Once

	// pumpsDone is closed when the write pump has drained and exited, so Close
	// can wait for the close frame to actually reach the wire.
	writePumpDone chan struct{}
}

// Upgrade performs the WebSocket handshake. Any error has already been written
// to w as a plain HTTP response, so the caller only logs it.
func Upgrade(w http.ResponseWriter, r *http.Request, log logger.Handler) (*Socket, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, ErrUpgrade(err)
	}
	return &Socket{
		ws:            ws,
		log:           log,
		outbound:      make(chan frame, outboundBufferFrames),
		stdin:         make(chan []byte, stdinBufferFrames),
		resize:        make(chan TerminalSize, resizeBufferFrames),
		done:          make(chan struct{}),
		writePumpDone: make(chan struct{}),
	}, nil
}

// Start launches the read and write pumps and returns a context that is
// cancelled when the socket dies for any reason — client disconnect, missed
// pong, or a call to Close. A driver should thread that context into its
// streaming call so a vanished client tears the remote console down promptly.
func (s *Socket) Start(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-s.done
		cancel()
	}()
	go s.readPump()
	go s.writePump()

	// Cancelling the parent context must also stop the pumps, otherwise a
	// server shutdown would leave them blocked on the socket.
	go func() {
		<-ctx.Done()
		s.shutdown()
	}()

	return ctx
}

// readPump is the sole reader of the connection. It classifies inbound frames
// and translates them into stdin bytes and resize events.
func (s *Socket) readPump() {
	defer s.shutdown()
	defer close(s.stdin)
	defer close(s.resize)

	s.ws.SetReadLimit(maxInboundMessageSize)
	// A client that never speaks — a terminal sitting idle, a log stream with
	// no input — must still be reaped when its browser vanishes. The read
	// deadline plus the write pump's pings turn that into a bounded detection
	// time rather than a socket that leaks until the process restarts.
	_ = s.ws.SetReadDeadline(time.Now().Add(pongWait))
	s.ws.SetPongHandler(func(string) error {
		return s.ws.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		messageType, payload, err := s.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.log.Debug("console socket closed unexpectedly: " + err.Error())
			}
			return
		}

		switch messageType {
		case websocket.BinaryMessage:
			select {
			case s.stdin <- payload:
			case <-s.done:
				return
			}
		case websocket.TextMessage:
			s.handleControl(payload)
		}
	}
}

// handleControl decodes a client control frame. Unknown or wrong-direction
// types are ignored rather than fatal, so a newer client stays compatible with
// an older server.
func (s *Socket) handleControl(payload []byte) {
	var msg ControlMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		s.log.Debug("discarding malformed console control frame: " + err.Error())
		return
	}
	if msg.Type != ControlResize {
		return
	}
	// The schema bounds cols/rows to [0, 65535], but a hand-rolled client is not
	// obliged to honour it, and a silent wrap would resize the pty to something
	// absurd. Reject out-of-range geometry instead of truncating it.
	if msg.Cols <= 0 || msg.Rows <= 0 || msg.Cols > math.MaxUint16 || msg.Rows > math.MaxUint16 {
		return
	}
	size := TerminalSize{Width: uint16(msg.Cols), Height: uint16(msg.Rows)}
	// Only the newest geometry matters. If the consumer has not drained the
	// previous one, drop the oldest and retry, so a burst of resize events
	// during a window drag cannot block the read pump.
	for {
		select {
		case s.resize <- size:
			return
		case <-s.done:
			return
		default:
		}
		select {
		case <-s.resize:
		case <-s.done:
			return
		default:
			// The consumer drained it between the two selects; retry the send.
		}
	}
}

// writePump is the sole writer of the connection. It serializes application
// frames and keepalive pings onto the socket.
func (s *Socket) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		close(s.writePumpDone)
	}()

	for {
		select {
		case f, ok := <-s.outbound:
			if !ok {
				return
			}
			_ = s.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.ws.WriteMessage(f.messageType, f.payload); err != nil {
				s.log.Debug("console socket write failed: " + err.Error())
				s.shutdown()
				return
			}
			if f.messageType == websocket.CloseMessage {
				return
			}
		case <-ticker.C:
			_ = s.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				s.shutdown()
				return
			}
		case <-s.done:
			return
		}
	}
}

// enqueue hands a frame to the write pump.
//
// A client that stops reading must not be able to wedge the goroutine copying
// output from the remote target: once the outbound buffer fills, the producer
// waits at most writeWait for room and then kills the console. Blocking
// indefinitely would leak a goroutine and an upstream exec stream per stalled
// browser tab.
func (s *Socket) enqueue(f frame) error {
	// Refuse outright once the socket is down. This cannot be folded into the
	// select below: when both `done` and `outbound` are ready, Go chooses
	// between them at random, so a dead socket would accept writes about half
	// the time.
	select {
	case <-s.done:
		return io.ErrClosedPipe
	default:
	}

	timer := time.NewTimer(writeWait)
	defer timer.Stop()

	select {
	case s.outbound <- f:
		return nil
	case <-s.done:
		return io.ErrClosedPipe
	case <-timer.C:
		s.log.Warn(ErrSlowConsumer())
		s.shutdown()
		return io.ErrClosedPipe
	}
}

// SendControl queues a control frame for the client.
func (s *Socket) SendControl(msg ControlMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return s.enqueue(frame{messageType: websocket.TextMessage, payload: payload})
}

// Stdin returns the stream of bytes the client has typed. It reports io.EOF
// once the client disconnects.
func (s *Socket) Stdin() io.Reader { return &stdinReader{socket: s} }

// Stdout returns a writer whose every Write becomes one binary frame.
func (s *Socket) Stdout() io.Writer { return &binaryWriter{socket: s} }

// Resize yields the client's terminal geometry each time it changes.
func (s *Socket) Resize() <-chan TerminalSize { return s.resize }

// Close sends a close frame with the given code and reason, waits briefly for
// it to reach the wire, and tears the socket down. It is safe to call more than
// once; only the first call sends a frame.
func (s *Socket) Close(code int, reason string) {
	select {
	case <-s.done:
		// Already torn down; the peer will see the TCP close.
	default:
		payload := websocket.FormatCloseMessage(code, reason)
		// Best-effort: if the queue is wedged, shutdown below still runs.
		select {
		case s.outbound <- frame{messageType: websocket.CloseMessage, payload: payload}:
			// Give the write pump a moment to flush the close frame, so the
			// client learns *why* the console ended instead of seeing an
			// abnormal closure.
			select {
			case <-s.writePumpDone:
			case <-time.After(writeWait):
			}
		case <-time.After(writeWait):
		}
	}
	s.shutdown()
	_ = s.ws.Close()
}

// shutdown closes done exactly once, unblocking every pump and producer.
func (s *Socket) shutdown() {
	s.closeOnce.Do(func() { close(s.done) })
}

// binaryWriter adapts the socket to io.Writer, emitting one binary frame per
// Write call.
type binaryWriter struct{ socket *Socket }

func (w *binaryWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	// The caller — client-go's stream copier — reuses p across calls, and the
	// frame sits in a queue until the write pump drains it. Handing the queue a
	// slice that aliases p would let the next read scribble over bytes that
	// have not been written yet, interleaving garbage into the terminal.
	payload := make([]byte, len(p))
	copy(payload, p)

	if err := w.socket.enqueue(frame{messageType: websocket.BinaryMessage, payload: payload}); err != nil {
		return 0, err
	}
	return len(p), nil
}

// stdinReader adapts the socket's inbound binary frames to io.Reader, carrying
// a remainder across calls when the consumer's buffer is smaller than a frame.
type stdinReader struct {
	socket    *Socket
	remainder []byte
}

func (r *stdinReader) Read(p []byte) (int, error) {
	if len(r.remainder) == 0 {
		select {
		case chunk, ok := <-r.socket.stdin:
			if !ok {
				return 0, io.EOF
			}
			r.remainder = chunk
		case <-r.socket.done:
			return 0, io.EOF
		}
	}
	n := copy(p, r.remainder)
	r.remainder = r.remainder[n:]
	return n, nil
}
