package helpers

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrErrNewDynamicClientGeneratorCode    = "replace_me"
	ErrInvalidK8SConfigCode                = "replace_me1"
	ErrErrCreateClientCode                 = "replace_me2"
	ErrFetchKubernetesNodesCode            = "replace_me3"
	ErrFetchNodesCode                      = "replace_me4"
	ErrFetchKubernetesVersionCode          = "replace_me5"
	ErrScanKubernetesCode                  = "replace_me6"
	ErrRetrievePodListCode                 = "replace_me7"
	ErrDetectServiceForDeploymentImageCode = "replace_me8"
	ErrRetrieveNamespacesListCode          = "replace_me9"
	ErrGetNamespaceDeploymentsCode         = "replace_me10"
	ErrDetectServiceWithNameCode           = "replace_me11"
	ErrGeneratingLoadTestCode              = "replace_me12"
	ErrRunningTestCode                     = "replace_me13"
	ErrConvertingResultToMapCode           = "replace_me14"
	ErrUnmarshalCode                       = "replace_me15"
	ErrGrpcSupportCode                     = "replace_me16"
	ErrStartingNighthawkServerCode         = "replace_me17"
	ErrTransformingDataCode                = "replace_me18"
	ErrRunningNighthawkServerCode          = "replace_me19"
	ErrAddAndValidateExtraHeaderCode       = "replace_me20"
)

var (
	ErrStartingNighthawkServer = errors.New(ErrStartingNighthawkServerCode, errors.Alert, []string{"unable to start nighthawk server"}, []string{}, []string{}, []string{})
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

func ErrScanKubernetes(err error) error {
	return errors.New(ErrScanKubernetesCode, errors.Alert, []string{"Unable to scan kubernetes"}, []string{err.Error()}, []string{}, []string{})
}

func ErrRetrievePodList(err error) error {
	return errors.New(ErrRetrievePodListCode, errors.Alert, []string{"Unable to retrieve pod list"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDetectServiceForDeploymentImage(err error) error {
	return errors.New(ErrDetectServiceForDeploymentImageCode, errors.Alert, []string{"Unable to detect service for deployment image"}, []string{err.Error()}, []string{}, []string{})
}

func ErrRetrieveNamespacesList(err error) error {
	return errors.New(ErrRetrieveNamespacesListCode, errors.Alert, []string{"unable to get the list of namespaces"}, []string{err.Error()}, []string{}, []string{})
}

func ErrGetNamespaceDeployments(err error, obj string) error {
	return errors.New(ErrGetNamespaceDeploymentsCode, errors.Alert, []string{"unable to get deployments in the ", obj, "namespace"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDetectServiceWithName(err error) error {
	return errors.New(ErrDetectServiceWithNameCode, errors.Alert, []string{"Unable to get services from the cluster with the name given in names parameter"}, []string{err.Error()}, []string{}, []string{})
}

func ErrGeneratingLoadTest(err error) error {
	return errors.New(ErrGeneratingLoadTestCode, errors.Alert, []string{}, []string{err.Error()}, []string{}, []string{})
}

func ErrRunningTest(err error) error {
	return errors.New(ErrRunningTestCode, errors.Alert, []string{"Unable to run test"}, []string{err.Error()}, []string{}, []string{})
}

func ErrConvertingResultToMap(err error) error {
	return errors.New(ErrConvertingResultToMapCode, errors.Alert, []string{"Unable to convert from the result to map"}, []string{err.Error()}, []string{}, []string{})
}

func ErrUnmarshal(err error, obj string) error {
	return errors.New(ErrUnmarshalCode, errors.Alert, []string{"Unable to unmarshal the : ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrGrpcSupport(err error, obj string) error {
	return errors.New(ErrGrpcSupportCode, errors.Alert, []string{obj, " does not support gRPC load testing"}, []string{err.Error()}, []string{}, []string{})
}

func ErrTransformingData(err error) error {
	return errors.New(ErrTransformingDataCode, errors.Alert, []string{"error while transforming data"}, []string{err.Error()}, []string{}, []string{})
}

func ErrRunningNighthawkServer(err error) error {
	return errors.New(ErrRunningNighthawkServerCode, errors.Alert, []string{"error while running nighthawk server"}, []string{err.Error()}, []string{}, []string{})
}

func ErrAddAndValidateExtraHeader(err error) error {
	return errors.New(ErrAddAndValidateExtraHeaderCode, errors.Alert, []string{}, []string{err.Error()}, []string{}, []string{})
}
