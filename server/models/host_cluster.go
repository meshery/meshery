package models

import (
	"context"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Meshery Server running inside *a* Kubernetes cluster says nothing about the
// cluster a given connection targets. A Server in one cluster routinely manages
// others through their kubeconfigs, and a broker exposed only on a ClusterIP is
// reachable directly only from inside its own cluster. Deciding "skip the
// managed broker port-forward" from KUBERNETES_SERVICE_HOST alone therefore left
// every remote cluster's broker unreachable whenever the Server itself ran in a
// Pod. hostCluster answers the per-target question instead: does this context
// reach the cluster Meshery Server runs in?

const (
	// hostClusterIDTimeout bounds the one API call that resolves the host
	// cluster's identity.
	hostClusterIDTimeout = 5 * time.Second
	// hostClusterIDRetry is how long a failed resolution is remembered before
	// it is attempted again, so a missing RBAC grant costs one call per window
	// rather than one per broker connection attempt.
	hostClusterIDRetry = time.Minute
)

// inClusterAPIHosts are the names the Kubernetes API server answers to from
// inside any cluster. A context naming one of them can only mean the cluster
// the caller runs in.
var inClusterAPIHosts = map[string]bool{
	"kubernetes":                           true,
	"kubernetes.default":                   true,
	"kubernetes.default.svc":               true,
	"kubernetes.default.svc.cluster.local": true,
}

// hostCluster identifies the Kubernetes cluster Meshery Server runs in.
type hostCluster struct {
	// apiHostPort is the in-cluster API server address
	// (KUBERNETES_SERVICE_HOST:KUBERNETES_SERVICE_PORT). Empty when Meshery
	// Server does not run in a cluster.
	apiHostPort string
	// resolveServerID returns the host cluster's Kubernetes server ID (the
	// kube-system namespace UID). Injected so tests need no cluster.
	resolveServerID func(context.Context) (string, error)
	now             func() time.Time

	mu          sync.Mutex
	serverID    string
	lastAttempt time.Time
}

// defaultHostCluster is the identity of the cluster this process runs in, read
// from the environment Kubernetes injects into every Pod.
var defaultHostCluster = newHostCluster(
	os.Getenv("KUBERNETES_SERVICE_HOST"),
	os.Getenv("KUBERNETES_SERVICE_PORT"),
	resolveInClusterServerID,
)

func newHostCluster(host, port string, resolve func(context.Context) (string, error)) *hostCluster {
	h := &hostCluster{resolveServerID: resolve, now: time.Now}
	if host != "" && port != "" {
		h.apiHostPort = net.JoinHostPort(host, port)
	}
	return h
}

// resolveInClusterServerID reads the kube-system namespace UID through the
// Pod's service account - the same identity NewK8sContextWithServerID assigns a
// connection, so the two are directly comparable.
func resolveInClusterServerID(ctx context.Context) (string, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return "", err
	}
	cfg.Timeout = hostClusterIDTimeout
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return "", err
	}
	ns, err := client.CoreV1().Namespaces().Get(ctx, "kube-system", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(ns.GetUID()), nil
}

// inCluster reports whether Meshery Server runs inside a Kubernetes cluster.
func (h *hostCluster) inCluster() bool {
	return h != nil && h.apiHostPort != ""
}

// id returns the host cluster's server ID, or "" when it is not known (yet).
// A success is cached for the life of the process; a failure is retried after
// hostClusterIDRetry.
func (h *hostCluster) id() string {
	if !h.inCluster() || h.resolveServerID == nil {
		return ""
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.serverID != "" {
		return h.serverID
	}
	now := h.now()
	if !h.lastAttempt.IsZero() && now.Sub(h.lastAttempt) < hostClusterIDRetry {
		return ""
	}
	h.lastAttempt = now
	ctx, cancel := context.WithTimeout(context.Background(), hostClusterIDTimeout)
	defer cancel()
	id, err := h.resolveServerID(ctx)
	if err != nil {
		return ""
	}
	h.serverID = id
	return id
}

// hosts reports whether k8sctx reaches the cluster Meshery Server runs in.
//
// Identity is decided by the Kubernetes server ID, which is what Meshery uses
// for cluster identity everywhere else and the only signal that survives a
// kubeconfig naming the host cluster by an external address - the usual case,
// since an uploaded kubeconfig almost never names the in-cluster ClusterIP.
// When either ID is unknown (a context registered while unreachable, or a
// service account without read access to namespaces), it falls back to the API
// server address, which positively identifies only contexts that name the
// in-cluster endpoint itself.
//
// Anything not positively identified is treated as a different cluster. That is
// the safe direction for its one consumer: a needless broker port-forward costs
// a hop through the API server, while a skipped one leaves a ClusterIP-only
// broker unreachable.
func (h *hostCluster) hosts(k8sctx K8sContext) bool {
	if !h.inCluster() {
		return false
	}
	if target := k8sctx.KubernetesServerID; target != nil && *target != uuid.Nil {
		if id := h.id(); id != "" {
			return strings.EqualFold(id, target.String())
		}
	}
	return h.namesInClusterAPIServer(k8sctx.Server)
}

// namesInClusterAPIServer reports whether server (a kubeconfig cluster server
// URL) is the in-cluster API endpoint.
func (h *hostCluster) namesInClusterAPIServer(server string) bool {
	u, err := url.Parse(server)
	if err != nil || u.Hostname() == "" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	if inClusterAPIHosts[host] {
		return true
	}
	port := u.Port()
	if port == "" {
		port = "443"
	}
	return net.JoinHostPort(host, port) == h.apiHostPort
}
