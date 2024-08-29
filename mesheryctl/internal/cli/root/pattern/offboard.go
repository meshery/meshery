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

package pattern

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var offboardCmd = &cobra.Command{
	Use:   "offboard",
	Short: "Offboard pattern",
	Long:  `Offboard pattern will trigger undeploy of pattern`,
	Example: `
// Offboard pattern by providing file path
mesheryctl pattern offboard -f [filepath]
	`,

	Args: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("file") && file == "" {
			errMsg := `Usage: mesheryctl pattern offboard -f [filepath]`
			return ErrOffboardPattern(errors.New(errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		var err error
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		pattern := ""
		isID := false
		if len(args) > 0 {
			pattern, isID, err = utils.ValidId(mctlCfg.GetBaseMesheryURL(), args[0], "pattern")
			if err != nil {
				utils.Log.Error(err)
				return nil
			}
		}

		// Delete the pattern using the id
		if isID {
			err := utils.DeleteConfiguration(mctlCfg.GetBaseMesheryURL(), pattern, "pattern")
			if err != nil {
				utils.Log.Error(err)
				return errors.Wrap(err, utils.PatternError(fmt.Sprintf("failed to delete pattern %s", args[0])))
			}
			utils.Log.Info("pattern ", args[0], " deleted successfully")
			return nil
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/deploy"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		// Read file
		if !govalidator.IsURL(file) {
			content, err := os.ReadFile(file)
			if err != nil {
				utils.Log.Error(utils.ErrFileRead(err))
				return utils.ErrFileRead(err)
			}

			patternFile = string(content)
		} else {
			utils.Log.Info("URLs are not currently supported")
		}

		// Convert pattern File into Pattern File
		jsonValues, _ := json.Marshal(map[string]interface{}{
			"K8sManifest": patternFile,
		})

		req, err = utils.NewRequest("POST", patternURL, bytes.NewBuffer(jsonValues))
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		defer resp.Body.Close()

		var response []*models.MesheryPattern

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return nil
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			utils.Log.Error(utils.ErrUnmarshal(err))
			return nil
		}

		if len(response) == 0 {
			return ErrPatternNotFound()
		}

		utils.Log.Debug("pattern file converted to pattern file")

		patternFile := response[0].PatternFile
		patternFileByt, err := yaml.Marshal(patternFile)

		if err != nil {
			return models.ErrMarshallingDesignIntoYAML(err)
		}

		req, err = utils.NewRequest("DELETE", deployURL, bytes.NewBuffer(patternFileByt))
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		res, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return nil
		}

		if res.StatusCode == 200 {
			utils.Log.Info("pattern successfully offboarded")
		}
		utils.Log.Info(string(body))

		return nil
	},
}

func init() {
	offboardCmd.Flags().StringVarP(&file, "file", "f", "", "Path to pattern file")
}
