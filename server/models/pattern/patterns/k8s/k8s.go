package k8s

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/converter"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta2/component"

	meshkube "github.com/meshery/meshkit/utils/kubernetes"
	"gopkg.in/yaml.v2"
	kubeerror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func Deploy(kubeClient *meshkube.Client, comp component.ComponentDefinition, isDel bool) error {

	resource := converter.CreateK8sResourceStructure(&comp)
	manifest, err := yaml.Marshal(resource)
	if err != nil {
		return err
	}

	// Define a function to extract namesapce, labels and annotations in the componetn definiotn
	namespace := getNamespaceForComponent(&comp)

	err = kubeClient.ApplyManifest(manifest, meshkube.ApplyOptions{
		Namespace: namespace,
		Update:    true,
		Delete:    isDel,
	})

	if err != nil {
		return meshkube.ErrApplyManifest(err)
	}
	return nil
}

// DryRunOptions encapsulates configuration parameters for dry-run validation,
// including a discovery RESTMapper for dynamic resource resolution and in-design declared namespaces.
type DryRunOptions struct {
	RESTMapper         meta.RESTMapper
	DeclaredNamespaces map[string]bool
}

// DryRunHelper performs server-side dry-run validation for a component definition against the client's cluster,
// using opts for discovery-based resource resolution and in-design declared namespaces.
func DryRunHelper(client *meshkube.Client, comp component.ComponentDefinition, isDelete bool, opts DryRunOptions) (st map[string]interface{}, success bool, err error) {
	if client == nil || client.KubeClient == nil {
		return nil, false, ErrDryRun(fmt.Errorf("kubernetes client is nil"), "")
	}
	resource := converter.CreateK8sResourceStructure(&comp)
	namespace := getNamespaceForComponent(&comp)
	return dryRun(client.KubeClient.RESTClient(), resource, namespace, isDelete, opts)
}

// dryRun performs a dry-run operation on the given Kubernetes resource using the provided REST client.
// It constructs the appropriate API request based on the resource type and operation (create/update or delete).
// Returns the status response from the Kubernetes server, a success flag, and any encountered error.
func dryRun(rClient rest.Interface, k8sResource map[string]interface{}, namespace string, isDelete bool, opts DryRunOptions) (st map[string]interface{}, success bool, err error) {
	if rClient == nil {
		return nil, false, ErrDryRun(fmt.Errorf("RESTClient is nil"), "")
	}

	st, success, reqErr := executeDryRunRequest(rClient, k8sResource, namespace, isDelete, opts.RESTMapper)

	// Suppress 404 NotFound on undeploy dry-runs only when the object being deleted is what is missing
	if isDelete {
		if isDeleteTargetNotFound(reqErr, st, k8sResource, namespace, opts.RESTMapper) {
			return map[string]interface{}{"status": "Success"}, true, nil
		}
		if !success && st != nil && st["status"] == "Failure" {
			return st, false, nil
		}
		return st, success, reqErr
	}

	// For create dry-runs, account for in-design namespaces
	if !success && opts.DeclaredNamespaces != nil && opts.DeclaredNamespaces[namespace] && isMissingNamespaceError(reqErr, st, namespace) {
		// Substitute "default" namespace on a copy of the unstructured object and re-run create dry-run
		resourceCopy := deepCopyMap(k8sResource)
		if resourceCopy == nil {
			resourceCopy = make(map[string]interface{})
		}
		metaMap, ok := resourceCopy["metadata"].(map[string]interface{})
		if !ok || metaMap == nil {
			metaMap = make(map[string]interface{})
			resourceCopy["metadata"] = metaMap
		}
		metaMap["namespace"] = "default"

		stDef, okDef, reqErrDef := executeDryRunRequest(rClient, resourceCopy, "default", false, opts.RESTMapper)
		if okDef {
			return stDef, true, nil
		}
		// The API server checks for a name conflict only after admission and schema validation have passed, so
		// an unrelated object of the same name in default still means the workload is valid. Return the
		// workload itself, so the component's configuration in the response is its spec, not a Status.
		if isFallbackNameConflict(reqErrDef, stDef, resourceCopy) {
			return k8sResource, true, nil
		}
		// With "default" missing too, the workload cannot be validated on this cluster. Its namespace is created
		// by this design, so this is not a reason to fail the dry-run: pass it, returning the workload itself so
		// the component keeps its original configuration instead of a synthetic Status.
		if isMissingNamespaceError(reqErrDef, stDef, "default") {
			return k8sResource, true, nil
		}
		if stDef != nil && stDef["status"] == "Failure" {
			annotateFallbackForbidden(stDef, namespace)
			return stDef, false, nil
		}
		return stDef, okDef, reqErrDef
	}

	if !success && st != nil && st["status"] == "Failure" {
		return st, false, nil
	}
	return st, success, reqErr
}

