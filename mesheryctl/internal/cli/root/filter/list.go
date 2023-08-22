// Copyright 2023 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	pageSize   = 25
	pageNumber int
	verbose    bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List filters",
	Long:  `Display list of all available filter files.`,
	Example: `
// List all WASM filter files present
mesheryctl filter list	(maximum 25 filters)

// Search for filter
mesheryctl filter list Test (maximum 25 filters)

// Search for filter with space
mesheryctl filter list 'Test Filter' (maximum 25 filters)
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		var searchString string
		if len(args) > 0 {
			searchString = strings.ReplaceAll(args[0], " ", "%20")
		}
		response, err := fetchFilters(mctlCfg.GetBaseMesheryURL(), searchString, pageSize, pageNumber-1)
		if err != nil {
			return err
		}

		if len(args) > 0 && len(response.Filters) == 0 {
			utils.Log.Info("No WASM Filter to display with name :", strings.Join(args, " "))
			return nil
		} else if len(response.Filters) == 0 {
			utils.Log.Info("No WASM Filter to display")
			return nil
		}

		tokenObj, err := utils.ReadToken(utils.TokenFlag)
		if err != nil {
			return errors.New(utils.FilterListError("error reading token\nUse 'mesheryctl filter list --help' to display usage guide\n" + err.Error()))
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
				utils.PrintToTableWithFooter([]string{"FILTER ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", ""})
				return nil
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
			utils.PrintToTableWithFooter([]string{"FILTER ID", "USER ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", "", ""})
			return nil
		}

		// Check if meshery provider is set
		if provider == "None" {
			for _, v := range response.Filters {
				FilterName := strings.Trim(v.Name, filepath.Ext(v.Name))
				FilterID := utils.TruncateID(v.ID.String())
				CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
				UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
				data = append(data, []string{FilterID, FilterName, CreatedAt, UpdatedAt})
			}
			utils.PrintToTableWithFooter([]string{"FILTER ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", ""})
			return nil
		}
		for _, v := range response.Filters {
			FilterID := utils.TruncateID(v.ID.String())
			var UserID string
			if v.UserID != nil {
				UserID = utils.TruncateID(*v.UserID)
			} else {
				UserID = "null"
			}
			FilterName := v.Name
			CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
			UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
			data = append(data, []string{FilterID, UserID, FilterName, CreatedAt, UpdatedAt})
		}
		utils.PrintToTableWithFooter([]string{"FILTER ID", "USER ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", "", ""})
		return nil
	},
}

// Pagination(making multiple requests) to retrieve filter Data in batches
func fetchFilters(baseURL, searchString string, pageSize, pageNumber int) (*models.FiltersAPIResponse, error) {
	var response *models.FiltersAPIResponse

	url := baseURL + "/api/filter"

	url = fmt.Sprintf("%s?pagesize=%d&page=%d", url, pageSize, pageNumber)
	if searchString != "" {
		url = url + "&search=" + searchString
	}

	utils.Log.Debug(url)

	req, err := utils.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, utils.FilterListError("failed to read response body"))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Server returned with status code: " + fmt.Sprint(resp.StatusCode) + "\n" + "Response: " + string(body))
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, ErrUnmarshal(err)
	}
	return response, nil
}
func init() {
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display full length user and filter file identifiers")
	listCmd.Flags().IntVarP(&pageNumber, "page", "p", 1, "(optional) List next set of filters with --page (default = 1)")
}
