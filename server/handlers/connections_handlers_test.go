package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	v1beta3comp "github.com/meshery/schemas/models/v1beta3/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- test helpers ---

func testSystemID(t *testing.T) *core.Uuid {
	t.Helper()
	id := core.Uuid(uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")))
	return &id
}

func newConnectionsHandler(t *testing.T, rm *registry.RegistryManager) *Handler {
	t.Helper()
	sysID := testSystemID(t)
	h := newTestHandler(t, map[string]models.Provider{}, "")
	h.SystemID = sysID
	h.registryManager = rm
	h.ConnectionToStateMachineInstanceTracker = &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
	h.config.EventBroadcaster = models.NewBroadcaster("test-events")
	h.config.OperatorTracker = models.NewOperatorTracker(true)
	h.MeshsyncDefaultDeploymentMode = connections.MeshsyncDeploymentModeEmbedded
	return h
}

func newLocalProviderWithConnectionsDB(t *testing.T) (*models.DefaultLocalProvider, *database.Handler) {
	t.Helper()
	db, err := database.New(database.Options{Filename: ":memory:", Engine: "sqlite"})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&connections.Connection{}, &events.Event{}))

	p := &models.DefaultLocalProvider{}
	p.Initialize()
	p.ConnectionPersister = &models.ConnectionPersister{DB: &db}
	p.EventsPersister = &models.EventsPersister{DB: &db}
	p.GenericPersister = &db
	return p, &db
}

func seedRegistryConnectionComponent(t *testing.T, rm *registry.RegistryManager, componentKind string) {
	t.Helper()
	seedRegistryComponentWithKind(t, rm, componentKind)
}

func seedRegistryComponentWithKind(t *testing.T, rm *registry.RegistryManager, componentKind string) {
	t.Helper()
	conn := connection.Connection{
		Name:    "test-registrant",
		Kind:    "kubernetes",
		Type:    "platform",
		SubType: "orchestration",
		Status:  connection.ConnectionStatusConnected,
	}
	enabled := v1beta3comp.Enabled
	comp := v1beta3comp.ComponentDefinition{
		DisplayName:   componentKind,
		SchemaVersion: "core.meshery.io/v1beta1",
		Status:        &enabled,
		Component: v1beta3comp.Component{
			Kind:    componentKind,
			Version: "v1",
			Schema:  `{"properties":{}}`,
		},
		Model: &model.ModelDefinition{
			Name:          "kubernetes",
			DisplayName:   "Kubernetes",
			SchemaVersion: "models.meshery.io/v1beta1",
			Version:       "v1.25.0",
			Model:         model.Model{Version: "v1.25.0"},
			Category:      category.CategoryDefinition{Name: "Orchestration"},
			Status:        model.Enabled,
		},
	}
	id, err := comp.GenerateID()
	require.NoError(t, err)
	comp.ID = id
	_, _, err = rm.RegisterEntity(conn, &comp)
	require.NoError(t, err)
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, fmt.Errorf("forced read error") }
func (errReader) Close() error             { return nil }

type failWriteResponseWriter struct {
	httptest.ResponseRecorder
}

func (w *failWriteResponseWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("forced write error")
}

type connectionSpyProvider struct {
	*models.DefaultLocalProvider
	saveConnectionErr   error
	updateConnectionErr error
	deleteConnectionErr error
	getConnectionByIDFn func(token string, id core.Uuid) (*connections.Connection, int, error)
	getK8sContextFn     func(token, id string) (models.K8sContext, error)
	getProviderTokenFn  func(*http.Request) (string, error)
	persistEventFn      func(events.Event, string) error
	persistEventCalls   atomic.Int32
	lastPersistedEvent  atomic.Pointer[events.Event]
}

func newConnectionSpyProvider(t *testing.T) *connectionSpyProvider {
	t.Helper()
	base, _ := newLocalProviderWithConnectionsDB(t)
	return &connectionSpyProvider{DefaultLocalProvider: base}
}

func (p *connectionSpyProvider) SaveConnection(conn *connections.ConnectionPayload, token string, skip bool) (*connections.Connection, error) {
	if p.saveConnectionErr != nil {
		return nil, p.saveConnectionErr
	}
	return p.DefaultLocalProvider.SaveConnection(conn, token, skip)
}

