package router

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"slices"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/handlers"
	"github.com/meshery/meshery/server/models"
)

// sseRoutes lists every Server-Sent Events route NewRouter serves whose
// handler serves until its client leaves. Each must end when the server
// starts shutting down.
var sseRoutes = []string{
	"/api/providers/stream",
	"/api/system/adapter/events",
	"/api/system/events/subscribe",
	"/api/system/controllers/status/subscribe",
}

// loadTestRoutes stream a load-test run until it ends. They must be drained
// on shutdown, not cut off, because the run saves its result as it ends.
// Sorted, to compare with streamingHandlers.openStreams.
var loadTestRoutes = []string{
	"/api/perf/profile",
	"/api/user/performance/profiles/1/run",
}

// streamingHandlers stands in for the real handlers behind NewRouter. Its
// middlewares pass requests straight through. As in the real handlers, an
// SSE handler serves until its request context is done, and a load-test
// handler serves until its run ends (here, until release is closed) or its
// request context is done.
type streamingHandlers struct {
	*handlers.Handler

	mu      sync.Mutex
	open    map[string]bool
	release chan struct{}
}

func newStreamingHandlers() *streamingHandlers {
	return &streamingHandlers{
		Handler: &handlers.Handler{},
		open:    map[string]bool{},
		release: make(chan struct{}),
	}
}

type sessionHandlerFunc = func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)

func (s *streamingHandlers) ProviderMiddleware(next http.Handler) http.Handler { return next }

func (s *streamingHandlers) AuthMiddleware(next http.Handler, _ models.AuthenticationMechanism) http.Handler {
	return next
}

func (s *streamingHandlers) SessionInjectorMiddleware(next sessionHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next(w, r, nil, nil, nil) })
}

func (s *streamingHandlers) KubernetesMiddleware(next sessionHandlerFunc) sessionHandlerFunc {
	return next
}

func (s *streamingHandlers) ProvidersStreamHandler(w http.ResponseWriter, r *http.Request) {
	s.serveStream(w, r, nil)
}

func (s *streamingHandlers) EventStreamHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	s.serveStream(w, r, nil)
}

func (s *streamingHandlers) SubscribeEventsHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	s.serveStream(w, r, nil)
}

func (s *streamingHandlers) SubscribeMesheryControllersStatusHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	s.serveStream(w, r, nil)
}

func (s *streamingHandlers) LoadTestHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	s.serveStream(w, r, s.release)
}

// serveStream opens an SSE response and holds it until the request context is
// done or, when runEnded is not nil, until runEnded is closed.
func (s *streamingHandlers) serveStream(w http.ResponseWriter, r *http.Request, runEnded <-chan struct{}) {
	s.setOpen(r.URL.Path, true)
	defer s.setOpen(r.URL.Path, false)

	w.Header().Set("Content-Type", "text/event-stream")
	w.WriteHeader(http.StatusOK)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	select {
	case <-r.Context().Done():
	case <-runEnded:
	}
}

func (s *streamingHandlers) setOpen(path string, open bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if open {
		s.open[path] = true
	} else {
		delete(s.open, path)
	}
}

