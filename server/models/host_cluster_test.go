package models

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/core"
)

const (
	hostClusterID   = "6f1f6a4e-6f8b-4d55-9a0e-2f0b8f6b1c01"
	remoteClusterID = "0b9d0c55-3a0e-4f7e-8f1e-7d2c7c0e4a02"
)

// inClusterAt is a Meshery Server running in the cluster whose kube-system
// UID is id, with its API server at the conventional ClusterIP.
func inClusterAt(id string, err error) *hostCluster {
	return newHostCluster("10.96.0.1", "443", func(context.Context) (string, error) {
		return id, err
	})
}

func contextFor(server, serverID string) K8sContext {
	k := K8sContext{ID: "ctx", Server: server}
	if serverID != "" {
		id := core.Uuid(uuid.FromStringOrNil(serverID))
		k.KubernetesServerID = &id
	}
	return k
}

func TestHostClusterHosts(t *testing.T) {
	tests := []struct {
		name   string
		host   *hostCluster
		target K8sContext
		want   bool
	}{
		{
			name:   "Meshery not running in any cluster",
			host:   newHostCluster("", "", nil),
			target: contextFor("https://10.96.0.1:443", hostClusterID),
			want:   false,
		},
		{
			// The production incident: Meshery in a Pod of one cluster, managing
			// another through its kubeconfig.
			name:   "a different cluster, identified by server ID",
			host:   inClusterAt(hostClusterID, nil),
			target: contextFor("https://remote.example.invalid:6443", remoteClusterID),
			want:   false,
		},
		{
			// The usual shape of a kubeconfig for Meshery's own cluster: an
			// external address, not the ClusterIP.
			name:   "the host cluster named by an external address",
			host:   inClusterAt(hostClusterID, nil),
			target: contextFor("https://api.host.example.invalid:6443", hostClusterID),
			want:   true,
		},
		{
			name:   "the in-cluster context, with no server ID yet",
			host:   inClusterAt(hostClusterID, nil),
			target: contextFor("https://10.96.0.1:443", ""),
			want:   true,
		},
		{
			name:   "the in-cluster service name, default port",
			host:   inClusterAt(hostClusterID, nil),
			target: contextFor("https://kubernetes.default.svc", ""),
			want:   true,
		},
		{
			name:   "a different cluster with no server ID",
			host:   inClusterAt(hostClusterID, nil),
			target: contextFor("https://remote.example.invalid:6443", ""),
			want:   false,
		},
		{
			// A service account that cannot read namespaces leaves only the
			// address to go on, and an unrecognized one is not the host.
			name:   "host server ID unresolvable, remote address",
			host:   inClusterAt("", errors.New("forbidden")),
			target: contextFor("https://api.host.example.invalid:6443", hostClusterID),
			want:   false,
		},
		{
			name:   "host server ID unresolvable, in-cluster address",
			host:   inClusterAt("", errors.New("forbidden")),
			target: contextFor("https://10.96.0.1", hostClusterID),
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.host.hosts(tt.target); got != tt.want {
				t.Fatalf("hosts(%s) = %v, want %v", tt.target.Server, got, tt.want)
			}
		})
	}
}

// TestHostClusterIDIsCachedAndFailuresRetried: the identity is one API call
// for the life of the process once it succeeds, and a failure costs one call
// per retry window rather than one per broker connection attempt.
func TestHostClusterIDIsCachedAndFailuresRetried(t *testing.T) {
	calls := 0
	fail := true
	h := newHostCluster("10.96.0.1", "443", func(context.Context) (string, error) {
		calls++
		if fail {
			return "", errors.New("forbidden")
		}
		return hostClusterID, nil
	})
	now := time.Unix(0, 0)
	h.now = func() time.Time { return now }

	if id := h.id(); id != "" || calls != 1 {
		t.Fatalf("id = %q after %d calls, want a failed first resolution", id, calls)
	}
	now = now.Add(hostClusterIDRetry / 2)
	if h.id(); calls != 1 {
		t.Fatalf("resolution retried after %v, inside the retry window (%d calls)", hostClusterIDRetry/2, calls)
	}
	fail = false
	now = now.Add(hostClusterIDRetry)
	if id := h.id(); id != hostClusterID || calls != 2 {
		t.Fatalf("id = %q after %d calls, want a retry once the window passed", id, calls)
	}
	now = now.Add(time.Hour)
	if h.id(); calls != 2 {
		t.Fatalf("a resolved identity was resolved again (%d calls)", calls)
	}
}

// portForwardingBroker is a broker controller that records whether Meshery
// asked it for a managed port-forward. It hands back an error rather than a
// forwarder, so no test reaches a cluster.
type portForwardingBroker struct {
	stubController
	requests int
}

func (b *portForwardingBroker) GetPortForwarder(logger.Handler) (*mesherykube.PortForwarder, error) {
	b.requests++
	return nil, errors.New("no cluster in unit tests")
}

// TestBrokerPortForwardIsDecidedPerTargetCluster is defect B: the port-forward
// was skipped whenever Meshery ran in any Pod, so a remote cluster's
// ClusterIP-only broker was unreachable from a Meshery Server running in a
// different cluster.
func TestBrokerPortForwardIsDecidedPerTargetCluster(t *testing.T) {
	t.Setenv("MESHERY_MANAGED_BROKER_PORTFORWARD", "")
	tests := []struct {
		name   string
		host   *hostCluster
		target K8sContext
		want   int
	}{
		{"Meshery in one cluster, connection to another", inClusterAt(hostClusterID, nil), contextFor("https://remote.example.invalid:6443", remoteClusterID), 1},
		{"Meshery in the connection's own cluster", inClusterAt(hostClusterID, nil), contextFor("https://api.host.example.invalid:6443", hostClusterID), 0},
		{"Meshery outside any cluster", newHostCluster("", "", nil), contextFor("https://remote.example.invalid:6443", remoteClusterID), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mch := newTestControllersHelper(t)
			mch.hostCluster = tt.host
			broker := &portForwardingBroker{stubController: stubController{status: controllers.Deployed}}

			mch.ensureBrokerPortForward(tt.target, broker)

			if broker.requests != tt.want {
				t.Fatalf("port-forward requested %d times, want %d", broker.requests, tt.want)
			}
		})
	}

	t.Run("disabled by the operator of the server", func(t *testing.T) {
		t.Setenv("MESHERY_MANAGED_BROKER_PORTFORWARD", "false")
		mch := newTestControllersHelper(t)
		mch.hostCluster = inClusterAt(hostClusterID, nil)
		broker := &portForwardingBroker{}
		mch.ensureBrokerPortForward(contextFor("https://remote.example.invalid:6443", remoteClusterID), broker)
		if broker.requests != 0 {
			t.Fatalf("port-forward requested %d times with MESHERY_MANAGED_BROKER_PORTFORWARD=false", broker.requests)
		}
	})
}
