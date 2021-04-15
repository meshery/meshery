package pattern

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	token   string
	allflag bool
)
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list pattern files",
	Long:  "Display list of all available pattern files",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		var response models.PatternsApiResponse

		client := &http.Client{}
		req, err := http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/experimental/patternfile", nil)
		if err != nil {
			return err
		}
		err = utils.AddAuthDetails(req, token)
		if err != nil {
			return err
		}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		json.Unmarshal(body, &response)
		if allflag == true {
			for _, v := range response.Patterns {
				headers := []string{"PATTERN ID", "NAME", "USER ID", "CREATED", "UPDATED"}
				data := [][]string{
					{v.ID.String(), v.Name, v.ID.String(), v.CreatedAt.Format(time.RFC3339Nano), v.UpdatedAt.Format(time.RFC3339Nano)},
				}
				utils.PrintToTable(headers, data)
			}
			return nil
		}
		for _, v := range response.Patterns {
			headers := []string{"NAME", "USER ID", "CREATED", "UPDATED"}
			data := [][]string{
				{v.Name, v.ID.String(), v.CreatedAt.Format(time.RFC3339Nano), v.UpdatedAt.Format(time.RFC3339Nano)},
			}
			utils.PrintToTable(headers, data)
		}
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&allflag, "all", "a", false, "Display full length user and pattern file identifiers")
	listCmd.MarkFlagRequired("token")
}
