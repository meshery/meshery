package handlers

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestIsPublicUnicastIP is a table-driven check of the address classes
// isPublicUnicastIP must reject: loopback, RFC1918/IPv6-ULA private ranges,
// link-local (including 169.254.169.254, the AWS/GCP/Azure metadata
// address that motivates this check), the unspecified address, and
// multicast — versus ordinary public unicast addresses, which it must
// allow.
func TestIsPublicUnicastIP(t *testing.T) {
	cases := []struct {
		name string
		ip   string
		want bool
	}{
		{"IPv4 loopback", "127.0.0.1", false},
		{"IPv6 loopback", "::1", false},
		{"RFC1918 10/8", "10.0.0.5", false},
		{"RFC1918 172.16/12", "172.16.5.1", false},
		{"RFC1918 192.168/16", "192.168.1.1", false},
		{"link-local cloud metadata address", "169.254.169.254", false},
		{"IPv6 link-local", "fe80::1", false},
		{"IPv6 unique-local (ULA)", "fd00::1", false},
		{"unspecified IPv4", "0.0.0.0", false},
		{"unspecified IPv6", "::", false},
		{"IPv4 multicast", "224.0.0.1", false},
		{"IPv4-mapped IPv6 loopback", "::ffff:127.0.0.1", false},
		{"IPv4-mapped IPv6 private", "::ffff:10.1.2.3", false},
		{"public IPv4", "93.184.216.34", true},
		{"public IPv4 (well-known resolver)", "8.8.8.8", true},
		{"public IPv6", "2606:4700:4700::1111", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ip := net.ParseIP(tc.ip)
			if ip == nil {
				t.Fatalf("net.ParseIP(%q) returned nil", tc.ip)
			}
			if got := isPublicUnicastIP(ip); got != tc.want {
				t.Errorf("isPublicUnicastIP(%s) = %v, want %v", tc.ip, got, tc.want)
			}
		})
	}
}

// TestValidateImportURLScheme covers the up-front scheme/host check that
// runs before any network activity.
func TestValidateImportURLScheme(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"https allowed", "https://example.com/design.yaml", false},
		{"http allowed", "http://example.com/design.yaml", false},
		{"file scheme rejected", "file:///etc/passwd", true},
		{"ftp scheme rejected", "ftp://example.com/design.yaml", true},
		{"gopher scheme rejected", "gopher://example.com", true},
		{"no scheme rejected", "example.com/design.yaml", true},
		{"no host rejected", "https:///design.yaml", true},
		{"unparseable URL rejected", "http://%zz", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateImportURLScheme(tc.url)
			if (err != nil) != tc.wantErr {
				t.Errorf("validateImportURLScheme(%q) error = %v, wantErr %v", tc.url, err, tc.wantErr)
			}
		})
	}
}

// TestDesignImportHTTPClient_BlocksLoopback proves the protection is
// actually wired into the client fetchFileFromURL uses, not just defined
// and unused. httptest.NewServer binds to a loopback address, so a
// successful request here would mean an attacker-supplied URL pointing at
// any service on the Meshery host's own loopback interface — or, by the
// same code path, at cluster-internal addresses or cloud metadata — would
// be fetched on the caller's behalf. This is exercised through the real
// package-level designImportHTTPClient, and through http.Client's redirect
// handling, so it also confirms a redirect to a loopback address is
// blocked rather than followed.
func TestDesignImportHTTPClient_BlocksLoopback(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("should never be reachable"))
	}))
	defer target.Close()

	t.Run("direct request to a loopback address is blocked", func(t *testing.T) {
		_, err := designImportHTTPClient.Get(target.URL)
		if err == nil {
			t.Fatal("expected the request to a loopback address to fail, it succeeded")
		}
		if !strings.Contains(err.Error(), "non-public address") {
			t.Errorf("expected a non-public-address error, got: %v", err)
		}
	})

	t.Run("redirect to a loopback address is blocked, not followed", func(t *testing.T) {
		redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, target.URL, http.StatusFound)
		}))
		defer redirector.Close()

		// The redirector itself is also a loopback address, so this fails at
		// the first hop already — which is correct (any address here is
		// off-limits), but it means this case cannot additionally prove the
		// second hop is checked without a real public origin to redirect
		// from. TestSafeOutboundTransport_DialContext below covers the
		// per-dial check directly instead.
		_, err := designImportHTTPClient.Get(redirector.URL)
		if err == nil {
			t.Fatal("expected the request to fail, it succeeded")
		}
	})
}

// TestSafeOutboundTransport_DialContext exercises the DialContext function
// directly against IP-literal addresses. Resolving an IP literal is a
// local, no-network operation (net.DefaultResolver.LookupIP recognizes the
// literal and returns it without a DNS query), so this covers the blocking
// decision for both loopback/private/link-local and public addresses
// without depending on real DNS or real network egress from the test
// environment. For the "allowed" case, the assertion stops at "dialing was
// attempted" (a connection error, not a "non-public address" error) rather
// than requiring a real reachable peer at that IP.
func TestSafeOutboundTransport_DialContext(t *testing.T) {
	transport := safeOutboundTransport()

	blocked := []string{
		"127.0.0.1:80",
		"169.254.169.254:80",
		"10.1.2.3:443",
		"[::1]:80",
	}
	for _, addr := range blocked {
		t.Run("blocks "+addr, func(t *testing.T) {
			conn, err := transport.DialContext(t.Context(), "tcp", addr)
			if err == nil {
				_ = conn.Close()
				t.Fatalf("DialContext(%q) succeeded, want a non-public-address error", addr)
			}
			if !strings.Contains(err.Error(), "non-public address") {
				t.Errorf("DialContext(%q) error = %v, want a non-public-address error", addr, err)
			}
		})
	}

	t.Run("does not reject a public address before dialing", func(t *testing.T) {
		// Port 0 on a public-address literal: the IP passes isPublicUnicastIP,
		// so this fails at the actual TCP dial (a real network error), not at
		// the address-class check — proving public addresses reach the dialer
		// instead of being rejected outright.
		conn, err := transport.DialContext(t.Context(), "tcp", "93.184.216.34:0")
		if err == nil {
			_ = conn.Close()
			t.Fatal("dial to port 0 unexpectedly succeeded")
		}
		if strings.Contains(err.Error(), "non-public address") {
			t.Errorf("a public address was rejected by the address-class check: %v", err)
		}
	})
}