func (p *connectionSpyProvider) UpdateConnectionById(token string, conn *connections.ConnectionPayload, connID string) (*connections.Connection, error) {
	if p.updateConnectionErr != nil {
		return nil, p.updateConnectionErr
	}
	return p.DefaultLocalProvider.UpdateConnectionById(token, conn, connID)
}

func (p *connectionSpyProvider) DeleteConnection(req *http.Request, id core.Uuid) (*connections.Connection, error) {
	if p.deleteConnectionErr != nil {
		return nil, p.deleteConnectionErr
	}
	return p.DefaultLocalProvider.DeleteConnection(req, id)
}

func (p *connectionSpyProvider) GetConnectionByID(token string, id core.Uuid) (*connections.Connection, int, error) {
	if p.getConnectionByIDFn != nil {
		return p.getConnectionByIDFn(token, id)
	}
	return p.DefaultLocalProvider.GetConnectionByID(token, id)
}

func (p *connectionSpyProvider) GetK8sContext(token, id string) (models.K8sContext, error) {
	if p.getK8sContextFn != nil {
		return p.getK8sContextFn(token, id)
	}
	return p.DefaultLocalProvider.GetK8sContext(token, id)
}

func (p *connectionSpyProvider) GetProviderToken(req *http.Request) (string, error) {
	if p.getProviderTokenFn != nil {
		return p.getProviderTokenFn(req)
	}
	return p.DefaultLocalProvider.GetProviderToken(req)
}

func (p *connectionSpyProvider) PersistEvent(event events.Event, token string) error {
	p.persistEventCalls.Add(1)
	cp := event
	p.lastPersistedEvent.Store(&cp)
	if p.persistEventFn != nil {
		return p.persistEventFn(event, token)
	}
	return p.DefaultLocalProvider.PersistEvent(event, token)
}

func seedConnectionInDB(t *testing.T, p *connectionSpyProvider, conn *connections.Connection) {
	t.Helper()
	saved, err := p.ConnectionPersister.SaveConnection(conn)
	require.NoError(t, err)
	*conn = *saved
}

func newTestMesheryCtrlsHelper(h *Handler, p models.Provider, db *database.Handler) *models.MesheryControllersHelper {
	return models.NewMesheryControllersHelper(
		h.log,
		models.NewOperatorDeploymentConfig(nil),
		db,
		h.config.EventBroadcaster,
		p,
		h.SystemID,
	)
}

type noopMachineAction struct{}

func (noopMachineAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
func (noopMachineAction) Execute(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
func (noopMachineAction) ExecuteOnExit(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func seedNoopMachineWithCtx(t *testing.T, h *Handler, p models.Provider, db *database.Handler, connID core.Uuid, currentState machines.StateType, k8sCtx models.K8sContext) *machines.StateMachine {
	t.Helper()
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	sm, err := kubernetes.New(connID.String(), userID, h.log)
	require.NoError(t, err)
	sm.Provider = p
	sm.CurrentState = currentState
	sm.Context = &kubernetes.MachineCtx{
		K8sContext:         k8sCtx,
		MesheryCtrlsHelper: newTestMesheryCtrlsHelper(h, p, db),
		OperatorTracker:    h.config.OperatorTracker,
	}
	noop := noopMachineAction{}
	for stateType, state := range sm.States {
		state.Action = noop
		sm.States[stateType] = state
	}
	h.ConnectionToStateMachineInstanceTracker.Add(connID, sm)
	return sm
}

func seedMachineWithCtx(t *testing.T, h *Handler, connID core.Uuid, ctrlHelper *models.MesheryControllersHelper, k8sCtx models.K8sContext) *machines.StateMachine {
	t.Helper()
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	sm, err := kubernetes.New(connID.String(), userID, h.log)
	require.NoError(t, err)
	sm.Provider = h.config.Providers["Local"]
	sm.Context = &kubernetes.MachineCtx{
		K8sContext:         k8sCtx,
		MesheryCtrlsHelper: ctrlHelper,
		OperatorTracker:    h.config.OperatorTracker,
	}
	h.ConnectionToStateMachineInstanceTracker.Add(connID, sm)
	return sm
}

// --- Phase 1: handleRegistrationInitEvent ---

func TestHandleRegistrationInitEvent_UnknownKind(t *testing.T) {
	rm, _ := newTestRegistryManager(t)
	h := newConnectionsHandler(t, rm)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", nil)
	rec := httptest.NewRecorder()
	payload := &connections.ConnectionPayload{Kind: "unknown-kind"}

	h.handleRegistrationInitEvent(rec, req, payload)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrUnknownConnectionKindCode)
}

func TestHandleRegistrationInitEvent_SuccessWithoutCredential(t *testing.T) {
	rm, _ := newTestRegistryManager(t)
	seedRegistryConnectionComponent(t, rm, "grafanaConnection")
	h := newConnectionsHandler(t, rm)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", nil)
	rec := httptest.NewRecorder()
	payload := &connections.ConnectionPayload{Kind: "grafana"}

	h.handleRegistrationInitEvent(rec, req, payload)

	require.Equal(t, http.StatusOK, rec.Code)
	var schema map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &schema))
	assert.NotNil(t, schema["connection"])
	assert.NotNil(t, schema["id"])
	assert.Nil(t, schema["credential"])
}