// executeDryRunRequest constructs and executes the REST HTTP request (POST for create or DELETE for delete) to the Kubernetes API,
// returning the server's status response map, success boolean, and any encountered request error.
func executeDryRunRequest(rClient rest.Interface, k8sResource map[string]interface{}, namespace string, isDelete bool, mapper meta.RESTMapper) (st map[string]interface{}, success bool, reqErr error) {
	basePath, resourcePath, err := resolveResourcePath(k8sResource, namespace, mapper)
	if err != nil {
		return nil, false, err
	}

	var req *rest.Request
	if isDelete {
		// For delete operations, we need the resource name
		metadata, ok := k8sResource["metadata"].(map[string]interface{})
		if !ok {
			return nil, false, ErrDryRun(fmt.Errorf("metadata not found"), "metadata is required for delete operations")
		}

		name, ok := metadata["name"].(string)
		if !ok || name == "" {
			return nil, false, ErrDryRun(fmt.Errorf("resource name not found"), "resource name is required for delete operations")
		}

		// Complete path with resource name for DELETE
		path := fmt.Sprintf("%s/%s/%s", basePath, resourcePath, name)
		req = rClient.Delete().AbsPath(path).SetHeader("Accept", "application/json").Param("dryRun", "All")
	} else {
		// For create/update operations
		path := fmt.Sprintf("%s/%s", basePath, resourcePath)

		data, err := json.Marshal(k8sResource)
		if err != nil {
			return nil, false, models.ErrMarshal(err, "k8s resource")
		}

		req = rClient.Post().
			AbsPath(path).
			Body(data).
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			Param("dryRun", "All").
			Param("fieldValidation", "Strict").
			Param("fieldManager", "meshery")
	}

	// Execute the request and process response
	res := req.Do(context.Background())
	resp, reqErr := res.Raw()
	st, success, formatErr := formatDryRunResponse(resp, reqErr)
	if formatErr != nil {
		return st, success, formatErr
	}
	return st, success, reqErr
}

