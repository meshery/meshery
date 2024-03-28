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

package workspaces

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	name        string
	description string
	orgID       string
	page        int

	maxRowsPerPage       = 25
	whiteBoardPrinter    = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
	availableSubcommands = []*cobra.Command{listWorkspaceCmd, CreateWorkspaceCmd}
)

var WorkSpaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "View list of workspace and detailed information of a specific relationship",
	Long:  "View list of workspace and detailed information of a specific relationship",
	Example: `
// To view list of components
mesheryctl exp workspace list

// To view a specific model
mesheryctl exp workspace view [model-name]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			utils.Log.Error(err)
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(system.ErrGetCurrentContext(err))
			return nil
		}
		err = ctx.ValidateVersion()
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl exp workspace --help' to display usage guide.\n", args[0]), "model"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = cmd.Usage()
		if err != nil {
			utils.Log.Error(err)
			return err
		}
		return nil
	},
}

func init() {
	listWorkspaceCmd.Flags().StringVarP(&orgID, "orgId", "o", "", "Organization ID")
	listWorkspaceCmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	CreateWorkspaceCmd.Flags().StringVarP(&orgID, "orgId", "o", "", "Organization ID")
	CreateWorkspaceCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the workspace")
	CreateWorkspaceCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the workspace")
	WorkSpaceCmd.AddCommand(availableSubcommands...)	
}