func TestHandleRegistrationInitEvent_SuccessWithCredential(t *testing.T) {
	rm, _ := newTestRegistryManager(t)
	seedRegistryConnectionComponent(t, rm, "prometheusConnection")
	seedRegistryConnectionComponent(t, rm, "prometheusCredential")
	h := newConnectionsHandler(t, rm)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", nil)
	rec := httptest.NewRecorder()
	payload := &connections.ConnectionPayload{Kind: "prometheus"}

	h.handleRegistrationInitEvent(rec, req, payload)

	require.Equal(t, http.StatusOK, rec.Code)
	var schema map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &schema))
	assert.NotNil(t, schema["connection"])
	assert.NotNil(t, schema["credential"])
	assert.NotNil(t, schema["id"])
}

func TestHandleRegistrationInitEvent_WriteResponseError(t *testing.T) {
	rm, _ := newTestRegistryManager(t)
	seedRegistryConnectionComponent(t, rm, "grafanaConnection")
	h := newConnectionsHandler(t, rm)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/register", nil)
	w := &failWriteResponseWriter{}
	payload := &connections.ConnectionPayload{Kind: "grafana"}

	h.handleRegistrationInitEvent(w, req, payload)
}

// --- Phase 2: handleMeshSyncDeploymentModeChange ---

func TestHandleMeshSyncDeploymentModeChange_NilNewConnection(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)

	_, _, changed, err := h.handleMeshSyncDeploymentModeChange(
		context.Background(), core.Uuid(uuid.Nil), nil, "token", core.Uuid(uuid.Nil), p,
	)
	require.Error(t, err)
	assert.False(t, changed)
}

func TestHandleMeshSyncDeploymentModeChange_NilSystemID(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	h.SystemID = nil
	p := newConnectionSpyProvider(t)

	_, _, changed, err := h.handleMeshSyncDeploymentModeChange(
		context.Background(), core.Uuid(uuid.Nil), &connections.ConnectionPayload{}, "token", core.Uuid(uuid.Nil), p,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "system ID is not configured")
	assert.False(t, changed)
}