// resolveResourcePath determines the API basePath and resourcePath for the given resource,
// using the provided RESTMapper if available, and falling back to naive pluralization otherwise.
func resolveResourcePath(k8sResource map[string]interface{}, namespace string, mapper meta.RESTMapper) (basePath string, resourcePath string, err error) {
	kind, _ := k8sResource["kind"].(string)
	apiVersion, _ := k8sResource["apiVersion"].(string)

	if kind == "" || apiVersion == "" {
		return "", "", ErrDryRun(fmt.Errorf("invalid resource or namespace not provided"), "\"kind\" and \"apiVersion\" cannot be empty")
	}

	var mapping *meta.RESTMapping
	if mapper != nil {
		gv, parseErr := schema.ParseGroupVersion(apiVersion)
		if parseErr == nil {
			gk := schema.GroupKind{Group: gv.Group, Kind: kind}
			m, mappingErr := mapper.RESTMapping(gk, gv.Version)
			if mappingErr != nil {
				// The kind is served, but not at the requested version: say so instead of sending the object to
				// another version's endpoint, which the API server rejects with a confusing version mismatch.
				if served, servedErr := mapper.RESTMapping(gk); servedErr == nil && served != nil {
					preferred := served.Resource.GroupVersion().String()
					return "", "", ErrDryRun(
						fmt.Errorf("apiVersion %s is not served by this cluster (preferred: %s)", apiVersion, preferred),
						fmt.Sprintf("update the component's apiVersion to %s", preferred),
					)
				}
			}
			if mappingErr == nil && m != nil {
				mapping = m
			}
		}
	}

	if mapping != nil {
		resourcePath = mapping.Resource.Resource
		if mapping.Resource.Group == "" {
			basePath = fmt.Sprintf("/api/%s", mapping.Resource.Version)
		} else {
			basePath = fmt.Sprintf("/apis/%s/%s", mapping.Resource.Group, mapping.Resource.Version)
		}
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace && namespace != "" {
			basePath = fmt.Sprintf("%s/namespaces/%s", basePath, namespace)
		}
	} else {
		// Fall back to existing behavior if discovery fails or mapper is nil
		aV := apiVersion
		apiString := "api"
		if len(strings.Split(aV, "/")) > 1 {
			// for non-core resources, the endpoint should use 'apis' instead of 'api'
			apiString += "s"
		}
		resourcePath = kindToResource(kind)
		basePath = fmt.Sprintf("/%s/%s", apiString, aV)
		if namespace != "" {
			basePath = fmt.Sprintf("%s/namespaces/%s", basePath, namespace)
		}
	}

	return basePath, resourcePath, nil
}

// kindToResource converts a Kubernetes Kind to its naive lowercased plural resource name.
func kindToResource(kind string) string {
	return strings.ToLower(kind) + "s"
}

// notFoundDetails returns details.kind and details.name of a NotFound Status. The decoded response body is
// read first: client-go can return a generic StatusError with empty details while the body carries the
// Status the API server sent. The structured error is the fallback. Messages are never matched.
func notFoundDetails(err error, st map[string]interface{}) (kind, name string, ok bool) {
	if reason, _ := st["reason"].(string); reason == string(metav1.StatusReasonNotFound) {
		if details, isMap := st["details"].(map[string]interface{}); isMap {
			kind, _ = details["kind"].(string)
			name, _ = details["name"].(string)
			return kind, name, true
		}
	}
	var statusErr *kubeerror.StatusError
	if errors.As(err, &statusErr) && kubeerror.IsNotFound(statusErr) && statusErr.ErrStatus.Details != nil {
		return statusErr.ErrStatus.Details.Kind, statusErr.ErrStatus.Details.Name, true
	}
	return "", "", false
}

// isDeleteTargetNotFound reports whether a delete dry-run failed only because the object being deleted, or
// the namespace holding it, does not exist. The NotFound Status must name that resource and object, so a 404
// for a missing CRD or an unresolvable resource path is not mistaken for "not deployed".
func isDeleteTargetNotFound(err error, st map[string]interface{}, k8sResource map[string]interface{}, namespace string, mapper meta.RESTMapper) bool {
	kind, detailName, ok := notFoundDetails(err, st)
	if !ok || detailName == "" {
		return false
	}
	if isMissingNamespaceError(err, st, namespace) {
		return true
	}
	_, resource, resolveErr := resolveResourcePath(k8sResource, namespace, mapper)
	if resolveErr != nil {
		return false
	}
	metadata, _ := k8sResource["metadata"].(map[string]interface{})
	name, _ := metadata["name"].(string)
	return strings.EqualFold(kind, resource) && detailName == name
}

// isFallbackNameConflict reports whether a create dry-run failed only because an object with the same name
// already exists. Like the NotFound checks, it trusts the structured Status, never the message.
func isFallbackNameConflict(err error, st map[string]interface{}, k8sResource map[string]interface{}) bool {
	metadata, _ := k8sResource["metadata"].(map[string]interface{})
	name, _ := metadata["name"].(string)
	if name == "" {
		return false
	}
	if reason, _ := st["reason"].(string); reason == string(metav1.StatusReasonAlreadyExists) {
		details, _ := st["details"].(map[string]interface{})
		detailName, _ := details["name"].(string)
		return detailName == name
	}
	var statusErr *kubeerror.StatusError
	if errors.As(err, &statusErr) && kubeerror.IsAlreadyExists(statusErr) && statusErr.ErrStatus.Details != nil {
		return statusErr.ErrStatus.Details.Name == name
	}
	return false
}

