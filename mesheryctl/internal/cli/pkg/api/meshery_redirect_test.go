package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
)

// A redirect that would re-send a DELETE as a GET must reach the caller as
// its own error, not be followed into the GET handler and not be reported as
// "Unable to connect to Meshery server".
func TestDeleteReportsMethodChangingRedirect(t *testing.T) {
	// TestMakeRequest_Failures activates httpmock through utils.InitTestEnvironment
	// and never deactivates it, so when it runs first httpmock would intercept
	// this test's requests.
	httpmock.Deactivate()

	var getHits, deleteHits atomic.Int32
	r := mux.NewRouter()
	r.HandleFunc("/api/items/{id}", func(w http.ResponseWriter, _ *http.Request) {
		getHits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"1"}`)
	}).Methods(http.MethodGet)
	r.HandleFunc("/api/items/{id}", func(w http.ResponseWriter, _ *http.Request) {
		deleteHits.Add(1)
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodDelete)
	server := httptest.NewServer(r)
	t.Cleanup(server.Close)

	configPath := filepath.Join(t.TempDir(), "config.yaml")
	meshconfig := fmt.Sprintf(`contexts:
  local:
    endpoint: %s
    token: default
    platform: docker
current-context: local
tokens:
- name: default
  location: auth.json
`, server.URL)
	if err := os.WriteFile(configPath, []byte(meshconfig), 0o600); err != nil {
		t.Fatal(err)
	}
	utils.SetupCustomContextEnv(t, configPath)
	utils.TokenFlag = utils.GetToken(t)
	t.Cleanup(func() { utils.TokenFlag = "Not Set" })
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	// The leading slash makes the request URL "<endpoint>//api/items/1", which
	// gorilla/mux answers with 301 to "/api/items/1".
	resp, err := Delete("/api/items/1")
	if resp != nil {
		_ = resp.Body.Close()
		t.Errorf("got response %d, want nil", resp.StatusCode)
	}
	if got := errors.GetCode(err); got != utils.ErrRedirectChangesMethodCode {
		t.Errorf("error code = %q (err: %v), want %q", got, err, utils.ErrRedirectChangesMethodCode)
	}
	if got := getHits.Load(); got != 0 {
		t.Errorf("the DELETE was re-sent as a GET and ran the GET handler %d time(s)", got)
	}
	if got := deleteHits.Load(); got != 0 {
		t.Errorf("DELETE handler ran %d time(s), want 0", got)
	}
}
