package model

import (
	"fmt"

	goerrors "errors"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrExportModelCode                      = "mesheryctl-1127"
	ErrTemplateFileNotPresentCode           = "mesheryctl-1131"
	ErrModelUnsupportedOutputFormatCode     = "mesheryctl-1146"
	ErrModelInitCode                        = "mesheryctl-1148"
	ErrModelUnsupportedVersionCode          = "mesheryctl-1149"
	ErrModelBuildCode                       = "mesheryctl-1151"
	ErrModelUnsupportedOutputTypeFormatCode = "mesheryctl-1155"
	ErrModelInvalidPageNumberCode           = "mesheryctl-1156"
	ErrFolderNotPresentCode                 = "mesheryctl-1157"
	ErrProcessingConfigCode                 = "mesheryctl-1158"
)

func ErrExportModel(err error, name string) error {
	return errors.New(ErrExportModelCode, errors.Fatal, []string{"Error exporting model"}, []string{fmt.Sprintf("Given model with name: %s could not be exported: %s", name, err)}, []string{"Model may not be present in the registry"}, []string{"Ensure that there are no typos in the model name"})
}

func ErrProcessingConfig(message string) error {
	return errors.New(ErrProcessingConfigCode, errors.Fatal, []string{"Error processing config"}, []string{fmt.Sprintf("Error processing config i.e %s", message)}, []string{"Check config again"}, []string{"Are you sure the configuration is correct, for surety check again form it correct"})
}

func ErrTemplateFileNotPresent() error {
	return errors.New(ErrTemplateFileNotPresentCode, errors.Fatal, []string{"error no template file provided"}, []string{"no template file is provided while using url for importing a model "}, []string{"template file not present"}, []string{"ensure that the template file is present in the given path"})
}
func ErrFolderNotPresent(message string) error {
	return errors.New(ErrFolderNotPresentCode, errors.Fatal, []string{"error folder not provided"}, []string{message}, []string{"folder not present"}, []string{"ensure that the folder is present in the given path"})
}
func ErrModelUnsupportedOutputFormat(message string) error {
	return errors.New(ErrModelUnsupportedOutputFormatCode, errors.Fatal, []string{"Error viewing model"}, []string{message}, []string{"Output format not supported"}, []string{"Ensure giving a valid format"})
}

func ErrModelUnsupportedOutputType(message string) error {
	return errors.New(ErrModelUnsupportedOutputTypeFormatCode, errors.Fatal, []string{"Error in output type"}, []string{message}, []string{"Output type not supported"}, []string{"Ensure giving a valid format"})
}

func ErrModelUnsupportedVersion(message string) error {
	return errors.New(ErrModelUnsupportedVersionCode, errors.Fatal, []string{"Error in model version format"}, []string{message}, []string{"Version format not supported"}, []string{"Ensure giving a semver version format"})
}
func ErrModelInvalidPageNumber(message string) error {
	return errors.New(ErrModelInvalidPageNumberCode, errors.Fatal, []string{"Error in page number"}, []string{message}, []string{"Page number format not supported"}, []string{"Ensure giving a valid page number ( i.e > 0 )"})
}
func ErrModelInitFromString(message string) error {
	return errors.New(ErrModelInitCode, errors.Fatal,
		[]string{"Error model init"},
		[]string{message},
		[]string{"Error during run of model init command"},
		[]string{"Ensure passing all params according to the command description"})
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
	return errors.New(ErrModelBuildCode, errors.Fatal,
		[]string{"Error model build"},
		[]string{err.Error()},
		[]string{"Error during run of model build command"},
		[]string{"Ensure passing all params according to the command description"})
}
func ErrNilRequest(modelName string) error {
	return errors.New(ErrExportModelCode, errors.Fatal,
		[]string{"Request creation failed"},
		[]string{fmt.Sprintf("Request returned nil for model %s", modelName)},
		[]string{"HTTP client returned nil request"},
		[]string{"Check the request building logic and ensure valid URL"})
}

func ErrNilResponse(modelName string) error {
	return errors.New(ErrExportModelCode, errors.Fatal,
		[]string{"Response handling failed"},
		[]string{fmt.Sprintf("Response returned nil for model %s", modelName)},
		[]string{"HTTP client returned nil response"},
		[]string{"Check the server availability and request execution"})
}

func ErrInvalidStatus(modelName string, statusCode int, statusText string) error {
	return errors.New(ErrExportModelCode, errors.Fatal,
		[]string{"Invalid response status"},
		[]string{fmt.Sprintf("Failed exporting model %s: status %d %s", modelName, statusCode, statusText)},
		[]string{"Server returned non-200 status"},
		[]string{"Verify the API endpoint and parameters"})
}
