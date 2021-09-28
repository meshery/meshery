package core

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrGetK8sComponentsCode = "10000"
)

func ErrGetK8sComponents(err error) error {
	return errors.New(ErrGetK8sComponentsCode, errors.Alert, []string{"Could not get K8s components for registeration"}, []string{err.Error()}, []string{"Got an invalid kubeconfig", "Filters passed incorrectly in config", "Could not fetch api resources from kubernetes server"}, []string{"Make sure that the configuration filters passed are in accordance with output from /openapi/v2"})
}
