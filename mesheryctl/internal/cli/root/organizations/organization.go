package organizations

import (
	"errors"
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands = []*cobra.Command{listOrgCmd}
	organizationsApiPath = "api/identity/orgs"
)

var OrgCmd = &cobra.Command{

	Use:   "organization",
	Short: "Interact with registered organizations",
	Long: `Interact with registered organizations to display detailed information
Find more information at: https://docs.meshery.io/reference/mesheryctl/organization`,
	Example: `
// Number of registered orgs
mesheryctl organization --count

// List registered orgs
mesheryctl organization list
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		count, err := cmd.Flags().GetBool("count")
		if err != nil {
			return err
		}
		if utils.IsValidCountUsage(count, args) {
			return nil
		}
		if len(args) == 0 {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New("please provide a subcommand with the command"))
		}
		return utils.ErrInvalidArgument(fmt.Errorf("'%s' is an invalid command. Use 'mesheryctl organization --help' to display usage guide", args[0]))
	},
	RunE: func(_ *cobra.Command, _ []string) error {
		orgs, err := api.Fetch[models.OrganizationsPage](fmt.Sprintf("%s?pagesize=all", organizationsApiPath))
		if err != nil {
			return err
		}
		utils.DisplayCount("organizations", int64(orgs.TotalCount))
		return nil
	},
}

func init() {
	OrgCmd.Flags().BoolP("count", "", false, "total number of registered organizations")
	OrgCmd.AddCommand(availableSubcommands...)

}
