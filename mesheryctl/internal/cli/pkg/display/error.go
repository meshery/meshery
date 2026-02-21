package display

import (
	"fmt"
	"strings"

	"github.com/meshery/meshkit/errors"
)

var (
	ErrListPaginationCode         = "mesheryctl-1157"
	ErrEncodingDataCode           = "mesheryctl-1183"
	ErrUnsupportedFormatCode      = "mesheryctl-1184"
	ErrOutputFileNotSpecifiedCode = "mesheryctl-1194"
	ErrInvalidOutputFormatCode    = "mesheryctl-1198"
	ErrNoResultsFoundCode         = "mesheryctl-1199"
	ErrPromptFailedCode           = "mesheryctl-1200"
	ErrSelectionCancelledCode     = "mesheryctl-1201"
)

func ErrNoResultsFound(name string) error {
	return errors.New(ErrNoResultsFoundCode, errors.Alert,
		[]string{fmt.Sprintf("No matching results found for %q", name)},
		[]string{fmt.Sprintf("The name %q did not match any available items", name)},
		[]string{"The name may be misspelled or the item may not exist in the registry"},
		[]string{
			"Check the spelling and try again",
			"Run 'mesheryctl model list' to see all available models",
		})
}

func ErrPromptFailed(err error) error {
	return errors.New(ErrPromptFailedCode, errors.Alert,
		[]string{"Unable to display selection menu"},
		[]string{err.Error()},
		[]string{"The selection menu could not be displayed in this terminal"},
		[]string{
			"Make sure you are running this command in a regular terminal window",
			"Try running the command again",
			"You can also use the item's ID directly to skip the selection menu",
		})
}

func ErrSelectionCancelled() error {
	return errors.New(ErrSelectionCancelledCode, errors.Alert,
		[]string{"Operation cancelled"},
		[]string{"No item was selected before the prompt was closed"},
		[]string{"The operation was cancelled before a selection was made"},
		[]string{
			"Run the command again to retry",
			"You can also provide the item's ID directly to skip the selection menu",
		})
}

func ErrorListPagination(err error, currentPage int) error {
	return errors.New(ErrListPaginationCode, errors.Alert,
		[]string{"Failed to fetch data from Meshery server."},
		[]string{fmt.Errorf("failed to fetch data for page %d: %w", currentPage, err).Error()},
		[]string{"While fetching data from Meshery server an error occurred."},
		[]string{"Please check if returned data is valid"})
}

func ErrEncodingData(err error, encoder string) error {
	return errors.New(ErrEncodingDataCode, errors.Alert, []string{fmt.Sprintf("error occurred while trying to encode data in %s", encoder)}, []string{fmt.Sprintf("Encoding the data provided failed in %s format", encoder)}, []string{"Non supported characters in the data"}, []string{"Ensure the content of the data provided does not contain invalid supported characters"})
}

func ErrUnsupportedFormat(format string) error {
	return errors.New(ErrUnsupportedFormatCode, errors.Alert, []string{fmt.Sprintf("The output format '%s' is not supported. ", format)}, []string{fmt.Sprintf("Output format '%s' is not supported. ", format)}, []string{"An unsupported output format was requested. "}, []string{"Specify a supported output format. Choices are 'json' or 'yaml'."})
}

func ErrOutputFileNotSpecified() error {
	return errors.New(ErrOutputFileNotSpecifiedCode, errors.Alert, []string{"Output file path is not specified."}, []string{"The output file path was not provided."}, []string{"An output file path must be specified to save the output."}, []string{"Provide a valid file path."})
}

func ErrInvalidOutputFormat(format string) error {
	return errors.New(
		ErrInvalidOutputFormatCode,
		errors.Alert,
		[]string{"Invalid Output Format"},
		[]string{fmt.Sprintf("Provided output format %q is invalid", format)},
		[]string{"The specified output format is not supported"},
		[]string{fmt.Sprintf("Ensure using [%s] as the output format", strings.Join(validOutputFormat, "|"))},
	)
}
