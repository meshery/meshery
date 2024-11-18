package registry

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

var (
	ErrUpdateModelCode    = "mesheryctl-1057"
	ErrUpdateRegistryCode = "mesheryctl-1059"
	ErrParsingSheetCode   = "mesheryctl-1128"
)

func ErrUpdateRegistry(err error, path string) error {
	return errors.New(ErrUpdateRegistryCode, errors.Alert, []string{"error updating registry at ", path}, []string{err.Error()}, []string{"Provided spreadsheet ID is incorrect", "Provided credential is incorrect"}, []string{"Ensure correct spreadsheet ID is provided", "Ensure correct credential is used"})
}
func ErrUpdateModel(err error, modelName string) error {
	return errors.New(ErrUpdateModelCode, errors.Alert, []string{"error updating model ", modelName}, []string{err.Error()}, []string{"Model does not exist"}, []string{"Ensure existence of model, check for typo in model name"})
}
func ErrParsingSheet(err error, obj string) error {
	return errors.New(ErrParsingSheetCode, errors.Alert, []string{fmt.Sprintf("error parsing %s sheet", obj)}, []string{fmt.Sprintf("while parsing the %s sheet encountered an error: %s", obj, err)}, []string{"provied sheet id for %s might be incorrect"}, []string{"ensure the sheet id is correct"})
}
