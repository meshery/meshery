package relationships

import (
	"github.com/meshery/meshkit/errors"
)

var (
	ErrSpreadsheetFlagMissingCode = "mesheryctl-1203"
	ErrEmptySheetDataCode         = "mesheryctl-1204"
)

var (
	errInvalidArg = "only one argument must be provided and needs to be enclosed by double quotes if it contains spaces (eg. \"model name\", modelName)"
	errMsg        = "both of [--spreadsheet-id, --spreadsheet-cred] is required\n\nUsage: mesheryctl exp relationship generate [--spreadsheet-id <sheet-ID>] [--spreadsheet-cred <sheet-CRED>]\nRun 'mesheryctl exp relationship generate --help'"
)

func ErrSpreadsheetFlagMissing(err error) error {
	return errors.New(
		ErrSpreadsheetFlagMissingCode,
		errors.Alert,
		[]string{"Spreadsheet flags missing"},
		[]string{err.Error()},
		[]string{"Either --spreadsheet-id or --spreadsheet-cred flag was not specified"},
		[]string{"Provide both --spreadsheet-id and --spreadsheet-cred flags to generate relationships"})
}

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
