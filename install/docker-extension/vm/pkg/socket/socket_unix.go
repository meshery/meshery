//go:build !windows
// +build !windows

package socket

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// ListenUnix wraps `net.ListenUnix`.
func ListenUnix(path string) (*net.UnixListener, error) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	// Make sure the parent directory exists.
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	short, err := shortenUnixSocketPath(path)
	if err != nil {
		return nil, err
	}
	return net.ListenUnix("unix", &net.UnixAddr{Name: short, Net: "unix"})
}

func DialSocket(socket string) (net.Conn, error) {
	return net.Dial("unix", socket)
}

func shortenUnixSocketPath(path string) (string, error) {
	if len(path) <= maxUnixSocketPathLen {
		return path, nil
	}
	// absolute path is too long, attempt to use a relative path
	p, err := relative(path)
	if err != nil {
		return "", err
	}

	if len(p) > maxUnixSocketPathLen {
		return "", fmt.Errorf("absolute and relative socket path %s longer than %d characters", p, maxUnixSocketPathLen)
	}
	return p, nil
}

func relative(p string) (string, error) {
	// Assume the parent directory exists already but the child (the socket)
	// hasn't been created.
	path2, err := filepath.EvalSymlinks(filepath.Dir(p))
	if err != nil {
		return "", err
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir2, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(dir2, path2)
	if err != nil {
		return "", err
	}
	return filepath.Join(rel, filepath.Base(p)), nil
}