func TestHandleMeshSyncDeploymentModeChange_GetConnectionError(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	p.getConnectionByIDFn = func(string, core.Uuid) (*connections.Connection, int, error) {
		return nil, http.StatusInternalServerError, fmt.Errorf("db down")
	}

	_, _, changed, err := h.handleMeshSyncDeploymentModeChange(
		context.Background(), connID, &connections.ConnectionPayload{}, "token", core.Uuid(uuid.Nil), p,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to retrieve existing connection")
	assert.False(t, changed)
}

func TestHandleMeshSyncDeploymentModeChange_NilExistingConnection(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	p.getConnectionByIDFn = func(string, core.Uuid) (*connections.Connection, int, error) {
		return nil, http.StatusOK, nil
	}

	_, _, changed, err := h.handleMeshSyncDeploymentModeChange(
		context.Background(), connID, &connections.ConnectionPayload{}, "token", core.Uuid(uuid.Nil), p,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "existing connection is nil")
	assert.False(t, changed)
}

func TestHandleMeshSyncDeploymentModeChange_NotKubernetesKind(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	p.getConnectionByIDFn = func(string, core.Uuid) (*connections.Connection, int, error) {
		return &connections.Connection{ID: connID, Kind: "grafana"}, http.StatusOK, nil
	}

	_, _, changed, err := h.handleMeshSyncDeploymentModeChange(
		context.Background(), connID, &connections.ConnectionPayload{}, "token", core.Uuid(uuid.Nil), p,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not of kind kubernetes")
	assert.False(t, changed)
}

func TestHandleMeshSyncDeploymentModeChange_NoModeChange(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	meta := core.Map{connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeOperator)}
	p.getConnectionByIDFn = func(string, core.Uuid) (*connections.Connection, int, error) {
		return &connections.Connection{ID: connID, Kind: "kubernetes", Metadata: meta}, http.StatusOK, nil
	}
	newConn := &connections.ConnectionPayload{MetaData: meta}

	old, newMode, changed, err := h.handleMeshSyncDeploymentModeChange(
		context.Background(), connID, newConn, "token", core.Uuid(uuid.Nil), p,
	)
	require.NoError(t, err)
	assert.False(t, changed)
	assert.Equal(t, connections.MeshsyncDeploymentModeOperator, old)
	assert.Equal(t, connections.MeshsyncDeploymentModeOperator, newMode)
}

func TestHandleMeshSyncDeploymentModeChange_UndefinedNewModeUsesDefault(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	h.MeshsyncDefaultDeploymentMode = connections.MeshsyncDeploymentModeEmbedded
	p := newConnectionSpyProvider(t)
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	existingMeta := core.Map{connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeOperator)}
	p.getConnectionByIDFn = func(string, core.Uuid) (*connections.Connection, int, error) {
		return &connections.Connection{ID: connID, Kind: "kubernetes", Metadata: existingMeta}, http.StatusOK, nil
	}

	ctrlHelper := newTestMesheryCtrlsHelper(h, p, p.GenericPersister)
	k8sCtx := models.K8sContext{ID: "ctx-1", Name: "test", ConnectionID: connID.String()}
	seedMachineWithCtx(t, h, connID, ctrlHelper, k8sCtx)

	old, newMode, changed, err := h.handleMeshSyncDeploymentModeChange(
		context.Background(), connID, &connections.ConnectionPayload{MetaData: core.Map{}}, "token", core.Uuid(uuid.Nil), p,
	)
	require.NoError(t, err)
	assert.True(t, changed)
	assert.Equal(t, connections.MeshsyncDeploymentModeOperator, old)
	assert.Equal(t, connections.MeshsyncDeploymentModeEmbedded, newMode)
}

func TestHandleMeshSyncDeploymentModeChange_NilInstanceTracker(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	h.ConnectionToStateMachineInstanceTracker = nil
	p := newConnectionSpyProvider(t)
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	existingMeta := core.Map{connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeOperator)}
	newMeta := core.Map{connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeEmbedded)}
	p.getConnectionByIDFn = func(string, core.Uuid) (*connections.Connection, int, error) {
		return &connections.Connection{ID: connID, Kind: "kubernetes", Metadata: existingMeta}, http.StatusOK, nil
	}

	_, _, changed, err := h.handleMeshSyncDeploymentModeChange(
		context.Background(), connID, &connections.ConnectionPayload{MetaData: newMeta}, "token", core.Uuid(uuid.Nil), p,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "instance tracker is nil")
	assert.False(t, changed)
}

func TestHandleMeshSyncDeploymentModeChange_MachineMissingInTracker(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	existingMeta := core.Map{connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeOperator)}
	newMeta := core.Map{connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeEmbedded)}
	p.getConnectionByIDFn = func(string, core.Uuid) (*connections.Connection, int, error) {
		return &connections.Connection{ID: connID, Kind: "kubernetes", Metadata: existingMeta}, http.StatusOK, nil
	}

	_, _, changed, err := h.handleMeshSyncDeploymentModeChange(
		context.Background(), connID, &connections.ConnectionPayload{MetaData: newMeta}, "token", core.Uuid(uuid.Nil), p,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not contain machine")
	assert.False(t, changed)
}

