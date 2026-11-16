package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func seedTestFilter(t *testing.T, provider *models.DefaultLocalProvider, id uuid.UUID) {
	t.Helper()
	name := "test-filter"
	filter := &models.MesheryFilter{
		ID:         &id,
		Name:       name,
		FilterFile: []byte("wasm-binary-content"),
		Visibility: models.Private,
	}
	if _, err := provider.MesheryFilterPersister.SaveMesheryFilter(filter); err != nil {
		t.Fatalf("seed filter: %v", err)
	}
}

// TestGetMesheryFilterHandler_LocalProvider drives GET /api/filter/{id}
// through the Local provider and a real database. A filter that does not
// exist must be reported as 404 with an error body; a malformed ID returns 400
// and a database failure is a 500 server error.
func TestGetMesheryFilterHandler_LocalProvider(t *testing.T) {
	tests := []struct {
		name       string
		pathID     string
		setup      func(t *testing.T, provider *models.DefaultLocalProvider, id uuid.UUID)
		wantStatus int
	}{
		{
			name:       "stored filter is returned",
			setup:      func(t *testing.T, p *models.DefaultLocalProvider, id uuid.UUID) { seedTestFilter(t, p, id) },
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing filter is not found",
			setup:      func(*testing.T, *models.DefaultLocalProvider, uuid.UUID) {},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "malformed ID is a bad request",
			pathID:     "not-a-uuid",
			setup:      func(*testing.T, *models.DefaultLocalProvider, uuid.UUID) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, provider, _ := newDeleteFilterFixture(t)
			id := uuid.Must(uuid.NewV4())
			tc.setup(t, provider, id)

			pathID := id.String()
			if tc.pathID != "" {
				pathID = tc.pathID
			}
			req := httptest.NewRequest(http.MethodGet, "/api/filter/"+pathID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": pathID})
			rec := httptest.NewRecorder()

			h.GetMesheryFilterHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantStatus, rec.Body.String())
			}

			if tc.wantStatus == http.StatusOK {
				var decoded struct {
					ID   *uuid.UUID `json:"id"`
					Name string     `json:"name"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
					t.Fatalf("body did not parse as JSON: %v (body=%q)", err, rec.Body.String())
				}
				if decoded.ID == nil || *decoded.ID != id {
					t.Errorf("filter id = %v, want %s (body=%q)", decoded.ID, id, rec.Body.String())
				}
			}
		})
	}
}

// TestGetMesheryFilterFileHandler_LocalProvider drives GET /api/filter/file/{id}
// ensuring valid WASM retrieval for stored filters, 404 for missing, and 400 for malformed IDs.
func TestGetMesheryFilterFileHandler_LocalProvider(t *testing.T) {
	tests := []struct {
		name       string
		pathID     string
		setup      func(t *testing.T, provider *models.DefaultLocalProvider, id uuid.UUID)
		wantStatus int
	}{
		{
			name:       "stored filter file is returned",
			setup:      func(t *testing.T, p *models.DefaultLocalProvider, id uuid.UUID) { seedTestFilter(t, p, id) },
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing filter file is not found",
			setup:      func(*testing.T, *models.DefaultLocalProvider, uuid.UUID) {},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "malformed ID is a bad request",
			pathID:     "not-a-valid-uuid",
			setup:      func(*testing.T, *models.DefaultLocalProvider, uuid.UUID) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, provider, _ := newDeleteFilterFixture(t)
			id := uuid.Must(uuid.NewV4())
			tc.setup(t, provider, id)

			pathID := id.String()
			if tc.pathID != "" {
				pathID = tc.pathID
			}
			req := httptest.NewRequest(http.MethodGet, "/api/filter/file/"+pathID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": pathID})
			rec := httptest.NewRecorder()

			h.GetMesheryFilterFileHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantStatus, rec.Body.String())
			}

			if tc.wantStatus == http.StatusOK {
				if rec.Header().Get("Content-Type") != "application/wasm" {
					t.Errorf("content-type = %q, want application/wasm", rec.Header().Get("Content-Type"))
				}
				if rec.Body.String() != "wasm-binary-content" {
					t.Errorf("body = %q, want %q", rec.Body.String(), "wasm-binary-content")
				}
			}
		})
	}
}

// TestCloneMesheryFilterHandler_LocalProvider drives POST /api/filter/clone/{id}
// ensuring malformed IDs, empty names, and missing filters return appropriate errors.
func TestCloneMesheryFilterHandler_LocalProvider(t *testing.T) {
	tests := []struct {
		name       string
		pathID     string
		body       string
		setup      func(t *testing.T, provider *models.DefaultLocalProvider, id uuid.UUID)
		wantStatus int
	}{
		{
			name:       "cloning missing filter returns 404",
			body:       `{"name":"cloned-filter"}`,
			setup:      func(*testing.T, *models.DefaultLocalProvider, uuid.UUID) {},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "malformed ID is a bad request",
			pathID:     "not-a-uuid",
			body:       `{"name":"cloned-filter"}`,
			setup:      func(*testing.T, *models.DefaultLocalProvider, uuid.UUID) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty name is a bad request",
			body:       `{"name":""}`,
			setup:      func(t *testing.T, p *models.DefaultLocalProvider, id uuid.UUID) { seedTestFilter(t, p, id) },
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, provider, _ := newDeleteFilterFixture(t)
			id := uuid.Must(uuid.NewV4())
			tc.setup(t, provider, id)

			pathID := id.String()
			if tc.pathID != "" {
				pathID = tc.pathID
			}
			req := httptest.NewRequest(http.MethodPost, "/api/filter/clone/"+pathID, bytes.NewBufferString(tc.body))
			req = mux.SetURLVars(req, map[string]string{"id": pathID})
			rec := httptest.NewRecorder()

			h.CloneMesheryFilterHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantStatus, rec.Body.String())
			}
		})
	}
}

// TestDeleteMesheryFilterHandler_MalformedUUID tests that DELETE /api/filter/{id}
// rejects non-UUID path variables with 400 Bad Request via requirePathUUID.
func TestDeleteMesheryFilterHandler_MalformedUUID(t *testing.T) {
	h, provider, _ := newDeleteFilterFixture(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/filter/not-a-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-uuid"})
	rec := httptest.NewRecorder()

	h.DeleteMesheryFilterHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}
