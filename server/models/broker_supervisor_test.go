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
	"github.com/meshery/schemas/models/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// liveBroker is the cluster's current Meshery Broker as its controller reports
// it: the endpoint the Broker resource publishes and the token in
// meshery-nats-auth. A test "recreates the broker" by changing both.
type liveBroker struct {
	stubController
	mu       sync.Mutex
	endpoint string
	token    string
}

func (b *liveBroker) GetPublicEndpoint() (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.endpoint == "" {
		return "", errors.New("broker has no published endpoint")
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
}

func (d *brokerDialer) dial(opts nats.Options) (meshkitBroker.Handler, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	c := &natsConn{testMeshsyncBroker: &testMeshsyncBroker{}, opts: opts, publishErr: d.publishErr}
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

	broker := &liveBroker{endpoint: oldBrokerEndpoint, token: "old-token"}
	mch.ctxControllerHandlers = map[MesheryController]controllers.IMesheryController{MesheryBroker: broker}
	dialer := &brokerDialer{}
	mch.newBrokerConnection = dialer.dial
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
	broker.recreate(newBrokerEndpoint, "new-token")

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
