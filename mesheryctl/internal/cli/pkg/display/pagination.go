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

		pagesQuerySearch := url.Values{}
		if !strings.Contains(displayData.UrlPath, "page") {
			pagesQuerySearch.Set("page", fmt.Sprintf("%d", currentPage))
		}

		if !strings.Contains(displayData.UrlPath, "pagesize") {
			pagesQuerySearch.Set("pagesize", fmt.Sprintf("%d", effectivePageSize))
		}

		if displayData.SearchTerm != "" {
			pagesQuerySearch.Set("search", displayData.SearchTerm)
		}

		if strings.Contains(displayData.UrlPath, "?") {
			urlPath = fmt.Sprintf("%s&%s", displayData.UrlPath, pagesQuerySearch.Encode())
		} else {
			urlPath = fmt.Sprintf("%s?%s", displayData.UrlPath, pagesQuerySearch.Encode())
		}

		data, err := api.Fetch[T](urlPath)
		if err != nil {
			errCode := errors.GetCode(err)
			if slices.Contains(serverAndNetworkErrors, errCode) || errCode == utils.ErrUnmarshalCode {
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
