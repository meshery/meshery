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

package relationships

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/mesheryctl/pkg/utils/format"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view relationships of a model by its name",
	Long:  "view a relationship queried by the model name",
	Example: `
// View relationships of a model
mesheryctl exp relationship view [model-name]
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "\n\nUsage: mesheryctl exp relationship view [model-name]\nRun 'mesheryctl exp relationship view --help' to see detailed help message"
		if len(args) == 0 {
			return errors.New(utils.RelationshipsError(fmt.Sprintf("[model-name] isn't specified%s", errMsg), "view"))
		} else if len(args) > 1 {
			return errors.New(utils.RelationshipsError(fmt.Sprintf("Too many arguments only [model-name] is expected%s", errMsg), "view"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		model := args[0]

		relationshipsResponse, err := api.Fetch[MeshmodelRelationshipsAPIResponse](fmt.Sprintf("api/meshmodels/models/%s/relationships?pagesize=all", model))
		if err != nil {
			return err
		}

		var selectedModel *relationship.RelationshipDefinition

		if relationshipsResponse.Count == 0 {
			utils.Log.Info("No relationship(s) found for the given name ", model)
			return nil
		} else if relationshipsResponse.Count == 1 {
			selectedModel = &relationshipsResponse.Relationships[0]
		} else {
			selectedModel = selectRelationshipPrompt(relationshipsResponse.Relationships)
		}
		var output []byte

		// user may pass flag in lower or upper case but we have to keep it lower
		// in order to make it consistent while checking output format
		outpoutFormat, _ := cmd.Flags().GetString("output-format")
		outFormatFlag := strings.ToLower(outpoutFormat)
		if outFormatFlag == "yaml" || outFormatFlag == "yml" {
			if output, err = yaml.Marshal(selectedModel); err != nil {
				return errors.Wrap(err, "failed to format output in YAML")
			}
			utils.Log.Info(string(output))
		} else if outFormatFlag == "json" {
			// return outputRelationshipJson(selectedModel)
			return format.OutputJson(selectedModel)
		} else {
			return errors.New("output-format choice invalid, use [json|yaml]")
		}

		return nil
	},
}

// selectModelPrompt lets user to select a relation if relations are more than one
func selectRelationshipPrompt(relationship []relationship.RelationshipDefinition) *relationship.RelationshipDefinition {
	relationshipNames := []string{}

	for _, _rel := range relationship {
		// here display Kind and EvaluationQuery as relationship name
		relationshipName := fmt.Sprintf("kind: %s, EvaluationPolicy: %s, SubType: %s", _rel.Kind, *_rel.EvaluationQuery, _rel.SubType)
		relationshipNames = append(relationshipNames, relationshipName)
	}

	prompt := promptui.Select{
		Label: "Select a relationship:",
		Items: relationshipNames,
	}

	for {
		i, _, err := prompt.Run()
		if err != nil {
			continue
		}

		return &relationship[i]
	}
}

func init() {
	viewCmd.Flags().StringP("output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
