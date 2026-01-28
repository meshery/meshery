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

	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/environment"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands = []*cobra.Command{listEnvironmentCmd, createEnvironmentCmd, deleteEnvironmentCmd, viewEnvironmentCmd}
	environmentApiPath   = "api/environments"
)

var EnvironmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "Manage environments",
	Long: `Create, delete, list of view details of environment(s) of a specific organization
Documentation for environment can be found at https://docs.meshery.io/concepts/logical/environments
	`,
	Example: `
// Create an environment in an organization
mesheryctl environment create --orgID [orgID] --name [name] --description [description]

// Delete an environment in an organization
mesheryctl environment delete environment-id

// List of registered environments in an organization
mesheryctl environment list --orgID [orgID]

// View a particular environment
mesheryctl environment view --orgID [orgID]
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New("Please provide a subcommand with the command"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return utils.ErrInvalidArgument(errors.New(utils.EnvironmentSubError(fmt.Sprintf("'%s' is an invalid command. Use 'mesheryctl environment --help' to display usage guide", args[0]), "environment")))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = cmd.Usage()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	EnvironmentCmd.AddCommand(availableSubcommands...)
}

// selectComponentPrompt lets user to select a model if models are more than one
func selectEnvironmentPrompt(environments []environment.Environment) environment.Environment {
	environmentNames := []string{}
	environmentArray := environments

	for _, environment := range environmentArray {
		environmentName := fmt.Sprintf("ID: %s, Name: %s, Owner: %s, Organization: %s", environment.ID, environment.Name, environment.Owner, environment.OrganizationID)
		environmentNames = append(environmentNames, environmentName)
	}

	prompt := promptui.Select{
		Label: "Select environment",
		Items: environmentNames,
	}

	for {
		i, _, err := prompt.Run()
		if err != nil {
			continue
		}

		return environmentArray[i]
	}
}
