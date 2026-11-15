// Copyright Meshery Authors
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
	"strings"
	"time"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type filterListFlags struct {
	Count    bool `json:"count" validate:"boolean"`
	Page     int  `json:"page" validate:"gte=1"`
	PageSize int  `json:"pagesize" validate:"gte=1"`
	Verbose  bool `json:"verbose" validate:"boolean"`
}

var filterListFlagsProvided filterListFlags

var listUsageMsg = "Usage: mesheryctl filter list [filter-name]\nRun 'mesheryctl filter list --help' to see detailed help message"

var listCmd = &cobra.Command{
	Use:   "list [filter-name]",
	Short: "List filters",
	Long: `Display list of all available filter files.
	Find more information at: https://docs.meshery.io/reference/mesheryctl/filter/list`,
	Example: `
// List all WASM filter files present (maximum 25 filters)
mesheryctl filter list

// Search for filter
mesheryctl filter list Test

// Search for filter with space
mesheryctl filter list 'Test Filter'
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &filterListFlagsProvided)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", listUsageMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		var searchString string
		if len(args) > 0 {
			searchString = strings.ReplaceAll(args[0], " ", "%20")
		}

		response, err := fetchFilters(mctlCfg.GetBaseMesheryURL(), searchString, filterListFlagsProvided.PageSize, filterListFlagsProvided.Page-1)
		if err != nil {
			utils.Log.Error(err)
			return ErrFetchFilter(err)
		}

		if filterListFlagsProvided.Count {
			utils.DisplayCount("filter", int64(response.TotalCount))
			return nil
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
			return utils.ErrReadToken(err)
		}
		provider := tokenObj["meshery-provider"]
		header, data, footer := processFilterData(response, filterListFlagsProvided.Verbose, utils.IsLocalProvider(provider))

		if cmd.Flags().Changed("page") {
			utils.PrintToTable(header, data, footer)
			return nil
		}
		err = utils.HandlePagination(filterListFlagsProvided.PageSize, "filter files", data, header, footer)
		if err != nil {
			return utils.ErrHandlePagination(err)
		}
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

	req, err := utils.NewRequest("GET", url, nil)
	if err != nil {
		return nil, utils.ErrCreatingRequest(err)
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, utils.ErrRequestResponse(err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, utils.FilterListError("failed to read response body"))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Server returned with status code: " + fmt.Sprint(resp.StatusCode) + "\n" + "Response: " + string(body))
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, utils.ErrUnmarshal(err)
	}
	return response, nil
}

func init() {
	listCmd.Flags().BoolVarP(&filterListFlagsProvided.Count, "count", "c", false, "(optional) Display count only")
	listCmd.Flags().IntVarP(&filterListFlagsProvided.Page, "page", "p", 1, "(optional) List next set of filters with --page (default = 1)")
	listCmd.Flags().IntVarP(&filterListFlagsProvided.PageSize, "pagesize", "s", 25, "(optional) List next set of filters with --pagesize (default = 25)")
	listCmd.Flags().BoolVarP(&filterListFlagsProvided.Verbose, "verbose", "v", false, "Display full length user and filter file identifiers")
}

func processFilterData(response *models.FiltersAPIResponse, isVerbose, isLocal bool) ([]string, [][]string, []string) {
	if response == nil {
		return nil, nil, nil
	}

	var data [][]string
	for _, v := range response.Filters {
		filterID := utils.FormatUUID(v.ID, "N/A")
		if !isVerbose || !isLocal {
			filterID = utils.TruncateID(filterID)
		}
		filterName := utils.TrimFileExtension(v.Name)
		createdAt := formatFilterTime(v.CreatedAt, isVerbose)
		updatedAt := formatFilterTime(v.UpdatedAt, isVerbose)

		if isLocal {
			data = append(data, []string{filterID, filterName, createdAt, updatedAt})
			continue
		}

		owner := "null"
		if v.Owner != nil {
			owner = *v.Owner
			if !isVerbose {
				owner = utils.TruncateID(owner)
			}
		}
		data = append(data, []string{filterID, owner, filterName, createdAt, updatedAt})
	}

	if isLocal {
		header := []string{"FILTER ID", "NAME", "CREATED", "UPDATED"}
		footer := []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", ""}
		return header, data, footer
	}

	header := []string{"FILTER ID", "OWNER", "NAME", "CREATED", "UPDATED"}
	footer := []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", "", ""}
	return header, data, footer
}

func formatFilterTime(t *time.Time, isVerbose bool) string {
	if t == nil {
		return "N/A"
	}
	if isVerbose {
		return fmt.Sprintf("%d-%d-%d %d:%d:%d", int(t.Month()), t.Day(), t.Year(), t.Hour(), t.Minute(), t.Second())
	}
	return fmt.Sprintf("%d-%d-%d", int(t.Month()), t.Day(), t.Year())
}
