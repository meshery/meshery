// Copyright 2024 Layer5, Inc.
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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models/environments"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	name          string
	description   string
	orgID         string
	outFormatFlag string
	saveFlag      bool

	maxRowsPerPage       = 25
	whiteBoardPrinter    = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
	availableSubcommands = []*cobra.Command{listEnvironmentCmd, CreateEnvironmentCmd, DeleteEnvironmentCmd, viewEnvironmentCmd}
)

var EnvironmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "View list of environments and detail of environments",
	Long:  "View list of environments and detailed information of a specific environments",
	Example: `
// To view a list environments
mesheryctl exp environment list [orgId]
// To create a environment
mesheryctl exp environment create [orgId]
// Documentation for environment can be found at:
https://docs.layer5.io/cloud/spaces/environments/
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return system.ErrGetCurrentContext(err)
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return errors.New(utils.EnvironmentSubError("Please provide a orgID", "environment"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.EnvironmentSubError(fmt.Sprintf("'%s' is an invalid command. Use 'mesheryctl exp environment --help' to display usage guide.'\n", args[0]), "environment"))
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
	viewEnvironmentCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewEnvironmentCmd.Flags().BoolVarP(&saveFlag, "save", "s", false, "(optional) save output as a JSON/YAML file")
	EnvironmentCmd.AddCommand(availableSubcommands...)
}

// selectComponentPrompt lets user to select a model if models are more than one
func selectEnvironmentPrompt(environment []environments.EnvironmentData) environments.EnvironmentData {
	environmentNames := []string{}
	environmentArray := []environments.EnvironmentData{}

	environmentArray = append(environmentArray, environment...)

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
