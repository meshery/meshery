package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/component"

	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
	_errors "github.com/pkg/errors"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/rest"
)

func Deploy(kubeClient *meshkube.Client, comp component.ComponentDefinition, isDel bool) error {
	var namespace string

	resource := createK8sResourceStructure(comp)
	manifest, err := yaml.Marshal(resource)
	if err != nil {
		return err
	}

	// Define a function to extract namesapce, labels and annotations in the componetn definiotn
	var confMetadata map[string]interface{}
	_confMetadata, ok := comp.Configuration["metadata"]
	if ok && !utils.IsInterfaceNil(_confMetadata) {
		confMetadata, err = utils.Cast[map[string]interface{}](_confMetadata)
		if err != nil {
			err = _errors.Wrapf(err, "unable to extract namespace from component configuration")
			fmt.Println("line 36 ;;;;;;;;;;;;;;;;;;;", err)
			return err
		}

		_namespace, ok := confMetadata["namespace"]
		if ok && !utils.IsInterfaceNil(_namespace) {
			namespace, _ = utils.Cast[string](_namespace)
		}
	}

	err = kubeClient.ApplyManifest(manifest, meshkube.ApplyOptions{
		Namespace: namespace,
		Update:    true,
		Delete:    isDel,
	})

	fmt.Println("line 51 ;;;;;;;;;;;;;;;;;;;", err)
	if err != nil {
		fmt.Println("line 54 ;;;;;;;;;;;;;;;;;;;", err)
		if isErrKubeStatusErr(err) {
			status, _ := json.Marshal(err)
			return formatKubeStatusErrToMeshkitErr(&status, comp.DisplayName)
		} else {
			return meshkube.ErrApplyManifest(err)
		}
	}
	fmt.Println("line 61 ;;;;;;;;;;;;;;;;;;;", err)
	return nil
}

func DryRunHelper(client *meshkube.Client, comp component.ComponentDefinition) (st map[string]interface{}, success bool, err error) {
	resource := createK8sResourceStructure(comp)
	// Define a function to extract namesapce, labels and annotations in the componetn definiotn
	_metadata, ok := comp.Configuration["metadata"]
	var namespace string
	if !ok {
		err = ErrDryRun(fmt.Errorf("unable to extract namespace from component configuration"), comp.Component.Kind)
		return
	}

	metadata, err := utils.Cast[map[string]interface{}](_metadata)
	if err != nil {
		err = _errors.Wrapf(err, "unable to extract namespace from component configuration")
		err = ErrDryRun(err, comp.Component.Kind)
		return nil, false, err
	}

	_namespace, ok := metadata["namespace"]
	if !ok {
		err = ErrDryRun(fmt.Errorf("unable to extract namespace from component configuration"), comp.Component.Kind)
		return
	}

	namespace, err = utils.Cast[string](_namespace)
	if err != nil {
		err = _errors.Wrapf(err, "unable to extract namespace from component configuration")
		err = ErrDryRun(err, comp.Component.Kind)
		return nil, false, err
	}

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

func createK8sResourceStructure(comp component.ComponentDefinition) map[string]interface{} {
	annotations := map[string]interface{}{}
	labels := map[string]interface{}{}

	_confMetadata, ok := comp.Configuration["metadata"]
	if ok {
		confMetadata, err := utils.Cast[map[string]interface{}](_confMetadata)
		if err == nil {

			_annotations, ok := confMetadata["annotations"]
			if ok {
				annotations, _ = utils.Cast[map[string]interface{}](_annotations)
			}

			_label, ok := confMetadata["labels"]

			if ok {
				labels, _ = utils.Cast[map[string]interface{}](_label)
			}
		}
	}

	component := map[string]interface{}{
		"apiVersion": comp.Component.Version,
		"kind":       comp.Component.Kind,
		"metadata": map[string]interface{}{
			"name":        comp.DisplayName,
			"annotations": annotations,
			"labels":      labels,
		},
	}

	// was/is this required, since previosuly in settings we stored kind, apiversion and metadata and namespace differently and now as well?

	for k, v := range comp.Configuration {
		if k == "apiVersion" || k == "kind" || k == "metadata" {
			continue
		}

		component[k] = v
	}
	return component
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
