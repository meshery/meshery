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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
)

func newCloneFilterFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider, database.Handler) {
	t.Helper()

	db, err := database.New(database.Options{
		Engine:   database.SQLITE,
		Filename: filepath.Join(t.TempDir(), "meshery.db"),
	})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&models.MesheryFilter{}); err != nil {
		t.Fatalf("migrate tables: %v", err)
	}

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	provider := &models.DefaultLocalProvider{
		MesheryFilterPersister: &models.MesheryFilterPersister{DB: &db},
		GenericPersister:       &db,
	}
	return h, provider, db
}

func cloneFilter(t *testing.T, h *Handler, provider models.Provider, filterID string, body []byte) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody *bytes.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}
	req := httptest.NewRequest(http.MethodPost, "/api/filter/clone/"+filterID, reqBody)
	req = mux.SetURLVars(req, map[string]string{"id": filterID})
	rec := httptest.NewRecorder()

	h.CloneMesheryFilterHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	return rec
}

func TestCloneFilterNilProvider(t *testing.T) {
	h, _, _ := newCloneFilterFixture(t)
	body, _ := json.Marshal(&models.MesheryCloneFilterRequestBody{Name: "cloned"})
	rec := cloneFilter(t, h, nil, uuid.Must(uuid.NewV4()).String(), body)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrNilProviderCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrNilProviderCode, rec.Body.String())
	}
}

func TestCloneFilterMalformedIDIsBadRequest(t *testing.T) {
	h, provider, _ := newCloneFilterFixture(t)
	body, _ := json.Marshal(&models.MesheryCloneFilterRequestBody{Name: "cloned"})
	rec := cloneFilter(t, h, provider, "not-a-valid-uuid", body)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
		t.Errorf("expected error body to contain %s, got: %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
	}
}

func TestCloneFilterNilUUIDIsBadRequest(t *testing.T) {
	h, provider, _ := newCloneFilterFixture(t)
	body, _ := json.Marshal(&models.MesheryCloneFilterRequestBody{Name: "cloned"})
	rec := cloneFilter(t, h, provider, uuid.Nil.String(), body)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
		t.Errorf("expected error body to contain %s, got: %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
	}
}

func TestCloneFilterNilOrEmptyBodyIsBadRequest(t *testing.T) {
	h, provider, _ := newCloneFilterFixture(t)
	validID := uuid.Must(uuid.NewV4()).String()

	// Empty body
	rec := cloneFilter(t, h, provider, validID, []byte(""))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("empty body status = %d, want %d. body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}

	// JSON "null"
	rec = cloneFilter(t, h, provider, validID, []byte("null"))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("null body status = %d, want %d. body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestCloneFilterNotFound(t *testing.T) {
	h, provider, _ := newCloneFilterFixture(t)
	body, _ := json.Marshal(&models.MesheryCloneFilterRequestBody{Name: "cloned"})
	randomUUID := uuid.Must(uuid.NewV4()).String()

	rec := cloneFilter(t, h, provider, randomUUID, body)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), models.ErrResultNotFoundCode) {
		t.Errorf("expected error body to contain %s, got: %s", models.ErrResultNotFoundCode, rec.Body.String())
	}
}

func TestCloneFilterSuccess(t *testing.T) {
	h, provider, db := newCloneFilterFixture(t)

	sourceID := uuid.Must(uuid.NewV4())
	sourceFilter := models.MesheryFilter{
		ID:         &sourceID,
		Name:       "original-filter",
		Visibility: models.Published,
		FilterFile: []byte("\x00asm wasm binary"),
	}
	if err := db.Save(&sourceFilter).Error; err != nil {
		t.Fatalf("seed filter: %v", err)
	}

	body, _ := json.Marshal(&models.MesheryCloneFilterRequestBody{Name: "cloned-filter"})
	rec := cloneFilter(t, h, provider, sourceID.String(), body)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var cloned []models.MesheryFilter
	if err := json.Unmarshal(rec.Body.Bytes(), &cloned); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(cloned) == 0 {
		t.Fatalf("expected non-empty response")
	}
	if cloned[0].ID == nil || *cloned[0].ID == sourceID {
		t.Errorf("cloned ID = %v, want new ID different from %s", cloned[0].ID, sourceID)
	}
	if cloned[0].Name != "cloned-filter" {
		t.Errorf("cloned name = %q, want %q", cloned[0].Name, "cloned-filter")
	}
	if cloned[0].Visibility != models.Private {
		t.Errorf("cloned visibility = %q, want %q", cloned[0].Visibility, models.Private)
	}
}
