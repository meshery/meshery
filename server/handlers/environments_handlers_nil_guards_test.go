package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func newEnvironmentNilGuardsFixture(t *testing.T) *Handler {
	t.Helper()
	return &Handler{
		log: newTestLogger(t),
	}
}

func TestGetEnvironments_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newEnvironmentNilGuardsFixture(t)

	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "fake-token")
	req := httptest.NewRequest(http.MethodGet, "/api/environments?orgId=11111111-1111-1111-1111-111111111111", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetEnvironments(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetEnvironmentByIDHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newEnvironmentNilGuardsFixture(t)
	envID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodGet, "/api/environments/"+envID+"?orgId=11111111-1111-1111-1111-111111111111", nil)
	req = mux.SetURLVars(req, map[string]string{"id": envID})
	rec := httptest.NewRecorder()

	h.GetEnvironmentByIDHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestSaveEnvironment_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newEnvironmentNilGuardsFixture(t)

	body := `{"name":"test-env","organizationId":"11111111-1111-1111-1111-111111111111"}`
	req := httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.SaveEnvironment(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestDeleteEnvironmentHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newEnvironmentNilGuardsFixture(t)
	envID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodDelete, "/api/environments/"+envID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": envID})
	rec := httptest.NewRecorder()

	h.DeleteEnvironmentHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestUpdateEnvironmentHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newEnvironmentNilGuardsFixture(t)
	envID := uuid.Must(uuid.NewV4()).String()

	body := `{"name":"updated-env"}`
	req := httptest.NewRequest(http.MethodPut, "/api/environments/"+envID, strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": envID})
	rec := httptest.NewRecorder()

	h.UpdateEnvironmentHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestAddConnectionToEnvironmentHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newEnvironmentNilGuardsFixture(t)
	envID := uuid.Must(uuid.NewV4()).String()
	connID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodPost, "/api/environments/"+envID+"/connections/"+connID, nil)
	req = mux.SetURLVars(req, map[string]string{"environmentID": envID, "connectionID": connID})
	rec := httptest.NewRecorder()

	h.AddConnectionToEnvironmentHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestRemoveConnectionFromEnvironmentHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newEnvironmentNilGuardsFixture(t)
	envID := uuid.Must(uuid.NewV4()).String()
	connID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodDelete, "/api/environments/"+envID+"/connections/"+connID, nil)
	req = mux.SetURLVars(req, map[string]string{"environmentID": envID, "connectionID": connID})
	rec := httptest.NewRecorder()

	h.RemoveConnectionFromEnvironmentHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetConnectionsOfEnvironmentHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newEnvironmentNilGuardsFixture(t)
	envID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodGet, "/api/environments/"+envID+"/connections", nil)
	req = mux.SetURLVars(req, map[string]string{"environmentID": envID})
	rec := httptest.NewRecorder()

	h.GetConnectionsOfEnvironmentHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestEnvironmentHandler_TypedNilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newEnvironmentNilGuardsFixture(t)
	var rp *models.RemoteProvider = nil
	var typedNil models.Provider = rp

	req := httptest.NewRequest(http.MethodGet, "/api/environments", nil)
	rec := httptest.NewRecorder()

	h.GetEnvironments(rec, req, nil, nil, typedNil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for typed-nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestEnvironmentHandler_NilLogger_DoesNotPanic(t *testing.T) {
	h := &Handler{log: nil}
	req := httptest.NewRequest(http.MethodGet, "/api/environments", nil)
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("GetEnvironments panicked with nil logger: %v", r)
		}
	}()

	h.GetEnvironments(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider with nil logger, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}
