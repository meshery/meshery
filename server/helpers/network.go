package helpers

import (
	"fmt"
	"net"
	"net/url"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrInvalidURLCode        = "meshery-server-1170"
	ErrBlockedIPCode         = "meshery-server-1171"
	ErrURLResolutionCode     = "meshery-server-1172"
)

func ErrInvalidURL(err error) error {
	return errors.New(ErrInvalidURLCode, errors.Alert,
		[]string{"Invalid or unsafe URL provided"},
		[]string{err.Error()},
		[]string{"The URL may have an unsupported scheme or missing host"},
		[]string{"Ensure the URL uses http or https and points to a valid external host"},
	)
}

func ErrBlockedIP(err error) error {
	return errors.New(ErrBlockedIPCode, errors.Alert,
		[]string{"Request to private or internal address blocked"},
		[]string{err.Error()},
		[]string{"The URL resolves to a loopback, private, or link-local address"},
		[]string{"Use a publicly accessible URL instead"},
	)
}

func ErrURLResolution(err error) error {
	return errors.New(ErrURLResolutionCode, errors.Alert,
		[]string{"Unable to resolve hostname in URL"},
		[]string{err.Error()},
		[]string{"The hostname may be invalid or DNS resolution failed"},
		[]string{"Check that the URL hostname is correct and reachable"},
	)
}

// ValidateURL checks that a URL is safe to use in outbound HTTP requests.
// It rejects loopback, private, and link-local addresses to prevent SSRF.
func ValidateURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ErrInvalidURL(err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrInvalidURL(fmt.Errorf("unsupported scheme: %s", parsed.Scheme))
	}

	hostname := parsed.Hostname()
	if hostname == "" {
		return ErrInvalidURL(fmt.Errorf("missing host in URL"))
	}

	ips, err := net.LookupHost(hostname)
	if err != nil {
		return ErrURLResolution(err)
	}

	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			continue
		}
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return ErrBlockedIP(fmt.Errorf("host %s resolves to a private address %s", hostname, ipStr))
		}
	}

	return nil
}