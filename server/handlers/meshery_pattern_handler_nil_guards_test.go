package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func TestPatternPOST_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	req := httptest.NewRequest(http.MethodPost, "/api/pattern", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.handlePatternPOST(rec, req, nil, nil, provider)

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

func TestPatternPOST_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/pattern", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.handlePatternPOST(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestPatternPOST_NilSystemID_DoesNotPanic(t *testing.T) {
	fixture := newPatternPOSTFixture(t)
	fixture.h.SystemID = nil
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	body := `{"name":"demo","pattern_file":"name: demo\nversion: 0.0.1\n"}`
	req := httptest.NewRequest(http.MethodPost, "/api/pattern", strings.NewReader(body))
	rec := httptest.NewRecorder()

	fixture.h.handlePatternPOST(rec, req, nil, authUser, fixture.provider)

	if rec.Code != http.StatusOK && rec.Code != http.StatusBadRequest && rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200, 400, or 500, got %d", rec.Code)
	}
}

func TestGetMesheryPatternsHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)

	req := httptest.NewRequest(http.MethodGet, "/api/pattern", nil)
	rec := httptest.NewRecorder()

	h.GetMesheryPatternsHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetCatalogMesheryPatternsHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)

	req := httptest.NewRequest(http.MethodGet, "/api/pattern/catalog", nil)
	rec := httptest.NewRecorder()

	h.GetCatalogMesheryPatternsHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestDeleteMesheryPatternHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	patID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodDelete, "/api/pattern/"+patID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String()})
	rec := httptest.NewRecorder()

	h.DeleteMesheryPatternHandler(rec, req, nil, nil, provider)

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

func TestDeleteMesheryPatternHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	patID := uuid.Must(uuid.NewV4())
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodDelete, "/api/pattern/"+patID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String()})
	rec := httptest.NewRecorder()

	h.DeleteMesheryPatternHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestDownloadMesheryPatternHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	patID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodGet, "/api/pattern/download/"+patID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String()})
	rec := httptest.NewRecorder()

	h.DownloadMesheryPatternHandler(rec, req, nil, nil, provider)

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

func TestDownloadMesheryPatternHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	patID := uuid.Must(uuid.NewV4())
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodGet, "/api/pattern/download/"+patID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String()})
	rec := httptest.NewRecorder()

	h.DownloadMesheryPatternHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestPublishCatalogPatternHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	req := httptest.NewRequest(http.MethodPost, "/api/pattern/catalog/publish", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.PublishCatalogPatternHandler(rec, req, nil, nil, provider)

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

func TestPublishCatalogPatternHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/pattern/catalog/publish", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.PublishCatalogPatternHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestUnPublishCatalogPatternHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	req := httptest.NewRequest(http.MethodPost, "/api/pattern/catalog/unpublish", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.UnPublishCatalogPatternHandler(rec, req, nil, nil, provider)

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

func TestUnPublishCatalogPatternHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/pattern/catalog/unpublish", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.UnPublishCatalogPatternHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestDeleteMultiMesheryPatternsHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodDelete, "/api/pattern/delete", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.DeleteMultiMesheryPatternsHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetMesheryPatternHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	patID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodGet, "/api/pattern/"+patID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String()})
	rec := httptest.NewRecorder()

	h.GetMesheryPatternHandler(rec, req, nil, nil, provider)

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

func TestGetMesheryPatternHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	patID := uuid.Must(uuid.NewV4())
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodGet, "/api/pattern/"+patID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String()})
	rec := httptest.NewRecorder()

	h.GetMesheryPatternHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestHandlePatternUpdate_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	req := httptest.NewRequest(http.MethodPut, "/api/pattern", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.handlePatternUpdate(rec, req, nil, nil, provider)

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

func TestHandlePatternUpdate_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPut, "/api/pattern", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.handlePatternUpdate(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestHandlePatternUpdate_NilSystemID_DoesNotPanic(t *testing.T) {
	fixture := newPatternPOSTFixture(t)
	fixture.h.SystemID = nil
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPut, "/api/pattern", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	fixture.h.handlePatternUpdate(rec, req, nil, authUser, fixture.provider)

	// Token retrieval fails without token, but must not panic on nil SystemID
	if rec.Code != http.StatusOK && rec.Code != http.StatusBadRequest && rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200, 400, or 500, got %d", rec.Code)
	}
}

func TestDownloadMesheryPatternHandler_NilSystemID_DoesNotPanic(t *testing.T) {
	fixture := newPatternPOSTFixture(t)
	fixture.h.SystemID = nil
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	patID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodGet, "/api/pattern/download/"+patID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String()})
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("DownloadMesheryPatternHandler panicked with nil SystemID: %v", r)
		}
	}()

	fixture.h.DownloadMesheryPatternHandler(rec, req, nil, authUser, fixture.provider)

	if rec.Code == http.StatusBadRequest {
		t.Fatalf("expected the handler to reach the nil-SystemID path, got 400 Bad Request (route id not set?)")
	}
}

func TestGetMesheryPatternHandler_NilSystemID_DoesNotPanic(t *testing.T) {
	fixture := newPatternPOSTFixture(t)
	fixture.h.SystemID = nil
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	patID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodGet, "/api/pattern/"+patID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String()})
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("GetMesheryPatternHandler panicked with nil SystemID: %v", r)
		}
	}()

	fixture.h.GetMesheryPatternHandler(rec, req, nil, authUser, fixture.provider)

	if rec.Code == http.StatusBadRequest {
		t.Fatalf("expected the handler to reach the nil-SystemID path, got 400 Bad Request (route id not set?)")
	}
}

func TestCloneMesheryPatternHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	patID := uuid.Must(uuid.NewV4())
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/pattern/clone/"+patID.String(), strings.NewReader(`{"name":"cloned"}`))
	req = mux.SetURLVars(req, map[string]string{"id": patID.String()})
	rec := httptest.NewRecorder()

	h.CloneMesheryPatternHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetMesheryPatternSourceHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	patID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodGet, "/api/pattern/download/"+patID.String()+"/k8smanifest", nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String(), "sourcetype": "k8smanifest"})
	rec := httptest.NewRecorder()

	h.GetMesheryPatternSourceHandler(rec, req, nil, nil, provider)

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

func TestGetMesheryPatternSourceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	patID := uuid.Must(uuid.NewV4())
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodGet, "/api/pattern/download/"+patID.String()+"/k8smanifest", nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String(), "sourcetype": "k8smanifest"})
	rec := httptest.NewRecorder()

	h.GetMesheryPatternSourceHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestPublishCatalogPatternHandler_NullBody_ReturnsBadRequest(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "fake-token")
	req := httptest.NewRequest(http.MethodPost, "/api/pattern/catalog/publish", strings.NewReader("null")).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.PublishCatalogPatternHandler(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestUnPublishCatalogPatternHandler_NullBody_ReturnsBadRequest(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "fake-token")
	req := httptest.NewRequest(http.MethodDelete, "/api/pattern/catalog/unpublish", strings.NewReader("null")).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.UnPublishCatalogPatternHandler(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}
