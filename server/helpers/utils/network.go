package utils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"time"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrInvalidURLCode    = "meshery-server-1443"
	ErrBlockedIPCode     = "meshery-server-1444"
	ErrURLResolutionCode = "meshery-server-1445"
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

// NewSafeHTTPClient returns an *http.Client that validates the resolved IP
// at dial time, preventing SSRF and DNS rebinding attacks by blocking
// loopback, private, and link-local addresses.
// DialControl is the dial-time guard used by NewSafeHTTPClient to block
// connections to private/internal IPs. It's a package-level variable so
// tests can override it to permit dialing httptest servers (which bind to
// loopback addresses), while production always uses safeDialControl.
var DialControl = safeDialControl

func safeDialControl(network, address string, c syscall.RawConn) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", host)
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return fmt.Errorf("connection to private/internal IP %s is blocked", ip)
	}
	return nil
}

func NewSafeHTTPClient(timeout time.Duration) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = (&net.Dialer{
		Timeout:   timeout,
		KeepAlive: 30 * time.Second,
		Control:   DialControl,
	}).DialContext
	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

// URLValidator is the function outbound HTTP clients use to validate
// user-supplied URLs before connecting. It defaults to
// ValidateURLForOutboundRequest. Tests may reassign it (and must restore
// the original afterward) to allow connecting to httptest servers, which
// bind to loopback addresses that are correctly rejected in production.
var URLValidator = ValidateURLForOutboundRequest

func ValidateURLForOutboundRequest(rawURL string) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ips, err := net.DefaultResolver.LookupHost(ctx, hostname)
	if err != nil {
		return ErrURLResolution(err)
	}

	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			continue
		}
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return ErrBlockedIP(fmt.Errorf("host %s resolves to an unsafe address %s", hostname, ipStr))
		}
	}

	return nil
}
