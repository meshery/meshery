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

package filter

import (
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
	// for formatting errors
	cmdUsed string
)

// FilterCmd represents the root command for filter commands
var FilterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Manage WebAssembly filters",
	Long: `Cloud Native Filter Management
Find more information at: https://docs.meshery.io/reference/mesheryctl#command-reference`,
	Example: `
// Base command for WASM filters:
mesheryctl filter [subcommands]
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return ErrInvalidFilterCommand(args[0])
		}
		return nil
	},
}

func init() {
	FilterCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token file default from current context")
	availableSubcommands = []*cobra.Command{viewCmd, deleteCmd, listCmd, importCmd}

	FilterCmd.AddCommand(availableSubcommands...)
}
