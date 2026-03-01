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
	// Define expectedVersion at the top level of the test so sub-tests can access it.
	expectedVersion := "v9.9.9"

	// 1. Save and restore the global Log state to avoid side effects on other tests.
	oldLog := Log
	defer func() { Log = oldLog }()

	// 2. Setup a Mock Server to simulate the Meshery release API.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(expectedVersion))
	}))
	defer server.Close()

	// 3. Initialize the logger for this test to capture output in our buffer.
	Log = SetupMeshkitLogger("test", false, &buf)

	t.Run("Outdated version notification", func(t *testing.T) {
		buf.Reset()
		// We pass an old version to trigger the "new release available" log path.
		CheckMesheryctlClientVersion("v0.0.1")
		output := buf.String()

		assert.NotEmpty(t, output, "Logger failed to write to buffer")
		assert.Contains(t, output, "Checking for latest version")
	})

	t.Run("Verify internal version fetcher", func(t *testing.T) {
		// Test the internal function directly by injecting the mock server URL.
		version, err := getLatestVersion(server.URL)
		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, version)
	})
}