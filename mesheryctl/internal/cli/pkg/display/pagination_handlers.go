package display

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

var (
	whiteBoardPrinter    = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
	nextPageKeyboardKeys = []keyboard.Key{keyboard.KeyEnter, keyboard.KeyArrowDown, keyboard.KeySpace}
	escapeKeyboardKeys   = []keyboard.Key{keyboard.KeyEsc, keyboard.KeyCtrlC}
)

// listPageHandler creates a pageHandler that displays each fetched page as a table.
// It tracks the current index across pages to determine when all data has been shown
func listPageHandler[T any](displayData DisplayDataAsync, processDataFunc listRowBuilder[T]) pageHandler[T] {
	startIndex := 0
	return func(data *T, currentPage int, pgSize int) (bool, error) {
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

		if int64(startIndex+pgSize) >= totalCount {
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
			startIndex += pgSize
			return true, nil
		} else {
			return false, nil
		}
	}
}

// promptPageHandler creates a pageHandler that presents each fetched page as a selection prompt.
func promptPageHandler[T any, R any](displayData DisplayDataAsync, processData promptLabelBuilder[R], extractItem itemExtractor[T, R], selectedItem *R) pageHandler[T] {
	return func(data *T, currentPage int, pgSize int) (bool, error) {
		rows, totalCount := extractItem(data)

		if len(rows) == 0 {
			if currentPage == 0 {
				return false, utils.ErrNotFound(fmt.Errorf("no results found"))
			}
			return false, nil
		}

		// Auto-select if only one result on the first page
		if len(rows) == 1 && currentPage == 0 {
			*selectedItem = rows[0]
			return false, nil
		}

		picked, itemSelected, err := SelectFromPagedResults(rows, processData, pgSize, currentPage, totalCount)
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
