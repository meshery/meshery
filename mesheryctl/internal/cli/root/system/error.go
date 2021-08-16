package system

import (
	"strconv"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrHealthCheckFailedCode        = "1000"
	ErrInvalidAdapterCode           = "1001"
	ErrDownloadFileCode             = "1002"
	ErrStopMesheryCode              = "1003"
	ErrResetMeshconfigCode          = "1004"
	ErrApplyManifestCode            = "1005"
	ErrApplyOperatorManifestCode    = "1006"
	ErrCreateDirCode                = "1007"
	ErrUnmarshalCode                = "1008"
	ErrUnsupportedPlatformCode      = "1009"
	ErrRetrievingCurrentContextCode = "replace me"
	ErrSettingTemporaryContextCode  = "replace me"
	ErrCreateManifestsFolderCode    = "replace me"
	ErrProcessingMctlConfigCode     = "replace me"
)

func ErrHealthCheckFailed(err error) error {
	return errors.New(ErrHealthCheckFailedCode, errors.Alert, []string{"Health checks failed"}, []string{err.Error()}, []string{}, []string{})
}

func ErrInvalidAdapter(err error, obj string) error {
	return errors.New(ErrInvalidAdapterCode, errors.Alert, []string{"Invalid adapter ", obj, " specified"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDownloadFile(err error, obj string) error {
	return errors.New(ErrDownloadFileCode, errors.Alert, []string{"Error downloading file ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrStopMeshery(err error) error {
	return errors.New(ErrStopMesheryCode, errors.Alert, []string{"Error stopping meshery"}, []string{err.Error()}, []string{}, []string{})
}

func ErrResetMeshconfig(err error) error {
	return errors.New(ErrResetMeshconfigCode, errors.Alert, []string{"Error resetting meshconfig"}, []string{err.Error()}, []string{}, []string{})
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
	return errors.New(ErrUnsupportedPlatformCode, errors.Alert, []string{"the platform ", platform, " is not supported. Supported platforms are:\n\n- docker\n- kubernetes\n\nVerify this setting in your meshconfig at ", config, " or verify by executing `mesheryctl system context view`"}, []string{err.Error()}, []string{}, []string{})
}

func ErrRetrievingCurrentContext(err error) error {
	return errors.New(ErrRetrievingCurrentContextCode, errors.Alert, []string{"Error retrieving current context"}, []string{err.Error()}, []string{}, []string{})
}

func ErrSettingTemporaryContext(err error) error {
	return errors.New(ErrSettingTemporaryContextCode, errors.Alert, []string{"Error setting temporary context"}, []string{err.Error()}, []string{}, []string{})
}

func ErrCreateManifestsFolder(err error) error {
	return errors.New(ErrCreateManifestsFolderCode, errors.Alert, []string{"Error creating manifest folder"}, []string{err.Error()}, []string{}, []string{})
}

func ErrProcessingMctlConfig(err error) error {
	return errors.New(ErrProcessingMctlConfigCode, errors.Alert, []string{"Error processing config"}, []string{err.Error()}, []string{}, []string{})
}
