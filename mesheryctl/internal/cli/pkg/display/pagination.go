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

		parsedURL, err := url.Parse(displayData.UrlPath)
		if err != nil {
			return ErrPagination(err, currentPage)
		}

		existingQuery := parsedURL.Query()
		pagesQuerySearch := url.Values{}

		if !existingQuery.Has("page") {
			pagesQuerySearch.Set("page", fmt.Sprintf("%d", currentPage))
		} else if currentPage > 0 {
			existingQuery.Set("page", fmt.Sprintf("%d", currentPage))
		}

		if !existingQuery.Has("pagesize") {
			pagesQuerySearch.Set("pagesize", fmt.Sprintf("%d", effectivePageSize))
		}

		urlPath := displayData.UrlPath
		if displayData.SearchTerm != "" {
			existingQuery.Set("search", displayData.SearchTerm)
		}

		if (existingQuery.Has("search") && displayData.SearchTerm != "") || (existingQuery.Has("page") && currentPage > 0) {
			parsedURL.RawQuery = existingQuery.Encode()
			urlPath = parsedURL.String()
		} else if displayData.SearchTerm != "" && !existingQuery.Has("search") {
			pagesQuerySearch.Set("search", displayData.SearchTerm)
		}

		if queryStr := pagesQuerySearch.Encode(); queryStr != "" {
			if strings.Contains(urlPath, "?") {
				urlPath = fmt.Sprintf("%s&%s", urlPath, queryStr)
			} else {
				urlPath = fmt.Sprintf("%s?%s", urlPath, queryStr)
			}
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
