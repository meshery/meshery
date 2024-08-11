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

package experimental

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/connections"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/environments"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/relationships"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/workspaces"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
)

// ExpCmd represents the Experimental commands
var ExpCmd = &cobra.Command{
	Use:   "exp",
	Short: "Experimental commands for mesheryctl",
	Long:  `Commands under the Experimental group are for testing and evaluation prior to promotion to general availability. Experimental commands are subject to change.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.ExpError(fmt.Sprintf("'%s' is an invalid command. Use 'mesheryctl exp --help' to display usage guide.'\n", args[0])))
		}
		return nil
	},
}

func init() {
	availableSubcommands = append(availableSubcommands, connections.ConnectionsCmd, relationships.RelationshipCmd, workspaces.WorkSpaceCmd, environments.EnvironmentCmd)

	ExpCmd.AddCommand(availableSubcommands...)
}
