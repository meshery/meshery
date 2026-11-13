package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	mhelpers "github.com/meshery/meshery/server/machines/helpers"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type mockContextsProvider struct {
	*models.DefaultLocalProvider
	getK8sContextFn func(token, id string) (models.K8sContext, error)
	persistEventFn  func(event events.Event, token string) error
	persistedEvents []events.Event
}

// GetK8sContext returns the stubbed context lookup, or an empty context when none is set.
func (m *mockContextsProvider) GetK8sContext(token, id string) (models.K8sContext, error) {
	if m.getK8sContextFn != nil {
		return m.getK8sContextFn(token, id)
	}
	return models.K8sContext{}, nil
}

// PersistEvent records the event and returns the stubbed persistence result.
func (m *mockContextsProvider) PersistEvent(event events.Event, token string) (models.PersistedEvent, error) {
	m.persistedEvents = append(m.persistedEvents, event)
	if m.persistEventFn != nil {
		if err := m.persistEventFn(event, token); err != nil {
			return models.PersistedEvent{}, err
		}
	}
	return models.PersistedEvent{Event: &event, Disposition: models.EventCreated}, nil
}

func (m *mockContextsProvider) GetConnectionByID(_ string, connectionID core.Uuid) (*connections.Connection, int, error) {
	return &connections.Connection{ID: connectionID, Kind: "kubernetes", Status: connections.DISCOVERED}, http.StatusOK, nil
}

func (m *mockContextsProvider) UpdateConnectionById(_ string, payload *connections.ConnectionPayload, _ string) (*connections.Connection, error) {
	return &connections.Connection{ID: payload.ID, Kind: payload.Kind, Status: payload.Status}, nil
}

// TestDeleteContext_MissingToken verifies that a request without a token is rejected with ErrFetchToken.
func TestDeleteContext_MissingToken(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	provider := &mockContextsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
	}
	contextID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+contextID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": contextID.String()})
	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), ErrFetchTokenCode) {
		t.Fatalf("expected error code %s in body, got %s", ErrFetchTokenCode, rec.Body.String())
	}
}

// TestDeleteContext_GetK8sContextError verifies that a failed context lookup returns a 500,
// persists one error event named after the context id, and broadcasts it.
func TestDeleteContext_GetK8sContextError(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	broadcaster := &models.Broadcast{}
	userID := uuid.Must(uuid.NewV4())
	ch, unsub := broadcaster.Subscribe(userID)
	defer unsub()

	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: broadcaster},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	provider := &mockContextsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getK8sContextFn: func(token, id string) (models.K8sContext, error) {
			return models.K8sContext{}, errors.New("context not found")
		},
	}

	contextID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+contextID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": contextID.String()})
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "test-token")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, &models.User{ID: userID}, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d. body: %s", http.StatusInternalServerError, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrGetK8sContextsCode) {
		t.Fatalf("expected error code %s in body, got %s", ErrGetK8sContextsCode, rec.Body.String())
	}

	// Verify exactly one event was persisted with Error severity
	if len(provider.persistedEvents) != 1 {
		t.Fatalf("expected 1 persisted event, got %d", len(provider.persistedEvents))
	}
	ev := provider.persistedEvents[0]
	if ev.Severity != events.Error {
		t.Fatalf("expected event severity %v, got %v", events.Error, ev.Severity)
	}
	expectedDesc := "Failed to delete connection for " + contextID.String()
	if ev.Description != expectedDesc {
		t.Fatalf("expected event description %q, got %q", expectedDesc, ev.Description)
	}
	if ev.Metadata == nil || ev.Metadata["error"] == nil {
		t.Fatalf("expected error metadata in event, got %v", ev.Metadata)
	}

	// Verify broadcast received
	select {
	case broadcastMsg := <-ch:
		bEvent, ok := broadcastMsg.(*events.Event)
		if !ok {
			t.Fatalf("expected *events.Event from broadcaster, got %T", broadcastMsg)
		}
		if bEvent.Severity != events.Error {
			t.Fatalf("expected broadcast event severity %v, got %v", events.Error, bEvent.Severity)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for broadcast event")
	}
}

