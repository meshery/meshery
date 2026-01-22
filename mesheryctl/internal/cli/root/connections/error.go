package connections

import (
	"fmt"
	"strings"

	"github.com/meshery/meshkit/errors"
)

var (
	ErrAwsCliNotFoundCode         = "mesheryctl-1174"
	ErrAwsEksGetCredentialsCode   = "mesheryctl-1186"
	ErrAzureAksGetCredentialsCode = "mesheryctl-1189"
	ErrAzureCliNotFoundCode       = "mesheryctl-1190"
	ErrConnectionTypeCode         = "mesheryctl-1182"
	ErrGcpGKEGetCredentialsCode   = "mesheryctl-1175"
	ErrReadKubeConfigCode         = "mesheryctl-1187"
	ErrWriteKubeConfigCode        = "mesheryctl-1188"

	invalidOutputFormatMsg = "output-format choice is invalid, use [json|yaml]"
)

func errInvalidConnectionType(connectionType string) error {
	return errors.New(ErrConnectionTypeCode, errors.Alert, []string{fmt.Sprintf("Invalid connection type %s provided", connectionType)}, []string{fmt.Sprintf("the provided connection type is not yet supported. Supported type(s) are: %s", strings.Join(supportedConnectionTypes, ", "))}, []string{}, []string{fmt.Sprintf("Please provide a valid connection type: %s", strings.Join(supportedConnectionTypes, ", "))})
}

func errReadKubeConfig(err error) error {
	return errors.New(ErrReadKubeConfigCode, errors.Alert, []string{"Unable to read kubeconfig file"}, []string{"There was an error reading the kubeconfig file"}, []string{err.Error()}, []string{"Please ensure that the kubeconfig file exists, content is valid and is accessible"})
}

func errWriteKubeConfig(err error) error {
	return errors.New(ErrWriteKubeConfigCode, errors.Alert, []string{"Unable to write kubeconfig file"}, []string{"There was an error writing the kubeconfig file"}, []string{err.Error()}, []string{"Please ensure that you have the necessary permissions to write to the kubeconfig file location"})
}

func errAzureCliNotFound(err error) error {
	return errors.New(ErrAzureCliNotFoundCode, errors.Alert, []string{"Azure CLI not found"}, []string{"The Azure CLI is required to create an AKS connection but was not found on your system"}, []string{err.Error()}, []string{"Please install the Azure CLI from https://learn.microsoft.com/en-us/cli/azure/install-azure-cli and ensure it is accessible in your system's PATH"})
}

func errAzureAksGetCredentials(err error) error {
	return errors.New(ErrAzureAksGetCredentialsCode, errors.Alert, []string{"Unable to get AKS cluster credentials"}, []string{"There was an error while fetching the AKS cluster credentials"}, []string{err.Error()}, []string{"Ensure that the AKS cluster name and resource group are correct, and that you have the necessary permissions to access the cluster"})
}

func errAwsCliNotFound(err error) error {
	return errors.New(ErrAwsCliNotFoundCode, errors.Alert, []string{"AWS CLI not found"}, []string{"The AWS CLI is required to create an EKS connection but was not found on your system"}, []string{err.Error()}, []string{"Please install the AWS CLI from https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html and ensure it is accessible in your system's PATH"})
}

func errAwsEksGetCredentials(err error) error {
	return errors.New(ErrAwsEksGetCredentialsCode, errors.Alert, []string{"Unable to get EKS cluster credentials"}, []string{"There was an error while fetching the EKS cluster credentials"}, []string{err.Error()}, []string{"Ensure that the EKS cluster name and region are correct, and that you have the necessary permissions to access the cluster"})
}

func errGcpGKEGetCredentials(err error) error {
	return errors.New(ErrGcpGKEGetCredentialsCode, errors.Alert, []string{"Unable to get GKE cluster credentials"}, []string{"There was an error while fetching the GKE cluster credentials"}, []string{err.Error()}, []string{"Ensure that the GKE cluster name, zone, and project ID are correct, and that you have the necessary permissions to access the cluster"})
}
