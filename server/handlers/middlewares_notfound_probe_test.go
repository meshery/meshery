package handlers

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// notFoundProbeTestProvider is a fake Provider for the NOTFOUND probe tests.
// It embeds recordingEventsProvider (PersistEvent + CorrelatesConditions)
// and stubs the three calls the probe path makes: GetK8sContexts (the
// NOTFOUND load), GetConnectionByID and UpdateConnectionById (the status
// write SendEvent performs on every Discovery).
//
// PersistEvent is overridden to serialize with the provider's own mutex:
// the probe runs on its own goroutine, and the test reads the persisted
// events while it runs.
type notFoundProbeTestProvider struct {
	recordingEventsProvider

	mu               sync.Mutex
	getContextsCalls []string // withStatus values seen, in order
	notFoundPage     []byte
	connection       *connections.Connection
}

func (p *notFoundProbeTestProvider) PersistEvent(event events.Event, token string) (models.PersistedEvent, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.recordingEventsProvider.PersistEvent(event, token)
}

func (p *notFoundProbeTestProvider) GetK8sContexts(_, _, _, _, _, withStatus string, _ bool) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.getContextsCalls = append(p.getContextsCalls, withStatus)
	return p.notFoundPage, nil
}

func (p *notFoundProbeTestProvider) GetConnectionByID(_ string, _ core.Uuid) (*connections.Connection, int, error) {
	return p.connection, 0, nil
}

func (p *notFoundProbeTestProvider) UpdateConnectionById(_ string, conn *connections.ConnectionPayload, _ string) (*connections.Connection, error) {
	// Mirror the write the real provider performs: the persisted status
	// becomes the state the machine settled in.
	updated := *p.connection
	updated.Status = conn.Status
	p.connection = &updated
	return p.connection, nil
}

func (p *notFoundProbeTestProvider) getContextsCallCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.getContextsCalls)
}

func (p *notFoundProbeTestProvider) persistedCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.persisted)
}

func (p *notFoundProbeTestProvider) persistedEvent(i int) events.Event {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.persisted[i]
}

// nextEventAction is a stub machine action that emits one follow-up event
// and never fails. It stands in for the real Register/Connect/NotFound
// actions, which need a live cluster.
type nextEventAction struct {
	next machines.EventType
}

