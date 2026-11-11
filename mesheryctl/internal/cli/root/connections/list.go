package connections

import (
	"fmt"
	"net/url"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta3/connection"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type connectionListFlags struct {
	Count    bool     `json:"count" validate:"boolean"`
	Kind     []string `json:"kind"`
	Status   []string `json:"status"`
	PageSize int      `json:"pagesize" validate:"gte=1"`
	Page     int      `json:"page" validate:"gte=1"`
}

var connectionListFlagsProvided connectionListFlags

var listConnectionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the connections",
	Long: `List all available connections.
Find more information at: https://docs.meshery.io/reference/mesheryctl/connection/list`,
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

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &connectionListFlagsProvided)
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		urlPath := connectionApiPath
		querySearch := url.Values{}

		// Filters are repeated query params (kind=a&kind=b, status=x&status=y),
		// the standard convention across the API — no JSON encoding.
		for _, kind := range connectionListFlagsProvided.Kind {
			querySearch.Add("kind", kind)
		}
		for _, status := range connectionListFlagsProvided.Status {
			querySearch.Add("status", status)
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
			Page:             connectionListFlagsProvided.Page,
			PageSize:         connectionListFlagsProvided.PageSize,
			IsPage:           cmd.Flags().Changed("page"),
			DisplayCountOnly: connectionListFlagsProvided.Count,
		}

		// Bypass interactive pagination if terminal is non-interactive (stdout redirected)
		if !utils.IsInteractiveTerminal() {
			data.PageSize = 0
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

	data[2] = connection.ConnectionType

	data[3] = connection.Kind

	data[4] = string(connection.Status)
	if string(connection.Status) == "" {
		data[4] = "N/A"
	}

	return data
}

func init() {
	listConnectionsCmd.Flags().BoolVarP(&connectionListFlagsProvided.Count, "count", "c", false, "Display the count of total available connections")
	listConnectionsCmd.Flags().StringSliceVarP(&connectionListFlagsProvided.Kind, "kind", "k", []string{}, "Filter connections by kind")
	listConnectionsCmd.Flags().IntVarP(&connectionListFlagsProvided.Page, "page", "p", 1, "Page number")
	listConnectionsCmd.Flags().IntVarP(&connectionListFlagsProvided.PageSize, "pagesize", "", 10, "Number of connections per page")
	listConnectionsCmd.Flags().StringSliceVarP(&connectionListFlagsProvided.Status, "status", "s", []string{}, "Filter connections by status")
}
