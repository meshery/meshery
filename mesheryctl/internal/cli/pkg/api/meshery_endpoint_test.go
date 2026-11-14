package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// TestDeleteWithTrailingSlashEndpointReachesDeleteHandler guards issue #23389 end to end. With an
// endpoint ending in "/", requests went to "//api/...". Meshery Server's gorilla/mux router answers that
// with a 301 to "/api/...", and Go's client re-sends a redirected DELETE as a GET, so a GET handler on
// the same path ran instead and the delete was reported as successful. The fake server here uses
// gorilla/mux with default settings, as server/router/server.go does.
func TestDeleteWithTrailingSlashEndpointReachesDeleteHandler(t *testing.T) {
	var (
		mu  sync.Mutex
		ran []string
	)
	record := func(handler string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			ran = append(ran, handler+" "+r.URL.Path)
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))
		}
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/environments/{id}", record("GET")).Methods(http.MethodGet)
	router.HandleFunc("/api/environments/{id}", record("DELETE")).Methods(http.MethodDelete)
	server := httptest.NewServer(router)
	defer server.Close()

	useTrailingSlashEndpoint(t, server.URL)

	resp, err := Delete("api/environments/22222222-2222-2222-2222-222222222222")
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	mu.Lock()
	defer mu.Unlock()
	want := "DELETE /api/environments/22222222-2222-2222-2222-222222222222"
	if len(ran) != 1 || ran[0] != want {
		t.Fatalf("handlers that ran = %v, want only %q", ran, want)
	}
}

// TestConnectionErrorShowsTrimmedEndpoint checks that the "Unable to connect" message names the
// endpoint the request was actually sent to, without the trailing slash it was stored with.
func TestConnectionErrorShowsTrimmedEndpoint(t *testing.T) {
	// A server that is closed straight away gives an address nothing listens on.
	server := httptest.NewServer(http.NotFoundHandler())
	unreachable := server.URL
	server.Close()

	useTrailingSlashEndpoint(t, unreachable)

	_, err := Delete("api/environments/22222222-2222-2222-2222-222222222222")
	if err == nil {
		t.Fatal("expected a connection error, got nil")
	}
	want := fmt.Sprintf("Unable to connect to Meshery server at %s (current context).", unreachable)
	if !strings.Contains(err.Error(), want) {
		t.Errorf("error = %q, want it to contain %q", err.Error(), want)
	}
}

// useTrailingSlashEndpoint points mesheryctl at baseURL through a private config whose endpoint ends in
// "/", the form `system context create --url http://host:9081/` stores, and attaches the test token.
func useTrailingSlashEndpoint(t *testing.T, baseURL string) {
	t.Helper()

	// TestMakeRequest_Failures activates httpmock through utils.InitTestEnvironment and never
	// deactivates it, so when it runs first (as with -shuffle) httpmock would intercept these requests.
	httpmock.Deactivate()

	configPath := filepath.Join(t.TempDir(), "config.yaml")
	config := fmt.Sprintf(`contexts:
  local:
    endpoint: %s/
    token: default
    platform: docker
current-context: local
tokens:
- name: default
  location: auth.json
`, baseURL)
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatalf("writing config: %v", err)
	}
	utils.SetupCustomContextEnv(t, configPath)
	utils.TokenFlag = utils.GetToken(t)
	t.Cleanup(func() { utils.TokenFlag = "Not Set" })
	_ = utils.SetupMeshkitLoggerTesting(t, false)
}
