package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
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
)

// newDeleteConnectionFixture wires the handler the way the router does, over an
// in-memory database. It deliberately goes through DefaultLocalProvider rather
// than a stub: the behaviour under test is what the handler makes of a
// provider's error, and the local provider is the one whose "not found" the
// remote provider was taught to match.
func newDeleteConnectionFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider) {
	t.Helper()

	// A file, not ":memory:": every pooled sqlite connection to ":memory:"
	// opens its own empty database, so a query racing a background Discovery
	// would find no tables.
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "meshery.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}, models.K8sContext{}); err != nil {
		t.Fatalf("migrate tables: %v", err)
	}

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     newTestLogger(t),
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: &machines.ConnectionToStateMachineInstanceTracker{ConnectToInstanceMap: map[uuid.UUID]*machines.StateMachine{}},
	}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister:        &models.ConnectionPersister{DB: &db},
		EventsPersister:            &models.EventsPersister{DB: &db},
		MesheryK8sContextPersister: &models.MesheryK8sContextPersister{DB: &db},
		GenericPersister:           &db,
	}
	return h, provider
}

// lookupFailingProvider simulates a transient connection lookup failure
// while preserving the underlying provider's normal delete behavior.
type lookupFailingProvider struct {
	models.Provider
}

func (p *lookupFailingProvider) GetConnectionByID(
	token string,
	connectionID uuid.UUID,
) (*connections.Connection, int, error) {
	return nil, http.StatusInternalServerError, errors.New("simulated temporary lookup failure")
}

func deleteConnection(t *testing.T, h *Handler, provider models.Provider, connectionID uuid.UUID) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodDelete, "/api/integrations/connections/"+connectionID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"connectionId": connectionID.String()})
	rec := httptest.NewRecorder()

	h.DeleteConnection(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	return rec
}

// TestDeleteConnectionNotFound pins the status `mesheryctl connection delete`
// keys off. A syntactically valid id that was never issued must come back 404,
// which the CLI turns into a "No connection with ID ..." warning and exit 0.
// It used to come back 500 "Failed to Save: .connection"
// (meshery-server-1051), because every provider error - not-found included -
// was wrapped in ErrFailToSave.
func TestDeleteConnectionNotFound(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	rec := deleteConnection(t, h, provider, uuid.Must(uuid.NewV4()))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
}

// TestDeleteConnectionExisting is the regression guard for the fix's blast
// radius: the not-found path must not swallow real deletes. An earlier attempt
// short-circuited on the handler's own database before calling the provider,
// which 404s every connection under a remote provider - those rows live at the
// provider, not in Meshery Server's database.
func TestDeleteConnectionExisting(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name:           "test-connection",
		Kind:           "kubernetes",
		ConnectionType: "platform",
	})
	if err != nil {
		t.Fatalf("save connection: %v", err)
	}

	rec := deleteConnection(t, h, provider, saved.ID)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	if _, err := provider.ConnectionPersister.GetConnection(saved.ID, ""); err == nil {
		t.Error("connection is still readable after a successful delete")
	}
}

// Deleting a Kubernetes connection while it is NOTFOUND resolves its open
// unreachable condition instead of leaving it open forever; deleting one that
// is not NOTFOUND resolves nothing.
func TestDeleteConnectionResolvesTheUnreachableCondition(t *testing.T) {
	for _, c := range []struct {
		status      connections.ConnectionStatus
		wantResolve bool
	}{
		{connections.ConnectionStatus(machines.NOTFOUND), true},
		{connections.CONNECTED, false},
	} {
		h, provider := newDeleteConnectionFixture(t)
		user := &models.User{ID: uuid.Must(uuid.NewV4())}
		saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
			Name: "test-connection", Kind: "kubernetes", ConnectionType: "platform", Status: c.status,
		})
		if err != nil {
			t.Fatalf("save connection: %v", err)
		}
		key := kubernetes.UnreachableConditionKey(saved.ID)
		if _, err := provider.PersistEvent(*events.NewEvent().FromOwner(user.ID).WithCategory("connection").
			WithAction("update").WithSeverity(events.Error).WithCorrelationKey(key).Build(), ""); err != nil {
			t.Fatalf("open the condition: %v", err)
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/integrations/connections/"+saved.ID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{"connectionId": saved.ID.String()})
		rec := httptest.NewRecorder()
		h.DeleteConnection(rec, req, nil, user, provider)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d. body: %s", c.status, rec.Code, rec.Body.String())
		}

		var open int64
		if err := provider.EventsPersister.DB.Model(&events.Event{}).
			Where("correlation_key = ? AND resolved_at IS NULL", key).Count(&open).Error; err != nil {
			t.Fatalf("count open conditions: %v", err)
		}
		if got := open == 0; got != c.wantResolve {
			t.Errorf("%s: %d open conditions after the delete, want resolved=%v", c.status, open, c.wantResolve)
		}
	}
}

