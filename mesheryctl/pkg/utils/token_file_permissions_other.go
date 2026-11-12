//go:build !windows

package utils

import "os"

func openTokenFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE, tokenFileMode)
}

func restrictTokenFileToCurrentUser(_ *os.File) error {
	return nil
}
