package utils

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// cgnatNet covers Carrier-Grade NAT (RFC 6598: 100.64.0.0/10).
	cgnatNet = net.IPNet{
		IP:   net.ParseIP("100.64.0.0"),
		Mask: net.CIDRMask(10, 32),
	}
	// benchNet covers Benchmarking (RFC 2544: 198.18.0.0/15).
	benchNet = net.IPNet{
		IP:   net.ParseIP("198.18.0.0"),
		Mask: net.CIDRMask(15, 32),
	}
)

// IsRestrictedIP checks whether an IP address belongs to a private, loopback,
// link-local, multicast, unspecified, CGNAT, or benchmarking range.
// It explicitly unpacks IPv4-mapped IPv6 addresses so addresses such as
// `::ffff:127.0.0.1` or `::ffff:169.254.169.254` are evaluated accurately.
func IsRestrictedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}

	// Unwrap IPv4-mapped IPv6 addresses to 4-byte IPv4 slices.
	if ip4 := ip.To4(); ip4 != nil {
		if ip4.IsLoopback() ||
			ip4.IsPrivate() ||
			ip4.IsLinkLocalUnicast() ||
			ip4.IsLinkLocalMulticast() ||
			ip4.IsMulticast() ||
			ip4.IsUnspecified() ||
			ip4[0] == 0 ||
			cgnatNet.Contains(ip4) ||
			benchNet.Contains(ip4) {
			return true
		}
		return false
	}

	// Pure IPv6 checks
	if ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() ||
		ip.IsUnspecified() {
		return true
	}

	return false
}

// IsRestrictedHost checks whether a hostname matches known internal or cloud
// metadata domains (e.g. localhost, metadata.google.internal).
func IsRestrictedHost(host string) bool {
	h := strings.ToLower(strings.TrimSpace(host))
	// Strip port if present
	if hostname, _, err := net.SplitHostPort(h); err == nil {
		h = hostname
	}
	h = strings.TrimSuffix(h, ".")

	if h == "" ||
		h == "localhost" || strings.HasSuffix(h, ".localhost") ||
		h == "metadata.google.internal" || strings.HasSuffix(h, ".metadata.google.internal") ||
		h == "metadata" || h == "instance-data" {
		return true
	}
	return false
}

// ValidateSafeURL ensures that rawURL has an allowed scheme (http or https),
// a non-empty host, is not a restricted hostname, and does not parse as a
// restricted IP literal.
func ValidateSafeURL(rawURL string) (*url.URL, error) {
	if strings.TrimSpace(rawURL) == "" {
		return nil, errors.New("URL cannot be empty")
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme %q: only http and https are permitted", parsed.Scheme)
	}

	hostname := parsed.Hostname()
	if hostname == "" {
		return nil, errors.New("URL host cannot be empty")
	}

	if IsRestrictedHost(hostname) {
		return nil, fmt.Errorf("access to host %q is restricted", hostname)
	}

	// If the host is an IP literal, validate it immediately.
	if ip := net.ParseIP(hostname); ip != nil {
		if IsRestrictedIP(ip) {
			return nil, fmt.Errorf("access to IP address %q is restricted", hostname)
		}
	}

	return parsed, nil
}

// SafeDialContext returns a dial function that resolves hostnames and verifies
// that every returned IP address is safe/public before establishing a TCP connection.
// If allowPrivate is true (for testing only), private/loopback address filtering is bypassed.
func SafeDialContext(timeout time.Duration, allowPrivate bool) func(ctx context.Context, network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: 30 * time.Second,
	}

	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("invalid network address %q: %w", addr, err)
		}

		if !allowPrivate && IsRestrictedHost(host) {
			return nil, fmt.Errorf("connection to host %q is restricted", host)
		}

		// Resolve IP addresses
		ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
		if err != nil {
			return nil, fmt.Errorf("DNS resolution failed for %q: %w", host, err)
		}
		if len(ips) == 0 {
			return nil, fmt.Errorf("no IP addresses found for host %q", host)
		}

		if !allowPrivate {
			// Verify every resolved IP. If ANY resolved IP is in a restricted range, abort.
			for _, ip := range ips {
				if IsRestrictedIP(ip) {
					return nil, fmt.Errorf("connection to restricted IP %q for host %q is forbidden", ip.String(), host)
				}
			}
		}

		// Connect to the first resolved safe IP to pin the connection to the validated address.
		targetAddr := net.JoinHostPort(ips[0].String(), port)
		return dialer.DialContext(ctx, network, targetAddr)
	}
}

// SafeCheckRedirect validates each hop in a redirect chain to prevent an attacker
// from redirecting a benign URL to a private/internal endpoint.
func SafeCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errors.New("stopped after 10 redirects")
	}
	if req == nil || req.URL == nil {
		return errors.New("nil redirect request")
	}
	_, err := ValidateSafeURL(req.URL.String())
	if err != nil {
		return fmt.Errorf("redirect blocked: %w", err)
	}
	return nil
}

// NewSafeHTTPClientWithAllowPrivate returns an *http.Client bounded by timeout
// and configured with a transport and redirect checker that block SSRF attempts.
// Setting allowPrivate to true should only be used in unit tests that spin up
// local httptest.Server instances.
func NewSafeHTTPClientWithAllowPrivate(timeout time.Duration, allowPrivate bool) *http.Client {
	transport := &http.Transport{
		Proxy:                 nil, // No proxy: a proxy would resolve the destination itself, bypassing SafeDialContext's destination validation.
		DialContext:           SafeDialContext(timeout, allowPrivate),
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport:     transport,
		Timeout:       timeout,
		CheckRedirect: SafeCheckRedirect,
	}
}

// NewSafeHTTPClient returns an *http.Client bounded by timeout and configured
// with a transport and redirect checker that block SSRF attempts.
func NewSafeHTTPClient(timeout time.Duration) *http.Client {
	return NewSafeHTTPClientWithAllowPrivate(timeout, false)
}
