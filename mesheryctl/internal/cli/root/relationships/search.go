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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	searchModelName string
	searchType      string
	searchSubType   string
	searchKind      string
)

// represents the mesheryctl exp relationship search [query-text] subcommand.
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search registered relationship(s)",
	Long:  "Search registred relationship(s) used by different models",
	Example: `
// Search for relationship using a query
mesheryctl exp relationship search [--kind <kind>] [--type <type>] [--subtype <subtype>] [--model <model>] [query-text]`,
	Args: func(cmd *cobra.Command, args []string) error {
		const usage = "mesheryctl exp relationship search [--kind <kind>] [--type <type>] [--subtype <subtype>] [--model <model>]"
		errMsg := fmt.Errorf("[--kind, --subtype or --type or --model] and [query-text] are required\n\nUsage: %s\nRun 'mesheryctl exp relationship search --help'", usage)

		if searchKind == "" && searchSubType == "" && searchType == "" && searchModelName == "" {
			err := utils.ErrInvalidArgument(errMsg)
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		relationshipResponse, err := api.Fetch[MeshmodelRelationshipsAPIResponse](buildSearchUrl())

		if err != nil {
			return err
		}
		rows := [][]string{}

		for _, relationship := range relationshipResponse.Relationships {
			if len(relationship.Type()) > 0 {
				evaluationQuery := ""
				if relationship.EvaluationQuery != nil {
					evaluationQuery = *relationship.EvaluationQuery
				}
				rows = append(rows, []string{string(relationship.Kind), relationship.SchemaVersion, relationship.Model.DisplayName, relationship.SubType, evaluationQuery})
			}
		}

		dataToDisplay := display.DisplayedData{
			DataType:         "relationship",
			Header:           []string{"kind", "apiVersion", "model-name", "subType", "regoQuery"},
			Rows:             rows,
			Count:            relationshipResponse.Count,
			DisplayCountOnly: false,
			IsPage:           cmd.Flags().Changed("page"),
		}

		return display.List(dataToDisplay)
	},
}

func init() {
	searchCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})
	searchCmd.Flags().StringVarP(&searchKind, "kind", "k", "", "search particular kind of relationships")
	searchCmd.Flags().StringVarP(&searchSubType, "subtype", "s", "", "search particular subtype of relationships")
	searchCmd.Flags().StringVarP(&searchModelName, "model", "m", "", "search relationships of particular model name")
	searchCmd.Flags().StringVarP(&searchType, "type", "t", "", "search particular type of relationships")
}

func buildSearchUrl() string {
	var searchUrl strings.Builder

	if searchModelName == "" {
		searchUrl.WriteString("api/meshmodels/relationships?")
	} else {
		escapeModelName := url.QueryEscape(searchModelName)
		searchUrl.WriteString(fmt.Sprintf("api/meshmodels/models/%s/relationships?", escapeModelName))
	}

	if searchType != "" {
		escapedType := url.QueryEscape(searchType)
		searchUrl.WriteString(fmt.Sprintf("type=%s&", escapedType))
	}

	if searchKind != "" {
		escapeKind := url.QueryEscape(searchKind)
		searchUrl.WriteString(fmt.Sprintf("kind=%s&", escapeKind))
	}

	if searchSubType != "" {
		escapeSubType := url.QueryEscape(searchSubType)
		searchUrl.WriteString(fmt.Sprintf("subType=%s&", escapeSubType))
	}

	searchUrl.WriteString("pagesize=all")

	return searchUrl.String()
}
