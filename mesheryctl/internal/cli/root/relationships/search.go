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

type relationshipSearchFlags struct {
	SearchKind      string
	SearchSubType   string
	SearchType      string
	SearchModelName string
}

var searchFlags relationshipSearchFlags

// represents the mesheryctl exp relationship search [query-text] subcommand.
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search registered relationship(s)",
	Long:  "Search registred relationship(s) used by different models",
	Example: `
// Search for relationship using a query
mesheryctl exp relationship search [--kind <kind>] [--type <type>] [--subtype <subtype>] [--model <model>] [query-text]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flagValidator, ok := cmd.Context().Value(mesheryctlflags.FlagValidatorKey).(*mesheryctlflags.FlagValidator)
		if !ok || flagValidator == nil {
			return utils.ErrCommandContextMissing("flags-validator")
		}
		err := flagValidator.Validate(searchFlags)
		if err != nil {
			return utils.ErrFlagsInvalid(err)
		}
		const usage = "mesheryctl exp relationship search [--kind <kind>] [--type <type>] [--subtype <subtype>] [--model <model>]"
		errMsg := fmt.Errorf("at least one of [--kind, --subtype, --type, --model] is required\n\nUsage: %s\nRun 'mesheryctl exp relationship search --help'", usage)
		if searchFlags.SearchKind == "" && searchFlags.SearchSubType == "" && searchFlags.SearchType == "" && searchFlags.SearchModelName == "" {
			return utils.ErrInvalidArgument(errMsg)
		}
		return nil

	},
	RunE: func(cmd *cobra.Command, args []string) error {
		dataToDisplay := display.DisplayDataAsync{
			UrlPath:          buildSearchUrl(),
			DataType:         "relationship",
			Header:           []string{"kind", "apiVersion", "model-name", "subType", "regoQuery"},
			Page:             1,
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
	searchCmd.Flags().StringVarP(&searchFlags.SearchKind, "kind", "k", "", "Search for a specific kind of relationship")
	searchCmd.Flags().StringVarP(&searchFlags.SearchType, "type", "t", "", "Search for a specific type of relationship")
	searchCmd.Flags().StringVarP(&searchFlags.SearchSubType, "subtype", "s", "", "Search for a specific subtype of relationship")
	searchCmd.Flags().StringVarP(&searchFlags.SearchModelName, "model", "m", "", "Search for a specific model of relationship")

}

func buildSearchUrl() string {
	var searchUrl strings.Builder

	if searchFlags.SearchModelName == "" {
		searchUrl.WriteString("api/meshmodels/relationships?")
	} else {
		escapeModelName := url.QueryEscape(searchFlags.SearchModelName)
		searchUrl.WriteString(fmt.Sprintf("api/meshmodels/models/%s/relationships?", escapeModelName))
	}

	if searchFlags.SearchType != "" {
		escapedType := url.QueryEscape(searchFlags.SearchType)
		searchUrl.WriteString(fmt.Sprintf("type=%s&", escapedType))
	}

	if searchFlags.SearchKind != "" {
		escapeKind := url.QueryEscape(searchFlags.SearchKind)
		searchUrl.WriteString(fmt.Sprintf("kind=%s&", escapeKind))
	}

	if searchFlags.SearchSubType != "" {
		escapeSubType := url.QueryEscape(searchFlags.SearchSubType)
		searchUrl.WriteString(fmt.Sprintf("subType=%s&", escapeSubType))
	}

	searchUrl.WriteString("pagesize=all")

	return searchUrl.String()
}
