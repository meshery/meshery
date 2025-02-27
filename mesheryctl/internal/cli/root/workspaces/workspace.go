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

package workspaces

import (
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
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

	maxRowsPerPage       = 25
	whiteBoardPrinter    = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
	availableSubcommands = []*cobra.Command{listWorkspaceCmd, createWorkspaceCmd}
)

var WorkSpaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "View list of workspaces and detail of workspaces",
	Long:  "View list of workspaces and detailed information of a specific workspaces",
	Example: `

// To view a list workspaces
mesheryctl exp workspace list --orgId [orgId]

// To create a workspace
mesheryctl exp workspace create --orgId [orgId] --name [name] --description [description]

// Documentation for workspace can be found at:
https://docs.layer5.io/cloud/spaces/workspaces/
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if err := cmd.Usage(); err != nil {
				return nil
			}
			return utils.ErrInvalidArgument(errors.New("Please provide a subcommand"))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return utils.ErrInvalidArgument(cmd.Usage())
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
	listWorkspaceCmd.Flags().StringVarP(&orgID, "orgId", "o", "", "Organization ID")
	createWorkspaceCmd.Flags().StringVarP(&orgID, "orgId", "o", "", "Organization ID")
	createWorkspaceCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the workspace")
	createWorkspaceCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the workspace")

	err := listWorkspaceCmd.MarkFlagRequired("orgId")
	if err != nil {
		utils.Log.Info(err)
	}
	err = createWorkspaceCmd.MarkFlagRequired("orgId")
	if err != nil {
		utils.Log.Info(err)
	}
	err = createWorkspaceCmd.MarkFlagRequired("name")
	if err != nil {
		utils.Log.Info(err)
	}
	err = createWorkspaceCmd.MarkFlagRequired("description")
	if err != nil {
		utils.Log.Info(err)
	}
	WorkSpaceCmd.AddCommand(availableSubcommands...)
}
