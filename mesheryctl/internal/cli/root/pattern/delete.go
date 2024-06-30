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
	"os"
	"strings"

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
	Long:  `Delete pattern file will trigger deletion of the pattern file`,
	Example: `
// delete a pattern file
mesheryctl pattern delete -f [file | URL]

// delete a pattern file by name
mesheryctl pattern delete [pattern name]
`,
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Provide a pattern name or use the --file flag\nRun 'mesheryctl pattern delete --help' to see detailed help message"

		// Check if the file flag is set
		fileFlagSet := cmd.Flags().Changed("file")

		// Ensure only one of the pattern name or file flag is set
		if len(args) == 1 && !fileFlagSet {
			return nil
		}
		if len(args) == 0 && fileFlagSet {
			return nil
		}

		return utils.ErrInvalidArgument(errors.New(errMsg))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		contextID, err := getContextID(mctlCfg)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/deploy"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		// If a file path or URL is provided, delete the pattern using the file
		if file != "" {
			return deleteByFile(file, patternURL, deployURL, mctlCfg, contextID)
		}

		// If the --file flag is not set, delete the pattern using the pattern name
		return deleteByName(args, patternURL, deployURL, mctlCfg, contextID)
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&file, "file", "f", "", "Path to pattern file")
}

func deleteByFile(file string, patternURL string, deployURL string, mctlCfg *config.MesheryCtlConfig, contextID string) error {
	var patternFile string

	// If file path not a valid URL, treat it like a local file path
	if !govalidator.IsURL(file) {
		content, err := os.ReadFile(file)
		if err != nil {
			utils.Log.Error(utils.ErrFileRead(errors.New(utils.PatternError(fmt.Sprintf("failed to read file %s. Ensure the filename or URL is valid", file)))))
			return nil
		}

		patternFile = string(content)
	} else {
		// Else treat it like a URL
		url, path, err := utils.ParseURLGithub(file)
		if err != nil {
			utils.Log.Error(utils.ErrParseGithubFile(err, file))
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

		req, err := utils.NewRequest("POST", patternURL, bytes.NewBuffer(jsonValues))
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
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return nil
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			utils.Log.Error(utils.ErrUnmarshal(err))
			return nil
		}

		patternFile = response[0].PatternFile
	}

	err := utils.DeleteConfiguration(mctlCfg.GetBaseMesheryURL(), patternFile, "pattern")
	if err != nil {
		utils.Log.Error(err)
		return errors.Wrap(err, utils.PatternError(fmt.Sprintf("failed to delete pattern %s", patternFile)))
	}

	req, err := utils.NewRequest("DELETE", deployURL, bytes.NewBuffer([]byte(patternFile)))
	if err != nil {
		utils.Log.Error(utils.ErrCreatingRequest(err))
		return nil
	}

	res, err := utils.MakeRequest(req)
	if err != nil {
		utils.Log.Error(utils.ErrRequestResponse(err))
		return nil
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		utils.Log.Error(utils.ErrReadResponseBody(err))
		return nil
	}

	utils.Log.Info(string(body))

	return nil
}

func deleteByName(args []string, patternURL, deployURL string, mctlCfg *config.MesheryCtlConfig, contextID string) error {
	patternName := ""
	if len(args) > 0 {
		patternName = strings.Join(args, "%20")
	}

	// Fetch the pattern using the pattern name
	utils.Log.Debug("Fetching patterns")
	req, err := utils.NewRequest("GET", patternURL+"?search="+patternName, nil)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		utils.Log.Error(utils.ErrFailRequest(err))
		return nil
	}

	var response *models.PatternsAPIResponse
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, utils.PatternError("failed to read response body"))
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	// Delete the pattern using the pattern ID
	patternID := response.Patterns[0].ID
	patternFile := response.Patterns[0].PatternFile
	err = utils.DeleteConfiguration(mctlCfg.GetBaseMesheryURL(), patternID.String(), "pattern")
	if err != nil {
		utils.Log.Error(err)
		return errors.Wrap(err, utils.PatternError(fmt.Sprintf("failed to delete pattern %s", patternName)))
	}
	utils.Log.Info("Pattern ", patternName, " deleted successfully")

	payload := models.MesheryPatternFileDeployPayload{
		PatternFile: patternFile,
		PatternID:   patternID.String(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	req, err = utils.NewRequest("DELETE", deployURL+"?context_id="+contextID, bytes.NewBuffer(payloadBytes))
	if err != nil {
		utils.Log.Error(err)
		return utils.ErrCreatingRequest(err)
	}

	res, err := utils.MakeRequest(req)
	if err != nil {
		utils.Log.Error(utils.ErrFailRequest(err))
		return nil
	}

	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		utils.Log.Error(utils.ErrReadResponseBody(err))
		return nil
	}

	utils.Log.Info(string(body))

	return nil
}
