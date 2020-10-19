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
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		return nil
	},
}

func init() {
	availableSubcommands = []*cobra.Command{validateCmd}
	MeshCmd.PersistentFlags().StringVarP(&serverURL, "serverURL", "u", "localhost:9081", "Server address where meshery is hosted. default 'localhost:9081'")
	_ = MeshCmd.MarkFlagRequired("serverURL")
	MeshCmd.PersistentFlags().StringVarP(&tokenPath, "tokenPath", "t", "", "Path to token for authenticating to Meshery API")
	_ = MeshCmd.MarkFlagRequired("tokenPath")
	MeshCmd.AddCommand(availableSubcommands...)
}
