package registry

import (
	goerrors "errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
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
	ErrModelGenerationFailedCode = "mesheryctl-1159"
	ErrInvalidOutputPathCode     = "mesheryctl-1167"
	ErrCreateOutputDirCode       = "mesheryctl-1168"
	ErrFailedToGetCWDCode        = "mesheryctl-1169"
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
	return errors.New(ErrParsingSheetCode, errors.Alert, []string{fmt.Sprintf("error parsing %s sheet", obj)}, []string{fmt.Sprintf("while parsing the %s sheet encountered an error: %s", obj, err)}, []string{"provided sheet id for %s might be incorrect"}, []string{"ensure the sheet id is correct"})
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

func ErrCreateOutputDir(err error, path string) error {
	return errors.New(ErrCreateOutputDirCode, errors.Alert,
		[]string{fmt.Sprintf("Failed to create output directory: %s", path)},
		[]string{err.Error()},
		[]string{
			"Insufficient filesystem permissions to create the directory.",
			"The parent directory does not exist or is not writable.",
			"The path points to an existing file, not a directory.",
		},
		[]string{
			fmt.Sprintf("Ensure you have write permissions for the parent directory: %s", filepath.Dir(path)),
			"Try specifying a different output location using --output.",
			"Check if a file with the same name already exists.",
		})
}

func ErrFailedToGetCWD(err error) error {
	return errors.New(
		ErrFailedToGetCWDCode,
		errors.Alert,
		[]string{"Failed to get current working directory"},
		[]string{fmt.Sprintf("Unable to determine the current working directory: %v", err)},
		[]string{
			"Insufficient permissions to access the directory",
			"The current directory has been deleted or moved",
			"File system errors or corruption",
		},
		[]string{
			"Ensure you have read permissions for the current directory",
			"Navigate to a valid directory and try again",
			"Check file system health and integrity",
		})
}

func ErrNoSourceSpecified() error {
	return goerrors.New(utils.RegistryError(
		"[ Spreadsheet ID | Registrant Connection Definition Path | Local Directory | Individual CSV files ] isn't specified\n\n"+
			"Usage: \n"+
			"mesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED\n"+
			"mesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED --model \"[model-name]\"\n"+
			"mesheryctl registry generate --model-csv [path] --component-csv [path] --relationship-csv [path]\n"+
			"Run 'mesheryctl registry generate --help' to see detailed help message",
		"generate"))
}

func ErrSpreadsheetCredRequired() error {
	return goerrors.New(utils.RegistryError(
		"Spreadsheet Credentials is required\n\n"+
			"Usage: \n"+
			"mesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED\n"+
			"mesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED --model \"[model-name]\"\n"+
			"Run 'mesheryctl registry generate --help'",
		"generate"))
}

func ErrRegistrantCredRequired() error {
	return goerrors.New(utils.RegistryError(
		"Registrant Credentials is required\n\n"+
			"Usage: mesheryctl registry generate --registrant-def [path to connection definition] --registrant-cred [path to credential definition]\n"+
			"Run 'mesheryctl registry generate --help'",
		"generate"))
}
