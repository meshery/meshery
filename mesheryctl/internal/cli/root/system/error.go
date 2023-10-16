// Copyright 2023 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"fmt"
	"strconv"

	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrHealthCheckFailedCode             = "1133"
	ErrDownloadFileCode                  = "1134"
	ErrStopMesheryCode                   = "1135"
	ErrResetMeshconfigCode               = "1136"
	ErrApplyManifestCode                 = "1137"
	ErrApplyOperatorManifestCode         = "1138"
	ErrCreateDirCode                     = "1139"
	ErrUnsupportedPlatformCode           = "1140"
	ErrRetrievingCurrentContextCode      = "1141"
	ErrSettingDefaultContextToConfigCode = "1142"
	ErrSettingTemporaryContextCode       = "1143"
	ErrCreateManifestsFolderCode         = "1144"
	ErrRestartMesheryCode                = "1145"
	ErrK8sQueryCode                      = "1146"
	ErrK8sConfigCode                     = "1147"
	ErrInitPortForwardCode               = "1148"
	ErrRunPortForwardCode                = "1149"
	ErrFailedGetEphemeralPortCode        = "1150"
	ErrUnmarshalDockerComposeCode        = "1151"
	ErrCreatingDockerClientCode          = "1152"
	ErrWriteConfigCode                   = "1153"
	ErrContextContentCode                = "1154"
	ErrSwitchChannelResponseCode         = "1155"
	ErrGetCurrentContextCode             = "1156"
	ErrSetCurrentContextCode             = "1157"
	ErrTokenContextCode                  = "1158"
	ErrProviderInfoCode                  = "1159"
	ErrValidProviderCode                 = "1160"
	ErrUnmarshallConfigCode              = "1161"
	ErrUploadFileParamsCode              = "1162"
)

var (
	cmdType     string
	contextdocs string = "See https://docs.meshery.io/reference/mesheryctl/system/context for usage details."
	contextDir  string = "see that you have a correct context in your  meshconfig at `$HOME/.meshery/config.yaml`."
)

// A Format reference that returns Mesheryctl's URL docs for system command and sub commands
func FormatErrorReference() string {
	baseURL := "https://docs.meshery.io/reference/mesheryctl/system"
	switch cmdType {
	case "channel":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/channel")
	case "context":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/context")
	case "config":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/context")
	case "dashboard":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/dashboard")
	case "login":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/login")
	case "logout":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/logout")
	case "logs":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/logs")
	case "provider":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/provider")
	case "reset":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/reset")
	case "restart":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/restart")
	case "start":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/start")
	case "status":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/status")
	case "stop":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/stop")
	case "token":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/token")
	case "update":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/update")
	}
	return fmt.Sprintf("\nSee %s for usage details\n", baseURL)
}

func ErrHealthCheckFailed(err error) error {
	return errors.New(ErrHealthCheckFailedCode,
		errors.Alert,
		[]string{"Health checks failed"},
		[]string{"Failed to initialize healthchecker" + err.Error()},
		[]string{"Health checks execution failed in starting Meshery server successfully"},
		[]string{"Ensure Mesheryctl is running and has the right configurations."})
}

func ErrDownloadFile(err error, obj string) error {
	return errors.New(ErrDownloadFileCode, errors.Alert, []string{"Error downloading file ", obj}, []string{err.Error()}, []string{"Failed to download docker-compose or manifest file due to system/config/network issues"}, []string{"Make sure docker-compose or manifest file is downloaded successfully"})
}

func ErrStopMeshery(err error) error {
	return errors.New(ErrStopMesheryCode, errors.Alert, []string{"Error stopping Meshery"}, []string{err.Error()}, []string{"Meshery server is not stopped, some of the docker containers are still running"}, []string{"Verify all docker containers of Meshery server are stopped"})
}

func ErrResetMeshconfig(err error) error {
	return errors.New(
		ErrResetMeshconfigCode,
		errors.Alert,
		[]string{"Error resetting meshconfig to default settings"},
		[]string{err.Error()},
		[]string{"Meshery server config file is not reset to default settings"},
		[]string{"Verify Meshery server config file is reset to default settings by executing `mesheryctl system context view`" + FormatErrorReference()})
}

func ErrApplyManifest(err error, deleteStatus, updateStatus bool) error {
	return errors.New(ErrApplyManifestCode, errors.Alert, []string{"Error applying manifest with update: ", strconv.FormatBool(updateStatus), " and delete: ", strconv.FormatBool(deleteStatus)}, []string{err.Error()}, []string{}, []string{})
}

