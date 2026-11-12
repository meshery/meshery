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
	"path/filepath"
	"strconv"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/pattern/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	skipSave   bool // skip saving a design
	designFile string
)

var linkDocPatternApply = map[string]string{
	"link":    "![pattern-apply-usage](../../../images/patternApply.png)",
	"caption": "Usage of mesheryctl design apply",
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply design file",
	Long: `Apply design will trigger deploy of the design file.
	Find more information at: https://docs.meshery.io/reference/mesheryctl/design/apply`,
	Example: `
// apply a design file
mesheryctl design apply -f [file | URL]

// deploy a saved design
mesheryctl design apply [design-name]
	`,
	Annotations: linkDocPatternApply,
	Args: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("file") && strings.TrimSpace(file) == "" {
			return utils.ErrInvalidArgument(fmt.Errorf("file path cannot be empty"))
		}
		hasDesignName := strings.TrimSpace(strings.Join(args, " ")) != ""
		hasFile := strings.TrimSpace(file) != ""

		if !hasDesignName && !hasFile {
			return utils.ErrInvalidArgument(fmt.Errorf("please provide a design name or specify a file using --file"))
		}
		if len(args) > 0 && hasFile {
			return utils.ErrInvalidArgument(fmt.Errorf("cannot specify both a design name and the --file flag"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/deploy"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		// pattern name has been passed
		if len(args) > 0 {
			// Merge args to get pattern-name
			patternName := strings.Join(args, "%20")

			// search and fetch patterns with pattern-name
			utils.Log.Debug("Fetching designs")

			req, err = utils.NewRequest("GET", patternURL+"?populate=pattern_file&search="+patternName, nil)
			if err != nil {
				return err
			}

			resp, err := utils.MakeRequest(req)
			if err != nil {
				return err
			}

			var response *models.PatternsAPIResponse
			defer func() { _ = resp.Body.Close() }()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return utils.ErrReadFromBody(err)
			}
			err = json.Unmarshal(body, &response)
			if err != nil {
				return utils.ErrInvalidAPIResponse(err)
			}

			index := 0
			if len(response.Patterns) == 0 {
				return ErrDesignNotFound(patternName)
			} else if len(response.Patterns) == 1 {
				designFile = response.Patterns[0].PatternFile
			} else {
				// Multiple patterns with same name
				index = multiplePatternsConfirmation(response.Patterns)
				designFile = response.Patterns[index].PatternFile
			}
		} else {
			var content []byte
			if utils.IsValidUrl(file) {
				content, err = fetchDesignFile(file)
				if err != nil {
					return err
				}
			} else {
				content, err = os.ReadFile(file)
				if err != nil {
					return utils.ErrFileRead(fmt.Errorf(errInvalidPathMsg, file))
				}
			}

			// if --skip-save is not passed we save the pattern first
			if !skipSave {
				fileName := designFileName(file)
				if _, err := importDesignFile(patternURL+"/import", fileName, fileName, content); err != nil {
					return err
				}
			}

			// setup pattern file
			designFile = string(content)
		}

		payload := models.MesheryPatternFileDeployPayload{
			PatternFile: designFile,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return utils.ErrMarshal(err)
		}

		req, err = utils.NewRequest("POST", deployURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return err
		}

		pf, err := core.NewPatternFile([]byte(designFile))
		if err != nil {
			return ErrParseDesignFile(err)
		}

		s := utils.CreateDefaultSpinner("Applying design "+pf.Name, "")
		s.Start()
		res, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}

		defer func() { _ = res.Body.Close() }()
		body, err := io.ReadAll(res.Body)
		s.Stop()
		if err != nil {
			return utils.ErrReadFromBody(err)
		}

		if res.StatusCode == 200 {
			utils.Log.Info("design applied")
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
		fmt.Printf("%s\n", a.PatternFile)
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
