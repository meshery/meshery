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

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type searchRelationshipFlags struct {
	Kind    string `json:"kind" validate:"omitempty"`
	SubType string `json:"subtype" validate:"omitempty"`
	Model   string `json:"model" validate:"omitempty"`
	Type    string `json:"type" validate:"omitempty"`
}

var (
	searchRelationshipFlagsProvided searchRelationshipFlags
)

// represents the mesheryctl exp relationship search subcommand.
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search registered relationship(s)",
	Long:  "Search registered relationship(s) used by different models",
	Example: `
// Search for a specific relationship
mesheryctl exp relationship search [--kind <kind>] [--type <type>] [--subtype <subtype>] [--model <model>]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flagValidator, ok := cmd.Context().Value(mesheryctlflags.FlagValidatorKey).(*mesheryctlflags.FlagValidator)
		if !ok || flagValidator == nil {
			return utils.ErrCommandContextMissing("flags-validator")
		}
		err := flagValidator.Validate(searchRelationshipFlagsProvided)
		if err != nil {
			return utils.ErrFlagsInvalid(err)
		}

		if searchRelationshipFlagsProvided.Kind == "" && searchRelationshipFlagsProvided.SubType == "" && searchRelationshipFlagsProvided.Type == "" && searchRelationshipFlagsProvided.Model == "" {
			return utils.ErrFlagsInvalid(fmt.Errorf("at least one of [--kind, --subtype, --type, --model] is required\n\nUsage: mesheryctl exp relationship search [--kind <kind>] [--type <type>] [--subtype <subtype>] [--model <model>]\nRun 'mesheryctl exp relationship search --help'"))
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
	searchCmd.Flags().StringVarP(&searchRelationshipFlagsProvided.Kind, "kind", "k", "", "search particular kind of relationships")
	searchCmd.Flags().StringVarP(&searchRelationshipFlagsProvided.SubType, "subtype", "s", "", "search particular subtype of relationships")
	searchCmd.Flags().StringVarP(&searchRelationshipFlagsProvided.Model, "model", "m", "", "search relationships of particular model name")
	searchCmd.Flags().StringVarP(&searchRelationshipFlagsProvided.Type, "type", "t", "", "search particular type of relationships")
}

func buildSearchUrl() string {
	var searchUrl strings.Builder

	if searchRelationshipFlagsProvided.Model == "" {
		searchUrl.WriteString("api/meshmodels/relationships?")
	} else {
		escapeModelName := url.QueryEscape(searchRelationshipFlagsProvided.Model)
		searchUrl.WriteString(fmt.Sprintf("api/meshmodels/models/%s/relationships?", escapeModelName))
	}

	if searchRelationshipFlagsProvided.Type != "" {
		escapedType := url.QueryEscape(searchRelationshipFlagsProvided.Type)
		searchUrl.WriteString(fmt.Sprintf("type=%s&", escapedType))
	}

	if searchRelationshipFlagsProvided.Kind != "" {
		escapeKind := url.QueryEscape(searchRelationshipFlagsProvided.Kind)
		searchUrl.WriteString(fmt.Sprintf("kind=%s&", escapeKind))
	}

	if searchRelationshipFlagsProvided.SubType != "" {
		escapeSubType := url.QueryEscape(searchRelationshipFlagsProvided.SubType)
		searchUrl.WriteString(fmt.Sprintf("subType=%s&", escapeSubType))
	}

	searchUrl.WriteString("pagesize=all")

	return searchUrl.String()
}
