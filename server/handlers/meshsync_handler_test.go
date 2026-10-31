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

func TestGetMeshSyncResourceByID_InvalidUUID(t *testing.T) {
	h := newMockHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/not-a-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-uuid"})
	rec := httptest.NewRecorder()

	h.GetMeshSyncResourceByID(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid UUID, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	assertJSONErrorEnvelope(t, rec.Result())
}

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

func TestDeleteMeshSyncResource_InvalidUUID(t *testing.T) {
	h := newMockHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/system/meshsync/resources/invalid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
	rec := httptest.NewRecorder()

	h.DeleteMeshSyncResource(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid UUID, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	assertJSONErrorEnvelope(t, rec.Result())
}

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
