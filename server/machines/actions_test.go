package machines

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/schemas/models/core"
)

type mockActionProvider struct {
	fakeProvider
	saveCredErr error
	saveConnErr error
	savedCred   *models.Credential
	savedConn   *connections.ConnectionPayload
}

func (m *mockActionProvider) SaveUserCredential(_ string, cred *models.Credential) (*models.Credential, error) {
	m.savedCred = cred
	if m.saveCredErr != nil {
		return nil, m.saveCredErr
	}
	credID := core.Uuid(uuid.Must(uuid.NewV4()))
	cred.ID = credID
	return cred, nil
}

func (m *mockActionProvider) SaveConnection(conn *connections.ConnectionPayload, _ string, _ bool) (*connections.Connection, error) {
	m.savedConn = conn
	if m.saveConnErr != nil {
		return nil, m.saveConnErr
	}
	return &connections.Connection{
		ID:     conn.ID,
		Kind:   conn.Kind,
		Status: conn.Status,
	}, nil
}

func TestDefaultConnectAction_MissingProvider(t *testing.T) {
	action := &DefaultConnectAction{}
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx := context.WithValue(context.Background(), models.SystemIDKey, &sysID)

	payload := connections.ConnectionPayload{
		Name: "test-conn",
		Kind: "grafana",
	}

	nextEvent, event, err := action.Execute(ctx, nil, payload)
	if err == nil {
		t.Fatal("expected error when provider is missing, got nil")
	}
	if nextEvent != NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event == nil {
		t.Fatal("expected non-nil event when provider is missing")
	}
	if event.Metadata == nil || event.Metadata["error"] == nil {
		t.Fatalf("expected error metadata on event, got %v", event.Metadata)
	}
}

func TestDefaultConnectAction_TypedNilProvider(t *testing.T) {
	action := &DefaultConnectAction{}
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	var fp *mockActionProvider
	var nilProvider models.Provider = fp
	ctx := context.WithValue(context.Background(), models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.ProviderCtxKey, nilProvider)

	payload := connections.ConnectionPayload{
		Name: "test-conn",
		Kind: "grafana",
	}

	nextEvent, event, err := action.Execute(ctx, nil, payload)
	if err == nil {
		t.Fatal("expected error on typed-nil provider, got nil")
	}
	if nextEvent != NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event == nil {
		t.Fatal("expected non-nil event")
	}
	if event.Metadata == nil || event.Metadata["error"] == nil {
		t.Fatalf("expected error metadata on event, got %v", event.Metadata)
	}
}

func TestDefaultConnectAction_MissingAttribution(t *testing.T) {
	action := &DefaultConnectAction{}
	mockProv := &mockActionProvider{}
	ctx := context.WithValue(context.Background(), models.ProviderCtxKey, mockProv)
	// Completely missing user and system IDs

	payload := connections.ConnectionPayload{
		Name: "test-conn",
		Kind: "grafana",
	}

	nextEvent, event, err := action.Execute(ctx, nil, payload)
	if err == nil {
		t.Fatal("expected error when event attribution is missing, got nil")
	}
	if nextEvent != NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event == nil {
		t.Fatal("expected non-nil event")
	}
}

func TestDefaultConnectAction_ExistingCredentialWithSystemAttribution(t *testing.T) {
	action := &DefaultConnectAction{}
	mockProv := &mockActionProvider{}

	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx := context.WithValue(context.Background(), models.ProviderCtxKey, mockProv)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)

	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	credID := uuid.Must(uuid.NewV4()).String()
	payload := connections.ConnectionPayload{
		ID:   connID,
		Name: "test-conn",
		Kind: "grafana",
		CredentialSecret: map[string]interface{}{
			"id": credID,
		},
	}

	nextEvent, event, err := action.Execute(ctx, nil, payload)
	if err != nil {
		t.Fatalf("unexpected error on system context with existing credential: %v", err)
	}
	if nextEvent != NoOp {
		t.Fatalf("expected next event NoOp, got %v", nextEvent)
	}
	if event != nil {
		t.Fatalf("expected nil event on success, got %v", event)
	}
	if mockProv.savedConn == nil || mockProv.savedConn.ID != connID {
		t.Fatalf("expected saved connection with ID %v", connID)
	}
}
