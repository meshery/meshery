package models

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// testRemoteFileClient allows local test servers and limits redirects while
// preserving the production client's HTTP(S)-only validation.
func testRemoteFileClient(timeout time.Duration) *http.Client {
	client := newRemoteFileClient(timeout, true)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > maxRemoteFileRedirects {
			return fmt.Errorf("stopped after %d redirects", maxRemoteFileRedirects)
		}
		return validateRemoteFileURL(req.URL)
	}
	return client
}

func TestFetchRemoteFileRejectsLoopback(t *testing.T) {
	var hit atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hit.Store(true)
		_, _ = w.Write([]byte("name: x"))
	}))
	defer srv.Close()

	_, err := fetchRemoteFile(srv.URL)
	if err == nil || !strings.Contains(err.Error(), "restricted") {
		t.Fatalf("expected restricted-address error, got %v", err)
	}
	if hit.Load() {
		t.Fatal("request reached the loopback server")
	}
}

func TestFetchRemoteFileRejectsUnsupportedScheme(t *testing.T) {
	for _, u := range []string{"file:///etc/passwd", "ftp://example.com/f", "gopher://example.com", "//example.com/f", "http:///f"} {
		if _, err := fetchRemoteFile(u); err == nil {
			t.Errorf("fetchRemoteFile(%q) succeeded, want error", u)
		}
	}
}

func TestFetchRemoteFileReturnsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("name: x"))
	}))
	defer srv.Close()

	body, err := fetchRemoteFileWith(testRemoteFileClient(5*time.Second), srv.URL, 1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != "name: x" {
		t.Fatalf("body = %q", body)
	}
}

func TestFetchRemoteFileNonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()

	if _, err := fetchRemoteFileWith(testRemoteFileClient(5*time.Second), srv.URL, 1024); err == nil || err.Error() != "file not found" {
		t.Fatalf("expected 'file not found', got %v", err)
	}
}

func TestFetchRemoteFileEnforcesSizeLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("a", 2048)))
	}))
	defer srv.Close()

	client := testRemoteFileClient(5 * time.Second)
	if _, err := fetchRemoteFileWith(client, srv.URL, 1024); err == nil || !strings.Contains(err.Error(), "maximum allowed size") {
		t.Fatalf("expected size error, got %v", err)
	}
	// Exactly at the limit is accepted.
	if _, err := fetchRemoteFileWith(client, srv.URL, 2048); err != nil {
		t.Fatalf("body at the limit should be accepted, got %v", err)
	}
}

func TestFetchRemoteFileTimesOut(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		<-release
	}))
	defer srv.Close()
	defer close(release)

	client := testRemoteFileClient(100 * time.Millisecond)
	start := time.Now()
	if _, err := fetchRemoteFileWith(client, srv.URL, 1024); err == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("fetch did not time out promptly: %v", elapsed)
	}
}

func TestFetchRemoteFileRedirects(t *testing.T) {
	client := testRemoteFileClient(5 * time.Second)

	t.Run("redirect to unsupported scheme is rejected", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "file:///etc/passwd", http.StatusFound)
		}))
		defer srv.Close()
		if _, err := fetchRemoteFileWith(client, srv.URL, 1024); err == nil || !strings.Contains(err.Error(), "unsupported URL scheme") {
			t.Fatalf("expected scheme error, got %v", err)
		}
	})

	// chain serves a redirect chain: /n redirects to /n-1 until /0, which
	// returns a body. Fetching /n therefore follows exactly n redirects.
	chain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if n == 0 {
			_, _ = w.Write([]byte("name: ok"))
			return
		}
		http.Redirect(w, r, "/"+strconv.Itoa(n-1), http.StatusFound)
	}))
	defer chain.Close()

	t.Run("redirects up to the limit are followed", func(t *testing.T) {
		u := chain.URL + "/" + strconv.Itoa(maxRemoteFileRedirects)
		if _, err := fetchRemoteFileWith(client, u, 1024); err != nil {
			t.Fatalf("expected %d redirects to be followed, got %v", maxRemoteFileRedirects, err)
		}
	})

	t.Run("one redirect past the limit is rejected", func(t *testing.T) {
		u := chain.URL + "/" + strconv.Itoa(maxRemoteFileRedirects+1)
		if _, err := fetchRemoteFileWith(client, u, 1024); err == nil || !strings.Contains(err.Error(), "redirects") {
			t.Fatalf("expected redirect limit error, got %v", err)
		}
	})

	t.Run("redirect loop is bounded", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, r.URL.Path+"x", http.StatusFound)
		}))
		defer srv.Close()
		if _, err := fetchRemoteFileWith(client, srv.URL, 1024); err == nil || !strings.Contains(err.Error(), "redirects") {
			t.Fatalf("expected redirect limit error, got %v", err)
		}
	})
}

func TestRemoteFileClientBlocksInternalRedirectTarget(t *testing.T) {
	internal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("name: secret"))
	}))
	defer internal.Close()

	entry := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, internal.URL, http.StatusFound)
	}))
	defer entry.Close()

	client := newRemoteFileClient(5*time.Second, true)
	if _, err := fetchRemoteFileWith(client, entry.URL, 1024); err == nil || !strings.Contains(err.Error(), "redirect blocked") {
		t.Fatalf("expected the redirect target to be refused, got %v", err)
	}
}
