package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type mockMiddlewareProvider struct {
	models.Provider
	conn   *connections.Connection
	status connections.ConnectionStatus
}

func (m *mockMiddlewareProvider) LoadAllK8sContext(token string) ([]*models.K8sContext, error) {
	if m.status == connections.DISCONNECTED {
		return []*models.K8sContext{
			{
				ID:           "test-ctx-id",
				Name:         "test-cluster",
				Server:       "http://localhost",
				ConnectionID: m.conn.ID.String(),
			},
		}, nil
	}
	return []*models.K8sContext{}, nil
}

func (m *mockMiddlewareProvider) SaveK8sContext(token string, k8sContext models.K8sContext, additionalMetadata map[string]any) (connections.Connection, error) {
	return *m.conn, nil
}

func (m *mockMiddlewareProvider) GetConnectionByID(token string, connectionID uuid.UUID) (*connections.Connection, int, error) {
	return m.conn, 200, nil
}

func (m *mockMiddlewareProvider) UpdateConnectionById(token string, conn *connections.ConnectionPayload, connId string) (*connections.Connection, error) {
	m.conn.Status = conn.Status
	m.status = conn.Status
	return m.conn, nil
}

func (m *mockMiddlewareProvider) PersistEvent(e events.Event, token string) error {
	return nil
}

func (m *mockMiddlewareProvider) GetGenericPersister() *database.Handler {
	return nil
}

// TestKubernetesMiddleware_SecondInvocationAfter403 verifies that a 403 Forbidden on the first discovery
// correctly transitions the connection to DISCONNECTED, and that a subsequent middleware execution
// does NOT cause an additional Kubernetes API request.
func TestKubernetesMiddleware_SecondInvocationAfter403(t *testing.T) {
	var requestCount int32

	// Create a channel to dictate what status code the mock server returns.
	// This satisfies CodeRabbit's feedback to not use requestCount for deciding 403,
	// while still allowing the first discovery request to succeed (200 OK) so the FSM is spawned,
	// and the subsequent FSM request to fail with 403.
	statusResponses := make(chan int, 2)
	statusResponses <- http.StatusOK
	statusResponses <- http.StatusForbidden

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/namespaces/kube-system" {
			status := <-statusResponses
			w.WriteHeader(status)
			if status == http.StatusForbidden {
				_, _ = w.Write([]byte(`{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure","message":"forbidden","reason":"Forbidden","code":403}`))
				return
			}
			_, _ = w.Write([]byte(`{"metadata":{"uid":"test-uid"}}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	instID := core.Uuid(uuid.Must(uuid.NewV4()))
	viper.Set("INSTANCE_ID", &instID)

	kubeConfigDir := t.TempDir()
	configContent := fmt.Sprintf(`
apiVersion: v1
clusters:
- cluster:
    server: %s
    insecure-skip-tls-verify: true
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
kind: Config
preferences: {}
users:
- name: test-user
  user:
    token: test-token
`, server.URL)

	err := os.WriteFile(filepath.Join(kubeConfigDir, "config"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("failed to write mock kubeconfig: %v", err)
	}

	connUUID := uuid.Must(uuid.NewV4())
	provider := &mockMiddlewareProvider{
		conn: &connections.Connection{
			ID:     connUUID,
			Kind:   "kubernetes",
			Status: connections.DISCOVERED,
		},
		status: connections.DISCOVERED,
	}

	log, err := logger.New("test", logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: int(logrus.DebugLevel),
	})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}
	handler := &Handler{
		log: log,
		config: &models.HandlerConfig{
			KubeConfigFolder:  kubeConfigDir,
			EventBroadcaster:  models.NewBroadcaster("test"),
			OperatorTracker:   models.NewOperatorTracker(false),
			K8scontextChannel: models.NewContextHelper(),
		},
		ConnectionToStateMachineInstanceTracker: &machines.ConnectionToStateMachineInstanceTracker{
			ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
		},
		MesheryCtrlsHelper: &models.MesheryControllersHelper{},
	}

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "mock-token")

	// Phase 1: First invocation
	// The middleware will execute, FSM transitions to DISCOVERED, calls AssignServerID -> 403 Forbidden
	// FSM will transition to DISCONNECTED.
	ctx = context.WithValue(ctx, models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.ProviderCtxKey, provider)
	ctx = context.WithValue(ctx, models.SystemIDKey, &instID)

	_, err = KubernetesMiddleware(ctx, handler, provider, user, []string{"all"})
	if err != nil {
		t.Fatalf("first middleware invocation failed: %v", err)
	}

	// Wait for the asynchronous FSM to complete its execution
	time.Sleep(200 * time.Millisecond)

	countAfterFirst := atomic.LoadInt32(&requestCount)
	if countAfterFirst == 0 {
		t.Fatalf("expected first invocation to make a Kubernetes API request, got 0")
	}

	if provider.status != connections.DISCONNECTED {
		t.Fatalf("expected connection status to transition to DISCONNECTED after 403, got %v", provider.status)
	}

	// Phase 2: Second invocation
	// Due to DISCONNECTED status, middleware should bypass discovery and NOT hit the K8s API.
	_, err = KubernetesMiddleware(ctx, handler, provider, user, []string{"all"})
	if err != nil {
		t.Fatalf("second middleware invocation failed: %v", err)
	}

	// Wait briefly to ensure no asynchronous discovery fires
	time.Sleep(100 * time.Millisecond)

	countAfterSecond := atomic.LoadInt32(&requestCount)
	if countAfterSecond > countAfterFirst {
		t.Fatalf("second invocation made %d additional Kubernetes API requests (expected 0). Total: %d", countAfterSecond-countAfterFirst, countAfterSecond)
	}
}
