package registry

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

var (
	ErrGenerateModelCode        = "mesheryctl-1055"
	ErrGenerateComponentCode    = "mesheryctl-1056"
	ErrUpdateModelCode          = "mesheryctl-1057"
	ErrUpdateComponentCode      = "mesheryctl-1058"
	ErrUpdateRegistryCode       = "mesheryctl-1059"
	ErrParsingSheetCode         = "mesheryctl-1128"
	ErrFindImageRefRequiredCode = "mesheryctl-1129"
	ErrFindInvalidOutputCode    = "mesheryctl-1130"
	ErrFindScanImageCode        = "mesheryctl-1131"
)

func ErrUpdateRegistry(err error, path string) error {
	return errors.New(ErrUpdateRegistryCode, errors.Alert, []string{"error updating registry at ", path}, []string{err.Error()}, []string{"Provided spreadsheet ID is incorrect", "Provided credential is incorrect"}, []string{"Ensure correct spreadsheet ID is provided", "Ensure correct credential is used"})
}

func ErrGenerateModel(err error, modelName string) error {
	return errors.New(ErrGenerateModelCode, errors.Alert, []string{fmt.Sprintf("error generating model: %s", modelName)}, []string{fmt.Sprintf("Error generating model: %s\n %s", modelName, err.Error())}, []string{"Registrant used for the model is not supported", "Verify the model's source URL.", "Failed to create a local directory in the filesystem for this model."}, []string{"Ensure that each kind of registrant used is a supported kind.", "Ensure correct model source URL is provided and properly formatted.", "Ensure sufficient permissions to allow creation of model directory."})
}

func ErrGenerateComponent(err error, modelName, compName string) error {
	return errors.New(ErrGenerateComponentCode, errors.Alert, []string{"error generating comp %s of model %s", compName, modelName}, []string{err.Error()}, []string{}, []string{})
}

func ErrUpdateModel(err error, modelName string) error {
	return errors.New(ErrUpdateModelCode, errors.Alert, []string{"error updating model ", modelName}, []string{err.Error()}, []string{"Model does not exist"}, []string{"Ensure existence of model, check for typo in model name"})
}

func ErrUpdateComponent(err error, modelName, compName string) error {
	return errors.New(ErrUpdateComponentCode, errors.Alert, []string{fmt.Sprintf("error updating component %s of model %s ", compName, modelName)}, []string{err.Error()}, []string{"Component does not exist", "Component definition is corrupted"}, []string{"Ensure existence of component, check for typo in component name", "Regenerate corrupted component"})
}
func ErrParsingSheet(err error, obj string) error {
	return errors.New(ErrParsingSheetCode, errors.Alert, []string{fmt.Sprintf("error parsing %s sheet", obj)}, []string{fmt.Sprintf("while parsing the %s sheet encountered an error: %s", obj, err)}, []string{"provied sheet id for %s might be incorrect"}, []string{"ensure the sheet id is correct"})
}

func ErrFindImageRefRequired() error {
	return errors.New(ErrFindImageRefRequiredCode, errors.Alert,
		[]string{"OCI image reference is required"},
		[]string{"No OCI image reference was provided to the find command"},
		[]string{"The find command requires an OCI image reference to scan"},
		[]string{"Provide an OCI image reference, e.g., 'mesheryctl registry find docker.io/crossplane/crossplane:v1.14.0'"})
}

func ErrFindInvalidOutputFormat(format string) error {
	return errors.New(ErrFindInvalidOutputCode, errors.Alert,
		[]string{fmt.Sprintf("Invalid output format: %s", format)},
		[]string{fmt.Sprintf("The specified output format '%s' is not supported", format)},
		[]string{"Only 'json' and 'yaml' output formats are supported"},
		[]string{"Use '-o json' or '-o yaml' for structured output, or omit the flag for table output"})
}

func ErrFindScanImage(err error, imageRef string) error {
	return errors.New(ErrFindScanImageCode, errors.Alert,
		[]string{fmt.Sprintf("Failed to scan image: %s", imageRef)},
		[]string{err.Error()},
		[]string{"The image reference may be invalid or inaccessible", "Authentication may be required to access the registry", "Network connectivity issues"},
		[]string{"Verify the image reference is correct", "Ensure you have access to the container registry", "Check your network connection"})
}
