package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// TestK8sEventMetadataHasError guards the decision that raises a Kubernetes
// connection receipt event to Error severity. A receipt whose per-context
// metadata records any failure must be reported as Error so it persists in the
// notification center and stays retrievable under the Error severity filter
// (issue #20725); a receipt describing only successful connections must remain
// Informational.
func TestK8sEventMetadataHasError(t *testing.T) {
	tests := []struct {
		name          string
		eventMetadata map[string]interface{}
		want          bool
	}{
		{
			name:          "nil metadata reports no error",
			eventMetadata: nil,
			want:          false,
		},
		{
			name:          "empty metadata reports no error",
			eventMetadata: map[string]interface{}{},
			want:          false,
		},
		{
			name: "present-but-nil error value reports no error",
			eventMetadata: map[string]interface{}{
				"prod": map[string]interface{}{
					"error": nil,
				},
			},
			want: false,
		},
		{
			name: "only successful contexts report no error",
			eventMetadata: map[string]interface{}{
				"prod": map[string]interface{}{
					"description": "Connection registered with kubernetes context \"prod\".",
				},
				"staging": map[string]interface{}{
					"description": "Connection already exists with Kubernetes context \"staging\".",
				},
			},
			want: false,
		},
		{
			name: "a single failed context reports an error",
			eventMetadata: map[string]interface{}{
				"unreachable": map[string]interface{}{
					"description": "Unable to establish connection with context \"unreachable\".",
					"error":       errors.New("api server unreachable"),
				},
			},
			want: true,
		},
		{
			name: "a failure mixed with successes reports an error",
			eventMetadata: map[string]interface{}{
				"prod": map[string]interface{}{
					"description": "Connection registered with kubernetes context \"prod\".",
				},
				"unreachable": map[string]interface{}{
					"error": errors.New("api server unreachable"),
				},
			},
			want: true,
		},
		{
			name: "non-map entries are ignored",
			eventMetadata: map[string]interface{}{
				"weird": "not a metadata map",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := k8sEventMetadataHasError(tt.eventMetadata); got != tt.want {
				t.Errorf("k8sEventMetadataHasError() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockK8sConfigProvider struct {
	*models.DefaultLocalProvider
	savedContexts []models.K8sContext
	mu            sync.Mutex
}

// SaveK8sContext records the discovered context in memory and returns a mock DISCOVERED connection.
func (m *mockK8sConfigProvider) SaveK8sContext(_ string, k8sContext models.K8sContext, _ map[string]any) (connections.Connection, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.savedContexts = append(m.savedContexts, k8sContext)
	connID := uuid.Must(uuid.NewV4())
	return connections.Connection{
		ID:     core.Uuid(connID),
		Status: connections.DISCOVERED,
	}, nil
}

// PersistEvent is a mock implementation that satisfies the Provider interface.
func (m *mockK8sConfigProvider) PersistEvent(_ events.Event, _ string) error {
	return nil
}

// GetConnectionByID returns a mock connection for the specified ID.
func (m *mockK8sConfigProvider) GetConnectionByID(_ string, id core.Uuid) (*connections.Connection, int, error) {
	return &connections.Connection{
		ID:     id,
		Status: connections.DISCOVERED,
	}, 1, nil
}

// UpdateConnectionStatusByID updates and returns the status of a mock connection.
func (m *mockK8sConfigProvider) UpdateConnectionStatusByID(_ string, id core.Uuid, status connections.ConnectionStatus) (*connections.Connection, int, error) {
	return &connections.Connection{
		ID:     id,
		Status: status,
	}, 1, nil
}

// UpdateConnectionById updates and returns a mock connection for the specified ID string.
func (m *mockK8sConfigProvider) UpdateConnectionById(_ string, _ *connections.ConnectionPayload, connID string) (*connections.Connection, error) {
	id := uuid.FromStringOrNil(connID)
	return &connections.Connection{
		ID:     core.Uuid(id),
		Status: connections.DISCOVERED,
	}, nil
}

// UpdateConnection satisfies the Provider interface for connection updates.
func (m *mockK8sConfigProvider) UpdateConnection(_ *http.Request, conn *connections.Connection) (*connections.Connection, error) {
	return conn, nil
}

// testMultiContextKubeconfig generates a synthetic kubeconfig with n distinct contexts.
func testMultiContextKubeconfig(n int) []byte {
	var clusters, contexts, users string
	for i := 0; i < n; i++ {
		clusters += fmt.Sprintf("- cluster:\n    server: https://127.0.0.1:%d\n  name: cluster-%d\n", 59900+i, i)
		contexts += fmt.Sprintf("- context:\n    cluster: cluster-%d\n    user: user-%d\n  name: ctx-%d\n", i, i, i)
		users += fmt.Sprintf("- name: user-%d\n  user:\n    token: token-%d\n", i, i)
	}
	return []byte(fmt.Sprintf(
		"apiVersion: v1\nkind: Config\ncurrent-context: ctx-0\nclusters:\n%scontexts:\n%susers:\n%s",
		clusters, contexts, users,
	))
}

// createK8sConfigUploadRequest builds a multipart HTTP request with the given kubeconfig and form parameters.
func createK8sConfigUploadRequest(t *testing.T, kubeconfig []byte, formParams map[string]string, user *models.User, systemID *core.Uuid, provider models.Provider) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("k8sfile", "kubeconfig.yaml")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := part.Write(kubeconfig); err != nil {
		t.Fatalf("failed to write kubeconfig: %v", err)
	}

	for k, v := range formParams {
		if err := writer.WriteField(k, v); err != nil {
			t.Fatalf("failed to write field %s: %v", k, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/system/kubernetes", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "test-user-token")
	ctx = context.WithValue(ctx, models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, systemID)
	ctx = context.WithValue(ctx, models.ProviderCtxKey, provider)
	return req.WithContext(ctx)
}

// TestAddK8SConfig_ContextSelection verifies that when a user uploads a multi-context
// kubeconfig, explicit context selection (via contextName fallback or selectedContexts array)
// imports ONLY the selected context and does NOT register unselected contexts.
func TestAddK8SConfig_ContextSelection(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	kubeconfig := testMultiContextKubeconfig(3)

	// Discover contexts using models helper to discover the generated context ID for ctx-1
	instanceID := core.Uuid(systemID)
	discovered := models.K8sContextsFromKubeconfigWithOptions(nil, user.ID.String(), nil, kubeconfig, &instanceID, map[string]interface{}{}, newTestLogger(t), true)
	if len(discovered) != 3 {
		t.Fatalf("expected 3 discovered contexts, got %d", len(discovered))
	}
	var ctx1ID string
	for _, ctx := range discovered {
		if ctx.Name == "ctx-1" {
			ctx1ID = ctx.ID
			break
		}
	}
	if ctx1ID == "" {
		t.Fatal("could not find discovered context ID for ctx-1")
	}

	tests := []struct {
		name          string
		formParams    map[string]string
		wantSaved     int
		wantSavedName string
	}{
		{
			name: "legacy contextName form parameter selects only the requested context",
			formParams: map[string]string{
				"contextName": "ctx-1",
			},
			wantSaved:     1,
			wantSavedName: "ctx-1",
		},
		{
			name: "selectedContexts array with context name selects only the requested context",
			formParams: map[string]string{
				"selectedContexts": `["ctx-1"]`,
			},
			wantSaved:     1,
			wantSavedName: "ctx-1",
		},
		{
			name: "selectedContexts array with context ID selects only the requested context (UI wizard style)",
			formParams: map[string]string{
				"selectedContexts": fmt.Sprintf(`[%q]`, ctx1ID),
			},
			wantSaved:     1,
			wantSavedName: "ctx-1",
		},
		{
			name: "both selectedContexts and contextName provided for same context registers only that context",
			formParams: map[string]string{
				"selectedContexts": `["ctx-1"]`,
				"contextName":      "ctx-1",
			},
			wantSaved:     1,
			wantSavedName: "ctx-1",
		},
		{
			name: "conflicting selectors: selectedContexts is authoritative over contextName fallback",
			formParams: map[string]string{
				"selectedContexts": `["ctx-1"]`,
				"contextName":      "ctx-0",
			},
			wantSaved:     1,
			wantSavedName: "ctx-1",
		},
		{
			name:          "omitted selection registers all contexts in kubeconfig (backward compatible default)",
			formParams:    map[string]string{},
			wantSaved:     3,
			wantSavedName: "",
		},
		{
			name: "nonexistent context selection registers zero contexts and does not fall back to all",
			formParams: map[string]string{
				"contextName": "ctx-nonexistent",
			},
			wantSaved:     0,
			wantSavedName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				config: &models.HandlerConfig{
					EventBroadcaster:  models.NewBroadcaster("test"),
					K8scontextChannel: models.NewContextHelper(),
				},
				log:      newTestLogger(t),
				SystemID: &systemID,
				ConnectionToStateMachineInstanceTracker: &machines.ConnectionToStateMachineInstanceTracker{
					ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
				},
			}
			provider := &mockK8sConfigProvider{
				DefaultLocalProvider: &models.DefaultLocalProvider{},
			}

			req := createK8sConfigUploadRequest(t, kubeconfig, tt.formParams, user, &instanceID, provider)
			rec := httptest.NewRecorder()

			h.K8SConfigHandler(rec, req, nil, user, provider)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected HTTP 200, got %d: %s", rec.Code, rec.Body.String())
			}

			var resp SaveK8sContextResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			totalReported := len(resp.RegisteredContexts) + len(resp.ConnectedContexts) + len(resp.IgnoredContexts)
			if len(provider.savedContexts) != tt.wantSaved {
				t.Fatalf("saved %d contexts in provider, want %d", len(provider.savedContexts), tt.wantSaved)
			}
			if totalReported != tt.wantSaved {
				t.Fatalf("response reported %d contexts, want %d", totalReported, tt.wantSaved)
			}

			if tt.wantSavedName != "" {
				if provider.savedContexts[0].Name != tt.wantSavedName {
					t.Errorf("saved context name = %q, want %q", provider.savedContexts[0].Name, tt.wantSavedName)
				}
			}
		})
	}
}
