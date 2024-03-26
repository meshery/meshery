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
	sourceType string // pattern file type (manifest / compose)
)

var linkDocpatternOnboard = map[string]string{
	"link":    "![pattern-onboard-usage](/assets/img/mesheryctl/pattern-onboard.png)",
	"caption": "Usage of mesheryctl pattern onboard",
}

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Onboard pattern",
	Long:  `Command will trigger deploy of pattern`,
	Example: `
// Onboard pattern by providing file path
mesheryctl pattern onboard -f [filepath] -s [source type]
mesheryctl pattern onboard -f ./pattern.yml -s "Kubernetes Manifest"
	`,
	Annotations: linkDocpatternOnboard,
	Args: func(_ *cobra.Command, args []string) error {

		if file == "" && len(args) == 0 {
			return ErrOnboardPattern()
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
			utils.Log.Debug("Fetching patterns")

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
				utils.Log.Error(utils.ErrReadResponseBody(err))
				return nil
			}
			err = json.Unmarshal(body, &response)
			if err != nil {
				utils.Log.Error(utils.ErrUnmarshal(err))
				return nil
			}

			index := 0
			if len(response.Patterns) == 0 {
				utils.Log.Error(utils.ErrNotFound(errors.New("no pattern found with the given name")))
				return nil
			} else if len(response.Patterns) == 1 {
				patternFile = response.Patterns[0].PatternFile
			} else {
				// Multiple patterns with same name
				index = multiplepatternsConfirmation(response.Patterns)
				patternFile = response.Patterns[index].PatternFile
			}
		} else {
			// Check if a valid source type is set
			if sourceType, err = getFullSourceType(sourceType); err != nil {
				return ErrInValidSource(sourceType, validSourceTypes)
			}
			pattern, err := importPattern(sourceType, file, patternURL, !skipSave)
			if err != nil {
				utils.Log.Error(err)
				return nil
			}

			patternFile = pattern.PatternFile
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

		res, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		if res.StatusCode == 200 {
			utils.Log.Info("pattern successfully onboarded")
		}
		utils.Log.Info(string(body))
		return nil
	},
}

func multiplepatternsConfirmation(profiles []models.MesheryPattern) int {
	reader := bufio.NewReader(os.Stdin)

	for index, a := range profiles {
		fmt.Printf("Index: %v\n", index)
		fmt.Printf("Name: %v\n", a.Name)
		fmt.Printf("ID: %s\n", a.ID.String())
		fmt.Printf("patternFile:\n")
		fmt.Printf(a.PatternFile)
		fmt.Println("---------------------")
	}

	for {
		fmt.Printf("Enter the index of pattern: ")
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

func getSourceTypes() error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		utils.Log.Error(err)
		return nil
	}
	validTypesURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/types"
	req, err := utils.NewRequest("GET", validTypesURL, nil)
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

	var response []*models.PatternSourceTypesAPIResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Log.Error(utils.ErrReadResponseBody(errors.Wrap(err, "couldn't read response from server. Please try again after some time")))
		return nil
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.Log.Error(utils.ErrUnmarshal(errors.Wrap(err, "couldn't process response received from server")))
		return nil
	}

	for _, apiResponse := range response {
		validSourceTypes = append(validSourceTypes, apiResponse.DesignType)
	}

	return nil
}

// returns full source name e.g. helm -> `Helm Chart`
// user passes only helm, compose or manifest but server accepts full source type
// e.g `Heml Chart`, `Docker Compose`, `Kubernetes Manifest`
func getFullSourceType(sType string) (string, error) {
	for _, validType := range validSourceTypes {
		lowerType := strings.ToLower(validType)
		// user may pass Pascal Case source e.g Helm
		sType = strings.ToLower(sType)
		if strings.Contains(lowerType, sType) {
			return validType, nil
		}
	}

	return sType, fmt.Errorf("no matching source type found")
}

func init() {
	onboardCmd.Flags().StringVarP(&file, "file", "f", "", "Path to pattern file")
	onboardCmd.Flags().BoolVarP(&skipSave, "skip-save", "", false, "Skip saving a pattern")
	onboardCmd.Flags().StringVarP(&sourceType, "source-type", "s", "", "Type of source file (ex. manifest / compose / helm)")
}