// TestDeleteContext_GetK8sContextError_WithNameFallback verifies that the error event
// names the context by its name when the failed lookup still returned one.
func TestDeleteContext_GetK8sContextError_WithNameFallback(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	provider := &mockContextsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getK8sContextFn: func(token, id string) (models.K8sContext, error) {
			return models.K8sContext{Name: "production-cluster"}, errors.New("provider failure")
		},
	}

	contextID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+contextID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": contextID.String()})
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "test-token")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if len(provider.persistedEvents) != 1 {
		t.Fatalf("expected 1 persisted event, got %d", len(provider.persistedEvents))
	}
	ev := provider.persistedEvents[0]
	if ev.Severity != events.Error {
		t.Fatalf("expected event severity %v, got %v", events.Error, ev.Severity)
	}
	expectedDesc := "Failed to delete connection for production-cluster"
	if ev.Description != expectedDesc {
		t.Fatalf("expected event description %q, got %q", expectedDesc, ev.Description)
	}
}

// A failed read must reach the client as an error, not as a 200 with an empty
// page. Uses the real local provider so the persister's error path is covered.
func TestGetAllContexts_DatabaseFailure(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to access underlying sql.DB: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	// The k8s_contexts table is deliberately never migrated.
	provider := &models.DefaultLocalProvider{
		MesheryK8sContextPersister: &models.MesheryK8sContextPersister{DB: &database.Handler{DB: db}},
	}
	h := &Handler{log: newTestLogger(t)}

	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/contexts", nil)
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
	rec := httptest.NewRecorder()

	h.GetAllContexts(rec, req, nil, nil, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d. body: %s", http.StatusInternalServerError, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrGetK8sContextsCode) {
		t.Fatalf("expected error code %s in body, got %s", ErrGetK8sContextsCode, rec.Body.String())
	}
}

