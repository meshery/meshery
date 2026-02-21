package display

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

const pageSize = 10

type DisplayedData struct {
	// Meshery Logical conmponent
	DataType string
	// Header of the table
	Header []string
	// Data to display
	Rows [][]string
	// Total Meshery logical component count available
	Count            int64
	DisplayCountOnly bool
	IsPage           bool
}

type DisplayDataAsync struct {
	UrlPath          string
	DataType         string
	Header           []string
	Page             int
	PageSize         int
	DisplayCountOnly bool
	IsPage           bool
}

type dataProcessor[T any] func(*T) ([][]string, int64)

func List(data DisplayedData) error {
	utils.DisplayCount(data.DataType, data.Count)

	// flag --count is set or no data available
	if data.DisplayCountOnly || data.Count == 0 {
		return nil
	}

	if data.IsPage {
		utils.PrintToTable(data.Header, data.Rows, nil)
	} else {
		maxRowsPerPage := 25
		err := utils.HandlePagination(maxRowsPerPage, data.DataType, data.Rows, data.Header)
		if err != nil {
			return utils.ErrHandlePagination(err)
		}
	}
	return nil
}

func ListAsyncPagination[T any](displayData DisplayDataAsync, processData dataProcessor[T]) error {
	effctivePageSize := pageSize
	if displayData.PageSize > 0 {
		effctivePageSize = displayData.PageSize
	}
	return HandlePaginationAsync(
		effctivePageSize,
		displayData,
		processData,
	)
}

func SelectFromPagedResults[T any](rows []T, formatLabel func([]T) []string, pageSize int) (T, int, error) {
	var zero T

	names := formatLabel(rows)
	if len(rows) < pageSize {
		noMoreLabel := color.New(color.FgHiBlack).Sprint("End of list")
		names = append(names, noMoreLabel)
	} else {
		loadMoreLabel := color.New(color.FgCyan, color.Bold).Sprint("Load More.....")
		names = append(names, loadMoreLabel)
	}

	itemCount := len(rows)

	prompt := promptui.Select{
		Label: "Select item",
		Items: names,
		Size:  5,
		Templates: &promptui.SelectTemplates{
			Help: "Use ↑/↓/←/→ to navigate, Ctrl+C or Esc to cancel",
		},
	}

	maxRetries := 3
	retries := 0
	for {
		i, _, err := prompt.Run()
		if err != nil {
			// Handle ctrl+c
			if err == promptui.ErrInterrupt {
				return zero, -1, fmt.Errorf("Selection cancelled")
			}
			retries++
			if retries >= maxRetries {
				return zero, -1, fmt.Errorf("prompt failed after %d attempts: %w", maxRetries, err)
			}
			continue
		}

		// Last item (Load More / End of list) selected
		if i == itemCount {
			// No more pages just re-show the prompt
			if len(rows) < pageSize {
				continue
			}
			return zero, i, nil
		}

		return rows[i], i, nil
	}
}
