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
	"os"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/mesheryctl/pkg/utils/format"
	"github.com/meshery/meshery/server/models/environments"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// represents the mesheryctl environment view [orgID] subcommand.
var viewEnvironmentCmd = &cobra.Command{
	Use:   "view",
	Short: "View registered environmnents",
	Long: `View details of an environment registered in Meshery Server
Documentation for environment can be found at https://docs.meshery.io/reference/mesheryctl/environment/view`,
	Example: `
// View details of a specific environment
mesheryctl environment view --orgID [orgID]
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		orgIDFlag, _ := cmd.Flags().GetString("orgID")

		if orgIDFlag == "" {
			const errMsg = "[ orgID ] isn't specified\n\nUsage: mesheryctl environment view --orgID [orgID]\nRun 'mesheryctl environment view --help' to see detailed help message"
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		orgID, _ := cmd.Flags().GetString("orgID")
		outFormat, _ := cmd.Flags().GetString("output-format")
		save, _ := cmd.Flags().GetBool("save")

		environmentResponse, err := api.Fetch[environments.EnvironmentPage](fmt.Sprintf("api/environments?orgID=%s", orgID))

		if err != nil {
			return err
		}

		var selectedEnvironment environments.EnvironmentData

		if environmentResponse.TotalCount == 0 {
			utils.Log.Info("No environment(s) found for the given ID: ", orgID)
			return nil
		} else if environmentResponse.TotalCount == 1 {
			selectedEnvironment = environmentResponse.Environments[0] // Update the type of selectedModel
		} else {
			selectedEnvironment = selectEnvironmentPrompt(environmentResponse.Environments)
		}

		var output []byte

		// user may pass flag in lower or upper case but we have to keep it lower
		// in order to make it consistent while checking output format
		outFormat = strings.ToLower(outFormat)

		if outFormat != "json" && outFormat != "yaml" {
			return utils.ErrInvalidArgument(errors.New("output-format choice is invalid or not provided, use [json|yaml]"))
		}
		// Get the home directory of the user to save the output file
		homeDir, _ := os.UserHomeDir()
		componentString := strings.ReplaceAll(fmt.Sprintf("%v", selectedEnvironment.Name), " ", "_")

		if outFormat == "yaml" || outFormat == "yml" {
			if output, err = yaml.Marshal(selectedEnvironment); err != nil {
				return utils.ErrMarshal(errors.Wrap(err, "failed to format output in YAML"))
			}
			if save {
				utils.Log.Info("Saving output as YAML file")
				err = os.WriteFile(homeDir+"/.meshery/component_"+componentString+".yaml", output, 0666)
				if err != nil {
					return utils.ErrMarshal(errors.Wrap(err, "failed to save output as YAML file"))
				}
				utils.Log.Info("Output saved as YAML file in ~/.meshery/component_" + componentString + ".yaml")
			} else {
				utils.Log.Info(string(output))
			}
		} else if outFormat == "json" {
			if save {
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
			return format.OutputJson(selectedEnvironment)
		} else {
			return utils.ErrInvalidArgument(errors.New("output-format choice is invalid or not provided, use [json|yaml]"))
		}

		return nil
	},
}

func init() {
	viewEnvironmentCmd.Flags().StringP("output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewEnvironmentCmd.Flags().BoolP("save", "s", false, "(optional) save output as a JSON/YAML file")
	viewEnvironmentCmd.Flags().StringP("orgID", "", "", "Organization ID")
}
