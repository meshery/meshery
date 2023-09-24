package helpers

import (
	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrErrNewDynamicClientGeneratorCode    = "1200"
	ErrInvalidK8SConfigCode                = "1201"
	ErrClientConfigCode                    = "1202"
	ErrFetchKubernetesNodesCode            = "1203"
	ErrFetchNodesCode                      = "1204"
	ErrFetchKubernetesVersionCode          = "1205"
	ErrScanKubernetesCode                  = "1206"
	ErrRetrievePodListCode                 = "1207"
	ErrDetectServiceForDeploymentImageCode = "1208"
	ErrRetrieveNamespacesListCode          = "1209"
	ErrGetNamespaceDeploymentsCode         = "1210"
	ErrDetectServiceWithNameCode           = "1211"
	ErrGeneratingLoadTestCode              = "1212"
	ErrRunningTestCode                     = "1213"
	ErrConvertingResultToMapCode           = "1214"
	ErrGrpcSupportCode                     = "1216"
	ErrStartingNighthawkServerCode         = "1217"
	ErrTransformingDataCode                = "1218"
	ErrRunningNighthawkServerCode          = "1219"
	ErrAddAndValidateExtraHeaderCode       = "1220"
	ErrInClusterConfigCode                 = "1221"
	ErrNewKubeClientGeneratorCode          = "1222"
	ErrRestConfigFromKubeConfigCode        = "1223"
	ErrNewKubeClientCode                   = "1224"
	ErrAdapterInsufficientInformationCode  = "1226"
	ErrDeployingAdapterInK8sEnvCode        = "1227"
	ErrUnDeployingAdapterInK8sEnvCode      = "1228"
	ErrDeployingAdapterInDockerEnvCode     = "1229"
	ErrUnDeployingAdapterInDockerEnvCode   = "1230"
	ErrDeployingAdapterCode                = "1231"
	ErrUnDeployingAdapterCode              = "1232"
	ErrClientSetCode                       = "1233"
)

