package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	mhelpers "github.com/meshery/meshery/server/machines/helpers"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// probeAction stands in for a state's action so a test can see whether, how
// often and against which machine context a transition ran, without a cluster.
type probeAction struct {
	provider *models.DefaultLocalProvider
	id       core.Uuid
	err      error         // returned from Execute, halting the transition
	entered  chan struct{} // closed on entry when non-nil
	release  chan struct{} // Execute waits on it when non-nil

	calls       atomic.Int32
	sawRow      atomic.Bool
	sawCtx      atomic.Value
	sawClusters atomic.Value // []*models.K8sContext under AllKubeClusterKey
	enterOnce   sync.Once
}

func (p *probeAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (p *probeAction) Execute(ctx context.Context, machineCtx interface{}, _ interface{}) (machines.EventType, *events.Event, error) {
	p.calls.Add(1)
	p.sawCtx.Store(machineCtx)
	if clusters, ok := ctx.Value(models.AllKubeClusterKey).([]*models.K8sContext); ok {
		p.sawClusters.Store(clusters)
	}
	if _, _, err := p.provider.GetConnectionByID("", p.id); err == nil {
		p.sawRow.Store(true)
	}
	if p.entered != nil {
		p.enterOnce.Do(func() { close(p.entered) })
	}
	if p.release != nil {
		<-p.release
	}
	return machines.NoOp, nil, p.err
}

func (p *probeAction) ExecuteOnExit(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

type lifecycleFixture struct {
	h        *Handler
	provider *models.DefaultLocalProvider
	tracker  *machines.ConnectionToStateMachineInstanceTracker
	userID   core.Uuid
}

func newLifecycleFixture(t *testing.T) *lifecycleFixture {
	t.Helper()
	h, provider := newDeleteConnectionFixture(t)
	// Retirement loads the connected contexts, as KubernetesMiddleware does.
	db := provider.ConnectionPersister.DB
	if err := db.AutoMigrate(models.K8sContext{}); err != nil {
		t.Fatalf("migrate k8s contexts: %v", err)
	}
	provider.MesheryK8sContextPersister = &models.MesheryK8sContextPersister{DB: db}
	tracker := &machines.ConnectionToStateMachineInstanceTracker{ConnectToInstanceMap: map[core.Uuid]*machines.StateMachine{}}
	h.ConnectionToStateMachineInstanceTracker = tracker
	return &lifecycleFixture{h: h, provider: provider, tracker: tracker, userID: uuid.Must(uuid.NewV4())}
}

// k8sConnectionID derives the id the way SaveK8sContext does, so a second
// import of the same context lands on the same id.
func k8sConnectionID(t *testing.T) core.Uuid {
	t.Helper()
	instanceID := uuid.Must(uuid.NewV4())
	id, err := models.K8sContextGenerateID(models.K8sContext{
		Name:              "kind-test",
		MesheryInstanceID: &instanceID,
		Cluster:           map[string]interface{}{"server": "https://127.0.0.1:6443"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return uuid.FromStringOrNil(id)
}

func (f *lifecycleFixture) saveRow(t *testing.T, id core.Uuid, status connections.ConnectionStatus) {
	t.Helper()
	if _, err := f.provider.ConnectionPersister.SaveConnection(&connections.Connection{
		ID: id, Name: "kind-test", Kind: "kubernetes", ConnectionType: "platform", Status: status,
	}); err != nil {
		t.Fatalf("save connection: %v", err)
	}
}

func (f *lifecycleFixture) rowStatus(id core.Uuid) (connections.ConnectionStatus, bool) {
	c, err := f.provider.ConnectionPersister.GetConnection(id, "")
	if err != nil {
		return "", false
	}
	return c.Status, true
}

// trackMachine registers a kubernetes machine in state, with probe as the
// action of its DELETED state. A nil machineCtx leaves the Context unassigned.
func (f *lifecycleFixture) trackMachine(t *testing.T, id core.Uuid, state machines.StateType, machineCtx interface{}, probe *probeAction) *machines.StateMachine {
	t.Helper()
	inst, err := kubernetes.New(id.String(), f.userID, f.h.log)
	if err != nil {
		t.Fatal(err)
	}
	inst.CurrentState = state
	inst.Context = machineCtx
	inst.Provider = f.provider
	deleted := inst.States[machines.DELETED]
	inst.States[machines.DELETED] = *deleted.RegisterAction(probe)
	f.tracker.Add(id, inst)
	return inst
}

// reimport runs what a kubeconfig re-import does for an existing id: obtain
// the machine through InitializeMachineWithContext.
func (f *lifecycleFixture) reimport(t *testing.T, id core.Uuid, newCtx interface{}) (*machines.StateMachine, int) {
	t.Helper()
	initCalls := 0
	initFn := func(context.Context, interface{}, logger.Handler) (interface{}, *events.Event, error) {
		initCalls++
		return newCtx, nil, nil
	}
	inst, err := mhelpers.InitializeMachineWithContext(&kubernetes.MachineCtx{}, f.requestCtx(), id, f.userID,
		f.tracker, f.h.log, f.provider, machines.DefaultState, "kubernetes", initFn)
	if err != nil {
		t.Errorf("InitializeMachineWithContext: %v", err) // may run off the test goroutine
	}
	return inst, initCalls
}

func (f *lifecycleFixture) requestCtx() context.Context {
	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: f.userID})
	return context.WithValue(ctx, models.SystemIDKey, f.h.SystemID)
}

// The hard-delete route has no KubernetesMiddleware. The Delete transition
// must still receive the connected contexts, including the one being deleted,
// or DeleteAction's MeshSync flush cannot run its last-context check.
func TestDeleteConnectionGivesTeardownTheConnectedContexts(t *testing.T) {
	f := newLifecycleFixture(t)
	id := k8sConnectionID(t)
	f.saveRow(t, id, connections.CONNECTED)
	contextID := strings.ReplaceAll(id.String(), "-", "")
	if _, err := f.provider.MesheryK8sContextPersister.SaveMesheryK8sContext(models.K8sContext{ID: contextID, Name: "kind-test"}); err != nil {
		t.Fatalf("save k8s context: %v", err)
	}
	probe := &probeAction{provider: f.provider, id: id}
	f.trackMachine(t, id, machines.CONNECTED, &kubernetes.MachineCtx{}, probe)

	if rec := deleteConnection(t, f.h, f.provider, id); rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200. body: %s", rec.Code, rec.Body.String())
	}

	clusters, _ := probe.sawClusters.Load().([]*models.K8sContext)
	found := false
	for _, c := range clusters {
		if c != nil && c.ID == contextID {
			found = true
		}
	}
	if !found {
		t.Fatalf("delete transition saw contexts %v under AllKubeClusterKey, want them to include %s", clusters, contextID)
	}
}

// Hard delete must run the Delete transition while the row is still there,
// persist the "deleted" status through it, and drop the machine.
func TestDeleteConnectionRetiresTrackedMachine(t *testing.T) {
	f := newLifecycleFixture(t)
	id := k8sConnectionID(t)
	f.saveRow(t, id, connections.CONNECTED)
	probe := &probeAction{provider: f.provider, id: id}
	old := f.trackMachine(t, id, machines.CONNECTED, &kubernetes.MachineCtx{}, probe)

	if rec := deleteConnection(t, f.h, f.provider, id); rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200. body: %s", rec.Code, rec.Body.String())
	}

	if n := probe.calls.Load(); n != 1 {
		t.Errorf("delete transition ran %d times, want 1", n)
	}
	if !probe.sawRow.Load() {
		t.Error("delete transition ran after the connection row was gone")
	}
	if old.CurrentState != machines.DELETED {
		t.Errorf("machine state = %q, want %q", old.CurrentState, machines.DELETED)
	}
	if _, ok := f.tracker.Get(id); ok {
		t.Error("tracker still holds the machine of a deleted connection")
	}
	if _, ok := f.rowStatus(id); ok {
		t.Error("connection row still present after delete")
	}
}

// Re-importing the same context after a hard delete must get a new machine
// that accepts Discovery, and must keep the new row's status.
func TestDeleteConnectionThenReimportGetsFreshMachine(t *testing.T) {
	f := newLifecycleFixture(t)
	id := k8sConnectionID(t)
	f.saveRow(t, id, connections.CONNECTED)
	oldCtx := &kubernetes.MachineCtx{}
	probe := &probeAction{provider: f.provider, id: id}
	old := f.trackMachine(t, id, machines.CONNECTED, oldCtx, probe)

	if rec := deleteConnection(t, f.h, f.provider, id); rec.Code != http.StatusOK {
		t.Fatalf("delete status = %d. body: %s", rec.Code, rec.Body.String())
	}

	f.saveRow(t, id, connections.DISCOVERED)
	newCtx := &kubernetes.MachineCtx{}
	fresh, initCalls := f.reimport(t, id, newCtx)
	if fresh == old {
		t.Fatal("re-import reused the deleted connection's machine")
	}
	if initCalls != 1 || fresh.CurrentState != machines.InitialState {
		t.Fatalf("re-import: init calls = %d, state = %q; want a newly initialized machine", initCalls, fresh.CurrentState)
	}

	// The old teardown ran against the old machine context only; the new
	// machine carries its own.
	if got := probe.sawCtx.Load(); got != interface{}(oldCtx) {
		t.Error("old teardown did not run against the old machine context")
	}
	if fresh.Context == interface{}(oldCtx) {
		t.Error("new machine shares the retired machine's context")
	}

	// DiscoverAction needs a reachable cluster; stub it so the test exercises
	// the transition and status write-back only.
	discover := &probeAction{provider: f.provider, id: id}
	discovered := fresh.States[machines.DISCOVERED]
	fresh.States[machines.DISCOVERED] = *discovered.RegisterAction(discover)

	if _, err := fresh.SendEvent(f.requestCtx(), mhelpers.StatusToEvent(connections.DISCOVERED), nil); err != nil {
		t.Fatalf("discovery on the re-imported connection failed: %v", err)
	}
	if fresh.CurrentState != machines.DISCOVERED {
		t.Errorf("state after discovery = %q, want %q", fresh.CurrentState, machines.DISCOVERED)
	}
	if status, _ := f.rowStatus(id); status != connections.DISCOVERED {
		t.Errorf("re-imported row status = %q, want %q", status, connections.DISCOVERED)
	}
}

// If the tracker entry is replaced while the delete transition is running
// (another path retired the machine and a re-import registered a new one),
// the delete must leave the replacement alone.
func TestDeleteConnectionKeepsReplacementMachine(t *testing.T) {
	f := newLifecycleFixture(t)
	id := k8sConnectionID(t)
	f.saveRow(t, id, connections.CONNECTED)
	probe := &probeAction{provider: f.provider, id: id, entered: make(chan struct{}), release: make(chan struct{})}
	f.trackMachine(t, id, machines.CONNECTED, &kubernetes.MachineCtx{}, probe)

	done := make(chan int)
	go func() { done <- deleteConnection(t, f.h, f.provider, id).Code }()

	select {
	case <-probe.entered:
	case code := <-done:
		close(probe.release)
		t.Fatalf("delete returned %d without running the delete transition", code)
	}
	replacement, err := kubernetes.New(id.String(), f.userID, f.h.log)
	if err != nil {
		t.Fatal(err)
	}
	f.tracker.Remove(id)
	f.tracker.Add(id, replacement)
	close(probe.release)

	if code := <-done; code != http.StatusOK {
		t.Fatalf("delete status = %d, want 200", code)
	}
	if got, ok := f.tracker.Get(id); !ok || got != replacement {
		t.Fatal("the delete removed a machine registered after it started")
	}
}

func TestDeleteConnectionUntrackedMachine(t *testing.T) {
	f := newLifecycleFixture(t)
	id := k8sConnectionID(t)
	f.saveRow(t, id, connections.CONNECTED)

	if rec := deleteConnection(t, f.h, f.provider, id); rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200. body: %s", rec.Code, rec.Body.String())
	}
	if _, ok := f.rowStatus(id); ok {
		t.Error("connection row still present after delete")
	}
	if _, ok := f.tracker.Get(id); ok {
		t.Error("delete registered a machine for an untracked connection")
	}
}

