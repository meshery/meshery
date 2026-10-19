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

// DryRunHelper performs server-side dry-run validation for a given component definition,
// extracting Kubernetes resource metadata and delegating to DryRunHelperWithOptions with optional DryRunOptions.
func DryRunHelper(client *meshkube.Client, comp component.ComponentDefinition, isDelete bool, opts ...DryRunOptions) (st map[string]interface{}, success bool, err error) {
	var opt DryRunOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return DryRunHelperWithOptions(client, comp, isDelete, opt)
}

// DryRunHelperWithOptions performs server-side dry-run validation for a component definition using the provided DryRunOptions.
func DryRunHelperWithOptions(client *meshkube.Client, comp component.ComponentDefinition, isDelete bool, opts DryRunOptions) (st map[string]interface{}, success bool, err error) {
	if client == nil || client.KubeClient == nil {
		return nil, false, ErrDryRun(fmt.Errorf("kubernetes client is nil"), "")
	}
	resource := converter.CreateK8sResourceStructure(&comp)
	namespace := getNamespaceForComponent(&comp)
	return dryRun(client.KubeClient.RESTClient(), resource, namespace, isDelete, opts)
}

// DryRunResource performs dry-run validation directly on an unstructured Kubernetes resource map using the provided REST client and options.
func DryRunResource(rClient rest.Interface, k8sResource map[string]interface{}, namespace string, isDelete bool, opts ...DryRunOptions) (st map[string]interface{}, success bool, err error) {
	var opt DryRunOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return dryRun(rClient, k8sResource, namespace, isDelete, opt)
}

// dryRun performs a dry-run operation on the given Kubernetes resource using the provided REST client.
// It constructs the appropriate API request based on the resource type and operation (create/update or delete).
// Returns the status response from the Kubernetes server, a success flag, and any encountered error.
func dryRun(rClient rest.Interface, k8sResource map[string]interface{}, namespace string, isDelete bool, opts DryRunOptions) (st map[string]interface{}, success bool, err error) {
	if rClient == nil {
		return nil, false, ErrDryRun(fmt.Errorf("RESTClient is nil"), "")
	}

	st, success, reqErr := executeDryRunRequest(rClient, k8sResource, namespace, isDelete, opts.RESTMapper)

	// Suppress 404 NotFound on undeploy dry-runs
	if isDelete {
		if isNotFoundError(reqErr, st) {
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
		// If validation against "default" also failed solely due to missing namespace, suppress it
		if isMissingNamespaceError(reqErrDef, stDef, "default") {
			return map[string]interface{}{"status": "Success"}, true, nil
		}
		// Otherwise return workload validation errors found when testing against default
		if stDef != nil && stDef["status"] == "Failure" {
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
	basePath, resourcePath, err := ResolveResourcePath(k8sResource, namespace, mapper)
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

// ResolveResourcePath determines the API basePath and resourcePath for the given resource,
// using the provided RESTMapper if available, and falling back to naive pluralization otherwise.
func ResolveResourcePath(k8sResource map[string]interface{}, namespace string, mapper meta.RESTMapper) (basePath string, resourcePath string, err error) {
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
			if mappingErr != nil && gv.Version != "" {
				m, mappingErr = mapper.RESTMapping(gk)
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

// isNotFoundError checks whether an error or status response represents a Kubernetes 404/NotFound error.
func isNotFoundError(err error, st map[string]interface{}) bool {
	if err != nil && kubeerror.IsNotFound(err) {
		return true
	}
	if st != nil {
		if reason, ok := st["reason"].(string); ok && reason == string(metav1.StatusReasonNotFound) {
			return true
		}
		if code, ok := st["code"].(float64); ok && int(code) == 404 {
			return true
		}
		if code, ok := st["code"].(int); ok && code == 404 {
			return true
		}
	}
	return false
}

// isMissingNamespaceError reports whether the failure is a NotFound response for the target namespace itself.
// Only the structured Status is trusted: reason NotFound, details.kind naming namespaces, and details.name
// equal to targetNamespace. Messages are never matched, so unrelated failures keep their original error.
func isMissingNamespaceError(err error, st map[string]interface{}, targetNamespace string) bool {
	if targetNamespace == "" {
		return false
	}
	var statusErr *kubeerror.StatusError
	if errors.As(err, &statusErr) && kubeerror.IsNotFound(statusErr) {
		if details := statusErr.ErrStatus.Details; details != nil &&
			isNamespaceKind(details.Kind) && details.Name == targetNamespace {
			return true
		}
	}
	if st != nil {
		reason, _ := st["reason"].(string)
		if reason != string(metav1.StatusReasonNotFound) {
			return false
		}
		details, ok := st["details"].(map[string]interface{})
		if !ok {
			return false
		}
		kind, _ := details["kind"].(string)
		name, _ := details["name"].(string)
		return isNamespaceKind(kind) && name == targetNamespace
	}
	return false
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
	if namespace == "" || namespace == "default" {
		if _namespace, ok := comp.Configuration["namespace"]; ok {
			if ns, _ := utils.Cast[string](_namespace); ns != "" {
				namespace = ns
			}
		}
	}
	return namespace
}
