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
	pageNo   int
	pageSize int
)

type orgsStruct struct {
	Orgs []models.Organization `json:"organizations"`
}

var OrgCmd = &cobra.Command{
	Use:   "org",
	Short: "show existing organizations with their name,id, date of creation",
	Long:  "show existing organiazations with their name, id, date of creation",
	Example: `
	// list all organizations
	mesheryctl exp org
	// list organizations (using flags)
	mesheryctl exp org --page [page_no] --size [size]
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(utils.ErrLoadConfig(err))
			return nil

		}
		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/identity/orgs?page=%d&pagesize=%d", baseUrl, pageNo, pageSize)

		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		defer resp.Body.Close()
		JsonData, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		var orgs orgsStruct
		err = json.Unmarshal(JsonData, &orgs)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		var orgsData [][]string
		columnNames := []string{"Name", "Id", "Created-At"}

		for _, org := range orgs.Orgs {
			orgsData = append(orgsData, []string{org.Name, org.ID.String(), org.CreatedAt.String()})
		}

		utils.PrintToTable(columnNames, orgsData)

		return nil
	},
}

func init() {
	OrgCmd.Flags().IntVarP(&pageNo, "page", "p", 0, "page number to fetch")
	OrgCmd.Flags().IntVarP(&pageSize, "page_size", "s", 10, "page size")
}
