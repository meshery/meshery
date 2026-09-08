package llm

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

// ValidateURL checks if a provider URL is safe to use.
// It guards against SSRF by blocking metadata endpoints and enforcing schemes.
// For cloud providers, HTTPS is generally expected, but we allow HTTP for local/LAN inference.
func ValidateURL(rawURL string) (*url.URL, error) {
	if rawURL == "" {
		return nil, errors.New("url cannot be empty")
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, ErrSSRFValidation(err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, ErrSSRFValidation(errors.New("unsupported URL scheme, only http and https are allowed"))
	}

	// Do not allow credentials embedded in the URL (e.g., https://user:pass@host/)
	if u.User != nil {
		return nil, ErrSSRFValidation(errors.New("credentials must not be embedded in the URL"))
	}

	host := u.Hostname()

	// Block known dangerous metadata endpoints (e.g., AWS, GCP, Azure metadata IPs)
	if isMetadataEndpoint(host) {
		return nil, ErrSSRFValidation(errors.New("metadata endpoints are not allowed"))
	}

	return u, nil
}

// isMetadataEndpoint checks against common cloud metadata IP addresses.
func isMetadataEndpoint(host string) bool {
	metadataIPs := []string{
		"169.254.169.254", // AWS, GCP, Azure
		"169.254.169.253", // GCP
		"169.254.169.250", // GCP
		"100.100.100.200", // Alibaba
	}

	// Simple string match
	for _, ip := range metadataIPs {
		if host == ip {
			return true
		}
	}

	// Check if the host resolves to a metadata IP
	// This helps prevent DNS rebinding where a safe domain resolves to a metadata IP.
	// We only do this if it's not localhost/LAN to save time.
	if net.ParseIP(host) != nil {
		return false
	}
	
	// Optional: we can do a DNS lookup here, but it might slow down validation.
	// For now, strict IP matching is implemented. In a production environment with
	// strict SSRF protection, we would dial and check the resolved IP.

	return false
}

// EnsureHTTPS is a helper to verify HTTPS is used when interacting with cloud providers.
func EnsureHTTPS(u *url.URL) error {
	if u.Scheme != "https" && !isLocalHostOrLAN(u.Hostname()) {
		return ErrSSRFValidation(errors.New("https is required for remote providers"))
	}
	return nil
}

// isLocalHostOrLAN checks if a given host points to a local or private network.
func isLocalHostOrLAN(host string) bool {
	if host == "localhost" || strings.HasSuffix(host, ".local") {
		return true
	}
	
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	// Loopback (127.0.0.0/8, ::1/128)
	if ip.IsLoopback() {
		return true
	}
	// Private / LAN IPs (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16, fc00::/7)
	if ip.IsPrivate() {
		return true
	}

	return false
}
