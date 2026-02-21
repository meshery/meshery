package relationships

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

var (
	ErrSpreadsheetFlagMissingCode = "mesheryctl-1202"
)

var (
	viewUsageMsg           = "\n\nUsage: mesheryctl exp relationship view [model-name]\nRun 'mesheryctl exp relationship view --help' to see detailed help message"
	errNoModelNameProvided = fmt.Errorf("[model-name] isn't specified%s", viewUsageMsg)
	errTooManyArgs         = fmt.Errorf("too many arguments, only [model-name] is expected%s", viewUsageMsg)
	errMsg                 = "Missing flag\nUsage: mesheryctl exp relationship generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED\nRun 'mesheryctl exp relationship generate --help' to see detailed help message"
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
