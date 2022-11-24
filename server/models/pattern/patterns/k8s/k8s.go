package k8s

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
	man "github.com/layer5io/meshkit/utils/manifests"
	"gopkg.in/yaml.v2"
)

// In case of any breaking change or bug caused by this, set this to false and the whitespace addition in schema generated/consumed would be removed(will go back to default behavior)
const Format prettifier = true

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

func getAPIVersionFromComponent(comp v1alpha1.Component) string {
	return comp.Annotations["pattern.meshery.io.k8s.k8sAPIVersion"]
}

func getKindFromComponent(comp v1alpha1.Component) string {
	kind := strings.TrimPrefix(comp.Annotations["pattern.meshery.io.k8s.k8sKind"], "/")

	return kind
}

type prettifier bool

// prettifyEndString will be true in cases of schema prettification where we want to prettify everything and will be false in
// cases of YAML inputs from users where we want the end input as it is before sending to external systems such as kubernetes.
// NOTE: For clients which want a complete prettified version, prettifyEndString will be passed true.
func (p prettifier) Prettify(m map[string]interface{}, prettifyEndString bool) map[string]interface{} {
	res := ConvertMapInterfaceMapString(m, true, prettifyEndString)
	out, ok := res.(map[string]interface{})
	if !ok {
		fmt.Println("failed to cast")
	}
	return out
}
func (p prettifier) DePrettify(m map[string]interface{}, deprettifyEndString bool) map[string]interface{} {
	res := ConvertMapInterfaceMapString(m, false, deprettifyEndString)
	out, ok := res.(map[string]interface{})
	if !ok {
		fmt.Println("failed to cast")
	}

	return out
}

// ConvertMapInterfaceMapString converts map[interface{}]interface{} => map[string]interface{}
//
// It will also convert []interface{} => []string
func ConvertMapInterfaceMapString(v interface{}, prettify bool, endString bool) interface{} {
	switch x := v.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v2 := range x {
			switch k2 := k.(type) {
			case string:
				if prettify {
					m[man.FormatToReadableString(k2)] = ConvertMapInterfaceMapString(v2, prettify, endString)
				} else {
					m[man.DeFormatReadableString(k2)] = ConvertMapInterfaceMapString(v2, prettify, endString)
				}
			default:
				m[fmt.Sprint(k)] = ConvertMapInterfaceMapString(v2, prettify, endString)
			}
		}
		return m

	case []interface{}:
		x2 := make([]interface{}, len(x))
		for i, v2 := range x {
			x2[i] = ConvertMapInterfaceMapString(v2, prettify, endString)
		}
		return x2
	case map[string]interface{}:
		m := map[string]interface{}{}
		foundFormatIntOrString := false
		for k, v2 := range x {
			if prettify {
				m[man.FormatToReadableString(k)] = ConvertMapInterfaceMapString(v2, prettify, endString)
			} else {
				m[man.DeFormatReadableString(k)] = ConvertMapInterfaceMapString(v2, prettify, endString)
			}
			//Apply this fix only when the format specifies string|int and type specifies string therefore when there is a contradiction
			if k == "format" && v2 == "int-or-string" {
				foundFormatIntOrString = true
			}
		}
		if x["type"] == "string" && foundFormatIntOrString {
			m["type"] = "integer"
		}
		return m
	case string:
		if endString {
			if prettify {
				return man.FormatToReadableString(x) //Whitespace formatting should be done at the time of prettification only
			}
			return man.DeFormatReadableString(x)
		}
	}
	return v
}
