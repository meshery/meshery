package helpers

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestValidateExternalURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
		errType string
	}{
		// Valid URLs
		{
			name:    "valid https URL",
			url:     "https://example.com/path",
			wantErr: false,
		},
		{
			name:    "valid http URL",
			url:     "http://example.com/path",
			wantErr: false,
		},
		{
			name:    "valid URL with port",
			url:     "https://example.com:8080/path",
			wantErr: false,
		},

		// Invalid schemes
		{
			name:    "file scheme blocked",
			url:     "file:///etc/passwd",
			wantErr: true,
			errType: "scheme",
		},
		{
			name:    "ftp scheme blocked",
			url:     "ftp://example.com/file",
			wantErr: true,
			errType: "scheme",
		},
		{
			name:    "gopher scheme blocked",
			url:     "gopher://example.com/",
			wantErr: true,
			errType: "scheme",
		},

		// Localhost variations
		{
			name:    "localhost blocked",
			url:     "http://localhost/admin",
			wantErr: true,
			errType: "internal",
		},
		{
			name:    "localhost with port blocked",
			url:     "http://localhost:8080/admin",
			wantErr: true,
			errType: "internal",
		},
		{
			name:    "subdomain of localhost blocked",
			url:     "http://api.localhost/admin",
			wantErr: true,
			errType: "internal",
		},

		// Loopback IPs
		{
			name:    "127.0.0.1 blocked",
			url:     "http://127.0.0.1/admin",
			wantErr: true,
			errType: "internal",
		},
		{
			name:    "127.0.0.1 with port blocked",
			url:     "http://127.0.0.1:8080/admin",
			wantErr: true,
			errType: "internal",
		},
		{
			name:    "127.x.x.x range blocked",
			url:     "http://127.1.2.3/admin",
			wantErr: true,
			errType: "internal",
		},

		// Private networks - Class A (10.x.x.x)
		{
			name:    "10.0.0.1 blocked",
			url:     "http://10.0.0.1/internal",
			wantErr: true,
			errType: "internal",
		},
		{
			name:    "10.255.255.255 blocked",
			url:     "http://10.255.255.255/internal",
			wantErr: true,
			errType: "internal",
		},

		// Private networks - Class B (172.16-31.x.x)
		{
			name:    "172.16.0.1 blocked",
			url:     "http://172.16.0.1/internal",
			wantErr: true,
			errType: "internal",
		},
		{
			name:    "172.31.255.255 blocked",
			url:     "http://172.31.255.255/internal",
			wantErr: true,
			errType: "internal",
		},

		// Private networks - Class C (192.168.x.x)
		{
			name:    "192.168.1.1 blocked",
			url:     "http://192.168.1.1/internal",
			wantErr: true,
			errType: "internal",
		},
		{
			name:    "192.168.0.1 blocked",
			url:     "http://192.168.0.1/internal",
			wantErr: true,
			errType: "internal",
		},

		// Cloud metadata endpoints
		{
			name:    "AWS metadata endpoint blocked",
			url:     "http://169.254.169.254/latest/meta-data/",
			wantErr: true,
			errType: "internal",
		},
		{
			name:    "Link-local address blocked",
			url:     "http://169.254.1.1/",
			wantErr: true,
			errType: "internal",
		},

		// IPv6 addresses
		{
			name:    "IPv6 loopback blocked",
			url:     "http://[::1]/admin",
			wantErr: true,
			errType: "internal",
		},

		// Invalid URLs
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
			errType: "invalid",
		},
		{
			name:    "malformed URL",
			url:     "not-a-valid-url",
			wantErr: true,
			errType: "scheme",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExternalURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateExternalURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestIsBlockedIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		blocked bool
	}{
		// Loopback
		{"loopback 127.0.0.1", "127.0.0.1", true},
		{"loopback 127.255.255.255", "127.255.255.255", true},

		// Private Class A
		{"private 10.0.0.1", "10.0.0.1", true},
		{"private 10.255.255.255", "10.255.255.255", true},

		// Private Class B
		{"private 172.16.0.1", "172.16.0.1", true},
		{"private 172.31.255.255", "172.31.255.255", true},
		{"public 172.15.0.1", "172.15.0.1", false}, // Just outside range
		{"public 172.32.0.1", "172.32.0.1", false}, // Just outside range

		// Private Class C
		{"private 192.168.0.1", "192.168.0.1", true},
		{"private 192.168.255.255", "192.168.255.255", true},

		// Link-local / Cloud metadata
		{"metadata 169.254.169.254", "169.254.169.254", true},
		{"link-local 169.254.1.1", "169.254.1.1", true},

		// Public IPs (should not be blocked)
		{"public 8.8.8.8", "8.8.8.8", false},
		{"public 1.1.1.1", "1.1.1.1", false},
		{"public 93.184.216.34", "93.184.216.34", false}, // example.com

		// IPv6
		{"ipv6 loopback", "::1", true},
		{"ipv6 public", "2001:4860:4860::8888", false}, // Google DNS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := parseIP(tt.ip)
			if ip == nil {
				t.Fatalf("Failed to parse IP: %s", tt.ip)
			}
			if got := isBlockedIP(ip); got != tt.blocked {
				t.Errorf("isBlockedIP(%s) = %v, want %v", tt.ip, got, tt.blocked)
			}
		})
	}
}

