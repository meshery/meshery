package connections

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta3/connection"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const connectionApiPath = "api/integrations/connections"

var (
	availableSubcommands = []*cobra.Command{listConnectionsCmd, deleteConnectionCmd, viewConnectionCmd, createConnectionCmd}
)

var ConnectionsCmd = &cobra.Command{
	Use:   "connection",
	Short: "Manage Meshery connections",
	Long: `View and manage your Meshery connection.
Find more information at: https://docs.meshery.io/reference/mesheryctl/connection`,
	Example: `
// Display total count of all available connections
mesheryctl connection --count

// Create a new Kubernetes connection using a specific type
mesheryctl connection create --type aks
mesheryctl connection create --type eks
mesheryctl connection create --type gke
mesheryctl connection create --type minikube

// List all the connection
mesheryctl connection list
mesheryctl connection list --count

// Delete a connection
mesheryctl connection delete [connection_id]
`,
	Args: func(cmd *cobra.Command, args []string) error {
		count, err := cmd.Flags().GetBool("count")
		if err != nil {
			return err
		}
		if utils.IsValidCountUsage(count, args) {
			return nil
		}
		if len(args) == 0 {
			return utils.ErrInvalidArgument(errors.New("no subcommand provided for connection"))
		}
		return utils.ErrInvalidArgument(fmt.Errorf("'%s' is an invalid subcommand. Use 'mesheryctl connection --help' to display usage", args[0]))
	},
	RunE: func(_ *cobra.Command, _ []string) error {
		connectionsResponse, err := api.Fetch[connection.ConnectionPage](connectionApiPath)

		if err != nil {
			return err
		}

		utils.DisplayCount("connection", int64(connectionsResponse.TotalCount))

		return nil
	},
}

func init() {
	ConnectionsCmd.AddCommand(availableSubcommands...)
	ConnectionsCmd.Flags().BoolP("count", "c", false, "Display the count of total available connections")
}
