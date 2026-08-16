package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteAuthTokenFile(t *testing.T) {
	t.Run("creates auth file with 0600 permissions", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "auth.json")

		err := WriteAuthTokenFile(path, []byte(`{"token":"test-token"}`))
		require.NoError(t, err)

		info, err := os.Stat(path)
		require.NoError(t, err)

		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	})

	t.Run("corrects existing insecure permissions to 0600", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "auth.json")

		err := os.WriteFile(path, []byte(`{"token":"test-token"}`), 0o644)
		require.NoError(t, err)

		err = os.Chmod(path, 0o644)
		require.NoError(t, err)

		err = WriteAuthTokenFile(path, []byte(`{"token":"updated-token"}`))
		assert.NoError(t, err)

		info, err := os.Stat(path)
		require.NoError(t, err)

		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())

		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Contains(t, string(content), "updated-token")
	})
	t.Run("leaves no temp file behind on success", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "auth.json")

		require.NoError(t, WriteAuthTokenFile(path, []byte(`{"token":"x"}`)))

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)

		assert.Len(t, entries, 1)
		assert.Equal(t, "auth.json", entries[0].Name())
	})
}
