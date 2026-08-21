package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/helpers"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// registerRouteHarness drives POST /api/integrations/connections/register all
// the way through the real state machine into machines.DefaultConnectAction and
// back out as an HTTP response, so a test asserts the status a client actually
// receives rather than the shape of an error value.
type registerRouteHarness struct {
	handler  *Handler
	provider *registerRouteProvider
	connID   core.Uuid
	userID   core.Uuid
}

// registerRouteProvider is a DefaultLocalProvider backed by an in-memory
// datastore. Only the credentials table is migrated: the credential write is
// the step under test and it fails before any connection row is written.
type registerRouteProvider struct {
	*models.DefaultLocalProvider
}

func (p *registerRouteProvider) PersistEvent(_ events.Event, _ string) error { return nil }

func newRegisterRouteHarness(t *testing.T) *registerRouteHarness {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("opening in-memory database: %v", err)
	}
	if err := db.AutoMigrate(&models.Credential{}); err != nil {
		t.Fatalf("migrating credential schema: %v", err)
	}

	base := &models.DefaultLocalProvider{GenericPersister: &database.Handler{DB: db}}
	provider := &registerRouteProvider{DefaultLocalProvider: base}

	systemID := uuid.Must(uuid.NewV4())
	connID := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())

	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: map[core.Uuid]*machines.StateMachine{},
	}

	h := newTestHandler(t, map[string]models.Provider{}, "")
	h.SystemID = &systemID
	h.ConnectionToStateMachineInstanceTracker = tracker
	h.config.EventBroadcaster = models.NewBroadcaster("test")

	// Seed the tracker with a machine already in REGISTERED. The handler always
	// initialises at DISCOVERED, from which `connect` has no edge; a real
	// registration reaches DefaultConnectAction only from REGISTERED, and
	// InitializeMachineWithContext reuses a tracked machine, so this is the
	// state a live server would be in when the connect event arrives.
	if _, err := helpers.InitializeMachineWithContext(
		map[string]string{}, context.Background(), connID, userID, tracker,
		newTestLogger(t), provider, machines.REGISTERED, "prometheus", nil,
	); err != nil {
		t.Fatalf("seeding the connection state machine: %v", err)
	}

	return &registerRouteHarness{handler: h, provider: provider, connID: connID, userID: userID}
}

// connect posts a `connected` registration event carrying credentialSecret and
// returns the recorded response.
func (h *registerRouteHarness) connect(t *testing.T, credentialSecret map[string]interface{}) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(map[string]interface{}{
		"id":               h.connID.String(),
		"kind":             "prometheus",
		"type":             "observability",
		"subType":          "monitoring",
		"name":             "prom",
		"status":           string(machines.Connect),
		"credentialSecret": credentialSecret,
	})
	if err != nil {
		t.Fatalf("marshalling the registration event: %v", err)
	}

	user := &models.User{ID: h.userID}
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", strings.NewReader(string(body)))
	ctx := context.WithValue(req.Context(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.handler.SystemID)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	h.handler.ProcessConnectionRegistration(rec, req, nil, user, h.provider)
	return rec
}

// TestConnectionRegisterRejectsReservedCredentialProperty is the positive half
// of the status contract: the same request body the credential endpoints answer
// with 400 must also be a 400 here. The connection-registration route is
// equally client-controlled - credentialSecret is an arbitrary client map - and
// before the error chain carried the tagged status it answered 500, blaming
// Meshery for a body the caller can fix.
func TestConnectionRegisterRejectsReservedCredentialProperty(t *testing.T) {
	h := newRegisterRouteHarness(t)

	rec := h.connect(t, map[string]interface{}{
		"__mesheryEncryptedSecret": "meshery.enc.v1:aa:bb",
		"apiKey":                   "real-token",
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body=%q)", rec.Code, rec.Body.String())
	}

	var resp struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decoding the error response: %v (body=%q)", err, rec.Body.String())
	}
	if resp.Code != ErrSendMachineEventCode {
		t.Errorf("error code = %q, want %q - the outer MeshKit error must still be the one rendered", resp.Code, ErrSendMachineEventCode)
	}

	var rows int64
	if err := h.provider.GetGenericPersister().Table("credentials").Count(&rows).Error; err != nil {
		t.Fatalf("counting credential rows: %v", err)
	}
	if rows != 0 {
		t.Errorf("%d credential rows persisted, want 0 - the plaintext apiKey reached the datastore", rows)
	}
}

// TestConnectionRegisterUntaggedCredentialFailuresStayInternal is the blast
// radius proof. Joining the cause into ErrPersistCredential and
// ErrSendMachineEvent makes every error reaching the handler introspectable,
// so the fallback must still hold for errors that carry no provider status.
// Both cases below are the other two ErrPersistCredential call sites in
// machines/actions.go, neither of which tags a status; only an error that
// explicitly carries one may move off 500.
func TestConnectionRegisterUntaggedCredentialFailuresStayInternal(t *testing.T) {
	for name, secret := range map[string]map[string]interface{}{
		"credential id is not a string": {"id": 42},
		"credential id is not a uuid":   {"id": "not-a-uuid"},
	} {
		t.Run(name, func(t *testing.T) {
			h := newRegisterRouteHarness(t)

			rec := h.connect(t, secret)

			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, want 500 - an untagged failure must stay on the fallback (body=%q)", rec.Code, rec.Body.String())
			}
		})
	}
}