func (a *nextEventAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (a *nextEventAction) Execute(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return a.next, nil, nil
}

func (a *nextEventAction) ExecuteOnExit(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

// resolvingProbeAction simulates a cluster that came back: Discovery
// succeeds, so the machine walks DISCOVERED -> REGISTERED -> CONNECTED and
// the probe resolves the connection's unreachable condition.
type resolvingProbeAction struct {
	probed chan struct{}
}

func (a *resolvingProbeAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (a *resolvingProbeAction) Execute(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	select {
	case a.probed <- struct{}{}:
	default:
	}
	return machines.Register, nil, nil
}

func (a *resolvingProbeAction) ExecuteOnExit(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

// newNotFoundProbeFixture builds a Handler with a fake clock, a tracker
// pre-seeded with a hand-built machine for one NOTFOUND connection, and a
// fake provider whose NOTFOUND page holds that connection's context. The
// machine mirrors the real kubernetes machine's shape (NOTFOUND has an
// action, ConditionComponent is "k8s-connection") with the cluster probe
// stubbed: probeAction runs in place of DiscoverAction.
//
// It also subscribes the user to the Handler's event broadcaster, so tests
// can assert what a probe notifies (a Recurrence must not broadcast; a
// resolution must).
func newNotFoundProbeFixture(t *testing.T, probeAction machines.Action) (*Handler, *notFoundProbeTestProvider, *fakeClock, core.Uuid, *models.User, chan interface{}) {
	t.Helper()

	connectionID := uuid.Must(uuid.NewV4())
	clock := &fakeClock{now: time.Now()}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	// The context the NOTFOUND load returns. It is deliberately NOT in
	// AllKubeClusterKey: that is the bug — Discovery never runs for it.
	notFoundCtx := &models.K8sContext{
		ID:           connectionID.String(),
		Name:         "dead-ctx",
		ConnectionID: connectionID.String(),
	}
	page, err := json.Marshal(models.MesheryK8sContextPage{
		Page:       0,
		PageSize:   25,
		TotalCount: 1,
		Contexts:   []*models.K8sContext{notFoundCtx},
	})
	if err != nil {
		t.Fatalf("marshal NOTFOUND page: %v", err)
	}

	provider := &notFoundProbeTestProvider{
		notFoundPage: page,
		connection: &connections.Connection{
			ID:     connectionID,
			Name:   "dead-cluster",
			Kind:   "kubernetes",
			Status: connections.NOTFOUND,
		},
	}

	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: map[uuid.UUID]*machines.StateMachine{},
	}
	tracker.Add(connectionID, &machines.StateMachine{
		ID:                 connectionID,
		Name:               "kubernetes",
		ConditionComponent: "k8s-connection",
		InitialState:       machines.InitialState,
		CurrentState:       machines.NOTFOUND,
		PreviousState:      machines.DISCOVERED,
		Log:                newTestLogger(t),
		Provider:           provider,
		Context:            &kubernetes.MachineCtx{MesheryCtrlsHelper: &models.MesheryControllersHelper{}},
		States: machines.States{
			machines.InitialState: machines.State{Events: machines.Events{machines.Discovery: machines.DISCOVERED}},
			machines.DISCOVERED: machines.State{
				Events: machines.Events{machines.Register: machines.REGISTERED, machines.NotFound: machines.NOTFOUND},
				Action: probeAction,
			},
			machines.REGISTERED: machines.State{
				Events: machines.Events{machines.Connect: machines.CONNECTED, machines.NotFound: machines.NOTFOUND},
				Action: &nextEventAction{next: machines.Connect},
			},
			machines.CONNECTED: machines.State{
				Events: machines.Events{machines.Disconnect: machines.DISCONNECTED, machines.NotFound: machines.NOTFOUND},
				Action: &nextEventAction{next: machines.NoOp},
			},
			machines.NOTFOUND: machines.State{
				Events: machines.Events{machines.Discovery: machines.DISCOVERED},
				Action: &nextEventAction{next: machines.NoOp},
			},
		},
	})

	broadcaster := models.NewBroadcaster("test")
	broadcastCh, unsubscribe := broadcaster.Subscribe(user.ID)
	t.Cleanup(unsubscribe)

	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: broadcaster},
		log:                                     newTestLogger(t),
		ConnectionToStateMachineInstanceTracker: tracker,
		notFoundProbeThrottle:                   newNotFoundProbeThrottle(clock.Now, notFoundProbeInterval),
	}

	return h, provider, clock, connectionID, user, broadcastCh
}

// fakeClock is a controllable clock for the throttle: tests advance it
// instead of sleeping.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// waitForProbes blocks until n probes have fired, or fails the test.
func waitForProbes(t *testing.T, probed chan struct{}, n int) {
	t.Helper()
	timeout := time.After(10 * time.Second)
	for i := 0; i < n; i++ {
		select {
		case <-probed:
		case <-timeout:
			t.Fatalf("timed out waiting for probe %d of %d", i+1, n)
		}
	}
}

// waitForPersisted blocks until the provider has persisted n events, or
// fails the test. Probes run on their own goroutine, so the persist lands
// after K8sFSMMiddleware returns.
func waitForPersisted(t *testing.T, provider *notFoundProbeTestProvider, n int) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for {
		if count := provider.persistedCount(); count >= n {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %d persisted events, got %d", n, provider.persistedCount())
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// waitForThrottleForget blocks until the throttle no longer holds an entry
// for id, or fails the test. The forget runs in the probe goroutine after
// the resolution is reported.
func waitForThrottleForget(t *testing.T, throttle *notFoundProbeThrottle, id core.Uuid) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for {
		throttle.mu.Lock()
		_, ok := throttle.lastProbed[id]
		throttle.mu.Unlock()
		if !ok {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for the throttle entry for %s to be forgotten", id)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// runFSM drives one K8sFSMMiddleware pass with an empty AllKubeClusterKey:
// no connected contexts at all, so the only Discovery that can fire is the
// NOTFOUND probe.
func runFSM(t *testing.T, h *Handler, provider models.Provider, user *models.User) {
	t.Helper()
	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "")
	ctx = context.WithValue(ctx, models.AllKubeClusterKey, []*models.K8sContext{})
	K8sFSMMiddleware(ctx, h, provider, user)
}

// TestK8sFSMMiddlewareProbesNotFoundConnections is the regression test for
// issue #22955: a NOTFOUND connection that is excluded from
// AllKubeClusterKey is still probed, at most once per minute, and a
// still-unreachable probe counts an occurrence on the open condition
// (a Recurrence) instead of notifying again.
func TestK8sFSMMiddlewareProbesNotFoundConnections(t *testing.T) {
	probed := make(chan struct{}, 16)
	action := &probeCountingAction{probed: probed}
	h, provider, clock, _, user, broadcastCh := newNotFoundProbeFixture(t, action)

	// First pass: the NOTFOUND list is loaded once, with status "not found",
	// and the connection is probed.
	runFSM(t, h, provider, user)
	waitForProbes(t, probed, 1)
	waitForPersisted(t, provider, 1)

	if got := provider.getContextsCallCount(); got != 1 {
		t.Fatalf("GetK8sContexts calls = %d, want 1", got)
	}
	provider.mu.Lock()
	status := provider.getContextsCalls[0]
	provider.mu.Unlock()
	persisted := provider.persistedEvent(0)
	if status != string(connections.NOTFOUND) {
		t.Errorf("NOTFOUND load status = %q, want %q", status, string(connections.NOTFOUND))
	}
	if persisted.CorrelationKey == nil || *persisted.CorrelationKey == "" {
		t.Error("still-unreachable probe was not counted as a Recurrence (no correlation key)")
	}
	if persisted.ResolvedAt != nil {
		t.Error("still-unreachable probe resolved the condition, want it left open")
	}
	if got := provider.connection.Status; got != connections.NOTFOUND {
		t.Errorf("connection status after failed probe = %q, want %q", got, connections.NOTFOUND)
	}
	// The Recurrence is counted on the open condition and sends no new
	// notification: nothing is broadcast for it.
	select {
	case data := <-broadcastCh:
		t.Errorf("still-unreachable probe broadcast %T, want no new notification", data)
	case <-time.After(500 * time.Millisecond):
	}

	// Second pass inside the same minute: no reload, no second probe.
	runFSM(t, h, provider, user)
	select {
	case <-probed:
		t.Error("connection probed twice within the throttle interval")
	case <-time.After(500 * time.Millisecond):
	}
	if got := provider.getContextsCallCount(); got != 1 {
		t.Errorf("GetK8sContexts calls after second pass = %d, want 1 (load gated)", got)
	}

	// A minute later: the list reloads and the connection is probed again.
	clock.advance(notFoundProbeInterval + time.Second)
	runFSM(t, h, provider, user)
	waitForProbes(t, probed, 1)
	waitForPersisted(t, provider, 2)
	if got := provider.getContextsCallCount(); got != 2 {
		t.Errorf("GetK8sContexts calls after a minute = %d, want 2", got)
	}
}

// TestK8sFSMMiddlewareForgetsResolvedNotFoundConnection checks the other
// half of #22955: when a probe finds the cluster reachable again, the
// condition resolves and the throttle entry is forgotten, so the connection
// is probed again on the next load cycle instead of waiting out the
// interval.
func TestK8sFSMMiddlewareForgetsResolvedNotFoundConnection(t *testing.T) {
	probed := make(chan struct{}, 16)
	action := &resolvingProbeAction{probed: probed}
	h, provider, clock, connectionID, user, broadcastCh := newNotFoundProbeFixture(t, action)

	runFSM(t, h, provider, user)
	waitForProbes(t, probed, 1)
	waitForPersisted(t, provider, 1)

	resolution := provider.persistedEvent(0)
	if resolution.ResolvedAt == nil {
		t.Fatal("reachable probe did not resolve the condition")
	}
	if resolution.CorrelationKey == nil || *resolution.CorrelationKey == "" {
		t.Error("resolution event has no correlation key")
	}
	if got := provider.connection.Status; got != connections.CONNECTED {
		t.Errorf("connection status after resolving probe = %q, want %q", got, connections.CONNECTED)
	}
	// The resolution is notified: the stored condition event is broadcast.
	broadcast := receiveEvent(t, broadcastCh)
	if broadcast.ResolvedAt == nil {
		t.Error("broadcast resolution has no ResolvedAt")
	}

	// The connection left NOTFOUND, so its throttle entry is forgotten.
	waitForThrottleForget(t, h.notFoundProbeThrottle, connectionID)

	// On the next load cycle the connection is probed again immediately:
	// the forgotten entry does not hold it back for the rest of the
	// interval. (The load itself is still gated, so the clock must move
	// past the interval first.)
	clock.advance(notFoundProbeInterval + time.Second)
	runFSM(t, h, provider, user)
	waitForProbes(t, probed, 1)
}

// TestNotFoundProbeThrottle is a unit test for the throttle itself: the load
// gate and the per-connection probe gate are atomic check-and-set, stamped
// at start, and forget/prune drop entries for connections that left
// NOTFOUND or were deleted.
func TestNotFoundProbeThrottle(t *testing.T) {
	clock := &fakeClock{now: time.Now()}
	throttle := newNotFoundProbeThrottle(clock.Now, time.Minute)
	id := uuid.Must(uuid.NewV4())
	other := uuid.Must(uuid.NewV4())

	if !throttle.allowLoad() {
		t.Fatal("first load must be allowed")
	}
	if throttle.allowLoad() {
		t.Error("second load inside the interval must be gated")
	}
	if !throttle.allowProbe(id) {
		t.Fatal("first probe must be allowed")
	}
	if throttle.allowProbe(id) {
		t.Error("second probe inside the interval must be gated")
	}
	if !throttle.allowProbe(other) {
		t.Error("a different connection must be probed independently")
	}

	// forget drops one entry; prune drops entries absent from the list.
	throttle.forget(id)
	if !throttle.allowProbe(id) {
		t.Error("probe must be allowed again after forget")
	}
	throttle.prune(map[core.Uuid]struct{}{other: {}})
	if !throttle.allowProbe(id) {
		t.Error("probe must be allowed again after prune dropped it")
	}
	if throttle.allowProbe(other) {
		t.Error("prune must keep entries still in the NOTFOUND list")
	}

	// Advancing the clock past the interval re-allows both gates.
	clock.advance(time.Minute + time.Second)
	if !throttle.allowLoad() {
		t.Error("load must be allowed again after the interval")
	}
	if !throttle.allowProbe(id) {
		t.Error("probe must be allowed again after the interval")
	}
	if !throttle.allowProbe(other) {
		t.Error("probe must be allowed again after the interval")
	}
}
