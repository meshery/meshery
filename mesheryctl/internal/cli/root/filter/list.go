package filter

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
	verbose bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List filters",
	Long:  `Display list of all available filter files.`,
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		var response models.FiltersAPIResponse
		client := &http.Client{}
		req, err := http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/experimental/filter", nil)
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
		err = json.Unmarshal(body, &response)

		if err != nil {
			return err
		}
		tokenObj, err := utils.ReadToken(token)
		if err != nil {
			return err
		}
		provider := tokenObj["meshery-provider"]
		var data [][]string

		if verbose {
			if provider == "None" {
				for _, v := range response.Filters {
					FilterID := v.ID.String()
					FilterName := v.Name
					CreatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year(), v.CreatedAt.Hour(), v.CreatedAt.Minute(), v.CreatedAt.Second())
					UpdatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year(), v.UpdatedAt.Hour(), v.UpdatedAt.Minute(), v.UpdatedAt.Second())
					data = append(data, []string{FilterID, FilterName, CreatedAt, UpdatedAt})
				}
				if len(data) > 0 {
					utils.PrintToTableWithFooter([]string{"FILTER ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", ""})
				}
				return errors.New("no filters available")
			}

			for _, v := range response.Filters {
				FilterID := utils.TruncateID(v.ID.String())
				var UserID string
				if v.UserID != nil {
					UserID = *v.UserID
				} else {
					UserID = "null"
				}
				FilterName := v.Name
				CreatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year(), v.CreatedAt.Hour(), v.CreatedAt.Minute(), v.CreatedAt.Second())
				UpdatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year(), v.UpdatedAt.Hour(), v.UpdatedAt.Minute(), v.UpdatedAt.Second())
				data = append(data, []string{FilterID, UserID, FilterName, CreatedAt, UpdatedAt})
			}
			if len(data) > 0 {
				utils.PrintToTableWithFooter([]string{"FILTER ID", "USER ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", "", ""})
			}
			return errors.New("no filters available")
		}

		// Check if messhery provider is set
		if provider == "None" {
			for _, v := range response.Filters {
				FilterName := fmt.Sprintf("%s", strings.Trim(v.Name, filepath.Ext(v.Name)))
				FilterID := utils.TruncateID(v.ID.String())
				CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
				UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
				data = append(data, []string{FilterID, FilterName, CreatedAt, UpdatedAt})
			}
			if len(data) > 0 {
				utils.PrintToTableWithFooter([]string{"FILTER ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", ""})
			}
			return errors.New("no filters available")
		}
		for _, v := range response.Filters {
			FilterID := utils.TruncateID(v.ID.String())
			var UserID string
			if v.UserID != nil {
				UserID = *v.UserID
			} else {
				UserID = "null"
			}
			FilterName := v.Name
			CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
			UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
			data = append(data, []string{FilterID, UserID, FilterName, CreatedAt, UpdatedAt})
		}
		if len(data) > 0 {
			utils.PrintToTableWithFooter([]string{"FILTER ID", "USER ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", "", ""})
		}
		return errors.New("no filters available")

	},
}

func init() {
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display full length user and filter file identifiers")
	listCmd.Flags().StringVarP(&token, "token", "t", "", "path to token")
	_ = listCmd.MarkFlagRequired("token")
}
