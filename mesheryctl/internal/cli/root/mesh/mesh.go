package mesh

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
)

// MeshCmd represents the Performance Management CLI command
var MeshCmd = &cobra.Command{
	Use:   "mesh",
	Short: "Service Mesh Lifecycle Management",
	Long:  `Provisioning, configuration, and on-going operational management of service meshes`,
	Example: `
// Lifecycle management of service meshes
mesheryctl mesh
	`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.MeshError(fmt.Sprintf("'%s' is an invalid command. Use 'mesheryctl mesh --help' to display usage guide.\n ", args[0])))
		}

		return nil
	},
}

func init() {
	availableSubcommands = []*cobra.Command{validateCmd, deployCmd, removeCmd}
	MeshCmd.AddCommand(availableSubcommands...)
}
