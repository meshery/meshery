package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

type lifecycleTestMockProvider struct {
	models.Provider
	getK8sContextErr error
	k8sContext       models.K8sContext
	connection       *connections.Connection
	updateConnErr    error
	events           []events.Event
	mu               sync.Mutex
}

func (m *lifecycleTestMockProvider) GetProviderToken(req *http.Request) (string, error) {
	return "test-token", nil
}

func (m *lifecycleTestMockProvider) UpdateConnectionById(token string, conn *connections.ConnectionPayload, connID string) (*connections.Connection, error) {
	if m.updateConnErr != nil {
		return nil, m.updateConnErr
	}
	name := "test-connection"
	if m.connection != nil && m.connection.Name != "" {
		name = m.connection.Name
	}
	return &connections.Connection{ID: uuid.FromStringOrNil(connID), Name: name, Status: conn.Status, Kind: "kubernetes"}, nil
}

func (m *lifecycleTestMockProvider) GetK8sContext(token string, contextID string) (models.K8sContext, error) {
	if m.getK8sContextErr != nil {
		return models.K8sContext{}, m.getK8sContextErr
	}
	return m.k8sContext, nil
}

func (m *lifecycleTestMockProvider) PersistEvent(event events.Event, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
	return nil
}

func (m *lifecycleTestMockProvider) PersistSystemEvent(event events.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
	return nil
}

func (m *lifecycleTestMockProvider) GetConnectionByID(token string, connectionID core.Uuid) (*connections.Connection, int, error) {
	return &connections.Connection{ID: uuid.UUID(connectionID), Name: "test-connection", Kind: "kubernetes"}, http.StatusOK, nil
}

func (m *lifecycleTestMockProvider) GetGenericPersister() *database.Handler {
	return nil
}

func newLifecycleTestHandler(t *testing.T, provider *lifecycleTestMockProvider) (*Handler, *machines.ConnectionToStateMachineInstanceTracker) {
	t.Helper()
	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
	sysID := uuid.Must(uuid.NewV4())
	log := newTestLogger(t)
	h := &Handler{
		ConnectionToStateMachineInstanceTracker: tracker,
		config:                                  &models.HandlerConfig{EventBroadcaster: models.NewBroadcaster("test")},
		log:                                     log,
		SystemID:                                &sysID,
		MesheryCtrlsHelper:                      models.NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, provider, &sysID),
	}
	return h, tracker
}

// TestUpdateConnectionById_FailsOnLifecycleError verifies that when NotifySmOfConnectionStatusChange fails,
// UpdateConnectionById propagates the error as HTTP 500 and persists/publishes the failure event.
func TestUpdateConnectionById_FailsOnLifecycleError(t *testing.T) {
	connID := uuid.Must(uuid.NewV4())
	provider := &lifecycleTestMockProvider{
		getK8sContextErr: fmt.Errorf("k8s context not found"),
	}
	h, _ := newLifecycleTestHandler(t, provider)

	body := strings.NewReader(`{"status":"connected"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), body)
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	reqCtx := req.Context()
	reqCtx = context.WithValue(reqCtx, models.TokenCtxKey, "test-token")
	req = req.WithContext(reqCtx)
	rec := httptest.NewRecorder()

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	h.UpdateConnectionById(rec, req, nil, user, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 on lifecycle transition error, got: %d (body: %s)", rec.Code, rec.Body.String())
	}

	provider.mu.Lock()
	defer provider.mu.Unlock()
	var errorEventFound bool
	for _, ev := range provider.events {
		if ev.Severity == events.Error {
			errorEventFound = true
			if ev.ActedUpon != core.Uuid(connID) {
				t.Fatalf("expected error event ActedUpon %s, got %s", connID, ev.ActedUpon)
			}
		}
	}
	if !errorEventFound {
		t.Fatal("expected failure event to be persisted on lifecycle transition error, but none was recorded")
	}
}

// TestUpdateConnectionById_StatusOnlyUsesURLConnectionID verifies that status-only PUT requests
// (with no ID in request body) populate connection.ID from the URL route parameter.
func TestUpdateConnectionById_StatusOnlyUsesURLConnectionID(t *testing.T) {
	connID := uuid.Must(uuid.NewV4())
	log := newTestLogger(t)
	provider := &lifecycleTestMockProvider{
		k8sContext: models.K8sContext{
			ID:           connID.String(),
			Name:         "test-k8s",
			ConnectionID: connID.String(),
		},
	}
	h, tracker := newLifecycleTestHandler(t, provider)

	sm, err := kubernetes.New(connID.String(), core.Uuid(connID), log)
	if err != nil {
		t.Fatalf("failed to create machine: %v", err)
	}
	sm.Provider = provider
	state := sm.States[machines.CONNECTED]
	sm.States[machines.CONNECTED] = *state.RegisterAction(&successfulLifecycleAction{})
	sm.CurrentState = machines.REGISTERED

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")

	machineCtx := &kubernetes.MachineCtx{
		K8sContext:  provider.k8sContext,
		ActionMutex: &sync.Mutex{},
	}
	_, err = sm.Start(ctx, machineCtx, log, func(c context.Context, mc interface{}, l logger.Handler) (interface{}, *events.Event, error) {
		return mc, nil, nil
	})
	if err != nil {
		t.Fatalf("start machine: %v", err)
	}

	tracker.Add(core.Uuid(connID), sm)

	body := strings.NewReader(`{"status":"connected"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), body)
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	reqCtx := req.Context()
	reqCtx = context.WithValue(reqCtx, models.TokenCtxKey, "test-token")
	reqCtx = context.WithValue(reqCtx, models.UserCtxKey, user)
	reqCtx = context.WithValue(reqCtx, models.SystemIDKey, &sysID)
	req = req.WithContext(reqCtx)
	rec := httptest.NewRecorder()

	h.UpdateConnectionById(rec, req, nil, user, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 for status-only PUT with valid URL connection ID, got: %d (body: %s)", rec.Code, rec.Body.String())
	}
}

