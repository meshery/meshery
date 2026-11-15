// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshsync/pkg/model"
)

// newMockHandler initializes a mock Handler instance configured for testing.
func newMockHandler(t *testing.T) *Handler {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	return &Handler{
		config: &models.HandlerConfig{
			Providers:          map[string]models.Provider{},
			ProviderCookieName: "meshery-provider",
		},
		log: log,
	}
}

// newMeshSyncTestFixture sets up an in-memory SQLite database, migrates MeshSync models,
// and returns an initialized Handler along with a mock provider and database handle.
func newMeshSyncTestFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider, database.Handler) {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		if strings.Contains(err.Error(), "go-sqlite3 requires cgo") {
			t.Skip("skipping SQLite test: go-sqlite3 requires cgo (disabled on this platform)")
		}
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() {
		if db.DB != nil {
			if sqlDB, err := db.DB.DB(); err == nil && sqlDB != nil {
				_ = sqlDB.Close()
			}
		}
	})

	if err := db.AutoMigrate(
		&model.KubernetesResource{},
		&model.KubernetesResourceObjectMeta{},
		&model.KubernetesKeyValue{},
		&model.KubernetesResourceSpec{},
		&model.KubernetesResourceStatus{},
	); err != nil {
		t.Fatalf("migrate meshsync tables: %v", err)
	}

	h := newMockHandler(t)

	provider := &models.DefaultLocalProvider{
		GenericPersister: &db,
	}

	return h, provider, db
}

// TestGetMeshSyncResourceByID_InvalidUUID verifies that passing an invalid UUID returns HTTP 400 Bad Request.
func TestGetMeshSyncResourceByID_InvalidUUID(t *testing.T) {
	h := newMockHandler(t)
	provider := &models.DefaultLocalProvider{}

	req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/not-a-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-uuid"})
	rec := httptest.NewRecorder()

	h.GetMeshSyncResourceByID(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid UUID, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	assertJSONErrorEnvelope(t, rec.Result())
}

// TestGetMeshSyncResourceByID_NotFound verifies that querying a non-existent resource returns HTTP 404 Not Found.
func TestGetMeshSyncResourceByID_NotFound(t *testing.T) {
	h, provider, _ := newMeshSyncTestFixture(t)

	randomID := uuid.Must(uuid.NewV4()).String()
	req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/"+randomID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": randomID})
	rec := httptest.NewRecorder()

	h.GetMeshSyncResourceByID(rec, req, nil, nil, provider)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent resource, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	assertJSONErrorEnvelope(t, rec.Result())
}