func ErrNewDynamicClientGenerator(err error) error {
	return errors.New(ErrErrNewDynamicClientGeneratorCode, errors.Alert, []string{"Unable to generate the dynamic client generator"}, []string{err.Error()}, []string{"Kubernetes config is not accessible to meshery or not valid"}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
}

func ErrInvalidK8SConfig(err error) error {
	return errors.New(ErrInvalidK8SConfigCode, errors.Alert, []string{"No valid kubernetes config found"}, []string{err.Error()}, []string{"Kubernetes config is not accessible to meshery or not valid"}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
}

func ErrClientConfig(err error) error {
	return errors.New(ErrClientConfigCode, errors.Alert, []string{"Unable to create client config"}, []string{err.Error()}, []string{"Kubernetes config is not accessible to meshery or not valid"}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
}

func ErrFetchKubernetesNodes(err error) error {
	return errors.New(ErrFetchKubernetesNodesCode, errors.Alert, []string{"Unable to fetch kubernetes nodes"}, []string{err.Error()}, []string{"Kubernetes API server might not be reachable from the Meshery server"}, []string{"Make sure kubernetes API server is reachable from meshery server"})
}

func ErrFetchNodes(err error) error {
	return errors.New(ErrFetchNodesCode, errors.Alert, []string{"Unable to get the list of nodes"}, []string{err.Error()}, []string{"Kubernetes API server might not be reachable from the Meshery server"}, []string{"Make sure kubernetes API server is reachable from meshery server"})
}

func ErrFetchKubernetesVersion(err error) error {
	return errors.New(ErrFetchKubernetesVersionCode, errors.Alert, []string{"Unable to fetch kubernetes version"}, []string{err.Error()}, []string{"Kubernetes API server might not be reachable from the Meshery server"}, []string{"Make sure kubernetes API server is reachable from meshery server"})
}

func ErrScanKubernetes(err error) error {
	return errors.New(ErrScanKubernetesCode, errors.Alert, []string{"Unable to scan kubernetes"}, []string{err.Error()}, []string{"Kubernetes API server might not be reachable from the Meshery server"}, []string{"Make sure kubernetes API server is reachable from meshery server"})
}

func ErrRetrievePodList(err error) error {
	return errors.New(ErrRetrievePodListCode, errors.Alert, []string{"Unable to retrieve pod list"}, []string{err.Error()}, []string{"Kubernetes API server might not be reachable from the Meshery server", "Requested resource might not be available"}, []string{"Make sure kubernetes API server is reachable from meshery server", "Make sure you are requesting for a valid resource"})
}

func ErrDetectServiceForDeploymentImage(err error) error {
	return errors.New(ErrDetectServiceForDeploymentImageCode, errors.Alert, []string{"Unable to detect service for deployment image"}, []string{err.Error()}, []string{"Kubernetes API server might not be reachable from the Meshery server", "Requested resource might not be available"}, []string{"Make sure kubernetes API server is reachable from meshery server", "Make sure you are requesting for a valid resource"})
}

func ErrRetrieveNamespacesList(err error) error {
	return errors.New(ErrRetrieveNamespacesListCode, errors.Alert, []string{"Unable to get the list of namespaces"}, []string{err.Error()}, []string{"Kubernetes API server might not be reachable from the Meshery server", "Requested resource might not be available"}, []string{"Make sure kubernetes API server is reachable from meshery server", "Make sure you are requesting for a valid resource"})
}

func ErrGetNamespaceDeployments(err error, obj string) error {
	return errors.New(ErrGetNamespaceDeploymentsCode, errors.Alert, []string{"Unable to get deployments in the ", obj, "namespace"}, []string{err.Error()}, []string{"Kubernetes API server might not be reachable from the Meshery server", "Requested resource might not be available"}, []string{"Make sure kubernetes API server is reachable from meshery server", "Make sure you are requesting for a valid resource"})
}

func ErrDetectServiceWithName(err error) error {
	return errors.New(ErrDetectServiceWithNameCode, errors.Alert, []string{"Unable to get services from the cluster with the name given in names parameter"}, []string{err.Error()}, []string{"Kubernetes API server might not be reachable from the Meshery server", "Requested resource might not reachable from Meshery server"}, []string{"Make sure kubernetes API server is reachable from meshery server", "Make sure the network connectivity is up between meshery server and the service endpoint"})
}

func ErrGeneratingLoadTest(err error) error {
	return errors.New(ErrGeneratingLoadTestCode, errors.Alert, []string{"Unable to generate load test"}, []string{err.Error()}, []string{}, []string{})
}

func ErrRunningTest(err error) error {
	return errors.New(ErrRunningTestCode, errors.Alert, []string{"Unable to run test"}, []string{err.Error()}, []string{}, []string{})
}

func ErrConvertingResultToMap(err error) error {
	return errors.New(ErrConvertingResultToMapCode, errors.Alert, []string{"Unable to convert from the result to map"}, []string{err.Error()}, []string{}, []string{})
}

func ErrGrpcSupport(err error, obj string) error {
	return errors.New(ErrGrpcSupportCode, errors.Alert, []string{obj, " does not support gRPC load testing"}, []string{err.Error()}, []string{}, []string{})
}

func ErrTransformingData(err error) error {
	return errors.New(ErrTransformingDataCode, errors.Alert, []string{"Error while transforming data"}, []string{err.Error()}, []string{}, []string{})
}

func ErrRunningNighthawkServer(err error) error {
	return errors.New(ErrRunningNighthawkServerCode, errors.Alert, []string{"Error while running nighthawk server"}, []string{err.Error()}, []string{}, []string{})
}

func ErrAddAndValidateExtraHeader(err error) error {
	return errors.New(ErrAddAndValidateExtraHeaderCode, errors.Alert, []string{"Unable to add and validate extra header"}, []string{err.Error()}, []string{}, []string{})
}

func ErrInClusterConfig(err error) error {
	return errors.New(ErrInClusterConfigCode, errors.Alert, []string{"Unable to load in-cluster kubeconfig"}, []string{err.Error()}, []string{}, []string{})
}

func ErrNewKubeClientGenerator(err error) error {
	return errors.New(ErrNewKubeClientGeneratorCode, errors.Alert, []string{"Unable to generate new kube dynamic client"}, []string{err.Error()}, []string{}, []string{})
}

func ErrRestConfigFromKubeConfig(err error) error {
	return errors.New(ErrRestConfigFromKubeConfigCode, errors.Alert, []string{"Unable to create rest config from kube congif"}, []string{err.Error()}, []string{}, []string{})
}

func ErrClientSet(err error) error {
	return errors.New(ErrClientSetCode, errors.Alert, []string{"Unable to create client set"}, []string{err.Error()}, []string{}, []string{})
}

func ErrStartingNighthawkServer(err error) error {
	return errors.New(ErrStartingNighthawkServerCode, errors.Alert, []string{"Unable to start the nighthawk server"}, []string{err.Error()}, []string{}, []string{})
}

func ErrNewKubeClient(err error) error {
	return errors.New(ErrNewKubeClientCode, errors.Alert, []string{"Unable to create new kube client"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDeployingAdapterInK8s(err error) error {
	return errors.New(ErrDeployingAdapterInK8sEnvCode, errors.Critical, []string{"Unable to deploy adapter in k8s env"}, []string{err.Error()}, []string{"Possible issues with Kubernetes cluster configuration or network connectivity."}, []string{"Check the Kubernetes cluster's configuration, ensure necessary resources are available, and verify network connectivity."})
}

func ErrUnDeployingAdapterInK8s(err error) error {
	return errors.New(ErrUnDeployingAdapterInK8sEnvCode, errors.Critical, []string{"Unable to undeploy adapter in k8s env"}, []string{err.Error()}, []string{"Possible issues with Kubernetes cluster configuration or network connectivity."}, []string{"Check the Kubernetes cluster's configuration, ensure necessary resources are available, and verify network connectivity."})
}

func ErrDeployingAdapterInDocker(err error) error {
	return errors.New(ErrDeployingAdapterInDockerEnvCode, errors.Critical, []string{"Unable to deploy adapter in k8s env"}, []string{err.Error()}, []string{"Possible issues with Docker configuration, container networking, or resource availability."}, []string{"Check Docker configuration settings, ensure containers have proper networking access, and verify available resources."})
}

func ErrUnDeployingAdapterInDocker(err error) error {
	return errors.New(ErrUnDeployingAdapterInDockerEnvCode, errors.Critical, []string{"Unable to undeploy Meshery Adapter in Kubernetes environment"}, []string{err.Error()}, []string{"Possible issues with Docker configuration, container networking, or resource availability."}, []string{"Check Docker configuration settings, ensure containers have proper networking access, and verify available resources."})
}

func ErrDeployingAdapterInUnknownPlatform(err error) error {
	return errors.New(ErrDeployingAdapterCode, errors.Critical, []string{"Unable to deploy Meshery Adapter in the current environment"}, []string{err.Error()}, []string{"Your platform is not supported for deploying Meshery Adapters"}, []string{"Consider using a supported platform for deploying Meshery Adapters"})
}

func ErrUnDeployingAdapterInUnknownPlatform(err error) error {
	return errors.New(ErrUnDeployingAdapterCode, errors.Critical, []string{"Unable to undeploy Meshery Adapter in the current environment"}, []string{err.Error()}, []string{"Current platform is not supported for undeploying Meshery Adapters"}, []string{"Consider using a supported platform for undeploying Meshery Adapters"})
}
