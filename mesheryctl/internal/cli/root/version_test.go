package root

import (
	"bytes"
	"errors"
	"fmt"
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

func TestVersionCmdArgs(t *testing.T) {
	tests := []struct {
		Name          string
		Args          []string
		ExpectError   bool
		ExpectedError error
	}{
		{
			Name:        "valid with no positional arguments",
			Args:        []string{},
			ExpectError: false,
		},
		{
			Name:          "invalid with unexpected positional argument",
			Args:          []string{"unexpected-arg"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("this command does not accept positional arguments\n\n%s", versionUsageMsg)),
		},
		{
			Name:          "invalid with multiple unexpected positional arguments",
			Args:          []string{"arg1", "arg2"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("this command does not accept positional arguments\n\n%s", versionUsageMsg)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := versionCmd.Args(versionCmd, tt.Args)
			if tt.ExpectError {
				assert.Error(t, err)
				utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

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
		{
			name:       "given server returns 500 Internal Server Error when version then show server version as unavailable and log status error",
			responder:  httpmock.NewStringResponder(http.StatusInternalServerError, `{"error":"internal error"}`),
			wantStdout: []string{"Client", "Server", "unavailable"},
			wantLog:    "Server returned with status code: 500",
		},
		{
			name:       "given server returns 502 Bad Gateway when version then show server version as unavailable and log status error",
			responder:  httpmock.NewStringResponder(http.StatusBadGateway, `<html>502 Bad Gateway</html>`),
			wantStdout: []string{"Client", "Server", "unavailable"},
			wantLog:    "Server returned with status code: 502",
		},
		{
			name:       "given server returns 401 Unauthorized when version then show server version as unavailable and log status error",
			responder:  httpmock.NewStringResponder(http.StatusUnauthorized, `{"message":"Unauthorized"}`),
			wantStdout: []string{"Client", "Server", "unavailable"},
			wantLog:    "Server returned with status code: 401",
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

func TestVersionPreRunERecoveryDoesNotPanicOnNilConfig(t *testing.T) {
	origConfigPath := utils.DefaultConfigPath
	origMctlCfg := mctlCfg
	origStdin := os.Stdin
	origLog := utils.Log
	t.Cleanup(func() {
		utils.DefaultConfigPath = origConfigPath
		mctlCfg = origMctlCfg
		os.Stdin = origStdin
		utils.Log = origLog
		viper.Reset()
	})

	configPath := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte("contexts: invalid-not-a-map\n"), 0o600))
	utils.DefaultConfigPath = configPath

	logs := utils.SetupMeshkitLoggerTesting(t, false)
	utils.Log.UpdateErrorLogOutput(logs)

	r, w, err := os.Pipe()
	require.NoError(t, err)
	_, err = w.WriteString("y\n")
	require.NoError(t, err)
	require.NoError(t, w.Close())
	t.Cleanup(func() { _ = r.Close() })
	os.Stdin = r

	// Force config.GetMesheryCtl(viper.GetViper()) to fail with an invalid contexts value
	// so PreRunE enters the recovery branch and hits the post-recreation GetMesheryCtl check.
	viper.Reset()
	viper.Set("contexts", "invalid-not-a-map")

	require.NotPanics(t, func() {
		err := versionCmd.PreRunE(versionCmd, []string{})
		require.Error(t, err)
		utils.AssertMeshkitErrorsEqual(t, err, ErrUnmarshallingConfigFile)
	})
}
