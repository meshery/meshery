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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/meshery/schemas/models/v1alpha3/relationship"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	outFormatFlag        string
	pageNumberFlag       int
	availableSubcommands = []*cobra.Command{ViewRelationshipsCmd, GenerateRelationshipDocsCmd, listRelationshipsCmd, SearchComponentsCmd}
	maxRowsPerPage       = 25
)

type MeshmodelRelationshipsAPIResponse struct {
	Page          int                                   `json:"page"`
	PageSize      int                                   `json:"page_size"`
	Count         int64                                 `json:"total_count"`
	Relationships []relationship.RelationshipDefinition `json:"relationships"`
}

var RelationshipCmd = &cobra.Command{
	Use:   "relationship",
	Short: "View list of relationships and details of relationship",
	Long:  "Meshery uses relationships to define how interconnected components interact. View list of relationships and detailed information of a specific relationship",
	Example: `
// To view list of relationships
mesheryctl exp relationships list

// To view a specific relationship
mesheryctl exp relationships view [model-name]

//To search a specific relationship
mesheryctl exp relationships search --[flag] [query-text]

	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			errMsg := "Usage: mesheryctl exp relationships [subcommand]\nRun 'mesheryctl exp relationships --help' to see detailed help message"
			return utils.ErrInvalidArgument(errors.New("missing required argument: [model-name]. " + errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.RelationshipsError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl exp relationships --help' to display usage guide.\n", args[0]), "relationship"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		err = cmd.Usage()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	ViewRelationshipsCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json| yaml]")
	RelationshipCmd.AddCommand(availableSubcommands...)
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
