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

package pattern

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

// PatternCmd represents the root command for pattern commands
var PatternCmd = &cobra.Command{
	Use:   "pattern",
	Short: "Cloud Native Patterns Management",
	Long: `Manage cloud and cloud native infrastructure using predefined patterns.
Find more information at: https://docs.meshery.io/reference/mesheryctl#command-reference`,
	Example: `
// Apply pattern file:
mesheryctl pattern apply --file [path to pattern file | URL of the file]

// Delete pattern file:
mesheryctl pattern delete --file [path to pattern file]

// View pattern file:
mesheryctl pattern view [pattern name | ID]

// List all patterns:
mesheryctl pattern list
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
				return errors.New(utils.PatternError(fmt.Sprintf("'%s' is an invalid command. \nDid you mean %v? \nUse 'mesheryctl pattern --help' to display usage guide.\n", args[0], suggestions)))
			}
			return errors.New(utils.PatternError(fmt.Sprintf("'%s' is an invalid command. Use 'mesheryctl pattern --help' to display usage guide.\n", args[0])))
		}
		return nil
	},
}

func init() {
	PatternCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token file default from current context")

	availableSubcommands = []*cobra.Command{applyCmd, deleteCmd, viewCmd, listCmd, importCmd, onboardCmd, offboardCmd, exportCmd}
	PatternCmd.AddCommand(availableSubcommands...)
}
