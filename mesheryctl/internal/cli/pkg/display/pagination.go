package display

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

var whiteBoardPrinter = color.New(color.FgHiBlack, color.BgWhite, color.Bold)

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
		utils.ClearLine()

		// Fetch data for the current page
		urlPath := fmt.Sprintf("%s?page=%d&pagesize=%d", displayData.UrlPath, currentPage, pageSize)
		data, err := api.Fetch[T](urlPath)
		if err != nil {
			return fmt.Errorf("failed to fetch data for page %d: %w", currentPage, err)
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
		whiteBoardPrinter.Fprint(os.Stdout, "Page: ", currentPage+1)
		fmt.Println()

		// Display the data in a table
		utils.PrintToTable(displayData.Header, rows)

		if displayData.IsPage {
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

		if event.Key == keyboard.KeyEsc || event.Key == keyboard.KeyCtrlC {
			break
		}

		if event.Key == keyboard.KeyEnter || event.Key == keyboard.KeyArrowDown {
			currentPage++
			startIndex += pageSize
		} else {
			break
		}

	}

	return nil
}
