package relationships

import (
	"github.com/meshery/meshkit/errors"
)

var (
	ErrEmptySheetDataCode = "mesheryctl-1203"
)

var (
	errInvalidArg = "only one argument must be provided and needs to be enclosed by double quotes if it contains spaces (eg. \"model name\", modelName)"
	errMsg        = "both of [--spreadsheet-id, --spreadsheet-cred] is required\n\nUsage: mesheryctl exp relationship generate [--spreadsheet-id <sheet-ID>] [--spreadsheet-cred <sheet-CRED>]\nRun 'mesheryctl exp relationship generate --help'"
)

func ErrEmptySheetData(err error) error {
	return errors.New(
		ErrEmptySheetDataCode,
		errors.Alert,
		[]string{"Invalid spreadsheet"},
		[]string{err.Error()},
		[]string{"Ensure the spreadsheet contains valid relationship entries before importing"},
		[]string{"Spreadsheet must contains two headers and values for it's headers"},
	)
}
