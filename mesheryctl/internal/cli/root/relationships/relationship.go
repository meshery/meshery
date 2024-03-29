// Copyright 2024 Layer5, Inc.
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
	"encoding/json"
	"fmt"
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/manifoldco/promptui"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	outFormatFlag string

	availableSubcommands = []*cobra.Command{ViewRelationshipsCmd}
)

var RelationshipCmd = &cobra.Command{
	Use:   "relationships",
	Short: "View list of relationships and detailed information of a specific relationship",
	Long:  "View list of relationships and detailed information of a specific relationship",
	Example: `
// To view list of components
mesheryctl exp relationships list

// To view a specific model
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
			utils.Log.Error(err)
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(system.ErrGetCurrentContext(err))
			return nil
		}
		err = ctx.ValidateVersion()
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl exp relationships --help' to display usage guide.\n", args[0]), "model"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		err = cmd.Usage()
		if err != nil {
			utils.Log.Error(err)
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
func selectRelationshipPrompt(relationship []v1alpha1.RelationshipDefinition) v1alpha1.RelationshipDefinition {
	relationshipArray := []v1alpha1.RelationshipDefinition{}
	relationshipNames := []string{}

	relationshipArray = append(relationshipArray, relationship...)

	for _, relationship := range relationshipArray {
		// here display Kind and EvaluationQuery as relationship name
		relationshipName := fmt.Sprintf("kind: %s, EvaluationPolicy: %s, SubType: %s", relationship.Kind, relationship.EvaluationQuery, relationship.SubType)
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

		return relationshipArray[i]
	}
}

func outputRelationshipJson(model v1alpha1.RelationshipDefinition) error {
	if err := prettifyRelationshipJson(model); err != nil {
		// if prettifyJson return error, marshal output in conventional way using json.MarshalIndent
		// but it doesn't convert unicode to its corresponding HTML string (it is default behavior)
		// e.g unicode representation of '&' will be printed as '\u0026'
		if output, err := json.MarshalIndent(model, "", "  "); err != nil {
			return errors.Wrap(err, "failed to format output in JSON")
		} else {
			fmt.Print(string(output))
		}
	}
	return nil
}

// prettifyJson takes a v1alpha1.RelationshipDefinition struct as input, marshals it into a nicely formatted JSON representation,
// and prints it to standard output with proper indentation and without escaping HTML entities.
func prettifyRelationshipJson(relationship v1alpha1.RelationshipDefinition) error {
	// Create a new JSON encoder that writes to the standard output (os.Stdout).
	enc := json.NewEncoder(os.Stdout)
	// Configure the JSON encoder settings.
	// SetEscapeHTML(false) prevents special characters like '<', '>', and '&' from being escaped to their HTML entities.
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	// Any errors during the encoding process will be returned as an error.
	return enc.Encode(relationship)
}
