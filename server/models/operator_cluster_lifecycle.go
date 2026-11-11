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
// Uninstalling when the last consumer leaves is deliberately not attempted.
// This process can only count its own connections: other Meshery Servers, and
// operators installed by Helm or mesheryctl, are invisible to it, so "the last
// consumer I know of" is not "the last consumer". An in-cluster consumer
// registry would need leases and heartbeats from every server to every cluster
// to stay correct, and still could not see out-of-band installs. Getting
// removal wrong is an outage for someone else; leaving an idle operator behind
// is visible and removable on request.

// operatorCluster is the per-cluster lifecycle state.
type operatorCluster struct {
	// helm serializes Deploy and Undeploy against this cluster. It is held for
	// the duration of a Helm action and is always the innermost lock taken.
	helm sync.Mutex
	// consumers is the set of operator-mode connections on this server that
	// currently use this cluster's operator. Guarded by OperatorTracker.mx.
	consumers map[string]struct{}
}

// operatorClusterKey identifies the cluster a context reaches, for keying
// cluster-scoped operator lifecycle state.
//
// The Kubernetes server ID (the kube-system namespace UID) is the cluster's
// identity everywhere else in Meshery - MeshSync rows, the delete-time flush
// refcount - and is the only key under which two different API server URLs (a
// load balancer and a node address, say) are recognised as one cluster. A
// context whose ID was never resolved falls back to its API server URL, which
// still unifies the common case of several kubeconfig contexts naming the same
// endpoint, and finally to the context ID, which is no coarser than the
// per-connection behavior this replaces.
func operatorClusterKey(k8sctx K8sContext) string {
	if id := k8sctx.KubernetesServerID; id != nil && *id != uuid.Nil {
		return "kubernetes-server-id:" + id.String()
	}
	if k8sctx.Server != "" {
		return "server:" + k8sctx.Server
	}
	return "context:" + k8sctx.ID
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

// cluster returns the lifecycle state for key, creating it on first use.
func (ot *OperatorTracker) cluster(key string) *operatorCluster {
	ot.mx.Lock()
	defer ot.mx.Unlock()
	if ot.clusters == nil {
		ot.clusters = make(map[string]*operatorCluster)
	}
	c, ok := ot.clusters[key]
	if !ok {
		c = &operatorCluster{consumers: make(map[string]struct{})}
		ot.clusters[key] = c
	}
	return c
}

// lockOperatorCluster serializes a Helm action against the cluster k8sctx
// reaches and returns the function that releases it. A nil tracker serializes
// nothing; every production path carries the server-wide tracker.
func (ot *OperatorTracker) lockOperatorCluster(k8sctx K8sContext) (unlock func()) {
	if ot == nil {
		return func() {}
	}
	c := ot.cluster(operatorClusterKey(k8sctx))
	c.helm.Lock()
	return c.helm.Unlock
}

// claimOperator records k8sctx's connection as a consumer of its cluster's
// operator. Idempotent.
func (ot *OperatorTracker) claimOperator(k8sctx K8sContext) {
	if ot == nil {
		return
	}
	c := ot.cluster(operatorClusterKey(k8sctx))
	ot.mx.Lock()
	c.consumers[operatorConsumerID(k8sctx)] = struct{}{}
	ot.mx.Unlock()
}

// releaseOperatorClaim removes k8sctx's connection from its cluster's
// consumers, reporting whether it held a claim and how many consumers remain.
// Idempotent; it never touches the cluster.
func (ot *OperatorTracker) releaseOperatorClaim(k8sctx K8sContext) (released bool, remaining int) {
	if ot == nil {
		return false, 0
	}
	c := ot.cluster(operatorClusterKey(k8sctx))
	id := operatorConsumerID(k8sctx)
	ot.mx.Lock()
	defer ot.mx.Unlock()
	_, released = c.consumers[id]
	delete(c.consumers, id)
	return released, len(c.consumers)
}

// OperatorConsumers returns, sorted, the connections on this server that use
// the operator of the cluster k8sctx reaches. It cannot know about other
// Meshery Servers or about operators installed out of band, so it is a lower
// bound on who depends on that operator, never the full set.
func (ot *OperatorTracker) OperatorConsumers(k8sctx K8sContext) []string {
	if ot == nil {
		return nil
	}
	c := ot.cluster(operatorClusterKey(k8sctx))
	ot.mx.Lock()
	defer ot.mx.Unlock()
	ids := make([]string, 0, len(c.consumers))
	for id := range c.consumers {
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
