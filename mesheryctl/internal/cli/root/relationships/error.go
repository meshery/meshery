package relationships

import (
	"github.com/meshery/meshkit/errors"
)

var (
	ErrEmptySheetDataCode = "mesheryctl-1204"
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