func TestHandleMeshSyncDeploymentModeChange_InvalidMachineContext(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	existingMeta := core.Map{connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeOperator)}
	newMeta := core.Map{connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeEmbedded)}
	p.getConnectionByIDFn = func(string, core.Uuid) (*connections.Connection, int, error) {
		return &connections.Connection{ID: connID, Kind: "kubernetes", Metadata: existingMeta}, http.StatusOK, nil
	}
	sm, err := kubernetes.New(connID.String(), core.Uuid(uuid.Must(uuid.NewV4())), h.log)
	require.NoError(t, err)
	sm.Context = "not-a-machine-ctx"
	h.ConnectionToStateMachineInstanceTracker.Add(connID, sm)

	_, _, changed, err := h.handleMeshSyncDeploymentModeChange(
		context.Background(), connID, &connections.ConnectionPayload{MetaData: newMeta}, "token", core.Uuid(uuid.Nil), p,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get machine context")
	assert.False(t, changed)
}

func TestHandleMeshSyncDeploymentModeChange_NilMesheryCtrlsHelper(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	existingMeta := core.Map{connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeOperator)}
	newMeta := core.Map{connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeEmbedded)}
	p.getConnectionByIDFn = func(string, core.Uuid) (*connections.Connection, int, error) {
		return &connections.Connection{ID: connID, Kind: "kubernetes", Metadata: existingMeta}, http.StatusOK, nil
	}
	sm, err := kubernetes.New(connID.String(), core.Uuid(uuid.Must(uuid.NewV4())), h.log)
	require.NoError(t, err)
	sm.Context = &kubernetes.MachineCtx{K8sContext: models.K8sContext{ID: "ctx-1", ConnectionID: connID.String()}}
	h.ConnectionToStateMachineInstanceTracker.Add(connID, sm)

	_, _, changed, err := h.handleMeshSyncDeploymentModeChange(
		context.Background(), connID, &connections.ConnectionPayload{MetaData: newMeta}, "token", core.Uuid(uuid.Nil), p,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MesheryCtrlsHelper")
	assert.False(t, changed)
}

// --- Phase 3: NotifySmOfConnectionStatusChange ---

func TestNotifySmOfConnectionStatusChange_EmptyStatus(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	conn := &connections.ConnectionPayload{ID: connID}

	event, err := h.NotifySmOfConnectionStatusChange(context.Background(), userID, p, "token", conn)
	require.NoError(t, err)
	assert.Equal(t, connID, event.ActedUpon)
}

func TestNotifySmOfConnectionStatusChange_GetK8sContextError(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	conn := &connections.ConnectionPayload{ID: connID, Status: connections.CONNECTED}
	p.getK8sContextFn = func(string, string) (models.K8sContext, error) {
		return models.K8sContext{}, fmt.Errorf("context not found")
	}

	event, err := h.NotifySmOfConnectionStatusChange(context.Background(), userID, p, "token", conn)
	require.Error(t, err)
	assert.Equal(t, events.Error, event.Severity)
}

func TestNotifySmOfConnectionStatusChange_InitializeMachineError(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	h.MesheryCtrlsHelper = newTestMesheryCtrlsHelper(h, p, p.GenericPersister)
	h.K8sCompRegHelper = &models.ComponentsRegistrationHelper{}
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	conn := &connections.ConnectionPayload{ID: connID, Status: connections.CONNECTED}
	p.getK8sContextFn = func(string, string) (models.K8sContext, error) {
		return models.K8sContext{
			Name:         "unreachable",
			Server:       "https://127.0.0.1:1",
			ConnectionID: connID.String(),
		}, nil
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userID})
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "token")

	event, err := h.NotifySmOfConnectionStatusChange(ctx, userID, p, "token", conn)
	require.Error(t, err)
	assert.Equal(t, events.Error, event.Severity)
}

func TestNotifySmOfConnectionStatusChange_HappyPathWithExistingMachine(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	conn := &connections.ConnectionPayload{ID: connID, Status: connections.DISCOVERED, Kind: "kubernetes"}

	seedConnectionInDB(t, p, &connections.Connection{
		ID:     connID,
		Name:   "test-conn",
		Kind:   "kubernetes",
		Status: connections.DISCOVERED,
	})

	k8sCtx := models.K8sContext{ID: "ctx-1", Name: "cluster-a", ConnectionID: connID.String()}
	seedNoopMachineWithCtx(t, h, p, p.GenericPersister, connID, machines.InitialState, k8sCtx)

	p.getK8sContextFn = func(string, string) (models.K8sContext, error) { return k8sCtx, nil }

	done := make(chan struct{}, 1)
	p.persistEventFn = func(events.Event, string) error {
		select {
		case done <- struct{}{}:
		default:
		}
		return nil
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userID})
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "token")

	event, err := h.NotifySmOfConnectionStatusChange(ctx, userID, p, "token", conn)
	require.NoError(t, err)
	assert.Equal(t, events.Informational, event.Severity)

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for async PersistEvent from state machine goroutine")
	}
}