func (s *streamingHandlers) openStreams() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	paths := make([]string, 0, len(s.open))
	for path := range s.open {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

// TestNewRouterEndsStreamsOnShutdown opens every SSE route NewRouter serves,
// plus both load-test routes with a run in progress, then shuts the router
// down.
//
// http.Server.Shutdown waits for active connections to go idle but never
// cancels their request contexts, and an SSE handler that serves until its
// client leaves ends only when its request context does. Every Meshery UI tab
// holds such a stream for as long as it is open, so an unwrapped stream holds
// Shutdown for its full timeout. Those streams must end as soon as Shutdown
// starts. A load-test run must still be drained, because it saves its result
// as it ends: Shutdown returns only once the runs finish.
func TestNewRouterEndsStreamsOnShutdown(t *testing.T) {
	h := newStreamingHandlers()
	endRuns := sync.OnceFunc(func() { close(h.release) })

	port, err := freePort()
	if err != nil {
		t.Fatalf("freePort() failed with error: %v", err)
	}
	r := NewRouter(context.Background(), h, port, http.NotFoundHandler(), http.NotFoundHandler())
	go func() { _ = r.Run() }()
	t.Cleanup(func() {
		// Only does anything if the test failed partway.
		endRuns()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = r.Shutdown(ctx)
	})
	addr := waitForListener(t, port)

	all := append(slices.Clone(sseRoutes), loadTestRoutes...)
	for _, path := range all {
		openStream(t, addr, path)
	}
	if got := h.openStreams(); len(got) != len(all) {
		t.Fatalf("open streams = %v, want all of %v", got, all)
	}

	shutdownErr := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		shutdownErr <- r.Shutdown(ctx)
	}()

	deadline := time.Now().Add(2 * time.Second)
	for !slices.Equal(h.openStreams(), loadTestRoutes) && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if open := h.openStreams(); !slices.Equal(open, loadTestRoutes) {
		t.Fatalf("open streams 2s after Shutdown started = %v, want %v: SSE streams must end when Shutdown starts, and load-test runs must be drained", open, loadTestRoutes)
	}

	select {
	case err := <-shutdownErr:
		t.Fatalf("Shutdown returned (%v) while load-test runs were still in progress, want them drained", err)
	case <-time.After(200 * time.Millisecond):
	}

	endRuns()
	select {
	case err := <-shutdownErr:
		if err != nil {
			t.Fatalf("Shutdown() failed with error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Shutdown() did not return after the load-test runs ended")
	}
	if open := h.openStreams(); len(open) > 0 {
		t.Fatalf("streams %v still open after Shutdown returned", open)
	}
}

// TestRouterShutdownClosesConnectionsOnTimeout holds a request that ends only
// when its context does, on a route that is not wrapped with endOnShutdown,
// past the Shutdown deadline. http.Server.Shutdown returns at the deadline but
// leaves that connection open and its handler running, while the caller goes
// on to close the provider session and database the handler uses.
// Router.Shutdown must close the connection, which cancels the request
// context.
func TestRouterShutdownClosesConnectionsOnTimeout(t *testing.T) {
	handlerDone := make(chan struct{})
	entered := make(chan struct{})
	gMux := mux.NewRouter()
	gMux.HandleFunc("/hold", func(w http.ResponseWriter, r *http.Request) {
		defer close(handlerDone)
		w.WriteHeader(http.StatusOK)
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		close(entered)
		<-r.Context().Done()
	})

	port, err := freePort()
	if err != nil {
		t.Fatalf("freePort() failed with error: %v", err)
	}
	r := &Router{S: gMux, port: port, srv: newHTTPServer(port, gMux)}
	go func() { _ = r.Run() }()
	addr := waitForListener(t, port)

	openStream(t, addr, "/hold")
	<-entered

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	if err := r.Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Shutdown() = %v, want %v", err, context.DeadlineExceeded)
	}

	select {
	case <-handlerDone:
	case <-time.After(2 * time.Second):
		t.Fatal("handler still running 2s after Shutdown returned, want its connection closed and its request context cancelled")
	}
}

// openStream sends GET path on a new connection and waits for a 200 response
// header, leaving the response body open. The connection is closed when the
// test ends.
func openStream(t *testing.T, addr, path string) {
	t.Helper()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("dial %s: %v", addr, err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	if _, err := fmt.Fprintf(conn, "GET %s HTTP/1.1\r\nHost: %s\r\nAccept: text/event-stream\r\n\r\n", path, addr); err != nil {
		t.Fatalf("writing GET %s: %v", path, err)
	}
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		t.Fatalf("reading GET %s response: %v", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s returned %d, want %d", path, resp.StatusCode, http.StatusOK)
	}
	_ = conn.SetReadDeadline(time.Time{})
}
