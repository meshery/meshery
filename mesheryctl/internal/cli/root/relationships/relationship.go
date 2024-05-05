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
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/manifoldco/promptui"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	outFormatFlag string

	availableSubcommands = []*cobra.Command{ViewRelationshipsCmd, GenerateRelationshipDocsCmd}
)

var RelationshipCmd = &cobra.Command{
	Use:   "relationship",
	Short: "View list of relationships and details of relationship",
	Long:  "Meshery uses relationships to define how interconnected components interact. View list of relationships and detailed information of a specific relationship",
	Example: `
// To view list of relationships
mesheryctl exp relationships list

// To view a specific relationship
mesheryctl exp relationships view [model-name]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return system.ErrGetCurrentContext(err)
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
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
func selectRelationshipPrompt(relationship []entity.Entity) *v1alpha2.RelationshipDefinition {
	relationshipArray := []v1alpha2.RelationshipDefinition{}
	relationshipNames := []string{}

	for _, rel := range relationship {
		_rel, err := mutils.Cast[*v1alpha2.RelationshipDefinition](rel)
		if err != nil {
			continue
		}

		// here display Kind and EvaluationQuery as relationship name
		relationshipName := fmt.Sprintf("kind: %s, EvaluationPolicy: %s, SubType: %s", _rel.Kind, _rel.EvaluationQuery, _rel.SubType)
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

		return &relationshipArray[i]
	}
}
