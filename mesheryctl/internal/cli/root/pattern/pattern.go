package pattern

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

// PatternCmd represents the root command for pattern commands
var PatternCmd = &cobra.Command{
	Use:   "pattern",
	Short: "Meshery Patterns Management",
	Long:  `Manage meshery patterns from the manage`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		return nil
	},
}

func init() {
	PatternCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Path to pattern file")
	_ = PatternCmd.MarkFlagRequired("file")

	tokenPath = os.Getenv("MESHERY_AUTH_TOKEN")

	availableSubcommands = []*cobra.Command{applyCmd, deleteCmd}
	PatternCmd.AddCommand(availableSubcommands...)
}
