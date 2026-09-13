package system

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkiterrors "github.com/meshery/meshkit/errors"
)

var update = flag.Bool("update", false, "update golden files")

// TestIsLocalMesheryEndpoint guards against #21696: runDockerHealthChecks used to
// split the endpoint on ":" and index into the result without checking its length,
// panicking with "index out of range [1] with length 1" whenever the configured
// endpoint had no port (e.g. "localhost") or was empty.
func TestIsLocalMesheryEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		want     bool
	}{
		{"scheme, host and port", "http://localhost:9081", true},
		{"scheme and host, no port", "http://localhost", true},
		{"bare host, no scheme or port", "localhost", true},
		{"host and port, no scheme", "localhost:9081", true},
		{"empty endpoint", "", false},
		{"loopback IP with port", "127.0.0.1:9081", true},
		{"remote host", "http://example.com:9081", false},
		{"non-http scheme", "tcp://localhost:2375", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLocalMesheryEndpoint(tt.endpoint); got != tt.want {
				t.Errorf("isLocalMesheryEndpoint(%q) = %v, want %v", tt.endpoint, got, tt.want)
			}
		})
	}
}

// TestRunDockerHealthChecksInvalidEndpoint pins the endpoint validation that
// guards runDockerHealthChecks. NewHealthChecker only validates the context
// name, so without this check a malformed endpoint reached the Docker error
// handling below instead of reporting ErrInvalidEndpoint.
func TestRunDockerHealthChecksInvalidEndpoint(t *testing.T) {
	invalidEndpoints := []struct {
		name     string
		endpoint string
	}{
		{"empty endpoint", ""},
		{"whitespace only", "   "},
		{"bare host, no port", "localhost"},
		{"scheme and host, no port", "http://localhost"},
		{"non-numeric port", "http://localhost:abcd"},
		{"port above the valid range", "http://localhost:65536"},
	}

	for _, tt := range invalidEndpoints {
		t.Run(tt.name, func(t *testing.T) {
			hc := &HealthChecker{
				Options: &HealthCheckOptions{},
				context: &config.Context{Endpoint: tt.endpoint},
			}

			err := hc.runDockerHealthChecks()
			if err == nil {
				t.Fatalf("runDockerHealthChecks() with endpoint %q: expected an error, got nil", tt.endpoint)
			}
			if code := meshkiterrors.GetCode(err); code != ErrInvalidEndpointCode {
				t.Errorf("runDockerHealthChecks() with endpoint %q: error code = %s, want %s", tt.endpoint, code, ErrInvalidEndpointCode)
			}
		})
	}
}

// This is an Integration test
func TestPreflightCmdIntegration(t *testing.T) {
	// skipping this integration test with --short flag
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	setupContextTestEnv(t)
	tests := []utils.CmdTestInput{
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--pre"},
			ExpectedResponse: "check.output.golden",
		},
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--preflight"},
			ExpectedResponse: "check.output.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// setting up log to grab logs
			buf := setupSystemOutCmdTest(t)
			defer func() {
				buf.Reset()
				// resetCmdFlags(SystemCmd, t)
			}()

			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			output := buf.String()
			actualResponse := output

			// get current directory
			_, filename, _, ok := runtime.Caller(0)
			if !ok {
				t.Fatal("Not able to get current working directory")
			}

			currDir := filepath.Dir(filename)
			testdataDir := filepath.Join(currDir, "testdata")
			// testdataDir := filepath.Join(filepath.Dir(filename), tf.dir, tf.name)
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, filepath.Join(testdataDir, "check"))

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			if actualResponse != expectedResponse {
				t.Errorf("expected response %q, got %q", expectedResponse, actualResponse)
			}
		})
		t.Log("PreflightCmdIntegration Test Passed")
	}
}
