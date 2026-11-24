package models

import (
	"sort"
	"sync"

	"github.com/gofrs/uuid"
)

// Meshery Operator is a cluster singleton: one Helm release, `meshery-operator`
// in namespace `meshery`, serves every Meshery connection - on this server or
// any other - whose kubeconfig reaches that cluster. Its lifecycle used to be
// driven per connection, so any one connection disconnecting, being deleted, or
// switching to embedded MeshSync ran `helm uninstall` underneath every other
// operator-mode connection to the same cluster, and two connections connecting
// at once raced each other into `helm install`.
//
// The lifecycle is therefore cluster-scoped, on two rules:
//
//  1. Implicit transitions never uninstall. Disconnect, delete and a switch to
//     embedded mode only release the connection's claim (ReleaseOperator).
//     Removal happens solely on an explicit user request
//     (SetOperatorDeployment with deploy=false), which names the other
//     connections on this server that it takes the operator away from.
//  2. Every Helm action against a cluster is serialized per cluster, and the
//     operator's live status is re-read under that lock, so the second of two
//     concurrent connects sees the first one's install instead of repeating it.
//
// A cluster is known by aliases rather than one key (operatorClusterKeys): its
// Kubernetes server ID and the API server URL a context names. Two connections
// sharing either alias serialize their Helm actions and see each other as
// consumers. A context with no server ID whose URL no other context names
// cannot be unified, and stays keyed on its own.
//
// Uninstalling when the last consumer leaves is deliberately not attempted.
// This process can only count its own connections: other Meshery Servers, and
// operators installed by Helm or mesheryctl, are invisible to it, so "the last
// consumer I know of" is not "the last consumer". An in-cluster consumer
// registry would need leases and heartbeats from every server to every cluster
// to stay correct, and still could not see out-of-band installs. Getting
// removal wrong is an outage for someone else; leaving an idle operator behind
// is visible and removable on request.

// operatorCluster is the lifecycle state recorded under one cluster key.
type operatorCluster struct {
	// helm serializes Deploy and Undeploy against this cluster. It is held only
	// around the Helm calls themselves and is always the innermost lock taken:
	// nothing acquires OperatorTracker.mx, logs or emits an event while holding it.
	helm sync.Mutex
	// consumers is the set of operator-mode connections on this server that
	// currently use this cluster's operator. Guarded by OperatorTracker.mx.
	consumers map[string]struct{}
}

// operatorClusterKeys returns, sorted, every key under which the cluster a
// context reaches is known. Each is an alias for the same cluster, and state is
// recorded and locked under all of them, so two contexts that share any one key
// share the operator lifecycle.
//
// The Kubernetes server ID (the kube-system namespace UID) is the cluster's
// identity everywhere else in Meshery - MeshSync rows, the delete-time flush
// refcount - and is the only key under which two different API server URLs (a
// load balancer and a node address, say) are recognised as one cluster. The
// API server URL is the second alias: a context's persisted row carries a
// server ID only if its cluster was reachable when it was imported, and that
// row is never rewritten, so the URL is what unifies such a context with the
// others naming the same endpoint. A context with neither falls back to its
// context ID, which is no coarser than the per-connection behavior this
// replaces.
//
// Deriving identity this way needs no cluster I/O, so nothing has to resolve
// or rewrite a context before it is keyed. The edge it leaves: aliasing is not
// transitive, so a context with no server ID whose API server URL no other
// context names cannot be unified with the cluster's other connections - and
// neither can two contexts that are only linked through a third context
// carrying both their keys.
func operatorClusterKeys(k8sctx K8sContext) []string {
	keys := make([]string, 0, 2)
	if id := k8sctx.KubernetesServerID; id != nil && *id != uuid.Nil {
		keys = append(keys, "kubernetes-server-id:"+id.String())
	}
	if k8sctx.Server != "" {
		keys = append(keys, "server:"+k8sctx.Server)
	}
	if len(keys) == 0 {
		keys = append(keys, "context:"+k8sctx.ID)
	}
	sort.Strings(keys)
	return keys
}

// operatorConsumerID identifies a connection as a consumer of its cluster's
// operator. Every kubeconfig context is its own connection, so two contexts
// reaching one cluster are two consumers.
func operatorConsumerID(k8sctx K8sContext) string {
	if k8sctx.ConnectionID != "" {
		return k8sctx.ConnectionID
	}
	return k8sctx.ID
}

