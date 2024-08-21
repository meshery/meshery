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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the mesheryctl exp relationships list command
var listRelationshipsCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered relationships",
	Long:  "List all relationships registered in Meshery Server",
	Example: `
	View list of relationship
    mesheryctl exp relationship list
    View list of relationship with specified page number (25 relationships per page)
    mesheryctl exp relationship list --page 2
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
		var url string
		if cmd.Flags().Changed("page") {
			url = fmt.Sprintf("%s/api/meshmodels/relationships?page=%d", baseUrl, pageNumberFlag)
		} else {
			url = fmt.Sprintf("%s/api/meshmodels/relationships?pagesize=all", baseUrl)
		}
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		relationshipsResponse := &MeshmodelRelationshipsAPIResponse{}
		err = json.Unmarshal(data, relationshipsResponse)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		header := []string{"kind", "API Version", "Model name", "Sub Type", "Evaluation Policy"}
		var rows [][]string

		for _, rel := range relationshipsResponse.Relationships {
			evaluationQuery := ""
			if rel.EvaluationQuery != nil {
				evaluationQuery = *rel.EvaluationQuery
			}
			if len(rel.GetEntityDetail()) > 0 {
				rows = append(rows, []string{string(rel.Kind), rel.Version, rel.Model.Name, rel.SubType, evaluationQuery})
			}
		}

		if len(rows) == 0 {
			// if no component is found
			fmt.Println("No relationship(s) found")
			return nil
		}

		if cmd.Flags().Changed("page") {
			utils.PrintToTable(header, rows)
		} else {
			maxRowsPerPage := 25
			err := utils.HandlePagination(maxRowsPerPage, "relationships", rows, header)
			if err != nil {
				utils.Log.Error(err)
				return err
			}
		}
		return nil
	},
}

func init() {
	// Add the new exp relationship commands to the listRelationshipsCmd
	listRelationshipsCmd.Flags().IntVarP(&pageNumberFlag, "page", "p", 1, "(optional) List next set of relationships with --page (default = 1)")
}
