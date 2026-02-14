package model

import (
	"fmt"

	goerrors "errors"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrExportModelCode                  = "mesheryctl-1127"
	ErrTemplateFileNotPresentCode       = "mesheryctl-1131"
	ErrModelUnsupportedOutputFormatCode = "mesheryctl-1146"
	ErrModelInitCode                    = "mesheryctl-1148"
	ErrModelUnsupportedVersionCode      = "mesheryctl-1149"
	ErrModelBuildCode                   = "mesheryctl-1151"
	ErrDeleteModelCode                  = "mesheryctl-1200"

	// Error Constants
	errBuildUsage                    = "Usage:\nmesheryctl model build [model-name]\nor\nmesheryctl model build [model-name]/[model-version]\n\nRun 'mesheryctl model build --help' to see detailed help message"
	errBuildFolderNotFound           = "\nfolder %s does not exist"
	errBuildMultiVersionNotSupported = "\nCommand does not support multiple versions build under one image"

	errInitOneArg            = "must provide only one argument: model name"
	errInitUnsupportedFormat = "[ %s ] are the only format supported"
	errInitInvalidVersion    = "version must follow a semver format, f.e. v1.2.3"
	errInitFolderExists      = "folder %s exists, please specify different model name or version"
	errInitInvalidModelName  = "invalid model name: name must match pattern ^[a-z0-9-]+$"

	errSearchUsage     = "Usage: mesheryctl model search [query-text]\nRun 'mesheryctl model search --help' to see detailed help message"
	errSearchModelName = "Please provide a model name. " + errSearchUsage
)

func ErrExportModel(err error, name string) error {
	return errors.New(ErrExportModelCode, errors.Fatal, []string{"Error exporting model"}, []string{fmt.Sprintf("Given model with name: %s could not be exported: %s", name, err)}, []string{"Model may not be present in the registry"}, []string{"Ensure that there are no typos in the model name"})
}

func ErrTemplateFileNotPresent() error {
	return errors.New(ErrTemplateFileNotPresentCode, errors.Fatal, []string{"error no template file provided"}, []string{"no template file is provided while using url for importing a model "}, []string{"template file not present"}, []string{"ensure that the template file is present in the given path"})
}

func ErrModelUnsupportedOutputFormat(message string) error {
	return errors.New(ErrModelUnsupportedOutputFormatCode, errors.Fatal, []string{"Error viewing model"}, []string{message}, []string{"Output format not supported"}, []string{"Ensure giving a valid format"})
}

func ErrModelUnsupportedVersion(message string) error {
	return errors.New(ErrModelUnsupportedVersionCode, errors.Fatal, []string{"Error in model version format"}, []string{message}, []string{"Version format not supported"}, []string{"Ensure giving a semver version format"})
}

func ErrModelInitFromString(message string) error {
	return errors.New(ErrModelInitCode, errors.Fatal, []string{"Error model init"}, []string{message}, []string{"Error during run of model init command"}, []string{"Ensure passing all params according to the command description"})
}

func ErrModelInit(err error) error {
	return ErrModelInitFromString(err.Error())
}

func ErrModelBuildFromStrings(message ...string) error {
	errs := make([]error, 0, len(message))
	for _, m := range message {
		errs = append(errs, goerrors.New(m))
	}
	return ErrModelBuild(goerrors.Join(errs...))
}

func ErrModelBuild(err error) error {
	return errors.New(ErrModelBuildCode, errors.Fatal, []string{"Error model build"}, []string{err.Error()}, []string{"Error during run of model build command"}, []string{"Ensure passing all params according to the command description"})
}

func ErrDeleteModel(err error, id string) error {
	return errors.New(ErrDeleteModelCode, errors.Alert, []string{"Failed to delete model"}, []string{fmt.Sprintf("failed to delete model with ID %s: %s", id, err.Error())}, []string{"The specified model ID may not exist"}, []string{"Verify the model ID using `mesheryctl model list`"})
}
