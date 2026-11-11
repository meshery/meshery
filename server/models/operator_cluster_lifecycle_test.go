package models

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// fakeCluster is the one piece of state every connection to a cluster shares:
// whether the Meshery Operator release is installed there. Each connection gets
// its own clusterOperator handler, exactly as each connection's helper builds
// its own meshkit handler, but they all read and write this.
type fakeCluster struct {
	mu          sync.Mutex
	status      controllers.MesheryControllerStatus
	installs    int
	uninstalls  int
	inFlight    int
	maxInFlight int
	// hold, when set, blocks every install until it is closed.
	hold chan struct{}
}

func newFakeCluster() *fakeCluster {
	return &fakeCluster{status: controllers.NotDeployed}
}

func (c *fakeCluster) snapshot() (status controllers.MesheryControllerStatus, installs, uninstalls, maxInFlight int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status, c.installs, c.uninstalls, c.maxInFlight
}

// helmAction records an overlapping Helm action the way the real release
// store would reject one: two at once against one release is the race.
func (c *fakeCluster) helmAction(apply func()) {
	c.mu.Lock()
	c.inFlight++
	if c.inFlight > c.maxInFlight {
		c.maxInFlight = c.inFlight
	}
	hold := c.hold
	c.mu.Unlock()

	if hold != nil {
		<-hold
	} else {
		// Long enough for a second, unserialized caller to overlap.
		time.Sleep(20 * time.Millisecond)
	}

	c.mu.Lock()
	apply()
	c.inFlight--
	c.mu.Unlock()
}

// clusterOperator is a meshkit operator handler backed by a fakeCluster. Like
// meshkit's, it never touches the network.
type clusterOperator struct{ cluster *fakeCluster }

func (o *clusterOperator) GetName() string { return "MesheryOperator" }
func (o *clusterOperator) GetStatus() controllers.MesheryControllerStatus {
	o.cluster.mu.Lock()
	defer o.cluster.mu.Unlock()
	return o.cluster.status
}
func (o *clusterOperator) Deploy(bool) error {
	o.cluster.helmAction(func() {
		o.cluster.installs++
		o.cluster.status = controllers.Deployed
	})
	return nil
}
func (o *clusterOperator) Undeploy() error {
	o.cluster.helmAction(func() {
		o.cluster.uninstalls++
		o.cluster.status = controllers.NotDeployed
	})
	return nil
}
func (o *clusterOperator) GetPublicEndpoint() (string, error)        { return "", nil }
func (o *clusterOperator) GetVersion() (string, error)               { return "v1.0.64", nil }
func (o *clusterOperator) GetEndpointForPort(string) (string, error) { return "", nil }

// connectionTo builds the K8sContext of a connection reaching the cluster
// identified by serverID. Each connection names the API server differently, so
// only the server ID can tell that they share a cluster - as with a load
// balancer address and a node address for one control plane.
func connectionTo(serverID core.Uuid, name string) K8sContext {
	id := serverID
	return K8sContext{
		ID:                 "ctx-" + name,
		Name:               name,
		ConnectionID:       uuid.Must(uuid.NewV4()).String(),
		Server:             "https://" + name + ".example.invalid:6443",
		KubernetesServerID: &id,
	}
}

// operatorModeConsumer is a connection's helper as the connect chain leaves it
// just before DeployUndeployedOperators: operator mode, handlers attached, and
// UpdateOperatorsStatusMap having read the operator as not deployed.
func operatorModeConsumer(t *testing.T, cluster *fakeCluster) *MesheryControllersHelper {
	t.Helper()
	mch := newTestControllersHelper(t)
	mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
	mch.ctxControllerHandlers = map[MesheryController]controllers.IMesheryController{
		MesheryOperator: &clusterOperator{cluster: cluster},
	}
	mch.ctxOperatorStatus = controllers.NotDeployed
	return mch
}

// TestConcurrentConnectsToOneClusterInstallTheOperatorOnce is the install
// race: every connection read "not deployed" before any had installed, and
// each ran helm install for the same release, which Helm rejects with "cannot
// re-use a name that is still in use". Serialized per cluster with the status
// re-read under the lock, exactly one install happens and none overlap.
func TestConcurrentConnectsToOneClusterInstallTheOperatorOnce(t *testing.T) {
	const consumers = 5
	ot := NewOperatorTracker(false)
	cluster := newFakeCluster()
	serverID := core.Uuid(uuid.Must(uuid.NewV4()))

	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < consumers; i++ {
		mch := operatorModeConsumer(t, cluster)
		k8sctx := connectionTo(serverID, "user-"+string(rune('a'+i)))
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			mch.DeployUndeployedOperators(ot, k8sctx)
		}()
	}
	close(start)
	wg.Wait()

	status, installs, _, maxInFlight := cluster.snapshot()
	if installs != 1 {
		t.Fatalf("helm install ran %d times for one cluster, want 1", installs)
	}
	if maxInFlight != 1 {
		t.Fatalf("%d Helm actions overlapped on one cluster, want them serialized", maxInFlight)
	}
	if status != controllers.Deployed {
		t.Fatalf("operator status = %v, want Deployed", status)
	}
	if got := ot.OperatorConsumers(connectionTo(serverID, "any")); len(got) != consumers {
		t.Fatalf("consumers = %v, want all %d connections recorded", got, consumers)
	}
}

