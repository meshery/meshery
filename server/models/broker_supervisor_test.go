package models

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	meshkitBroker "github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/broker/nats"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/controllers"
	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// liveBroker is the cluster's current Meshery Broker as its controller reports
// it: the endpoint the Broker resource publishes and the token in
// meshery-nats-auth. A test "recreates the broker" by changing both. While up,
// it answers at that endpoint and nowhere else.
type liveBroker struct {
	stubController
	mu       sync.Mutex
	endpoint string
	token    string
	up       bool
	// acceptedToken is the token the broker actually accepts. It usually
	// tracks token, but a test separates them to model a broker that rejects
	// the token Meshery resolved (for example, the secret was rotated again
	// after Meshery read it).
	acceptedToken string
	// onResolve runs inside GetPublicEndpoint, after the endpoint is read.
	// A test uses it to change the world while the supervisor is
	// re-resolving, between its stale check and its rebuild.
	onResolve func()
}

// accepts reports whether the broker accepts a token for a new connection.
func (b *liveBroker) accepts(token string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return token == b.acceptedToken
}

// setAcceptedToken changes which token the broker accepts, without changing
// what its controller reports.
func (b *liveBroker) setAcceptedToken(token string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.acceptedToken = token
}

func (b *liveBroker) answers(endpoint string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.up && endpoint == b.endpoint
}

func (b *liveBroker) setUp(v bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.up = v
}

func (b *liveBroker) GetPublicEndpoint() (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.endpoint == "" {
		return "", errors.New("broker has no published endpoint")
	}
	if b.onResolve != nil {
		b.onResolve()
	}
	return b.endpoint, nil
}

func (b *liveBroker) GetToken() (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.token, nil
}

func (b *liveBroker) recreate(endpoint, token string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.endpoint, b.token = endpoint, token
	b.acceptedToken = token
}

// natsConn is one NATS connection Meshery opened, standing in for the real
// client: it never connects anywhere, and reports whatever the test says.
type natsConn struct {
	*testMeshsyncBroker
	opts       nats.Options
	mu         sync.Mutex
	connected  bool
	publishErr error
}

func (c *natsConn) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected && !c.isClosed()
}

func (c *natsConn) setConnected(v bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connected = v
}

func (c *natsConn) Publish(string, *meshkitBroker.Message) error { return c.publishErr }

// brokerDialer records every connection Meshery opens to the broker.
type brokerDialer struct {
	mu         sync.Mutex
	conns      []*natsConn
	publishErr error
	// acceptToken reports whether the broker accepts a token for a new
	// connection. A nil func accepts every token. A connection dialed with an
	// accepted token reports connected, standing in for the NATS client
	// connecting; one dialed with a rejected token never connects.
	acceptToken func(string) bool
}

func (d *brokerDialer) dial(opts nats.Options) (meshkitBroker.Handler, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	c := &natsConn{testMeshsyncBroker: &testMeshsyncBroker{}, opts: opts, publishErr: d.publishErr}
	if d.acceptToken == nil || d.acceptToken(opts.Token) {
		c.connected = true
	}
	d.conns = append(d.conns, c)
	return c, nil
}

func (d *brokerDialer) opened() []*natsConn {
	d.mu.Lock()
	defer d.mu.Unlock()
	return slices.Clone(d.conns)
}

const (
	oldBrokerEndpoint = "10.0.0.10:4222"
	newBrokerEndpoint = "10.0.0.20:4222"
)

