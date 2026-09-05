package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
)

// newRegisterConnectionFixture wires ProcessConnectionRegistration the way the
// router does, over an in-memory database, mirroring
// newDeleteConnectionFixture in delete_connection_test.go.
func newRegisterConnectionFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider) {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}); err != nil {
		t.Fatalf("migrate tables: %v", err)
	}

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
		ConnectionToStateMachineInstanceTracker: &machines.ConnectionToStateMachineInstanceTracker{
			ConnectToInstanceMap: map[core.Uuid]*machines.StateMachine{},
		},
	}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister: &models.ConnectionPersister{DB: &db},
		EventsPersister:     &models.EventsPersister{DB: &db},
	}
	return h, provider
}

// TestProcessConnectionRegistration_PersistsConnectionBeforeDrivingStateMachine
// is the regression guard for #21813: registering a connection must persist
// its row before the state machine is driven, so a later
// GetConnectionByID-based status update can find it. Before the fix, no
// registration path ever called SaveConnection, so the connection was never
// persisted and the state machine's status-update lookup always failed with
// "record not found" — regardless of the connection kind or whether the
// state machine itself succeeded.
//
// This uses an arbitrary, unmapped connection kind (no registerActionForKind
// entry, e.g. not kubernetes/grafana/prometheus) so the transition into
// REGISTERED deterministically fails on a missing Action — a fast, in-memory
// failure with no network calls — while still exercising the exact code
// path (SaveConnection -> InitializeMachineWithContext -> SendEvent) that
// the fix changes.
func TestProcessConnectionRegistration_PersistsConnectionBeforeDrivingStateMachine(t *testing.T) {
	h, provider := newRegisterConnectionFixture(t)

	connID := uuid.Must(uuid.NewV4())
	coreID := core.Uuid(connID)
	body := connectionv1beta3.ConnectionRegistrationEvent{
		ID:                         &coreID,
		Kind:                       "artifacthub",
		Name:                       "test-connection",
		Status:                     connectionv1beta3.ConnectionRegistrationEventStatusRegister,
		SkipCredentialVerification: true,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal registration event: %v", err)
	}

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(bodyBytes)).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.ProcessConnectionRegistration(rec, req, nil, user, provider)

	// The specific symptom of #21813: the state machine's status-update
	// lookup failing because nothing was ever persisted. That failure must
	// not recur, however the rest of the (unmapped-kind) transition fares.
	if strings.Contains(rec.Body.String(), "record not found") {
		t.Fatalf("got the pre-fix 'record not found' failure; connection was not persisted before the state machine ran. body: %s", rec.Body.String())
	}

	if _, err := provider.ConnectionPersister.GetConnection(coreID, "artifacthub"); err != nil {
		t.Fatalf("expected the connection to be persisted before the state machine runs, but it was not found: %v", err)
	}
}

// TestProcessConnectionRegistration_PersistsWithDiscoveredStatus guards
// against a second defect in the same code path: registrationEventToConnectionPayload
// maps the wire event's Status field (a ConnectionRegistrationEventStatus like
// "register") directly into connections.ConnectionPayload.Status via a bare
// type cast. That value is a state-machine event name, not a valid
// connections.ConnectionStatus (discovered/registered/connected/disconnected).
// The persisted row must record DISCOVERED — the correct initial status for
// a freshly-registered connection — not the raw event string.
func TestProcessConnectionRegistration_PersistsWithDiscoveredStatus(t *testing.T) {
	h, provider := newRegisterConnectionFixture(t)

	connID := uuid.Must(uuid.NewV4())
	coreID := core.Uuid(connID)
	body := connectionv1beta3.ConnectionRegistrationEvent{
		ID:                         &coreID,
		Kind:                       "artifacthub",
		Name:                       "test-connection",
		Status:                     connectionv1beta3.ConnectionRegistrationEventStatusRegister,
		SkipCredentialVerification: true,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal registration event: %v", err)
	}

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(bodyBytes)).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.ProcessConnectionRegistration(rec, req, nil, user, provider)

	conn, err := provider.ConnectionPersister.GetConnection(coreID, "artifacthub")
	if err != nil {
		t.Fatalf("expected the connection to be persisted, but it was not found: %v", err)
	}
	if conn.Status != connections.DISCOVERED {
		t.Fatalf("expected persisted status to be %q, got %q (the raw wire event status must not leak into the connection's persisted status)", connections.DISCOVERED, conn.Status)
	}
}

// TestProcessConnectionRegistration_RejectsNilConnectionID guards against a
// registration event with a missing/nil connection ID being silently
// persisted under a freshly-minted ID while the state machine continues to
// operate on uuid.Nil — a divergence that reintroduces the #21813
// "record not found" symptom under a different cause. A nil ID must be
// rejected with 400 before any persistence or state-machine work happens.
func TestProcessConnectionRegistration_RejectsNilConnectionID(t *testing.T) {
	h, provider := newRegisterConnectionFixture(t)

	body := connectionv1beta3.ConnectionRegistrationEvent{
		ID:                         nil,
		Kind:                       "artifacthub",
		Name:                       "test-connection",
		Status:                     connectionv1beta3.ConnectionRegistrationEventStatusRegister,
		SkipCredentialVerification: true,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal registration event: %v", err)
	}

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", bytes.NewReader(bodyBytes)).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.ProcessConnectionRegistration(rec, req, nil, user, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for a nil connection ID, got %d. body: %s", rec.Code, rec.Body.String())
	}

	if _, err := provider.ConnectionPersister.GetConnection(core.Uuid(uuid.Nil), "artifacthub"); err == nil {
		t.Fatalf("expected no connection to be persisted under uuid.Nil, but one was found")
	}
}
