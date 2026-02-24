package display

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

type DisplayedData struct {
	// Meshery Logical component
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
	// Core fields
	UrlPath  string
	PageSize int
	Page     int
	// List-only fields
	DataType         string
	Header           []string
	DisplayCountOnly bool
	IsPage           bool
	// Prompt-only field
	SearchTerm string
}

type (
	listRowBuilder[T any]       func(data *T) (rows [][]string, totalCount int64)
	promptLabelBuilder[R any]   func(rows []R) (names []string)
	pageHandler[T any]          func(data *T, currentPage int, pgSize int) (shouldContinue bool, err error)
	itemExtractor[T any, R any] func(data *T) (rows []R)
)

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

func ListAsyncPagination[T any](displayData DisplayDataAsync, processData listRowBuilder[T]) error {
	return HandlePaginationAsync(
		displayData,
		listPageHandler(displayData, processData),
	)
}

func PromptAsyncPagination[T any, R any](displayData DisplayDataAsync, processData promptLabelBuilder[R], extractItems itemExtractor[T, R], selectedItem *R) error {
	return HandlePaginationAsync(
		displayData,
		promptPageHandler(displayData, processData, extractItems, selectedItem),
	)
}

func SelectFromPagedResults[T any](rows []T, formatLabel func([]T) []string, pageSize int) (selected T,
	itemSelected bool,
	err error,
) {
	var zero T
	itemCount := len(rows)

	names := formatLabel(rows)
	if itemCount < pageSize {
		noMoreLabel := color.New(color.FgHiBlack).Sprint("End of list")
		names = append(names, noMoreLabel)
	} else {
		loadMoreLabel := color.New(color.FgCyan, color.Bold).Sprint("Load More.....")
		names = append(names, loadMoreLabel)
	}

	prompt := promptui.Select{
		Label: "Select item",
		Items: names,
		Size:  6,
		Templates: &promptui.SelectTemplates{
			Help: "Use ↑/↓/←/→ to navigate, Ctrl+C to cancel",
		},
	}

	maxRetries := 3
	retries := 0
	for {
		i, _, err := prompt.Run()
		if err != nil {
			// Handle ctrl+c
			if err == promptui.ErrInterrupt {
				return zero, false, fmt.Errorf("selection cancelled")
			}
			retries++
			if retries >= maxRetries {
				return zero, false, fmt.Errorf("prompt failed after %d attempts: %w", maxRetries, err)
			}
			continue
		}

		// Last item (Load More | End of list) selected
		if i == itemCount {
			// No more items to show
			if itemCount < pageSize {
				return zero, false, fmt.Errorf("no more items available")
			}
			return zero, false, nil
		}

		return rows[i], true, nil
	}
}