// supervisedConnection attaches an operator-mode MeshSync data handler whose
// broker connection is recorded by the returned dialer. The supervisor's own
// ticker is parked; tests drive recheckBroker directly with a fake clock.
func supervisedConnection(t *testing.T) (*MesheryControllersHelper, *liveBroker, *brokerDialer) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	mch := newTestControllersHelper(t)
	mch.dbHandler = &database.Handler{DB: db}
	mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
	mch.hostCluster = newHostCluster("", "", nil)
	mch.brokerRecheckInterval = time.Hour
	mch.brokerStaleAfter = time.Minute

	broker := &liveBroker{endpoint: oldBrokerEndpoint, token: "old-token", acceptedToken: "old-token"}
	mch.ctxControllerHandlers = map[MesheryController]controllers.IMesheryController{MesheryBroker: broker}
	mch.brokerReachable = broker.answers
	dialer := &brokerDialer{acceptToken: broker.accepts}
	mch.newBrokerConnection = dialer.dial
	// A rebuild waits only briefly for its replacement to connect, so a test
	// that rejects the replacement does not sit out the production wait.
	mch.brokerConnectWait = 50 * time.Millisecond
	t.Cleanup(func() { mch.RemoveMeshSyncDataHandler(context.Background(), "ctx-1") })
	return mch, broker, dialer
}

func attach(t *testing.T, mch *MesheryControllersHelper) *brokerTarget {
	t.Helper()
	mch.AddMeshsyncDataHandlers(context.Background(), K8sContext{ID: "ctx-1", Name: "remote"}, core.Uuid(uuid.Nil), core.Uuid(uuid.Nil), nil)
	mch.meshsyncMutex.Lock()
	defer mch.meshsyncMutex.Unlock()
	if mch.ctxMeshsyncDataHandler == nil || mch.brokerTarget == nil {
		t.Fatal("no supervised MeshSync data handler was attached")
	}
	return mch.brokerTarget
}

func stopped(target *brokerTarget) bool {
	select {
	case <-target.stop:
		return true
	default:
		return false
	}
}

// TestRecreatedBrokerIsReconnectedWithItsNewEndpointAndToken is defect C: the
// broker was recreated with a new endpoint and token, and the handler kept
// retrying the old ones (the token rejection is terminal for the NATS client)
// until the server was restarted.
func TestRecreatedBrokerIsReconnectedWithItsNewEndpointAndToken(t *testing.T) {
	mch, broker, dialer := supervisedConnection(t)
	target := attach(t, mch)
	first := dialer.opened()[0]
	if first.opts.Token != "old-token" || !slices.Contains(first.opts.URLS, oldBrokerEndpoint) {
		t.Fatalf("first connection built from %v / %q, want the broker as resolved", first.opts.URLS, first.opts.Token)
	}

	broker.recreate(newBrokerEndpoint, "new-token")
	// The recreated broker drops the old connection.
	first.setConnected(false)
	var since time.Time
	t0 := time.Unix(1000, 0)
	if mch.recheckBroker(target, t0, &since) {
		t.Fatal("the supervisor finished on the first sign of an outage")
	}
	if mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter-time.Second), &since) {
		t.Fatal("the supervisor rebuilt before the outage was sustained")
	}
	if n := len(dialer.opened()); n != 1 {
		t.Fatalf("%d connections opened before the outage was sustained, want 1", n)
	}

	if !mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter), &since) {
		t.Fatal("a sustained outage against a recreated broker did not rebuild the connection")
	}
	conns := dialer.opened()
	if len(conns) != 2 {
		t.Fatalf("%d connections opened, want the original and one rebuilt", len(conns))
	}
	rebuilt := conns[1]
	if rebuilt.opts.Token != "new-token" || !slices.Contains(rebuilt.opts.URLS, newBrokerEndpoint) {
		t.Fatalf("rebuilt connection uses %v / %q, want the recreated broker's endpoint and token", rebuilt.opts.URLS, rebuilt.opts.Token)
	}
	if !first.isClosed() {
		t.Fatal("the stale connection was left reconnecting in the background")
	}
	if !stopped(target) {
		t.Fatal("the stale connection's supervisor was not stopped")
	}
	mch.meshsyncMutex.Lock()
	defer mch.meshsyncMutex.Unlock()
	if mch.ctxMeshsyncDataHandler == nil || mch.ctxMeshsyncDataHandler.GetBrokerHandler() != meshkitBroker.Handler(rebuilt) {
		t.Fatal("the attached data handler is not the rebuilt one")
	}
	if mch.brokerTarget == nil || mch.brokerTarget == target || mch.brokerTarget.token != "new-token" {
		t.Fatal("the rebuilt connection is not supervised against the new broker")
	}
}

