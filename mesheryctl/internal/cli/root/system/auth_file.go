package system

import "os"

func writeAuthTokenFile(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}

	return os.Chmod(path, 0o600)
}
