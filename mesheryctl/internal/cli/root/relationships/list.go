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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the mesheryctl exp relationships list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered relationships",
	Long:  "List all relationships registered in Meshery Server",
	Example: `
// List of relationships
mesheryctl exp relationship list

// List of relationships for a specified page
mesheryctl relationship list --page [page-number]

// Display number of available relationships in Meshery
mesheryctl relationship list --count
`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp relationship list \nRun 'mesheryctl exp relationship list --help' to see detailed help message"
		if len(args) != 0 {
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.RelationshipsError("this command takes no arguments\n", "list"))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		page, _ := cmd.Flags().GetInt("page")
		url := fmt.Sprintf("%s/api/meshmodels/relationships?%s", baseUrl, utils.GetPageQueryParameter(cmd, page))

		relationships, err := api.Fetch[MeshmodelRelationshipsAPIResponse](url)

		if err != nil {
			return err
		}

		count, _ := cmd.Flags().GetBool("count")
		var rows [][]string

		for _, rel := range relationships.Relationships {
			evaluationQuery := ""
			if rel.EvaluationQuery != nil {
				evaluationQuery = *rel.EvaluationQuery
			}
			if len(rel.GetEntityDetail()) > 0 {
				rows = append(rows, []string{string(rel.Kind), rel.Version, rel.Model.Name, rel.SubType, evaluationQuery})
			}
		}

		dataToDisplay := display.DisplayedData{
			DataType:         "relationship",
			Header:           []string{"kind", "API Version", "Model name", "Sub Type", "Evaluation Policy"},
			Rows:             rows,
			Count:            relationships.Count,
			DisplayCountOnly: count,
			IsPage:           cmd.Flags().Changed("page"),
		}

		return display.List(dataToDisplay)
	},
}

func init() {
	// Add the new exp relationship commands to the listRelationshipsCmd
	listCmd.Flags().IntP("page", "p", 1, "(optional) List next set of relationships with --page (default = 1)")
	listCmd.Flags().BoolP("count", "c", false, "(optional) Get the number of relationship(s) in total")
}