func TestSafeHTTPClientRedirectValidation(t *testing.T) {
	// Create a test server that redirects to an internal IP
	redirectCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectCount++
		if redirectCount == 1 {
			// First request - redirect to internal IP
			http.Redirect(w, r, "http://127.0.0.1/admin", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := SafeHTTPClient(5 * time.Second)

	// This should fail because the redirect goes to an internal IP
	_, err := client.Get(server.URL)
	if err == nil {
		t.Error("Expected error for redirect to internal IP, got nil")
	}
}

// TestDNSRebindingProtection tests that hostnames resolving to blocked IPs are rejected.
// This simulates a DNS rebinding attack where an attacker controls a domain that
// resolves to internal IPs.
func TestDNSRebindingProtection(t *testing.T) {
	// These domains are known to resolve to 127.0.0.1 for testing purposes
	// They simulate an attacker-controlled domain pointing to internal IPs
	rebindingDomains := []struct {
		name   string
		domain string
	}{
		{"localtest.me resolves to 127.0.0.1", "http://localtest.me/"},
		{"vcap.me resolves to 127.0.0.1", "http://vcap.me/"},
		{"127.0.0.1.nip.io resolves to 127.0.0.1", "http://127.0.0.1.nip.io/"},
	}

	for _, tt := range rebindingDomains {
		t.Run(tt.name, func(t *testing.T) {
			// First verify the domain actually resolves to a blocked IP
			// Skip test if DNS doesn't work as expected (network issues, etc.)
			ips, err := net.LookupIP(tt.domain[7 : len(tt.domain)-1]) // Extract hostname
			if err != nil {
				t.Skipf("Skipping test: DNS lookup failed for %s: %v", tt.domain, err)
				return
			}

			// Check if any resolved IP is blocked (should be 127.0.0.1)
			hasBlockedIP := false
			for _, ip := range ips {
				if isBlockedIP(ip) {
					hasBlockedIP = true
					break
				}
			}

			if !hasBlockedIP {
				t.Skipf("Skipping test: %s does not resolve to a blocked IP", tt.domain)
				return
			}

			// Now test that ValidateExternalURL blocks this domain
			err = ValidateExternalURL(tt.domain)
			if err == nil {
				t.Errorf("ValidateExternalURL(%s) should have blocked DNS rebinding attack, got nil error", tt.domain)
			}
		})
	}
}

// TestIPv4MappedIPv6Blocking tests that IPv4-mapped IPv6 addresses are properly blocked.
// Attackers may try to bypass filters using ::ffff:127.0.0.1 notation.
func TestIPv4MappedIPv6Blocking(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		blocked bool
	}{
		{"IPv4-mapped loopback", "::ffff:127.0.0.1", true},
		{"IPv4-mapped private 10.x", "::ffff:10.0.0.1", true},
		{"IPv4-mapped private 192.168.x", "::ffff:192.168.1.1", true},
		{"IPv4-mapped metadata", "::ffff:169.254.169.254", true},
		{"IPv4-mapped public", "::ffff:8.8.8.8", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if ip == nil {
				t.Fatalf("Failed to parse IP: %s", tt.ip)
			}
			if got := isBlockedIP(ip); got != tt.blocked {
				t.Errorf("isBlockedIP(%s) = %v, want %v", tt.ip, got, tt.blocked)
			}
		})
	}
}

// Helper function to parse IP for tests
func parseIP(s string) net.IP {
	return net.ParseIP(s)
}
