package models

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	serverutils "github.com/meshery/meshery/server/helpers/utils"
)

const (
	// maxRemoteFileSize bounds response bodies for filter and pattern URL imports.
	maxRemoteFileSize int64 = 10 * 1024 * 1024

	// remoteFileFetchTimeout bounds the whole request, including redirects and
	// reading the response body.
	remoteFileFetchTimeout = 30 * time.Second

	// maxRemoteFileRedirects is the number of redirects a fetch may follow.
	maxRemoteFileRedirects = 5
)

// newRemoteFileClient builds the HTTP client used to fetch user-supplied URLs.
// The shared safe transport validates resolved addresses and pins the dial to a
// validated public IP. This client adds the redirect limit used by filter and
// pattern imports. allowPrivate is only enabled by tests using httptest servers.
func newRemoteFileClient(timeout time.Duration, allowPrivate bool) *http.Client {
	client := serverutils.NewSafeHTTPClientWithAllowPrivate(timeout, allowPrivate)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// via includes the original request, so len(via)-1 is the redirect count.
		if len(via) > maxRemoteFileRedirects {
			return fmt.Errorf("stopped after %d redirects", maxRemoteFileRedirects)
		}
		if req == nil || req.URL == nil {
			return fmt.Errorf("redirect request has no URL")
		}
		if _, err := serverutils.ValidateSafeURL(req.URL.String()); err != nil {
			return fmt.Errorf("redirect blocked: %w", err)
		}
		return nil
	}
	return client
}

func validateRemoteFileURL(u *url.URL) error {
	if u == nil {
		return fmt.Errorf("URL is missing")
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("unsupported URL scheme %q: only http and https are allowed", u.Scheme)
	}
	if u.Hostname() == "" {
		return fmt.Errorf("URL has no host")
	}
	return nil
}

var remoteFileFetchClient = newRemoteFileClient(remoteFileFetchTimeout, false)

// SetRemoteFileFetchClientForTest overrides the client and returns a function
// that restores the previous client. Use only in tests with local HTTP servers.
func SetRemoteFileFetchClientForTest(client *http.Client) func() {
	previous := remoteFileFetchClient
	remoteFileFetchClient = client
	return func() { remoteFileFetchClient = previous }
}

// fetchRemoteFile downloads a user-supplied URL for import. It only uses public
// HTTP(S) destinations, follows a bounded number of redirects, times out, and
// reads at most maxRemoteFileSize bytes.
func fetchRemoteFile(rawURL string) ([]byte, error) {
	return fetchRemoteFileWith(remoteFileFetchClient, rawURL, maxRemoteFileSize)
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
