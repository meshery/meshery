package registry

import (
	"fmt"
	"time"

	"github.com/meshery/meshkit/errors"
)

var (
	ErrGenerateModelCode         = "mesheryctl-1055"
	ErrGenerateComponentCode     = "mesheryctl-1056"
	ErrUpdateModelCode           = "mesheryctl-1057"
	ErrUpdateComponentCode       = "mesheryctl-1058"
	ErrUpdateRegistryCode        = "mesheryctl-1059"
	ErrParsingSheetCode          = "mesheryctl-1128"
	ErrGenerationTimeoutCode     = "mesheryctl-1129"
	ErrParsingRelationshipCode   = "mesheryctl-1130"
	ErrModelGenerationFailedCode = "mesheryctl-1131"
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

func ErrGenerationTimeout(duration time.Duration) error {
	return errors.New(ErrGenerationTimeoutCode, errors.Alert, []string{fmt.Sprintf("model generation timed out after %v", duration)}, []string{fmt.Sprintf("The model generation process exceeded the maximum allowed time of %v", duration)}, []string{"The model source may be unresponsive", "Network connectivity issues", "Large number of components to generate"}, []string{"Try generating a specific model using --model flag", "Check network connectivity", "Review the generation logs for problematic models"})
}

func ErrParsingRelationship(err error, modelName string) error {
	return errors.New(ErrParsingRelationshipCode, errors.Alert, []string{fmt.Sprintf("error parsing relationships for model: %s", modelName)}, []string{err.Error()}, []string{"Invalid relationship definition format", "Missing required relationship fields", "Malformed selector in relationship"}, []string{"Verify the relationship CSV format", "Check that all required fields are present", "Validate the selector JSON syntax"})
}

func ErrModelGenerationFailed(modelName string, err error) error {
	return errors.New(ErrModelGenerationFailedCode, errors.Alert, []string{fmt.Sprintf("failed to generate model: %s", modelName)}, []string{err.Error()}, []string{"Invalid model source URL", "Unsupported registrant type", "Network issues fetching model data"}, []string{"Verify the model's source URL is accessible", "Check the registrant type is supported", "Review network connectivity"})
}
