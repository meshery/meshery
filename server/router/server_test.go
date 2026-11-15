package router

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestClose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	t.Log("Need to run Close() skipping")
	//err := r.Close()
	//if err != nil {
	//	t.Errorf("Close() failed with error: %s", err)
	//}
}

// TestRouterShutdownDrainsListener locks the graceful-shutdown contract:
// Run serves requests until Shutdown is called, Shutdown waits for an
// in-flight request to finish instead of cutting it off, and Run then
// returns http.ErrServerClosed instead of an error the caller would treat
// as a crash.
func TestRouterShutdownDrainsListener(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	release := make(chan struct{})
	gMux := mux.NewRouter()
	gMux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		<-release
		w.WriteHeader(http.StatusNoContent)
	})

	port, err := freePort()
	if err != nil {
		t.Fatalf("freePort() failed with error: %v", err)
	}
	r := &Router{S: gMux, port: port, srv: newHTTPServer(port, gMux)}

	runErr := make(chan error, 1)
	go func() {
		runErr <- r.Run()
	}()

	// Wait until the listener accepts connections, so the request below is
	// genuinely in flight and not refused by a listener that is not up yet.
	addr := waitForListener(t, port)

	// Start a request and hold it in flight while Shutdown runs.
	reqErr := make(chan error, 1)
	go func() {
		resp, err := http.Get("http://" + addr + "/slow")
		if err != nil {
			reqErr <- err
			return
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusNoContent {
			reqErr <- fmt.Errorf("GET /slow returned status %d, want %d", resp.StatusCode, http.StatusNoContent)
			return
		}
		reqErr <- nil
	}()

	// Give the request time to reach the handler before shutting down.
	time.Sleep(200 * time.Millisecond)

	shutdownErr := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownErr <- r.Shutdown(ctx)
	}()

	// Shutdown must wait for the in-flight request instead of returning
	// while it is still running.
	select {
	case err := <-shutdownErr:
		t.Fatalf("Shutdown returned before the in-flight request finished: %v", err)
	case <-time.After(500 * time.Millisecond):
	}

	close(release)

	select {
	case err := <-shutdownErr:
		if err != nil {
			t.Fatalf("Shutdown() failed with error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Shutdown() did not return after the in-flight request finished")
	}

	if err := <-reqErr; err != nil {
		t.Fatalf("in-flight request failed: %v", err)
	}

	select {
	case err := <-runErr:
		if !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("Run() returned %v, want %v", err, http.ErrServerClosed)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return after Shutdown()")
	}
}

// TestNewHTTPServerTimeouts pins which timeouts the server sets. Meshery
// serves Server-Sent Event streams that stay open for minutes to hours (a load
// test run, the provider status stream, adapter events) and accepts uploads
// whose duration depends on the client's link, not the server. ReadTimeout and
// WriteTimeout are absolute deadlines measured from the start of each request,
// so either one cuts those requests off partway through: the final frame of a
// load test longer than WriteTimeout is never delivered, and an upload still
// arriving after ReadTimeout fails with an i/o timeout.
//
// Slow clients are bounded by ReadHeaderTimeout instead, and IdleTimeout has
// to be set explicitly: with ReadTimeout unset it would otherwise fall back to
// no idle timeout at all.
func TestNewHTTPServerTimeouts(t *testing.T) {
	srv := newHTTPServer(9081, mux.NewRouter())

	if srv.ReadTimeout != 0 {
		t.Errorf("ReadTimeout = %v, want 0: it fails request bodies still uploading when it expires", srv.ReadTimeout)
	}
	if srv.WriteTimeout != 0 {
		t.Errorf("WriteTimeout = %v, want 0: it cuts off SSE streams and long-running responses when it expires", srv.WriteTimeout)
	}
	if srv.ReadHeaderTimeout <= 0 {
		t.Errorf("ReadHeaderTimeout = %v, want > 0 so a client cannot hold a connection open by trickling headers", srv.ReadHeaderTimeout)
	}
	if srv.IdleTimeout <= 0 {
		t.Errorf("IdleTimeout = %v, want > 0 so idle keep-alive connections are closed", srv.IdleTimeout)
	}
}

// TestRouterAcceptsSlowRequestBody sends a request body that takes longer to
// arrive than the 5s ReadTimeout the server used to set, as a model or design
// upload does over a slow link or kubectl port-forward, and checks that the
// handler receives all of it instead of failing with an i/o timeout.
func TestRouterAcceptsSlowRequestBody(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	const (
		bodyBytes = 12
		interval  = 500 * time.Millisecond // 12 x 500ms = 6s of upload
	)

	gMux := mux.NewRouter()
	gMux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "reading request body: %v", err)
			return
		}
		_, _ = fmt.Fprintf(w, "%d", len(body))
	})

	port, err := freePort()
	if err != nil {
		t.Fatalf("freePort() failed with error: %v", err)
	}
	r := &Router{S: gMux, port: port, srv: newHTTPServer(port, gMux)}
	go func() { _ = r.Run() }()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.Shutdown(ctx)
	})
	addr := waitForListener(t, port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("dial %s: %v", addr, err)
	}
	defer func() { _ = conn.Close() }()

	if _, err := fmt.Fprintf(conn, "POST /upload HTTP/1.1\r\nHost: %s\r\nContent-Length: %d\r\nConnection: close\r\n\r\n", addr, bodyBytes); err != nil {
		t.Fatalf("writing request headers: %v", err)
	}
	go func() {
		for i := 0; i < bodyBytes; i++ {
			time.Sleep(interval)
			if _, err := conn.Write([]byte{'a'}); err != nil {
				// The server has already answered; the response read below
				// reports why.
				return
			}
		}
	}()

	_ = conn.SetReadDeadline(time.Now().Add(bodyBytes*interval + 5*time.Second))
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		t.Fatalf("reading response: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST /upload returned %d (%s), want %d", resp.StatusCode, strings.TrimSpace(string(got)), http.StatusOK)
	}
	if want := fmt.Sprint(bodyBytes); string(got) != want {
		t.Fatalf("handler read %s body bytes, want %s", got, want)
	}
}

// waitForListener blocks until the server accepts connections on port, so a
// test request is not refused by a listener that is not up yet, and returns
// the address to dial.
func waitForListener(t *testing.T, port int) string {
	t.Helper()
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	deadline := time.Now().Add(5 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return addr
		}
		if time.Now().After(deadline) {
			t.Fatalf("server did not start listening on :%d within 5s: %v", port, err)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func freePort() (int, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer func() { _ = ln.Close() }()
	return ln.Addr().(*net.TCPAddr).Port, nil
}
