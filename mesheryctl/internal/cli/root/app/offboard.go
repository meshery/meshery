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

package app

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

var offboardCmd = &cobra.Command{
	Use:   "offboard",
	Short: "Offboard application",
	Long:  `Offboard application will trigger undeploy of application`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// Offboard application by providing file path
mesheryctl app offboard -f [filepath]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("file") && file == "" {
			const errMsg = `Usage: mesheryctl app offboard -f [filepath]`
			return fmt.Errorf("no file path provided \n\n%v", errMsg)
		}
		var req *http.Request
		var err error
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(utils.ErrProcessingConfig(err))
			return utils.ErrProcessingConfig(err)
		}

		app := ""
		isID := false
		if len(args) > 0 {
			app, isID, err = utils.ValidId(args[0], "application")
			if err != nil {
				utils.Log.Error(utils.ErrInvalidNameOrID(err))
				return utils.ErrInvalidNameOrID(err)
			}
		}

		// Delete the app using the id
		if isID {
			err := utils.DeleteConfiguration(mctlCfg.GetBaseMesheryURL(), app, "application")
			if err != nil {
				utils.Log.Error(err)
				return errors.Wrap(err, utils.AppError(fmt.Sprintf("failed to delete application %s", args[0])))
			}
			utils.Log.Info("Application ", args[0], " deleted successfully")
			return nil
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/application/deploy"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		// Read file
		if !govalidator.IsURL(file) {
			content, err := os.ReadFile(file)
			if err != nil {
				utils.Log.Error(utils.ErrFileRead(err))
				return utils.ErrFileRead(err)
			}

			appFile = string(content)
		} else {
			utils.Log.Info("URLs are not currently supported")
		}

		// Convert App File into Pattern File
		jsonValues, _ := json.Marshal(map[string]interface{}{
			"K8sManifest": appFile,
		})

		req, err = utils.NewRequest("POST", patternURL, bytes.NewBuffer(jsonValues))
		if err != nil {
			utils.Log.Error(utils.ErrCreatingRequest(err))
			return utils.ErrCreatingRequest(err)
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(utils.ErrCreatingRequest(err))
			return utils.ErrCreatingRequest(err)
		}
		defer resp.Body.Close()

		var response []*models.MesheryPattern

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return utils.ErrReadResponseBody(err)
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			utils.Log.Error(utils.ErrUnmarshal(err))
			return utils.ErrUnmarshal(err)
		}

		utils.Log.Debug("application file converted to pattern file")

		patternFile := response[0].PatternFile

		req, err = utils.NewRequest("DELETE", deployURL, bytes.NewBuffer([]byte(patternFile)))
		if err != nil {
			utils.Log.Error(utils.ErrCreatingRequest(err))
			return utils.ErrCreatingRequest(err)
		}

		res, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(utils.ErrCreatingRequest(err))
			return utils.ErrCreatingRequest(err)
		}

		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return utils.ErrReadResponseBody(err)
		}

		if res.StatusCode == 200 {
			utils.Log.Info("app successfully offboarded")
		}
		utils.Log.Info(string(body))

		return nil
	},
}

func init() {
	offboardCmd.Flags().StringVarP(&file, "file", "f", "", "Path to app file")
}
