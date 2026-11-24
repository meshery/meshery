package system

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// A server, proxy or provider can answer 200 with the JSON value null (what
// json.Marshal writes for a nil pointer). Decoding that into a *config.Version
// leaves it nil, and the next call, GetBuild, dereferenced it: mesheryctl system
// check panicked instead of reporting that the version could not be read.
func TestMesheryVersionHealthCheck_NullVersionResponseIsAnError(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("null"))
	}))
	t.Cleanup(server.Close)

	hc := &HealthChecker{
		Options: &HealthCheckOptions{},
		mctlCfg: &config.MesheryCtlConfig{
			Contexts:       map[string]config.Context{"test": {Endpoint: server.URL}},
			CurrentContext: "test",
		},
	}

	err := hc.runMesheryVersionHealthChecks()
	if err == nil {
		t.Fatal("a null version response must be reported as an error")
	}
	if !strings.Contains(err.Error(), "no version information") {
		t.Errorf("error = %q, want it to say the server returned no version information", err)
	}
}