// A machine that cannot complete its delete transition must not survive the
// connection: the row is still deleted, the entry is still dropped, and the
// failure is reported on the delete event.
func TestDeleteConnectionTeardownFailure(t *testing.T) {
	cases := []struct {
		name       string
		state      machines.StateType
		machineCtx interface{}
		probeErr   error
		wantCalls  int32
		wantWarn   bool
	}{
		{name: "delete action fails", state: machines.CONNECTED, machineCtx: &kubernetes.MachineCtx{}, probeErr: errors.New("undeploy failed"), wantCalls: 1, wantWarn: true},
		{name: "state has no delete edge", state: machines.REGISTERED, machineCtx: &kubernetes.MachineCtx{}, wantCalls: 0, wantWarn: true},
		{name: "machine context never assigned", state: machines.CONNECTED, machineCtx: nil, wantCalls: 0, wantWarn: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newLifecycleFixture(t)
			id := k8sConnectionID(t)
			f.saveRow(t, id, connections.ConnectionStatus(tc.state))
			probe := &probeAction{provider: f.provider, id: id, err: tc.probeErr}
			f.trackMachine(t, id, tc.state, tc.machineCtx, probe)

			if rec := deleteConnection(t, f.h, f.provider, id); rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200. body: %s", rec.Code, rec.Body.String())
			}
			if n := probe.calls.Load(); n != tc.wantCalls {
				t.Errorf("delete action ran %d times, want %d", n, tc.wantCalls)
			}
			if _, ok := f.tracker.Get(id); ok {
				t.Error("tracker still holds the machine after its connection was deleted")
			}
			if _, ok := f.rowStatus(id); ok {
				t.Error("connection row still present after delete")
			}

			var persisted []events.Event
			if err := f.provider.EventsPersister.DB.Where("action = ?", "delete").Find(&persisted).Error; err != nil {
				t.Fatal(err)
			}
			if len(persisted) != 1 {
				t.Fatalf("persisted %d delete events, want 1", len(persisted))
			}
			// Read back from the database: the reason must survive JSON.
			teardownErr, _ := persisted[0].Metadata["teardownError"].(string)
			gotWarn := persisted[0].Severity == events.Warning && teardownErr != ""
			if gotWarn != tc.wantWarn {
				t.Errorf("delete event severity = %q, teardownError = %q; want warning = %v",
					persisted[0].Severity, teardownErr, tc.wantWarn)
			}
			if tc.probeErr != nil && !strings.Contains(teardownErr, tc.probeErr.Error()) {
				t.Errorf("persisted teardownError = %q, want it to carry %q", teardownErr, tc.probeErr)
			}
		})
	}
}