// TestMovedBrokerWithTheSameTokenIsReconnected: a broker that answers at a
// different endpoint while the connection is down has moved, even though its
// token did not change.
func TestMovedBrokerWithTheSameTokenIsReconnected(t *testing.T) {
	mch, broker, dialer := supervisedConnection(t)
	target := attach(t, mch)
	first := dialer.opened()[0]
	broker.recreate(newBrokerEndpoint, "old-token")
	broker.setUp(true)
	// The recreated broker drops the old connection.
	first.setConnected(false)

	var since time.Time
	t0 := time.Unix(1000, 0)
	mch.recheckBroker(target, t0, &since)
	if !mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter), &since) {
		t.Fatal("a sustained outage against a moved broker did not rebuild the connection")
	}
	conns := dialer.opened()
	if len(conns) != 2 || !slices.Contains(conns[1].opts.URLS, newBrokerEndpoint) {
		t.Fatalf("connections opened: %d, want one rebuilt against the moved broker's endpoint", len(conns))
	}
}

// TestBrokerSupervisorLeavesWhatItShouldNotTouch covers every case in which
// the supervisor must not rebuild: the NATS client recovers a merely-down
// broker on its own, and tearing down a live connection loses data.
func TestBrokerSupervisorLeavesWhatItShouldNotTouch(t *testing.T) {
	tests := []struct {
		name string
		// arrange runs after attach and before the clock passes staleAfter.
		arrange func(*MesheryControllersHelper, *liveBroker, *natsConn)
	}{
		{
			name:    "broker down, same endpoint and token",
			arrange: func(*MesheryControllersHelper, *liveBroker, *natsConn) {},
		},
		{
			// A NodePort, LoadBalancer or host.docker.internal connection: once
			// nothing answers, the resolved endpoint falls back to the internal
			// one although the broker has not moved.
			name: "broker down, resolves to a different candidate, same token",
			arrange: func(_ *MesheryControllersHelper, b *liveBroker, _ *natsConn) {
				b.recreate(newBrokerEndpoint, "old-token")
			},
		},
		{
			name: "connection live, broker token rotated",
			arrange: func(_ *MesheryControllersHelper, b *liveBroker, c *natsConn) {
				c.setConnected(true)
				b.recreate(newBrokerEndpoint, "new-token")
			},
		},
		{
			name: "new broker's token not provisioned yet",
			arrange: func(_ *MesheryControllersHelper, b *liveBroker, _ *natsConn) {
				b.recreate(newBrokerEndpoint, "")
			},
		},
		{
			name: "no endpoint resolvable, same token",
			arrange: func(_ *MesheryControllersHelper, b *liveBroker, _ *natsConn) {
				b.recreate("", "old-token")
			},
		},
		{
			// The new token is provisioned before the new Broker publishes an
			// endpoint: rebuilding now would find nothing to connect through.
			name: "new token, no endpoint published yet",
			arrange: func(_ *MesheryControllersHelper, b *liveBroker, _ *natsConn) {
				b.recreate("", "new-token")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mch, broker, dialer := supervisedConnection(t)
			target := attach(t, mch)
			tt.arrange(mch, broker, dialer.opened()[0])

			var since time.Time
			t0 := time.Unix(1000, 0)
			for _, at := range []time.Time{t0, t0.Add(mch.brokerStaleAfter), t0.Add(10 * mch.brokerStaleAfter)} {
				if mch.recheckBroker(target, at, &since) {
					t.Fatal("the supervisor finished although its connection was neither replaced nor removed")
				}
			}
			if n := len(dialer.opened()); n != 1 {
				t.Fatalf("%d connections opened, want the original left alone", n)
			}
		})
	}
}

