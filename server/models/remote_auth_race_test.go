//go:build race

package models

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// TestRemoteProviderRefreshToken_TokenStoreConcurrentEviction guards against a
// regression of the data race between refreshToken's time.AfterFunc eviction
// callback (delete(l.TokenStore, ...)) and the lock-protected readers
// ExtractToken / UpdateToken. All TokenStore access must hold TokenStoreMut.
//
// It spins tight loops to widen the race window, so it is gated behind the
// `race` build tag and only runs under `go test -race`.
func TestRemoteProviderRefreshToken_TokenStoreConcurrentEviction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/refresh" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"token":"refreshed-token"}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	// Per-instance delay; no global state mutation.
	provider.TokenDeletionDelay = time.Millisecond

	const oldToken = "old-token-string"

	newReaderRequest := func() *http.Request {
		req := httptest.NewRequest(http.MethodGet, "http://localhost/api/token", nil)
		req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: oldToken})
		return req
	}

	var wg sync.WaitGroup
	stop := make(chan struct{})

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := newReaderRequest()
			for {
				select {
				case <-stop:
					return
				default:
				}
				rec := httptest.NewRecorder()
				provider.ExtractToken(rec, req)
			}
		}()
	}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				_, _ = provider.refreshToken(oldToken)
				time.Sleep(100 * time.Microsecond)
			}
		}()
	}

	time.Sleep(300 * time.Millisecond)
	close(stop)
	wg.Wait()
	// Allow any in-flight eviction timers to run before the test exits.
	time.Sleep(5 * provider.TokenDeletionDelay)
}
