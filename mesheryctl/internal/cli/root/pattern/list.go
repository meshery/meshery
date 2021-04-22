package pattern

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

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
	Long: "Display list of all available pattern files",
	Args: cobra.MinimumNArgs(0),
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
		tokenObj, err := utils.ReadToken(token)
		if err != nil {
			return err
		}
		provider := tokenObj["meshery-provider"]
		var data [][]string

		if allflag == true {
			if provider == "None" {
				for _, v := range response.Patterns {
					PatternId := v.ID.String()
					PatterName := v.Name
					CreatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year(), v.CreatedAt.Hour(), v.CreatedAt.Minute(), v.CreatedAt.Second())
					UpdatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year(), v.UpdatedAt.Hour(), v.UpdatedAt.Minute(), v.UpdatedAt.Second())
					data = append(data, []string{PatternId, PatterName, CreatedAt, UpdatedAt})
				}
				utils.PrintToTable([]string{"PATTERN ID", "NAME", "CREATED", "UPDATED"}, data)
				return nil
			}

			for _, v := range response.Patterns {
				PatternId := utils.TruncateID(v.ID.String())
				var UserId string
				if v.UserID != nil {
					UserId = *v.UserID
				} else {
					UserId = "null"
				}
				PatterName := v.Name
				CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
				UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
				data = append(data, []string{PatternId, UserId, PatterName, CreatedAt, UpdatedAt})
			}
			utils.PrintToTable([]string{"PATTERN ID", "USER ID", "NAME", "CREATED", "UPDATED"}, data)

			return nil
		}

		// Check if messhery provider is set
		if provider == "None" {
			for _, v := range response.Patterns {
				PatterName := fmt.Sprintf("%s", strings.Trim(v.Name, filepath.Ext(v.Name)))
				PatternId := utils.TruncateID(v.ID.String())
				CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
				UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
				data = append(data, []string{PatternId, PatterName, CreatedAt, UpdatedAt})
			}
			utils.PrintToTable([]string{"PATTERN ID", "NAME", "CREATED", "UPDATED"}, data)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&allflag, "all", "a", false, "Display full length user and pattern file identifiers")
	listCmd.Flags().StringVarP(&token, "token", "t", "", "path to token")
	listCmd.MarkFlagRequired("token")
}
