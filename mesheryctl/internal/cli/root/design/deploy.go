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
	pattern "github.com/meshery/schemas/models/v1beta3/design"
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
	"link":    "![pattern-onboard-usage](../../../images/app-onboard.png)",
	"caption": "Usage of mesheryctl design deploy",
}

var deployDesignCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy design",
	Long: `Command will trigger deploy of design.
	Find more information at: https://docs.meshery.io/reference/references/mesheryctl/design/deploy`,
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

			urlPath := getDesignDeploySearchURLPath(patternURLPath, patternName)
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
			pattern, err := importPattern(designDeployFlags.SourceType, designDeployFlags.File, patternImportURL)
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

// getDesignDeploySearchURLPath builds the design lookup URL for a name search.
//
// The name is handed to url.Values unescaped: Encode() percent-encodes each
// value exactly once, and the server reads it back with r.URL.Query(). Escaping
// it beforehand with url.QueryEscape encoded it twice, so "My Design" reached
// the server as the literal "My+Design" and the LIKE never matched - reported
// as `design not found` rather than as an encoding fault. Every name that
// percent-encodes was affected: spaces, "&", "+", and all non-ASCII.
// getDesignViewUrlPath in view.go does the same single encoding.
//
// The name is passed through exactly as the caller built it. Normalising it
// here - trimming, for instance - would widen the LIKE and can turn a
// one-match deploy into a multi-match one, which is a lookup-semantics change
// rather than an encoding fix. See the discussion on #21847.
func getDesignDeploySearchURLPath(baseAPIPath, designName string) string {
	queryParams := url.Values{}
	queryParams.Set("populate", "pattern_file")
	queryParams.Set("search", designName)

	return fmt.Sprintf("%s?%s", baseAPIPath, queryParams.Encode())
}

func init() {
	deployDesignCmd.Flags().StringVarP(&designDeployFlags.File, "file", "f", "", "Path to design file")
	deployDesignCmd.Flags().BoolVarP(&designDeployFlags.SkipSave, "skip-save", "", false, "Skip saving a design")
	deployDesignCmd.Flags().StringVarP(&designDeployFlags.SourceType, "source-type", "s", "", "Type of source file (ex. manifest / compose / helm)")
}
