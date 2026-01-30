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
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type relationshipViewFlags struct {
	outputFormat string
	save         bool
}

var (
	validOutputFormat             = []string{"json", "yaml"}
	relationshipViewFlagsProvided relationshipViewFlags
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view relationships of a model by its name",
	Long:  "view a relationship queried by the model name",
	Example: `
// View relationships of a model in default format yaml
mesheryctl exp relationship view [model-name]

// View relationships of a model in JSON format
mesheryctl exp relationship view [model-name] --output-format json

// View relationships of a model in json format and save it to a file
mesheryctl exp relationship view [model-name] --output-format json --save
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "\n\nUsage: mesheryctl exp relationship view [model-name]\nRun 'mesheryctl exp relationship view --help' to see detailed help message"
		if len(args) == 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("[model-name] isn't specified%s", errMsg))
		} else if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("Too many arguments only [model-name] is expected%s", errMsg))
		}

		// Validate output-format
		if !slices.Contains(validOutputFormat, strings.ToLower(relationshipViewFlagsProvided.outputFormat)) {
			return utils.ErrInvalidArgument(errors.New("output-format choice is invalid, use [json|yaml]"))
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

		switch relationshipsResponse.Count {
		case 0:
			utils.Log.Info("No relationship(s) found for the given name ", model)
			return nil
		case 1:
			selectedModel = &relationshipsResponse.Relationships[0]
		default:
			selectedModel = selectRelationshipPrompt(relationshipsResponse.Relationships)
		}

		outputFormatterFactory := display.OutputFormatterFactory[relationship.RelationshipDefinition]{}
		outputFormatter, err := outputFormatterFactory.New(relationshipViewFlagsProvided.outputFormat, *selectedModel)
		if err != nil {
			return err
		}

		err = outputFormatter.Display()
		if err != nil {
			return err
		}

		if relationshipViewFlagsProvided.save {
			// Get the home directory of user
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return utils.ErrRetrieveHomeDir(errors.Wrap(err, "failed to determine user home directory"))
			}

			shortID := selectedModel.Id.String()[:8]
			fileName := fmt.Sprintf("relationship_%s_%s", selectedModel.Model.Name, shortID)
			file := filepath.Join(homeDir, ".meshery", fileName)

			outputFormatterSaverFactory := display.OutputFormatterSaverFactory[relationship.RelationshipDefinition]{}
			outputFormatterSaver, err := outputFormatterSaverFactory.New(relationshipViewFlagsProvided.outputFormat, outputFormatter)
			if err != nil {
				return err
			}

			outputFormatterSaver = outputFormatterSaver.WithFilePath(file)
			err = outputFormatterSaver.Save()
			if err != nil {
				return err
			}

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
	viewCmd.Flags().StringVarP(&relationshipViewFlagsProvided.outputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewCmd.Flags().BoolVarP(&relationshipViewFlagsProvided.save, "save", "s", false, "(optional) save output as a JSON/YAML file")
}