func TestNotifySmOfConnectionStatusChange_DeletedStatusRemovesTracker(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	conn := &connections.ConnectionPayload{ID: connID, Status: connections.DELETED, Kind: "kubernetes"}

	seedConnectionInDB(t, p, &connections.Connection{
		ID:     connID,
		Name:   "delete-me",
		Kind:   "kubernetes",
		Status: connections.CONNECTED,
	})

	k8sCtx := models.K8sContext{ID: "ctx-2", Name: "cluster-b", ConnectionID: connID.String()}
	seedNoopMachineWithCtx(t, h, p, p.GenericPersister, connID, machines.CONNECTED, k8sCtx)
	p.getK8sContextFn = func(string, string) (models.K8sContext, error) { return k8sCtx, nil }

	var wg sync.WaitGroup
	wg.Add(1)
	p.persistEventFn = func(events.Event, string) error {
		wg.Done()
		return nil
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userID})
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "token")

	_, err := h.NotifySmOfConnectionStatusChange(ctx, userID, p, "token", conn)
	require.NoError(t, err)

	waitDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitDone)
	}()
	select {
	case <-waitDone:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for delete goroutine")
	}

	_, ok := h.ConnectionToStateMachineInstanceTracker.Get(connID)
	assert.False(t, ok, "machine should be removed from tracker after DELETED status")
}

// --- SaveConnection ---

func TestSaveConnection_ReadBodyError(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections", errReader{})
	rec := httptest.NewRecorder()
	h.SaveConnection(rec, req, nil, user, p)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSaveConnection_InvalidJSON(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections", bytes.NewBufferString("{bad"))
	rec := httptest.NewRecorder()
	h.SaveConnection(rec, req, nil, user, p)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSaveConnection_ProviderError(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	p.saveConnectionErr = fmt.Errorf("save failed")
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}

	body := `{"name":"conn-a","kind":"grafana","type":"platform","status":"discovered"}`
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections", bytes.NewBufferString(body))
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "token"))
	rec := httptest.NewRecorder()
	h.SaveConnection(rec, req, nil, user, p)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.GreaterOrEqual(t, p.persistEventCalls.Load(), int32(1))
}

func TestSaveConnection_SuccessWithContextToken(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}

	body := `{"name":"conn-b","kind":"grafana","type":"platform","status":"discovered"}`
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections", bytes.NewBufferString(body))
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "ctx-token"))
	rec := httptest.NewRecorder()
	h.SaveConnection(rec, req, nil, user, p)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.GreaterOrEqual(t, p.persistEventCalls.Load(), int32(1))
}

func TestSaveConnection_SuccessWithCookieToken(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}

	body := `{"name":"conn-c","kind":"grafana","type":"platform","status":"discovered"}`
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections", bytes.NewBufferString(body))
	req.AddCookie(&http.Cookie{Name: models.TokenCookieName, Value: "cookie-token"})
	rec := httptest.NewRecorder()
	h.SaveConnection(rec, req, nil, user, p)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

// --- UpdateConnectionById ---

