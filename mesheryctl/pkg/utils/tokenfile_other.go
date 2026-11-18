//go:build !windows

package utils

import "os"

// restrictTokenFile limits an open auth token file to its owner. On Unix the
// permission bits are the access control, so chmod-ing the descriptor to
// tokenFileMode (0600) is enough and also tightens a file an older mesheryctl
// left readable by others.
func restrictTokenFile(f *os.File) error {
	return f.Chmod(tokenFileMode)
}
