package router

import (
	"io"
	"net/http"
	"time"
)

// requestBodyIdleTimeout bounds how long the server waits for the next bytes
// of a request body. It replaces a whole-request ReadTimeout: an upload that
// keeps arriving, however slowly, is never cut off, but one that stops
// arriving fails after this long instead of holding its handler and
// connection open indefinitely.
const requestBodyIdleTimeout = 30 * time.Second

// boundRequestBodyReads arms a read deadline on every request that declares a
// body and pushes it back each time body bytes arrive, so reading the body
// fails once the client stops sending it for longer than idle.
//
// This also covers handlers that never read the body. Before it sends the
// response headers, net/http reads and discards a small unread body. Without
// a deadline, a client that declares a body and never sends it would hold
// that handler at its first write or flush forever, as an SSE stream's first
// flush would be.
//
// Requests without a body are left alone. Once a body has been read to EOF,
// net/http clears the read deadline itself, so it never reaches a
// long-running handler or stream after that point.
func boundRequestBodyReads(next http.Handler, idle time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength != 0 && r.Body != nil && r.Body != http.NoBody {
			rc := http.NewResponseController(w)
			// SetReadDeadline reports ErrNotSupported for writers that don't
			// implement it; there body reads stay unbounded, as before.
			_ = rc.SetReadDeadline(time.Now().Add(idle))
			// Wrap the body on a shallow copy, never on r itself. r is the
			// request net/http keeps for writing the response, and it checks
			// the body's own type there: an Expect: 100-continue request
			// rejected before its body is read is only answered at once while
			// net/http can still see its expectContinueReader.
			wrapped := r.WithContext(r.Context())
			wrapped.Body = &idleDeadlineBody{ReadCloser: r.Body, rc: rc, idle: idle}
			r = wrapped
		}
		next.ServeHTTP(w, r)
	})
}

// idleDeadlineBody pushes the connection's read deadline back after every
// read that returns body bytes.
type idleDeadlineBody struct {
	io.ReadCloser
	rc   *http.ResponseController
	idle time.Duration
}

func (b *idleDeadlineBody) Read(p []byte) (int, error) {
	n, err := b.ReadCloser.Read(p)
	// Only extend while the body is still open. net/http returns the last
	// bytes together with io.EOF after it has already cleared the read
	// deadline for its background read; re-arming it then would cancel the
	// request's context once idle elapsed, long after the body was done.
	if n > 0 && err == nil {
		_ = b.rc.SetReadDeadline(time.Now().Add(b.idle))
	}
	return n, err
}
