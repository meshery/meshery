package display

import (
	"fmt"
	"net/url"
	"slices"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
)

const pageSize = 10

var serverAndNetworkErrors = []string{
	utils.ErrUnauthenticatedCode,
	utils.ErrInvalidTokenCode,
	utils.ErrAttachAuthTokenCode,
	utils.ErrFailRequestCode,
}

// HandlePaginationAsync fetches and processes paginated results from the Meshery API,
// filling missing or empty page, pagesize, and search query values while preserving
// non-empty values on displayData.UrlPath.
func HandlePaginationAsync[T any](
	displayData DisplayDataAsync,
	pageHandlerFunc pageHandler[T],
) error {
	effectivePageSize := pageSize
	if displayData.PageSize > 0 {
		effectivePageSize = displayData.PageSize
	}

	// Adjust the page number to be zero-based
	currentPage := displayData.Page - 1

	if currentPage < 0 {
		currentPage = 0
	}

	for {
		// Clear the terminal screen
		if currentPage > 0 {
			utils.ClearLine()
		}

		parsedURL, err := url.Parse(displayData.UrlPath)
		if err != nil {
			return ErrPagination(err, currentPage)
		}

		query := parsedURL.Query()
		if query.Get("page") == "" {
			query.Set("page", fmt.Sprintf("%d", currentPage))
		}

		if query.Get("pagesize") == "" {
			query.Set("pagesize", fmt.Sprintf("%d", effectivePageSize))
		}

		if displayData.SearchTerm != "" && query.Get("search") == "" {
			query.Set("search", displayData.SearchTerm)
		}

		parsedURL.RawQuery = query.Encode()
		urlPath := parsedURL.String()

		data, err := api.Fetch[T](urlPath)
		if err != nil {
			errCode := errors.GetCode(err)
			if slices.Contains(serverAndNetworkErrors, errCode) || errCode == utils.ErrUnmarshalCode || errCode == utils.ErrNotFoundCode {
				return err
			}

			return ErrPagination(err, currentPage)
		}

		// Process the fetched data
		shouldContinue, err := pageHandlerFunc(data, currentPage, effectivePageSize)
		if err != nil {
			return err
		}

		if shouldContinue {
			currentPage++
		} else {
			break
		}

	}

	return nil
}
