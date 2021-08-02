package filter

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

// FilterCmd represents the root command for filter commands
var FilterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Service Mesh Filter Management",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		return nil
	},
}

func init() {
	FilterCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Path to filter file")
	FilterCmd.PersistentFlags().StringVarP(&tokenPath, "token", "t", "", "Path to token file")
	_ = FilterCmd.MarkFlagRequired("file")

	tokenPath = os.Getenv("MESHERY_AUTH_TOKEN")

	availableSubcommands = []*cobra.Command{applyCmd, deleteCmd, viewCmd, listCmd}
	FilterCmd.AddCommand(availableSubcommands...)
}