// TestFailedRebuildKeepsTheConnectionSupervised: a rebuild that attaches no
// handler must not leave the connection without one, nor without a supervisor
// to try again. It must also not announce the rebuild or reset the connected
// event flag: those belong to a rebuild that actually replaced the handler,
// or a failing rebuild repeats the "Broker changed" event on every retry tick.
func TestFailedRebuildKeepsTheConnectionSupervised(t *testing.T) {
	mch, broker, dialer := supervisedConnection(t)
	target := attach(t, mch)
	first := dialer.opened()[0]
	// The connected event was already emitted for this handler session.
	mch.meshsyncConnectedEventEmitted.Store(true)
	broker.recreate(newBrokerEndpoint, "new-token")
	// The recreated broker drops the old connection.
	first.setConnected(false)
	dialer.mu.Lock()
	dialer.publishErr = errors.New("publish failed")
	dialer.mu.Unlock()

	var since time.Time
	t0 := time.Unix(1000, 0)
	mch.recheckBroker(target, t0, &since)
	if mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter), &since) {
		t.Fatal("the supervisor finished although its rebuild attached nothing")
	}
	if !mch.meshsyncConnectedEventEmitted.Load() {
		t.Fatal("a failed rebuild reset the connected event flag, so the stale handler's connected event would be emitted again")
	}
	conns := dialer.opened()
	if len(conns) != 2 || !conns[1].isClosed() {
		t.Fatalf("%d connections opened, want one failed rebuild, closed", len(conns))
	}
	if first.isClosed() || stopped(target) {
		t.Fatal("a failed rebuild tore down the connection it was replacing")
	}
	mch.meshsyncMutex.Lock()
	kept := mch.brokerTarget == target && mch.ctxMeshsyncDataHandler != nil &&
		mch.ctxMeshsyncDataHandler.GetBrokerHandler() == meshkitBroker.Handler(first)
	mch.meshsyncMutex.Unlock()
	if !kept {
		t.Fatal("a failed rebuild left the connection without its handler or supervisor")
	}

	dialer.mu.Lock()
	dialer.publishErr = nil
	dialer.mu.Unlock()
	// Clear the flag so the retry's emission is observable: after the retry
	// the flag must be true because the rebuild emitted the connected event
	// for the new handler, not because it was left over from the old session.
	mch.meshsyncConnectedEventEmitted.Store(false)
	if !mch.recheckBroker(target, t0.Add(2*mch.brokerStaleAfter), &since) {
		t.Fatal("the supervisor did not retry the rebuild")
	}
	if !mch.meshsyncConnectedEventEmitted.Load() {
		t.Fatal("a successful rebuild did not emit the connected event for the new handler")
	}
	conns = dialer.opened()
	if len(conns) != 3 || conns[2].opts.Token != "new-token" {
		t.Fatalf("%d connections opened, want a retried rebuild with the new token", len(conns))
	}
	if !first.isClosed() || !stopped(target) {
		t.Fatal("the replaced connection or its supervisor outlived a successful rebuild")
	}
}

