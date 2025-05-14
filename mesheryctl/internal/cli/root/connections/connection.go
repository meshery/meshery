package connections

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const connectionApiPath = "api/integrations/connections"

var (
	availableSubcommands = []*cobra.Command{listConnectionsCmd, deleteConnectionCmd}

	pageNumberFlag int
)

var ConnectionsCmd = &cobra.Command{
	Use:   "connection",
	Short: "Manage Meshery connection",
	Long: `View and manage your Meshery connection.
Documentation for connection can be found at https://docs.meshery.io/reference/mesheryctl/exp/connection`,
	Example: `
// List all the connection
mesheryctl exp connection list

// Delete a connection
mesheryctl exp connection delete [connection_id]

`,
	Args: func(cmd *cobra.Command, args []string) error {
		countFlag, _ := cmd.Flags().GetBool("count")
		if len(args) == 0 && !countFlag {
			if err := cmd.Usage(); err != nil {
				return nil
			}
			return errors.New("please provide a subcommand")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		countFlag, _ := cmd.Flags().GetBool("count")
		if countFlag {
			connectionsResponse, err := api.Fetch[connections.ConnectionPage](connectionApiPath)

			if err != nil {
				return err
			}

			utils.DisplayCount("connections", int64(connectionsResponse.TotalCount))

			return nil
		}

		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl exp connection --help' to display usage guide.\n", args[0]), "connection"))
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