// TestGetMeshSyncResourceByID_Success verifies that querying an existing resource returns HTTP 200 OK with JSON data.
func TestGetMeshSyncResourceByID_Success(t *testing.T) {
	h, provider, db := newMeshSyncTestFixture(t)

	resID := uuid.Must(uuid.NewV4()).String()
	meta := model.KubernetesResourceObjectMeta{
		ID:   resID,
		Name: "test-resource",
		UID:  resID,
	}
	if err := db.Create(&meta).Error; err != nil {
		t.Fatalf("failed to seed resource meta: %v", err)
	}

	resource := model.KubernetesResource{
		ID:        resID,
		ClusterID: "test-cluster",
		ComponentMetadata: map[string]interface{}{
			"name": "test-resource",
		},
	}
	// Insert via Table to preserve exact UUID primary key without BeforeCreate hook overriding it
	if err := db.Table("kubernetes_resources").Create(&resource).Error; err != nil {
		t.Fatalf("failed to seed resource: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/"+resID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": resID})
	rec := httptest.NewRecorder()

	h.GetMeshSyncResourceByID(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

// TestDeleteMeshSyncResource_InvalidUUID verifies that passing an invalid resource ID returns HTTP 400 Bad Request.
func TestDeleteMeshSyncResource_InvalidUUID(t *testing.T) {
	h := newMockHandler(t)
	provider := &models.DefaultLocalProvider{}

	req := httptest.NewRequest(http.MethodDelete, "/api/system/meshsync/resources/invalid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
	rec := httptest.NewRecorder()

	h.DeleteMeshSyncResource(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid UUID, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	assertJSONErrorEnvelope(t, rec.Result())
}

// TestDeleteMeshSyncResource_NotFound verifies that attempting to delete a non-existent resource returns HTTP 404 Not Found.
func TestDeleteMeshSyncResource_NotFound(t *testing.T) {
	h, provider, _ := newMeshSyncTestFixture(t)

	randomID := uuid.Must(uuid.NewV4()).String()
	req := httptest.NewRequest(http.MethodDelete, "/api/system/meshsync/resources/"+randomID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": randomID})
	rec := httptest.NewRecorder()

	h.DeleteMeshSyncResource(rec, req, nil, nil, provider)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent resource, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	assertJSONErrorEnvelope(t, rec.Result())
}

// TestDeleteMeshSyncResource_Success verifies that deleting an existing resource removes it and returns HTTP 200 OK.
func TestDeleteMeshSyncResource_Success(t *testing.T) {
	h, provider, db := newMeshSyncTestFixture(t)

	resID := uuid.Must(uuid.NewV4()).String()
	resource := model.KubernetesResource{
		ID:        resID,
		ClusterID: "test-cluster",
	}
	if err := db.Table("kubernetes_resources").Create(&resource).Error; err != nil {
		t.Fatalf("failed to seed resource: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/system/meshsync/resources/"+resID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": resID})
	rec := httptest.NewRecorder()

	h.DeleteMeshSyncResource(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}

	var response struct {
		Deleted bool `json:"deleted"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal delete response: %v (body=%q)", err, rec.Body.String())
	}
	if !response.Deleted {
		t.Errorf("expected deleted to be true, got false")
	}

	// Verify the entity is truly deleted from the database
	var check model.KubernetesResource
	err := db.Where("id = ?", resID).First(&check).Error
	if err == nil {
		t.Errorf("expected resource to be deleted from database, but it still exists")
	}
}

// TestGetMeshSyncResourceByID_Base64ID_Success verifies that querying a resource by a valid base64 composite ID succeeds.
func TestGetMeshSyncResourceByID_Base64ID_Success(t *testing.T) {
	h, provider, db := newMeshSyncTestFixture(t)

	base64ID := base64.StdEncoding.EncodeToString([]byte("test-cluster.Pod.v1.default.test-base64-resource"))
	meta := model.KubernetesResourceObjectMeta{
		ID:   base64ID,
		Name: "test-base64-resource",
		UID:  base64ID,
	}
	if err := db.Create(&meta).Error; err != nil {
		t.Fatalf("failed to seed resource meta: %v", err)
	}

	resource := model.KubernetesResource{
		ID:        base64ID,
		ClusterID: "test-cluster",
		ComponentMetadata: map[string]interface{}{
			"name": "test-base64-resource",
		},
	}
	if err := db.Table("kubernetes_resources").Create(&resource).Error; err != nil {
		t.Fatalf("failed to seed resource: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/"+base64ID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": base64ID})
	rec := httptest.NewRecorder()

	h.GetMeshSyncResourceByID(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

// TestDeleteMeshSyncResource_Base64ID_Success verifies that deleting a resource by a valid base64 composite ID succeeds.
func TestDeleteMeshSyncResource_Base64ID_Success(t *testing.T) {
	h, provider, db := newMeshSyncTestFixture(t)

	base64ID := base64.StdEncoding.EncodeToString([]byte("test-cluster.Pod.v1.default.test-delete-base64"))
	resource := model.KubernetesResource{
		ID:        base64ID,
		ClusterID: "test-cluster",
	}
	if err := db.Table("kubernetes_resources").Create(&resource).Error; err != nil {
		t.Fatalf("failed to seed resource: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/system/meshsync/resources/"+base64ID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": base64ID})
	rec := httptest.NewRecorder()

	h.DeleteMeshSyncResource(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}

	var response struct {
		Deleted bool `json:"deleted"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal delete response: %v (body=%q)", err, rec.Body.String())
	}
	if !response.Deleted {
		t.Errorf("expected deleted to be true, got false")
	}
}

// TestGetMeshSyncResources_OrderBy verifies that the order query parameter properly
// applies ORDER BY sorting to the returned resources.
func TestGetMeshSyncResources_OrderBy(t *testing.T) {
	h, provider, db := newMeshSyncTestFixture(t)

	idA := uuid.Must(uuid.NewV4()).String()
	metaA := model.KubernetesResourceObjectMeta{
		ID:   idA,
		Name: "alpha-resource",
		UID:  idA,
	}
	if err := db.Create(&metaA).Error; err != nil {
		t.Fatalf("failed to seed metaA: %v", err)
	}
	resA := model.KubernetesResource{
		ID:        idA,
		ClusterID: "test-cluster",
		ComponentMetadata: map[string]interface{}{
			"name": "alpha-resource",
		},
	}
	if err := db.Table("kubernetes_resources").Create(&resA).Error; err != nil {
		t.Fatalf("failed to seed resA: %v", err)
	}

	idB := uuid.Must(uuid.NewV4()).String()
	metaB := model.KubernetesResourceObjectMeta{
		ID:   idB,
		Name: "beta-resource",
		UID:  idB,
	}
	if err := db.Create(&metaB).Error; err != nil {
		t.Fatalf("failed to seed metaB: %v", err)
	}
	resB := model.KubernetesResource{
		ID:        idB,
		ClusterID: "test-cluster",
		ComponentMetadata: map[string]interface{}{
			"name": "beta-resource",
		},
	}
	if err := db.Table("kubernetes_resources").Create(&resB).Error; err != nil {
		t.Fatalf("failed to seed resB: %v", err)
	}

	// Test 1: Order by name descending (beta should come before alpha)
	reqDesc := httptest.NewRequest(http.MethodGet, `/api/system/meshsync/resources?clusterIds=["test-cluster"]&order=name%20desc`, nil)
	recDesc := httptest.NewRecorder()
	h.GetMeshSyncResources(recDesc, reqDesc, nil, nil, provider)

	if recDesc.Code != http.StatusOK {
		t.Fatalf("expected status 200 for desc query, got %d (body=%q)", recDesc.Code, recDesc.Body.String())
	}

	var respDesc models.MeshSyncResourcesAPIResponse
	if err := json.Unmarshal(recDesc.Body.Bytes(), &respDesc); err != nil {
		t.Fatalf("failed to unmarshal desc response: %v (body=%q)", err, recDesc.Body.String())
	}
	if len(respDesc.Resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(respDesc.Resources))
	}
	if respDesc.Resources[0].ID != idB || respDesc.Resources[1].ID != idA {
		t.Errorf("expected [beta, alpha] for name desc, got [%s, %s]", respDesc.Resources[0].ID, respDesc.Resources[1].ID)
	}

	// Test 2: Order by name ascending (alpha should come before beta)
	reqAsc := httptest.NewRequest(http.MethodGet, `/api/system/meshsync/resources?clusterIds=["test-cluster"]&order=name%20asc`, nil)
	recAsc := httptest.NewRecorder()
	h.GetMeshSyncResources(recAsc, reqAsc, nil, nil, provider)

	if recAsc.Code != http.StatusOK {
		t.Fatalf("expected status 200 for asc query, got %d (body=%q)", recAsc.Code, recAsc.Body.String())
	}

	var respAsc models.MeshSyncResourcesAPIResponse
	if err := json.Unmarshal(recAsc.Body.Bytes(), &respAsc); err != nil {
		t.Fatalf("failed to unmarshal asc response: %v (body=%q)", err, recAsc.Body.String())
	}
	if len(respAsc.Resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(respAsc.Resources))
	}
	if respAsc.Resources[0].ID != idA || respAsc.Resources[1].ID != idB {
		t.Errorf("expected [alpha, beta] for name asc, got [%s, %s]", respAsc.Resources[0].ID, respAsc.Resources[1].ID)
	}

	// Test 3: No order param (exercises empty-guard behavior)
	reqNoOrder := httptest.NewRequest(http.MethodGet, `/api/system/meshsync/resources?clusterIds=["test-cluster"]`, nil)
	recNoOrder := httptest.NewRecorder()
	h.GetMeshSyncResources(recNoOrder, reqNoOrder, nil, nil, provider)

	if recNoOrder.Code != http.StatusOK {
		t.Fatalf("expected status 200 for no-order query, got %d (body=%q)", recNoOrder.Code, recNoOrder.Body.String())
	}

	var respNoOrder models.MeshSyncResourcesAPIResponse
	if err := json.Unmarshal(recNoOrder.Body.Bytes(), &respNoOrder); err != nil {
		t.Fatalf("failed to unmarshal no-order response: %v (body=%q)", err, recNoOrder.Body.String())
	}
	if len(respNoOrder.Resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(respNoOrder.Resources))
	}

	// Test 4: Bare order param (exercises empty-guard behavior with bare order=)
	reqBareOrder := httptest.NewRequest(http.MethodGet, `/api/system/meshsync/resources?clusterIds=["test-cluster"]&order=`, nil)
	recBareOrder := httptest.NewRecorder()
	h.GetMeshSyncResources(recBareOrder, reqBareOrder, nil, nil, provider)

	if recBareOrder.Code != http.StatusOK {
		t.Fatalf("expected status 200 for bare-order query, got %d (body=%q)", recBareOrder.Code, recBareOrder.Body.String())
	}

	var respBareOrder models.MeshSyncResourcesAPIResponse
	if err := json.Unmarshal(recBareOrder.Body.Bytes(), &respBareOrder); err != nil {
		t.Fatalf("failed to unmarshal bare-order response: %v (body=%q)", err, recBareOrder.Body.String())
	}
	if len(respBareOrder.Resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(respBareOrder.Resources))
	}
}

// TestMeshSyncResources_OrderInputSanitization verifies that allowed sorting columns
// are safely translated to snake_case while invalid columns and injection attempts are rejected.
func TestMeshSyncResources_OrderInputSanitization(t *testing.T) {
	validColumns := []string{"creation_timestamp", "name", "kind", "model", "api_version", "namespace"}

	tests := []struct {
		input    string
		expected string
	}{
		{"name desc", "name desc"},
		{"name asc", "name asc"},
		{"creationTimestamp desc", "creation_timestamp desc"},
		{"apiVersion asc", "api_version asc"},
		{"kind desc", "kind desc"},
		{"model asc", "model asc"},
		{"namespace desc", "namespace desc"},
		{"invalidColumn asc", ""},
		{"name; DROP TABLE kubernetes_resources; -- desc", ""},
		{"", ""},
		{"singleword", ""},
		{"three words here", ""},
	}

	for _, tc := range tests {
		got := models.SanitizeOrderInput(tc.input, validColumns)
		if got != tc.expected {
			t.Errorf("SanitizeOrderInput(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
