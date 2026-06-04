package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/events"
)

type modelProviderConnectionSpyProvider struct {
	*models.DefaultLocalProvider
	observedSave      atomic.Pointer[connections.ConnectionPayload]
	observedUpdate    atomic.Pointer[connections.ConnectionPayload]
	observedListTypes atomic.Value
	observedListKinds atomic.Value
	connection        *connections.Connection
}

func newModelProviderConnectionSpyProvider(connection *connections.Connection) *modelProviderConnectionSpyProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &modelProviderConnectionSpyProvider{
		DefaultLocalProvider: base,
		connection:           connection,
	}
}

func (m *modelProviderConnectionSpyProvider) SaveConnection(conn *connections.ConnectionPayload, _ string, _ bool) (*connections.Connection, error) {
	cp := *conn
	m.observedSave.Store(&cp)
	return m.connection, nil
}

func (m *modelProviderConnectionSpyProvider) GetConnections(_ *http.Request, _ string, page, pageSize int, _ string, _ string, _ string, _ []string, kind []string, connType []string, _ string) (*connections.ConnectionPage, error) {
	m.observedListKinds.Store(append([]string{}, kind...))
	m.observedListTypes.Store(append([]string{}, connType...))
	return &connections.ConnectionPage{
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  1,
		Connections: []*connections.Connection{m.connection},
	}, nil
}

func (m *modelProviderConnectionSpyProvider) UpdateConnectionById(_ string, conn *connections.ConnectionPayload, _ string) (*connections.Connection, error) {
	cp := *conn
	m.observedUpdate.Store(&cp)
	return m.connection, nil
}

func (m *modelProviderConnectionSpyProvider) PersistEvent(_ events.Event, _ string) error {
	return nil
}

func newModelProviderConnectionTestHandler(t *testing.T) *Handler {
	t.Helper()

	h := newTestHandler(t, map[string]models.Provider{}, "")
	systemID := uuid.Must(uuid.FromString("99999999-9999-9999-9999-999999999999"))
	h.SystemID = &systemID
	h.config.EventBroadcaster = models.NewBroadcaster("test")
	return h
}

func TestSaveConnectionAcceptsCloudModelProviderWithCredential(t *testing.T) {
	credentialID := uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))
	provider := newModelProviderConnectionSpyProvider(&connections.Connection{Name: "OpenAI Production"})
	h := newModelProviderConnectionTestHandler(t)
	authUser := &models.User{ID: uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))}

	body := `{
		"name": "OpenAI Production",
		"kind": "openai",
		"type": "model-provider",
		"subType": "inference",
		"status": "registered",
		"credentialId": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		"metadata": {
			"baseUrl": "https://api.openai.com/v1",
			"defaultModel": "gpt-4o"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections", bytes.NewBufferString(body))
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
	rec := httptest.NewRecorder()

	h.SaveConnection(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	saved := provider.observedSave.Load()
	if saved == nil {
		t.Fatal("provider SaveConnection was not invoked")
	}
	if saved.Kind != connections.ModelProviderKindOpenAI || saved.Type != connections.ModelProviderConnectionType || saved.SubType != connections.ModelProviderConnectionSubType {
		t.Fatalf("unexpected connection classification: kind=%q type=%q subType=%q", saved.Kind, saved.Type, saved.SubType)
	}
	if saved.CredentialID == nil || *saved.CredentialID != credentialID {
		t.Fatalf("credentialId was not preserved: got %v, want %s", saved.CredentialID, credentialID)
	}
}

func TestSaveConnectionAcceptsLocalModelProviderWithoutCredential(t *testing.T) {
	provider := newModelProviderConnectionSpyProvider(&connections.Connection{Name: "Ollama Local"})
	h := newModelProviderConnectionTestHandler(t)
	authUser := &models.User{ID: uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))}

	body := `{
		"name": "Ollama Local",
		"kind": "ollama",
		"type": "model-provider",
		"subType": "inference",
		"status": "registered",
		"metadata": {
			"baseUrl": "http://localhost:11434",
			"defaultModel": "llama3.1"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections", bytes.NewBufferString(body))
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
	rec := httptest.NewRecorder()

	h.SaveConnection(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	saved := provider.observedSave.Load()
	if saved == nil {
		t.Fatal("provider SaveConnection was not invoked")
	}
	if saved.CredentialID != nil {
		t.Fatalf("expected no credentialId for local no-auth provider, got %s", saved.CredentialID.String())
	}
	if saved.MetaData[connections.ModelProviderMetadataBaseURL] != "http://localhost:11434" {
		t.Fatalf("baseUrl metadata was not preserved: %#v", saved.MetaData)
	}
}

func TestGetConnectionsSupportsModelProviderTypeFilter(t *testing.T) {
	provider := newModelProviderConnectionSpyProvider(&connections.Connection{
		Name: "Anthropic Production",
		Kind: connections.ModelProviderKindAnthropic,
		Type: connections.ModelProviderConnectionType,
	})
	h := newModelProviderConnectionTestHandler(t)
	authUser := &models.User{ID: uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))}

	req := httptest.NewRequest(http.MethodGet, `/api/integrations/connections?type=["model-provider"]&kind=["anthropic"]`, nil)
	rec := httptest.NewRecorder()

	h.GetConnections(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	gotTypes, _ := provider.observedListTypes.Load().([]string)
	if !reflect.DeepEqual(gotTypes, []string{connections.ModelProviderConnectionType}) {
		t.Fatalf("got connection type filters %v, want [model-provider]", gotTypes)
	}
	gotKinds, _ := provider.observedListKinds.Load().([]string)
	if !reflect.DeepEqual(gotKinds, []string{connections.ModelProviderKindAnthropic}) {
		t.Fatalf("got connection kind filters %v, want [anthropic]", gotKinds)
	}
}

func TestUpdateConnectionAcceptsModelProviderCredentialLink(t *testing.T) {
	connectionID := uuid.Must(uuid.FromString("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"))
	credentialID := uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))
	provider := newModelProviderConnectionSpyProvider(&connections.Connection{Name: "AWS Bedrock Production"})
	h := newModelProviderConnectionTestHandler(t)
	authUser := &models.User{ID: uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))}

	body := `{
		"id": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		"name": "AWS Bedrock Production",
		"kind": "aws-bedrock",
		"type": "model-provider",
		"subType": "inference",
		"credentialId": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		"metadata": {
			"region": "us-east-1",
			"defaultModel": "anthropic.claude-3-5-sonnet-20240620-v1:0"
		}
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connectionID.String(), bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connectionID.String()})
	rec := httptest.NewRecorder()

	h.UpdateConnectionById(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	updated := provider.observedUpdate.Load()
	if updated == nil {
		t.Fatal("provider UpdateConnectionById was not invoked")
	}
	if updated.CredentialID == nil || *updated.CredentialID != credentialID {
		t.Fatalf("credentialId was not preserved on update: got %v, want %s", updated.CredentialID, credentialID)
	}
}
