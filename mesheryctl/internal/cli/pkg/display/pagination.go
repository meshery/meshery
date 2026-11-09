package display

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

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
// appending page, pagesize, and search query parameters when they are not already
// populated with non-empty values on displayData.UrlPath.
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

		urlPath := ""

		var existingQuery url.Values
		if parsedURL, err := url.Parse(displayData.UrlPath); err == nil {
			existingQuery = parsedURL.Query()
		}

		pagesQuerySearch := url.Values{}
		if existingQuery.Get("page") == "" {
			pagesQuerySearch.Set("page", fmt.Sprintf("%d", currentPage))
		}

		if existingQuery.Get("pagesize") == "" {
			pagesQuerySearch.Set("pagesize", fmt.Sprintf("%d", effectivePageSize))
		}

		if displayData.SearchTerm != "" && existingQuery.Get("search") == "" {
			pagesQuerySearch.Set("search", displayData.SearchTerm)
		}

		encodedQuery := pagesQuerySearch.Encode()
		switch {
		case encodedQuery == "":
			urlPath = displayData.UrlPath
		case strings.Contains(displayData.UrlPath, "?"):
			urlPath = fmt.Sprintf("%s&%s", displayData.UrlPath, encodedQuery)
		default:
			urlPath = fmt.Sprintf("%s?%s", displayData.UrlPath, encodedQuery)
		}

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
