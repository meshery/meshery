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
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/components"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ViewRelationshipsCmd = &cobra.Command{
	Use:   "view",
	Short: "view relationships of a model by its name",
	Long:  "view a relationship queried by the model name",
	Example: `
// View relationships of a model
mesheryctl exp relationship view [model-name]
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp relationships view [model-name]\nRun 'mesheryctl exp relationships view --help' to see detailed help message"
		if len(args) == 0 {
			return utils.ErrInvalidArgument(errors.New("missing required argument: [model-name]. " + errMsg))
		} else if len(args) > 1 {
			return errors.New(utils.RelationshipsError(fmt.Sprintf("'%s' is an invalid subcommand. %s\n", args[0], errMsg), "view"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		model := args[0]

		url := fmt.Sprintf("%s/api/meshmodels/models/%s/relationships?pagesize=all", baseUrl, model)
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		relationshipsResponse := &MeshmodelRelationshipsAPIResponse{}

		err = json.Unmarshal(data, relationshipsResponse)
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
		outFormatFlag = strings.ToLower(outFormatFlag)
		if outFormatFlag == "yaml" || outFormatFlag == "yml" {
			if output, err = yaml.Marshal(selectedModel); err != nil {
				return errors.Wrap(err, "failed to format output in YAML")
			}
			utils.Log.Info(string(output))
		} else if outFormatFlag == "json" {
			// return outputRelationshipJson(selectedModel)
			return components.OutputJson(selectedModel)
		} else {
			return errors.New("output-format choice invalid, use [json|yaml]")
		}

		return nil
	},
}
