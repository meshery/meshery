package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/converter"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/component"

	meshkube "github.com/meshery/meshkit/utils/kubernetes"
	"gopkg.in/yaml.v2"
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

func DryRunHelper(client *meshkube.Client, comp component.ComponentDefinition, isDelete bool) (st map[string]interface{}, success bool, err error) {
	resource := converter.CreateK8sResourceStructure(&comp)
	// Define a function to extract namesapce, labels and annotations in the componetn definiotn
	namespace := getNamespaceForComponent(&comp)
	return dryRun(client.KubeClient.RESTClient(), resource, namespace, isDelete)
}

// dryRun performs a dry-run operation on the given Kubernetes resource using the provided REST client.
// It constructs the appropriate API request based on the resource type and operation (create/update or delete).
// Returns the status response from the Kubernetes server, a success flag, and any encountered error.
// TODO: add more tests for this function
func dryRun(rClient rest.Interface, k8sResource map[string]interface{}, namespace string, isDelete bool) (st map[string]interface{}, success bool, err error) {
	// Validate resource has required fields
	if k8sResource["kind"] == "" || k8sResource["apiVersion"] == "" {
		err = ErrDryRun(fmt.Errorf("invalid resource or namespace not provided"), "\"kind\" and \"apiVersion\" cannot be empty")
		return
	}

	// Determine API path prefix (api vs apis)
	aV := k8sResource["apiVersion"].(string)
	apiString := "api"
	if len(strings.Split(aV, "/")) > 1 {
		// for non-core resources, the endpoint should use 'apis' instead of 'api'
		apiString += "s"
	}

	// Build the base resource path
	resourcePath := kindToResource(k8sResource["kind"].(string))
	basePath := fmt.Sprintf("/%s/%s", apiString, aV)

	// Add namespace to path if provided
	if namespace != "" {
		basePath = fmt.Sprintf("%s/namespaces/%s", basePath, namespace)
	}

	var req *rest.Request

	if isDelete {
		// For delete operations, we need the resource name
		metadata, ok := k8sResource["metadata"].(map[string]interface{})
		if !ok {
			err = ErrDryRun(fmt.Errorf("metadata not found"), "metadata is required for delete operations")
			return
		}

		name, ok := metadata["name"].(string)
		if !ok || name == "" {
			err = ErrDryRun(fmt.Errorf("resource name not found"), "resource name is required for delete operations")
			return
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
	resp, err := res.Raw()
	return formatDryRunResponse(resp, err)
}

func kindToResource(kind string) string {
	return strings.ToLower(kind) + "s"
}

func formatDryRunResponse(resp []byte, err error) (status map[string]interface{}, success bool, meshkiterr error) {
	e := json.Unmarshal(resp, &status)
	if e != nil {
		meshkiterr = models.ErrMarshal(e, fmt.Sprintf("Cannot serialize Status object from the server: %s", e.Error()))
		return
	}

	if status == nil || status["kind"] == nil {
		meshkiterr = ErrDryRun(fmt.Errorf("nil response to dry run request to Kubernetes"), "")
	}
	if status["status"] == "Failure" { // The dryRun returns errors in the form of Status
		success = false
		return
	}
	success = true
	return
}

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
