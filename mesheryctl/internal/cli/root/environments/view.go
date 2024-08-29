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

package environments

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/components"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models/environments"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// represents the mesheryctl exp environment view [orgId] subcommand.
var viewEnvironmentCmd = &cobra.Command{
	Use:   "view",
	Short: "view registered environmnents",
	Long:  "view a environments registered in Meshery Server",
	Example: `
// View details of a specific environment
mesheryctl exp environment view --orgID [orgId]
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		orgIdFlag, _ := cmd.Flags().GetString("orgId")

		if orgIdFlag == "" {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New("Please provide a --orgId flag"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()

		orgid := args[0]
		url := fmt.Sprintf("%s/api/environments?orgID=%s", baseUrl, orgid)
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		environmentResponse := &environments.EnvironmentPage{}
		err = json.Unmarshal(data, environmentResponse)
		if err != nil {
			return err
		}

		var selectedEnvironment environments.EnvironmentData

		if environmentResponse.TotalCount == 0 {
			utils.Log.Info("No environment(s) found for the given ID: ", orgid)
			return nil
		} else if environmentResponse.TotalCount == 1 {
			selectedEnvironment = environmentResponse.Environments[0] // Update the type of selectedModel
		} else {
			selectedEnvironment = selectEnvironmentPrompt(environmentResponse.Environments)
		}

		var output []byte

		// user may pass flag in lower or upper case but we have to keep it lower
		// in order to make it consistent while checking output format
		outFormatFlag = strings.ToLower(outFormatFlag)

		if outFormatFlag != "json" && outFormatFlag != "yaml" {
			return utils.ErrInvalidArgument(errors.New("output-format choice is invalid or not provided, use [json|yaml]"))
		}
		// Get the home directory of the user to save the output file
		homeDir, _ := os.UserHomeDir()
		componentString := strings.ReplaceAll(fmt.Sprintf("%v", selectedEnvironment.Name), " ", "_")

		if outFormatFlag == "yaml" || outFormatFlag == "yml" {
			if output, err = yaml.Marshal(selectedEnvironment); err != nil {
				return utils.ErrMarshal(errors.Wrap(err, "failed to format output in YAML"))
			}
			if saveFlag {
				utils.Log.Info("Saving output as YAML file")
				err = os.WriteFile(homeDir+"/.meshery/component_"+componentString+".yaml", output, 0666)
				if err != nil {
					return utils.ErrMarshal(errors.Wrap(err, "failed to save output as YAML file"))
				}
				utils.Log.Info("Output saved as YAML file in ~/.meshery/component_" + componentString + ".yaml")
			} else {
				utils.Log.Info(string(output))
			}
		} else if outFormatFlag == "json" {
			if saveFlag {
				utils.Log.Info("Saving output as JSON file")
				output, err = json.MarshalIndent(selectedEnvironment, "", "  ")
				if err != nil {
					return utils.ErrMarshal(errors.Wrap(err, "failed to format output in JSON"))
				}
				err = os.WriteFile(homeDir+"/.meshery/component_"+componentString+".json", output, 0666)
				if err != nil {
					return utils.ErrMarshal(errors.Wrap(err, "failed to save output as JSON file"))
				}
				utils.Log.Info("Output saved as JSON file in ~/.meshery/component_" + componentString + ".json")
				return nil
			}
			return components.OutputJson(selectedEnvironment)
		} else {
			return utils.ErrInvalidArgument(errors.New("output-format choice is invalid or not provided, use [json|yaml]"))
		}

		return nil
	},
}
