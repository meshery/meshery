package display

import (
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
)

var (
	whiteBoardPrinter    = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
	nextPageKeyboardKeys = []keyboard.Key{keyboard.KeyEnter, keyboard.KeyArrowDown, keyboard.KeySpace}
	escapeKeyboardKeys   = []keyboard.Key{keyboard.KeyEsc, keyboard.KeyCtrlC}
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
	processDataFunc pageHandler[T],
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
			if meshkitErr, ok := err.(*errors.Error); ok {
				if slices.Contains(serverAndNetworkErrors, meshkitErr.Code) {
					return err
				}
				return ErrorListPagination(err, currentPage)
			}
			return err
		}

		// Process the fetched data
		shouldContiue, err := processDataFunc(data, currentPage, effectivePageSize)
		if err != nil {
			return err
		}

		if shouldContiue {
			currentPage++
		} else {
			break
		}

	}

	return nil
}

func listPageCallback[T any](displayData DisplayDataAsync, processDataFunc listRowBuilder[T]) pageHandler[T] {
	startIndex := 0
	return func(data *T, currentPage int, pageSize int) (bool, error) {
		rows, totalCount := processDataFunc(data)

		// Display the total count and current page
		utils.DisplayCount(displayData.DataType, totalCount)

		if len(rows) == 0 {
			return false, nil
		}

		if displayData.DisplayCountOnly {
			return false, nil
		}

		// Display the current page number to be one-based
		_, _ = whiteBoardPrinter.Fprint(os.Stdout, "Page: ", currentPage+1)
		fmt.Println()

		// Display the data in a table
		utils.PrintToTable(displayData.Header, rows, nil)

		if displayData.IsPage {
			return false, nil
		}

		// If the URL already contains "pagesize=all", it means all data has been fetched in one go,
		// so we can break the loop without waiting for user input
		if strings.Contains(displayData.UrlPath, "pagesize=all") {
			return false, nil
		}

		if int64(startIndex+pageSize) >= totalCount {
			return false, nil
		}

		// Wait for user input to navigate pages
		keysEvents, err := keyboard.GetKeys(10)
		if err != nil {
			return false, err
		}

		defer func() {
			_ = keyboard.Close()
		}()

		event := <-keysEvents
		if event.Err != nil {
			utils.Log.Error(fmt.Errorf("unable to capture keyboard events"))
			return false, event.Err
		}

		if slices.Contains(escapeKeyboardKeys, event.Key) {
			return false, nil
		}

		if slices.Contains(nextPageKeyboardKeys, event.Key) {
			startIndex += pageSize
			return true, nil
		} else {
			return false, nil
		}
	}
}

func selectPageCallback[T any, R any](displayData DisplayDataAsync, processData selectRowBuilder[R], extractItem extractItems[T, R], selectedItem *R) pageHandler[T] {
	return func(data *T, currentPage int, pageSize int) (bool, error) {
		rows := extractItem(data)

		switch len(rows) {
		case 0:
			if currentPage == 0 {
				return false, utils.ErrNotFound(fmt.Errorf("no results for %s", displayData.SearchTerm))
			}
			return false, nil
		case 1:
			*selectedItem = rows[0]
			return false, nil
		default:
			picked, itemSelected, err := SelectFromPagedResults(rows, processData, pageSize)
			if err != nil {
				return false, err
			}
			if itemSelected {
				*selectedItem = picked
				return false, nil
			}
			return true, nil
		}
	}
}
