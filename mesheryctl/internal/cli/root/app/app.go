package app

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

// AppCmd represents the root command for app commands
var AppCmd = &cobra.Command{
	Use:   "app",
	Short: "Service Mesh Apps Management",
	Long:  `Manage apps using service mesh`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		return nil
	},
}

func init() {
	AppCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Path to app file")
	AppCmd.PersistentFlags().StringVarP(&tokenPath, "token", "t", "", "Path to token file")
	_ = AppCmd.MarkFlagRequired("file")

	tokenPath = os.Getenv("MESHERY_AUTH_TOKEN")

	availableSubcommands = []*cobra.Command{onboardCmd, viewCmd, offboardCmd, listCmd}
	AppCmd.AddCommand(availableSubcommands...)
}
