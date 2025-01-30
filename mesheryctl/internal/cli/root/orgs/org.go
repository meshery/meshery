package orgs

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands = []*cobra.Command{ListOrgCmd}
)
var OrgCmd = &cobra.Command{

	Use:   "org",
	Short: "view list of registered orgs",
	Long:  "view list of registered orgs with their name,id and date of creation",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				utils.Log.Error(err)
				return nil
			}
			return nil
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			utils.Log.Error(fmt.Errorf("invalid subcommand"))
			return nil
		}
		return nil
	},
}

func init() {
	OrgCmd.AddCommand(availableSubcommands...)

}
