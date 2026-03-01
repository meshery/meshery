package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckMesheryctlClientVersion(t *testing.T) {
	var buf bytes.Buffer

	// 1. Setup Mock Server
	expectedVersion := "v9.9.9"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(expectedVersion))
	}))
	defer server.Close()

	// 2. Override URL and Logger
	oldURL := LatestVersionURL
	LatestVersionURL = server.URL
	defer func() { LatestVersionURL = oldURL }()

	Log = SetupMeshkitLogger("test", false, &buf)

	t.Run("Outdated version notification", func(t *testing.T) {
		buf.Reset()
		CheckMesheryctlClientVersion("v0.0.1")
		output := buf.String()

		assert.NotEmpty(t, output)
		assert.Contains(t, output, "A new release of mesheryctl is available")
	})

	t.Run("Latest version notification", func(t *testing.T) {
		buf.Reset()
		CheckMesheryctlClientVersion(expectedVersion)
		output := buf.String()

		assert.Contains(t, output, "is the latest release")
	})
}