// TestRejectedTokenKeepsTheOldHandler: the broker was recreated and Meshery
// resolved a new token, but the broker rejects it - the secret was rotated
// again after Meshery read it. The replacement attaches but never connects,
// so the rebuild must discard it and keep the old handler and its supervisor:
// stopping the old handler would leave the connection down until a restart,
// because later ticks see the same endpoint and token and never rebuild again.
func TestRejectedTokenKeepsTheOldHandler(t *testing.T) {
	mch, broker, dialer := supervisedConnection(t)
	target := attach(t, mch)
	first := dialer.opened()[0]
	// The connected event was already emitted for this handler session.
	mch.meshsyncConnectedEventEmitted.Store(true)
	broker.recreate(newBrokerEndpoint, "new-token")
	// The broker rejects the resolved token.
	broker.setAcceptedToken("old-token")
	// The recreated broker drops the old connection.
	first.setConnected(false)

	var since time.Time
	t0 := time.Unix(1000, 0)
	mch.recheckBroker(target, t0, &since)
	if mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter), &since) {
		t.Fatal("the supervisor finished although its rebuild never connected")
	}
	conns := dialer.opened()
	if len(conns) != 2 {
		t.Fatalf("%d connections opened, want the original and one rejected rebuild", len(conns))
	}
	if !conns[1].isClosed() {
		t.Fatal("the replacement the broker rejected was left reconnecting in the background")
	}
	if first.isClosed() || stopped(target) {
		t.Fatal("a rejected rebuild tore down the connection it was replacing")
	}
	if !mch.meshsyncConnectedEventEmitted.Load() {
		t.Fatal("a rejected rebuild reset the connected event flag, so the kept handler's connected event would be emitted again")
	}
	mch.meshsyncMutex.Lock()
	kept := mch.brokerTarget == target && mch.ctxMeshsyncDataHandler != nil &&
		mch.ctxMeshsyncDataHandler.GetBrokerHandler() == meshkitBroker.Handler(first)
	mch.meshsyncMutex.Unlock()
	if !kept {
		t.Fatal("a rejected rebuild left the connection without its handler or supervisor")
	}

	// The next tick retries the rebuild instead of leaving the connection down.
	if mch.recheckBroker(target, t0.Add(2*mch.brokerStaleAfter), &since) {
		t.Fatal("the supervisor finished although its rebuild never connected")
	}
	conns = dialer.opened()
	if len(conns) != 3 || !conns[2].isClosed() {
		t.Fatalf("%d connections opened, want a retried rebuild, rejected and closed", len(conns))
	}
	if first.isClosed() || stopped(target) {
		t.Fatal("a retried rejected rebuild tore down the connection it was replacing")
	}
	mch.meshsyncMutex.Lock()
	kept = mch.brokerTarget == target && mch.ctxMeshsyncDataHandler != nil &&
		mch.ctxMeshsyncDataHandler.GetBrokerHandler() == meshkitBroker.Handler(first)
	mch.meshsyncMutex.Unlock()
	if !kept {
		t.Fatal("a retried rejected rebuild left the connection without its handler or supervisor")
	}
}

// TestRebuildPausesTheOldHandlerUntilRollback: during the replacement wait the
// old handler is detached but not stopped, and its NATS client reconnects
// forever - if it kept its MeshSync subscriptions it would receive the same
// messages as the replacement and persist them twice. The rebuild must pause
// the old handler's subscriptions for the wait and resume them only when the
// rebuild rolls back.
func TestRebuildPausesTheOldHandlerUntilRollback(t *testing.T) {
	mch, broker, dialer := supervisedConnection(t)
	// The wait runs long enough for the test to observe the paused state.
	mch.brokerConnectWait = 2 * time.Second
	target := attach(t, mch)
	first := dialer.opened()[0]
	// The broker rejects the resolved token, so the replacement never
	// connects and the rebuild rolls back to the old handler.
	broker.recreate(newBrokerEndpoint, "new-token")
	broker.setAcceptedToken("old-token")
	// The recreated broker drops the old connection.
	first.setConnected(false)

	var since time.Time
	t0 := time.Unix(1000, 0)
	mch.recheckBroker(target, t0, &since)
	done := make(chan bool, 1)
	go func() { done <- mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter), &since) }()

	// The rebuild is waiting for its replacement to connect: the old handler
	// must have no MeshSync subscription in that window.
	paused := false
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if !first.subscribed(meshsyncEventsSubject) && !first.subscribed(MeshsyncStoreUpdatesSubject) {
			paused = true
			break
		}
		time.Sleep(time.Millisecond)
	}
	if !paused {
		t.Fatal("the old handler kept its MeshSync subscriptions during the replacement wait")
	}
	if finished := <-done; finished {
		t.Fatal("the supervisor finished although its rebuild never connected")
	}
	// The rebuild rolled back: the old handler is attached again, and its
	// subscriptions are back so it receives MeshSync data.
	if !first.subscribed(meshsyncEventsSubject) || !first.subscribed(MeshsyncStoreUpdatesSubject) {
		t.Fatal("the rollback did not resume the old handler's MeshSync subscriptions")
	}
	if first.isClosed() || stopped(target) {
		t.Fatal("a rejected rebuild tore down the connection it was replacing")
	}
	mch.meshsyncMutex.Lock()
	kept := mch.brokerTarget == target && mch.ctxMeshsyncDataHandler != nil &&
		mch.ctxMeshsyncDataHandler.GetBrokerHandler() == meshkitBroker.Handler(first)
	mch.meshsyncMutex.Unlock()
	if !kept {
		t.Fatal("a rejected rebuild left the connection without its handler or supervisor")
	}
}