func TestDeleteConnectionResolvesUnreachableConditionWhenLookupFails(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name:           "unreachable-cluster",
		Kind:           "kubernetes",
		ConnectionType: "platform",
		Status:         connections.ConnectionStatus(machines.NOTFOUND),
	})
	if err != nil {
		t.Fatalf("save connection: %v", err)
	}

	key := kubernetes.UnreachableConditionKey(saved.ID)
	if _, err := provider.PersistEvent(*events.NewEvent().
		FromOwner(user.ID).
		WithCategory("connection").
		WithAction("update").
		WithSeverity(events.Error).
		WithCorrelationKey(key).
		Build(), ""); err != nil {
		t.Fatalf("open unreachable condition: %v", err)
	}

	failingProvider := &lookupFailingProvider{Provider: provider}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/integrations/connections/"+saved.ID.String(),
		nil,
	)
	req = mux.SetURLVars(req, map[string]string{"connectionId": saved.ID.String()})
	rec := httptest.NewRecorder()

	h.DeleteConnection(rec, req, nil, user, failingProvider)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var open int64
	if err := provider.EventsPersister.DB.Model(&events.Event{}).
		Where("correlation_key = ? AND resolved_at IS NULL", key).
		Count(&open).Error; err != nil {
		t.Fatalf("count open conditions: %v", err)
	}

	if open != 0 {
		t.Errorf("unreachable condition remains open after successful deletion; open count = %d", open)
	}
}

// probeCountingAction fails every Discovery the way an unreachable cluster
// does, and reports each run so a test can tell whether a probe happened.
type probeCountingAction struct{ probed chan struct{} }

func (a *probeCountingAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (a *probeCountingAction) Execute(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	a.probed <- struct{}{}
	return machines.NotFound, events.NewEvent().WithSeverity(events.Error).WithDescription("cluster unreachable").Build(), errors.New("cluster unreachable")
}

func (a *probeCountingAction) ExecuteOnExit(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

// Deleting a NOTFOUND Kubernetes connection stops the background probing of
// it. The Kubernetes middlewares used to keep running Discovery on the
// deleted connection's context on every request, and each run recorded a new
// "Failed to retrieve the connection" event: the local provider kept the
// context row and the handler kept the tracked machine.
func TestDeleteConnectionStopsProbingTheDeletedConnection(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name: "dead-cluster", Kind: "kubernetes", ConnectionType: "platform",
		Status: connections.ConnectionStatus(machines.NOTFOUND),
	})
	if err != nil {
		t.Fatalf("save connection: %v", err)
	}
	if err := provider.MesheryK8sContextPersister.DB.Create(&models.K8sContext{
		ID: strings.ReplaceAll(saved.ID.String(), "-", ""), Name: "dead-ctx", ConnectionID: saved.ID.String(),
	}).Error; err != nil {
		t.Fatalf("save context: %v", err)
	}

	action := &probeCountingAction{probed: make(chan struct{}, 16)}
	h.ConnectionToStateMachineInstanceTracker.Add(saved.ID, &machines.StateMachine{
		ID: saved.ID, Name: "kubernetes", InitialState: machines.InitialState,
		CurrentState: machines.NOTFOUND, PreviousState: machines.DISCOVERED,
		Log: h.log, Provider: provider,
		Context: &kubernetes.MachineCtx{MesheryCtrlsHelper: &models.MesheryControllersHelper{}},
		States: machines.States{
			machines.InitialState: machines.State{Events: machines.Events{machines.Discovery: machines.DISCOVERED}},
			machines.DISCOVERED:   machines.State{Events: machines.Events{machines.NotFound: machines.NOTFOUND}, Action: action},
			machines.NOTFOUND:     machines.State{Events: machines.Events{machines.Discovery: machines.DISCOVERED}},
		},
	})

	if rec := deleteConnection(t, h, provider, saved.ID); rec.Code != http.StatusOK {
		t.Fatalf("delete: status = %d. body: %s", rec.Code, rec.Body.String())
	}
	if _, ok := h.ConnectionToStateMachineInstanceTracker.Get(saved.ID); ok {
		t.Error("the deleted connection's machine is still tracked")
	}

	countEvents := func() int64 {
		var n int64
		if err := provider.EventsPersister.DB.Model(&events.Event{}).Count(&n).Error; err != nil {
			t.Fatalf("count events: %v", err)
		}
		return n
	}
	before := countEvents()

	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "")
	for i := 0; i < 3; i++ {
		syncCtx, err := KubernetesMiddleware(ctx, h, provider, user, []string{"all"})
		if err != nil {
			t.Fatalf("sync %d: KubernetesMiddleware: %v", i, err)
		}
		K8sFSMMiddleware(syncCtx, h, provider, user)
	}

	// Discovery runs on its own goroutine; give a probe time to arrive and
	// finish persisting before asserting there was none.
	select {
	case <-action.probed:
		time.Sleep(200 * time.Millisecond)
		t.Errorf("the deleted connection was probed; %d new events recorded", countEvents()-before)
	case <-time.After(500 * time.Millisecond):
	}
	if after := countEvents(); after != before {
		t.Errorf("syncs after the delete recorded %d new events, want 0", after-before)
	}
}
