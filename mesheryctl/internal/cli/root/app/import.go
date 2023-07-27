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
	"path"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import app manifests",
	Long:  `Import the app manifest into Meshery`,
	Example: `
// Import app manifest
mesheryctl app import -f [file/URL] -s [source-type]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = `Usage: mesheryctl app import -f [file/url] -s [source-type]
Example: mesheryctl app import -f ./application.yml -s "Kubernetes Manifest"`

		if file == "" {
			utils.Log.Debug("manifest path not provided")
			return fmt.Errorf("manifest path not provided. Provide the path to the app manifest. \n\n%v", errMsg)
		}

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return getSourceTypes()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log, err := logger.New("app", logger.Options{
			Format:     logger.SyslogLogFormat,
			DebugLevel: true,
		})
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Error(err)
			return nil
		}

		appURL := mctlCfg.GetBaseMesheryURL() + "/api/application"

		// If app file is passed via flags
		if !isValidSource(sourceType) {
			return errors.Errorf("Application source type (-s) invalid or not passed.\nAllowed source types: %s", strings.Join(validSourceTypes, ", "))
		}

		app, err := importApp(sourceType, file, appURL, true)

		if err != nil {
			log.Error(err)
			return nil
		}

		fmt.Printf("App file imported successfully. \nID of the app: %s \n", utils.TruncateID(app.ID.String()))

		return nil
	},
}

func importApp(sourceType string, file string, appURL string, save bool) (*models.MesheryApplication, error) {
	var req *http.Request
	var app *models.MesheryApplication

	// Check if the app manifest is file or URL
	if validURL := govalidator.IsURL(file); !validURL {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, utils.ErrFileRead(err)
		}
		text := string(content)

		jsonValues, err := json.Marshal(map[string]interface{}{
			"application_data": map[string]interface{}{
				"name":             path.Base(file),
				"application_file": text,
			},
			"save": save,
		})
		if err != nil {
			return nil, utils.ErrMarhalling(err)
		}
		req, err = utils.NewRequest("POST", appURL+"/"+sourceType, bytes.NewBuffer(jsonValues))
		if err != nil {
			return nil, err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return nil, err
		}
		utils.Log.Debug("App file saved")
		var response []*models.MesheryApplication
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Debug("failed to read response body")
			return nil, utils.ErrReadResponseBody(err)
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			utils.Log.Debug("failed to unmarshal JSON response")
			return nil, utils.ErrUnmarshal(err)
		}
		// set app
		app = response[0]
	} else {
		var jsonValues []byte
		url, path, err := utils.ParseURLGithub(file)
		if err != nil {
			return nil, err
		}

		utils.Log.Debug(url)
		utils.Log.Debug(path)

		// save the app with Github URL
		if path != "" {
			jsonValues, _ = json.Marshal(map[string]interface{}{
				"url":  url,
				"path": path,
				"save": save,
			})
		} else {
			jsonValues, _ = json.Marshal(map[string]interface{}{
				"url":  url,
				"save": save,
			})
		}

		req, err = utils.NewRequest("POST", appURL+"/"+sourceType, bytes.NewBuffer(jsonValues))
		if err != nil {
			return nil, err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return nil, err
		}
		utils.Log.Debug("remote hosted app request success")
		var response []*models.MesheryApplication
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Debug("failed to read response body")
			return nil, utils.ErrReadResponseBody(err)
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			utils.Log.Debug("failed to unmarshal JSON response")
			return nil, utils.ErrUnmarshal(err)
		}

		// set app
		app = response[0]
	}

	return app, nil
}

func init() {
	importCmd.Flags().StringVarP(&file, "file", "f", "", "Path/URL to app file")
	importCmd.Flags().StringVarP(&sourceType, "source-type", "s", "", "Type of source file (ex. manifest / compose / helm)")
}
