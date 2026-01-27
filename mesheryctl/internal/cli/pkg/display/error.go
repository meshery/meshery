package display

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

var ErrListPaginationCode = "mesheryctl-1157"
var ErrEncodingDataCode = "mesheryctl-1183"
var ErrUnsupportedFormatCode = "mesheryctl-1184"
var ErrOutputFileNotSpecifiedCode = "mesheryctl-1194"

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
