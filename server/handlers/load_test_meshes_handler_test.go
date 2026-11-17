package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSMPTestRegistryManager(t *testing.T) (*registry.RegistryManager, *database.Handler) {
	t.Helper()
	db, err := database.New(database.Options{
		Filename: ":memory:",
		Engine:   "sqlite",
	})
	if err != nil {
		if strings.Contains(err.Error(), "requires cgo") || strings.Contains(err.Error(), "go-sqlite3 requires cgo") {
			t.Skip("skipping SQLite test: go-sqlite3 requires cgo (disabled on this platform)")
		}
		require.NoError(t, err, "failed to create in-memory database")
	}
	rm, err := registry.NewRegistryManager(&db)
	require.NoError(t, err, "failed to create registry manager")
	return rm, &db
}

func seedTestSMPModel(t *testing.T, rm *registry.RegistryManager, modelName, displayName string) {
	t.Helper()
	conn := connection.Connection{
		Name:    "test-registrant-" + modelName,
		Kind:    modelName,
		Type:    "platform",
		SubType: "service-mesh",
		Status:  connection.ConnectionStatusConnected,
	}
	m := &model.ModelDefinition{
		Name:          modelName,
		DisplayName:   displayName,
		SchemaVersion: "models.meshery.io/v1beta1",
		Version:       "v1.0.0",
		Model:         model.Model{Version: "v1.0.0"},
		Category:      category.CategoryDefinition{Name: "Service Mesh"},
		Status:        model.Enabled,
		Metadata:      &model.ModelDefinition_Metadata{},
	}
	id, err := m.GenerateID()
	require.NoError(t, err)
	m.ID = id
	_, _, err = rm.RegisterEntity(registry.RegistrantHostToV1beta3(conn), m)
	require.NoError(t, err, "seedTestSMPModel: RegisterEntity failed")
}

func TestGetSMPServiceMeshes_MethodValidation(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{"POST method not allowed", http.MethodPost, http.StatusMethodNotAllowed},
		{"PUT method not allowed", http.MethodPut, http.StatusMethodNotAllowed},
		{"DELETE method not allowed", http.MethodDelete, http.StatusMethodNotAllowed},
		{"PATCH method not allowed", http.MethodPatch, http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			req := httptest.NewRequest(tt.method, "/api/mesh", nil)
			rec := httptest.NewRecorder()

			h.GetSMPServiceMeshes(rec, req, nil, nil, nil)
			resp := rec.Result()
			defer func() {
				_ = resp.Body.Close()
			}()

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestGetSMPServiceMeshes_NilRegistryManager_ReturnsInternalServerError(t *testing.T) {
	h := &Handler{
		registryManager: nil,
		log:             newTestLogger(t),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/mesh", nil)
	rec := httptest.NewRecorder()

	h.GetSMPServiceMeshes(rec, req, nil, nil, nil)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrRetrieveDataCode, errBody.Code)
	assert.NotEmpty(t, errBody.Error)
	assert.NotEmpty(t, errBody.ProbableCause)
	assert.NotEmpty(t, errBody.SuggestedRemediation)
}

func TestGetSMPServiceMeshes_NilLogger_DoesNotPanic(t *testing.T) {
	h := &Handler{
		registryManager: nil,
		log:             nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/mesh", nil)
	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		h.GetSMPServiceMeshes(rec, req, nil, nil, nil)
	})
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrRetrieveDataCode, errBody.Code)
}

func TestGetSMPServiceMeshes_EmptyRegistry_Success(t *testing.T) {
	rm, db := newSMPTestRegistryManager(t)
	t.Cleanup(func() { _ = db.DBClose() })

	h := &Handler{
		registryManager: rm,
		log:             newTestLogger(t),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/mesh", nil)
	rec := httptest.NewRecorder()

	h.GetSMPServiceMeshes(rec, req, nil, nil, nil)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var meshes SMPMeshes
	err := json.Unmarshal(rec.Body.Bytes(), &meshes)
	require.NoError(t, err)
	assert.Empty(t, meshes.AvailableMeshes)
}

func TestGetSMPServiceMeshes_ModelExtractionAndSorting(t *testing.T) {
	rm, db := newSMPTestRegistryManager(t)
	t.Cleanup(func() { _ = db.DBClose() })

	// Seed multiple models with display names, fallbacks, and duplicates
	seedTestSMPModel(t, rm, "istio", "Istio")
	seedTestSMPModel(t, rm, "consul", "Consul")
	seedTestSMPModel(t, rm, "linkerd", "")        // DisplayName is empty; falls back to Name ("linkerd")
	seedTestSMPModel(t, rm, "istio-dup", "Istio") // Duplicate display name "Istio"
	seedTestSMPModel(t, rm, "app-mesh", "AWS App Mesh")

	h := &Handler{
		registryManager: rm,
		log:             newTestLogger(t),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/mesh", nil)
	rec := httptest.NewRecorder()

	h.GetSMPServiceMeshes(rec, req, nil, nil, nil)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var meshes SMPMeshes
	err := json.Unmarshal(rec.Body.Bytes(), &meshes)
	require.NoError(t, err)

	// Expected sorted list: "AWS App Mesh", "Consul", "Istio", "linkerd"
	// "Istio" should be deduplicated (present once), "linkerd" fallback used
	expected := []string{"AWS App Mesh", "Consul", "Istio", "linkerd"}
	assert.Equal(t, expected, meshes.AvailableMeshes)
}
