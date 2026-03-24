package relationships

import (
	"github.com/meshery/meshkit/errors"
)

var (
	ErrEmptySheetDataCode = "mesheryctl-1204"
	ErrEmptyCSVDataCode   = "mesheryctl-1207"
)

const (
	errInvalidArg = "only one argument must be provided and needs to be enclosed by double quotes if it contains spaces (eg. \"model name\", modelName)"
)

func ErrEmptySheetData(err error) error {
	return errors.New(
		ErrEmptySheetDataCode,
		errors.Alert,
		[]string{"Invalid or empty spreadsheet provided"},
		[]string{err.Error()},
		[]string{"The spreadsheet must contain valid relationship entries with the required headers."},
		[]string{"Ensure the spreadsheet contains at least two headers and corresponding values for it's headers. Refer to the Meshery relationship csv for the expected spreadsheet format \"https://github.com/meshery/meshery/blob/master/mesheryctl/templates/template-csvs/Relationships.csv\""},
	)
}

func ErrEmptyCSVData(err error) error {
	return errors.New(
		ErrEmptyCSVDataCode,
		errors.Alert,
		[]string{"Invalid or empty CSV file provided"},
		[]string{err.Error()},
		[]string{"The CSV file must contain valid relationship entries with the required headers."},
		[]string{"Ensure the CSV file contains at least two header rows and data rows with a minimum of 15 columns. Refer to the Meshery relationship CSV template: \"https://github.com/meshery/meshery/blob/master/mesheryctl/templates/template-csvs/Relationships.csv\""},
	)
}
