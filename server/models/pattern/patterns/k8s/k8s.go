package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/converter"
	"github.com/layer5io/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/component"

	meshkube "github.com/layer5io/meshkit/utils/kubernetes"

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
		if isErrKubeStatusErr(err) {
			status, _ := json.Marshal(err)
			return formatKubeStatusErrToMeshkitErr(&status, comp.DisplayName)
		} else {
			return meshkube.ErrApplyManifest(err)
		}
	}
	return nil
}

func DryRunHelper(client *meshkube.Client, comp component.ComponentDefinition) (st map[string]interface{}, success bool, err error) {
	resource := converter.CreateK8sResourceStructure(&comp)
	// Define a function to extract namesapce, labels and annotations in the componetn definiotn
	namespace := getNamespaceForComponent(&comp)
	return dryRun(client.KubeClient.RESTClient(), resource, namespace)
}

// does dry-run on the kubernetes server for the given k8s resource and returns the Status object returned by the k8s server
// TODO: add more tests for this function
func dryRun(rClient rest.Interface, k8sResource map[string]interface{}, namespace string) (st map[string]interface{}, success bool, err error) {
	if k8sResource["kind"] == "" || k8sResource["apiVersion"] == "" {
		err = ErrDryRun(fmt.Errorf("invalid resource or namespace not provided"), "\"kind\" and \"apiVersion\" cannot be empty")
		return
	}

	aV := k8sResource["apiVersion"].(string)
	// for non-core resources, the endpoint should use 'apis' instead of 'api'
	apiString := "api"
	if len(strings.Split(aV, "/")) > 1 {
		apiString = apiString + "s"
	}

	var path string

	if namespace != "" {
		path = fmt.Sprintf("/%s/%s/namespaces/%s/%s", apiString, aV, namespace, kindToResource(k8sResource["kind"].(string)))
	} else {
		path = fmt.Sprintf("/%s/%s/%s", apiString, aV, kindToResource(k8sResource["kind"].(string)))
	}

	data, err := json.Marshal(k8sResource)
	if err != nil {
		err = models.ErrMarshal(err, "k8s resource")
		return
	}

	req := rClient.Post().AbsPath(path).Body(data).SetHeader("Content-Type", "application/json").SetHeader("Accept", "application/json").Param("dryRun", "All").Param("fieldValidation", "Strict").Param("fieldManager", "meshery")
	res := req.Do(context.Background())

	// ignoring the error since this client-go treats failure of dryRun as an error
	resp, err := res.Raw()

	st, success, err = formatDryRunResponse(resp, err)

	return
}

func kindToResource(kind string) string {
	return strings.ToLower(kind) + "s"
}

func formatDryRunResponse(resp []byte, err error) (status map[string]interface{}, success bool, meshkiterr error) {

	e := json.Unmarshal(resp, &status)
	if e != nil {
		meshkiterr = models.ErrMarshal(err, fmt.Sprintf("cannot serialize Status object from the server: %s", e.Error()))
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
