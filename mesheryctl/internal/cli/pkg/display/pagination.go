package display

import (
	"fmt"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var whiteBoardPrinter = color.New(color.FgHiBlack, color.BgWhite, color.Bold)

func HandlePaginationAsync[T any](
	pageSize int,
	displayData DisplayDataAsync,
	processDataFunc func(*T) ([][]string, int64),
) error {
	startIndex := 0
	currentPage := displayData.Page

	for {
		// Clear the terminal screen
		utils.ClearLine()

		// Fetch data for the current page
		urlPath := fmt.Sprintf("%s?page=%d&page_size=%d", displayData.UrlPath, currentPage, pageSize)
		data, err := api.Fetch[T](urlPath)
		if err != nil {
			return fmt.Errorf("failed to fetch data for page %d: %w", currentPage, err)
		}

		// Process the fetched data
		rows, totalCount := processDataFunc(data)

		if len(rows) == 0 {
			whiteBoardPrinter.Print("No data available.")
			fmt.Println()
			break
		}

		// Display the total count and current page
		whiteBoardPrinter.Print("Total number of ", displayData.DataType, ":", totalCount)
		fmt.Println()

		if displayData.DisplayCountOnly {
			return nil
		}

		whiteBoardPrinter.Print("Page: ", currentPage)
		fmt.Println()

		// Display the data in a table
		utils.PrintToTable(displayData.Header, rows)

		if displayData.IsPage {
			break
		}

		// Check if there is more data to display
		if len(rows) < pageSize {
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
