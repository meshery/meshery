package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func newConnectionActionTestHandler(t *testing.T) *Handler {
	t.Helper()
	return &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}
}

func TestPerformConnectionAction_InvalidConnectionID(t *testing.T) {
	h := newConnectionActionTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections//actions",
		strings.NewReader(`{"action":"setMeshsyncMode","mode":"operator"}`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": ""})
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for missing connection id; body: %s", rec.Code, rec.Body.String())
	}
}

func TestPerformConnectionAction_UnsupportedAction(t *testing.T) {
	h := newConnectionActionTestHandler(t)

	connID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/"+connID.String()+"/actions",
		strings.NewReader(`{"action":"noSuchAction"}`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for unsupported action; body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "unsupported connection action") {
		t.Fatalf("expected unsupported-action message, got: %s", rec.Body.String())
	}
}

func TestPerformConnectionAction_MalformedBody(t *testing.T) {
	h := newConnectionActionTestHandler(t)

	connID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/"+connID.String()+"/actions",
		strings.NewReader(`{not json`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for malformed body; body: %s", rec.Code, rec.Body.String())
	}
}