// TestHelmLockIsPerClusterNotGlobal: serializing every cluster behind one lock
// would let a slow install on one cluster stall every other connection's.
func TestHelmLockIsPerClusterNotGlobal(t *testing.T) {
	ot := NewOperatorTracker(false)
	slow := newFakeCluster()
	slow.hold = make(chan struct{})
	defer close(slow.hold)
	fast := newFakeCluster()

	slowConsumer := operatorModeConsumer(t, slow)
	fastConsumer := operatorModeConsumer(t, fast)
	go slowConsumer.DeployUndeployedOperators(ot, connectionTo(core.Uuid(uuid.Must(uuid.NewV4())), "slow"))

	done := make(chan struct{})
	go func() {
		fastConsumer.DeployUndeployedOperators(ot, connectionTo(core.Uuid(uuid.Must(uuid.NewV4())), "fast"))
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("an install on one cluster blocked an install on another")
	}
	if _, installs, _, _ := fast.snapshot(); installs != 1 {
		t.Fatalf("the unblocked cluster's operator installed %d times, want 1", installs)
	}
}

// TestImplicitTransitionsLeaveASharedOperatorInstalled is the production
// incident: one connection switching to embedded mode ran helm uninstall on
// the operator every other operator-mode connection to that cluster depended
// on, and disconnecting or deleting a connection did the same. Each of those
// transitions now only releases the connection's claim.
func TestImplicitTransitionsLeaveASharedOperatorInstalled(t *testing.T) {
	transitions := []struct {
		name string
		run  func(*MesheryControllersHelper, *OperatorTracker, K8sContext)
	}{
		{
			// reconcileMeshsyncDeploymentMode and the connection update handler.
			name: "switches to embedded",
			run: func(mch *MesheryControllersHelper, ot *OperatorTracker, k8sctx K8sContext) {
				mch.DetachFromCluster(context.Background(), ot, k8sctx)
				mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeEmbedded).
					DeployUndeployedOperators(ot, k8sctx)
			},
		},
		{
			// DisconnectAction.Execute.
			name: "disconnects",
			run: func(mch *MesheryControllersHelper, ot *OperatorTracker, k8sctx K8sContext) {
				mch.DetachFromCluster(context.Background(), ot, k8sctx)
			},
		},
		{
			// DeleteAction.Execute.
			name: "is deleted",
			run: func(mch *MesheryControllersHelper, ot *OperatorTracker, k8sctx K8sContext) {
				mch.DetachFromCluster(context.Background(), ot, k8sctx)
			},
		},
	}

	for _, tt := range transitions {
		t.Run(tt.name, func(t *testing.T) {
			ot := NewOperatorTracker(false)
			cluster := newFakeCluster()
			serverID := core.Uuid(uuid.Must(uuid.NewV4()))
			first, firstCtx := operatorModeConsumer(t, cluster), connectionTo(serverID, "first")
			second, secondCtx := operatorModeConsumer(t, cluster), connectionTo(serverID, "second")
			first.DeployUndeployedOperators(ot, firstCtx)
			second.DeployUndeployedOperators(ot, secondCtx)

			tt.run(first, ot, firstCtx)

			status, installs, uninstalls, _ := cluster.snapshot()
			if uninstalls != 0 {
				t.Fatalf("a connection that %s uninstalled the shared operator (%d uninstalls)", tt.name, uninstalls)
			}
			if status != controllers.Deployed || installs != 1 {
				t.Fatalf("operator status = %v after %d installs, want it still Deployed from the one install", status, installs)
			}
			if got, want := ot.OperatorConsumers(secondCtx), []string{secondCtx.ConnectionID}; !reflect.DeepEqual(got, want) {
				t.Fatalf("consumers = %v, want only the connection still using the operator %v", got, want)
			}
		})
	}
}

// TestLastConsumerLeavingLeavesTheOperatorInstalled pins the design choice.
// This server can count only its own connections; other Meshery Servers and
// out-of-band (Helm, mesheryctl) installations are invisible to it, so its last
// consumer leaving is not evidence that nobody needs the operator. Removal is
// left to an explicit request.
func TestLastConsumerLeavingLeavesTheOperatorInstalled(t *testing.T) {
	ot := NewOperatorTracker(false)
	cluster := newFakeCluster()
	serverID := core.Uuid(uuid.Must(uuid.NewV4()))
	first, firstCtx := operatorModeConsumer(t, cluster), connectionTo(serverID, "first")
	second, secondCtx := operatorModeConsumer(t, cluster), connectionTo(serverID, "second")
	first.DeployUndeployedOperators(ot, firstCtx)
	second.DeployUndeployedOperators(ot, secondCtx)

	first.DetachFromCluster(context.Background(), ot, firstCtx)
	second.DetachFromCluster(context.Background(), ot, secondCtx)

	if status, _, uninstalls, _ := cluster.snapshot(); uninstalls != 0 || status != controllers.Deployed {
		t.Fatalf("operator status = %v after %d uninstalls, want it left installed", status, uninstalls)
	}
	if got := ot.OperatorConsumers(firstCtx); len(got) != 0 {
		t.Fatalf("consumers = %v, want none once every connection has left", got)
	}
}