func TestDeleteContext_SendEventError(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())
	contextID := uuid.Must(uuid.NewV4())

	broadcaster := &models.Broadcast{}
	broadcastCh, unsubscribe := broadcaster.Subscribe(userID)
	defer unsubscribe()

	asyncErrorPersisted := make(chan events.Event, 1)

	provider := &mockContextsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getK8sContextFn: func(token, id string) (models.K8sContext, error) {
			return models.K8sContext{Name: "test-cluster"}, nil
		},
		persistEventFn: func(event events.Event, token string) error {
			if event.Severity == events.Error {
				asyncErrorPersisted <- event
			}
			return nil
		},
	}

	log := newTestLogger(t)
	inst, err := kubernetes.New(contextID.String(), userID, log)
	if err != nil {
		t.Fatalf("failed to create Kubernetes state machine: %v", err)
	}

	// A non-nil but invalid context passes HasMachineContext and makes
	// DeleteAction.Execute fail its Kubernetes MachineCtx type assertion.
	// Use a non-zero context that passes HasMachineContext but fails
	// DeleteAction.Execute's *kubernetes.MachineCtx type assertion.
	inst.Context = struct {
		marker int
	}{marker: 1}

	if !mhelpers.HasMachineContext(inst) {
		t.Fatal("test setup failed: machine context is not detected")
	}

	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: broadcaster},
		log:                                     log,
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: trackerWith(contextID, inst),
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/system/kubernetes/contexts/"+contextID.String(),
		nil,
	)
	req = mux.SetURLVars(req, map[string]string{"id": contextID.String()})
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))

	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, &models.User{ID: userID}, provider)

	// Wait for the asynchronous SendEvent failure to be persisted.
	var errorEvent events.Event
	select {
	case errorEvent = <-asyncErrorPersisted:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for asynchronous DeleteAction error event to be persisted")
	}

	if errorEvent.Severity != events.Error {
		t.Fatalf("expected Error severity, got %v", errorEvent.Severity)
	}
	if errorEvent.Metadata == nil || errorEvent.Metadata["error"] == nil {
		t.Fatalf("expected error metadata, got %v", errorEvent.Metadata)
	}

	// The initial informational delete-request event and the asynchronous
	// failure event should both have been persisted.
	if len(provider.persistedEvents) != 2 {
		t.Fatalf("expected 2 persisted events, got %d", len(provider.persistedEvents))
	}

	// The asynchronous failure should also be broadcast to the user.
	select {
	case message := <-broadcastCh:
		broadcastEvent, ok := message.(*events.Event)
		if !ok {
			t.Fatalf("expected *events.Event from broadcaster, got %T", message)
		}
		if broadcastEvent.Severity != events.Error {
			t.Fatalf("expected broadcast Error severity, got %v", broadcastEvent.Severity)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for asynchronous DeleteAction error broadcast")
	}

	// The machine must be evicted from the instance tracker even when SendEvent fails,
	// preventing deleted context instances from leaking in memory and being polled by middlewares.
	// The deferred Remove runs when the teardown goroutine returns, after the
	// broadcast above, so poll for eviction instead of asserting it immediately.
	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, ok := h.ConnectionToStateMachineInstanceTracker.Get(contextID); !ok {
			break
		}
		if time.Now().After(deadline) {
			t.Errorf("expected state machine %s to be evicted from tracker on SendEvent failure, but it was leaked", contextID)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

type contextCheckingAction struct {
	receivedCtx            context.Context
	errAtExecution         error
	hasDeadlineAtExecution bool
	done                   chan struct{}
}

func (a *contextCheckingAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (a *contextCheckingAction) Execute(ctx context.Context, _ interface{}, _ interface{}) (machines.EventType, *events.Event, error) {
	a.receivedCtx = ctx
	a.errAtExecution = ctx.Err()
	_, a.hasDeadlineAtExecution = ctx.Deadline()
	close(a.done)
	return machines.NoOp, nil, nil
}

func (a *contextCheckingAction) ExecuteOnExit(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

// TestDeleteContext_DetachesContextFromRequest verifies that background
// SendEvent receives a detached context immune to HTTP request cancellation,
// while preserving context values and enforcing a bounded execution deadline.
func TestDeleteContext_DetachesContextFromRequest(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: map[uuid.UUID]*machines.StateMachine{},
	}
	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     newTestLogger(t),
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: tracker,
	}

	contextID := uuid.Must(uuid.NewV4())
	provider := &mockContextsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getK8sContextFn: func(token, id string) (models.K8sContext, error) {
			return models.K8sContext{
				ID:           strings.ReplaceAll(contextID.String(), "-", ""),
				Name:         "test-cluster",
				ConnectionID: contextID.String(),
			}, nil
		},
	}

	action := &contextCheckingAction{done: make(chan struct{})}
	sm := &machines.StateMachine{
		ID:            core.Uuid(contextID),
		Name:          "kubernetes",
		InitialState:  machines.InitialState,
		CurrentState:  machines.DISCOVERED,
		PreviousState: machines.InitialState,
		Log:           h.log,
		Provider:      provider,
		Context:       &kubernetes.MachineCtx{MesheryCtrlsHelper: &models.MesheryControllersHelper{}},
		States: machines.States{
			machines.DISCOVERED: machines.State{
				Events: machines.Events{machines.Delete: machines.DELETED},
			},
			machines.DELETED: machines.State{
				Action: action,
			},
		},
	}
	tracker.Add(contextID, sm)

	reqCtx, cancel := context.WithCancel(context.Background())
	reqCtx = context.WithValue(reqCtx, models.TokenCtxKey, "secret-token")
	// Cancel the request context immediately to simulate the client disconnecting
	// or the HTTP server unwinding the request handler.
	cancel()

	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+contextID.String(), nil)
	req = req.WithContext(reqCtx)
	req = mux.SetURLVars(req, map[string]string{"id": contextID.String()})
	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	select {
	case <-action.done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for background SendEvent to execute")
	}

	if action.errAtExecution != nil {
		t.Fatalf("expected detached context during execution to have nil Err, got %v", action.errAtExecution)
	}
	if !action.hasDeadlineAtExecution {
		t.Fatal("expected detached context to have a bounded deadline during execution")
	}
	if tok, ok := action.receivedCtx.Value(models.TokenCtxKey).(string); !ok || tok != "secret-token" {
		t.Fatalf("expected TokenCtxKey to be preserved, got %v", tok)
	}

	// Wait for the background goroutine to finish SendEvent and evict the
	// state machine from the tracker, instead of assuming a fixed delay.
	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, ok := tracker.Get(contextID); !ok {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("expected state machine to be evicted from tracker after deletion")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestDeleteContext_InvalidContextID(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	provider := &mockContextsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/invalid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
		t.Fatalf("expected error code %s in body, got %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
	}
}

func TestDeleteContext_NilContextID(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	provider := &mockContextsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+uuid.Nil.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": uuid.Nil.String()})
	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
		t.Fatalf("expected error code %s in body, got %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
	}
}

func TestDeleteContext_NilProvider(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}

	contextID := uuid.Must(uuid.NewV4())
	reqCtx := context.WithValue(context.Background(), models.TokenCtxKey, "secret-token")
	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+contextID.String(), nil)
	req = req.WithContext(reqCtx)
	req = mux.SetURLVars(req, map[string]string{"id": contextID.String()})
	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestDeleteContext_NilUser_DoesNotPanic(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	provider := &mockContextsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
	}

	contextID := uuid.Must(uuid.NewV4())
	reqCtx := context.WithValue(context.Background(), models.TokenCtxKey, "secret-token")
	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+contextID.String(), nil)
	req = req.WithContext(reqCtx)
	req = mux.SetURLVars(req, map[string]string{"id": contextID.String()})
	rec := httptest.NewRecorder()

	// Must not panic with user == nil
	h.DeleteContext(rec, req, nil, nil, provider)
}

