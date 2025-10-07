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
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
)

var (
	verbose  bool
	page     int
	pageSize int
	provider string
)

var linkDocPatternList = map[string]string{
	"link":    "![pattern-list-usage](/assets/img/mesheryctl/patternList.png)",
	"caption": "Usage of mesheryctl design list",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List designs",
	Long: `Display list of all available designs.
Documentation for design can be found at https://docs.meshery.io/reference/mesheryctl/design/list
`,
	Args: cobra.MinimumNArgs(0),
	Example: `
// Display a list of all available designs
mesheryctl design list

// Display a list of all available designs with verbose output
mesheryctl design list --verbose

// Display a list of all available designs with specified page number (10 designs per page by default)
mesheryctl design list --page [pange-number]

// Display a list of all available designs with custom page size (10 designs per page by default)
mesheryctl design list --pagesize [page-size]

// Display only the count of all available designs
mesheryctl design list --count
	`,
	Annotations: linkDocPatternList,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// token flag is not provided
		if utils.TokenFlag == "" {
			// retrieve token from default location
			tokenPath, err := utils.GetCurrentAuthToken()
			if err != nil {
				return err
			}
			utils.TokenFlag = tokenPath
		}

		tokenObj, err := utils.ReadToken(utils.TokenFlag)
		if err != nil {
			return err
		}
		provider = tokenObj["meshery-provider"]
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		header := func(provider string) []string {
			if provider == "None" {
				return []string{"DESIGN ID", "NAME", "CREATED", "UPDATED"}
			}
			return []string{"DESIGN ID", "USER ID", "NAME", "CREATED", "UPDATED"}
		}(provider)

		designData := display.DisplayDataAsync{
			UrlPath:          "api/pattern",
			DataType:         "design",
			Header:           header,
			Page:             page,
			PageSize:         pageSize,
			IsPage:           cmd.Flags().Changed("page"),
			DisplayCountOnly: cmd.Flags().Changed("count"),
		}

		return display.ListAsyncPagination(designData, processDesignData)
	},
}

func processDesignData(data *models.PatternsAPIResponse) ([][]string, int64) {
	var displayData [][]string
	for _, v := range data.Patterns {
		designId := func(id *uuid.UUID, isVerbose bool) string {
			if isVerbose {
				return id.String()
			}
			return utils.TruncateID(id.String())
		}(v.ID, verbose)

		designName := strings.Trim(v.Name, filepath.Ext(v.Name))
		createdAt := formatTimeToString(v.CreatedAt, verbose)
		updatedAt := formatTimeToString(v.UpdatedAt, verbose)

		if provider != "None" {
			userID := func(userID *string) string {
				if userID != nil {
					if verbose {
						return *userID
					}
					return utils.TruncateID(*userID)
				}
				return "null"
			}(v.UserID)
			displayData = append(displayData, []string{designId, userID, designName, createdAt, updatedAt})
		} else {
			displayData = append(displayData, []string{designId, designName, createdAt, updatedAt})
		}
	}
	return displayData, int64(data.TotalCount)
}

func formatTimeToString(t *time.Time, isVerbose bool) string {
	if isVerbose {
		return t.Format("01-02-2006 15:04:05") // M-D-YYYY HH:MM:SS
	}
	return t.Format("01-02-2006") // MM-DD-YYYY
}

func init() {
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "(optional) Display full length user and design file identifiers")
	listCmd.Flags().IntVarP(&page, "page", "p", 1, "(optional) List next set of designs with --page")
	listCmd.Flags().IntVarP(&pageSize, "pagesize", "", 10, "(optional) Number of designs to be displayed per page")
	listCmd.Flags().BoolP("count", "c", false, "(optional) Display count only")
}
