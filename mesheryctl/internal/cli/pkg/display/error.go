package display

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

var ErrListPaginationCode = "mesheryctl-1157"

func ErrorListPagination(err error, currentPage int) error {
	return errors.New(ErrListPaginationCode, errors.Alert,
		[]string{"Failed to fetch data from Meshery server."},
		[]string{fmt.Errorf("failed to fetch data for page %d: %w", currentPage, err).Error()},
		[]string{"While fetching data from Meshery server an error occurred."},
		[]string{"Please check if returned data is valid"})
}
