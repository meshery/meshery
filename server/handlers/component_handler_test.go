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
	// The registry types seeded below are only available as v1beta1 here: the
	// v1beta3 ComponentDefinition's Model field is a *v1beta1 ModelDefinition,
	// the v1beta3 category and model packages are not generated, and
	// registry.RegistrantHostToV1beta3 takes a v1beta1 Connection. So these
	// stay v1beta1 by necessity, matching the existing registry-backed tests in
	// policy_relationship_handler_test.go.
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	v1beta3comp "github.com/meshery/schemas/models/v1beta3/component"
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
	if err != nil {
		t.Fatalf("seedComponentInModel(%s/%s): generate id: %v", modelName, kind, err)
	}
	comp.ID = id
	if _, _, err := rm.RegisterEntity(registry.RegistrantHostToV1beta3(conn), &comp); err != nil {
		t.Fatalf("seedComponentInModel(%s/%s): RegisterEntity: %v", modelName, kind, err)
	}
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
// route is bound to), sets any mux path vars, and returns the status, the
// decoded response, the raw top-level keys, and the raw body for diagnostics.
// The body is decoded unconditionally so an empty or non-JSON response fails the
// test rather than passing as zero-valued fields.
func callComponents(t *testing.T, rawURL string, vars map[string]string, handlerFn http.HandlerFunc) (int, apiComponentsPage, map[string]json.RawMessage, string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, rawURL, nil)
	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}
	rec := httptest.NewRecorder()
	handlerFn(rec, req)

	body := rec.Body.Bytes()
	var page apiComponentsPage
	if err := json.Unmarshal(body, &page); err != nil {
		t.Fatalf("decode response for %s: %v; body: %q", rawURL, err, string(body))
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(body, &keys); err != nil {
		t.Fatalf("decode response envelope for %s: %v; body: %q", rawURL, err, string(body))
	}
	return rec.Code, page, keys, string(body)
}

// assertEnvelopeKeys checks the response carries every top-level envelope key,
// so an omitted field is caught rather than passing as a zero value.
func assertEnvelopeKeys(t *testing.T, keys map[string]json.RawMessage, body string) {
	t.Helper()
	for _, k := range []string{"page", "pageSize", "totalCount", "components"} {
		if _, ok := keys[k]; !ok {
			t.Errorf("response envelope is missing the %q key; body: %s", k, body)
		}
	}
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

	code, page, keys, body := callComponents(t, "/api/registry/components", nil, h.GetAllMeshmodelComponents)

	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", code, body)
	}
	assertEnvelopeKeys(t, keys, body)
	if page.TotalCount != 1 {
		t.Errorf("totalCount = %d, want 1", page.TotalCount)
	}
	if len(page.Components) != 1 {
		t.Errorf("components length = %d, want 1", len(page.Components))
	}
	if page.Page != 0 {
		t.Errorf("page = %d, want 0", page.Page)
	}
	// No page size requested, so the handler reports the default (utils.go).
	if page.PageSize != defaultPageSize {
		t.Errorf("pageSize = %d, want default %d", page.PageSize, defaultPageSize)
	}
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
			code, page, keys, body := callComponents(t, "/api/registry/components?"+param, nil, h.GetAllMeshmodelComponents)
			if code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body: %s", code, body)
			}
			assertEnvelopeKeys(t, keys, body)
			if page.TotalCount != 3 {
				t.Errorf("totalCount = %d, want 3 (all matches, not just the page)", page.TotalCount)
			}
			if page.PageSize != 2 {
				t.Errorf("pageSize = %d, want 2", page.PageSize)
			}
			if len(page.Components) != 2 {
				t.Errorf("components length = %d, want 2 (limited to the page size)", len(page.Components))
			}
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
			code, page, keys, body := callComponents(t, "/api/registry/components?"+param, nil, h.GetAllMeshmodelComponents)
			if code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body: %s", code, body)
			}
			assertEnvelopeKeys(t, keys, body)
			if page.TotalCount != 3 {
				t.Errorf("totalCount = %d, want 3", page.TotalCount)
			}
			if len(page.Components) != 3 {
				t.Errorf("components length = %d, want 3 (all results)", len(page.Components))
			}
			if page.PageSize != int(page.TotalCount) {
				t.Errorf("pageSize = %d, want %d (pageSize=all reports pageSize as totalCount)", page.PageSize, page.TotalCount)
			}
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
		code, page, keys, body := callComponents(t, "/api/registry/models/kubernetes/components",
			map[string]string{"model": "kubernetes"}, h.GetMeshmodelComponentByModel)
		if code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body: %s", code, body)
		}
		assertEnvelopeKeys(t, keys, body)
		if page.TotalCount != 1 {
			t.Errorf("totalCount = %d, want 1", page.TotalCount)
		}
		if len(page.Components) != 1 {
			t.Fatalf("components length = %d, want 1", len(page.Components))
		}
		if page.Components[0].DisplayName != "Job" {
			t.Errorf("component displayName = %q, want %q", page.Components[0].DisplayName, "Job")
		}
		if page.Components[0].Model.Name != "kubernetes" {
			t.Errorf("component model = %q, want %q", page.Components[0].Model.Name, "kubernetes")
		}
	})

	t.Run("a different model returns only its own components", func(t *testing.T) {
		code, page, keys, body := callComponents(t, "/api/registry/models/aws/components",
			map[string]string{"model": "aws"}, h.GetMeshmodelComponentByModel)
		if code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body: %s", code, body)
		}
		assertEnvelopeKeys(t, keys, body)
		if page.TotalCount != 1 {
			t.Errorf("totalCount = %d, want 1", page.TotalCount)
		}
		if len(page.Components) != 1 {
			t.Fatalf("components length = %d, want 1", len(page.Components))
		}
		if page.Components[0].DisplayName != "Bucket" {
			t.Errorf("component displayName = %q, want %q", page.Components[0].DisplayName, "Bucket")
		}
	})

	t.Run("an unknown model returns an empty list", func(t *testing.T) {
		code, page, keys, body := callComponents(t, "/api/registry/models/nonexistent/components",
			map[string]string{"model": "nonexistent"}, h.GetMeshmodelComponentByModel)
		if code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body: %s", code, body)
		}
		assertEnvelopeKeys(t, keys, body)
		if page.TotalCount != 0 {
			t.Errorf("totalCount = %d, want 0", page.TotalCount)
		}
		if len(page.Components) != 0 {
			t.Errorf("components length = %d, want 0", len(page.Components))
		}
	})
}

// TestGetAllMeshmodelComponents_EmptyResult pins that no matches is a valid,
// well-formed 200 whose envelope still carries every key, not an error or a
// blank body.
func TestGetAllMeshmodelComponents_EmptyResult(t *testing.T) {
	h, _ := newComponentsTestHandler(t)

	code, page, keys, body := callComponents(t, "/api/registry/components", nil, h.GetAllMeshmodelComponents)

	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", code, body)
	}
	assertEnvelopeKeys(t, keys, body)
	if page.TotalCount != 0 {
		t.Errorf("totalCount = %d, want 0", page.TotalCount)
	}
	if len(page.Components) != 0 {
		t.Errorf("components length = %d, want 0", len(page.Components))
	}
}
