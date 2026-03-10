// Copyright 2024 Meshery Authors
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

package relationships

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type searchRelationshipFlags struct {
	Kind    string `json:"kind"     validate:"omitempty"`
	SubType string `json:"subtype"  validate:"omitempty"`
	Model   string `json:"model"    validate:"omitempty"`
	Type    string `json:"type"     validate:"omitempty"`
	Page    int    `json:"page"     validate:"omitempty,gte=1"`
}

var (
	searchRelationshipFlagsProvided searchRelationshipFlags
)

// represents the mesheryctl relationship search subcommand.
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search registered relationship(s)",
	Long:  "Search registered relationship(s) used by different models",
	Example: `
// Search for a specific relationship
mesheryctl relationship search [--kind <kind>] [--type <type>] [--subtype <subtype>] [--model <model>]

// Search a relationship for a specified page
mesheryctl relationship search [--kind <kind>] [--page <int>]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := mesheryctlflags.ValidateCmdFlags(cmd, &searchRelationshipFlagsProvided)
		if err != nil {
			return err
		}

		if searchRelationshipFlagsProvided.Kind == "" &&
			searchRelationshipFlagsProvided.SubType == "" &&
			searchRelationshipFlagsProvided.Type == "" &&
			searchRelationshipFlagsProvided.Model == "" {
			return utils.ErrFlagsInvalid(fmt.Errorf(
				"at least one of [--kind, --subtype, --type, --model] is required\n\n" +
					"Usage: mesheryctl exp relationship search [--kind <kind>] [--type <type>] [--subtype <subtype>] [--model <model>]\n" +
					"Run 'mesheryctl exp relationship search --help'",
			))
		}

		return nil
	},

	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 0 {
			errMsg := "Usage: mesheryctl exp relationship search\nRun 'mesheryctl exp relationship search --help' to see detailed help message"
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		dataToDisplay := display.DisplayDataAsync{
			UrlPath:          buildSearchUrl(),
			DataType:         "relationship",
			Header:           []string{"ID", "Kind", "API Version", "Model Name", "Sub Type", "Type"},
			Page:             searchRelationshipFlagsProvided.Page,
			PageSize:         10,
			DisplayCountOnly: false,
			IsPage:           cmd.Flags().Changed("page"),
		}
		return display.ListAsyncPagination(dataToDisplay, generateRelationshipDataToDisplay)
	},
}

func init() {
	searchCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})

	searchCmd.Flags().StringVarP(&searchRelationshipFlagsProvided.Kind, "kind", "k", "", "(optional) Search relationships of a particular kind")
	searchCmd.Flags().StringVarP(&searchRelationshipFlagsProvided.SubType, "subtype", "s", "", "(optional) Search relationships of a particular subtype")
	searchCmd.Flags().StringVarP(&searchRelationshipFlagsProvided.Model, "model", "m", "", "(optional) Search relationships of a particular model name")
	searchCmd.Flags().StringVarP(&searchRelationshipFlagsProvided.Type, "type", "t", "", "(optional) Search relationships of a particular type")
	searchCmd.Flags().IntVarP(&searchRelationshipFlagsProvided.Page, "page", "p", 1, "(optional) Page number of results to fetch (default = 1)")

}

// buildSearchUrl constructs the API URL for the relationship search endpoint.
func buildSearchUrl() string {
	base := "api/meshmodels/relationships"

	if searchRelationshipFlagsProvided.Model != "" {
		base = fmt.Sprintf("api/meshmodels/models/%s/relationships",
			url.PathEscape(searchRelationshipFlagsProvided.Model))
	}

	params := url.Values{}

	if searchRelationshipFlagsProvided.Type != "" {
		params.Set("type", searchRelationshipFlagsProvided.Type)
	}
	if searchRelationshipFlagsProvided.Kind != "" {
		params.Set("kind", searchRelationshipFlagsProvided.Kind)
	}
	if searchRelationshipFlagsProvided.SubType != "" {
		params.Set("subType", searchRelationshipFlagsProvided.SubType)
	}

	return base + "?" + params.Encode()
}
