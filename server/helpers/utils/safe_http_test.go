package utils

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestIsRestrictedIP(t *testing.T) {
	tests := []struct {
		name       string
		ipStr      string
		restricted bool
	}{
		// Loopback
		{"IPv4 Loopback 127.0.0.1", "127.0.0.1", true},
		{"IPv4 Loopback 127.1.2.3", "127.1.2.3", true},
		{"IPv6 Loopback ::1", "::1", true},
		{"IPv4-mapped IPv6 Loopback", "::ffff:127.0.0.1", true},

		// Link-Local (Cloud Metadata AWS / Azure / GCP)
		{"AWS/GCP IMDS 169.254.169.254", "169.254.169.254", true},
		{"AWS ECS Metadata 169.254.170.2", "169.254.170.2", true},
		{"IPv4 Link-local generic", "169.254.1.1", true},
		{"IPv4-mapped IPv6 Link-local", "::ffff:169.254.169.254", true},
		{"IPv6 Link-local fe80::1", "fe80::1", true},

		// Private Networks (RFC 1918 & RFC 4193)
		{"RFC 1918 Class A 10.0.0.1", "10.0.0.1", true},
		{"RFC 1918 Class A 10.255.255.254", "10.255.255.254", true},
		{"RFC 1918 Class B 172.16.0.1", "172.16.0.1", true},
		{"RFC 1918 Class B 172.31.255.254", "172.31.255.254", true},
		{"RFC 1918 Class C 192.168.1.1", "192.168.1.1", true},
		{"RFC 4193 IPv6 ULA fc00::1", "fc00::1", true},
		{"RFC 4193 IPv6 ULA fd12:3456:789a::1", "fd12:3456:789a::1", true},

		// Unspecified / Broadcast / 0.0.0.0/8
		{"IPv4 Unspecified 0.0.0.0", "0.0.0.0", true},
		{"Current network 0.0.0.0/8 0.1.2.3", "0.1.2.3", true},
		{"IPv6 Unspecified ::", "::", true},

		// CGNAT / Benchmarking
		{"CGNAT 100.64.0.1", "100.64.0.1", true},
		{"CGNAT 100.127.255.254", "100.127.255.254", true},
		{"Benchmarking 198.18.0.1", "198.18.0.1", true},

		// Multicast
		{"IPv4 Multicast 224.0.0.1", "224.0.0.1", true},
		{"IPv6 Multicast ff02::1", "ff02::1", true},

		// Public & Safe IPs
		{"Google Public DNS 8.8.8.8", "8.8.8.8", false},
		{"Cloudflare DNS 1.1.1.1", "1.1.1.1", false},
		{"GitHub IP 140.82.121.4", "140.82.121.4", false},
		{"Cloudflare IPv6 2606:4700:4700::1111", "2606:4700:4700::1111", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ipStr)
			if ip == nil {
				t.Fatalf("failed to parse IP %q", tt.ipStr)
			}
			got := IsRestrictedIP(ip)
			if got != tt.restricted {
				t.Errorf("IsRestrictedIP(%q) = %v, want %v", tt.ipStr, got, tt.restricted)
			}
		})
	}
}

func TestIsRestrictedHost(t *testing.T) {
	tests := []struct {
		host       string
		restricted bool
	}{
		{"localhost", true},
		{"LOCALHOST", true},
		{"localhost.", true},
		{"test.localhost", true},
		{"localhost:8080", true},
		{"metadata.google.internal", true},
		{"sub.metadata.google.internal", true},
		{"metadata.google.internal:80", true},
		{"metadata", true},
		{"instance-data", true},
		{"example.com", false},
		{"raw.githubusercontent.com", false},
		{"meshery.io", false},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			got := IsRestrictedHost(tt.host)
			if got != tt.restricted {
				t.Errorf("IsRestrictedHost(%q) = %v, want %v", tt.host, got, tt.restricted)
			}
		})
	}
}

func TestValidateSafeURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		// Valid URLs
		{"Valid HTTPS URL", "https://raw.githubusercontent.com/meshery/meshery/master/README.md", false},
		{"Valid HTTP URL", "http://example.com/manifest.yaml", false},

		// Invalid Schemes
		{"File scheme", "file:///etc/passwd", true},
		{"FTP scheme", "ftp://example.com/file.yaml", true},
		{"Gopher scheme", "gopher://127.0.0.1:70", true},
		{"Javascript scheme", "javascript:alert(1)", true},

		// Malformed URLs
		{"Empty string", "", true},
		{"Empty host", "http:///path", true},
		{"Whitespace only", "   ", true},

		// Restricted Hosts
		{"Localhost host", "http://localhost/api", true},
		{"Subdomain of localhost", "http://dev.localhost/manifest.yaml", true},
		{"GCP metadata domain", "http://metadata.google.internal/computeMetadata/v1/", true},
		{"AWS instance-data domain", "http://instance-data/latest/meta-data/", true},

		// Restricted IP Literals
		{"Loopback IPv4 literal", "http://127.0.0.1:9081/api", true},
		{"Loopback IPv6 literal", "http://[::1]:9081/api", true},
		{"AWS/GCP IMDS IP literal", "http://169.254.169.254/latest/meta-data/", true},
		{"IPv4-mapped IPv6 IMDS", "http://[::ffff:169.254.169.254]/latest/meta-data/", true},
		{"Private class A literal", "http://10.0.0.1/design.yaml", true},
		{"Private class B literal", "http://172.16.0.1/design.yaml", true},
		{"Private class C literal", "http://192.168.1.1/design.yaml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateSafeURL(tt.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSafeURL(%q) err = %v, wantErr %v", tt.rawURL, err, tt.wantErr)
			}
		})
	}
}

func TestSafeHTTPClient_BlocksLoopback(t *testing.T) {
	// A local test server listens on 127.0.0.1 (loopback).
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("internal secret data"))
	}))
	defer server.Close()

	client := NewSafeHTTPClient(2 * time.Second)

	// Fetching server.URL must fail because 127.0.0.1 is restricted.
	resp, err := client.Get(server.URL)
	if err == nil {
		if resp != nil {
			_ = resp.Body.Close()
		}
		t.Fatalf("expected NewSafeHTTPClient to block request to %s, but request succeeded", server.URL)
	}

	// Fetching with allowPrivate = true must succeed (verifying the test harness works).
	testClient := NewSafeHTTPClientWithAllowPrivate(2*time.Second, true)
	testResp, err := testClient.Get(server.URL)
	if err != nil {
		t.Fatalf("expected allowPrivate client to connect to %s, got error: %v", server.URL, err)
	}
	defer func() {
		_ = testResp.Body.Close()
	}()
	if testResp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", testResp.StatusCode)
	}
}

func TestSafeHTTPClient_BlocksRedirectToRestrictedTarget(t *testing.T) {
	// Target server on loopback
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("sensitive metadata"))
	}))
	defer targetServer.Close()

	// Redirect server that bounces requests to the loopback server
	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, targetServer.URL, http.StatusFound)
	}))
	defer redirectServer.Close()

	// Use the allowPrivate client so the first hop connects to the local
	// redirect server, and the redirect to the loopback target is what gets blocked.
	client := NewSafeHTTPClientWithAllowPrivate(2*time.Second, true)
	_, err := client.Get(redirectServer.URL)
	if err == nil {
		t.Fatal("expected request following redirect to restricted host to fail, but it succeeded")
	}
	if !strings.Contains(err.Error(), "redirect blocked") && !strings.Contains(err.Error(), "restricted") {
		t.Errorf("expected redirect blocking error, got: %v", err)
	}
}

func TestSafeHTTPClient_DisablesEnvironmentProxy(t *testing.T) {
	client := NewSafeHTTPClient(2 * time.Second)
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", client.Transport)
	}
	if transport.Proxy != nil {
		t.Errorf("expected transport.Proxy to be nil to prevent proxy bypass, got non-nil")
	}
}
