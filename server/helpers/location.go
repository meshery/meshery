package helpers

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ParseAdapterLocation safely extracts host and port from location strings
// such as "localhost:10000", "http://localhost:11434", or "[::1]:8080".
func ParseAdapterLocation(location string) (host string, port string, err error) {
	location = strings.TrimSpace(location)
	if location == "" {
		return "", "", fmt.Errorf("adapter location cannot be empty")
	}

	rawURL := location
	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid adapter location format: %w", err)
	}

	host, port, err = net.SplitHostPort(parsedURL.Host)
	if err != nil {
		// No port explicitly provided in location string
		return parsedURL.Hostname(), "", nil
	}

	return host, port, nil
}
