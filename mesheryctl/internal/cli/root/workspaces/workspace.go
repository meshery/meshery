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
	"fmt"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	availableSubcommands = []*cobra.Command{listWorkspaceCmd, createWorkspaceCmd}
	workspacesApiPath    = "api/workspaces"
)

var WorkSpaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Managge workspaces under an organization",
	Long: `Create, list of workspaces under an organization
Find more information at: https://docs.meshery.io/reference/mesheryctl/exp/workspace`,
	Example: `

// To view a list workspaces
mesheryctl exp workspace list --orgId [orgId]

// To create a workspace
mesheryctl exp workspace create --orgId [orgId] --name [name] --description [description]
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || len(args) > 1 {
			return returnInvalidArgumentProvidedError()
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return utils.ErrInvalidArgument(fmt.Errorf("'%s' is an invalid subcommand. Please provide required options from [create, list]. Use 'mesheryctl exp workspace --help' to display usage guide", args[0]))
		}

		return nil
	},
}

func returnInvalidArgumentProvidedError() error {
	errMsg := "Usage: mesheryctl exp workspace [subcommand]\nRun 'mesheryctl exp workspace --help' to see detailed help message"
	return utils.ErrInvalidArgument(fmt.Errorf("no subcommand specified. %s", errMsg))
}

func init() {
	WorkSpaceCmd.Flags().BoolP("count", "", false, "total number of registered workspaces")
	WorkSpaceCmd.AddCommand(availableSubcommands...)
}
