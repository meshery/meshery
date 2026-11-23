package root

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// The local context in pkg/utils/TestConfig.yaml points at this endpoint.
	versionTestServerURL = "http://localhost:9081/api/system/version"
	// version always checks for the latest mesheryctl release; keep it offline.
	versionTestLatestReleaseURL = "https://docs.meshery.io/project/releases/latest"
)

// setupVersionTest loads the shared test meshconfig into mctlCfg, activates
// HTTP mocking, and returns the buffer that log and error output is written to.
func setupVersionTest(t *testing.T) *bytes.Buffer {
	t.Helper()

	// utils.SetupContextEnv resolves the shared fixture four levels up, which
	// suits command groups such as root/design; this package is one level
	// higher, so point at the fixture directly. Copy it first so the command
	// under test never writes to the shared file.
	fixture := filepath.Join("..", "..", "..", "pkg", "utils", "TestConfig.yaml")
	utils.SetupCustomContextEnv(t, utils.CopyMeshconfigFixture(t, fixture))
	var err error
	mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
	require.NoError(t, err)

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	httpmock.RegisterResponder(http.MethodGet, versionTestLatestReleaseURL,
		httpmock.NewStringResponder(http.StatusOK, "v0.0.0"))

	// Warnings go to the logger's regular output and errors to its error
	// output (stderr by default), so collect both in one buffer.
	logs := utils.SetupMeshkitLoggerTesting(t, false)
	utils.Log.UpdateErrorLogOutput(logs)
	t.Cleanup(func() { utils.Log.UpdateErrorLogOutput(os.Stderr) })
	return logs
}

// captureVersionStdout returns what fn writes to stdout, where the version
// table is printed.
func captureVersionStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	out := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		out <- buf.String()
	}()
	// Close the pipe writer even if fn panics, so the reader goroutine
	// finishes instead of hanging the test on <-out.
	defer func() { _ = w.Close() }()

	fn()
	_ = w.Close()
	got := <-out
	_ = r.Close()
	return got
}

func TestVersionCmd(t *testing.T) {
	tests := []struct {
		name        string
		responder   httpmock.Responder
		wantStdout  []string
		notInStdout []string
		wantLog     string
	}{
		{
			name: "given server returns its version when version then show server build and commit sha",
			responder: httpmock.NewStringResponder(http.StatusOK,
				`{"build":"v0.8.200","commitsha":"abc1234","release_channel":"stable"}`),
			wantStdout:  []string{"Client", "Server", "v0.8.200", "abc1234"},
			notInStdout: []string{"unavailable"},
		},
		{
			name:       "given server is unreachable when version then show server version as unavailable",
			responder:  httpmock.NewErrorResponder(errors.New("connection refused")),
			wantStdout: []string{"Client", "Server", "unavailable"},
			wantLog:    "Unable to communicate with Meshery server",
		},
		{
			name:       "given server returns invalid JSON when version then show server version as unavailable",
			responder:  httpmock.NewStringResponder(http.StatusOK, `not json`),
			wantStdout: []string{"Client", "Server", "unavailable"},
			wantLog:    "Error processing json API data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs := setupVersionTest(t)
			httpmock.RegisterResponder(http.MethodGet, versionTestServerURL, tt.responder)

			stdout := captureVersionStdout(t, func() {
				versionCmd.Run(versionCmd, []string{})
			})

			for _, want := range tt.wantStdout {
				assert.Contains(t, stdout, want)
			}
			for _, notWant := range tt.notInStdout {
				assert.NotContains(t, stdout, notWant)
			}
			if tt.wantLog != "" {
				assert.Contains(t, logs.String(), tt.wantLog)
			}
			assert.Equal(t, 1, httpmock.GetCallCountInfo()["GET "+versionTestServerURL],
				"version should query the server exactly once")
		})
	}
}
