// Copyright 2023 Layer5, Inc.
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
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete pattern file",
	Long:  `delete pattern file will trigger deletion of the pattern file`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// delete a pattern file
mesheryctl pattern delete [file | URL]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		client := &http.Client{}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		pattern := ""
		isID := false
		if len(args) > 0 {
			pattern, isID, err = utils.Valid(args[0], "pattern")
			if err != nil {
				return err
			}
		}

		// Delete the pattern using the id
		if isID {
			err := utils.DeleteConfiguration(pattern, "pattern")
			if err != nil {
				return errors.Wrap(err, utils.PatternError(fmt.Sprintf("failed to delete pattern %s", args[0])))
			}
			utils.Log.Info("Pattern ", args[0], " deleted successfully")
			return nil
		}
		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/deploy"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		// If file path not a valid URL, treat it like a local file path
		if !govalidator.IsURL(file) {
			content, err := os.ReadFile(file)
			if err != nil {
				return errors.New(utils.PatternError(fmt.Sprintf("failed to read file %s. Ensure the filename or URL is valid", file)))
			}

			patternFile = string(content)
		} else {
			// Else treat it like a URL
			url, path, err := utils.ParseURLGithub(file)
			if err != nil {
				return err
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
				return err
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			utils.Log.Debug("remote hosted pattern request success")
			var response []*models.MesheryPattern
			// If API returns a non 200 status, return error
			if resp.StatusCode != 200 {
				return errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrap(err, utils.PerfError("failed to read response body"))
			}

			err = json.Unmarshal(body, &response)
			if err != nil {
				return errors.Wrap(err, "couldn't process JSON response from Meshery Server")
			}

			patternFile = response[0].PatternFile
		}

		req, err = utils.NewRequest("DELETE", deployURL, bytes.NewBuffer([]byte(patternFile)))
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		utils.Log.Info(string(body))

		return nil
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&file, "file", "f", "", "Path to pattern file")
}
