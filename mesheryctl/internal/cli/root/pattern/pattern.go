// Copyright 2023 Layer5, Inc.
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

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
	file                 string
)

// PatternCmd represents the root command for pattern commands
var PatternCmd = &cobra.Command{
	Use:   "pattern",
	Short: "Service Mesh Patterns Management",
	Long:  `Manage service meshes using predefined patterns`,
	Example: `
// Apply pattern file
mesheryctl pattern apply --file [path to pattern file | URL of the file]

// Delete pattern file
mesheryctl pattern delete --file [path to pattern file]

// View pattern file
mesheryctl pattern view [pattern name | ID]

// List all patterns
mesheryctl pattern list
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.PatternError(fmt.Sprintf("'%s' is a invalid command.  Use 'mesheryctl pattern --help' to display usage guide.\n", args[0])))
		}
		return nil
	},
}

func init() {
	PatternCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token file default from current context")

	availableSubcommands = []*cobra.Command{applyCmd, deleteCmd, viewCmd, listCmd}
	PatternCmd.AddCommand(availableSubcommands...)
}
