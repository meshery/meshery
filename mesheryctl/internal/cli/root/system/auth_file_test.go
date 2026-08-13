package system

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

		err := writeAuthTokenFile(path, []byte(`{"token":"test-token"}`))
		require.NoError(t, err)

		info, err := os.Stat(path)
		require.NoError(t, err)

		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	})

	t.Run("corrects existing insecure permissions to 0600", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "auth.json")

		err := os.WriteFile(path, []byte(`{"token":"test-token"}`), 0o644)
		require.NoError(t, err)

		err = writeAuthTokenFile(path, []byte(`{"token":"updated-token"}`))
		require.NoError(t, err)

		info, err := os.Stat(path)
		require.NoError(t, err)

		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	})
}
