package model

import (
	"encoding/json"
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

type fakeRegistryModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Version     string `json:"version"`
}

// newFakeModelRegistry serves the model routes that model delete uses, matching names the way
// Meshery Server does: GET /api/registry/models?search= is a substring match (name LIKE '%x%'),
// GET /api/registry/models/{model} is an exact name match. It returns the server URL and a
// function reporting the DELETE requests the server received.
func newFakeModelRegistry(t *testing.T, registered []fakeRegistryModel) (string, func() []string) {
	t.Helper()
	var (
		mu      sync.Mutex
		deletes []string
	)
	page := func(w http.ResponseWriter, matches []fakeRegistryModel) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"page": 0, "pageSize": 10, "page_size": 10,
			"totalCount": len(matches), "total_count": len(matches),
			"models": matches,
		})
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/registry/models", func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		var matches []fakeRegistryModel
		for _, m := range registered {
			if strings.Contains(m.Name, search) || strings.Contains(m.DisplayName, search) {
				matches = append(matches, m)
			}
		}
		page(w, matches)
	}).Methods(http.MethodGet)
	router.HandleFunc("/api/registry/models/{model}", func(w http.ResponseWriter, r *http.Request) {
		var matches []fakeRegistryModel
		for _, m := range registered {
			if m.Name == mux.Vars(r)["model"] {
				matches = append(matches, m)
			}
		}
		page(w, matches)
	}).Methods(http.MethodGet)
	router.HandleFunc("/api/registry/models/{id}", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		deletes = append(deletes, mux.Vars(r)["id"])
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodDelete)
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	return server.URL, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return append([]string(nil), deletes...)
	}
}

// useModelServer points mesheryctl at baseURL through a private meshconfig.
func useModelServer(t *testing.T, baseURL string) {
	t.Helper()
	// Tests run through utils.InitTestEnvironment leave httpmock active, which would intercept
	// the requests sent to the fake server.
	httpmock.DeactivateAndReset()

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
`, baseURL)
	if err := os.WriteFile(configPath, []byte(meshconfig), 0o600); err != nil {
		t.Fatalf("writing meshconfig: %v", err)
	}
	utils.SetupCustomContextEnv(t, configPath)
	utils.TokenFlag = utils.GetToken(t)
	t.Cleanup(func() { utils.TokenFlag = "Not Set" })
}

func TestDeleteModelByNameDeletesOnlyThatModel(t *testing.T) {
	const (
		kubernetesID      = "11111111-1111-4111-8111-111111111111"
		secretGeneratorID = "22222222-2222-4222-8222-222222222222"
	)
	secretGenerator := fakeRegistryModel{ID: secretGeneratorID, Name: "kubernetes-secret-generator", DisplayName: "Kubernetes Secret Generator", Version: "v1.0.0"}
	kubernetes := fakeRegistryModel{ID: kubernetesID, Name: "kubernetes", DisplayName: "Kubernetes", Version: "v1.30.0"}

	tests := []struct {
		name        string
		registered  []fakeRegistryModel
		wantDeletes []string
		wantOutput  string
	}{
		{
			name:        "given only a model whose name contains the requested name when running model delete then nothing is deleted",
			registered:  []fakeRegistryModel{secretGenerator},
			wantDeletes: nil,
			wantOutput:  "No model(s) found with the name: kubernetes",
		},
		{
			name:        "given the requested model and another model containing its name when running model delete then only the requested model is deleted",
			registered:  []fakeRegistryModel{secretGenerator, kubernetes},
			wantDeletes: []string{kubernetesID},
			wantOutput:  fmt.Sprintf("Model 'Kubernetes' (ID: %s) has been deleted", kubernetesID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseURL, deletes := newFakeModelRegistry(t, tt.registered)
			useModelServer(t, baseURL)
			out := utils.SetupMeshkitLoggerTesting(t, false)
			defer utils.ResetCommandFlags(ModelCmd, t)

			ModelCmd.SetArgs([]string{"delete", "kubernetes"})
			if err := ModelCmd.Execute(); err != nil {
				t.Fatalf("model delete kubernetes returned error: %v", err)
			}

			if got := deletes(); fmt.Sprint(got) != fmt.Sprint(tt.wantDeletes) {
				t.Errorf("DELETE requests for model IDs %v, want %v", got, tt.wantDeletes)
			}
			if !strings.Contains(out.String(), tt.wantOutput) {
				t.Errorf("output = %q, want it to contain %q", out.String(), tt.wantOutput)
			}
		})
	}
}
