package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/httputil"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
)

// newDeleteFilterFixture wires the filter handler over an isolated in-memory
// SQLite database using DefaultLocalProvider.
func newDeleteFilterFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider, database.Handler) {
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

func deleteFilter(t *testing.T, h *Handler, provider models.Provider, filterID string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodDelete, "/api/filter/"+filterID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": filterID})
	rec := httptest.NewRecorder()

	h.DeleteMesheryFilterHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	return rec
}

type mockRemoteFilterProvider struct {
	models.DefaultLocalProvider
	err error
}

func (m *mockRemoteFilterProvider) DeleteMesheryFilter(_ *http.Request, _ string) ([]byte, error) {
	return nil, m.err
}

// TestDeleteFilterMalformedIDIsBadRequest verifies that a malformed non-UUID
// filter ID returns HTTP 400 Bad Request instead of being parsed as uuid.Nil
// and returning HTTP 200.
func TestDeleteFilterMalformedIDIsBadRequest(t *testing.T) {
	h, provider, _ := newDeleteFilterFixture(t)

	rec := deleteFilter(t, h, provider, "not-a-valid-uuid")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

// TestDeleteFilterNilUUIDIsBadRequest verifies that passing the zero UUID
// (uuid.Nil) is rejected as HTTP 400 Bad Request.
func TestDeleteFilterNilUUIDIsBadRequest(t *testing.T) {
	h, provider, _ := newDeleteFilterFixture(t)

	rec := deleteFilter(t, h, provider, uuid.Nil.String())

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

// TestDeleteFilterNotFound verifies that attempting to delete a filter with a
// valid UUID that does not exist in the database returns HTTP 404 Not Found.
func TestDeleteFilterNotFound(t *testing.T) {
	h, provider, _ := newDeleteFilterFixture(t)

	randomUUID := uuid.Must(uuid.NewV4()).String()
	rec := deleteFilter(t, h, provider, randomUUID)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
}

// TestDeleteFilterRemoteProviderNotFound verifies that when a remote provider
// returns a 404 status, the handler propagates HTTP 404 Not Found rather than
// masking it as HTTP 500 Internal Server Error.
func TestDeleteFilterRemoteProviderNotFound(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	remoteErr := httputil.WithProviderStatus(fmt.Errorf("filter not found on remote provider"), http.StatusNotFound)
	mockProvider := &mockRemoteFilterProvider{err: remoteErr}

	randomUUID := uuid.Must(uuid.NewV4()).String()
	rec := deleteFilter(t, h, mockProvider, randomUUID)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
}

// TestDeleteFilterExisting verifies that deleting an existing filter succeeds
// with HTTP 200 OK and removes the filter record from the database.
func TestDeleteFilterExisting(t *testing.T) {
	h, provider, db := newDeleteFilterFixture(t)

	filterID := core.Uuid(uuid.Must(uuid.NewV4()))
	filter := models.MesheryFilter{
		ID:         &filterID,
		Name:       "test-filter",
		Visibility: models.Private,
	}
	if err := db.Create(&filter).Error; err != nil {
		t.Fatalf("seed filter: %v", err)
	}

	rec := deleteFilter(t, h, provider, filterID.String())

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var count int64
	if err := db.Model(&models.MesheryFilter{}).Where("id = ?", filterID).Count(&count).Error; err != nil {
		t.Fatalf("count filters: %v", err)
	}
	if count != 0 {
		t.Errorf("filter count = %d, want 0 after deletion", count)
	}
}
