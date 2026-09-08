package llm

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

// ValidateURL checks if a provider URL is safe to use.
// It guards against SSRF by resolving the hostname and checking all IP addresses.
// If allowLocal is false, private network IPs and loopback are blocked.
func ValidateURL(rawURL string, allowLocal bool) (*url.URL, error) {
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

	if u.User != nil {
		return nil, ErrSSRFValidation(errors.New("credentials must not be embedded in the URL"))
	}

	host := u.Hostname()

	// Resolve the host to IPs to prevent DNS rebinding and evaluate actual targets
	ips, err := net.LookupIP(host)
	if err != nil {
		// If we can't resolve it, we shouldn't trust it.
		// However, for certain local environments, it might just be a local hostname.
		// We'll return an error since DNS resolution is required for SSRF validation.
		return nil, ErrSSRFValidation(fmt.Errorf("could not resolve hostname: %v", err))
	}

	for _, ip := range ips {
		if err := validateIP(ip, allowLocal); err != nil {
			return nil, ErrSSRFValidation(fmt.Errorf("host %s resolved to unsafe IP %s: %v", host, ip.String(), err))
		}
	}

	return u, nil
}

func validateIP(ip net.IP, allowLocal bool) error {
	if ip.IsUnspecified() {
		return errors.New("unspecified IP addresses are not allowed")
	}
	if ip.IsMulticast() || ip.IsLinkLocalMulticast() {
		return errors.New("multicast IP addresses are not allowed")
	}
	if ip.IsLinkLocalUnicast() {
		// This blocks 169.254.x.x (AWS/GCP metadata) and fe80::/10
		return errors.New("link-local/metadata IP addresses are not allowed")
	}

	isPrivateOrLoopback := ip.IsLoopback() || ip.IsPrivate()

	if isPrivateOrLoopback && !allowLocal {
		return errors.New("private/loopback IPs are not allowed for this provider")
	}

	// Block specific known cloud metadata IPs just in case
	metadataIPs := []string{
		"169.254.169.254", // AWS, GCP, Azure
		"169.254.169.253", // GCP
		"169.254.169.250", // GCP
		"100.100.100.200", // Alibaba
	}
	for _, mIP := range metadataIPs {
		if ip.Equal(net.ParseIP(mIP)) {
			return errors.New("metadata endpoints are strictly forbidden")
		}
	}

	return nil
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
