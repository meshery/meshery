package connections

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands = []*cobra.Command{listConnectionsCmd, deleteConnectionCmd}

	pageNumberFlag int
	maxRowsPerPage int
)

var ConnectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "Manage connections",
	Long:  `Manage connections`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// List all the connections
mesheryctl exp connections list
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl exp connections --help' to display usage guide.\n", args[0]), "model"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
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