// refusingProvider fails the row delete the way a remote provider refusing
// it would, after the machine has already been retired.
type refusingProvider struct{ *models.DefaultLocalProvider }

func (refusingProvider) DeleteConnection(*http.Request, core.Uuid) (*connections.Connection, error) {
	return nil, errors.New("provider refused the delete")
}

// The machine cannot be restored once retired, so a refused row delete must
// say the connection was already taken down.
func TestDeleteConnectionProviderRefusalAfterRetirement(t *testing.T) {
	f := newLifecycleFixture(t)
	id := k8sConnectionID(t)
	f.saveRow(t, id, connections.CONNECTED)
	probe := &probeAction{provider: f.provider, id: id}
	f.trackMachine(t, id, machines.CONNECTED, &kubernetes.MachineCtx{}, probe)

	rec := deleteConnection(t, f.h, refusingProvider{f.provider}, id)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500. body: %s", rec.Code, rec.Body.String())
	}
	if status, ok := f.rowStatus(id); !ok || status != connections.DELETED {
		t.Errorf("row = (%q, present %v), want it kept with status %q", status, ok, connections.DELETED)
	}
	if _, ok := f.tracker.Get(id); ok {
		t.Error("tracker still holds the retired machine")
	}

	var persisted []events.Event
	if err := f.provider.EventsPersister.DB.Where("action = ?", "delete").Find(&persisted).Error; err != nil {
		t.Fatal(err)
	}
	if len(persisted) != 1 {
		t.Fatalf("persisted %d delete events, want 1", len(persisted))
	}
	if retired, _ := persisted[0].Metadata["connectionRetired"].(bool); !retired || persisted[0].Severity != events.Error {
		t.Errorf("failure event = (%q, %q, connectionRetired %v); want an error that reports the retirement",
			persisted[0].Severity, persisted[0].Description, persisted[0].Metadata["connectionRetired"])
	}
	if v, ok := persisted[0].Metadata["teardownError"]; ok {
		t.Errorf("failure event carries teardownError %v although the teardown succeeded", v)
	}
}

