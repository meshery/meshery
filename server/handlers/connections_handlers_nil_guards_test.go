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
	"github.com/meshery/meshery/server/models/connections"
)

func TestProcessConnectionRegistration_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.ProcessConnectionRegistration(rec, req, nil, nil, provider)

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

func TestProcessConnectionRegistration_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.ProcessConnectionRegistration(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestSaveConnection_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.SaveConnection(rec, req, nil, nil, provider)

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

func TestSaveConnection_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.SaveConnection(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestSaveConnection_NilSystemID_DoesNotPanic(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	h.SystemID = nil
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	body := `{"name":"test-conn","kind":"kubernetes","type":"platform"}`
	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "fake-token")
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections", strings.NewReader(body)).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.SaveConnection(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusCreated && rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 201 Created or 500 InternalServerError, got %d; body: %s", rec.Code, rec.Body.String())
	}
}

func TestGetConnections_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections", nil)
	rec := httptest.NewRecorder()

	h.GetConnections(rec, req, nil, nil, provider)

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

func TestGetConnections_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections", nil)
	rec := httptest.NewRecorder()

	h.GetConnections(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetConnectionsByKind_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections/kubernetes", nil)
	req = mux.SetURLVars(req, map[string]string{"connectionKind": "kubernetes"})
	rec := httptest.NewRecorder()

	h.GetConnectionsByKind(rec, req, nil, nil, provider)

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

func TestGetConnectionsByKind_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections/kubernetes", nil)
	req = mux.SetURLVars(req, map[string]string{"connectionKind": "kubernetes"})
	rec := httptest.NewRecorder()

	h.GetConnectionsByKind(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetConnectionByID_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	connID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections/"+connID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()

	h.GetConnectionByID(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestUpdateConnectionById_NilUser_ReturnsUnauthorized(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	connID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), strings.NewReader(`{}`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()

	h.UpdateConnectionById(rec, req, nil, nil, provider)

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

func TestUpdateConnectionById_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h, _ := newDeleteConnectionFixture(t)
	connID := uuid.Must(uuid.NewV4())
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), strings.NewReader(`{}`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()

	h.UpdateConnectionById(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestUpdateConnectionById_NilSystemID_DoesNotPanic(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	h.SystemID = nil
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	// Seed a connection so the ID lookup finds it
	connID := uuid.Must(uuid.NewV4())
	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		ID:             connID,
		Name:           "existing-conn",
		Kind:           "kubernetes",
		ConnectionType: "platform",
		Status:         connections.DISCOVERED,
	})
	if err != nil {
		t.Fatalf("seed connection: %v", err)
	}

	body := `{"name":"renamed-conn","kind":"kubernetes","type":"platform"}`
	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "fake-token")
	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+saved.ID.String(), strings.NewReader(body)).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"connectionId": saved.ID.String()})
	rec := httptest.NewRecorder()

	h.UpdateConnectionById(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200 OK or 500 InternalServerError, got %d; body: %s", rec.Code, rec.Body.String())
	}
}
