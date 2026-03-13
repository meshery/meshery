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
	"sort"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands = []*cobra.Command{viewCmd, generateCmd, listCmd, searchCmd, validateCmd}
	relationshipApiPath  = "api/meshmodels/relationships"
)

type MeshmodelRelationshipsAPIResponse struct {
	Page          int                                   `json:"page"`
	PageSize      int                                   `json:"page_size"`
	Count         int64                                 `json:"total_count"`
	Relationships []relationship.RelationshipDefinition `json:"relationships"`
}

var RelationshipCmd = &cobra.Command{
	Use:   "relationship",
	Short: "Manage relationships",
	Long: `Generate, list, search and view relationship(s) and detailed information
Meshery uses relationships to define how interconnected components interact.
`,
	Example: `
// Display number of available relationships in Meshery
mesheryctl relationship --count

// Generate a relationship documentation 
mesheryctl relationship generate [flags]

// List available relationship(s)
mesheryctl relationship list [flags]

// Search for a specific relationship
mesheryctl relationship search [--kind <kind>] [--type <type>] [--subtype <subtype>] [--model <model>]

// View a specific relationship
mesheryctl relationship view [model-name]

// Validate a relationship definition file
mesheryctl relationship validate --file ./relationship.yaml
`,
	Args: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetBool("count")
		if len(args) == 0 && !count {
			errMsg := "Usage: mesheryctl relationship [subcommand]\nRun 'mesheryctl relationship --help' to see detailed help message"
			return utils.ErrInvalidArgument(fmt.Errorf("no command specified. %s", errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		countFlag, _ := cmd.Flags().GetBool("count")
		if countFlag {
			models, err := api.Fetch[MeshmodelRelationshipsAPIResponse](relationshipApiPath)

			if err != nil {
				return err
			}

			utils.DisplayCount("relationships", models.Count)

			return nil
		}

		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return utils.ErrInvalidArgument(fmt.Errorf(
				"'%s' is an invalid subcommand. Please use one of [%s]. Use 'mesheryctl relationship --help' to display usage guide",
				args[0],
				strings.Join(relationshipSubcommandNames(), ", "),
			))
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
	RelationshipCmd.AddCommand(availableSubcommands...)
	RelationshipCmd.Flags().BoolP("count", "c", false, "(optional) Get the number of relationship(s) in total")
}

func relationshipSubcommandNames() []string {
	names := make([]string, 0, len(availableSubcommands))
	for _, subcommand := range availableSubcommands {
		names = append(names, subcommand.Name())
	}

	sort.Strings(names)

	return names
}

func generateRelationshipDataToDisplay(relationshipResponse *MeshmodelRelationshipsAPIResponse) ([][]string, int64) {
	defaultIfEmpty := func(value string) string {
		if value == "" {
			return "N/A"
		}
		return value
	}

	rows := make([][]string, 0, len(relationshipResponse.Relationships))
	for _, rel := range relationshipResponse.Relationships {
		rows = append(rows, []string{
			defaultIfEmpty(rel.Id.String()),
			string(rel.Kind),
			rel.Version,
			defaultIfEmpty(rel.Model.Name),
			rel.SubType,
			rel.RelationshipType,
		})
	}
	return rows, relationshipResponse.Count
}
