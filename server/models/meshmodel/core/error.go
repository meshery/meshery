package core

import (
	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrGetK8sComponentsCode = "2212"
)

func ErrGetK8sComponents(err error) error {
	return errors.New(ErrGetK8sComponentsCode, errors.Alert, []string{"Could not get K8s components for registration"}, []string{err.Error()}, []string{"Invalid kubeconfig", "Filters passed incorrectly in config", "Could not fetch API resources from Kubernetes server"}, []string{"Make sure that the configuration filters passed are in accordance with output from /openapi/v2"})
}
