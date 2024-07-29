package registration

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrDirPkgUnitParseFailCode = "replace_me"
	ErrGetEntityCode = "replace_me"
	ErrRegisterEntityCode = "replace_me"
	ErrInvalidModelDefinitionCode = "replace_me"
	ErrImportFailureCode = "replace_me"
)


func ErrRegisterEntity(err error, name, entity string) error {
	return errors.New(ErrRegisterEntityCode, errors.Alert, []string{fmt.Sprintf("failed to register %s %s", name, entity)}, []string{err.Error()}, []string{fmt.Sprintf("%s definition violates the definition schema", entity), fmt.Sprintf("%s might be missing model details", entity)}, []string{fmt.Sprintf("ensure the %s definition follows the correct schema", entity), fmt.Sprintf("ensure %s definition belongs to correct model", entity)})
}

func ErrGetEntity(err error) error {
	return errors.New(ErrGetEntityCode, errors.Alert, []string{err.Error()}, []string{err.Error()}, []string{}, []string{})
}

func ErrDirPkgUnitParseFail(dirpath string, err error) error {
	return errors.New(ErrDirPkgUnitParseFailCode, errors.Alert, []string{fmt.Sprintf("Directory %s is not a valid packaging unit for registration", dirpath)}, []string{fmt.Sprintf("Could not parse the given directory %s into a packaging unit:  %s", dirpath, err.Error())}, []string{}, []string{} )
}

func ErrInvalidModelDefinition(path string, err error) error {
	return errors.New(ErrInvalidModelDefinitionCode, errors.Alert, []string{fmt.Sprintf("The file at path %s is not a valid Model definition", path)}, []string{fmt.Sprintf("Could not parse Model:  %s", err.Error())}, []string{}, []string{} )
}

func ErrImportFailure(hostname string, failedMsg string) error {
	return errors.New(ErrImportFailureCode, errors.Alert, []string{}, []string{fmt.Sprintf("The import process for a registrant %s encountered difficulties,due to which %s. Specific issues during the import process resulted in certain entities not being successfully registered in the table.", hostname, failedMsg)}, []string{}, []string{fmt.Sprintf("Visit docs with the error code %s", "https://docs.meshery.io/reference/error-codes")})
}
