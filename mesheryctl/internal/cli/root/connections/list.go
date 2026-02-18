package connections

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type connectionListFlags struct {
	count    bool
	kind     []string
	status   []string
	pageSize int
	page     int
}

var connectionListFlagsProvided connectionListFlags

var listConnectionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the connections",
	Long: `List all available connections.
Documentation for connection can be found at https://docs.meshery.io/reference/mesheryctl/connection/list`,
	Example: `
// List all the connections
mesheryctl connection list

// List all the connections with page number
mesheryctl connection list --page [page-number]

// List all the connections matching a specific kind and status
mesheryctl connection list --kind [kind] --status [status]

// List all the connections matching a set of kinds and statuses
mesheryctl connection list --kind [kind] --kind [kind] --status [status] --status [status]

// Display total count of all available connections
mesheryctl connection list --count
`,

	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl connection list \nRun 'mesheryctl connection list --help' to see detailed help message"
		if len(args) != 0 {
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		urlPath := connectionApiPath
		querySearch := url.Values{}

		kindQuery, err := json.Marshal(connectionListFlagsProvided.kind)
		if err != nil {
			return utils.ErrMarshal(err)
		}
		if len(connectionListFlagsProvided.kind) > 0 {
			utils.Log.Debug("Adding kind to query: ", string(kindQuery))
			querySearch.Add("kind", string(kindQuery))
		}

		statusQuery, err := json.Marshal(connectionListFlagsProvided.status)
		if err != nil {
			return utils.ErrMarshal(err)
		}

		if len(connectionListFlagsProvided.status) > 0 {
			utils.Log.Debug("Adding status to query: ", string(statusQuery))
			querySearch.Add("status", string(statusQuery))
		}

		if len(querySearch) > 0 {
			urlPath += fmt.Sprintf("?%s", querySearch.Encode())
		}
		utils.Log.Debug("Final URL: ", urlPath)

		header := []string{"id", "Name", "Type", "Kind", "Status"}

		data := display.DisplayDataAsync{
			UrlPath:          urlPath,
			DataType:         "connection",
			Header:           header,
			Page:             connectionListFlagsProvided.page,
			PageSize:         connectionListFlagsProvided.pageSize,
			IsPage:           cmd.Flags().Changed("page"),
			DisplayCountOnly: connectionListFlagsProvided.count,
		}

		return display.ListAsyncPagination(data, processConnectionData)
	},
}

func processConnectionData(connectionsResponse *connection.ConnectionPage) ([][]string, int64) {
	rows := [][]string{}
	for _, connection := range connectionsResponse.Connections {
		row := getConnectionDetail(connection)
		rows = append(rows, row)
	}
	return rows, int64(connectionsResponse.TotalCount)
}

func getConnectionDetail(connection *connection.Connection) []string {
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
	listConnectionsCmd.Flags().BoolVarP(&connectionListFlagsProvided.count, "count", "c", false, "Display the count of total available connections")
	listConnectionsCmd.Flags().StringSliceVarP(&connectionListFlagsProvided.kind, "kind", "k", []string{}, "Filter connections by kind")
	listConnectionsCmd.Flags().IntVarP(&connectionListFlagsProvided.page, "page", "p", 1, "Page number")
	listConnectionsCmd.Flags().IntVarP(&connectionListFlagsProvided.pageSize, "pagesize", "", 10, "Number of connections per page")
	listConnectionsCmd.Flags().StringSliceVarP(&connectionListFlagsProvided.status, "status", "s", []string{}, "Filter connections by status")
}
