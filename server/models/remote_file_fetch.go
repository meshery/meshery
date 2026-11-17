package models

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"syscall"
	"time"

)

const (
	// maxRemoteFileSize bounds how much of a user-supplied URL's response body is
	// read into memory when importing a design or filter by URL (10 MiB).
	maxRemoteFileSize int64 = 10 * 1024 * 1024

	// remoteFileFetchTimeout bounds the whole request, including redirects and
	// reading the body.
	remoteFileFetchTimeout = 30 * time.Second

	// maxRemoteFileRedirects is the number of redirects a fetch may follow.
	maxRemoteFileRedirects = 5
)

// blockedRemoteFilePrefixes lists address ranges that are not covered by the
// netip.Addr predicates used in isBlockedRemoteFileAddr.
var blockedRemoteFilePrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),       // "this network"
	netip.MustParsePrefix("100.64.0.0/10"),   // carrier-grade NAT
	netip.MustParsePrefix("192.0.0.0/24"),    // IETF protocol assignments
	netip.MustParsePrefix("198.18.0.0/15"),   // benchmarking
	netip.MustParsePrefix("240.0.0.0/4"),     // reserved
	netip.MustParsePrefix("192.0.2.0/24"),    // TEST-NET-1 documentation
	netip.MustParsePrefix("198.51.100.0/24"), // TEST-NET-2 documentation
	netip.MustParsePrefix("203.0.113.0/24"),  // TEST-NET-3 documentation
	netip.MustParsePrefix("2001:db8::/32"),   // IPv6 documentation
}

// isBlockedRemoteFileAddr reports whether the server must refuse to connect to
// addr when fetching a user-supplied URL: loopback, private, link-local (which
// includes cloud metadata endpoints such as 169.254.169.254), multicast,
// unspecified and other non-public ranges.
func isBlockedRemoteFileAddr(addr netip.Addr) bool {
	// Drop any IPv6 zone first: netip's predicates compare against zoneless
	// addresses, so "::%lo" would otherwise not count as unspecified even though
	// connecting to it reaches loopback. Then unmap so ::ffff:127.0.0.1 is
	// judged as 127.0.0.1.
	addr = addr.WithZone("").Unmap()
	if !addr.IsValid() ||
		addr.IsLoopback() ||
		addr.IsPrivate() ||
		addr.IsLinkLocalUnicast() ||
		addr.IsLinkLocalMulticast() ||
		addr.IsInterfaceLocalMulticast() ||
		addr.IsMulticast() ||
		addr.IsUnspecified() {
		return true
	}
	for _, prefix := range blockedRemoteFilePrefixes {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}

// newRemoteFileClient builds the HTTP client used to fetch user-supplied URLs.
//
// The destination address is validated in the dialer's Control hook, i.e. on the
// IP address actually about to be connected to. That covers every redirect hop
// and defeats DNS rebinding, which a check of the URL host before the request
// would not. Proxies from the environment are deliberately not used: through a
// proxy the dialer only sees the proxy's address, so the check would not apply
// to the real destination.
func newRemoteFileClient(timeout time.Duration, blocked func(netip.Addr) bool) *http.Client {
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
		Control: func(_, address string, _ syscall.RawConn) error {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				return err
			}
			addr, err := netip.ParseAddr(host)
			if err != nil {
				return fmt.Errorf("refusing to connect to %q: not an IP address", host)
			}
			if blocked(addr) {
				return fmt.Errorf("refusing to connect to non-public address %s", addr)
			}
			return nil
		},
	}

	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy:                 nil,
			DialContext:           dialer.DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: timeout,
			MaxIdleConns:          1,
			DisableKeepAlives:     true,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// via holds every request made so far, including the original one,
			// so len(via) is the number of redirects already followed plus one.
			if len(via) > maxRemoteFileRedirects {
				return fmt.Errorf("stopped after %d redirects", maxRemoteFileRedirects)
			}
			return validateRemoteFileURL(req.URL)
		},
	}
}

// validateRemoteFileURL accepts only absolute http and https URLs.
func validateRemoteFileURL(u *url.URL) error {
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme %q: only http and https are allowed", u.Scheme)
	}
	if u.Hostname() == "" {
		return fmt.Errorf("URL has no host")
	}
	return nil
}

// fetchRemoteFile downloads a user-supplied URL for import. It only talks to
// public http(s) destinations (unless ALLOW_PRIVATE_REMOTE_FILE_IMPORTS is set),
// follows a bounded number of redirects, times out,
// and reads at most maxRemoteFileSize bytes.
func fetchRemoteFile(rawURL string) ([]byte, error) {
	return fetchRemoteFileWith(newRemoteFileClient(remoteFileFetchTimeout, isBlockedRemoteFileAddr), rawURL, maxRemoteFileSize)
}

func fetchRemoteFileWith(client *http.Client, rawURL string, maxBytes int64) ([]byte, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if err := validateRemoteFileURL(parsed); err != nil {
		return nil, err
	}

	resp, err := client.Get(parsed.String())
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("file not found")
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > maxBytes {
		return nil, fmt.Errorf("remote file exceeds the maximum allowed size of %d bytes", maxBytes)
	}
	return body, nil
}
