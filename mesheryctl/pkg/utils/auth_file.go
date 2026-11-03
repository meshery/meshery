package utils

import (
	"os"
	"path/filepath"
)

func WriteAuthTokenFile(path string, data []byte) error {
	dir := filepath.Dir(path)

	tmp, err := os.CreateTemp(dir, ".auth-*.tmp")
	if err != nil {
		return ErrWriteAuthTokenFile(err)
	}

	tmpPath := tmp.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return ErrWriteAuthTokenFile(err)
	}

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return ErrWriteAuthTokenFile(err)
	}

	if err := tmp.Close(); err != nil {
		return ErrWriteAuthTokenFile(err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return ErrWriteAuthTokenFile(err)
	}

	return nil
}
