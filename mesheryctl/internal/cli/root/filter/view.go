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
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/server/models"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/viewoutput"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type filterViewFlags struct {
	ViewAllFlag  bool   `json:"all" validate:"boolean"`
	OutputFormat string `json:"output-format" validate:"required,oneof=json yaml"`
	Save         bool   `json:"save" validate:"boolean"`
}

var filterViewFlagsProvided filterViewFlags

// filterNameSearchPageSize is the page size used when a filter name is looked
// up. The server's search matches substrings, so the lookup needs enough
// results to see the filter that carries the whole name.
const filterNameSearchPageSize = 10000

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View filter(s)",
	Long: `Displays the contents of a specific filter based on name or id.
	Find more information at: https://docs.meshery.io/reference/mesheryctl/filter/view`,
	Example: `
// View the specified WASM filter
// A unique prefix of the name or ID can also be provided. If the prefix is not unique, the first match will be returned.
mesheryctl filter view "[filter-name | ID]"

// View all filter files
mesheryctl filter view --all

// View all filter files in json
mesheryctl filter view --all --output-format json

// View all filter files in json and save it to a file
mesheryctl filter view --all --output-format json -s

// View multi-word named filter files. Multi-word filter names should be enclosed in quotes
mesheryctl filter view "filter name"
        `,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl filter view [filter-name | ID]\nRun 'mesheryctl filter view --help' to see detailed help message"
		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("accepts at most 1 arg, received %d\n\n%s", len(args), errMsg))
		}
		if len(args) > 0 && filterViewFlagsProvided.ViewAllFlag {
			return ErrViewAllWithName("view")
		}
		if (len(args) == 0 || args[0] == "") && !filterViewFlagsProvided.ViewAllFlag {
			return utils.ErrInvalidNameOrID(errors.New(errFilterNameOrIDNotProvided))
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &filterViewFlagsProvided)
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		filter := ""
		isID := false
		if len(args) > 0 {
			filter, isID, err = utils.ValidId(mctlCfg.GetBaseMesheryURL(), args[0], "filter")
			if err != nil {
				return utils.ErrInvalidNameOrID(err)
			}
		}

		urlString := ""
		if len(filter) == 0 {
			urlString = "api/filter?pagesize=10000"
		} else if isID {
			// if filter is a valid uuid, then directly fetch the filter
			urlString = "api/filter/" + filter
		} else {
			// else search filter by name
			urlString = fmt.Sprintf("api/filter?search=%s&pagesize=%d", url.QueryEscape(filter), filterNameSearchPageSize)
		}

		var selectedFilter *models.MesheryFilter
		var filterPage *models.MesheryFilterPage
		if isID {
			// Fetch single filter
			fetchedFilter, err := api.Fetch[models.MesheryFilter](urlString)
			if err != nil {
				return err
			}
			selectedFilter = fetchedFilter
		} else {
			page, err := api.Fetch[models.MesheryFilterPage](urlString)
			if err != nil {
				return err
			}
			filterPage = page

			if !filterViewFlagsProvided.ViewAllFlag {
				resolved, err := resolveFilterByName(filterPage, filter)
				if err != nil {
					return err
				}
				selectedFilter = resolved
			}
		}

		var data any

		if filterViewFlagsProvided.ViewAllFlag {
			data = filterPage
		} else {
			data = selectedFilter
		}

		savePath := ""

		if filterViewFlagsProvided.Save {
			var fileName string

			if filterViewFlagsProvided.ViewAllFlag {
				fileName = fmt.Sprintf(
					"filters_all.%s",
					strings.ToLower(filterViewFlagsProvided.OutputFormat),
				)
			} else {
				shortID := utils.TruncateID(utils.FormatUUID(selectedFilter.ID, "unknown"))

				sanitizedName := utils.SanitizeFileName(selectedFilter.Name)
				if sanitizedName == "" {
					sanitizedName = shortID
				}

				fileName = fmt.Sprintf(
					"filter_%s_%s.%s",
					sanitizedName,
					shortID,
					strings.ToLower(filterViewFlagsProvided.OutputFormat),
				)
			}

			savePath = filepath.Join(utils.MesheryFolder, fileName)
		}

		return viewoutput.HandleViewOutput(
			filterViewFlagsProvided.OutputFormat,
			data,
			cmd.OutOrStdout(),
			savePath,
		)

	},
}

