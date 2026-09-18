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

// DefaultPageSize is the number of items fetched per page when the caller
// does not specify --pagesize. Exported so other call sites that build their
// own single-page request (e.g. a structured --output-format path that
// bypasses HandlePaginationAsync, see organizations/list.go) share this
// default instead of hard-coding the same magic number a second time.
const DefaultPageSize = 10

// NormalizePagination converts a 1-based --page value and an optional
// --pagesize value into the zero-based page index and effective page size
// used to build a request. This is the exact conversion HandlePaginationAsync
// applies below, extracted so other call sites don't have to reimplement (and
// risk drifting from) the same logic.
func NormalizePagination(page, pagesize int) (currentPage, effectivePageSize int) {
	effectivePageSize = DefaultPageSize
	if pagesize > 0 {
		effectivePageSize = pagesize
	}

	// Adjust the page number to be zero-based
	currentPage = page - 1
	if currentPage < 0 {
		currentPage = 0
	}

	return currentPage, effectivePageSize
}

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
	currentPage, effectivePageSize := NormalizePagination(displayData.Page, displayData.PageSize)

	for {
		// Clear the terminal screen
		if currentPage > 0 {
			utils.ClearLine()
		}

		urlPath := ""

		pagesQuerySearch := url.Values{}
		if !strings.Contains(displayData.UrlPath, "page=") {
			pagesQuerySearch.Set("page", fmt.Sprintf("%d", currentPage))
		}

		if !strings.Contains(displayData.UrlPath, "pagesize=") {
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
