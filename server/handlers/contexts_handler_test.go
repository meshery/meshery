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

func TestDeleteContext_TrackerRemovedOnStateMachineFailure(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())
	contextID := uuid.Must(uuid.NewV4())
	testLog := newTestLogger(t)

	// Pre-populate tracker with a state machine instance that has a context
	// but no configured transitions, simulating an invalid transition failure in SendEvent.
	inst := &machines.StateMachine{
		ID:           contextID,
		CurrentState: machines.DefaultState,
		States:       machines.States{},
		Context:      &kubernetes.MachineCtx{},
		Log:          testLog,
	}
	tracker := trackerWith(contextID, inst)

	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     testLog,
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: tracker,
	}

	provider := &mockContextsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getK8sContextFn: func(token, id string) (models.K8sContext, error) {
			return models.K8sContext{
				ID:   id,
				Name: "test-cluster",
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+contextID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": contextID.String()})
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "test-token")
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userID})
	ctx = context.WithValue(ctx, models.SystemIDKey, &systemID)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, &models.User{ID: userID}, provider)

	// Verify tracker entry is removed despite SendEvent failure due to defer cleanup (#22431)
	deadline := time.Now().Add(2 * time.Second)
	removed := false
	for time.Now().Before(deadline) {
		if _, ok := tracker.Get(contextID); !ok {
			removed = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if !removed {
		t.Fatalf("expected connection %s to be removed from tracker even on SendEvent error", contextID)
	}
}

func TestDeleteContext_TrackerRemovedWhenNoMachineContext(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())
	contextID := uuid.Must(uuid.NewV4())

	inst := &machines.StateMachine{
		ID:      contextID,
		Context: nil,
	}
	tracker := trackerWith(contextID, inst)

	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     newTestLogger(t),
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: tracker,
	}

	provider := &mockContextsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getK8sContextFn: func(token, id string) (models.K8sContext, error) {
			return models.K8sContext{
				ID:   id,
				Name: "uninitialized-cluster",
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+contextID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": contextID.String()})
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "test-token")
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userID})
	ctx = context.WithValue(ctx, models.SystemIDKey, &systemID)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.DeleteContext(rec, req, nil, &models.User{ID: userID}, provider)

	if _, ok := tracker.Get(contextID); ok {
		t.Fatalf("expected connection %s to be removed from tracker when machine has no context", contextID)
	}
}
