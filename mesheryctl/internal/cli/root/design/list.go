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

package design

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	verbose           bool
	pageNumber        int
	pageSize          = 25
	whiteBoardPrinter = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
)

var linkDocPatternList = map[string]string{
	"link":    "![pattern-list-usage](/assets/img/mesheryctl/patternList.png)",
	"caption": "Usage of mesheryctl design list",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List designs",
	Long:  `Display list of all available design files.`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// list all available designs
mesheryctl design list
	`,
	Annotations: linkDocPatternList,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/pattern?pagesize=%d&page=%d", baseUrl, pageSize, pageNumber)
		var response models.PatternsAPIResponse
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		res, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return nil
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			utils.Log.Error(utils.ErrUnmarshal(err))
			return nil
		}

		tokenObj, err := utils.ReadToken(utils.TokenFlag)
		if err != nil {
			utils.Log.Error(utils.ErrReadToken(err))
			return nil
		}

		provider := tokenObj["meshery-provider"]
		var data [][]string
		provider_header := []string{"DESIGN ID", "USER ID", "NAME", "CREATED", "UPDATED"}
		non_provider_header := []string{"DESIGN ID", "NAME", "CREATED", "UPDATED"}

		if verbose {
			if provider == "None" {
				for _, v := range response.Patterns {
					PatternID := v.ID.String()
					PatterName := v.Name
					CreatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year(), v.CreatedAt.Hour(), v.CreatedAt.Minute(), v.CreatedAt.Second())
					UpdatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year(), v.UpdatedAt.Hour(), v.UpdatedAt.Minute(), v.UpdatedAt.Second())
					data = append(data, []string{PatternID, PatterName, CreatedAt, UpdatedAt})
				}
				processData(cmd, data, non_provider_header, int64(response.TotalCount))
				return nil
			}

			for _, v := range response.Patterns {
				PatternID := utils.TruncateID(v.ID.String())
				var UserID string
				if v.UserID != nil {
					UserID = utils.TruncateID(*v.UserID)
				} else {
					UserID = "null"
				}
				PatterName := v.Name
				CreatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year(), v.CreatedAt.Hour(), v.CreatedAt.Minute(), v.CreatedAt.Second())
				UpdatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year(), v.UpdatedAt.Hour(), v.UpdatedAt.Minute(), v.UpdatedAt.Second())
				data = append(data, []string{PatternID, UserID, PatterName, CreatedAt, UpdatedAt})
			}
			processData(cmd, data, provider_header, int64(response.TotalCount))

			return nil
		}

		// Check if messhery provider is set
		if provider == "None" {
			for _, v := range response.Patterns {
				PatterName := strings.Trim(v.Name, filepath.Ext(v.Name))
				PatternID := utils.TruncateID(v.ID.String())
				CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
				UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
				data = append(data, []string{PatternID, PatterName, CreatedAt, UpdatedAt})
			}
			processData(cmd, data, non_provider_header, int64(response.TotalCount))
			return nil
		}

		for _, v := range response.Patterns {
			PatternID := utils.TruncateID(v.ID.String())
			var UserID string
			if v.UserID != nil {
				UserID = utils.TruncateID(*v.UserID)
			} else {
				UserID = "null"
			}
			PatterName := v.Name
			CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
			UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
			data = append(data, []string{PatternID, UserID, PatterName, CreatedAt, UpdatedAt})
		}
		processData(cmd, data, provider_header, int64(response.TotalCount))

		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display full length user and design file identifiers")
	listCmd.Flags().IntVarP(&pageNumber, "page", "p", 1, "(optional) List next set of designs with --page (default = 1)")
}

func processData(cmd *cobra.Command, data [][]string, header []string, totalCount int64) {
	if len(data) == 0 {
		whiteBoardPrinter.Println("No pattern(s) found")
		return
	}
	utils.DisplayCount("patterns", totalCount)
	if cmd.Flags().Changed("page") {
		utils.PrintToTable(header, data)
	} else {
		err := utils.HandlePagination(pageSize, "patterns", data, header)
		if err != nil {
			utils.Log.Error(err)
		}
	}
}
