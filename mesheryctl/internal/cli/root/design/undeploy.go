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
	"net/url"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var designUndeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploy design",
	Long:  `Undeploy design will trigger undeploy of design`,
	Example: `
// Undeploy a deployed design by its name or ID
mesheryctl design undeploy [design name | ID]

// Undeploy a design by providing a design file path
mesheryctl design undeploy -f [filepath]
	`,

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && file == "" {
			return ErrUndeployDesign(
				fmt.Errorf("provide either a design name/ID or -f [filepath]"),
			)
		}

		if cmd.Flags().Changed("file") && file == "" {
			errMsg := `Usage: mesheryctl design undeploy -f [filepath]`
			return ErrUndeployDesign(fmt.Errorf("%s", errMsg))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/deploy"

		// Resolve the design file content to undeploy. A design can be referenced
		// either by its name/ID (a previously saved design) or by a design file
		// path. In every case the design is undeployed; the saved design record
		// is never deleted.
		var patternFileContent string
		if len(args) > 0 {
			patternFileContent, err = getUndeployPatternFileByNameOrID(mctlCfg.GetBaseMesheryURL(), args[0])
		} else {
			patternFileContent, err = getUndeployPatternFileFromFile(file)
		}
		if err != nil {
			return err
		}

		return undeployPatternFile(deployURL, patternFileContent)
	},
}

// getUndeployPatternFileByNameOrID resolves a saved design, referenced by name
// or ID, to its design file content. It only reads the design; it must never
// delete it.
func getUndeployPatternFileByNameOrID(baseMesheryURL, arg string) (string, error) {
	design, isID, err := utils.ValidId(baseMesheryURL, arg, "pattern")
	if err != nil {
		return "", err
	}

	if isID {
		resp, err := api.Fetch[models.MesheryPattern](fmt.Sprintf("api/pattern/%s", url.PathEscape(design)))
		if err != nil {
			return "", err
		}
		if resp.PatternFile == "" {
			return "", ErrDesignNotFound(arg)
		}
		return resp.PatternFile, nil
	}

	queryParams := url.Values{}
	queryParams.Set("populate", "pattern_file")
	queryParams.Set("search", design)
	urlPath := fmt.Sprintf("api/pattern?%s", queryParams.Encode())

	resp, err := api.Fetch[models.PatternsAPIResponse](urlPath)
	if err != nil {
		return "", err
	}
	if len(resp.Patterns) == 0 {
		return "", ErrDesignNotFound(design)
	}
	return resp.Patterns[0].PatternFile, nil
}

// getUndeployPatternFileFromFile imports the given design file into a saved
// design and returns the resulting design file content.
func getUndeployPatternFileFromFile(filePath string) (string, error) {
	if govalidator.IsURL(filePath) {
		return "", ErrUndeployDesign(fmt.Errorf("URLs are not currently supported"))
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", utils.ErrFileRead(err)
	}

	jsonValues, _ := json.Marshal(map[string]interface{}{
		"K8sManifest": string(content),
	})

	resp, err := api.Add("api/pattern", bytes.NewBuffer(jsonValues), nil)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", utils.ErrReadFromBody(err)
	}

	var response []*models.MesheryPattern
	if err := json.Unmarshal(body, &response); err != nil {
		return "", utils.ErrUnmarshal(err)
	}
	if len(response) == 0 {
		return "", ErrDesignNotFound(filePath)
	}

	utils.Log.Debug("design file converted to design")
	return response[0].PatternFile, nil
}

// undeployPatternFile sends the resolved design file to the server's undeploy
// endpoint using the same request contract as `design deploy`: a JSON-encoded
// models.MesheryPatternFileDeployPayload sent to DELETE /api/pattern/deploy.
func undeployPatternFile(deployURL, patternFileContent string) error {
	payload := models.MesheryPatternFileDeployPayload{
		PatternFile: patternFileContent,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return utils.ErrMarshal(err)
	}

	req, err := utils.NewRequest("DELETE", deployURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	res, err := utils.MakeRequest(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return utils.ErrReadFromBody(err)
	}

	if res.StatusCode == http.StatusOK {
		utils.Log.Info("design undeployed")
	}
	utils.Log.Info(string(body))

	return nil
}

func init() {
	designUndeployCmd.Flags().StringVarP(&file, "file", "f", "", "Path to design file")
}
