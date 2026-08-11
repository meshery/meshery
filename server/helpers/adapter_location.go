package helpers

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// AdapterPort returns the TCP port from an adapter location string.
// Supported forms: host:port, http(s)://host:port, and [ipv6]:port.
// Returns an error instead of panicking when the location has no port
// or is otherwise malformed.
func AdapterPort(location string) (string, error) {
	hostPort := strings.TrimSpace(location)
	if hostPort == "" {
		return "", fmt.Errorf("adapter location is empty")
	}

	if strings.Contains(hostPort, "://") {
		u, err := url.Parse(hostPort)
		if err != nil {
			return "", fmt.Errorf("invalid adapter location %q: %w", location, err)
		}
		if u.Host == "" {
			return "", fmt.Errorf("adapter location %q has no host", location)
		}
		hostPort = u.Host
	}

	_, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", fmt.Errorf("invalid adapter location %q: %w", location, err)
	}
	if port == "" {
		return "", fmt.Errorf("adapter location %q has no port", location)
	}
	return port, nil
}
