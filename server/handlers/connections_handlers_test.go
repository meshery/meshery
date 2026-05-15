package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/schemas/models/core"
)

type connectionPingSpyProvider struct {
	*models.DefaultLocalProvider
	connection           *connections.Connection
	statusCode           int
	err                  error
	calls                int
	observedToken        string
	observedConnectionID core.Uuid
}

func newConnectionPingSpyProvider(connection *connections.Connection, statusCode int, err error) *connectionPingSpyProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &connectionPingSpyProvider{
		DefaultLocalProvider: base,
		connection:           connection,
		statusCode:           statusCode,
		err:                  err,
	}
}

func (m *connectionPingSpyProvider) GetConnectionByID(token string, connectionID core.Uuid) (*connections.Connection, int, error) {
	m.calls++
	m.observedToken = token
	m.observedConnectionID = connectionID
	return m.connection, m.statusCode, m.err
}

func TestConnectionPingHandler_InvalidConnectionIdReturnsBadRequest(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newConnectionPingSpyProvider(nil, http.StatusOK, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections/not-a-uuid/ping", nil)
	req = mux.SetURLVars(req, map[string]string{"connectionId": "not-a-uuid"})
	rec := httptest.NewRecorder()

	h.ConnectionPingHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	if provider.calls != 0 {
		t.Fatalf("provider GetConnectionByID called %d times, want 0", provider.calls)
	}
}

func TestConnectionPingHandler_DefaultKindReturnsStoredStatus(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	connectionID := uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))
	provider := newConnectionPingSpyProvider(&connections.Connection{
		ID:     connectionID,
		Name:   "Acme Meshery",
		Kind:   "meshery",
		Status: connections.CONNECTED,
	}, http.StatusOK, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections/"+connectionID.String()+"/ping", nil)
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connectionID.String()})
	rec := httptest.NewRecorder()

	h.ConnectionPingHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	if provider.calls != 1 {
		t.Fatalf("provider GetConnectionByID called %d times, want 1", provider.calls)
	}
	if provider.observedToken != "test-token" {
		t.Fatalf("provider token=%q, want %q", provider.observedToken, "test-token")
	}
	if provider.observedConnectionID != connectionID {
		t.Fatalf("provider connectionID=%v, want %v", provider.observedConnectionID, connectionID)
	}

	var got map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("response did not decode as JSON: %v", err)
	}
	want := map[string]string{
		"name":   "Acme Meshery",
		"kind":   "meshery",
		"status": string(connections.CONNECTED),
	}
	for key, wantValue := range want {
		if got[key] != wantValue {
			t.Fatalf("response[%q]=%q, want %q", key, got[key], wantValue)
		}
	}
}

func TestConnectionPingHandler_GrafanaWithoutCredentialReturnsBadRequest(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	connectionID := uuid.Must(uuid.FromString("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"))
	provider := newConnectionPingSpyProvider(&connections.Connection{
		ID:       connectionID,
		Name:     "Acme Grafana",
		Kind:     "grafana",
		Metadata: map[string]interface{}{"url": "https://grafana.example.com"},
		Status:   connections.CONNECTED,
	}, http.StatusOK, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections/"+connectionID.String()+"/ping", nil)
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connectionID.String()})
	rec := httptest.NewRecorder()

	h.ConnectionPingHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d (body=%q)", rec.Code, rec.Body.String())
	}
}
