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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	skipSave    bool // skip saving a design
	patternFile string
)

var linkDocPatternApply = map[string]string{
	"link":    "![pattern-apply-usage](/assets/img/mesheryctl/patternApply.png)",
	"caption": "Usage of mesheryctl design apply",
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply design file",
	Long:  `Apply design will trigger deploy of the design file`,
	Example: `
// apply a design file
mesheryctl design apply -f [file | URL]

// deploy a saved design
mesheryctl design apply [design-name]
	`,
	Annotations: linkDocPatternApply,
	Args:        cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/deploy"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		// pattern name has been passed
		if len(args) > 0 {
			// Merge args to get pattern-name
			patternName := strings.Join(args, "%20")

			// search and fetch patterns with pattern-name
			utils.Log.Debug("Fetching designs")

			req, err = utils.NewRequest("GET", patternURL+"?search="+patternName, nil)
			if err != nil {
				utils.Log.Error(err)
				return nil
			}

			resp, err := utils.MakeRequest(req)
			if err != nil {
				utils.Log.Error(err)
				return nil
			}

			var response *models.PatternsAPIResponse
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrap(err, utils.DesignError("failed to read response body"))
			}
			err = json.Unmarshal(body, &response)
			if err != nil {
				utils.Log.Error(err)
				return nil
			}

			index := 0
			if len(response.Patterns) == 0 {
				utils.Log.Error(ErrDesignNotFound())
				return nil
			} else if len(response.Patterns) == 1 {
				patternFile = response.Patterns[0].PatternFile
			} else {
				// Multiple patterns with same name
				index = multiplePatternsConfirmation(response.Patterns)
				patternFile = response.Patterns[index].PatternFile
			}
		} else {
			// Method to check if the entered file is a URL or not
			validURL := strings.HasPrefix(file, "https://github.com") || strings.HasPrefix(file, "https://raw.githubusercontent.com")
			if !validURL {
				content, err := os.ReadFile(file)
				if err != nil {
					utils.Log.Error(utils.ErrFileRead(errors.Errorf("file path %s is invalid. Enter a valid path ", file)))
					return nil
				}

				// if --skip-save is not passed we save the pattern first
				if !skipSave {
					jsonValues, err := json.Marshal(map[string]interface{}{
						"pattern_data": map[string]interface{}{
							"name":         path.Base(file),
							"pattern_file": content,
						},
						"save": true,
					})
					if err != nil {
						return err
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

					var response []*models.MesheryPattern
					defer resp.Body.Close()

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return errors.Wrap(err, utils.DesignError("failed to read response body"))
					}
					err = json.Unmarshal(body, &response)
					if err != nil {
						utils.Log.Error(utils.ErrUnmarshal(err))
						return nil
					}
				}

				// setup pattern file
				patternFile = string(content)
			} else {
				var jsonValues []byte
				url, path, err := utils.ParseURLGithub(file)
				if err != nil {
					utils.Log.Error(utils.ErrParseGithubFile(err, file))
					return nil
				}

				utils.Log.Debug(url)
				utils.Log.Debug(path)

				// save the pattern with Github URL
				if !skipSave {
					if path != "" {
						jsonValues, _ = json.Marshal(map[string]interface{}{
							"url":  url,
							"path": path,
							"save": true,
						})
					} else {
						jsonValues, _ = json.Marshal(map[string]interface{}{
							"url":  url,
							"save": true,
						})
					}
				} else { // we don't save the pattern
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
					utils.Log.Error(utils.ErrUnmarshal(errors.Wrap(err, "failed to unmarshal response body")))
					return nil
				}

				// setup pattern file here
				patternFile = response[0].PatternFile
			}

		}

		payload := models.MesheryPatternFileDeployPayload{
			PatternFile: patternFile,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		req, err = utils.NewRequest("POST", deployURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		pf, err := core.NewPatternFile([]byte(patternFile))
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		s := utils.CreateDefaultSpinner("Applying design "+pf.Name, "")
		s.Start()
		res, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		s.Stop()
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return nil
		}

		if res.StatusCode == 200 {
			utils.Log.Info("design successfully applied")
		}
		utils.Log.Info(string(body))
		return nil
	},
}

func multiplePatternsConfirmation(profiles []models.MesheryPattern) int {
	reader := bufio.NewReader(os.Stdin)

	for index, a := range profiles {
		fmt.Printf("Index: %v\n", index)
		fmt.Printf("Name: %v\n", a.Name)
		fmt.Printf("ID: %s\n", a.ID.String())
		fmt.Printf("DesignFile:\n")
		fmt.Printf(a.PatternFile)
		fmt.Println("---------------------")
	}

	for {
		fmt.Printf("Enter the index of profile: ")
		response, err := reader.ReadString('\n')
		if err != nil {
			utils.Log.Warn(err)
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
	applyCmd.Flags().StringVarP(&file, "file", "f", "", "Path to design file")
	applyCmd.Flags().BoolVarP(&skipSave, "skip-save", "", false, "Skip saving a design")
}
