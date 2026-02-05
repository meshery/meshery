package helpers

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrInvalidURLCode        = "meshery-server-1380"
	ErrInvalidURLSchemeCode  = "meshery-server-1381"
	ErrInternalIPBlockedCode = "meshery-server-1382"
	ErrDNSResolutionCode     = "meshery-server-1383"
)

var (
	ErrInvalidURL = errors.New(
		ErrInvalidURLCode,
		errors.Alert,
		[]string{"Invalid URL format"},
		[]string{"The provided URL could not be parsed"},
		[]string{"Check URL syntax", "Ensure URL is properly encoded"},
		[]string{"https://docs.meshery.io"},
	)

	ErrInvalidURLScheme = errors.New(
		ErrInvalidURLSchemeCode,
		errors.Alert,
		[]string{"Invalid URL scheme"},
		[]string{"Only http and https schemes are allowed"},
		[]string{"Use http:// or https:// URLs"},
		[]string{"https://docs.meshery.io"},
	)

	ErrInternalIPBlocked = errors.New(
		ErrInternalIPBlockedCode,
		errors.Alert,
		[]string{"Internal IP address blocked"},
		[]string{"Requests to internal/private IP addresses are not allowed"},
		[]string{"Use a publicly accessible URL"},
		[]string{"https://docs.meshery.io"},
	)
)

// blockedCIDRs contains private/internal IP ranges that should not be accessible
var blockedCIDRs = []string{
	"127.0.0.0/8",    // Loopback
	"10.0.0.0/8",     // Private Class A
	"172.16.0.0/12",  // Private Class B
	"192.168.0.0/16", // Private Class C
	"169.254.0.0/16", // Link-local (AWS/Cloud metadata)
	"0.0.0.0/8",      // Current network
	"100.64.0.0/10",  // Shared address space (CGN)
	"192.0.0.0/24",   // IETF Protocol assignments
	"192.0.2.0/24",   // TEST-NET-1
	"198.51.100.0/24", // TEST-NET-2
	"203.0.113.0/24", // TEST-NET-3
	"224.0.0.0/4",    // Multicast
	"240.0.0.0/4",    // Reserved
	"::1/128",   // IPv6 loopback
	"fc00::/7",  // IPv6 unique local
	"fe80::/10", // IPv6 link-local
	// Note: IPv4-mapped IPv6 addresses (::ffff:x.x.x.x) are handled by
	// converting to IPv4 in isBlockedIP() and checking against IPv4 CIDRs
}

// parsedBlockedCIDRs holds the parsed CIDR networks
var parsedBlockedCIDRs []*net.IPNet

func init() {
	// Parse CIDRs at init time to avoid repeated parsing
	for _, cidr := range blockedCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			// Log error but continue - don't fail initialization
			fmt.Printf("Warning: Failed to parse CIDR %s: %v\n", cidr, err)
			continue
		}
		parsedBlockedCIDRs = append(parsedBlockedCIDRs, network)
	}
}

// ValidateExternalURL validates that a URL is safe for server-side requests.
// It checks for:
// - Valid URL format
// - Allowed schemes (http/https only)
// - Non-internal/private IP addresses
// - DNS resolution to non-blocked IPs
func ValidateExternalURL(urlStr string) error {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return ErrInvalidURL
	}

	// Only allow http/https schemes
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return ErrInvalidURLScheme
	}

	// Block localhost variations
	hostname := strings.ToLower(parsed.Hostname())
	if hostname == "localhost" || strings.HasSuffix(hostname, ".localhost") {
		return ErrInternalIPBlocked
	}

	// Check if hostname is already an IP
	if ip := net.ParseIP(hostname); ip != nil {
		if isBlockedIP(ip) {
			return ErrInternalIPBlocked
		}
		return nil
	}

	// Resolve DNS with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", hostname)
	if err != nil {
		return errors.New(
			ErrDNSResolutionCode,
			errors.Alert,
			[]string{"DNS resolution failed"},
			[]string{fmt.Sprintf("Could not resolve hostname: %s", hostname)},
			[]string{"Check if the hostname is correct", "Ensure the domain exists"},
			[]string{"https://docs.meshery.io"},
		)
	}

	// Check all resolved IPs
	for _, ip := range ips {
		if isBlockedIP(ip) {
			return ErrInternalIPBlocked
		}
	}

	return nil
}

// isBlockedIP checks if an IP address falls within any blocked CIDR range
func isBlockedIP(ip net.IP) bool {
	// Handle IPv4-mapped IPv6 addresses (e.g., ::ffff:127.0.0.1)
	// Convert to IPv4 if possible for consistent checking
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}

	for _, network := range parsedBlockedCIDRs {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

// SafeHTTPClient returns an HTTP client configured to prevent SSRF attacks.
// It validates redirect URLs before following them.
func SafeHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Limit redirect chain length
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}

			// Validate each redirect URL
			if err := ValidateExternalURL(req.URL.String()); err != nil {
				return fmt.Errorf("redirect blocked: %w", err)
			}

			return nil
		},
	}
}

// SafeGet performs an HTTP GET request with SSRF protection.
// It validates the URL before making the request and prevents redirects to internal IPs.
func SafeGet(urlStr string) (*http.Response, error) {
	// Validate URL before making request
	if err := ValidateExternalURL(urlStr); err != nil {
		return nil, err
	}

	client := SafeHTTPClient(30 * time.Second)
	return client.Get(urlStr)
}
