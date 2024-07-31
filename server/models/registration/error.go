package registration

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrDirPkgUnitParseFailCode = "replace_me"
	ErrGetEntityCode = "replace_me"
	ErrRegisterEntityCode = "replace_me"
	ErrImportFailureCode = "replace_me"
	ErrMissingRegistrantCode = "replace_me"
	ErrSeedingComponentsCode = "replace-me"
)



func ErrSeedingComponents(err error) error {
	return errors.New(
		ErrSeedingComponentsCode,
		errors.Alert,
		[]string{"Failed to register the given models into meshery's registry"},
		[]string{err.Error()},
		[]string{"Given models may not be in accordance with Meshery's schema", "Internal(OS level) error while reading files" },
		[]string{"Make sure the models being seeded are valid in accordance with Meshery's schema", "If it is an internal error, please try again after some time"},
	)
}

func ErrMissingRegistrant(modelName string) error {
	return errors.New(
		ErrMissingRegistrantCode,
		errors.Alert,
		[]string{fmt.Sprintf("Model with name: %s does not have registrant information", modelName)},
		[]string{"Meshery models are always registered in context of a registrant."},
		// there is only one cause for this error
		[]string{""},
		[]string{"Make sure that the registrant information is present in the model definition"},
	)
}

func ErrRegisterEntity(err error, name, entity string) error {
	return errors.New(
		ErrRegisterEntityCode,
		errors.Alert,
		[]string{fmt.Sprintf("Failed to register an entity of name: %s and type: %s into Meshery's registry", name, entity)},
		[]string{err.Error()},
		[]string{fmt.Sprintf("%s definition might be violating the definition schema", entity), fmt.Sprintf("%s might be missing model details", entity)},
		[]string{fmt.Sprintf("ensure the %s definition follows the correct schema", entity), fmt.Sprintf("ensure %s definition belongs to correct model", entity)},
	)
}

func ErrGetEntity(err error) error {
	return errors.New(
		ErrGetEntityCode,
		errors.Alert,
		[]string{"Could not parse the given data into any Meshery entity"},
		[]string{err.Error()},
		[]string{"Entity definition might not be in accordance with it's corresponding schema", "Might be an invalid/unsupported schemaVersion"},
		[]string{"Ensure that the definition given is in accordance with it's schema", "Ensure that the schemaVersion used is valid and is supported by Meshery"},
	)
}

func ErrDirPkgUnitParseFail(dirpath string, err error) error {
	return errors.New(
		ErrDirPkgUnitParseFailCode,
		errors.Alert,
		[]string{fmt.Sprintf("Directory at path: %s cannot be registered into Meshery", dirpath)},
		[]string{fmt.Sprintf(err.Error())},
		[]string{"The directory might not have a valid model definition", "Might be some internal issues while walking the file tree"},
		[]string{"Make sure that there is a valid model definition present in the directory. Meshery's registration pipeline currently does not support nested models, therefore the behaviour might be unexpected if it contains nested models.", "If there is an internal error, please try again after some time"},
	)
}

func ErrImportFailure(hostname string, failedMsg string) error {
	return errors.New(
		ErrImportFailureCode,
		errors.Alert,
		[]string{fmt.Sprintf("Errors while registering entities for registrant: %s", hostname)},
		[]string{failedMsg},
		[]string{"Entity definition might not be valid in accordance with schema", "Entity version might not be supported by Meshery"},
		[]string{fmt.Sprintf("Visit docs with the error code %s", "https://docs.meshery.io/reference/error-codes"),
	})
}