// When the teardown fails and the provider then refuses the row delete, the
// failure event must still carry why the teardown failed - as a string, so it
// survives persistence.
func TestDeleteConnectionProviderRefusalAfterTeardownFailure(t *testing.T) {
	f := newLifecycleFixture(t)
	id := k8sConnectionID(t)
	f.saveRow(t, id, connections.CONNECTED)
	probe := &probeAction{provider: f.provider, id: id, err: errors.New("undeploy failed")}
	f.trackMachine(t, id, machines.CONNECTED, &kubernetes.MachineCtx{}, probe)

	rec := deleteConnection(t, f.h, refusingProvider{f.provider}, id)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500. body: %s", rec.Code, rec.Body.String())
	}
	if n := probe.calls.Load(); n != 1 {
		t.Errorf("delete action ran %d times, want 1", n)
	}
	if _, ok := f.tracker.Get(id); ok {
		t.Error("tracker still holds the machine after a failed teardown")
	}

	var persisted []events.Event
	if err := f.provider.EventsPersister.DB.Where("action = ?", "delete").Find(&persisted).Error; err != nil {
		t.Fatal(err)
	}
	if len(persisted) != 1 {
		t.Fatalf("persisted %d delete events, want 1", len(persisted))
	}
	if persisted[0].Severity != events.Error {
		t.Errorf("failure event severity = %q, want %q", persisted[0].Severity, events.Error)
	}
	// Read back from the database: the reason must survive JSON.
	if teardownErr, _ := persisted[0].Metadata["teardownError"].(string); !strings.Contains(teardownErr, "undeploy failed") {
		t.Errorf("persisted teardownError = %#v, want a string carrying %q", persisted[0].Metadata["teardownError"], "undeploy failed")
	}
	if _, ok := persisted[0].Metadata["connectionRetired"]; ok {
		t.Error("failure event reports a retirement although the delete transition did not complete")
	}
}

