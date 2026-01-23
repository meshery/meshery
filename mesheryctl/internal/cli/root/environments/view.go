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
	"slices"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/environment"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type environmentViewFlags struct {
	orgID        string
	outputFormat string
	save         bool
}

var environmentViewFlagsProvided environmentViewFlags
var environmentViewOutputFormats = []string{"json", "yaml"}

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

		if environmentViewFlagsProvided.orgID == "" {
			const errMsg = "[ orgID ] isn't specified\n\nUsage: mesheryctl environment view --orgID [orgID]\nRun 'mesheryctl environment view --help' to see detailed help message"
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}

		if !slices.Contains(environmentViewOutputFormats, strings.ToLower(environmentViewFlagsProvided.outputFormat)) {
			return utils.ErrInvalidArgument(errors.New("output-format choice is invalid or not provided, use [json|yaml]"))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		environmentResponse, err := api.Fetch[environment.EnvironmentPage](fmt.Sprintf("api/environments?orgID=%s", environmentViewFlagsProvided.orgID))

		if err != nil {
			return err
		}

		var selectedEnvironment environment.Environment

		switch environmentResponse.TotalCount {
		case 0:
			utils.Log.Info("No environment(s) found for the given ID: ", environmentViewFlagsProvided.orgID)
			return nil
		case 1:
			selectedEnvironment = environmentResponse.Environments[0] // Update the type of selectedModel
		default:
			selectedEnvironment = selectEnvironmentPrompt(environmentResponse.Environments)
		}

		var output []byte

		// user may pass flag in lower or upper case but we have to keep it lower
		// in order to make it consistent while checking output format
		outputFormat := strings.ToLower(environmentViewFlagsProvided.outputFormat)

		outputFormatterFactory := display.OutputFormatterFactory[environment.Environment]{}
		outputFormatter, err := outputFormatterFactory.New(outputFormat, selectedEnvironment)
		if err != nil {
			return err
		}

		err = outputFormatter.Display()
		if err != nil {
			return err
		}

		// Get the home directory of the user to save the output file
		homeDir, _ := os.UserHomeDir()
		componentString := strings.ReplaceAll(fmt.Sprintf("%v", selectedEnvironment.Name), " ", "_")

		// TODO: Add support for YAML/JSON output saving in outputFormatter itself
		switch outputFormat {
		case "yaml", "yml":
			if output, err = yaml.Marshal(selectedEnvironment); err != nil {
				return utils.ErrMarshal(errors.Wrap(err, "failed to format output in YAML"))
			}
			if environmentViewFlagsProvided.save {
				utils.Log.Info("Saving output as YAML file")
				err = os.WriteFile(homeDir+"/.meshery/component_"+componentString+".yaml", output, 0666)
				if err != nil {
					return utils.ErrMarshal(errors.Wrap(err, "failed to save output as YAML file"))
				}
				utils.Log.Info("Output saved as YAML file in ~/.meshery/component_" + componentString + ".yaml")
			}
		case "json":
			if environmentViewFlagsProvided.save {
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
		}

		return nil
	},
}

func init() {
	viewEnvironmentCmd.Flags().StringVarP(&environmentViewFlagsProvided.outputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewEnvironmentCmd.Flags().BoolVarP(&environmentViewFlagsProvided.save, "save", "s", false, "(optional) save output as a JSON/YAML file")
	viewEnvironmentCmd.Flags().StringVarP(&environmentViewFlagsProvided.orgID, "orgID", "", "", "Organization ID")
}