// clusters returns the lifecycle state for each of k8sctx's cluster keys, in
// key order, creating any on first use.
func (ot *OperatorTracker) clusters(k8sctx K8sContext) []*operatorCluster {
	keys := operatorClusterKeys(k8sctx)
	ot.mx.Lock()
	defer ot.mx.Unlock()
	if ot.clustersByKey == nil {
		ot.clustersByKey = make(map[string]*operatorCluster)
	}
	cs := make([]*operatorCluster, 0, len(keys))
	for _, key := range keys {
		c, ok := ot.clustersByKey[key]
		if !ok {
			c = &operatorCluster{consumers: make(map[string]struct{})}
			ot.clustersByKey[key] = c
		}
		cs = append(cs, c)
	}
	return cs
}

// lockOperatorCluster serializes a Helm action against the cluster k8sctx
// reaches and returns the function that releases it. It takes the Helm lock of
// every one of the context's cluster keys, in sorted key order, and releases
// them in reverse; the order is the same for every caller, so two contexts
// sharing some of their keys cannot deadlock. A nil tracker serializes nothing;
// every production path carries the server-wide tracker.
func (ot *OperatorTracker) lockOperatorCluster(k8sctx K8sContext) (unlock func()) {
	if ot == nil {
		return func() {}
	}
	cs := ot.clusters(k8sctx)
	for _, c := range cs {
		c.helm.Lock()
	}
	return func() {
		for i := len(cs) - 1; i >= 0; i-- {
			cs[i].helm.Unlock()
		}
	}
}

// withOperatorClusterLock runs helmAction, and only helmAction, under k8sctx's
// cluster Helm lock, releasing it even if helmAction panics. Callers keep
// consumer lookups, error recording, logging and event emission outside it.
func (ot *OperatorTracker) withOperatorClusterLock(k8sctx K8sContext, helmAction func()) {
	unlock := ot.lockOperatorCluster(k8sctx)
	defer unlock()
	helmAction()
}

// claimOperator records k8sctx's connection as a consumer of its cluster's
// operator, under every one of its cluster keys. Idempotent.
func (ot *OperatorTracker) claimOperator(k8sctx K8sContext) {
	if ot == nil {
		return
	}
	cs := ot.clusters(k8sctx)
	id := operatorConsumerID(k8sctx)
	ot.mx.Lock()
	defer ot.mx.Unlock()
	for _, c := range cs {
		c.consumers[id] = struct{}{}
	}
}

// releaseOperatorClaim removes k8sctx's connection from its cluster's
// consumers, reporting whether it held a claim and how many consumers remain.
// Idempotent; it never touches the cluster.
func (ot *OperatorTracker) releaseOperatorClaim(k8sctx K8sContext) (released bool, remaining int) {
	if ot == nil {
		return false, 0
	}
	cs := ot.clusters(k8sctx)
	id := operatorConsumerID(k8sctx)
	ot.mx.Lock()
	defer ot.mx.Unlock()
	for _, c := range cs {
		if _, ok := c.consumers[id]; ok {
			released = true
			delete(c.consumers, id)
		}
	}
	return released, len(consumerUnion(cs))
}

// OperatorConsumers returns, sorted, the connections on this server that use
// the operator of the cluster k8sctx reaches: the consumers recorded under any
// of its cluster keys. It cannot know about other Meshery Servers or about
// operators installed out of band, so it is a lower bound on who depends on
// that operator, never the full set.
func (ot *OperatorTracker) OperatorConsumers(k8sctx K8sContext) []string {
	if ot == nil {
		return nil
	}
	cs := ot.clusters(k8sctx)
	ot.mx.Lock()
	defer ot.mx.Unlock()
	return consumerUnion(cs)
}

// consumerUnion is the sorted, de-duplicated union of the consumers of cs.
// The caller holds OperatorTracker.mx.
func consumerUnion(cs []*operatorCluster) []string {
	seen := make(map[string]struct{})
	for _, c := range cs {
		for id := range c.consumers {
			seen[id] = struct{}{}
		}
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// otherOperatorConsumers is OperatorConsumers without k8sctx's own connection:
// the connections an action taken on k8sctx's behalf would also affect.
func (ot *OperatorTracker) otherOperatorConsumers(k8sctx K8sContext) []string {
	self := operatorConsumerID(k8sctx)
	all := ot.OperatorConsumers(k8sctx)
	others := make([]string, 0, len(all))
	for _, id := range all {
		if id != self {
			others = append(others, id)
		}
	}
	return others
}
