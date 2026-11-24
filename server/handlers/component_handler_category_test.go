package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"os"
	"path/filepath"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta3/component"
	"github.com/stretchr/testify/require"
)

// registerMesheryCoreModelWithComponent registers the meshery-core model and
// its WASMFilter component the way a model import does: model first, so the
// model has its own registry entry and is listed.
func registerMesheryCoreModelWithComponent(t *testing.T, rm *registry.RegistryManager) {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	path := filepath.Join(filepath.Dir(filename), "..", "..", "models", "meshery-core", "0.7.2", "v1.0.0", "components", "WASMFilter.json")
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	var comp component.ComponentDefinition
	require.NoError(t, json.Unmarshal(content, &comp))
	comp.ID, err = comp.GenerateID()
	require.NoError(t, err)
	comp.Model.ID, err = comp.Model.GenerateID()
	require.NoError(t, err)

	host := registry.RegistrantHostToV1beta3(connection.Connection{
		Name:   "meshery",
		Kind:   "meshery",
		Type:   "registry",
		Status: connection.ConnectionStatusConnected,
	})
	_, _, err = rm.RegisterEntity(host, comp.Model)
	require.NoError(t, err, "registering the meshery-core model")
	_, _, err = rm.RegisterEntity(host, &comp)
	require.NoError(t, err, "registering the WASMFilter component")
}

// The category model routes document `components` and `relationships`
// ("When `true`, embeds the model's ... definitions in the response"), as the
// other model routes do. They must embed the components when asked.
func TestGetModelsByCategoryEmbedsComponentsWhenAsked(t *testing.T) {
	rm, _ := newTestRegistryDB(t)
	registerMesheryCoreModelWithComponent(t, rm)
	log, err := logger.New("component_handler_category_test", logger.Options{})
	require.NoError(t, err)
	h := &Handler{registryManager: rm, log: log}
	category := "Orchestration & Management"

	cases := []struct {
		name    string
		vars    map[string]string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{"category models", map[string]string{"category": category}, h.GetMeshmodelModelsByCategories},
		{"category model by name", map[string]string{"category": category, "model": "meshery-core"}, h.GetMeshmodelModelsByCategoriesByModel},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/registry/categories/x/models?components=true", nil)
			req = mux.SetURLVars(req, tc.vars)
			rw := httptest.NewRecorder()

			tc.handler(rw, req)

			require.Equal(t, http.StatusOK, rw.Code)
			var page struct {
				Models []struct {
					Name       string            `json:"name"`
					Components []json.RawMessage `json:"components"`
				} `json:"models"`
			}
			require.NoError(t, json.Unmarshal(rw.Body.Bytes(), &page))
			require.Len(t, page.Models, 1)
			require.Equal(t, "meshery-core", page.Models[0].Name)
			require.Len(t, page.Models[0].Components, 1, "?components=true must embed the model's WASMFilter component")
		})
	}
}
