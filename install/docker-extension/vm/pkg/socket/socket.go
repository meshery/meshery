package socket

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ListenFD listens to a file descriptor.
func ListenFD(filedes string) (net.Listener, error) {
	fd, err := strconv.Atoi(filedes)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse file descriptor: %s", filedes)
	}
	file := os.NewFile(uintptr(fd), fmt.Sprintf("fd:%v", filedes))
	res, err := net.FileListener(file)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot convert fd %v to net.Listener", fd)
	}
	return res, nil
}

// Listen wraps net.Listen, preferring ListenUnix when applicable, and
// offers support for the "fd" network type, in which case address is
// a file descriptor.
func Listen(network, address string) (net.Listener, error) {
	switch network {
	case "fd":
		return ListenFD(address)
	case "unix":
		return ListenUnix(address)
	default:
		return net.Listen(network, address)
	}
}

// ListenOn splits a "network:address" string to invoke Listen.
func ListenOn(addr string) (net.Listener, error) {
	ss := strings.SplitN(addr, ":", 2)
	if len(ss) < 2 {
		return nil, errors.Errorf("invalid listener address: %v", addr)
	}
	return Listen(ss[0], ss[1])
}

// Client Creates a client connected to the specified socket
func Client(addr string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return DialSocket(addr)
			},
		},
	}
}
