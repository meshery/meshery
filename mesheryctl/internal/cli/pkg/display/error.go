package display

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

var ErrListPaginationCode = "mesheryctl-1157"
var ErrEncodingDataCode = "mesheryctl-1181"

func ErrorListPagination(err error, currentPage int) error {
	return errors.New(ErrListPaginationCode, errors.Alert,
		[]string{"Failed to fetch data from Meshery server."},
		[]string{fmt.Errorf("failed to fetch data for page %d: %w", currentPage, err).Error()},
		[]string{"While fetching data from Meshery server an error occurred."},
		[]string{"Please check if returned data is valid"})
}

func ErrEncodingData(err error, encoder string) error {
	return errors.New(ErrEncodingDataCode, errors.Alert, []string{fmt.Sprintf("error occured while trying to encode data in %s", encoder)}, []string{fmt.Sprintf("Encoding the data provided failed in %s format", encoder)}, []string{"Non supported characters in the data"}, []string{"Ensure the content of the data provide does not contains invalid supported characters"})
}
