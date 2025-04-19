package display

import (
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

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
