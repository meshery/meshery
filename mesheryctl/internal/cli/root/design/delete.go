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
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete design file",
	Long:  `delete design file will trigger deletion of the design file`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// delete a design file
mesheryctl design delete [file | URL]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
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
				utils.Log.Error(ErrPatternInvalidNameOrID(err))
				return nil
			}
		}

		// Delete the pattern using the id
		if isID {
			err := utils.DeleteConfiguration(mctlCfg.GetBaseMesheryURL(), pattern, "pattern")
			if err != nil {
				return errors.Wrap(err, utils.DesignError(fmt.Sprintf("failed to delete design %s", args[0])))
			}
			utils.Log.Info("Design ", args[0], " deleted successfully")
			return nil
		}
		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/deploy"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		var patternFileByt []byte
		// If file path not a valid URL, treat it like a local file path
		if !govalidator.IsURL(file) {
			content, err := os.ReadFile(file)
			if err != nil {
				return utils.ErrFileRead(errors.New(utils.DesignError(fmt.Sprintf("failed to read file %s. Ensure the filename or URL is valid", file))))
			}

			patternFileByt = content
		} else {
			// Else treat it like a URL
			url, path, err := utils.ParseURLGithub(file)
			if err != nil {
				return utils.ErrParseGithubFile(err, file)
			}

			utils.Log.Debug(url)
			utils.Log.Debug(path)

			var jsonValues []byte

			// Send the URL and path to the server and let it fetch the patternfile and delete
			// the components
			if path != "" {
				jsonValues, _ = json.Marshal(map[string]interface{}{
					"url":  url,
					"path": path,
					"save": false,
				})
			} else {
				jsonValues, _ = json.Marshal(map[string]interface{}{
					"url":  url,
					"save": false,
				})
			}

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
			utils.Log.Debug("remote hosted pattern request success")
			var response []*models.MesheryPattern
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				utils.Log.Error(utils.ErrReadResponseBody(errors.Wrap(err, "failed to read response body")))
				return nil
			}

			err = json.Unmarshal(body, &response)
			if err != nil {
				utils.Log.Error(utils.ErrUnmarshal(err))
				return nil
			}

			if len(response) == 0 {
				return ErrDesignNotFound()
			}

			patternFileByt, err = yaml.Marshal(response[0].PatternFile)
			if err != nil {
				return models.ErrMarshallingDesignIntoYAML(err)
			}
		}

		req, err = utils.NewRequest("DELETE", deployURL, bytes.NewBuffer(patternFileByt))
		if err != nil {
			return utils.ErrCreatingRequest(err)
		}

		res, err := utils.MakeRequest(req)
		if err != nil {
			return utils.ErrRequestResponse(err)
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return nil
		}

		utils.Log.Info(string(body))

		return nil
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&file, "file", "f", "", "Path to design file")
}
