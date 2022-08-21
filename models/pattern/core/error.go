package core

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrGetK8sComponentsCode     = "2212"
	ErrParseK8sManifestCode     = "2213"
	ErrCreatePatternServiceCode = "2214"
)

func ErrGetK8sComponents(err error) error {
	return errors.New(ErrGetK8sComponentsCode, errors.Alert, []string{"Could not get K8s components for registration"}, []string{err.Error()}, []string{"Invalid kubeconfig", "Filters passed incorrectly in config", "Could not fetch API resources from Kubernetes server"}, []string{"Make sure that the configuration filters passed are in accordance with output from /openapi/v2"})
}

func ErrParseK8sManifest(err error) error {
	return errors.New(ErrParseK8sManifestCode, errors.Alert, []string{"Failed to parse manifest into JSON"}, []string{err.Error()}, []string{"Ensure manifests are valid and also use the restricted YAML features like using only \"strings\" in the field names"}, []string{"Ensure YAML is a valid Kubernetes Manifest"})
}

func ErrCreatePatternService(err error) error {
	return errors.New(ErrParseK8sManifestCode, errors.Alert, []string{"Failed to create pattern service from Manifest"}, []string{err.Error()}, []string{"Invalid Manifest", "Meshery doesn't identifies the Resource mentioned in the Manifest"}, []string{"Check if all of the meshery adapters are running", "Check if Meshery has successfully identified and registered Kubernetes components"})
}
