package connections

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands = []*cobra.Command{listConnectionsCmd, deleteConnectionCmd}

	pageNumberFlag int
)

var ConnectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "Manage Meshery connections",
	Long:  `View and manage your Meshery connections`,
	Example: `
// List all the connections
mesheryctl exp connections list

// Delete a connection
mesheryctl exp connections delete [connection_id]

`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if err := cmd.Usage(); err != nil {
				return nil
			}
			return errors.New("please provide a subcommand")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl exp connections --help' to display usage guide.\n", args[0]), "model"))
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
	listConnectionsCmd.Flags().BoolP("count", "c", false, "Display the count of models")
	listConnectionsCmd.Flags().IntVarP(&pageNumberFlag, "page", "p", 1, "Page number")
	deleteConnectionCmd.Flags().StringP("id", "i", "", "ID of the connection to be deleted")
	ConnectionsCmd.AddCommand(availableSubcommands...)
}
