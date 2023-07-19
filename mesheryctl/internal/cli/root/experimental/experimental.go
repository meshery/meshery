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

package experimental

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
)

// ExpCmd represents the Performance Management CLI command
var ExpCmd = &cobra.Command{
	Use:   "exp",
	Short: "Experimental commands for mesheryctl",
	Long:  `List of experimental commands for testing and evaluation purpose.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.ExpError(fmt.Sprintf("'%s' is a invalid command. Use 'mesheryctl exp --help' to display usage guide.'\n", args[0])))
		}
		return nil
	},
}

func init() {
	// availableSubcommands = []*cobra.Command{filter.FilterCmd}
	// ExpCmd.AddCommand(availableSubcommands...)
}
