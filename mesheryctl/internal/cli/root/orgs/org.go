package orgs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
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
			return cmd.Help()
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if count {
			orgs, err := getAllOrgs()
			if err != nil {
				utils.Log.Error(err)
				return nil
			}
			fmt.Printf("Total registered orgs : %v", orgs.TotalCount)
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
	OrgCmd.Flags().BoolVarP(&count, "count", "", false, "total number of registered orgs")
	OrgCmd.AddCommand(availableSubcommands...)

}

func getAllOrgs() (*models.OrganizationsPage, error) {
	var orgPage models.OrganizationsPage
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return nil, utils.ErrLoadConfig(err)
	}
	baseUrl := mctlCfg.GetBaseMesheryURL()
	url := fmt.Sprintf("%s/api/identity/orgs?all=true", baseUrl)

	req, err := utils.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := utils.MakeRequest(req)
	if err != nil {
		return nil, utils.ErrFailRequest(err)
	}
	defer res.Body.Close()

	jsonBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &orgPage)
	if err != nil {
		return nil, err
	}

	return &orgPage, nil

}
