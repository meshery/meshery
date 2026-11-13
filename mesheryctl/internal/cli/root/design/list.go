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
	"fmt"
	"time"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
)

const listUsageMsg = "Usage: mesheryctl design list\nRun 'mesheryctl design list --help' to see detailed help message"

type designListFlags struct {
	Page     int `json:"page" validate:"gte=1"`
	PageSize int `json:"pagesize" validate:"gte=1"`
}

var (
	verbose              bool
	designListFlagsInput designListFlags
	provider             string
)

var linkDocPatternList = map[string]string{
	"link":    "![pattern-list-usage](../../../images/patternList.png)",
	"caption": "Usage of mesheryctl design list",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List designs",
	Long: `Display list of all available designs.
Find more information at: https://docs.meshery.io/reference/mesheryctl/design/list
`,
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", listUsageMsg))
		}
		return nil
	},
	Example: `
// Display a list of all available designs
mesheryctl design list

// Display a list of all available designs with verbose output
mesheryctl design list --verbose

// Display a list of all available designs with specified page number (10 designs per page by default)
mesheryctl design list --page [page-number]

// Display a list of all available designs with custom page size (10 designs per page by default)
mesheryctl design list --pagesize [page-size]

// Display only the count of all available designs
mesheryctl design list --count
	`,
	Annotations: linkDocPatternList,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := mesheryctlflags.ValidateCmdFlags(cmd, &designListFlagsInput); err != nil {
			return err
		}

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
			if utils.IsLocalProvider(provider) {
				return []string{"DESIGN ID", "NAME", "CREATED", "UPDATED"}
			}
			return []string{"DESIGN ID", "OWNER", "NAME", "CREATED", "UPDATED"}
		}(provider)

		designData := display.DisplayDataAsync{
			UrlPath:          "api/pattern",
			DataType:         "design",
			Header:           header,
			Page:             designListFlagsInput.Page,
			PageSize:         designListFlagsInput.PageSize,
			IsPage:           cmd.Flags().Changed("page"),
			DisplayCountOnly: cmd.Flags().Changed("count"),
		}

		return display.ListAsyncPagination(designData, processDesignData)
	},
}

func processDesignData(data *models.PatternsAPIResponse) ([][]string, int64) {
	if data == nil {
		return nil, 0
	}
	var displayData [][]string
	for _, v := range data.Patterns {
		designId := utils.FormatUUID(v.ID, "N/A")

		designName := utils.TrimFileExtension(v.Name)
		createdAt := formatTimeToString(v.CreatedAt, verbose)
		updatedAt := formatTimeToString(v.UpdatedAt, verbose)

		if !utils.IsLocalProvider(provider) {
			owner := utils.FormatUUID(v.UserID, "null")
			if !verbose {
				owner = utils.TruncateID(owner)
			}
			displayData = append(displayData, []string{designId, owner, designName, createdAt, updatedAt})
		} else {
			displayData = append(displayData, []string{designId, designName, createdAt, updatedAt})
		}
	}
	return displayData, int64(data.TotalCount)
}

func formatTimeToString(t *time.Time, isVerbose bool) string {
	if isVerbose {
		return utils.FormatTime(t, "01-02-2006 15:04:05", "N/A")
	}
	return utils.FormatTime(t, "01-02-2006", "N/A")
}

func init() {
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "(optional) Display full length owner identifiers and detailed timestamps")
	listCmd.Flags().IntVarP(&designListFlagsInput.Page, "page", "p", 1, "(optional) List next set of designs with --page")
	listCmd.Flags().IntVarP(&designListFlagsInput.PageSize, "pagesize", "", 10, "(optional) Number of designs to be displayed per page")
	listCmd.Flags().BoolP("count", "c", false, "(optional) Display count only")
}
