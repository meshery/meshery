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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	skipSave   bool   // skip saving the app
	appFile    string // app file
	sourceType string // app file type (manifest / compose)
)

var linkDocAppOnboard = map[string]string{
	"link":    "![app-onboard-usage](/assets/img/mesheryctl/app-onboard.png)",
	"caption": "Usage of mesheryctl app onboard",
}

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Onboard application",
	Long:  `Command will trigger deploy of application`,
	Example: `
// Onboard application by providing file path
mesheryctl app onboard -f [filepath] -s [source type]
mesheryctl app onboard -f ./application.yml -s "Kubernetes Manifest"
	`,
	Annotations: linkDocAppOnboard,
	Args: func(_ *cobra.Command, args []string) error {

		if file == "" && len(args) == 0 {
			return ErrOnboardApp()
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return getSourceTypes()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		var err error

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/application/deploy"
		appURL := mctlCfg.GetBaseMesheryURL() + "/api/application"

		// app name has been passed
		if len(args) > 0 {
			// Merge args to get app-name
			appName := strings.Join(args, "%20")

			// search and fetch apps with app-name
			utils.Log.Debug("Fetching apps")

			req, err = utils.NewRequest("GET", appURL+"?search="+appName, nil)
			if err != nil {
				utils.Log.Error(utils.ErrCreatingRequest(err))
				return utils.ErrCreatingRequest(err)
			}

			resp, err := utils.MakeRequest(req)
			if err != nil {
				utils.Log.Error(utils.ErrRequestResponse(err))
				return utils.ErrRequestResponse(err)
			}

			var response *models.ApplicationsAPIResponse
			defer resp.Body.Close()
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

			index := 0
			if len(response.Applications) == 0 {
				utils.Log.Error(utils.ErrNotFound(errors.New("no app found with the given name")))
				return ErrAppFound()
			} else if len(response.Applications) == 1 {
				appFile = response.Applications[0].ApplicationFile
			} else {
				// Multiple apps with same name
				index = multipleApplicationsConfirmation(response.Applications)
				appFile = response.Applications[index].ApplicationFile
			}
		} else {
			// Check if a valid source type is set
			if !isValidSource(sourceType) {
				return ErrValidSource(validSourceTypes)
			}
			app, err := importApp(sourceType, file, appURL, !skipSave)
			if err != nil {
				utils.Log.Error(ErrImportApp(err))
				return ErrImportApp(err)
			}

			appFile = app.ApplicationFile
		}

		req, err = utils.NewRequest("POST", deployURL, bytes.NewBuffer([]byte(appFile)))
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
			return utils.ErrReadResponseBody(err)
		}

		if res.StatusCode == 200 {
			utils.Log.Info("app successfully onboarded")
		}
		utils.Log.Info(string(body))
		return nil
	},
}

func multipleApplicationsConfirmation(profiles []models.MesheryApplication) int {
	reader := bufio.NewReader(os.Stdin)

	for index, a := range profiles {
		fmt.Printf("Index: %v\n", index)
		fmt.Printf("Name: %v\n", a.Name)
		fmt.Printf("ID: %s\n", a.ID.String())
		fmt.Printf("ApplicationFile:\n")
		fmt.Printf(a.ApplicationFile)
		fmt.Println("---------------------")
	}

	for {
		fmt.Printf("Enter the index of app: ")
		response, err := reader.ReadString('\n')
		if err != nil {
			utils.Log.Info(err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		index, err := strconv.Atoi(response)
		if err != nil {
			utils.Log.Info(err)
		}
		if index < 0 || index >= len(profiles) {
			utils.Log.Info("Invalid index")
		} else {
			return index
		}
	}
}

func init() {
	onboardCmd.Flags().StringVarP(&file, "file", "f", "", "Path to app file")
	onboardCmd.Flags().BoolVarP(&skipSave, "skip-save", "", false, "Skip saving a app")
	onboardCmd.Flags().StringVarP(&sourceType, "source-type", "s", "", "Type of source file (ex. manifest / compose / helm)")
}
