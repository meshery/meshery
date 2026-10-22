package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAuthTokenFile(t *testing.T) {
	t.Run("creates auth file with 0600 permissions", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "auth.json")

		if err := WriteAuthTokenFile(path, []byte(`{"token":"test-token"}`)); err != nil {
			t.Fatalf("WriteAuthTokenFile() error = %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("os.Stat() error = %v", err)
		}

		if got := info.Mode().Perm(); got != 0o600 {
			t.Errorf("file permissions = %o, want %o", got, 0o600)
		}
	})

	t.Run("corrects existing insecure permissions to 0600", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "auth.json")

		if err := os.WriteFile(path, []byte(`{"token":"test-token"}`), 0o644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}

		if err := os.Chmod(path, 0o644); err != nil {
			t.Fatalf("os.Chmod() error = %v", err)
		}

		if err := WriteAuthTokenFile(path, []byte(`{"token":"updated-token"}`)); err != nil {
			t.Fatalf("WriteAuthTokenFile() error = %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("os.Stat() error = %v", err)
		}

		if got := info.Mode().Perm(); got != 0o600 {
			t.Errorf("file permissions = %o, want %o", got, 0o600)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("os.ReadFile() error = %v", err)
		}

		if got := string(content); got != `{"token":"updated-token"}` {
			t.Errorf("file content = %q, want %q", got, `{"token":"updated-token"}`)
		}
	})

	t.Run("leaves no temp file behind on success", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "auth.json")

		if err := WriteAuthTokenFile(path, []byte(`{"token":"x"}`)); err != nil {
			t.Fatalf("WriteAuthTokenFile() error = %v", err)
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("os.ReadDir() error = %v", err)
		}

		if len(entries) != 1 {
			t.Errorf("directory contains %d entries, want 1", len(entries))
		}

		if entries[0].Name() != "auth.json" {
			t.Errorf("file name = %q, want %q", entries[0].Name(), "auth.json")
		}
	})
}
