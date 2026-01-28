package connections

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const connectionApiPath = "api/integrations/connections"

var (
	availableSubcommands = []*cobra.Command{listConnectionsCmd, deleteConnectionCmd, viewConnectionCmd, createConnectionCmd}
)

var ConnectionsCmd = &cobra.Command{
	Use:   "connection",
	Short: "Manage Meshery connections",
	Long: `View and manage your Meshery connection.
Documentation for connection can be found at https://docs.meshery.io/reference/mesheryctl/connection`,
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
		countFlag, _ := cmd.Flags().GetBool("count")
		if len(args) == 0 && !countFlag {
			return utils.ErrInvalidArgument(errors.New("no subcommand provided for connection"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		countFlag, _ := cmd.Flags().GetBool("count")
		if countFlag {
			connectionsResponse, err := api.Fetch[connection.ConnectionPage](connectionApiPath)

			if err != nil {
				return err
			}

			utils.DisplayCount("connection", int64(connectionsResponse.TotalCount))

			return nil
		}

		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return utils.ErrInvalidArgument(fmt.Errorf("'%s' is an invalid subcommand. Use 'mesheryctl connection --help' to display usage.\n", args[0]))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = cmd.Usage()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	ConnectionsCmd.AddCommand(availableSubcommands...)
	ConnectionsCmd.Flags().BoolP("count", "c", false, "Display the count of total available connections")
}