func ErrApplyOperatorManifest(err error, deleteStatus, updateStatus bool) error {
	return errors.New(ErrApplyOperatorManifestCode, errors.Alert, []string{"Error applying operator manifests with update: ", strconv.FormatBool(updateStatus), " and delete: ", strconv.FormatBool(deleteStatus)}, []string{err.Error()}, []string{}, []string{})
}

func ErrCreateDir(err error, obj string) error {
	return errors.New(ErrCreateDirCode, errors.Alert, []string{"Error creating directory ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrUnmarshalDockerCompose(err error, obj string) error {
	return errors.New(
		ErrUnmarshalDockerComposeCode,
		errors.Alert,
		[]string{"Error processing JSON response from Meshery Server", obj},
		[]string{err.Error()},
		[]string{"Either the JSON response is invalid or the Response is distorted"},
		[]string{"Ensure Meshery Server is running and you have a strong newtwork connection"})
}

func ErrUnsupportedPlatform(platform string, config string) error {
	return errors.New(
		ErrUnsupportedPlatformCode,
		errors.Alert,
		[]string{"Unsupported platform"},
		[]string{"The provided platform is not supprted"},
		[]string{"The platform ", platform, " is not supported for the deployment of Meshery. "},
		[]string{"Supported platforms are:\n\n- docker\n- kubernetes\n\n Verify this setting in your meshconfig at ", config, " or verify by executing `mesheryctl system context view`"})
}

func ErrRetrievingCurrentContext(err error) error {
	return errors.New(
		ErrRetrievingCurrentContextCode,
		errors.Alert,
		[]string{"Error retrieving current-context"},
		[]string{err.Error()},
		[]string{"current-context is not retrieved successfully"},
		[]string{"Verify current-context is retrieved successfully and valid" + FormatErrorReference()})
}

func ErrSettingDefaultContextToConfig(err error) error {
	return errors.New(ErrSettingDefaultContextToConfigCode, errors.Alert, []string{"Error setting default context to config"}, []string{err.Error()}, []string{"Mesheryctl config file may not exist or is invalid"}, []string{"Make sure the Mesheryctl config file exists"})
}

func ErrSettingTemporaryContext(err error) error {
	return errors.New(ErrSettingTemporaryContextCode, errors.Alert, []string{"Error setting temporary context"}, []string{err.Error()}, []string{"temporary context is not set properly"}, []string{"Verify the temporary context is set properly using the -c flag provided"})
}

func ErrCreateManifestsFolder(err error) error {
	return errors.New(ErrCreateManifestsFolderCode, errors.Alert, []string{"Error creating manifest folder"}, []string{err.Error()}, []string{"system error in creating manifest folder"}, []string{"Make sure manifest folder (.meshery/manifests) is created properly"})
}

func ErrRestartMeshery(err error) error {
	return errors.New(ErrRestartMesheryCode, errors.Alert, []string{"Error restarting Meshery"}, []string{err.Error()}, []string{"Meshery is not running"}, []string{"Restart Meshery instance"})
}

func ErrK8sConfig(err error) error {
	return errors.New(ErrK8sConfigCode, errors.Alert, []string{"The Kubernetes cluster is not accessible."}, []string{err.Error(), "\nThe Kubernetes cluster is not accessible", " Please confirm that the cluster is running", " See https://docs.meshery.io/installation/quick-start for additional instructions."}, []string{"Kubernetes cluster isn't running or inaccessible"}, []string{"Verify kubernetes and Meshery connectivity or Verify kubeconfig certificates."})
}

func ErrK8SQuery(err error) error {
	return errors.New(ErrK8sQueryCode, errors.Alert, []string{"The Kubernetes cluster is not accessible."}, []string{err.Error(), "\nThe Kubernetes cluster is not accessible", " Please confirm that the token is valid", " See https://docs.meshery.io/installation/quick-start for additional instructions"}, []string{"Kubernetes cluster is unavailable and that the token is invalid"}, []string{"Please confirm that your cluster is available and that the token is valid. See https://docs.meshery.io/installation/quick-start for additional instructions"})
}

func ErrInitPortForward(err error) error {
	return errors.New(
		ErrInitPortForwardCode,
		errors.Alert,
		[]string{"Failed to initialize port-forward"},
		[]string{err.Error(), "Failed to create new Port Forward instance"},
		[]string{"Connection not established"},
		[]string{"Ensure you have a strong network connection"})
}

func ErrRunPortForward(err error) error {
	return errors.New(
		ErrRunPortForwardCode,
		errors.Fatal,
		[]string{"Failed to run port-forward"},
		[]string{err.Error(), "Error running port-forward for Meshery"},
		[]string{"Meshery pod is not in running phase", "mesheryctl can't connect to kubernetes with client-go"},
		[]string{"Make sure Meshery pod exists and is in running state",
			"Check if mesheryctl is connected to kubernetes with `mesheryctl system check`"},
	)
}

func ErrFailedGetEphemeralPort(err error) error {
	return errors.New(
		ErrFailedGetEphemeralPortCode,
		errors.Fatal,
		[]string{"Failed to get a free port"},
		[]string{err.Error(), "Failed to start port-forwarding"},
		[]string{"Unable to provide a free port to connect to a kuberentes cluster"},
		[]string{"Ensure your Meshery Server is running,", "Ensure mesheryctl is connected to kubernetes cluster with `mesheryctl system check`"},
	)
}

func ErrCreatingDockerClient(err error) error {
	return errors.New(
		ErrCreatingDockerClientCode,
		errors.Critical,
		[]string{"Failed to create Docker client"},
		[]string{"Error occurred while creating Docker client from config file", err.Error()},
		[]string{"Missing or invalid docker config"},
		[]string{"Please check the Docker config file for any errors or missing information. Make sure it is correctly formatted and contains all the required fields."})
}

func ErrContextContent() error {
	return errors.New(
		ErrContextContentCode,
		errors.Fatal,
		[]string{"Failed to detect context"},
		[]string{"Unable to detect current-context"},
		[]string{"Error while trying to fetch current-context in YML file"},
		[]string{"Ensure a valid context name is provided"})
}

func ErrWriteConfig(err error) error {
	return errors.New(
		ErrWriteConfigCode,
		errors.Fatal,
		[]string{"Error in writing config"},
		[]string{err.Error()},
		[]string{"Unable to write to config file"},
		[]string{"Ensure the right context is set. " + FormatErrorReference()})
}

func ErrSwitchChannelResponse() error {
	return errors.New(
		ErrSwitchChannelResponseCode,
		errors.Alert,
		[]string{"Unable to exectute command"},
		[]string{"Channel switch aborted"},
		[]string{"No user response provided"},
		[]string{"Provide a response or use the -y flag for confirmation. " + FormatErrorReference()})
}

func ErrGetCurrentContext(err error) error {
	return errors.New(
		ErrGetCurrentContextCode,
		errors.Fatal,
		[]string{"Unable to get current-context"},
		[]string{err.Error()},
		[]string{"Invalid context name provided"},
		[]string{"Ensure a valid context name is provided. " + contextdocs + "Also " + contextDir})
}

func ErrSetCurrentContext(err error) error {
	return errors.New(
		ErrSetCurrentContextCode,
		errors.Fatal,
		[]string{"Unable to set current-context"},
		[]string{err.Error()},
		[]string{"Invalid context name provided"},
		[]string{"Ensure a valid context name is provided. " + contextdocs + "Also " + contextDir})
}

func ErrTokenContext(err error) error {
	return errors.New(
		ErrTokenContextCode,
		errors.Fatal,
		[]string{"Unable to get token"},
		[]string{err.Error()},
		[]string{"No token found for the Current context"},
		[]string{"Ensure your Meshconfig file has valid token provided." + FormatErrorReference()})
}

func ErrProviderInfo(err error) error {
	return errors.New(
		ErrProviderInfoCode,
		errors.Fatal,
		[]string{"Unable to verify provider"},
		[]string{err.Error()},
		[]string{"Unable to verify provider  as Meshery server was unreachable"},
		[]string{"Start Meshery to verify provider. Run `mesheryctl system provider set [provider] --force` to force set the provider" + FormatErrorReference()})
}

func ErrValidProvider() error {
	return errors.New(
		ErrValidProviderCode,
		errors.Fatal,
		[]string{"Invalid provider"},
		[]string{"Unable to validate provider"},
		[]string{"A wrong provider was specified"},
		[]string{"Specify a valid provider" + FormatErrorReference()})
}

func ErrUnmarshallConfig(err error) error {
	return errors.New(
		ErrUnmarshallConfigCode,
		errors.Fatal,
		[]string{"Invalid config"},
		[]string{err.Error()},
		[]string{"Unable to decode Meshconfig."},
		[]string{"Ensure you have the right configuration set in your Meshconfig file." + FormatErrorReference()})
}

func ErrUploadFileParams(err error) error {
	return errors.New(
		ErrUploadFileParamsCode,
		errors.Fatal,
		[]string{"Unable to upload"},
		[]string{err.Error()},
		[]string{"Unable to upload parameters from config file with provided context"},
		[]string{"Ensure you have a strong network connection and the right configuration set in your Meshconfig file." + FormatErrorReference()})

}