// TestRemovalDuringRebuildStopsBothHandlers: if the data handler is removed
// while a rebuild waits for its replacement to connect, removal must cancel
// the wait and stop both handlers - the old handler is only stopped when the
// wait expires otherwise, and a connection that recovers in that window would
// persist updates after removal returned.
func TestRemovalDuringRebuildStopsBothHandlers(t *testing.T) {
	mch, broker, dialer := supervisedConnection(t)
	// The wait would run long without a cancellation, so a removal that
	// waits it out is observable.
	mch.brokerConnectWait = 5 * time.Second
	target := attach(t, mch)
	first := dialer.opened()[0]
	broker.recreate(newBrokerEndpoint, "new-token")
	// The replacement never connects, so the rebuild stays in its wait.
	broker.setAcceptedToken("old-token")
	// The recreated broker drops the old connection.
	first.setConnected(false)

	var since time.Time
	t0 := time.Unix(1000, 0)
	mch.recheckBroker(target, t0, &since)
	done := make(chan bool, 1)
	go func() { done <- mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter), &since) }()

	// Wait until the rebuild attached its replacement and is waiting.
	deadline := time.Now().Add(2 * time.Second)
	for len(dialer.opened()) < 2 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if n := len(dialer.opened()); n != 2 {
		t.Fatalf("%d connections opened, want the original and one rebuilding", n)
	}

	started := time.Now()
	mch.RemoveMeshSyncDataHandler(context.Background(), "ctx-1")
	if elapsed := time.Since(started); elapsed > 2*time.Second {
		t.Fatalf("removal took %s, want it to cancel the replacement wait instead of waiting it out", elapsed)
	}
	replacement := dialer.opened()[1]
	if !first.isClosed() {
		t.Fatal("removal returned with the old handler still running")
	}
	if !replacement.isClosed() {
		t.Fatal("removal returned with the replacement handler still running")
	}
	mch.meshsyncMutex.Lock()
	attached := mch.ctxMeshsyncDataHandler
	mch.meshsyncMutex.Unlock()
	if attached != nil {
		t.Fatal("removal left a data handler attached")
	}

	select {
	case finished := <-done:
		if !finished {
			t.Fatal("the interrupted rebuild did not finish")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the interrupted rebuild did not return after removal")
	}
}

