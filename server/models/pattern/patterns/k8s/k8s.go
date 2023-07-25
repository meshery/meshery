package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/rest"
)

func Deploy(kubeClient *meshkube.Client, oamComp v1alpha1.Component, _ v1alpha1.Configuration, isDel bool) error {
	resource := createK8sResourceStructure(oamComp)
	manifest, err := yaml.Marshal(resource)
	if err != nil {
		return err
	}
	return kubeClient.ApplyManifest(manifest, meshkube.ApplyOptions{
		Namespace: oamComp.Namespace,
		Update:    true,
		Delete:    isDel,
	})
}

func DryRunHelper(client *meshkube.Client, comp v1alpha1.Component) (st map[string]interface{}, success bool, err error) {
	resource := createK8sResourceStructure(comp)
	return dryRun(client.KubeClient.RESTClient(), resource, comp.Namespace)
}

// does dry-run on the kubernetes server for the given k8s resource and returns the Status object returned by the k8s server
// TODO: add more tests for this function
func dryRun(rClient rest.Interface, k8sResource map[string]interface{}, namespace string) (st map[string]interface{}, success bool, err error) {
	if k8sResource["kind"] == "" || k8sResource["apiVersion"] == "" {
		err = fmt.Errorf("invalid resource or namespace not provided")
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
		return
	}
	req := rClient.Post().AbsPath(path).Body(data).SetHeader("Content-Type", "application/json").SetHeader("Accept", "application/json").Param("dryRun", "All").Param("fieldValidation", "Strict").Param("fieldManager", "meshery")
	res := req.Do(context.Background())
	if res.Error() != nil {
		err = res.Error()
		return
	}
	// ignoring the error since this client-go treats failure of dryRun as an error
	resp, _ := res.Raw()
	e := json.Unmarshal(resp, &st)
	if e != nil {
		err = fmt.Errorf("cannot serialize Status object from the server: %s", e.Error())
		return
	}
	if st == nil || st["kind"] == nil {
		err = fmt.Errorf("nil response for dryRun from kubernetes")
	}
	if st["status"] == "Failure" { //The dryRun returned errors in the form of Status
		success = false
		return
	}
	success = true
	return
}

func kindToResource(kind string) string {
	return strings.ToLower(kind) + "s"
}

func createK8sResourceStructure(comp v1alpha1.Component) map[string]interface{} {
	apiVersion := v1alpha1.GetAPIVersionFromComponent(comp)
	kind := v1alpha1.GetKindFromComponent(comp)

	component := map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":        comp.ObjectMeta.Name,
			"annotations": comp.ObjectMeta.Annotations,
			"labels":      comp.ObjectMeta.Labels,
		},
	}

	for k, v := range comp.Spec.Settings {
		if k == "apiVersion" || k == "kind" || k == "metadata" {
			continue
		}

		component[k] = v
	}
	return component
}
