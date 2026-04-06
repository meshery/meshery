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

package design

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var designUndeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploy design",
	Long:  `Undeploy design will trigger undeploy of design`,
	Example: `
// Undeploy design by providing file path
mesheryctl design undeploy -f [filepath]
	`,

	Args: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("file") && file == "" {
			errMsg := `Usage: mesheryctl design undeploy -f [filepath]`
			return ErrUndeployDesign(fmt.Errorf("%s", errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		var err error
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		design := ""
		isID := false
		if len(args) > 0 {
			design, isID, err = utils.ValidId(mctlCfg.GetBaseMesheryURL(), args[0], "pattern")
			if err != nil {
				return err
			}
		}

		// Delete the design using the id
		if isID {
			err := utils.DeleteConfiguration(mctlCfg.GetBaseMesheryURL(), design, "pattern")
			if err != nil {
				return ErrDeleteDesign(err, args[0])
			}
			utils.Log.Infof("design %s deleted", args[0])
			return nil
		}

		designURLPath := "api/pattern"

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/deploy"

		if !govalidator.IsURL(file) {
			content, err := os.ReadFile(file)
			if err != nil {
				return utils.ErrFileRead(err)
			}
			designFile = string(content)
		} else {
			utils.Log.Info("URLs are not currently supported")
		}

		jsonValues, _ := json.Marshal(map[string]interface{}{
			"K8sManifest": designFile,
		})

		resp, err := api.Add(designURLPath, bytes.NewBuffer(jsonValues), nil)
		if err != nil {
			return err
		}
		defer func() { _ = resp.Body.Close() }()

		var response []*models.MesheryPattern

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return utils.ErrReadFromBody(err)
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return utils.ErrUnmarshal(err)
		}

		if len(response) == 0 {
			return ErrDesignNotFound(file)
		}

		utils.Log.Debug("design file converted to design")

		patternFile := response[0].PatternFile
		patternFileByt, err := yaml.Marshal(patternFile)
		if err != nil {
			return models.ErrMarshallingDesignIntoYAML(err)
		}

		req, err = utils.NewRequest("DELETE", deployURL, bytes.NewBuffer(patternFileByt))
		if err != nil {
			return err
		}

		res, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}

		defer func() { _ = res.Body.Close() }()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return utils.ErrReadFromBody(err)
		}

		if res.StatusCode == 200 {
			utils.Log.Info("design undeployed")
		}
		utils.Log.Info(string(body))

		return nil
	},
}

func init() {
	designUndeployCmd.Flags().StringVarP(&file, "file", "f", "", "Path to design file")
}
