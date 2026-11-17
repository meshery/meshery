package models

import (
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

)

func allowAllAddrs(netip.Addr) bool { return false }

func TestIsBlockedRemoteFileAddr(t *testing.T) {
	tests := []struct {
		addr    string
		blocked bool
	}{
		{"127.0.0.1", true},
		{"::1", true},
		{"::ffff:127.0.0.1", true},
		{"10.1.2.3", true},
		{"172.16.0.1", true},
		{"192.168.1.1", true},
		{"169.254.169.254", true},
		{"fe80::1", true},
		{"fd00:ec2::254", true},
		{"0.0.0.0", true},
		{"0.1.2.3", true},
		{"::", true},
		{"100.64.0.1", true},
		{"224.0.0.1", true},
		{"240.0.0.1", true},
		{"192.0.2.1", true},
		{"198.51.100.1", true},
		{"203.0.113.1", true},
		{"2001:db8::1", true},
		// Zoned addresses must be judged without their zone.
		{"::%lo", true},
		{"::1%lo", true},
		{"fe80::1%eth0", true},
		{"2001:db8::1%eth0", true},
		{"2606:4700:4700::1111%eth0", false},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
		{"2606:4700:4700::1111", false},
	}
	for _, tc := range tests {
		t.Run(tc.addr, func(t *testing.T) {
			if got := isBlockedRemoteFileAddr(netip.MustParseAddr(tc.addr)); got != tc.blocked {
				t.Fatalf("isBlockedRemoteFileAddr(%s) = %v, want %v", tc.addr, got, tc.blocked)
			}
		})
	}
}

func TestFetchRemoteFileRejectsLoopback(t *testing.T) {
	var hit atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hit.Store(true)
		_, _ = w.Write([]byte("name: x"))
	}))
	defer srv.Close()

	_, err := fetchRemoteFile(srv.URL)
	if err == nil || !strings.Contains(err.Error(), "non-public address") {
		t.Fatalf("expected non-public address error, got %v", err)
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

	client := newRemoteFileClient(5*time.Second, allowAllAddrs)
	body, err := fetchRemoteFileWith(client, srv.URL, 1024)
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

	client := newRemoteFileClient(5*time.Second, allowAllAddrs)
	if _, err := fetchRemoteFileWith(client, srv.URL, 1024); err == nil || err.Error() != "file not found" {
		t.Fatalf("expected 'file not found', got %v", err)
	}
}

func TestFetchRemoteFileEnforcesSizeLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("a", 2048)))
	}))
	defer srv.Close()

	client := newRemoteFileClient(5*time.Second, allowAllAddrs)
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

	client := newRemoteFileClient(100*time.Millisecond, allowAllAddrs)
	start := time.Now()
	if _, err := fetchRemoteFileWith(client, srv.URL, 1024); err == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("fetch did not time out promptly: %v", elapsed)
	}
}

func TestFetchRemoteFileRedirects(t *testing.T) {
	client := newRemoteFileClient(5*time.Second, allowAllAddrs)

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

// A redirect from a public host to an internal address must be refused by the
// dialer, not just by inspecting the initial URL. The first dial (the entry
// hop) is allowed and every later dial (the redirect hop) is refused.
func TestRemoteFileClientBlocksInternalRedirectTarget(t *testing.T) {
	internal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("name: secret"))
	}))
	defer internal.Close()

	entry := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, internal.URL, http.StatusFound)
	}))
	defer entry.Close()

	dials := 0
	client := newRemoteFileClient(5*time.Second, func(netip.Addr) bool {
		dials++
		return dials > 1
	})
	if _, err := fetchRemoteFileWith(client, entry.URL, 1024); err == nil || !strings.Contains(err.Error(), "non-public address") {
		t.Fatalf("expected the redirect hop to be refused, got %v", err)
	}
}
