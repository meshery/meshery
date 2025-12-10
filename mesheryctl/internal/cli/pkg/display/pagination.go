package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
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
			// Check if this is an authentication error by looking at the error message
			errMsg := strings.ToLower(err.Error())
			if strings.Contains(errMsg, "authentication") ||
				strings.Contains(errMsg, "token") ||
				strings.Contains(errMsg, "unauthorized") ||
				strings.Contains(errMsg, "login") {
				return err
			}

			// If this is a meshkit error, check for authentication/network errors
			if meshkitErr, ok := err.(*errors.Error); ok {
				// Authentication errors: return as-is so caller can handle login/token flows
				if meshkitErr.Code == utils.ErrUnauthenticatedCode ||
					meshkitErr.Code == utils.ErrInvalidTokenCode ||
					meshkitErr.Code == utils.ErrAttachAuthTokenCode {
					return err
				}

				// Network/request failures: return as-is
				if meshkitErr.Code == utils.ErrFailRequestCode {
					return err
				}

				// For other meshkit errors, wrap with pagination context
				return ErrorListPagination(err, currentPage)
			}

			// Non-meshkit errors: wrap with pagination context
			return ErrorListPagination(err, currentPage)
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
		utils.PrintToTable(displayData.Header, rows, nil)

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
