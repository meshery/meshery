package pattern

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

// PatternCmd represents the root command for pattern commands
var PatternCmd = &cobra.Command{
	Use:   "pattern",
	Short: "Service Mesh Patterns Management",
	Long:  `Manage service meshes using predefined patterns`,
	Args:  cobra.MinimumNArgs(1),
	Example: `
// Apply pattern file
mesheryctl pattern apply --file [path to pattern file | URL of the file]

// Deprovision pattern file
mesheryctl pattern delete --file [path to pattern file]

// View pattern file
mesheryctl pattern view [pattern name | ID]

// List all patterns
mesheryctl pattern list
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		return nil
	},
}

func init() {
	PatternCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token file default from current context")

	availableSubcommands = []*cobra.Command{applyCmd, deleteCmd, viewCmd, listCmd}
	PatternCmd.AddCommand(availableSubcommands...)
}
