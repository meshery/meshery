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

package organizations

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
)

var listOrgCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered organizations",
	Long: `List all registered organizations with their id, name and date of creation. Organization will be returned based on provider you logged in.
Find more information at: https://docs.meshery.io/reference/references/mesheryctl/organization/list`,
	Example: `
// list all organizations
mesheryctl organization list

// list organizations for a specified page
mesheryctl organization list --page [page-number]

// Display number of available organizations
mesheryctl organization list --count

// list organizations as structured output, for use in scripts
mesheryctl organization list --output-format json
mesheryctl organization list --output-format yaml
	`,
	Args: func(cmd *cobra.Command, _ []string) error {
		outputFormat, _ := cmd.Flags().GetString("output-format")
		if outputFormat == "" {
			// Default (unset) means "render the interactive table" - only
			// validate when the caller actually asked for a machine-readable
			// format.
			return nil
		}
		return display.ValidateOutputFormat(outputFormat)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		pagesize, _ := cmd.Flags().GetInt("pagesize")
		count, _ := cmd.Flags().GetBool("count")
		outputFormat, _ := cmd.Flags().GetString("output-format")

		if outputFormat != "" {
			return listOrgsAsStructuredOutput(outputFormat, page, pagesize)
		}

		data := display.DisplayDataAsync{
			UrlPath:          organizationsApiPath,
			Page:             page,
			PageSize:         pagesize,
			DataType:         "organizations",
			Header:           []string{"ID", "NAME", "CREATED-AT"},
			DisplayCountOnly: count,
			IsPage:           cmd.Flags().Changed("page"),
		}

		return display.ListAsyncPagination(data, processOrgData)
	},
}

// listOrgsAsStructuredOutput fetches a single page of organizations and
// renders it as JSON or YAML, instead of the interactive table used by the
// default `organization list` path. It builds the same zero-based
// page/pagesize query as display.HandlePaginationAsync so `--page`/
// `--pagesize` mean the same thing regardless of `--output-format`.
func listOrgsAsStructuredOutput(outputFormat string, page, pagesize int) error {
	currentPage := page - 1
	if currentPage < 0 {
		currentPage = 0
	}
	effectivePageSize := pagesize
	if effectivePageSize <= 0 {
		effectivePageSize = 10
	}

	urlPath := fmt.Sprintf("%s?page=%d&pagesize=%d", organizationsApiPath, currentPage, effectivePageSize)
	orgsData, err := api.Fetch[models.OrganizationsPage](urlPath)
	if err != nil {
		return err
	}

	outputFormatterFactory := display.OutputFormatterFactory[*models.OrganizationsPage]{}
	outputFormatter, err := outputFormatterFactory.New(outputFormat, orgsData)
	if err != nil {
		return err
	}

	return outputFormatter.Display()
}

func processOrgData(orgs *models.OrganizationsPage) ([][]string, int64) {
	var rows [][]string
	for _, org := range orgs.Organizations {
		rows = append(rows, []string{
			org.ID.String(),
			org.Name,
			org.CreatedAt.Format("01-02-2006"),
		})

	}
	return rows, int64(orgs.TotalCount)
}

func init() {
	listOrgCmd.Flags().IntP("page", "p", 1, "(optional) Page number of paginated results")
	listOrgCmd.Flags().IntP("pagesize", "s", 10, "(optional) Number of organizations per page")
	listOrgCmd.Flags().BoolP("count", "", false, "total number of registered orgs")
	listOrgCmd.Flags().StringP("output-format", "o", "", "(optional) format to display in [json|yaml]. Defaults to a human-readable table")
}