// TestRecreatedBrokerWithNewTokenAndNoEndpointRebuildsOverPortForward: the
// broker was recreated with a new token but publishes no endpoint, and the
// connection is carried by a managed port-forward. The token change must still
// rebuild the connection, through the forward - without one, the same change
// is (correctly) left alone; see "new token, no endpoint published yet" in
// TestBrokerSupervisorLeavesWhatItShouldNotTouch.
func TestRecreatedBrokerWithNewTokenAndNoEndpointRebuildsOverPortForward(t *testing.T) {
	mch, broker, dialer := supervisedConnection(t)
	// Stand in a managed port-forward. Only LocalAddr is ever read on this
	// path - GetBrokerPortForwardAddr, and ensureBrokerPortForward returns it
	// without starting anything while one is recorded - so a zero-value
	// forwarder carries the test without a Kubernetes client. Nil it before
	// teardown so RemoveMeshSyncDataHandler never calls Stop on it.
	mch.brokerPortForward = &mesherykube.PortForwarder{}
	t.Cleanup(func() { mch.brokerPortForward = nil })

	target := attach(t, mch)
	first := dialer.opened()[0]
	if !slices.Contains(first.opts.URLS, "127.0.0.1:0") {
		t.Fatalf("first connection built from %v, want it carried by the port-forward", first.opts.URLS)
	}
	broker.recreate("", "new-token")
	// The recreated broker drops the old connection.
	first.setConnected(false)

	var since time.Time
	t0 := time.Unix(1000, 0)
	mch.recheckBroker(target, t0, &since)
	if !mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter), &since) {
		t.Fatal("a sustained outage against a recreated broker did not rebuild the connection although a port-forward carries it")
	}
	conns := dialer.opened()
	if len(conns) != 2 {
		t.Fatalf("%d connections opened, want the original and one rebuilt", len(conns))
	}
	rebuilt := conns[1]
	if rebuilt.opts.Token != "new-token" || !slices.Contains(rebuilt.opts.URLS, "127.0.0.1:0") {
		t.Fatalf("rebuilt connection uses %v / %q, want the new token through the port-forward", rebuilt.opts.URLS, rebuilt.opts.Token)
	}
	if !first.isClosed() {
		t.Fatal("the stale connection was left reconnecting in the background")
	}
	if !stopped(target) {
		t.Fatal("the stale connection's supervisor was not stopped")
	}
	mch.meshsyncMutex.Lock()
	defer mch.meshsyncMutex.Unlock()
	if mch.ctxMeshsyncDataHandler == nil || mch.ctxMeshsyncDataHandler.GetBrokerHandler() != meshkitBroker.Handler(rebuilt) {
		t.Fatal("the attached data handler is not the rebuilt one")
	}
	if mch.brokerTarget == nil || mch.brokerTarget == target || mch.brokerTarget.token != "new-token" {
		t.Fatal("the rebuilt connection is not supervised against the new broker")
	}
}

// TestOutageClockResetsOnReconnect: a connection that flaps is not a sustained
// outage, so time spent connected must not count toward staleAfter.
func TestOutageClockResetsOnReconnect(t *testing.T) {
	mch, broker, dialer := supervisedConnection(t)
	target := attach(t, mch)
	conn := dialer.opened()[0]
	broker.recreate(newBrokerEndpoint, "new-token")

	var since time.Time
	t0 := time.Unix(1000, 0)
	mch.recheckBroker(target, t0, &since)
	conn.setConnected(true)
	mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter/2), &since)
	conn.setConnected(false)
	mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter), &since)
	if mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter+mch.brokerStaleAfter/2), &since) {
		t.Fatal("an outage interrupted by a reconnect was treated as sustained")
	}
	if n := len(dialer.opened()); n != 1 {
		t.Fatalf("%d connections opened, want no rebuild yet", n)
	}
}

