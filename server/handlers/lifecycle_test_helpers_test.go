package handlers

import (
	"net/http"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
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
	m.mu.Lock()
	if conn != nil {
		m.connection = &connections.Connection{
			ID:     conn.ID,
			Name:   conn.Name,
			Status: conn.Status,
		}
	}
	m.mu.Unlock()

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
