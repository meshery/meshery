package k8s

import (
	"strings"

	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
	"gopkg.in/yaml.v2"
)

func Deploy(kubeClient *meshkube.Client, oamComp v1alpha1.Component, oamConfig v1alpha1.Configuration, isDel bool) error {
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

func createK8sResourceStructure(comp v1alpha1.Component) map[string]interface{} {
	apiVersion := getAPIVersionFromComponent(comp)
	kind := getKindFromComponent(comp)
	if kind == "Namespace" {
		comp.ObjectMeta.Name = comp.Namespace // For namespace the only significance of this field is to provide the name of "Namespace". Originally this field is the service name, that should not be the behavior.
	}
	component := map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":        comp.ObjectMeta.Name,
			"annotations": comp.Annotations,
			"labels":      comp.Labels,
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

func getAPIVersionFromComponent(comp v1alpha1.Component) string {
	return comp.Annotations["pattern.meshery.io.k8s.k8sAPIVersion"]
}

func getKindFromComponent(comp v1alpha1.Component) string {
	kind := strings.TrimPrefix(comp.Annotations["pattern.meshery.io.k8s.k8sKind"], "/")

	return kind
}
