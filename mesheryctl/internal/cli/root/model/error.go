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

	errInitOneArg       = "model init requires one argument: 'Model name' using kebab-cased (eg. model-name, aws-model-name)"
	errInitFolderExists = "folder %s exists, please specify different model name or version"

	errSearchUsage     = "Usage: mesheryctl model search [query-text]\nRun 'mesheryctl model search --help' to see detailed help message"
	errSearchModelName = "Please provide a model name. " + errSearchUsage

	errInvalidArg = "only one argument must be provided and needs to be enclosed by double quotes if it contains spaces (eg. \"model name\", modelName)"

	errDeleteInvalidArg = "[ model-id | model-name ] is required\n\nUsage: mesheryctl model delete [model-id | model-name]\nRun 'mesheryctl model delete --help' to see detailed help message"

	errGenerateUsageMsg = "Usage: mesheryctl model generate [ file | filePath | URL ] path\nRun 'mesheryctl model generate --help' to see detailed help message"
	errImportUsageMsg   = "Usage: mesheryctl model import [ file | filePath | URL ]\nRun 'mesheryctl model import --help' to see detailed help message"

	errGenerateMissingArgsMsg = "either --file flag, a URL or a path as argument must be specified\n\n%s"
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

func ErrDeleteModel(err error, nameOrID string) error {
	return errors.New(ErrDeleteModelCode, errors.Alert,
		[]string{"Failed to delete model"},
		[]string{fmt.Sprintf("Failed to delete model with name or ID '%s': %s", nameOrID, err.Error())},
		[]string{"The specified model name or ID may not exist"},
		[]string{"Verify the model name or ID using 'mesheryctl model list' and try again"})
}
