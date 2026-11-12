package router

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
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
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	deadline := time.Now().Add(5 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("server did not start listening on :%d within 5s: %v", port, err)
		}
		time.Sleep(50 * time.Millisecond)
	}

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

func freePort() (int, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer func() { _ = ln.Close() }()
	return ln.Addr().(*net.TCPAddr).Port, nil
}
