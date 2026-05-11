package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLatestVersionForMesheryctl(t *testing.T) {
	originalURL := latestVersionURL
	originalClient := latestVersionHTTPClient
	t.Cleanup(func() {
		latestVersionURL = originalURL
		latestVersionHTTPClient = originalClient
	})

	t.Run("returns latest version from server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, "v1.2.3\n")
		}))
		defer server.Close()

		latestVersionURL = server.URL
		latestVersionHTTPClient = server.Client()

		got, err := GetLatestVersionForMesheryctl()
		require.NoError(t, err)
		assert.Equal(t, "v1.2.3", got)
	})

	t.Run("returns empty version when the request times out", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			_, _ = fmt.Fprint(w, "late response")
		}))
		defer server.Close()

		latestVersionURL = server.URL
		latestVersionHTTPClient = &http.Client{Timeout: 1 * time.Millisecond}

		got, err := GetLatestVersionForMesheryctl()
		require.NoError(t, err)
		assert.Empty(t, got)
	})
}
