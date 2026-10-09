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
