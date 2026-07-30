package handlers

// These tests cover the MeshModel component read handlers in
// component_handler.go at the HTTP boundary: the response envelope, the two
// accepted page-size spellings, pageSize=all, model-scoped lookup, and the
// empty-result shape. Route registration is already covered by
// server/router/registry_routes_test.go; the handler behaviour those routes
// depend on was not, so a filter or pagination regression would ship silently.
//
// The harness matches the registry-backed tests in
// policy_relationship_handler_test.go and connections_seeded_api_test.go: an
// in-memory registry seeded through the registry's own RegisterEntity API (no
// raw SQL), driven through httptest, asserting on the bytes a client actually
// receives rather than on the database.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	v1beta3comp "github.com/meshery/schemas/models/v1beta3/component"
	"github.com/stretchr/testify/require"
)

// seedComponentInModel registers one component belonging to modelName into rm
// using the registry's own RegisterEntity API, mirroring seedTestComponent in
// policy_relationship_handler_test.go but parameterised so a test can seed
// components across more than one model.
func seedComponentInModel(t *testing.T, rm *registry.RegistryManager, modelName, categoryName, kind, apiVersion, displayName string) {
	t.Helper()
	conn := connection.Connection{
		Name:    "test-registrant",
		Kind:    "kubernetes",
		Type:    "platform",
		SubType: "orchestration",
		Status:  connection.ConnectionStatusConnected,
	}
	enabled := v1beta3comp.Enabled
	bgColor := "#123456"
	shape := core.Shape("round-rectangle")
	comp := v1beta3comp.ComponentDefinition{
		DisplayName:   core.InputString(displayName),
		SchemaVersion: "core.meshery.io/v1beta1",
		Status:        &enabled,
		Component: v1beta3comp.Component{
			Kind:    kind,
			Version: apiVersion,
			Schema:  `{"properties":{}}`, // non-empty schema is required by Create
		},
		Model: &model.ModelDefinition{
			Name:          modelName,
			DisplayName:   modelName,
			SchemaVersion: "models.meshery.io/v1beta1",
			Version:       "v1.0.0",
			Model:         model.Model{Version: "v1.0.0"},
			Category:      category.CategoryDefinition{Name: category.CategoryDefinitionName(categoryName)},
			Status:        model.Enabled,
		},
		Styles: &core.ComponentStyles{
			BackgroundColor: &bgColor,
			PrimaryColor:    "#123456",
			Shape:           &shape,
		},
	}
	id, err := comp.GenerateID()
	require.NoError(t, err)
	comp.ID = id
	_, _, err = rm.RegisterEntity(registry.RegistrantHostToV1beta3(conn), &comp)
	require.NoError(t, err, "seedComponentInModel(%s/%s): RegisterEntity failed", modelName, kind)
}

// apiComponentsPage decodes the fields of the components response envelope that
// these tests assert on, using the canonical camelCase wire keys.
type apiComponentsPage struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
	Components []struct {
		DisplayName string `json:"displayName"`
		Model       struct {
			Name string `json:"name"`
		} `json:"model"`
	} `json:"components"`
}

// callComponents issues a request straight to handlerFn (the same func the
// route is bound to), sets any mux path vars, and returns the status and the
// decoded response together with the raw body for diagnostics.
func callComponents(t *testing.T, rawURL string, vars map[string]string, handlerFn http.HandlerFunc) (int, apiComponentsPage, string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, rawURL, nil)
	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}
	rec := httptest.NewRecorder()
	handlerFn(rec, req)

	var page apiComponentsPage
	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
			t.Fatalf("decode response for %s: %v; body: %s", rawURL, err, rec.Body.String())
		}
	}
	return rec.Code, page, rec.Body.String()
}

func newComponentsTestHandler(t *testing.T) (*Handler, *registry.RegistryManager) {
	t.Helper()
	rm, _ := newTestRegistryManager(t)
	h := &Handler{registryManager: rm, log: newTestLogger(t)}
	return h, rm
}

// TestGetAllMeshmodelComponents_ResponseEnvelope pins the response envelope: a
// consumer must get page, pageSize, totalCount, and a populated components list.
func TestGetAllMeshmodelComponents_ResponseEnvelope(t *testing.T) {
	h, rm := newComponentsTestHandler(t)
	seedComponentInModel(t, rm, "kubernetes", "Orchestration", "Job", "batch/v1", "Job")

	code, page, raw := callComponents(t, "/api/registry/components", nil, h.GetAllMeshmodelComponents)

	require.Equal(t, http.StatusOK, code, "body: %s", raw)
	require.Equal(t, int64(1), page.TotalCount)
	require.Len(t, page.Components, 1)
	require.Equal(t, 0, page.Page)
	// No page size requested, so the handler reports the default (utils.go).
	require.Equal(t, defaultPageSize, page.PageSize)
}