// TestRecoveredConnectionIsNotStoppedByASustainedOutage: the NATS client's own
// reconnect loop can restore the connection while the supervisor is
// re-resolving the endpoint and token of a recreated broker. The rebuild must
// not stop a handler that is live again.
func TestRecoveredConnectionIsNotStoppedByASustainedOutage(t *testing.T) {
	mch, broker, dialer := supervisedConnection(t)
	target := attach(t, mch)
	first := dialer.opened()[0]
	broker.recreate(newBrokerEndpoint, "new-token")
	// The connection recovers while the supervisor re-resolves the endpoint
	// and token, after the stale check and before the rebuild.
	broker.onResolve = func() { first.setConnected(true) }

	var since time.Time
	t0 := time.Unix(1000, 0)
	if mch.recheckBroker(target, t0, &since) {
		t.Fatal("the supervisor finished on the first sign of an outage")
	}
	if mch.recheckBroker(target, t0.Add(mch.brokerStaleAfter), &since) {
		t.Fatal("the supervisor finished although the connection had recovered")
	}
	if !since.IsZero() {
		t.Fatal("the outage clock was not reset when the connection recovered")
	}
	if n := len(dialer.opened()); n != 1 {
		t.Fatalf("%d connections opened, want the recovered one left alone", n)
	}
	if first.isClosed() || stopped(target) {
		t.Fatal("a recovered connection or its supervisor was stopped")
	}
	mch.meshsyncMutex.Lock()
	defer mch.meshsyncMutex.Unlock()
	if mch.ctxMeshsyncDataHandler == nil || mch.brokerTarget != target {
		t.Fatal("a recovered connection lost its handler or supervisor")
	}
}

// TestRemovedHandlerEndsItsSupervisor: disconnect and delete remove the data
// handler; its supervisor must neither outlive it nor rebuild it.
func TestRemovedHandlerEndsItsSupervisor(t *testing.T) {
	mch, broker, dialer := supervisedConnection(t)
	target := attach(t, mch)
	broker.recreate(newBrokerEndpoint, "new-token")

	mch.RemoveMeshSyncDataHandler(context.Background(), "ctx-1")

	if !stopped(target) {
		t.Fatal("removing the data handler left its supervisor running")
	}
	var since time.Time
	if !mch.recheckBroker(target, time.Unix(1000, 0), &since) {
		t.Fatal("a supervisor whose handler was removed did not finish")
	}
	if n := len(dialer.opened()); n != 1 {
		t.Fatalf("%d connections opened, want no rebuild of a removed handler", n)
	}
}

// TestSupervisorRebuildsOnItsOwnTicker drives the real goroutine end to end.
func TestSupervisorRebuildsOnItsOwnTicker(t *testing.T) {
	mch, broker, dialer := supervisedConnection(t)
	mch.brokerRecheckInterval = 5 * time.Millisecond
	mch.brokerStaleAfter = 20 * time.Millisecond
	attach(t, mch)
	first := dialer.opened()[0]
	broker.recreate(newBrokerEndpoint, "new-token")
	// The recreated broker drops the old connection.
	first.setConnected(false)

	deadline := time.Now().Add(5 * time.Second)
	for len(dialer.opened()) < 2 {
		if time.Now().After(deadline) {
			t.Fatal("the supervisor never rebuilt the connection to the recreated broker")
		}
		time.Sleep(5 * time.Millisecond)
	}
	if got := dialer.opened()[1].opts.Token; got != "new-token" {
		t.Fatalf("rebuilt connection token = %q, want the new one", got)
	}
}

// TestFailedDataHandlerStartClosesItsBrokerConnection: a handler that fails to
// start is discarded, and its connection - which reconnects forever - used to
// be discarded with it still running, one more per connect cycle.
func TestFailedDataHandlerStartClosesItsBrokerConnection(t *testing.T) {
	mch, _, dialer := supervisedConnection(t)
	dialer.publishErr = errors.New("publish failed")

	mch.AddMeshsyncDataHandlers(context.Background(), K8sContext{ID: "ctx-1"}, core.Uuid(uuid.Nil), core.Uuid(uuid.Nil), nil)

	conns := dialer.opened()
	if len(conns) != 1 {
		t.Fatalf("%d connections opened, want 1", len(conns))
	}
	if !conns[0].isClosed() {
		t.Fatal("the broker connection of a handler that failed to start was left open")
	}
	mch.meshsyncMutex.Lock()
	defer mch.meshsyncMutex.Unlock()
	if mch.ctxMeshsyncDataHandler != nil || mch.brokerTarget != nil {
		t.Fatal("a handler that failed to start was attached or supervised")
	}
}
