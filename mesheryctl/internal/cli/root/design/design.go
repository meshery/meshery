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

package design

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
	file                 string
	validSourceTypes     []string
)

// DesignCmd represents the root command for design commands
var DesignCmd = &cobra.Command{
	Use:   "design",
	Short: "Cloud Native Designs Management",
	Long: `Manage cloud and cloud native infrastructure using predefined designs.
Find more information at: https://docs.meshery.io/reference/mesheryctl#command-reference`,
	Example: `
// Apply design file:
mesheryctl design apply --file [path to design file | URL of the file]

// Delete design file:
mesheryctl design delete --file [path to design file]

// View design file:
mesheryctl design view [design name | ID]

// List all designs:
mesheryctl design list
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			suggestions := make([]string, 0)
			for _, subcmd := range availableSubcommands {
				if strings.HasPrefix(subcmd.Name(), args[0]) {
					suggestions = append(suggestions, subcmd.Name())
				}
			}
			if len(suggestions) > 0 {
				return errors.New(utils.DesignError(fmt.Sprintf("'%s' is an invalid command. \nDid you mean %v? \nUse 'mesheryctl design --help' to display usage guide.\n", args[0], suggestions)))
			}
			return errors.New(utils.DesignError(fmt.Sprintf("'%s' is an invalid command. Use 'mesheryctl design --help' to display usage guide.\n", args[0])))
		}
		return nil
	},
}

func init() {
	DesignCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token file default from current context")

	availableSubcommands = []*cobra.Command{applyCmd, deleteCmd, viewCmd, listCmd, importCmd, onboardCmd, offboardCmd}
	DesignCmd.AddCommand(availableSubcommands...)
}
