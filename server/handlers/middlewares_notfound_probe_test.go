package handlers

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// notFoundProbeTestProvider is a fake Provider for the NOTFOUND probe tests.
// It embeds recordingEventsProvider (PersistEvent + CorrelatesConditions)
// and stubs the calls the probe path makes: GetK8sContexts (the NOTFOUND
// load), GetConnectionByID and UpdateConnectionById (the status write
// SendEvent performs on every Discovery), plus the provider-identity calls
// the throttle scoping and the local-provider status check use.
//
// It models a local provider (GetProviderType returns LocalProviderType, so
// the persisted-status check runs), and by default a status-filtered load:
// with filterByStatus set, GetK8sContexts returns the page only when the
// connection's current status matches the requested one, like a real
// status-filtered provider. A test that models the local provider's
// unfiltered load clears filterByStatus, so the page is returned whatever
// the connection's status is.
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
	filterByStatus   bool
}

func (p *notFoundProbeTestProvider) Name() string {
	return "test-provider"
}

func (p *notFoundProbeTestProvider) GetProviderURL() string {
	return "test://provider"
}

func (p *notFoundProbeTestProvider) GetProviderType() models.ProviderType {
	return models.LocalProviderType
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
	if p.filterByStatus && (p.connection == nil || string(p.connection.Status) != withStatus) {
		empty, err := json.Marshal(models.MesheryK8sContextPage{
			Page:       0,
			PageSize:   25,
			TotalCount: 0,
			Contexts:   []*models.K8sContext{},
		})
		return empty, err
	}
	return p.notFoundPage, nil
}

func (p *notFoundProbeTestProvider) GetConnectionByID(_ string, _ core.Uuid) (*connections.Connection, int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.connection, 0, nil
}

func (p *notFoundProbeTestProvider) UpdateConnectionById(_ string, conn *connections.ConnectionPayload, _ string) (*connections.Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// Mirror the write the real provider performs: the persisted status
	// becomes the state the machine settled in.
	updated := *p.connection
	updated.Status = conn.Status
	p.connection = &updated
	return p.connection, nil
}

// connectionStatus returns the fake connection's current status, serialized
// with the provider mutex: the probe goroutine writes it through
// UpdateConnectionById while the test reads it, so an unlocked read would
// be a data race.
func (p *notFoundProbeTestProvider) connectionStatus() connections.ConnectionStatus {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.connection.Status
}

// setConnectionStatus sets the fake connection's status, serialized with
// the provider mutex.
func (p *notFoundProbeTestProvider) setConnectionStatus(status connections.ConnectionStatus) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.connection.Status = status
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
		notFoundPage:   page,
		filterByStatus: true,
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
	runFSMWithConnected(t, h, provider, user, []*models.K8sContext{})
}

// runFSMWithConnected drives one K8sFSMMiddleware pass with the given
// connected contexts in AllKubeClusterKey.
func runFSMWithConnected(t *testing.T, h *Handler, provider models.Provider, user *models.User, connected []*models.K8sContext) {
	t.Helper()
	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "")
	ctx = context.WithValue(ctx, models.AllKubeClusterKey, connected)
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
	if got := provider.connectionStatus(); got != connections.NOTFOUND {
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
	if got := provider.connectionStatus(); got != connections.CONNECTED {
		t.Errorf("connection status after resolving probe = %q, want %q", got, connections.CONNECTED)
	}
	// The resolution is notified: the stored condition event is broadcast.
	broadcast := receiveEvent(t, broadcastCh)
	if broadcast.ResolvedAt == nil {
		t.Error("broadcast resolution has no ResolvedAt")
	}

	// The connection left NOTFOUND, so its throttle entry is forgotten.
	waitForThrottleForget(t, h.notFoundProbeThrottle, connectionID)

	// On the next load cycle the load still happens (the load itself is
	// still gated, so the clock must move past the interval first), but a
	// real status-filtered provider no longer returns a CONNECTED
	// connection in a NOTFOUND query — so no second probe fires and no
	// second event is persisted. The forgotten throttle entry is what the
	// forget above already asserted; the load returning an empty page is
	// what keeps the resolved connection from being probed again.
	clock.advance(notFoundProbeInterval + time.Second)
	runFSM(t, h, provider, user)
	if got := provider.getContextsCallCount(); got != 2 {
		t.Errorf("GetK8sContexts calls after a minute = %d, want 2 (load happened again)", got)
	}
	select {
	case <-probed:
		t.Error("resolved connection probed again after it left NOTFOUND")
	case <-time.After(500 * time.Millisecond):
	}
	if got := provider.persistedCount(); got != 1 {
		t.Errorf("persisted events = %d, want 1 (no second probe, no second event)", got)
	}
}

