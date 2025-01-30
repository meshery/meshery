package orgs

import (
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands = []*cobra.Command{ListOrgCmd}
	count                bool
)
var OrgCmd = &cobra.Command{

	Use:   "org",
	Short: "view list of registered orgs",
	Long:  "view list of registered orgs with their name,id and date of creation",
	Example: `
	// Number of  registered orgs
	mesheryctl org --count 

	// List registerd orgs
	mesheryctl org list
	
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && !count {
			cmd.Usage()
			utils.ErrInvalidArgument(fmt.Errorf("No arguments passed, provide a subcommand"))
			return nil
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if count {

		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			utils.Log.Error(fmt.Errorf("invalid subcommand"))
			return nil
		}
		return nil
	},
}

func init() {
	OrgCmd.Flags().BoolVarP(&count, "count", "", false, "get the number of registered orgs, --count")
	OrgCmd.AddCommand(availableSubcommands...)

}

func getAllOrgs(countOnly bool) (int, error) {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return 0, utils.ErrLoadConfig(err)
	}
	baseUrl := mctlCfg.GetBaseMesheryURL()
	url := fmt.Sprintf("%s/api/identity/orgs", baseUrl)

	req, err := utils.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	res, err := utils.MakeRequest(req)
	if err != nil {
		return 0, utils.ErrFailRequest(err)
	}
	jsonBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

}