// subscribeSystemEvents returns the events mch publishes for the system owner,
// which is who owns a controller event raised with no user in scope.
func subscribeSystemEvents(t *testing.T, mch *MesheryControllersHelper) <-chan interface{} {
	t.Helper()
	systemID := core.Uuid(uuid.Must(uuid.NewV4()))
	mch.eventBroadcaster = NewBroadcaster("test")
	mch.systemID = &systemID
	ch, unsubscribe := mch.eventBroadcaster.Subscribe(systemID)
	t.Cleanup(unsubscribe)
	return ch
}

// TestExplicitUndeployOfASharedOperatorNamesTheAffectedConnections: the one
// remaining way to uninstall is an explicit request, and since it removes the
// operator for everyone on the cluster it says who else it removed it from.
func TestExplicitUndeployOfASharedOperatorNamesTheAffectedConnections(t *testing.T) {
	ot := NewOperatorTracker(false)
	cluster := newFakeCluster()
	serverID := core.Uuid(uuid.Must(uuid.NewV4()))
	first, firstCtx := operatorModeConsumer(t, cluster), connectionTo(serverID, "first")
	second, secondCtx := operatorModeConsumer(t, cluster), connectionTo(serverID, "second")
	first.DeployUndeployedOperators(ot, firstCtx)
	second.DeployUndeployedOperators(ot, secondCtx)
	published := subscribeSystemEvents(t, first)

	if err := first.SetOperatorDeployment(ot, firstCtx, false); err != nil {
		t.Fatalf("explicit undeploy failed: %v", err)
	}

	if _, _, uninstalls, _ := cluster.snapshot(); uninstalls != 1 {
		t.Fatalf("explicit undeploy ran %d uninstalls, want 1", uninstalls)
	}
	select {
	case msg := <-published:
		event, ok := msg.(*events.Event)
		if !ok {
			t.Fatalf("published %T, want *events.Event", msg)
		}
		if event.Severity != events.Warning {
			t.Fatalf("severity = %v, want a warning", event.Severity)
		}
		affected, _ := event.Metadata["affectedConnectionIDs"].([]string)
		if want := []string{secondCtx.ConnectionID}; !reflect.DeepEqual(affected, want) {
			t.Fatalf("affectedConnectionIDs = %v, want the other connection using the operator %v", affected, want)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("removing an operator other connections use published no event naming them")
	}
}

// TestExplicitUndeployOfAnUnsharedOperatorIsQuiet: with no other connection on
// this server using the operator there is no one to name, and a warning on
// every intended removal would be noise.
func TestExplicitUndeployOfAnUnsharedOperatorIsQuiet(t *testing.T) {
	ot := NewOperatorTracker(false)
	cluster := newFakeCluster()
	only, onlyCtx := operatorModeConsumer(t, cluster), connectionTo(core.Uuid(uuid.Must(uuid.NewV4())), "only")
	only.DeployUndeployedOperators(ot, onlyCtx)
	published := subscribeSystemEvents(t, only)

	if err := only.SetOperatorDeployment(ot, onlyCtx, false); err != nil {
		t.Fatalf("explicit undeploy failed: %v", err)
	}
	if _, _, uninstalls, _ := cluster.snapshot(); uninstalls != 1 {
		t.Fatalf("explicit undeploy ran %d uninstalls, want 1", uninstalls)
	}
	select {
	case msg := <-published:
		t.Fatalf("published %v for a removal that affected no other connection", msg)
	default:
	}
}

func TestOperatorClusterKey(t *testing.T) {
	serverID := core.Uuid(uuid.Must(uuid.NewV4()))
	viaLoadBalancer := connectionTo(serverID, "lb")
	viaNode := connectionTo(serverID, "node")
	if operatorClusterKey(viaLoadBalancer) != operatorClusterKey(viaNode) {
		t.Fatal("two API server addresses for one cluster must share its operator lifecycle")
	}

	unresolved := K8sContext{ID: "ctx-a", Server: "https://shared.example.invalid:6443"}
	sameServer := K8sContext{ID: "ctx-b", Server: "https://shared.example.invalid:6443"}
	if operatorClusterKey(unresolved) != operatorClusterKey(sameServer) {
		t.Fatal("contexts without a server ID must still share a key when they name the same API server")
	}
	nilID := core.Uuid(uuid.Nil)
	unresolved.KubernetesServerID = &nilID
	if operatorClusterKey(unresolved) != operatorClusterKey(sameServer) {
		t.Fatal("a nil server ID must fall back to the API server, not key every such context together")
	}

	if operatorClusterKey(K8sContext{ID: "ctx-a"}) == operatorClusterKey(K8sContext{ID: "ctx-b"}) {
		t.Fatal("contexts with nothing in common must not share a key")
	}
}
