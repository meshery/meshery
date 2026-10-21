package kubernetes

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/core"
)

type fakeK8sProvider struct {
	models.Provider
	connection *connections.Connection
	connErr    error
}

func (f *fakeK8sProvider) GetConnectionByID(_ string, id core.Uuid) (*connections.Connection, int, error) {
	if f.connErr != nil {
		return nil, 500, f.connErr
	}
	if f.connection != nil {
		return f.connection, 200, nil
	}
	return &connections.Connection{
		ID:     id,
		Kind:   "kubernetes",
		Status: connections.DISCONNECTED,
	}, 200, nil
}

func (f *fakeK8sProvider) SaveK8sContext(_ string, _ models.K8sContext, _ map[string]any) (connections.Connection, error) {
	return connections.Connection{ID: core.Uuid(uuid.Must(uuid.NewV4()))}, nil
}


func newTestMachineCtx(t *testing.T) *MachineCtx {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build logger: %v", err)
	}
	connID := uuid.Must(uuid.NewV4()).String()
	return &MachineCtx{
		K8sContext: models.K8sContext{
			ID:           "test-ctx-id",
			Name:         "test-cluster",
			Server:       "https://127.0.0.1:6443",
			ConnectionID: connID,
		},
		log: log,
	}
}

func TestConnectAction_MissingProvider(t *testing.T) {
	action := &ConnectAction{}
	mCtx := newTestMachineCtx(t)

	// Context with user, systemID, token, but NO Provider
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: userID})
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")

	nextEvent, event, err := action.Execute(ctx, mCtx, nil)
	if err == nil {
		t.Fatal("expected error when provider is missing, got nil")
	}
	if nextEvent != machines.NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event == nil {
		t.Fatal("expected non-nil event on error")
	}
	if event.Metadata == nil || event.Metadata["error"] == nil {
		t.Fatalf("expected error metadata on event, got %v", event.Metadata)
	}
}

func TestConnectAction_TypedNilProvider(t *testing.T) {
	action := &ConnectAction{}
	mCtx := newTestMachineCtx(t)

	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	var fp *fakeK8sProvider
	var nilProvider models.Provider = fp
	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: userID})
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")
	ctx = context.WithValue(ctx, models.ProviderCtxKey, nilProvider)

	nextEvent, event, err := action.Execute(ctx, mCtx, nil)
	if err == nil {
		t.Fatal("expected error on typed-nil provider, got nil")
	}
	if nextEvent != machines.NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event == nil {
		t.Fatal("expected non-nil event on typed-nil provider")
	}
	if event.Metadata == nil || event.Metadata["error"] == nil {
		t.Fatalf("expected error metadata on event, got %v", event.Metadata)
	}
}

func TestConnectAction_MissingToken(t *testing.T) {
	action := &ConnectAction{}
	mCtx := newTestMachineCtx(t)
	prov := &fakeK8sProvider{}

	ctx := context.WithValue(context.Background(), models.ProviderCtxKey, prov)
	// Deliberately omit TokenCtxKey

	nextEvent, event, err := action.Execute(ctx, mCtx, nil)
	if err == nil {
		t.Fatal("expected error when token is missing, got nil")
	}
	if nextEvent != machines.NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event == nil {
		t.Fatal("expected non-nil event on token error")
	}
}

func TestConnectAction_EmptyContextDoesNotPanic(t *testing.T) {
	action := &ConnectAction{}
	mCtx := newTestMachineCtx(t)

	// Empty context.Background() without user, provider, token, or sysID
	nextEvent, event, err := action.Execute(context.Background(), mCtx, nil)
	if err == nil {
		t.Fatal("expected error on empty context, got nil")
	}
	if nextEvent != machines.NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event == nil {
		t.Fatal("expected non-nil event on error")
	}
}

func TestDisconnectAction_DegradedContextDoesNotPanic(t *testing.T) {
	action := &DisconnectAction{}
	mCtx := newTestMachineCtx(t)

	// Completely empty context
	nextEvent, event, err := action.Execute(context.Background(), mCtx, nil)
	if err != nil {
		t.Fatalf("unexpected error on disconnect action: %v", err)
	}
	if nextEvent != machines.NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event != nil {
		t.Fatalf("expected nil event on success, got %v", event)
	}
}

func TestDeleteAction_DegradedContextDoesNotPanic(t *testing.T) {
	action := &DeleteAction{}
	mCtx := newTestMachineCtx(t)

	// Completely empty context
	nextEvent, event, err := action.Execute(context.Background(), mCtx, nil)
	if err != nil {
		t.Fatalf("unexpected error on delete action: %v", err)
	}
	if nextEvent != machines.NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event != nil {
		t.Fatalf("expected nil event on success, got %v", event)
	}
}

func TestDiscoverAction_MissingProvider(t *testing.T) {
	action := &DiscoverAction{}
	mCtx := newTestMachineCtx(t)

	// Context without provider
	nextEvent, event, err := action.Execute(context.Background(), mCtx, nil)
	if err == nil {
		t.Fatal("expected error on missing provider, got nil")
	}
	if nextEvent != machines.NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event == nil {
		t.Fatal("expected non-nil event")
	}
	if event.Metadata == nil || event.Metadata["error"] == nil {
		t.Fatalf("expected error metadata on event, got %v", event.Metadata)
	}
}

func TestDiscoverAction_TypedNilProvider(t *testing.T) {
	action := &DiscoverAction{}
	mCtx := newTestMachineCtx(t)

	var fp *fakeK8sProvider
	var nilProvider models.Provider = fp
	ctx := context.WithValue(context.Background(), models.ProviderCtxKey, nilProvider)

	nextEvent, event, err := action.Execute(ctx, mCtx, nil)
	if err == nil {
		t.Fatal("expected error on typed-nil provider, got nil")
	}
	if nextEvent != machines.NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event == nil {
		t.Fatal("expected non-nil event")
	}
	if event.Metadata == nil || event.Metadata["error"] == nil {
		t.Fatalf("expected error metadata on event, got %v", event.Metadata)
	}
}

func TestRegisterAction_MissingProvider(t *testing.T) {
	action := &RegisterAction{}
	mCtx := newTestMachineCtx(t)

	// Context without provider
	nextEvent, event, err := action.Execute(context.Background(), mCtx, nil)
	// Execute checks provider or performs ping test
	if err == nil {
		t.Fatal("expected error on degraded context, got nil")
	}
	if event == nil {
		t.Fatal("expected non-nil event on error")
	}
	if event.Metadata == nil || event.Metadata["error"] == nil {
		t.Fatalf("expected error metadata on event, got %v", event.Metadata)
	}
	_ = nextEvent
}

func TestAssignInitialCtx_DegradedContextDoesNotPanic(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build logger: %v", err)
	}
	mCtx := newTestMachineCtx(t)

	// Degraded context without user or system ID
	_, _, err = AssignInitialCtx(context.Background(), mCtx, log)
	// ClientSet assignment will fail on unreachable localhost server, but it must NOT panic on nil user/sysID
	if err == nil {
		t.Fatal("expected error on unreachable cluster context")
	}
}