// annotateFallbackForbidden prefixes the message of a Forbidden Status returned by the retry in "default" with
// the declared namespace and the retry, since the rejection happened in "default", a namespace the user never
// declared. Reason and code are left untouched, so callers still see a Forbidden.
func annotateFallbackForbidden(st map[string]interface{}, declaredNamespace string) {
	if reason, _ := st["reason"].(string); reason != string(metav1.StatusReasonForbidden) {
		return
	}
	message, _ := st["message"].(string)
	st["message"] = fmt.Sprintf(
		"namespace %q declared by this design does not exist on the cluster, so the dry-run was retried and rejected in namespace \"default\": %s",
		declaredNamespace, message,
	)
}

// isMissingNamespaceError reports whether the failure is a NotFound response for the target namespace itself.
// Only the structured Status is trusted: reason NotFound, details.kind naming namespaces, and details.name
// equal to targetNamespace. Messages are never matched, so unrelated failures keep their original error.
func isMissingNamespaceError(err error, st map[string]interface{}, targetNamespace string) bool {
	if targetNamespace == "" {
		return false
	}
	kind, name, ok := notFoundDetails(err, st)
	return ok && isNamespaceKind(kind) && name == targetNamespace
}

// isNamespaceKind reports whether a Status details.kind names the Namespace resource.
func isNamespaceKind(kind string) bool {
	kind = strings.ToLower(kind)
	return kind == "namespaces" || kind == "namespace"
}

// deepCopyMap creates a deep copy of an unstructured map, falling back to shallow copy on serialization error.
func deepCopyMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}
	data, err := json.Marshal(m)
	if err != nil {
		copied := make(map[string]interface{}, len(m))
		for k, v := range m {
			copied[k] = v
		}
		if metaMap, ok := m["metadata"].(map[string]interface{}); ok {
			metaCopy := make(map[string]interface{}, len(metaMap))
			for mk, mv := range metaMap {
				metaCopy[mk] = mv
			}
			copied["metadata"] = metaCopy
		}
		return copied
	}
	var copied map[string]interface{}
	_ = json.Unmarshal(data, &copied)
	return copied
}

// formatDryRunResponse parses the raw response bytes from Kubernetes into a status map and evaluates success or failure.
func formatDryRunResponse(resp []byte, err error) (status map[string]interface{}, success bool, meshkiterr error) {
	if len(resp) == 0 {
		if err != nil {
			return nil, false, ErrDryRun(err, "")
		}
		return nil, false, ErrDryRun(fmt.Errorf("nil response to dry run request to Kubernetes"), "")
	}

	e := json.Unmarshal(resp, &status)
	if e != nil {
		meshkiterr = models.ErrMarshal(e, fmt.Sprintf("Cannot serialize Status object from the server: %s", e.Error()))
		return
	}

	if status == nil || status["kind"] == nil {
		meshkiterr = ErrDryRun(fmt.Errorf("nil response to dry run request to Kubernetes"), "")
		return
	}
	if status["status"] == "Failure" { // The dryRun returns errors in the form of Status
		success = false
		return
	}
	success = true
	return
}

// getNamespaceForComponent extracts the effective namespace for a component definition, inspecting metadata and configuration.
func getNamespaceForComponent(comp *component.ComponentDefinition) string {
	namespace := ""
	isNamespaced, ok := comp.Metadata.AdditionalProperties["isNamespaced"].(bool)
	if ok && isNamespaced {
		namespace = "default"
	}

	_metadata, ok := comp.Configuration["metadata"]
	if ok {
		metadata, err := utils.Cast[map[string]interface{}](_metadata)
		if err == nil {
			_namespace, ok := metadata["namespace"]
			if ok {
				ns, _ := utils.Cast[string](_namespace)
				if ns != "" {
					namespace = ns
				}
			}

		}
	}
	return namespace
}
