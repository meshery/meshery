package helpers

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrErrNewDynamicClientGeneratorCode = "replace_me"
	ErrInvalidK8SConfigCode             = "replace_me1"
	ErrErrCreateClientCode              = "replace_me2"
	ErrFetchKubernetesNodesCode         = "replace_me3"
	ErrFetchNodesCode                   = "replace_me4"
	ErrFetchKubernetesVersionCode       = "replace_me5"
)

func ErrNewDynamicClientGenerator(err error) error {
	return errors.New(ErrErrNewDynamicClientGeneratorCode, errors.Alert, []string{"Unable to generate the dynamic client generator"}, []string{err.Error()}, []string{}, []string{})
}

func ErrInvalidK8SConfig(err error) error {
	return errors.New(ErrInvalidK8SConfigCode, errors.Alert, []string{"No valid kubernetes config found"}, []string{err.Error()}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
}

func ErrCreateClient(err error, obj string) error {
	return errors.New(ErrErrCreateClientCode, errors.Alert, []string{"Unable to create client ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrFetchKubernetesNodes(err error) error {
	return errors.New(ErrFetchKubernetesNodesCode, errors.Alert, []string{"Unable to fetch kubernetes nodes"}, []string{err.Error()}, []string{}, []string{})
}

func ErrFetchNodes(err error) error {
	return errors.New(ErrFetchNodesCode, errors.Alert, []string{"Unable to get the list of nodes"}, []string{err.Error()}, []string{}, []string{})
}

func ErrFetchKubernetesVersion(err error) error {
	return errors.New(ErrFetchKubernetesVersionCode, errors.Alert, []string{"Unable to fetch kubernetes version"}, []string{err.Error()}, []string{}, []string{})
}
