package router

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// testBodyIdle stands in for requestBodyIdleTimeout so these tests run in
// well under a second each.
const testBodyIdle = 200 * time.Millisecond

// serveWithBodyBound starts a real HTTP server (so connection read deadlines
// apply) with only a header timeout, like newHTTPServer, and handler wrapped
// in boundRequestBodyReads with the test idle timeout. It returns the address
// to dial.
func serveWithBodyBound(t *testing.T, handler http.HandlerFunc) string {
	t.Helper()
	return serveWithBodyBoundIdle(t, testBodyIdle, handler)
}

func serveWithBodyBoundIdle(t *testing.T, idle time.Duration, handler http.HandlerFunc) string {
	t.Helper()
	srv := httptest.NewUnstartedServer(boundRequestBodyReads(handler, idle))
	srv.Config.ReadHeaderTimeout = 5 * time.Second
	srv.Start()
	t.Cleanup(srv.Close)
	return srv.Listener.Addr().String()
}

// sendRaw dials addr and writes raw request bytes, so a test controls exactly
// when, or whether, the request body arrives.
func sendRaw(t *testing.T, addr, raw string) net.Conn {
	t.Helper()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("dial %s: %v", addr, err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	if _, err := io.WriteString(conn, raw); err != nil {
		t.Fatalf("write request: %v", err)
	}
	return conn
}

type bodyReadResult struct {
	n   int
	err error
}

// readBodyHandler reads the whole request body and reports what it got.
func readBodyHandler(results chan<- bodyReadResult) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		results <- bodyReadResult{n: len(b), err: err}
	}
}

func TestBoundRequestBodyReads_StalledBodyFails(t *testing.T) {
	results := make(chan bodyReadResult, 1)
	addr := serveWithBodyBound(t, readBodyHandler(results))

	start := time.Now()
	// Declare 6 bytes, send 2, then stop.
	sendRaw(t, addr, "POST /upload HTTP/1.1\r\nHost: x\r\nContent-Length: 6\r\n\r\naa")

	select {
	case got := <-results:
		if got.err == nil {
			t.Fatalf("reading a stalled body returned %d bytes and no error, want an i/o timeout", got.n)
		}
		if elapsed := time.Since(start); elapsed < testBodyIdle {
			t.Fatalf("stalled body failed after %v, before the %v idle timeout", elapsed, testBodyIdle)
		}
	case <-time.After(10 * testBodyIdle):
		t.Fatal("reading a stalled body did not fail; the handler is held by a client that stopped sending")
	}
}

func TestBoundRequestBodyReads_ProgressingBodyIsReadInFull(t *testing.T) {
	results := make(chan bodyReadResult, 1)
	addr := serveWithBodyBound(t, readBodyHandler(results))

	// 6 bytes, each arriving within the idle timeout of the last, so the
	// upload as a whole takes three times as long as the timeout.
	conn := sendRaw(t, addr, "POST /upload HTTP/1.1\r\nHost: x\r\nContent-Length: 6\r\n\r\n")
	go func() {
		for i := 0; i < 6; i++ {
			time.Sleep(testBodyIdle / 2)
			if _, err := conn.Write([]byte{'a'}); err != nil {
				return
			}
		}
	}()

	select {
	case got := <-results:
		if got.err != nil || got.n != 6 {
			t.Fatalf("read %d bytes, err %v; want all 6 bytes of an upload that kept arriving", got.n, got.err)
		}
	case <-time.After(20 * testBodyIdle):
		t.Fatal("handler never finished reading the upload")
	}
}

// TestBoundRequestBodyReads_WithheldBodyReleasesStream covers a handler that
// never reads the body, such as an SSE stream. net/http drains an unread body
// before sending the response headers, so a client that declares a body and
// never sends it would otherwise hold the handler at its first flush forever.
func TestBoundRequestBodyReads_WithheldBodyReleasesStream(t *testing.T) {
	released := make(chan struct{})
	addr := serveWithBodyBound(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
		close(released)
	})

	sendRaw(t, addr, "GET /stream HTTP/1.1\r\nHost: x\r\nContent-Length: 1\r\n\r\n")

	select {
	case <-released:
	case <-time.After(10 * testBodyIdle):
		t.Fatal("stream is still held by a client that declared a body and never sent it")
	}
}

// TestBoundRequestBodyReads_LeavesHandlerAloneAfterBody checks that once the
// body has been read in full, the deadline no longer applies: a handler that
// keeps running afterwards, as a long request or a stream does, must not have
// its context cancelled when the idle timeout passes.
func TestBoundRequestBodyReads_LeavesHandlerAloneAfterBody(t *testing.T) {
	outcome := make(chan error, 1)
	addr := serveWithBodyBound(t, func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.ReadAll(r.Body); err != nil {
			outcome <- fmt.Errorf("reading body: %w", err)
			return
		}
		select {
		case <-r.Context().Done():
			outcome <- fmt.Errorf("request context cancelled after the body was read: %w", r.Context().Err())
		case <-time.After(3 * testBodyIdle):
			outcome <- nil
		}
	})

	sendRaw(t, addr, "POST /upload HTTP/1.1\r\nHost: x\r\nContent-Length: 3\r\n\r\nabc")

	select {
	case err := <-outcome:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(20 * testBodyIdle):
		t.Fatal("handler did not finish")
	}
}

// TestBoundRequestBodyReads_EarlyRejectionOfExpectContinue sends a request
// with Expect: 100-continue and holds its body back, as a client does until
// the server says to go ahead. A handler that rejects it without reading the
// body must answer at once: net/http recognises the unread 100-continue body
// and closes the connection instead of waiting for it, as long as the
// middleware leaves the server's own request untouched.
func TestBoundRequestBodyReads_EarlyRejectionOfExpectContinue(t *testing.T) {
	// A long idle timeout, so waiting for the body would show as a clear delay.
	const idle = 5 * time.Second
	addr := serveWithBodyBoundIdle(t, idle, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	start := time.Now()
	conn := sendRaw(t, addr, "POST /upload HTTP/1.1\r\nHost: x\r\nContent-Length: 5\r\nExpect: 100-continue\r\n\r\n")
	_ = conn.SetReadDeadline(time.Now().Add(2 * idle))
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		t.Fatalf("reading response: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
	}
	if elapsed := time.Since(start); elapsed >= idle/2 {
		t.Fatalf("early rejection took %v; the server waited for a body the client was holding back", elapsed)
	}
}

func TestBoundRequestBodyReads_LeavesRequestsWithoutBodyAlone(t *testing.T) {
	outcome := make(chan error, 1)
	addr := serveWithBodyBound(t, func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			outcome <- fmt.Errorf("request context cancelled for a request without a body: %w", r.Context().Err())
		case <-time.After(3 * testBodyIdle):
			outcome <- nil
		}
	})

	sendRaw(t, addr, "GET /stream HTTP/1.1\r\nHost: x\r\n\r\n")

	select {
	case err := <-outcome:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(20 * testBodyIdle):
		t.Fatal("handler did not finish")
	}
}