func init() {
	viewCmd.Flags().BoolVarP(&filterViewFlagsProvided.ViewAllFlag, "all", "a", false, "(optional) view all filters available")
	viewCmd.Flags().StringVarP(&filterViewFlagsProvided.OutputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewCmd.Flags().BoolVarP(&filterViewFlagsProvided.Save, "save", "s", false, "(optional) save output as a JSON/YAML file")
}

// resolveFilterByName resolves a filter from filterPage by its whole name,
// prioritizing exact matches and falling back to case-insensitive matches.
// When multiple filters share the name, it prompts the user to select one.
func resolveFilterByName(filterPage *models.MesheryFilterPage, name string) (*models.MesheryFilter, error) {
	exact, folded := filtersNamed(filterPage.Filters, name)
	truncated := filterPage.TotalCount > len(filterPage.Filters)
	matches := exact
	if len(matches) == 0 {
		// A case-insensitive match is only safe when the result set is
		// complete; otherwise the exact match may be on a later page.
		if truncated {
			return nil, errFilterNameLookupTruncated(name)
		}
		matches = folded
	}

	switch len(matches) {
	case 0:
		return nil, ErrFilterNotFound(name, "view")
	case 1:
		return matches[0], nil
	default:
		return selectFilterPrompt(matches)
	}
}

// errFilterNameLookupTruncated reports a name lookup whose result set was too
// large to read in full, so the name cannot be resolved without guessing.
func errFilterNameLookupTruncated(filterName string) error {
	return utils.ErrInvalidArgument(fmt.Errorf("too many filters match %q to resolve the name; pass the filter ID instead", filterName))
}

// filtersNamed splits the filters into those named exactly name and those that
// differ only in case. The server's search is a case-insensitive substring
// match, so a filter whose name merely contains name is not a match: viewing
// it would display a different filter.
func filtersNamed(filters []*models.MesheryFilter, name string) (exact, folded []*models.MesheryFilter) {
	for _, f := range filters {
		if f == nil {
			continue
		}
		switch {
		case f.Name == name:
			exact = append(exact, f)
		case strings.EqualFold(f.Name, name):
			folded = append(folded, f)
		}
	}
	return exact, folded
}

func selectFilterPrompt(filters []*models.MesheryFilter) (*models.MesheryFilter, error) {
	columns := []string{"Filter Name", "Created At", "Updated At", "Filter ID"}
	widths := []int{24, 20, 20, 10}

	headingParts := make([]string, len(columns))
	for i, col := range columns {
		headingParts[i] = centerAlign(col, widths[i])
	}
	heading := strings.Join(headingParts, " ")

	filterInfos := make([]string, len(filters)+1)
	filterInfos[0] = heading
	for i, f := range filters {
		createdAt := utils.FormatTime(f.CreatedAt, "2006-01-02 15:04:05", "N/A")
		updatedAt := utils.FormatTime(f.UpdatedAt, "2006-01-02 15:04:05", "N/A")
		id := utils.FormatUUID(f.ID, "N/A")
		lastFourID := id
		if len(id) > 4 {
			lastFourID = id[len(id)-4:]
		}

		rowParts := []string{
			centerAlign(f.Name, widths[0]),
			centerAlign(createdAt, widths[1]),
			centerAlign(updatedAt, widths[2]),
			centerAlign(lastFourID, widths[3]),
		}
		filterInfos[i+1] = strings.Join(rowParts, " ")
	}

	prompt := promptui.Select{
		Label: "Select a filter",
		Items: filterInfos,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "{{ if eq . \"" + heading + "\" }}{{ . | bold }}{{ else }}{{ . | cyan }}{{ end }}",
			Inactive: "{{ if eq . \"" + heading + "\" }}{{ . | bold }}{{ else }}{{ . }}{{ end }}",
			Selected: "{{ if eq . \"" + heading + "\" }}{{ . | bold }}{{ else }}{{ . | green }}{{ end }}",
		},
		StartInSearchMode: false,
	}

	for {
		i, _, err := prompt.Run()
		if err != nil {
			return nil, utils.ErrPromptCancelled()
		}
		if i == 0 {
			continue // skip header row, re-prompt
		}
		return filters[i-1], nil
	}
}

func centerAlign(text string, width int) string {
	if len(text) >= width {
		return text
	}
	leftPadding := (width - len(text)) / 2
	rightPadding := width - len(text) - leftPadding
	return strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding)
}