func TestDeleteContext_NilSystemID_DoesNotPanic(t *testing.T) {
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: nil,
	}
	provider := &mockContextsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
	}

	contextID := uuid.Must(uuid.NewV4())
	reqCtx := context.WithValue(context.Background(), models.TokenCtxKey, "secret-token")
	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+contextID.String(), nil)
	req = req.WithContext(reqCtx)
	req = mux.SetURLVars(req, map[string]string{"id": contextID.String()})
	rec := httptest.NewRecorder()

	// Must not panic with SystemID == nil
	h.DeleteContext(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
}

func TestGetContext_InvalidContextID(t *testing.T) {
	h := &Handler{
		log: newTestLogger(t),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/contexts/invalid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
	rec := httptest.NewRecorder()

	h.GetContext(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
		t.Fatalf("expected error code %s in body, got %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
	}
}

func TestGetContext_NilContextID(t *testing.T) {
	h := &Handler{
		log: newTestLogger(t),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/contexts/"+uuid.Nil.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": uuid.Nil.String()})
	rec := httptest.NewRecorder()

	h.GetContext(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
		t.Fatalf("expected error code %s in body, got %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
	}
}

func TestGetContext_NilProvider(t *testing.T) {
	h := &Handler{
		log: newTestLogger(t),
	}

	contextID := uuid.Must(uuid.NewV4())
	reqCtx := context.WithValue(context.Background(), models.TokenCtxKey, "secret-token")
	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/contexts/"+contextID.String(), nil)
	req = req.WithContext(reqCtx)
	req = mux.SetURLVars(req, map[string]string{"id": contextID.String()})
	rec := httptest.NewRecorder()

	h.GetContext(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestGetAllContexts_NilProvider(t *testing.T) {
	h := &Handler{
		log: newTestLogger(t),
	}

	reqCtx := context.WithValue(context.Background(), models.TokenCtxKey, "secret-token")
	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/contexts", nil)
	req = req.WithContext(reqCtx)
	rec := httptest.NewRecorder()

	h.GetAllContexts(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
