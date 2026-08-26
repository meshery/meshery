package display

import (
	"fmt"
	"strings"

	"github.com/meshery/meshkit/errors"
)

var (
	ErrPaginationCode             = "mesheryctl-1157"
	ErrEncodingDataCode           = "mesheryctl-1183"
	ErrUnsupportedFormatCode      = "mesheryctl-1184"
	ErrOutputFileNotSpecifiedCode = "mesheryctl-1194"
	ErrInvalidOutputFormatCode    = "mesheryctl-1198"
	ErrAmbiguousSelectionCode     = "mesheryctl-1253"
)

func ErrPagination(err error, currentPage int) error {
	return errors.New(ErrPaginationCode, errors.Alert,
		[]string{"Failed to fetch paginated data from Meshery server."},
		[]string{fmt.Errorf("failed to fetch data for page %d: %w", currentPage, err).Error()},
		[]string{"An error occurred while fetching paginated data from Meshery server."},
		[]string{"Please check if the server is running and the returned data is valid."})
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

// ErrAmbiguousSelection reports a lookup that matched more than one item while
// running without a terminal to disambiguate it on. Selecting one on the user's
// behalf would act on a resource they did not name, so this is a hard failure -
// but it must say so plainly, rather than surfacing the selection library's
// "open /dev/tty" as if the device were the problem.
func ErrAmbiguousSelection(matches int64) error {
	return errors.New(
		ErrAmbiguousSelectionCode,
		errors.Alert,
		[]string{"Ambiguous match, and no terminal to resolve it on"},
		[]string{fmt.Sprintf("%d items matched, which requires an interactive selection, but mesheryctl is not attached to a terminal", matches)},
		[]string{"The command is running in a script, pipeline or CI job where stdin or stdout is redirected"},
		[]string{"Identify the item by its unique ID instead of its name", "Narrow the search so exactly one item matches"},
	)
}