// TestNotFoundProbeThrottle is a unit test for the throttle itself: the load
// gate is per provider+user and the per-connection probe gate is atomic
// check-and-set, stamped at start, and forget/prune drop entries for
// connections that left NOTFOUND or were deleted.
func TestNotFoundProbeThrottle(t *testing.T) {
	clock := &fakeClock{now: time.Now()}
	throttle := newNotFoundProbeThrottle(clock.Now, time.Minute)
	id := uuid.Must(uuid.NewV4())
	other := uuid.Must(uuid.NewV4())
	scope := notFoundProbeLoadScope{providerName: "test-provider", providerURL: "test://provider", user: uuid.Must(uuid.NewV4())}
	otherScope := notFoundProbeLoadScope{providerName: "test-provider", providerURL: "test://provider", user: uuid.Must(uuid.NewV4())}

	if !throttle.allowLoad(scope) {
		t.Fatal("first load must be allowed")
	}
	if throttle.allowLoad(scope) {
		t.Error("second load inside the interval must be gated for the same scope")
	}
	if !throttle.allowLoad(otherScope) {
		t.Error("a different user must load independently of the first")
	}
	if throttle.allowLoad(otherScope) {
		t.Error("second load inside the interval must be gated for the other scope too")
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

	// forget drops one entry; prune drops only expired entries, so a fresh
	// probe entry survives a prune from another scope whose load did not
	// list that connection.
	throttle.forget(id)
	if !throttle.allowProbe(id) {
		t.Error("probe must be allowed again after forget")
	}
	throttle.prune()
	if throttle.allowProbe(id) {
		t.Error("prune must keep a fresh probe entry, even one absent from this scope's load")
	}
	if throttle.allowProbe(other) {
		t.Error("prune must keep a fresh probe entry for another connection too")
	}

	// Prune drops entries older than the interval, but keeps fresh ones:
	// scope loaded at the start, freshScope thirty seconds later, then
	// forty more seconds pass — scope's load entry is expired, freshScope's
	// is still within the interval. The probe entries are expired by then
	// too, so prune drops them as well.
	freshScope := notFoundProbeLoadScope{providerName: "test-provider", providerURL: "test://provider", user: uuid.Must(uuid.NewV4())}
	clock.advance(30 * time.Second)
	if !throttle.allowLoad(freshScope) {
		t.Fatal("a new scope must load independently")
	}
	clock.advance(40 * time.Second)
	throttle.prune()
	if !throttle.allowLoad(scope) {
		t.Error("load must be allowed again after prune dropped the expired entry")
	}
	if throttle.allowLoad(freshScope) {
		t.Error("prune must keep a load entry that is still within the interval")
	}
	if !throttle.allowProbe(id) {
		t.Error("probe must be allowed again after prune dropped the expired entry")
	}
	if !throttle.allowProbe(other) {
		t.Error("probe must be allowed again after prune dropped the expired entry")
	}

	// Advancing the clock past the interval re-allows both gates.
	clock.advance(time.Minute + time.Second)
	if !throttle.allowLoad(scope) {
		t.Error("load must be allowed again after the interval")
	}
	if !throttle.allowLoad(otherScope) {
		t.Error("load must be allowed again after the interval for the other scope")
	}
	if !throttle.allowLoad(freshScope) {
		t.Error("load must be allowed again after the interval for the fresh scope")
	}
	if !throttle.allowProbe(id) {
		t.Error("probe must be allowed again after the interval")
	}
	if !throttle.allowProbe(other) {
		t.Error("probe must be allowed again after the interval")
	}
}

// TestK8sFSMMiddlewareProbeLoadIsScopedPerUser checks the load gate is
// scoped to provider+user: a second user on the same Handler loads the
// NOTFOUND list again within the same minute, instead of being blocked by
// the first user's load. The per-connection probe limit stays global, so the
// second pass still does not re-probe the connection.
func TestK8sFSMMiddlewareProbeLoadIsScopedPerUser(t *testing.T) {
	probed := make(chan struct{}, 16)
	action := &probeCountingAction{probed: probed}
	h, provider, _, _, user, _ := newNotFoundProbeFixture(t, action)
	otherUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	runFSM(t, h, provider, user)
	waitForProbes(t, probed, 1)
	waitForPersisted(t, provider, 1)
	if got := provider.getContextsCallCount(); got != 1 {
		t.Fatalf("GetK8sContexts calls after first user = %d, want 1", got)
	}

	// A different user on the same Handler loads again immediately.
	runFSM(t, h, provider, otherUser)
	if got := provider.getContextsCallCount(); got != 2 {
		t.Errorf("GetK8sContexts calls after second user = %d, want 2 (load is scoped per user)", got)
	}
	// ...but the connection is not probed twice within the interval: the
	// per-connection probe limit is global.
	select {
	case <-probed:
		t.Error("connection probed twice within the throttle interval")
	case <-time.After(500 * time.Millisecond):
	}
	if got := provider.persistedCount(); got != 1 {
		t.Errorf("persisted events = %d, want 1 (no second probe)", got)
	}

	// And the same user again is still gated: no third load.
	runFSM(t, h, provider, user)
	if got := provider.getContextsCallCount(); got != 2 {
		t.Errorf("GetK8sContexts calls after first user again = %d, want 2 (same scope gated)", got)
	}
}

// TestK8sFSMMiddlewareSkipsDisconnectedConnection checks the local-provider
// guard: when the unfiltered load returns a context whose connection is
// DISCONNECTED, the probe loop skips it instead of probing — Discovery must
// not move a disconnected connection to CONNECTED without a user reconnect.
func TestK8sFSMMiddlewareSkipsDisconnectedConnection(t *testing.T) {
	probed := make(chan struct{}, 16)
	action := &probeCountingAction{probed: probed}
	h, provider, _, _, user, _ := newNotFoundProbeFixture(t, action)

	// Model the local provider's unfiltered load: the page is returned
	// whatever the connection's status is, and the connection is
	// DISCONNECTED — the state a user reconnect must clear, not a probe.
	provider.filterByStatus = false
	provider.setConnectionStatus(connections.DISCONNECTED)

	runFSM(t, h, provider, user)
	if got := provider.getContextsCallCount(); got != 1 {
		t.Fatalf("GetK8sContexts calls = %d, want 1 (the load happened)", got)
	}
	select {
	case <-probed:
		t.Error("DISCONNECTED connection was probed, want it skipped")
	case <-time.After(500 * time.Millisecond):
	}
	if got := provider.persistedCount(); got != 0 {
		t.Errorf("persisted events = %d, want 0 (no probe ran)", got)
	}
}

// contextCapturingAction records the Server of the machine context the probe
// runs with, then behaves like resolvingProbeAction: the cluster came back,
// so Discovery succeeds and the machine walks to CONNECTED.
type contextCapturingAction struct {
	probed chan struct{}
	seen   chan string
}

func (a *contextCapturingAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (a *contextCapturingAction) Execute(_ context.Context, machineCtx interface{}, _ interface{}) (machines.EventType, *events.Event, error) {
	if ctx, ok := machineCtx.(*kubernetes.MachineCtx); ok {
		select {
		case a.seen <- ctx.K8sContext.Server:
		default:
		}
	}
	select {
	case a.probed <- struct{}{}:
	default:
	}
	return machines.Register, nil, nil
}

func (a *contextCapturingAction) ExecuteOnExit(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

// setNotFoundPage replaces the fake provider's NOTFOUND page with one holding
// the given contexts.
func setNotFoundPage(t *testing.T, p *notFoundProbeTestProvider, contexts ...*models.K8sContext) {
	t.Helper()
	page, err := json.Marshal(models.MesheryK8sContextPage{
		Page:       0,
		PageSize:   25,
		TotalCount: len(contexts),
		Contexts:   contexts,
	})
	if err != nil {
		t.Fatalf("marshal NOTFOUND page: %v", err)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.notFoundPage = page
}

// TestK8sFSMMiddlewareProbeRefreshesChangedCredentials checks the credential
// refresh: the tracked machine was built with an old server and token, the
// load returns the same connection with a new server and token, and the
// probe must run against the new material — not the clientset built at
// machine creation.
func TestK8sFSMMiddlewareProbeRefreshesChangedCredentials(t *testing.T) {
	probed := make(chan struct{}, 16)
	seen := make(chan string, 1)
	action := &contextCapturingAction{probed: probed, seen: seen}
	h, provider, _, connectionID, user, _ := newNotFoundProbeFixture(t, action)

	// GenerateKubeConfig consumes the cluster/auth blocks as kubeconfig
	// named entries: the cluster block carries the server under a nested
	// "cluster" key and the auth block carries credentials under a nested
	// "user" key, so the clientset can be built from them.
	newCtx := &models.K8sContext{
		ID:           connectionID.String(),
		Name:         "dead-ctx",
		ConnectionID: connectionID.String(),
		Server:       "https://new:6443",
		Auth:         sql.Map{"name": "new-user", "user": sql.Map{"token": "new-token"}},
		Cluster:      sql.Map{"name": "new-cluster", "cluster": sql.Map{"server": "https://new:6443", "insecure-skip-tls-verify": true}},
	}
	setNotFoundPage(t, provider, newCtx)

	// The machine was built with the old connection material.
	machine, ok := h.ConnectionToStateMachineInstanceTracker.Get(connectionID)
	if !ok {
		t.Fatal("no tracked machine for the connection")
	}
	stored, ok := machine.Context.(*kubernetes.MachineCtx)
	if !ok {
		t.Fatal("tracked machine context is not a *kubernetes.MachineCtx")
	}
	stored.K8sContext = models.K8sContext{
		ID:           connectionID.String(),
		Name:         "dead-ctx",
		ConnectionID: connectionID.String(),
		Server:       "https://old:6443",
		Auth:         sql.Map{"name": "old-user", "user": sql.Map{"token": "old-token"}},
		Cluster:      sql.Map{"name": "old-cluster", "cluster": sql.Map{"server": "https://old:6443", "insecure-skip-tls-verify": true}},
	}

	runFSM(t, h, provider, user)
	waitForProbes(t, probed, 1)
	waitForPersisted(t, provider, 1)

	select {
	case server := <-seen:
		if server != "https://new:6443" {
			t.Errorf("probe ran against server %q, want %q (the refreshed credentials)", server, "https://new:6443")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("probe action never reported the server it ran against")
	}
	if got := provider.connectionStatus(); got != connections.CONNECTED {
		t.Errorf("connection status after resolving probe = %q, want %q", got, connections.CONNECTED)
	}
}

// TestK8sFSMMiddlewareProbeSkipsRefreshFailure checks the failed-refresh
// guard: when the newly loaded credentials cannot build a clientset, the
// probe is skipped instead of reporting recovery through the old clientset,
// which the stored credentials can no longer achieve.
func TestK8sFSMMiddlewareProbeSkipsRefreshFailure(t *testing.T) {
	probed := make(chan struct{}, 16)
	action := &probeCountingAction{probed: probed}
	h, provider, _, connectionID, user, _ := newNotFoundProbeFixture(t, action)

	// The new credentials carry CA data that is not valid base64, so no
	// clientset can be built from them.
	newCtx := &models.K8sContext{
		ID:           connectionID.String(),
		Name:         "dead-ctx",
		ConnectionID: connectionID.String(),
		Server:       "https://new:6443",
		Auth:         sql.Map{"name": "new-user", "user": sql.Map{"token": "new-token"}},
		Cluster:      sql.Map{"name": "new-cluster", "cluster": sql.Map{"server": "https://new:6443", "insecure-skip-tls-verify": true, "certificate-authority-data": "!!!not-base64!!!"}},
	}
	setNotFoundPage(t, provider, newCtx)

	// The machine was built with the old connection material.
	machine, ok := h.ConnectionToStateMachineInstanceTracker.Get(connectionID)
	if !ok {
		t.Fatal("no tracked machine for the connection")
	}
	stored, ok := machine.Context.(*kubernetes.MachineCtx)
	if !ok {
		t.Fatal("tracked machine context is not a *kubernetes.MachineCtx")
	}
	stored.K8sContext = models.K8sContext{
		ID:           connectionID.String(),
		Name:         "dead-ctx",
		ConnectionID: connectionID.String(),
		Server:       "https://old:6443",
		Auth:         sql.Map{"name": "old-user", "user": sql.Map{"token": "old-token"}},
		Cluster:      sql.Map{"name": "old-cluster", "cluster": sql.Map{"server": "https://old:6443", "insecure-skip-tls-verify": true}},
	}

	runFSM(t, h, provider, user)
	select {
	case <-probed:
		t.Error("probe ran after a failed credential refresh, want it skipped")
	case <-time.After(500 * time.Millisecond):
	}
	if got := provider.persistedCount(); got != 0 {
		t.Errorf("persisted events = %d, want 0 (no probe ran)", got)
	}
	if got := provider.connectionStatus(); got != connections.NOTFOUND {
		t.Errorf("connection status after skipped probe = %q, want %q", got, connections.NOTFOUND)
	}
}

// TestK8sFSMMiddlewareSkipsConnectedNotFoundContext checks the connectedIDs
// guard: when the NOTFOUND page holds a context that is also in
// AllKubeClusterKey — the local-provider case, where the provider returns
// contexts unfiltered by status, so one context can appear in both lists
// inside a single request — the NOTFOUND probe loop must not probe it.
//
// Note the exact assertion: the context IS probed once, through the normal
// connected path (it is in AllKubeClusterKey, so Discovery runs for it
// there, exactly like any connected context). What the guard prevents is a
// SECOND probe from the NOTFOUND loop for the same connection in the same
// request. So the test asserts exactly one probe and one persisted event,
// not zero.
func TestK8sFSMMiddlewareSkipsConnectedNotFoundContext(t *testing.T) {
	probed := make(chan struct{}, 16)
	action := &probeCountingAction{probed: probed}
	h, provider, _, connectionID, user, _ := newNotFoundProbeFixture(t, action)

	connected := []*models.K8sContext{{
		ID:           connectionID.String(),
		Name:         "dead-ctx",
		ConnectionID: connectionID.String(),
	}}
	runFSMWithConnected(t, h, provider, user, connected)

	// The NOTFOUND list still loads once.
	if got := provider.getContextsCallCount(); got != 1 {
		t.Fatalf("GetK8sContexts calls = %d, want 1 (the load happened)", got)
	}
	// Exactly one probe: the connected path's. A second probe would be the
	// NOTFOUND loop probing the same connection again.
	waitForProbes(t, probed, 1)
	waitForPersisted(t, provider, 1)
	select {
	case <-probed:
		t.Error("connection probed twice in one request: the NOTFOUND loop probed a context already in AllKubeClusterKey")
	case <-time.After(500 * time.Millisecond):
	}
	if got := provider.persistedCount(); got != 1 {
		t.Errorf("persisted events = %d, want 1 (one probe, one Recurrence)", got)
	}
}
