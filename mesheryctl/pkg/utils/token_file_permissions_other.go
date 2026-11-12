//go:build !windows

package utils

import "os"

func restrictTokenFileToCurrentUser(_ *os.File) error {
	return nil
}
