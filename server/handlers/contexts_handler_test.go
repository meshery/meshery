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
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
)

type mockContextsProvider struct {
	*models.DefaultLocalProvider
	getK8sContextFn func(token, id string) (models.K8sContext, error)
	persistEventFn  func(event events.Event, token string) error
	persistedEvents []events.Event
}

func (m *mockContextsProvider) GetK8sContext(token, id string) (models.K8sContext, error) {
	if m.getK8sContextFn != nil {
		return m.getK8sContextFn(token, id)
	}
	return models.K8sContext{}, nil
}

func (m *mockContextsProvider) PersistEvent(event events.Event, token string) error {
	m.persistedEvents = append(m.persistedEvents, event)
	if m.persistEventFn != nil {
		return m.persistEventFn(event, token)
	}
	return nil
}

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

// blockingDeleteAction stands in for the kubernetes DeleteAction. It waits until
// the test releases it, so the request context can be cancelled first, then
// reports the context error it observed and returns the configured error.
type blockingDeleteAction struct {
	release <-chan struct{}
	ctxErr  chan error
	err     error
}

func (a *blockingDeleteAction) ExecuteOnEntry(ctx context.Context, _ interface{}, _ interface{}) (machines.EventType, *events.Event, error) {
	<-a.release
	a.ctxErr <- ctx.Err()
	return machines.NoOp, nil, a.err
}

func (a *blockingDeleteAction) Execute(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (a *blockingDeleteAction) ExecuteOnExit(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

// runDeleteContextWithMachine seeds the tracker with an initialized machine for
// the connection, so DeleteContext takes the cache-hit path and drives Delete
// through the given action, then cancels the request context once the handler
// has returned - as net/http does - before letting the action run.
func runDeleteContextWithMachine(t *testing.T, action *blockingDeleteAction, release chan struct{}) (*machines.ConnectionToStateMachineInstanceTracker, uuid.UUID) {
	t.Helper()

	systemID := uuid.Must(uuid.NewV4())
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	connectionID := uuid.Must(uuid.NewV4())

	inst := &machines.StateMachine{
		ID:           connectionID,
		Name:         "kubernetes",
		InitialState: machines.InitialState,
		CurrentState: machines.InitialState,
		States: machines.States{
			machines.InitialState: {Events: machines.Events{machines.Delete: machines.DELETED}},
			machines.DELETED:      {Action: action},
		},
		Context: &kubernetes.MachineCtx{},
		Log:     newTestLogger(t),
	}
	tracker := trackerWith(connectionID, inst)

	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     newTestLogger(t),
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: tracker,
	}
	provider := &mockContextsProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}

	reqCtx, cancel := context.WithCancel(context.Background())
	reqCtx = context.WithValue(reqCtx, models.TokenCtxKey, "test-token")
	reqCtx = context.WithValue(reqCtx, models.UserCtxKey, user)
	reqCtx = context.WithValue(reqCtx, models.SystemIDKey, &systemID)
	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+connectionID.String(), nil).WithContext(reqCtx)
	req = mux.SetURLVars(req, map[string]string{"id": connectionID.String()})

	h.DeleteContext(httptest.NewRecorder(), req, nil, user, provider)
	cancel()
	close(release)

	return tracker, connectionID
}

func waitForTrackerRemoval(t *testing.T, tracker *machines.ConnectionToStateMachineInstanceTracker, id uuid.UUID) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if _, ok := tracker.Get(id); !ok {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("tracker entry for connection %s was not removed", id)
}

// The delete runs in a goroutine that outlives the request, so it must not
// inherit the request's cancellation - otherwise operator undeploy and MeshSync
// teardown abort with context.Canceled as soon as the handler returns.
func TestDeleteContext_DetachesFromRequestContext(t *testing.T) {
	release := make(chan struct{})
	action := &blockingDeleteAction{release: release, ctxErr: make(chan error, 1)}

	tracker, connectionID := runDeleteContextWithMachine(t, action, release)

	select {
	case err := <-action.ctxErr:
		if err != nil {
			t.Fatalf("delete action saw a cancelled context: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for the delete action to run")
	}
	waitForTrackerRemoval(t, tracker, connectionID)
}

func TestDeleteContext_RemovesTrackerEntryWhenSendEventFails(t *testing.T) {
	release := make(chan struct{})
	action := &blockingDeleteAction{release: release, ctxErr: make(chan error, 1), err: errors.New("undeploy failed")}

	tracker, connectionID := runDeleteContextWithMachine(t, action, release)

	waitForTrackerRemoval(t, tracker, connectionID)
}
