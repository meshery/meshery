// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/utils/kubernetes/describe"
)

// describeResourceMap maps user-friendly kind strings to the describe library's
// DescribeType enum. Mirrors the same map in the GraphQL resolver
// (kubernetes.go:41-70) — keep them in sync until the resolver is deleted.
var describeResourceMap = map[string]describe.DescribeType{
	"pod":                       describe.Pod,
	"deployment":                describe.Deployment,
	"job":                       describe.Job,
	"cronjob":                   describe.CronJob,
	"statefulset":               describe.StatefulSet,
	"daemonset":                 describe.DaemonSet,
	"replicaset":                describe.ReplicaSet,
	"secret":                    describe.Secret,
	"service":                   describe.Service,
	"serviceaccount":            describe.ServiceAccount,
	"node":                      describe.Node,
	"limitrange":                describe.LimitRange,
	"resourcequota":             describe.ResourceQuota,
	"persistentvolume":          describe.PersistentVolume,
	"persistentvolumeclaim":     describe.PersistentVolumeClaim,
	"namespace":                 describe.Namespace,
	"endpoints":                 describe.Endpoints,
	"configmap":                 describe.ConfigMap,
	"priorityclass":             describe.PriorityClass,
	"ingress":                   describe.Ingress,
	"role":                      describe.Role,
	"clusterrole":               describe.ClusterRole,
	"rolebinding":               describe.RoleBinding,
	"clusterrolebinding":        describe.ClusterRoleBinding,
	"networkpolicy":             describe.NetworkPolicy,
	"replicationcontroller":     describe.ReplicationController,
	"certificatesigningrequest": describe.CertificateSigningRequest,
	"endpointslice":             describe.EndpointSlice,
}

// KubectlDescribeHandler returns `kubectl describe <kind>/<name> -n <namespace>`
// for the supplied resource in the user's active Kubernetes context. It
// replaces the `getKubectlDescribe` GraphQL query and fixes the resolver's
// bug at kubernetes.go:78 where `meshkitKube.New([]byte(""))` was constructed
// with an empty kubeconfig and effectively ignored the user's context. The
// REST endpoint requires a `contextId` query parameter and loads the
// kubeconfig for that context (mirrors KubernetesPingHandler's lookup flow).
func (h *Handler) KubectlDescribeHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	q := req.URL.Query()
	name := q.Get("name")
	kind := q.Get("kind")
	namespace := q.Get("namespace")
	contextID := q.Get("contextId")

	if name == "" || kind == "" {
		writeJSONError(w, "name and kind query parameters are required", http.StatusBadRequest)
		return
	}
	if contextID == "" {
		writeJSONError(w, "contextId query parameter is required", http.StatusBadRequest)
		return
	}

	dType, ok := describeResourceMap[strings.ToLower(kind)]
	if !ok {
		writeJSONError(w, "unsupported resource kind", http.StatusBadRequest)
		return
	}

	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		writeJSONError(w, "no auth token", http.StatusUnauthorized)
		return
	}

	k8sContext, err := provider.GetK8sContext(token, contextID)
	if err != nil {
		writeMeshkitError(w, ErrInvalidKubeContext(err, contextID), http.StatusNotFound)
		return
	}
	kubeclient, err := k8sContext.GenerateKubeHandler()
	if err != nil {
		writeMeshkitError(w, ErrInvalidKubeConfig(err, ""), http.StatusBadRequest)
		return
	}

	details, err := describe.Describe(kubeclient, describe.DescriberOptions{
		Name:      name,
		Namespace: namespace,
		Type:      dType,
	})
	if err != nil {
		// Log the full error server-side for debugging, but return a
		// generic message to the client — kubectl-describe errors can
		// embed cluster paths, namespaces, RBAC details, and other
		// implementation specifics that don't belong in a 500 response.
		h.log.Error(err)
		writeJSONError(w, "failed to describe resource", http.StatusInternalServerError)
		return
	}

	resp := &model.KctlDescribeDetails{
		Describe: &details,
		Ctxid:    &contextID,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(models.ErrMarshal(err, "kubectl describe response"))
		return
	}
}
