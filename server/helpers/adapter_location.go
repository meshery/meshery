package helpers

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// AdapterPort returns the TCP port from an adapter location string.
// Supported forms: host:port, http(s)://host:port, and [ipv6]:port.
// Returns an error instead of panicking when the location has no port
// or is otherwise malformed. Rejects empty hosts and non-TCP port values
// (non-numeric, zero, or > 65535).
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

	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", fmt.Errorf("invalid adapter location %q: %w", location, err)
	}
	if host == "" {
		return "", fmt.Errorf("adapter location %q has no host", location)
	}
	portNum, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return "", fmt.Errorf("adapter location %q has invalid port %q: %w", location, port, err)
	}
	if portNum == 0 {
		return "", fmt.Errorf("adapter location %q has invalid port %q", location, port)
	}
	return port, nil
}
