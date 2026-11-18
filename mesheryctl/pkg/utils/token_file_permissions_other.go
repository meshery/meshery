//go:build !windows

package utils

import "os"

func writeTokenFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, tokenFileMode)
	if err != nil {
		return err
	}
	if err := f.Chmod(tokenFileMode); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Truncate(0); err != nil {
		_ = f.Close()
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
