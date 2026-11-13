package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

	// Other tests in this package activate httpmock, which would intercept the request.
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
`, server.URL)
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatalf("writing config: %v", err)
	}
	utils.SetupCustomContextEnv(t, configPath)
	utils.TokenFlag = utils.GetToken(t)
	defer func() { utils.TokenFlag = "Not Set" }()
	_ = utils.SetupMeshkitLoggerTesting(t, false)

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
