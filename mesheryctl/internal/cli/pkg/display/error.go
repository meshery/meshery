package display

import (
	"fmt"
	"strings"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
)

var ErrListPaginationCode = "mesheryctl-1157"

func ErrorListPagination(err error, currentPage int) error {
	// Check if this is an authentication error by looking at the error message
	errMsg := strings.ToLower(err.Error())
	if strings.Contains(errMsg, "authentication") ||
		strings.Contains(errMsg, "token") ||
		strings.Contains(errMsg, "unauthorized") ||
		strings.Contains(errMsg, "login") {
		return err
	}

	// If this is a meshkit error representing authentication or token issues,
	// return it directly so the caller can act on the specific auth error.
	if meshkitErr, ok := err.(*errors.Error); ok {
		if meshkitErr.Code == utils.ErrUnauthenticatedCode ||
			meshkitErr.Code == utils.ErrInvalidTokenCode ||
			meshkitErr.Code == utils.ErrAttachAuthTokenCode {
			return err
		}
	}

	// Wrap non-authentication errors with pagination context
	return errors.New(ErrListPaginationCode, errors.Alert,
		[]string{"Failed to fetch data from Meshery server."},
		[]string{fmt.Errorf("failed to fetch data for page %d: %w", currentPage, err).Error()},
		[]string{"While fetching data from Meshery server an error occurred."},
		[]string{"Please check if returned data is valid"})
}
