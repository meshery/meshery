package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta3/component"
	"github.com/meshery/schemas/models/v1beta3/connection"
	"github.com/stretchr/testify/require"
)

// registerCategoryTestModel registers the model before its component and
// relationship, as a model import does.
func registerCategoryTestModel(t *testing.T, rm *registry.RegistryManager) {
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

	host := connection.Connection{
		Name:           "meshery",
		Kind:           "meshery",
		ConnectionType: "registry",
		Status:         connection.ConnectionStatusConnected,
	}
	_, _, err = rm.RegisterEntity(host, comp.Model)
	require.NoError(t, err, "registering the meshery-core model")
	_, _, err = rm.RegisterEntity(host, &comp)
	require.NoError(t, err, "registering the WASMFilter component")
	_, _, err = rm.RegisterEntity(host, &relationship.RelationshipDefinition{
		Kind:             "edge",
		RelationshipType: "binding",
		SubType:          "reference",
		ModelId:          &comp.Model.ID,
	})
	require.NoError(t, err, "registering the meshery-core relationship")
}

func TestGetModelsByCategoryEmbedsDefinitionsWhenAsked(t *testing.T) {
	log, err := logger.New("component_handler_category_test", logger.Options{})
	require.NoError(t, err)
	category := "Orchestration & Management"

	cases := []struct {
		name    string
		vars    map[string]string
		handler func(*Handler, http.ResponseWriter, *http.Request)
	}{
		{"category models", map[string]string{"category": category}, (*Handler).GetMeshmodelModelsByCategories},
		{"category model by name", map[string]string{"category": category, "model": "meshery-core"}, (*Handler).GetMeshmodelModelsByCategoriesByModel},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, query := range []struct {
				param         string
				child         any
				components    int
				relationships int
			}{
				{param: "components", child: &component.ComponentDefinition{}, components: 1},
				{param: "relationships", child: &relationship.RelationshipDefinition{}, relationships: 1},
			} {
				t.Run(query.param, func(t *testing.T) {
					rm, db := newTestRegistryDB(t)
					registerCategoryTestModel(t, rm)
					h := &Handler{registryManager: rm, log: log}
					req := httptest.NewRequest(http.MethodGet, "/api/registry/categories/x/models?"+query.param+"=true", nil)
					req = mux.SetURLVars(req, tc.vars)
					rw := httptest.NewRecorder()

					tc.handler(h, rw, req)

					require.Equal(t, http.StatusOK, rw.Code)
					var page struct {
						Models []struct {
							Name          string            `json:"name"`
							Components    []json.RawMessage `json:"components"`
							Relationships []json.RawMessage `json:"relationships"`
						} `json:"models"`
					}
					require.NoError(t, json.Unmarshal(rw.Body.Bytes(), &page))
					require.Len(t, page.Models, 1)
					require.Equal(t, "meshery-core", page.Models[0].Name)
					require.Len(t, page.Models[0].Components, query.components)
					require.Len(t, page.Models[0].Relationships, query.relationships)

					require.NoError(t, db.Migrator().DropTable(query.child))
					rw = httptest.NewRecorder()
					tc.handler(h, rw, req)
					require.Equal(t, http.StatusInternalServerError, rw.Code)
					var failure struct {
						Error string `json:"error"`
					}
					require.NoError(t, json.Unmarshal(rw.Body.Bytes(), &failure))
					require.Equal(t, "could not get meshmodel entities", failure.Error)
				})
			}
		})
	}
}
