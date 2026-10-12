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
	"io"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/schemas/models/v1beta2/organization"
	"github.com/spf13/cobra"
)

type cmdOrgListFlags struct {
	OutputFormat string `json:"output-format" validate:"omitempty,oneof=json yaml"`
}

var orgListFlags cmdOrgListFlags

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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &orgListFlags)
	},
	Args: func(cmd *cobra.Command, _ []string) error {
		if orgListFlags.OutputFormat == "" {
			// Default (unset) means "render the interactive table" - the
			// --count conflict below only applies once a structured format
			// was actually requested.
			return nil
		}
		if count, _ := cmd.Flags().GetBool("count"); count {
			// --count's only defined meaning is "print the total for the
			// table view"; it has no meaning once the output is a structured
			// document, so reject the combination rather than silently
			// dropping one of the two flags.
			return ErrCountWithOutputFormat()
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		pagesize, _ := cmd.Flags().GetInt("pagesize")
		count, _ := cmd.Flags().GetBool("count")

		if orgListFlags.OutputFormat != "" {
			return listOrgsAsStructuredOutput(cmd.OutOrStdout(), orgListFlags.OutputFormat, page, pagesize)
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
// renders it as JSON or YAML to out, instead of the interactive table used by
// the default `organization list` path. It builds the same zero-based
// page/pagesize query as display.HandlePaginationAsync (via the shared
// display.NormalizePagination helper), so --page/--pagesize mean the same
// thing regardless of --output-format.
//
// The page type comes from the schemas package (schema-driven), which already
// carries matching json and yaml tags on every field.
func listOrgsAsStructuredOutput(out io.Writer, outputFormat string, page, pagesize int) error {
	currentPage, effectivePageSize := display.NormalizePagination(page, pagesize)

	urlPath := fmt.Sprintf("%s?page=%d&pagesize=%d", organizationsApiPath, currentPage, effectivePageSize)
	orgsData, err := api.Fetch[organization.OrganizationsPage](urlPath)
	if err != nil {
		return err
	}

	outputFormatterFactory := display.OutputFormatterFactory[*organization.OrganizationsPage]{}
	outputFormatter, err := outputFormatterFactory.New(outputFormat, orgsData)
	if err != nil {
		return err
	}

	return outputFormatter.WithOutput(out).Display()
}

func processOrgData(orgs *organization.OrganizationsPage) ([][]string, int64) {
	var rows [][]string
	if orgs.Organizations != nil {
		for _, org := range *orgs.Organizations {
			id := ""
			if org.ID != nil {
				id = org.ID.String()
			}
			rows = append(rows, []string{
				id,
				org.Name,
				org.CreatedAt.Format("01-02-2006"),
			})
		}
	}

	var total int64
	if orgs.TotalCount != nil {
		total = int64(*orgs.TotalCount)
	}
	return rows, total
}

func init() {
	listOrgCmd.Flags().IntP("page", "p", 1, "(optional) Page number of paginated results")
	listOrgCmd.Flags().IntP("pagesize", "s", 10, "(optional) Number of organizations per page")
	listOrgCmd.Flags().BoolP("count", "", false, "total number of registered orgs")
	listOrgCmd.Flags().StringVarP(&orgListFlags.OutputFormat, "output-format", "o", "", "(optional) format to display in [json|yaml]. Defaults to a human-readable table")
}
