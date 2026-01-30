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

var whiteBoardPrinter = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
var nextPageKeyboardKeys = []keyboard.Key{keyboard.KeyEnter, keyboard.KeyArrowDown, keyboard.KeySpace}
var escapeKeyboardKeys = []keyboard.Key{keyboard.KeyEsc, keyboard.KeyCtrlC}

var serverAndNetworkErrors = []string{
	utils.ErrUnauthenticatedCode,
	utils.ErrInvalidTokenCode,
	utils.ErrAttachAuthTokenCode,
	utils.ErrFailRequestCode,
}

func HandlePaginationAsync[T any](
	pageSize int,
	displayData DisplayDataAsync,
	processDataFunc func(*T) ([][]string, int64),
) error {
	startIndex := 0
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
			pagesQuerySearch.Set("pagesize", fmt.Sprintf("%d", pageSize))
		}

		if strings.Contains(displayData.UrlPath, "?") {
			urlPath = fmt.Sprintf("%s&%s", displayData.UrlPath, pagesQuerySearch.Encode())
		} else {
			urlPath = fmt.Sprintf("%s?%s", displayData.UrlPath, pagesQuerySearch.Encode())
		}

		utils.Log.Debug("Fetching data from URL: ", urlPath)

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
		rows, totalCount := processDataFunc(data)

		// Display the total count and current page
		utils.DisplayCount(displayData.DataType, totalCount)

		if len(rows) == 0 {
			break
		}

		if displayData.DisplayCountOnly {
			return nil
		}

		// Display the current page number to be one-based
		_, _ = whiteBoardPrinter.Fprint(os.Stdout, "Page: ", currentPage+1)
		fmt.Println()

		// Display the data in a table
		utils.PrintToTable(displayData.Header, rows, nil)

		if displayData.IsPage {
			break
		}

		if int64(startIndex+pageSize) >= totalCount {
			break
		}

		// Wait for user input to navigate pages
		keysEvents, err := keyboard.GetKeys(10)
		if err != nil {
			return err
		}

		defer func() {
			_ = keyboard.Close()
		}()

		event := <-keysEvents
		if event.Err != nil {
			utils.Log.Error(fmt.Errorf("unable to capture keyboard events"))
			break
		}

		if slices.Contains(escapeKeyboardKeys, event.Key) {
			break
		}

		if slices.Contains(nextPageKeyboardKeys, event.Key) {
			currentPage++
			startIndex += pageSize
		} else {
			break
		}

	}

	return nil
}
