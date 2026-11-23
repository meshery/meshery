package prometheus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewUsesGuardedTransport pins that the production constructor refuses to
// fetch a private/loopback address. Connection URLs reach this client from
// user-supplied connection metadata, so the transport is what has to enforce
// the destination policy.
func TestNewUsesGuardedTransport(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits++
	}))
	defer srv.Close()

	if _, err := New(srv.URL, "", nil).Query(context.Background(), "up", ""); err == nil {
		t.Fatal("New() fetched a loopback address; expected the guarded transport to refuse it")
	}
	if hits != 0 {
		t.Fatalf("guarded transport still reached the server %d time(s)", hits)
	}
}

// TestNewTestClientReachesLoopback keeps the test-only constructor honest: the
// package's own tests must still be able to drive a local httptest server.
func TestNewTestClientReachesLoopback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"status":"success"}`))
	}))
	defer srv.Close()

	if _, err := newTestClient(srv.URL, "", nil).Query(context.Background(), "up", ""); err != nil {
		t.Fatalf("newTestClient could not reach the local server: %v", err)
	}
}

// TestBaseURLFragmentDoesNotOverrideAPIPath pins that a query or fragment in the
// stored base URL cannot swallow the API path appended per request.
func TestBaseURLFragmentDoesNotOverrideAPIPath(t *testing.T) {
	for _, base := range []string{"/metrics#", "/metrics?x=1", "/metrics/"} {
		var got string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got = r.URL.Path
			_, _ = w.Write([]byte(`{"status":"success"}`))
		}))
		if _, err := newTestClient(srv.URL+base, "", nil).Query(context.Background(), "up", ""); err != nil {
			t.Fatalf("base %q: %v", base, err)
		}
		srv.Close()
		if want := "/metrics/api/v1/query"; got != want {
			t.Fatalf("base %q: server saw %q, want %q", base, got, want)
		}
	}
}
