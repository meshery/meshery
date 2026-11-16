package handlers

import (
	"context"
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
	"github.com/meshery/meshkit/models/events"
)

func newFilterNilGuardsFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider) {
	t.Helper()

	db, err := database.New(database.Options{
		Engine:   database.SQLITE,
		Filename: filepath.Join(t.TempDir(), "meshery.db"),
	})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&models.MesheryFilter{}, &events.Event{}); err != nil {
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
		EventsPersister:        &models.EventsPersister{DB: &db},
		GenericPersister:       &db,
	}
	return h, provider
}

func TestGetMesheryFilterFileHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newFilterNilGuardsFixture(t)
	filterID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodGet, "/api/filter/file/"+filterID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": filterID})
	rec := httptest.NewRecorder()

	h.GetMesheryFilterFileHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestHandleFilterPOST_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newFilterNilGuardsFixture(t)

	req := httptest.NewRequest(http.MethodPost, "/api/filter", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.handleFilterPOST(rec, req, nil, nil, provider)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshaling error response: %v", err)
	}
	if body["error"] != "user unauthorized" {
		t.Fatalf("expected 'user unauthorized', got %q", body["error"])
	}
}

func TestHandleFilterPOST_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newFilterNilGuardsFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/filter", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.handleFilterPOST(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestHandleFilterPOST_NilSystemID_DoesNotPanic(t *testing.T) {
	h, provider := newFilterNilGuardsFixture(t)
	h.SystemID = nil
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	body := `{"filter_data":{"name":"test-filter","filter_file":"test"},"save":false}`
	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "fake-token")
	req := httptest.NewRequest(http.MethodPost, "/api/filter", strings.NewReader(body)).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.handleFilterPOST(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusOK && rec.Code != http.StatusBadRequest && rec.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", rec.Code)
	}
}

func TestGetMesheryFiltersHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newFilterNilGuardsFixture(t)

	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "fake-token")
	req := httptest.NewRequest(http.MethodGet, "/api/filter", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetMesheryFiltersHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetMesheryFiltersHandler_MissingToken_ReturnsInternalServerError(t *testing.T) {
	h, provider := newFilterNilGuardsFixture(t)

	req := httptest.NewRequest(http.MethodGet, "/api/filter", nil)
	rec := httptest.NewRecorder()

	h.GetMesheryFiltersHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for missing token, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrFetchTokenCode {
		t.Fatalf("expected %s, got %s", ErrFetchTokenCode, decoded.Code)
	}
}

func TestGetCatalogMesheryFiltersHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newFilterNilGuardsFixture(t)

	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "fake-token")
	req := httptest.NewRequest(http.MethodGet, "/api/filter/catalog", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetCatalogMesheryFiltersHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetCatalogMesheryFiltersHandler_MissingToken_ReturnsInternalServerError(t *testing.T) {
	h, provider := newFilterNilGuardsFixture(t)

	req := httptest.NewRequest(http.MethodGet, "/api/filter/catalog", nil)
	rec := httptest.NewRecorder()

	h.GetCatalogMesheryFiltersHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for missing token, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrFetchTokenCode {
		t.Fatalf("expected %s, got %s", ErrFetchTokenCode, decoded.Code)
	}
}

func TestDeleteMesheryFilterHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newFilterNilGuardsFixture(t)
	filterID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodDelete, "/api/filter/"+filterID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": filterID})
	rec := httptest.NewRecorder()

	h.DeleteMesheryFilterHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshaling error response: %v", err)
	}
	if body["error"] != "user unauthorized" {
		t.Fatalf("expected 'user unauthorized', got %q", body["error"])
	}
}

func TestDeleteMesheryFilterHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newFilterNilGuardsFixture(t)
	filterID := uuid.Must(uuid.NewV4()).String()
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodDelete, "/api/filter/"+filterID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": filterID})
	rec := httptest.NewRecorder()

	h.DeleteMesheryFilterHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestCloneMesheryFilterHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newFilterNilGuardsFixture(t)
	filterID := uuid.Must(uuid.NewV4()).String()

	body := `{"name":"cloned-filter"}`
	req := httptest.NewRequest(http.MethodPost, "/api/filter/clone/"+filterID, strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": filterID})
	rec := httptest.NewRecorder()

	h.CloneMesheryFilterHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	var res map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("unmarshaling error response: %v", err)
	}
	if res["error"] != "user unauthorized" {
		t.Fatalf("expected 'user unauthorized', got %q", res["error"])
	}
}

func TestCloneMesheryFilterHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newFilterNilGuardsFixture(t)
	filterID := uuid.Must(uuid.NewV4()).String()
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	body := `{"name":"cloned-filter"}`
	req := httptest.NewRequest(http.MethodPost, "/api/filter/clone/"+filterID, strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": filterID})
	rec := httptest.NewRecorder()

	h.CloneMesheryFilterHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestPublishCatalogFilterHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newFilterNilGuardsFixture(t)

	req := httptest.NewRequest(http.MethodPost, "/api/filter/catalog/publish", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.PublishCatalogFilterHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshaling error response: %v", err)
	}
	if body["error"] != "user unauthorized" {
		t.Fatalf("expected 'user unauthorized', got %q", body["error"])
	}
}

func TestPublishCatalogFilterHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newFilterNilGuardsFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/filter/catalog/publish", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.PublishCatalogFilterHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestPublishCatalogFilterHandler_NilSystemID_DoesNotPanic(t *testing.T) {
	h, provider := newFilterNilGuardsFixture(t)
	h.SystemID = nil
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/filter/catalog/publish", strings.NewReader(`invalid-json`))
	rec := httptest.NewRecorder()

	h.PublishCatalogFilterHandler(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusBadRequest && rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 400 or 500, got %d", rec.Code)
	}
}

func TestUnPublishCatalogFilterHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newFilterNilGuardsFixture(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/filter/catalog/unpublish", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.UnPublishCatalogFilterHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshaling error response: %v", err)
	}
	if body["error"] != "user unauthorized" {
		t.Fatalf("expected 'user unauthorized', got %q", body["error"])
	}
}

func TestUnPublishCatalogFilterHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newFilterNilGuardsFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodDelete, "/api/filter/catalog/unpublish", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.UnPublishCatalogFilterHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestUnPublishCatalogFilterHandler_NilSystemID_DoesNotPanic(t *testing.T) {
	h, provider := newFilterNilGuardsFixture(t)
	h.SystemID = nil
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodDelete, "/api/filter/catalog/unpublish", strings.NewReader(`invalid-json`))
	rec := httptest.NewRecorder()

	h.UnPublishCatalogFilterHandler(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusBadRequest && rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 400 or 500, got %d", rec.Code)
	}
}

func TestGetMesheryFilterHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newFilterNilGuardsFixture(t)
	filterID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodGet, "/api/filter/"+filterID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": filterID})
	rec := httptest.NewRecorder()

	h.GetMesheryFilterHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGenerateFilterComponent_NilRegistryManager_DoesNotPanic(t *testing.T) {
	h, _ := newFilterNilGuardsFixture(t)
	h.registryManager = nil

	res, err := h.generateFilterComponent("")
	if err != nil {
		t.Fatalf("expected nil error for empty config with nil registryManager, got %v", err)
	}
	if res != "" {
		t.Fatalf("expected empty string, got %q", res)
	}

	_, err = h.generateFilterComponent("some-config")
	if err == nil {
		t.Fatalf("expected error for non-empty config with nil registryManager")
	}
}
