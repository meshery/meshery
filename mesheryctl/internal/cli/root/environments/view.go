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
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta3/environment"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type environmentViewFlags struct {
	orgId        string
	outputFormat string
	save         bool
}

var environmentViewFlagsProvided environmentViewFlags

func formatEnvironmentLabel(rows []environment.Environment) []string {
	labels := []string{}
	for _, e := range rows {
		labels = append(labels, fmt.Sprintf("%s (ID: %s)", e.Name, e.ID.String()))
	}
	return labels
}

var viewEnvironmentCmd = &cobra.Command{
	Use:   "view [environment-name|environment-id]",
	Short: "View registered environments",
	Long: `View details of an environment registered in Meshery Server for a specific organization
Find more information at: https://docs.meshery.io/reference/references/mesheryctl/environment/view`,
	Example: `
// View details of a specific environment by ID
mesheryctl environment view [environment-id] --orgId [orgId]

// View details of an environment by name (prompts if several match)
mesheryctl environment view [environment-name] --orgId [orgId]

// Select from all environments of an organization
mesheryctl environment view --orgId [orgId]
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			const errMsg = "please provide at most one environment name or ID\n\nUsage: mesheryctl environment view [environment-name|environment-id] --orgId [orgId]\nRun 'mesheryctl environment view --help' to see detailed help message"
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if environmentViewFlagsProvided.orgId == "" {
			const errMsg = "[ orgId ] isn't specified\n\nUsage: mesheryctl environment view [environment-name|environment-id] --orgId [orgId]\nRun 'mesheryctl environment view --help' to see detailed help message"
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}

		if !utils.IsUUID(environmentViewFlagsProvided.orgId) {
			return utils.ErrInvalidUUID(fmt.Errorf("invalid orgId: %s", environmentViewFlagsProvided.orgId))
		}

		return display.ValidateOutputFormat(environmentViewFlagsProvided.outputFormat)
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		orgQuery := url.Values{}
		orgQuery.Set("orgId", environmentViewFlagsProvided.orgId)

		var selectedEnvironment environment.Environment

		if len(args) == 1 && utils.IsUUID(args[0]) {
			urlPath := fmt.Sprintf("%s/%s?%s", environmentApiPath, url.PathEscape(args[0]), orgQuery.Encode())
			fetchedEnvironment, err := api.Fetch[environment.Environment](urlPath)
			if err != nil {
				return err
			}
			// Not every provider scopes the lookup by orgId, so check it here.
			if !strings.EqualFold(fetchedEnvironment.OrganizationID.String(), environmentViewFlagsProvided.orgId) {
				return utils.ErrNotFound(fmt.Errorf("No environment found with ID %s in organization: %s", args[0], environmentViewFlagsProvided.orgId))
			}
			selectedEnvironment = *fetchedEnvironment
		} else {
			searchTerm := ""
			notFoundMsg := fmt.Sprintf("No environment(s) found in organization: %s", environmentViewFlagsProvided.orgId)
			if len(args) == 1 {
				searchTerm = args[0]
				notFoundMsg = fmt.Sprintf("No environment(s) found with name: %s", searchTerm)
			}

			err := display.PromptAsyncPagination(
				display.DisplayDataAsync{
					UrlPath:        fmt.Sprintf("%s?%s", environmentApiPath, orgQuery.Encode()),
					SearchTerm:     searchTerm,
					ErrNotFoundMsg: notFoundMsg,
				},
				formatEnvironmentLabel,
				func(data *environment.EnvironmentPage) ([]environment.Environment, int64) {
					return data.Environments, int64(data.TotalCount)
				},
				&selectedEnvironment,
			)
			if err != nil {
				return err
			}
		}

		outputFormat := strings.ToLower(environmentViewFlagsProvided.outputFormat)

		outputFormatterFactory := display.OutputFormatterFactory[environment.Environment]{}
		outputFormatter, err := outputFormatterFactory.New(outputFormat, selectedEnvironment)
		if err != nil {
			return err
		}
		outputFormatter = outputFormatter.WithOutput(cmd.OutOrStdout())

		err = outputFormatter.Display()
		if err != nil {
			return err
		}

		// Get the home directory of the user to save the output file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return utils.ErrRetrieveHomeDir(errors.Wrap(err, "failed to determine user home directory"))
		}
		environmentString := strings.ReplaceAll(fmt.Sprintf("%v", selectedEnvironment.Name), " ", "_")

		if environmentViewFlagsProvided.save {
			fileName := fmt.Sprintf("environment_%s.%s", environmentString, strings.ToLower(environmentViewFlagsProvided.outputFormat))
			file := filepath.Join(homeDir, ".meshery", fileName)
			outputFormatterSaverFactory := display.OutputFormatterSaverFactory[environment.Environment]{}
			outputFormatterSaver, err := outputFormatterSaverFactory.New(environmentViewFlagsProvided.outputFormat, outputFormatter)
			if err != nil {
				return err
			}

			outputFormatterSaver = outputFormatterSaver.WithFilePath(file)
			err = outputFormatterSaver.Save()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	viewEnvironmentCmd.Flags().StringVarP(&environmentViewFlagsProvided.outputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewEnvironmentCmd.Flags().BoolVarP(&environmentViewFlagsProvided.save, "save", "s", false, "(optional) save output as a JSON/YAML file")
	viewEnvironmentCmd.Flags().StringVarP(&environmentViewFlagsProvided.orgId, "orgId", "", "", "Organization ID")
}
