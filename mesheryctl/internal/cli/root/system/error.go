package system

import (
	"strconv"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrHealthCheckFailedCode             = "1000"
	ErrInvalidComponentCode              = "1001"
	ErrDownloadFileCode                  = "1002"
	ErrStopMesheryCode                   = "1003"
	ErrResetMeshconfigCode               = "1004"
	ErrApplyManifestCode                 = "1005"
	ErrApplyOperatorManifestCode         = "1006"
	ErrCreateDirCode                     = "1007"
	ErrUnmarshalCode                     = "1008"
	ErrUnsupportedPlatformCode           = "1009"
	ErrRetrievingCurrentContextCode      = "1022"
	ErrSettingDefaultContextToConfigCode = "1059"
	ErrSettingTemporaryContextCode       = "1023"
	ErrCreateManifestsFolderCode         = "1024"
	ErrProcessingMctlConfigCode          = "1025"
	ErrRestartMesheryCode                = "1026"
	ErrK8sQueryCode                      = "1041"
	ErrK8sConfigCode                     = "1042"
	ErrInitPortForwardCode               = "1047"
	ErrRunPortForwardCode                = "1048"
	ErrFailedGetEphemeralPortCode        = "1049"
)

func ErrHealthCheckFailed(err error) error {
	return errors.New(ErrHealthCheckFailedCode, errors.Alert, []string{"Health checks failed"}, []string{err.Error()}, []string{"Health checks execution failed"}, []string{"Health checks execution should passed to start Meshery server successfully"})
}

func ErrInvalidComponent(err error, obj string) error {
	return errors.New(ErrInvalidComponentCode, errors.Alert, []string{"Invalid component ", obj, " specified"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDownloadFile(err error, obj string) error {
	return errors.New(ErrDownloadFileCode, errors.Alert, []string{"Error downloading file ", obj}, []string{err.Error()}, []string{"Failed to download docker-compose or manifest file due to system/config/network issues"}, []string{"Make sure docker-compose or manifest file is downloaded successfully"})
}

func ErrStopMeshery(err error) error {
	return errors.New(ErrStopMesheryCode, errors.Alert, []string{"Error stopping Meshery"}, []string{err.Error()}, []string{"Meshery server is not stopped, some of the docker containers are still running"}, []string{"Verify all docker containers of Meshery server are stopped"})
}

func ErrResetMeshconfig(err error) error {
	return errors.New(ErrResetMeshconfigCode, errors.Alert, []string{"Error resetting meshconfig"}, []string{err.Error()}, []string{"Meshery server config file is not reset to default settings"}, []string{"Verify Meshery server config file is reset to default settings by executing `mesheryctl system context view`"})
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

func ErrUnmarshal(err error, obj string) error {
	return errors.New(ErrUnmarshalCode, errors.Alert, []string{"Error unmarshalling file ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrUnsupportedPlatform(platform string, config string) error {
	return errors.New(ErrUnsupportedPlatformCode, errors.Alert, []string{}, []string{"the platform ", platform, " is not supported. Supported platforms are:\n\n- docker\n- kubernetes\n\nVerify this setting in your meshconfig at ", config, " or verify by executing `mesheryctl system context view`"}, []string{}, []string{})
}

func ErrRetrievingCurrentContext(err error) error {
	return errors.New(ErrRetrievingCurrentContextCode, errors.Alert, []string{"Error retrieving current context"}, []string{err.Error()}, []string{"current context is not retrieved successfully"}, []string{"Verify current context is retrieved successfully and valid"})
}

func ErrSettingDefaultContextToConfig(err error) error {
	return errors.New(ErrRetrievingCurrentContextCode, errors.Alert, []string{"Error setting default context to config"}, []string{err.Error()}, []string{"Mesheryctl config file may not exist or is invalid"}, []string{"Make sure the Mesheryctl config file exists"})
}

func ErrSettingTemporaryContext(err error) error {
	return errors.New(ErrSettingTemporaryContextCode, errors.Alert, []string{"Error setting temporary context"}, []string{err.Error()}, []string{"temporary context is not set properly"}, []string{"Verify the temporary context is set properly using the -c flag provided"})
}

func ErrCreateManifestsFolder(err error) error {
	return errors.New(ErrCreateManifestsFolderCode, errors.Alert, []string{"Error creating manifest folder"}, []string{err.Error()}, []string{"system error in creating manifest folder"}, []string{"Make sure manifest folder (.meshery/manifests) is created properly"})
}

func ErrProcessingMctlConfig(err error) error {
	return errors.New(ErrProcessingMctlConfigCode, errors.Alert, []string{"Error processing config"}, []string{err.Error()}, []string{"Error due to invalid format of Mesheryctl config"}, []string{"Make sure Mesheryctl config is in correct format and valid"})
}

func ErrRestartMeshery(err error) error {
	return errors.New(ErrRestartMesheryCode, errors.Alert, []string{"Error restarting Meshery"}, []string{err.Error()}, []string{"Meshery is not running"}, []string{"Restart Meshery instance"})
}

func ErrK8SQuery(err error) error {
	return errors.New(ErrK8sQueryCode, errors.Alert, []string{"The Kubernetes cluster is not accessible."}, []string{err.Error(), " The Kubernetes cluster is not accessible", " Please confirm that the cluster is running", " See https://docs.meshery.io/installation/quick-start for additional instructions"}, []string{"Kubernetes cluster isn't running or inaccessible"}, []string{"Verify kubernetes and Meshery connectivity or Verify kubeconfig certificates"})
}

func ErrK8sConfig(err error) error {
	return errors.New(ErrK8sConfigCode, errors.Alert, []string{"The Kubernetes cluster is not accessible."}, []string{err.Error(), "<br />The Kubernetes cluster is not accessible", " Please confirm that the token is valid", " See https://docs.meshery.io/installation/quick-start for additional instructions"}, []string{"Kubernetes cluster is unavailable and that the token is invalid"}, []string{"Please confirm that your cluster is available and that the token is valid. See https://docs.meshery.io/installation/quick-start for additional instructions"})
}

func ErrInitPortForward(err error) error {
	return errors.New(
		ErrInitPortForwardCode,
		errors.Alert, []string{"Failed to initialize port-forward"},
		[]string{err.Error(), "Failed to create new Port Forward instance"},
		nil, nil,
	)
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
		nil, nil,
	)
}
