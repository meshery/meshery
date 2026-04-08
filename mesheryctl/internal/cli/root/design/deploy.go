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
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/patterns"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type cmdDesignDeployFlags struct {
	File       string `json:"file" validate:"omitempty,file"`
	SourceType string `json:"source-type" validate:"omitempty,design-source-type"`
	SkipSave   bool   `json:"skip-save" validate:"boolean"`
}

var designDeployFlags cmdDesignDeployFlags

var linkDocDesignDeploy = map[string]string{
	"link":    "![pattern-onboard-usage](/reference/images/pattern-onboard.png)",
	"caption": "Usage of mesheryctl design deploy",
}

var deployDesignCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy design",
	Long:  `Command will trigger deploy of design`,
	Example: `
// Deploy design by providing file path
mesheryctl design deploy -f [filepath] -s [source type]
	`,
	Annotations: linkDocDesignDeploy,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		flagValidator := mesheryctlflags.GetFlagValidator()
		designDeployValidSourceTypes, err := getDesignSourceTypes()
		if err != nil {
			return err
		}

		err = flagValidator.Validator.RegisterValidation("design-source-type", func(fl validator.FieldLevel) bool {
			if sourceType, ok := fl.Field().Interface().(string); ok {
				for _, validType := range designDeployValidSourceTypes {
					if strings.EqualFold(sourceType, validType) {
						return true
					}
				}
			}
			return false
		})
		if err != nil {
			return err
		}

		flagValidator.CustomErrors["design-source-type"] = fmt.Sprintf("Invalid value for --source-type '%v': valid values are %s", designDeployFlags.SourceType, strings.Join(designDeployValidSourceTypes, ", "))

		return flagValidator.Validate(&designDeployFlags)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if designDeployFlags.File == "" && len(args) == 0 {
			return ErrDeployDesign()
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var patternFile *pattern.PatternFile

		deployURLPath := "api/pattern/deploy"
		patternURLPath := "api/pattern"

		// pattern name has been passed
		if len(args) > 0 {
			// Merge args to get pattern-name
			patternName := strings.Join(args, " ")

			queryParams := url.Values{}
			queryParams.Set("populate", "pattern_file")
			queryParams.Set("search", url.QueryEscape(patternName))
			urlPath := fmt.Sprintf("%s?%s", patternURLPath, queryParams.Encode())
			// search and fetch patterns with pattern-name
			utils.Log.Debug("Fetching designs")

			response, err := api.Fetch[models.PatternsAPIResponse](urlPath)
			if err != nil {
				return err
			}

			index := 0
			if len(response.Patterns) == 0 {
				return ErrDesignNotFound(patternName)
			} else if len(response.Patterns) == 1 {
				patternFile, _ = patterns.GetPatternFormat(response.Patterns[0].PatternFile)
			} else {
				// Multiple patterns with same name
				index = multiplepatternsConfirmation(response.Patterns)
				patternFile, _ = patterns.GetPatternFormat(response.Patterns[index].PatternFile)
			}
		} else if designDeployFlags.SkipSave {
			patternFile, err = readPatternFromFile(designDeployFlags.File)
			if err != nil {
				return err
			}
		} else {
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return err
			}
			patternImportURL := fmt.Sprintf("%s/%s/import", mctlCfg.GetBaseMesheryURL(), patternURLPath)
			pattern, err := importPattern(designDeployFlags.SourceType, designDeployFlags.File, patternImportURL, !designDeployFlags.SkipSave)
			if err != nil {
				return err
			}

			patternFile, _ = patterns.GetPatternFormat(pattern.PatternFile)
		}

		patternFileByt, _ := yaml.Marshal(patternFile)
		payload := models.MesheryPatternFileDeployPayload{
			PatternFile: string(patternFileByt),
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return utils.ErrMarshal(err)
		}

		res, err := api.Add(deployURLPath, bytes.NewBuffer(payloadBytes), nil)
		if err != nil {
			if errors.GetCode(err) == utils.ErrMesheryServerInternalErrorCode {
				return ErrDeployDesign()
			}
			return err
		}

		defer func() { _ = res.Body.Close() }()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return utils.ErrReadFromBody(err)
		}

		if res.StatusCode == 200 {
			utils.Log.Info("design deployed")
		}
		utils.Log.Info(string(body))
		return nil

	},
}

func multiplepatternsConfirmation(profiles []models.MesheryPattern) int {
	reader := bufio.NewReader(os.Stdin)

	patternFileByt, _ := yaml.Marshal(designFile)

	for index, a := range profiles {
		fmt.Printf("Index: %v\n", index)
		fmt.Printf("Name: %v\n", a.Name)
		fmt.Printf("ID: %s\n", a.ID.String())
		fmt.Printf("patternFile:\n")
		fmt.Println(string(patternFileByt))
		fmt.Println("---------------------")
	}

	for {
		fmt.Printf("Enter the index of design: ")
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
	deployDesignCmd.Flags().StringVarP(&designDeployFlags.File, "file", "f", "", "Path to design file")
	deployDesignCmd.Flags().BoolVarP(&designDeployFlags.SkipSave, "skip-save", "", false, "Skip saving a design")
	deployDesignCmd.Flags().StringVarP(&designDeployFlags.SourceType, "source-type", "s", "", "Type of source file (ex. manifest / compose / helm)")
}
