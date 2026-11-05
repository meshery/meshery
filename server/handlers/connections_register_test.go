package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
)

// The Create Connection wizard sends {"status":"register"} for every
// non-Kubernetes kind. Kinds without a verification step (github,
// artifacthub) must advance to REGISTERED instead of failing the transition.
//
// The row is saved up front to cover a connection that already exists (e.g.
// re-registration): register must keep it and must not create a second one.
func newRegistrationTestHarness(t *testing.T) (*Handler, *models.DefaultLocalProvider, *database.Handler) {
	t.Helper()
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "test.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}, models.Credential{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
		ConnectionToStateMachineInstanceTracker: &machines.ConnectionToStateMachineInstanceTracker{
			ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
		},
	}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister: &models.ConnectionPersister{DB: &db},
		EventsPersister:     &models.EventsPersister{DB: &db},
		GenericPersister:    &db,
	}
	return h, provider, &db
}

// sendRegistrationEvent posts one wizard event to the registration endpoint.
func sendRegistrationEvent(t *testing.T, h *Handler, provider models.Provider, user *models.User, body map[string]interface{}) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal registration event: %v", err)
	}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(raw)).WithContext(ctx)
	rec := httptest.NewRecorder()
	h.ProcessConnectionRegistration(rec, req, nil, user, provider)
	return rec
}

func TestProcessConnectionRegistrationRegistersKindsWithoutVerifyAction(t *testing.T) {
	h, provider, _ := newRegistrationTestHarness(t)

	for _, kind := range []string{"github", "artifacthub"} {
		t.Run(kind, func(t *testing.T) {
			id := core.Uuid(uuid.Must(uuid.NewV4()))
			if _, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
				ID: id, Name: "my-" + kind, Kind: kind, Status: connections.DISCOVERED,
			}); err != nil {
				t.Fatalf("save connection: %v", err)
			}

			body, err := json.Marshal(connectionv1beta3.ConnectionRegistrationEvent{
				ID:                         &id,
				Kind:                       kind,
				Name:                       "my-" + kind,
				Status:                     connectionv1beta3.ConnectionRegistrationEventStatusRegister,
				SkipCredentialVerification: true,
			})
			if err != nil {
				t.Fatalf("marshal registration event: %v", err)
			}
			user := &models.User{ID: uuid.Must(uuid.NewV4())}
			ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
			ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
			req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(body)).WithContext(ctx)
			rec := httptest.NewRecorder()

			h.ProcessConnectionRegistration(rec, req, nil, user, provider)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
			}
			conn, err := provider.ConnectionPersister.GetConnection(id, kind)
			if err != nil {
				t.Fatalf("get connection: %v", err)
			}
			if conn.Status != connections.REGISTERED {
				t.Fatalf("persisted status = %q, want %q", conn.Status, connections.REGISTERED)
			}
		})
	}
}

// The wizard flow from #20983: register then connect for a connection with no
// pre-existing row. Register must persist the connection (and the new
// credential) as REGISTERED; connect must move that same row to CONNECTED
// without saving the credential a second time.
func TestProcessConnectionRegistrationPersistsOnRegister(t *testing.T) {
	h, provider, db := newRegistrationTestHarness(t)

	for _, kind := range []string{"prometheus", "grafana", "github"} {
		t.Run(kind, func(t *testing.T) {
			id := core.Uuid(uuid.Must(uuid.NewV4()))
			user := &models.User{ID: uuid.Must(uuid.NewV4())}
			event := map[string]interface{}{
				"id":   id.String(),
				"kind": kind,
				"name": "my-" + kind,
				// Unroutable on purpose: bypassed verification must not probe it.
				"metadata":                   map[string]interface{}{"url": "http://127.0.0.1:1"},
				"credentialSecret":           map[string]interface{}{"name": "my-" + kind + "-cred"},
				"skipCredentialVerification": true,
			}

			event["status"] = "register"
			if rec := sendRegistrationEvent(t, h, provider, user, event); rec.Code != http.StatusOK {
				t.Fatalf("register: status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
			}
			conn, err := provider.ConnectionPersister.GetConnection(id, kind)
			if err != nil {
				t.Fatalf("register: connection not persisted: %v", err)
			}
			if conn.Status != connections.REGISTERED {
				t.Fatalf("register: persisted status = %q, want %q", conn.Status, connections.REGISTERED)
			}
			if conn.CredentialID == nil || uuid.UUID(*conn.CredentialID) == uuid.Nil {
				t.Fatalf("register: connection has no credential attached")
			}
			credentialID := *conn.CredentialID

			event["status"] = "connect"
			if rec := sendRegistrationEvent(t, h, provider, user, event); rec.Code != http.StatusOK {
				t.Fatalf("connect: status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
			}
			conn, err = provider.ConnectionPersister.GetConnection(id, kind)
			if err != nil {
				t.Fatalf("connect: get connection: %v", err)
			}
			if conn.Status != connections.CONNECTED {
				t.Fatalf("connect: persisted status = %q, want %q", conn.Status, connections.CONNECTED)
			}
			if conn.CredentialID == nil || *conn.CredentialID != credentialID {
				t.Fatalf("connect: credential changed from %s to %v", credentialID, conn.CredentialID)
			}

			var credentials int64
			if err := db.Model(&models.Credential{}).Where("type = ?", kind).Count(&credentials).Error; err != nil {
				t.Fatalf("count credentials: %v", err)
			}
			if credentials != 1 {
				t.Fatalf("credentials saved = %d, want 1", credentials)
			}
		})
	}
}

// Cancelling a registration abandoned between register and connect removes the
// REGISTERED row; cancelling after connect keeps the connection.
func TestCancelConnectionRegisterDeletesOnlyRegisteredRows(t *testing.T) {
	h, provider, _ := newRegistrationTestHarness(t)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	cancel := func(id core.Uuid) {
		t.Helper()
		req := httptest.NewRequest(http.MethodDelete, "/api/integrations/connections/register/"+id.String(), nil)
		req = mux.SetURLVars(req, map[string]string{"registrationId": id.String()})
		rec := httptest.NewRecorder()
		h.CancelConnectionRegister(rec, req, nil, user, provider)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("cancel: status = %d, want %d. body: %s", rec.Code, http.StatusNoContent, rec.Body.String())
		}
	}

	for _, tt := range []struct {
		status     connections.ConnectionStatus
		wantExists bool
	}{
		{status: connections.REGISTERED, wantExists: false},
		{status: connections.CONNECTED, wantExists: true},
	} {
		t.Run(string(tt.status), func(t *testing.T) {
			id := core.Uuid(uuid.Must(uuid.NewV4()))
			if _, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
				ID: id, Name: "my-github", Kind: "github", Status: tt.status,
			}); err != nil {
				t.Fatalf("save connection: %v", err)
			}

			cancel(id)

			_, err := provider.ConnectionPersister.GetConnection(id, "")
			if exists := err == nil; exists != tt.wantExists {
				t.Fatalf("after cancel: row exists = %v, want %v (err: %v)", exists, tt.wantExists, err)
			}
		})
	}

	// Unknown ids stay a no-op.
	cancel(core.Uuid(uuid.Must(uuid.NewV4())))
}