// TestDeleteContext_FailsOnSendEventError verifies that when SendEvent fails during DeleteContext,
// the error is propagated with HTTP 500 and the machine is NOT removed from smInstanceTracker.
func TestDeleteContext_FailsOnSendEventError(t *testing.T) {
	connID := uuid.Must(uuid.NewV4())
	log := newTestLogger(t)

	provider := &lifecycleTestMockProvider{
		k8sContext: models.K8sContext{
			ID:           connID.String(),
			Name:         "test-k8s",
			ConnectionID: connID.String(),
		},
	}
	h, tracker := newLifecycleTestHandler(t, provider)

	// Create and register a kubernetes machine in tracker
	sm, err := kubernetes.New(connID.String(), core.Uuid(connID), log)
	if err != nil {
		t.Fatalf("failed to create machine: %v", err)
	}
	sm.Provider = provider
	// Replace delete action with a failing action
	state := sm.States[machines.DELETED]
	sm.States[machines.DELETED] = *state.RegisterAction(&failingLifecycleAction{})
	sm.CurrentState = machines.DISCOVERED

	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))})
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")

	machineCtx := &kubernetes.MachineCtx{
		K8sContext:  provider.k8sContext,
		ActionMutex: &sync.Mutex{},
	}
	_, err = sm.Start(ctx, machineCtx, log, func(c context.Context, mc interface{}, l logger.Handler) (interface{}, *events.Event, error) {
		return mc, nil, nil
	})
	if err != nil {
		t.Fatalf("start machine: %v", err)
	}

	tracker.Add(core.Uuid(connID), sm)

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+connID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": connID.String()})
	reqCtx := req.Context()
	reqCtx = context.WithValue(reqCtx, models.TokenCtxKey, "test-token")
	reqCtx = context.WithValue(reqCtx, models.UserCtxKey, user)
	reqCtx = context.WithValue(reqCtx, models.SystemIDKey, &sysID)
	req = req.WithContext(reqCtx)
	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, user, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 on SendEvent error, got: %d (body: %s)", rec.Code, rec.Body.String())
	}

	// Verify failure event was persisted
	provider.mu.Lock()
	var deleteErrorFound bool
	for _, ev := range provider.events {
		if ev.Severity == events.Error {
			deleteErrorFound = true
		}
	}
	provider.mu.Unlock()
	if !deleteErrorFound {
		t.Fatal("expected failure event to be persisted on DeleteContext SendEvent failure, but none was found")
	}

	// Verify machine was NOT removed from tracker
	if _, ok := tracker.Get(core.Uuid(connID)); !ok {
		t.Fatal("expected machine to remain in tracker after SendEvent failure, but it was removed")
	}
}

// TestDeleteContext_SuccessRemovesTracker verifies that on successful Delete transition,
// the machine is removed from smInstanceTracker and HTTP 200 is returned.
func TestDeleteContext_SuccessRemovesTracker(t *testing.T) {
	connID := uuid.Must(uuid.NewV4())
	log := newTestLogger(t)

	provider := &lifecycleTestMockProvider{
		k8sContext: models.K8sContext{
			ID:           connID.String(),
			Name:         "test-k8s",
			ConnectionID: connID.String(),
		},
	}
	h, tracker := newLifecycleTestHandler(t, provider)

	sm, err := kubernetes.New(connID.String(), core.Uuid(connID), log)
	if err != nil {
		t.Fatalf("failed to create machine: %v", err)
	}
	sm.Provider = provider
	// Replace delete action with a successful action
	state := sm.States[machines.DELETED]
	sm.States[machines.DELETED] = *state.RegisterAction(&successfulLifecycleAction{})
	sm.CurrentState = machines.DISCOVERED

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")

	machineCtx := &kubernetes.MachineCtx{
		K8sContext:  provider.k8sContext,
		ActionMutex: &sync.Mutex{},
	}
	_, err = sm.Start(ctx, machineCtx, log, func(c context.Context, mc interface{}, l logger.Handler) (interface{}, *events.Event, error) {
		return mc, nil, nil
	})
	if err != nil {
		t.Fatalf("start machine: %v", err)
	}

	tracker.Add(core.Uuid(connID), sm)

	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+connID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": connID.String()})
	reqCtx := req.Context()
	reqCtx = context.WithValue(reqCtx, models.TokenCtxKey, "test-token")
	reqCtx = context.WithValue(reqCtx, models.UserCtxKey, user)
	reqCtx = context.WithValue(reqCtx, models.SystemIDKey, &sysID)
	req = req.WithContext(reqCtx)
	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, user, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on successful delete, got: %d (body: %s)", rec.Code, rec.Body.String())
	}

	// Verify machine was removed from tracker
	if _, ok := tracker.Get(core.Uuid(connID)); ok {
		t.Fatal("expected machine to be removed from tracker after successful Delete, but it is still present")
	}
}

type failingLifecycleAction struct{}

func (a *failingLifecycleAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, events.NewEvent().WithSeverity(events.Error).WithDescription("action failed").Build(), fmt.Errorf("simulated action failure")
}

func (a *failingLifecycleAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (a *failingLifecycleAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

type successfulLifecycleAction struct{}

func (a *successfulLifecycleAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, events.NewEvent().WithSeverity(events.Informational).WithDescription("deleted successfully").Build(), nil
}

func (a *successfulLifecycleAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (a *successfulLifecycleAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
