package application

import (
	"fmt"
	"os"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
	file                 string
	tokenPath            string
)

// ApplicationCmd represents the root command for application commands
var ApplicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Service Mesh Applications Management",
	Long:  `Manage service meshes using predefined applications`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		return nil
	},
}

func init() {
	ApplicationCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Path to application file")
	ApplicationCmd.PersistentFlags().StringVarP(&tokenPath, "token", "t", "", "Path to token file")
	_ = ApplicationCmd.MarkFlagRequired("file")

	tokenPath = os.Getenv("MESHERY_AUTH_TOKEN")

	availableSubcommands = []*cobra.Command{applyCmd, deleteCmd, viewCmd, listCmd}
	ApplicationCmd.AddCommand(availableSubcommands...)
}
