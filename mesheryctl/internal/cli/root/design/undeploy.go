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
	"os"
	"path/filepath"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/patterns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var designUndeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploy design",
	Long: `Undeploy design will trigger undeploy of design.
	With -f, Meshery imports and saves the design file before removal, as it does with design deploy -f.
	Find more information at: https://docs.meshery.io/reference/mesheryctl/design/undeploy`,
	Example: `
// Undeploy design by providing file path
mesheryctl design undeploy -f [filepath]

// Undeploy design without an interactive prompt
mesheryctl design undeploy -f [filepath] --yes
	`,

	Args: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("file") && strings.TrimSpace(file) == "" {
			errMsg := `Usage: mesheryctl design undeploy -f [filepath]`
			return ErrUndeployDesign(fmt.Errorf("%s", errMsg))
		}

		if len(args) == 0 && file == "" {
			return ErrUndeployDesign(
				fmt.Errorf("provide either a design ID or -f [filepath]"),
			)
		}

		if len(args) > 0 && strings.TrimSpace(args[0]) == "" {
			return utils.ErrInvalidArgument(fmt.Errorf("target cannot be empty"))
		}

		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("expected at most one positional argument, got %d", len(args)))
		}

		if len(args) > 0 && file != "" {
			return utils.ErrInvalidArgument(fmt.Errorf("cannot specify both a positional argument and the --file flag"))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		design := ""
		isID := false
		if len(args) > 0 {
			design, isID, err = utils.ValidId(mctlCfg.GetBaseMesheryURL(), args[0], "pattern")
			if err != nil {
				return err
			}
		}

		// Delete the design using the id
		if isID {
			err := utils.DeleteConfiguration(mctlCfg.GetBaseMesheryURL(), design, "pattern")
			if err != nil {
				return ErrDeleteDesign(err, args[0])
			}
			utils.Log.Infof("design %s deleted", args[0])
			return nil
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/deploy?" + utils.AllK8sContextsQuery

		if !utils.IsValidUrl(file) {
			content, err := os.ReadFile(file)
			if err != nil {
				return utils.ErrFileRead(err)
			}
			designFile = string(content)
		} else {
			return utils.ErrInvalidArgument(fmt.Errorf("URLs are not currently supported"))
		}

		confirmed, err := confirmClusterRemoval(file, []byte(designFile))
		if err != nil {
			return err
		}
		if !confirmed {
			utils.Log.Info("Undeploy cancelled")
			return nil
		}

		importURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/import"
		fileName := filepath.Base(file)
		patternName := name
		if patternName == "" {
			patternName = fileName
		}
		imported, err := importDesignFile(importURL, patternName, fileName, []byte(designFile))
		if err != nil {
			return err
		}

		pf, err := patterns.GetPatternFormat(imported.PatternFile)
		if err != nil {
			return ErrParseDesignFile(err)
		}

		pfByt, err := yaml.Marshal(pf)
		if err != nil {
			return utils.ErrMarshal(err)
		}

		payloadBytes, err := json.Marshal(models.MesheryPatternFileDeployPayload{
			PatternFile: string(pfByt),
		})
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

		if res.StatusCode == 200 {
			utils.Log.Info("design undeployed")
		}
		utils.Log.Info(string(body))

		return nil
	},
}

func init() {
	designUndeployCmd.Flags().BoolVarP(&utils.SilentFlag, "yes", "y", false, "(optional) assume yes for user interactive prompts.")
	designUndeployCmd.Flags().StringVarP(&file, "file", "f", "", "Path to design file")
}
