package filter

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
	file                 string
)

// FilterCmd represents the root command for filter commands
var FilterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Service Mesh Filter Management",
	Long:  ``,
	Example: `
// Base command for WASM filters (experimental feature)
mesheryctl exp filter [subcommands]	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.FilterError(fmt.Sprintf("'%s' is a invalid command.  Use 'mesheryctl exp filter --help' to display usage guide.\n", args[0])))
		}
		return nil
	},
}

func init() {
	FilterCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token file default from current context")

	availableSubcommands = []*cobra.Command{applyCmd, viewCmd, deleteCmd, listCmd}
	FilterCmd.AddCommand(availableSubcommands...)
}
