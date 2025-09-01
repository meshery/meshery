package connections

import (
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var listConnectionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the connections",
	Long: `List all available connections.
Documentation for connection can be found at https://docs.meshery.io/reference/mesheryctl/exp/connection/list`,
	Example: `
// List all the connections
mesheryctl exp connection list

// List all the connections with page number
mesheryctl exp connection list --page [page-number]

// Display total count of all available connections
mesheryctl exp connection list --count
`,

	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp connection list \nRun 'mesheryctl exp connection list --help' to see detailed help message"
		if len(args) != 0 {
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		connectionsResponse, err := api.Fetch[connections.ConnectionPage](connectionApiPath)

		if err != nil {
			return err
		}

		header := []string{"id", "Name", "Type", "Kind", "Status"}
		rows := [][]string{}
		for _, connection := range connectionsResponse.Connections {
			rows = append(rows, getConnectionDetail(connection))
		}

		count, _ := cmd.Flags().GetBool("count")
		dataToDisplay := display.DisplayedData{
			DataType:         "connection",
			Header:           header,
			Rows:             rows,
			Count:            int64(connectionsResponse.TotalCount),
			DisplayCountOnly: count,
			IsPage:           cmd.Flags().Changed("page"),
		}
		err = display.List(dataToDisplay)
		if err != nil {
			return err
		}

		return nil
	},
}

func getConnectionDetail(connection *connections.Connection) []string {
	data := make([]string, 5)

	data[0] = connection.ID.String()

	data[1] = connection.Name
	if connection.Name == "" {
		data[1] = "N/A"
	}

	data[2] = connection.Type

	data[3] = connection.Kind

	data[4] = string(connection.Status)
	if string(connection.Status) == "" {
		data[4] = "N/A"
	}

	return data
}

func init() {
	listConnectionsCmd.Flags().BoolP("count", "c", false, "Display the count of total available connections")
	listConnectionsCmd.Flags().IntVarP(&pageNumberFlag, "page", "p", 1, "Page number")
}
