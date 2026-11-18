package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	k8smachine "github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshkit/database"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

func newConnectionActionTestHandler(t *testing.T) *Handler {
	t.Helper()
	return &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}
}

func TestPerformConnectionAction_InvalidConnectionID(t *testing.T) {
	h := newConnectionActionTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/not-a-valid-uuid/actions",
		strings.NewReader(`{"action":"setMeshsyncMode","mode":"operator"}`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": "not-a-valid-uuid"})
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for invalid connection id; body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
		t.Errorf("expected error body to contain %s, got: %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
	}
}

func TestPerformConnectionAction_NilConnectionID(t *testing.T) {
	h := newConnectionActionTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/"+uuid.Nil.String()+"/actions",
		strings.NewReader(`{"action":"setMeshsyncMode","mode":"operator"}`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": uuid.Nil.String()})
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for nil connection id; body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
		t.Errorf("expected error body to contain %s, got: %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
	}
}

func TestPerformConnectionAction_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newConnectionActionTestHandler(t)

	connID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/"+connID.String()+"/actions",
		strings.NewReader(`{"action":"setMeshsyncMode","mode":"operator"}`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500 for nil provider; body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got: %s", rec.Body.String())
	}
}

func TestPerformConnectionAction_FlushMeshsync_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newConnectionActionTestHandler(t)

	connID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/"+connID.String()+"/actions",
		strings.NewReader(`{"action":"flushMeshsync"}`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500 for nil provider; body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got: %s", rec.Body.String())
	}
}

func TestPerformConnectionAction_UnsupportedAction(t *testing.T) {
	h := newConnectionActionTestHandler(t)

	connID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/"+connID.String()+"/actions",
		strings.NewReader(`{"action":"noSuchAction"}`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for unsupported action; body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "unsupported connection action") {
		t.Fatalf("expected unsupported-action message, got: %s", rec.Body.String())
	}
}

func TestPerformConnectionAction_MalformedBody(t *testing.T) {
	h := newConnectionActionTestHandler(t)

	connID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/"+connID.String()+"/actions",
		strings.NewReader(`{not json`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for malformed body; body: %s", rec.Code, rec.Body.String())
	}
}

// lockedEventsProvider records persisted events under a mutex: the mode-change
// handler persists from two goroutines, and the controllers-config fake's
// unsynchronized slice would otherwise be reported by the race detector.
type lockedEventsProvider struct {
	*controllersConfigFakeProvider
	mu        sync.Mutex
	persisted []events.Event
}

func (p *lockedEventsProvider) PersistEvent(event events.Event, _ string) (models.PersistedEvent, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.persisted = append(p.persisted, event)
	return models.PersistedEvent{Event: &event, Disposition: models.EventCreated}, nil
}

func (p *lockedEventsProvider) snapshot() []events.Event {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]events.Event(nil), p.persisted...)
}

// A mode change whose redeploy fails (the connection has no state machine)
// emits a success event from the handler and an error event from the redeploy
// goroutine. They must be independent: before #22729 both were built from one
// EventBuilder, so both broadcasts carried the same *events.Event and the error
// overwrote the success notification. Arrival order is not guaranteed.
func TestSetMeshsyncMode_FailedRedeployEmitsIndependentEvents(t *testing.T) {
	const (
		changed = "MeshSync deployment mode changed from 'operator' to 'embedded'"
		failed  = "MeshSync mode saved as 'embedded' but redeploy failed"
	)

	for i := 0; i < 50; i++ {
		h, fake := newControllersConfigTestHandler(t)
		h.ConnectionToStateMachineInstanceTracker = &machines.ConnectionToStateMachineInstanceTracker{
			ConnectToInstanceMap: map[core.Uuid]*machines.StateMachine{},
		}
		provider := &lockedEventsProvider{controllersConfigFakeProvider: fake}
		connID := seedKubernetesConnection(fake, nil)
		user := testUser()
		sub, unsubscribe := h.config.EventBroadcaster.Subscribe(user.ID)

		req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/"+connID.String()+"/actions",
			strings.NewReader(`{"action":"setMeshsyncMode","mode":"embedded"}`))
		req = withToken(mux.SetURLVars(req, map[string]string{"connectionId": connID.String()}))
		rec := httptest.NewRecorder()
		h.PerformConnectionAction(rec, req, nil, user, provider)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body.String())
		}

		var broadcast []*events.Event
		for len(broadcast) < 2 {
			select {
			case msg := <-sub:
				broadcast = append(broadcast, msg.(*events.Event))
			case <-time.After(2 * time.Second):
				t.Fatalf("received %d broadcast events, want 2", len(broadcast))
			}
		}
		unsubscribe()
		if broadcast[0] == broadcast[1] {
			t.Fatalf("both broadcasts carry the same *events.Event")
		}
		if broadcast[0].OperationID != broadcast[1].OperationID {
			t.Errorf("live events of one mode change no longer share an OperationID: %s vs %s",
				broadcast[0].OperationID, broadcast[1].OperationID)
		}

		var persisted []events.Event
		for deadline := time.Now().Add(2 * time.Second); len(persisted) < 2 && time.Now().Before(deadline); {
			persisted = provider.snapshot()
			time.Sleep(time.Millisecond)
		}
		if len(persisted) != 2 {
			t.Fatalf("persisted %d events, want 2", len(persisted))
		}

		for _, got := range [][]events.Event{{*broadcast[0], *broadcast[1]}, persisted} {
			var informational, errored int
			for _, e := range got {
				switch {
				case e.Severity == events.Informational && strings.HasPrefix(e.Description, changed):
					informational++
				case e.Severity == events.Error && strings.HasPrefix(e.Description, failed):
					errored++
				default:
					t.Errorf("unexpected event: severity=%s description=%q", e.Severity, e.Description)
				}
			}
			if informational != 1 || errored != 1 {
				t.Fatalf("want one informational and one error event, got %d and %d", informational, errored)
			}
		}
		if persisted[0].OperationID != persisted[1].OperationID {
			t.Errorf("events of one mode change no longer share an OperationID: %s vs %s", persisted[0].OperationID, persisted[1].OperationID)
		}
	}
}