func TestUpdateConnectionById_ReadBodyError(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	connID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), errReader{})
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()
	h.UpdateConnectionById(rec, req, nil, user, p)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateConnectionById_InvalidJSON(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	connID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), bytes.NewBufferString("not-json"))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()
	h.UpdateConnectionById(rec, req, nil, user, p)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateConnectionById_MeshsyncModeChangeError(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	connID := uuid.Must(uuid.NewV4())
	p.getConnectionByIDFn = func(string, core.Uuid) (*connections.Connection, int, error) {
		return nil, http.StatusInternalServerError, fmt.Errorf("lookup failed")
	}

	body := fmt.Sprintf(`{"id":"%s","kind":"kubernetes","metadata":{"%s":"embedded"}}`, connID, connections.MeshsyncDeploymentModeMetadataKey)
	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "token"))
	rec := httptest.NewRecorder()
	h.UpdateConnectionById(rec, req, nil, user, p)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateConnectionById_MeshsyncModeChangedEvent(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	existingMeta := core.Map{connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeOperator)}
	p.getConnectionByIDFn = func(string, core.Uuid) (*connections.Connection, int, error) {
		return &connections.Connection{ID: connID, Kind: "kubernetes", Name: "k8s", Metadata: existingMeta}, http.StatusOK, nil
	}
	ctrlHelper := newTestMesheryCtrlsHelper(h, p, p.GenericPersister)
	k8sCtx := models.K8sContext{ID: "ctx-3", Name: "cluster-c", ConnectionID: connID.String()}
	seedMachineWithCtx(t, h, connID, ctrlHelper, k8sCtx)

	body := fmt.Sprintf(`{"id":"%s","name":"k8s","kind":"kubernetes","metadata":{"%s":"embedded"}}`, connID, connections.MeshsyncDeploymentModeMetadataKey)
	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "token"))
	rec := httptest.NewRecorder()
	h.UpdateConnectionById(rec, req, nil, user, p)

        assert.NotEqual(t, 0, rec.Code)
	assert.GreaterOrEqual(t, p.persistEventCalls.Load(), int32(2))
}

func TestUpdateConnectionById_GetProviderTokenError(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	connID := uuid.Must(uuid.NewV4())
	p.getProviderTokenFn = func(*http.Request) (string, error) {
		return "", fmt.Errorf("no token")
	}

	body := fmt.Sprintf(`{"id":"%s","name":"k8s","kind":"kubernetes"}`, connID)
	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()
	h.UpdateConnectionById(rec, req, nil, user, p)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateConnectionById_UpdateError(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	p.updateConnectionErr = fmt.Errorf("update failed")
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	connID := uuid.Must(uuid.NewV4())

	body := fmt.Sprintf(`{"id":"%s","name":"k8s","kind":"kubernetes"}`, connID)
	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()
	h.UpdateConnectionById(rec, req, nil, user, p)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateConnectionById_SuccessWithStatusNotify(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	connID := core.Uuid(uuid.Must(uuid.NewV4()))

	seedConnectionInDB(t, p, &connections.Connection{
		ID: connID, Name: "k8s", Kind: "kubernetes", Status: connections.DISCOVERED,
	})
	k8sCtx := models.K8sContext{ID: "ctx-4", Name: "cluster-d", ConnectionID: connID.String()}
	seedNoopMachineWithCtx(t, h, p, p.GenericPersister, connID, machines.InitialState, k8sCtx)
	p.getK8sContextFn = func(string, string) (models.K8sContext, error) { return k8sCtx, nil }

	body := fmt.Sprintf(`{"id":"%s","name":"k8s","kind":"kubernetes","status":"discovered"}`, connID)
	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "token")
	ctx = context.WithValue(ctx, models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h.UpdateConnectionById(rec, req, nil, user, p)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- DeleteConnection ---

func TestDeleteConnection_GetProviderTokenError(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	p.getProviderTokenFn = func(*http.Request) (string, error) {
		return "", fmt.Errorf("missing token")
	}
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	connID := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodDelete, "/api/integrations/connections/"+connID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()
	h.DeleteConnection(rec, req, nil, user, p)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteConnection_NotFound(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	connID := uuid.Must(uuid.NewV4())
	p.deleteConnectionErr = models.ErrResultNotFound(fmt.Errorf("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/api/integrations/connections/"+connID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()
	h.DeleteConnection(rec, req, nil, user, p)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteConnection_ProviderError(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	connID := uuid.Must(uuid.NewV4())
	p.deleteConnectionErr = fmt.Errorf("db error")

	req := httptest.NewRequest(http.MethodDelete, "/api/integrations/connections/"+connID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()
	h.DeleteConnection(rec, req, nil, user, p)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteConnection_Success(t *testing.T) {
	h := newConnectionsHandler(t, nil)
	p := newConnectionSpyProvider(t)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	seedConnectionInDB(t, p, &connections.Connection{
		ID: connID, Name: "gone", Kind: "grafana", Status: connections.DISCOVERED,
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/integrations/connections/"+connID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()
	h.DeleteConnection(rec, req, nil, user, p)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.GreaterOrEqual(t, p.persistEventCalls.Load(), int32(1))
}