// A delete racing a re-import of the same id must never leave the deleted
// connection's machine tracked, and must run its teardown exactly once.
func TestDeleteConnectionConcurrentReimport(t *testing.T) {
	f := newLifecycleFixture(t)
	id := k8sConnectionID(t)

	for round := 0; round < 100; round++ {
		f.tracker.Remove(id)
		f.saveRow(t, id, connections.CONNECTED)
		probe := &probeAction{provider: f.provider, id: id}
		old := f.trackMachine(t, id, machines.CONNECTED, &kubernetes.MachineCtx{}, probe)

		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			if rec := deleteConnection(t, f.h, f.provider, id); rec.Code != http.StatusOK {
				t.Errorf("round %d: delete status = %d", round, rec.Code)
			}
		}()
		go func() {
			defer wg.Done()
			<-start
			// Only the machine lookup races here; the in-memory SQLite
			// fixture is not safe for concurrent connections.
			_, _ = f.reimport(t, id, &kubernetes.MachineCtx{})
		}()
		close(start)
		wg.Wait()

		if n := probe.calls.Load(); n != 1 {
			t.Fatalf("round %d: old teardown ran %d times, want 1", round, n)
		}
		if got, ok := f.tracker.Get(id); ok && got == old {
			t.Fatalf("round %d: deleted connection's machine is still tracked", round)
		}
	}
}