func TestPerformConnectionAction_NilUser_DoesNotPanic(t *testing.T) {
	h, fake := newControllersConfigTestHandler(t)
	h.ConnectionToStateMachineInstanceTracker = &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: map[core.Uuid]*machines.StateMachine{},
	}
	connID := seedKubernetesConnection(fake, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/"+connID.String()+"/actions",
		strings.NewReader(`{"action":"setMeshsyncMode","mode":"operator"}`))
	req = withToken(mux.SetURLVars(req, map[string]string{"connectionId": connID.String()}))
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, nil, fake)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body.String())
	}
}

func TestPerformConnectionAction_NilSystemID_DoesNotPanic(t *testing.T) {
	h, fake := newControllersConfigTestHandler(t)
	h.SystemID = nil
	h.ConnectionToStateMachineInstanceTracker = &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: map[core.Uuid]*machines.StateMachine{},
	}
	connID := seedKubernetesConnection(fake, nil)
	user := testUser()

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections/"+connID.String()+"/actions",
		strings.NewReader(`{"action":"setMeshsyncMode","mode":"operator"}`))
	req = withToken(mux.SetURLVars(req, map[string]string{"connectionId": connID.String()}))
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, user, fake)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body.String())
	}
}

// connectionActionFlushTestProvider supplies the in-memory database used by
// the flush action while reusing the existing fake connection provider.
type connectionActionFlushTestProvider struct {
	*controllersConfigFakeProvider
	db *database.Handler
}

func (p *connectionActionFlushTestProvider) GetGenericPersister() *database.Handler {
	return p.db
}

func TestPerformConnectionAction_FlushMeshsync_DisconnectedDoesNotDeleteCache(t *testing.T) {
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("failed to get SQL database: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })

	db := &database.Handler{DB: gdb}
	if err := db.AutoMigrate(
		&meshsyncmodel.KubernetesKeyValue{},
		&meshsyncmodel.KubernetesResource{},
		&meshsyncmodel.KubernetesResourceSpec{},
		&meshsyncmodel.KubernetesResourceStatus{},
		&meshsyncmodel.KubernetesResourceObjectMeta{},
	); err != nil {
		t.Fatalf("failed to migrate MeshSync tables: %v", err)
	}

	serverID := uuid.Must(uuid.NewV4())
	resourceID := uuid.Must(uuid.NewV4()).String()
	resource := &meshsyncmodel.KubernetesResource{
		ID:         resourceID,
		ClusterID:  serverID.String(),
		Kind:       "Service",
		APIVersion: "v1",
	}
	if err := db.Create(resource).Error; err != nil {
		t.Fatalf("failed to seed MeshSync resource: %v", err)
	}

	baseProvider := newControllersConfigFakeProvider()
	connectionID := seedKubernetesConnection(baseProvider, nil)
	provider := &connectionActionFlushTestProvider{
		controllersConfigFakeProvider: baseProvider,
		db:                            db,
	}

	h := newConnectionActionTestHandler(t)
	machine, err := k8smachine.New(connectionID.String(), uuid.Must(uuid.NewV4()), h.log)
	if err != nil {
		t.Fatalf("failed to create Kubernetes state machine: %v", err)
	}
	machine.Context = &k8smachine.MachineCtx{
		K8sContext: models.K8sContext{
			ID:                 "test-context",
			KubernetesServerID: &serverID,
		},
	}

	_, _, err = (&k8smachine.DisconnectAction{}).Execute(
		context.Background(), machine.GetContext(), nil,
	)
	if err != nil {
		t.Fatalf("failed to mark machine disconnected: %v", err)
	}

	h.ConnectionToStateMachineInstanceTracker = &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
	h.ConnectionToStateMachineInstanceTracker.Add(connectionID, machine)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/integrations/connections/"+connectionID.String()+"/actions",
		strings.NewReader(`{"action":"flushMeshsync"}`),
	)
	req = mux.SetURLVars(req, map[string]string{"connectionId": connectionID.String()})
	rec := httptest.NewRecorder()

	h.PerformConnectionAction(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "disconnected") {
		t.Errorf("expected disconnected error, got: %s", rec.Body.String())
	}

	var count int64
	if err := db.Model(&meshsyncmodel.KubernetesResource{}).
		Where("cluster_id = ?", serverID.String()).
		Count(&count).Error; err != nil {
		t.Fatalf("failed to count cached MeshSync resources: %v", err)
	}
	if count != 1 {
		t.Errorf("cached MeshSync resource count = %d, want 1; disconnected flush must not delete cache", count)
	}
}
