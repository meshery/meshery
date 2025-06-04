package display

import (
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
		utils.PrintToTable(data.Header, data.Rows)
	} else {
		maxRowsPerPage := 25
		err := utils.HandlePagination(maxRowsPerPage, data.DataType, data.Rows, data.Header)
		if err != nil {
			return err
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
