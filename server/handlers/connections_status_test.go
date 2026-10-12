package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

func newConnectionStatusFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider, *database.Handler) {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "test.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}, models.K8sContext{}, models.Credential{}); err != nil {
		t.Fatalf("migrate tables: %v", err)
	}

	systemID := uuid.Must(uuid.NewV4())
	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     newTestLogger(t),
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: tracker,
	}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister:        &models.ConnectionPersister{DB: &db},
		EventsPersister:            &models.EventsPersister{DB: &db},
		MesheryK8sContextPersister: &models.MesheryK8sContextPersister{DB: &db},
		GenericPersister:           &db,
	}
	return h, provider, &db
}

func updateConnectionStatusReq(t *testing.T, h *Handler, provider models.Provider, connectionID uuid.UUID, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connectionID.String(), strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connectionID.String()})
	rec := httptest.NewRecorder()

	h.UpdateConnectionById(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	return rec
}

// TestUpdateConnectionStatus_NonKubernetes_GitHub verifies that status updates for
// non-Kubernetes connections (e.g. GitHub) do not attempt to look up a Kubernetes context,
// do not emit false-positive Error events, correctly update persisted status, and
// initialize the state machine with the connection's actual kind.
func TestUpdateConnectionStatus_NonKubernetes_GitHub(t *testing.T) {
	h, provider, db := newConnectionStatusFixture(t)

	connID := uuid.Must(uuid.NewV4())
	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		ID:             connID,
		Name:           "my-github-connection",
		Kind:           "github",
		ConnectionType: "scm",
		Status:         connections.DISCOVERED,
	})
	if err != nil {
		t.Fatalf("save initial connection: %v", err)
	}

	rec := updateConnectionStatusReq(t, h, provider, saved.ID, `{"status":"connected"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	// Poll until the state-machine transition completes and persists the confirmation event
	var infoEvents []events.Event
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		_ = db.Where("acted_upon = ? AND description LIKE ?", saved.ID, "%connection changed to connected%").Find(&infoEvents).Error
		if len(infoEvents) > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if len(infoEvents) == 0 {
		t.Fatalf("timed out waiting for connection %s to transition to connected", saved.ID)
	}

	updated, err := provider.ConnectionPersister.GetConnection(saved.ID, "")
	if err != nil {
		t.Fatalf("get updated connection: %v", err)
	}
	if updated.Status != connections.CONNECTED {
		t.Errorf("connection status = %q, want %q", updated.Status, connections.CONNECTED)
	}

	// Verify that no Error events were persisted
	var errorEvents []events.Event
	if err := db.Where("severity = ?", events.Error).Find(&errorEvents).Error; err != nil {
		t.Fatalf("query error events: %v", err)
	}
	if len(errorEvents) != 0 {
		t.Errorf("expected 0 error events for non-k8s update, got %d: %+v", len(errorEvents), errorEvents[0])
	}

	// Verify that the state machine was tracked with the correct kind and settled into CONNECTED
	inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(saved.ID)
	if !ok {
		t.Fatalf("expected state machine to be registered in tracker for connection %s", saved.ID)
	}
	if inst.Name != "github" {
		t.Errorf("state machine name = %q, want %q", inst.Name, "github")
	}
	if inst.CurrentState != machines.CONNECTED {
		t.Errorf("state machine current state = %q, want %q", inst.CurrentState, machines.CONNECTED)
	}
}

// TestUpdateConnectionStatus_NonKubernetes_Prometheus verifies that a Prometheus connection
// can transition from connected to disconnected without invoking Kubernetes-specific context logic.
func TestUpdateConnectionStatus_NonKubernetes_Prometheus(t *testing.T) {
	h, provider, db := newConnectionStatusFixture(t)

	connID := uuid.Must(uuid.NewV4())
	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		ID:             connID,
		Name:           "my-prometheus",
		Kind:           "prometheus",
		ConnectionType: "observability",
		Status:         connections.CONNECTED,
	})
	if err != nil {
		t.Fatalf("save initial connection: %v", err)
	}

	rec := updateConnectionStatusReq(t, h, provider, saved.ID, `{"status":"disconnected"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	// Poll until the state-machine transition completes and persists the confirmation event
	var promInfoEvents []events.Event
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		_ = db.Where("acted_upon = ? AND description LIKE ?", saved.ID, "%connection changed to disconnected%").Find(&promInfoEvents).Error
		if len(promInfoEvents) > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if len(promInfoEvents) == 0 {
		t.Fatalf("timed out waiting for prometheus connection %s to transition to disconnected", saved.ID)
	}

	updated, err := provider.ConnectionPersister.GetConnection(saved.ID, "")
	if err != nil {
		t.Fatalf("get updated connection: %v", err)
	}
	if updated.Status != connections.DISCONNECTED {
		t.Errorf("connection status = %q, want %q", updated.Status, connections.DISCONNECTED)
	}

	var errorEvents []events.Event
	if err := db.Where("severity = ?", events.Error).Find(&errorEvents).Error; err != nil {
		t.Fatalf("query error events: %v", err)
	}
	if len(errorEvents) != 0 {
		t.Errorf("expected 0 error events for prometheus update, got %d", len(errorEvents))
	}

	inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(saved.ID)
	if !ok {
		t.Fatalf("expected state machine to be registered in tracker for connection %s", saved.ID)
	}
	if inst.Name != "prometheus" {
		t.Errorf("state machine name = %q, want %q", inst.Name, "prometheus")
	}
	if inst.CurrentState != machines.DISCONNECTED {
		t.Errorf("state machine current state = %q, want %q", inst.CurrentState, machines.DISCONNECTED)
	}
}

// TestUpdateConnectionStatus_Kubernetes_PreservesExistingBehavior verifies that Kubernetes
// connections continue to follow the Kubernetes-specific path and fail with an Error event
// when the underlying Kubernetes context is missing.
func TestUpdateConnectionStatus_Kubernetes_PreservesExistingBehavior(t *testing.T) {
	h, provider, db := newConnectionStatusFixture(t)

	connID := uuid.Must(uuid.NewV4())
	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		ID:             connID,
		Name:           "my-k8s-cluster",
		Kind:           "kubernetes",
		ConnectionType: "cluster",
		Status:         connections.DISCOVERED,
	})
	if err != nil {
		t.Fatalf("save initial connection: %v", err)
	}

	// Updating a Kubernetes connection without a backing k8s context record in meshery_k8s_contexts
	// must invoke GetK8sContext, which errors out and persists an Error event.
	rec := updateConnectionStatusReq(t, h, provider, saved.ID, `{"status":"connected"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	// Poll until the expected Error event is persisted
	var errorEvents []events.Event
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		_ = db.Where("severity = ?", events.Error).Find(&errorEvents).Error
		if len(errorEvents) > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if len(errorEvents) == 0 {
		t.Fatalf("timed out waiting for expected error event on missing kubernetes context")
	}
	if !strings.Contains(errorEvents[0].Description, "Failed to update connection status for") {
		t.Errorf("unexpected error description: %s", errorEvents[0].Description)
	}
}

// TestNotifySmOfConnectionStatusChange_NilConnection ensures calling NotifySmOfConnectionStatusChange
// with a nil connection does not panic and returns safely.
func TestNotifySmOfConnectionStatusChange_NilConnection(t *testing.T) {
	h, provider, _ := newConnectionStatusFixture(t)
	systemID := *h.SystemID
	userID := uuid.Must(uuid.NewV4())

	ev, err := h.NotifySmOfConnectionStatusChange(context.Background(), userID, provider, "test-token", nil)
	if err != nil {
		t.Fatalf("expected nil error on nil connection, got: %v", err)
	}
	if ev.SystemID != systemID {
		t.Errorf("event systemID = %v, want %v", ev.SystemID, systemID)
	}
}

// TestNotifySmOfConnectionStatusChange_UnspecifiedKind verifies that updating a connection
// with an empty/unresolvable kind fails gracefully with an informative error event.
func TestNotifySmOfConnectionStatusChange_UnspecifiedKind(t *testing.T) {
	h, provider, _ := newConnectionStatusFixture(t)
	userID := uuid.Must(uuid.NewV4())
	connID := uuid.Must(uuid.NewV4())

	payload := &connections.ConnectionPayload{
		ID:     connID,
		Status: connections.CONNECTED,
		Kind:   "", // unspecified and not in DB
	}

	ev, err := h.NotifySmOfConnectionStatusChange(context.Background(), userID, provider, "test-token", payload)
	if err == nil {
		t.Fatalf("expected error for unspecified connection kind, got nil")
	}
	if ev.Severity != events.Error {
		t.Errorf("expected event severity %v, got %v", events.Error, ev.Severity)
	}
	if !strings.Contains(ev.Description, "connection kind is unspecified") {
		t.Errorf("unexpected event description: %s", ev.Description)
	}
}

// TestUpdateConnectionStatus_InvalidUUID ensures invalid connection UUIDs return an error.
func TestUpdateConnectionStatus_InvalidUUID(t *testing.T) {
	h, provider, _ := newConnectionStatusFixture(t)

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/invalid-uuid", strings.NewReader(`{"status":"connected"}`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": "invalid-uuid"})
	rec := httptest.NewRecorder()

	h.UpdateConnectionById(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "meshery-server-1051") {
		t.Errorf("expected ErrFailToSave code in body, got: %s", rec.Body.String())
	}
}

// TestUpdateConnectionStatus_MalformedJSON ensures malformed JSON body returns 400 Bad Request.
func TestUpdateConnectionStatus_MalformedJSON(t *testing.T) {
	h, provider, _ := newConnectionStatusFixture(t)
	connID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), strings.NewReader(`{malformed`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()

	h.UpdateConnectionById(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}