// TestGetAllMeshmodelComponents_PageSizeSpellings covers the camelCase pageSize
// and the legacy pagesize, both of which getPaginationParams honours, and
// asserts the page is actually limited while totalCount still reflects the whole
// set.
func TestGetAllMeshmodelComponents_PageSizeSpellings(t *testing.T) {
	h, rm := newComponentsTestHandler(t)
	seedComponentInModel(t, rm, "kubernetes", "Orchestration", "Job", "batch/v1", "Job")
	seedComponentInModel(t, rm, "kubernetes", "Orchestration", "Pod", "v1", "Pod")
	seedComponentInModel(t, rm, "kubernetes", "Orchestration", "Service", "v1", "Service")

	for _, param := range []string{"pageSize=2", "pagesize=2"} {
		t.Run(param, func(t *testing.T) {
			code, page, raw := callComponents(t, "/api/registry/components?"+param, nil, h.GetAllMeshmodelComponents)
			require.Equal(t, http.StatusOK, code, "body: %s", raw)
			require.Equal(t, int64(3), page.TotalCount, "totalCount must count all matches, not just the page")
			require.Equal(t, 2, page.PageSize)
			require.Len(t, page.Components, 2, "results must be limited to the page size")
		})
	}
}

// TestGetAllMeshmodelComponents_PageSizeAll covers the "all" sentinel: every
// result is returned and pageSize is reported as the total count.
func TestGetAllMeshmodelComponents_PageSizeAll(t *testing.T) {
	h, rm := newComponentsTestHandler(t)
	seedComponentInModel(t, rm, "kubernetes", "Orchestration", "Job", "batch/v1", "Job")
	seedComponentInModel(t, rm, "kubernetes", "Orchestration", "Pod", "v1", "Pod")
	seedComponentInModel(t, rm, "kubernetes", "Orchestration", "Service", "v1", "Service")

	for _, param := range []string{"pagesize=all", "pageSize=all"} {
		t.Run(param, func(t *testing.T) {
			code, page, raw := callComponents(t, "/api/registry/components?"+param, nil, h.GetAllMeshmodelComponents)
			require.Equal(t, http.StatusOK, code, "body: %s", raw)
			require.Equal(t, int64(3), page.TotalCount)
			require.Len(t, page.Components, 3, "all results must be returned")
			require.Equal(t, int(page.TotalCount), page.PageSize, "pageSize=all reports pageSize as totalCount")
		})
	}
}

// TestGetMeshmodelComponentByModel_ScopesToRequestedModel is the guarantee the
// future model-scoped retrieval path depends on: /models/{model}/components must
// return only the requested model's components.
func TestGetMeshmodelComponentByModel_ScopesToRequestedModel(t *testing.T) {
	h, rm := newComponentsTestHandler(t)
	seedComponentInModel(t, rm, "kubernetes", "Orchestration", "Job", "batch/v1", "Job")
	seedComponentInModel(t, rm, "aws", "Cloud", "Bucket", "s3/v1", "Bucket")

	t.Run("returns only the requested model's components", func(t *testing.T) {
		code, page, raw := callComponents(t, "/api/registry/models/kubernetes/components",
			map[string]string{"model": "kubernetes"}, h.GetMeshmodelComponentByModel)
		require.Equal(t, http.StatusOK, code, "body: %s", raw)
		require.Equal(t, int64(1), page.TotalCount)
		require.Len(t, page.Components, 1)
		require.Equal(t, "Job", page.Components[0].DisplayName)
		require.Equal(t, "kubernetes", page.Components[0].Model.Name)
	})

	t.Run("a different model returns only its own components", func(t *testing.T) {
		code, page, raw := callComponents(t, "/api/registry/models/aws/components",
			map[string]string{"model": "aws"}, h.GetMeshmodelComponentByModel)
		require.Equal(t, http.StatusOK, code, "body: %s", raw)
		require.Equal(t, int64(1), page.TotalCount)
		require.Len(t, page.Components, 1)
		require.Equal(t, "Bucket", page.Components[0].DisplayName)
	})

	t.Run("an unknown model returns an empty list", func(t *testing.T) {
		code, page, raw := callComponents(t, "/api/registry/models/nonexistent/components",
			map[string]string{"model": "nonexistent"}, h.GetMeshmodelComponentByModel)
		require.Equal(t, http.StatusOK, code, "body: %s", raw)
		require.Equal(t, int64(0), page.TotalCount)
		require.Empty(t, page.Components)
	})
}

// TestGetAllMeshmodelComponents_EmptyResult pins that no matches is a valid,
// well-formed 200 with an empty list, not an error or a broken envelope.
func TestGetAllMeshmodelComponents_EmptyResult(t *testing.T) {
	h, _ := newComponentsTestHandler(t)

	code, page, raw := callComponents(t, "/api/registry/components", nil, h.GetAllMeshmodelComponents)

	require.Equal(t, http.StatusOK, code, "body: %s", raw)
	require.Equal(t, int64(0), page.TotalCount)
	require.Empty(t, page.Components)
}
