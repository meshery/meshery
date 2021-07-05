package auth

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
)

// AuthCmd represents the root command for auth commands
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authtication",
	Long:  `Authentication`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		return nil
	},
}

func init() {
	availableSubcommands = []*cobra.Command{loginCmd}
	AuthCmd.AddCommand(availableSubcommands...)
}
