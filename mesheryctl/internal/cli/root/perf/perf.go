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

package perf

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
	outputFormatFlag     string
	// setting up for error formatting
	cmdUsed string
)

// PerfCmd represents the Performance Management CLI command
var PerfCmd = &cobra.Command{
	Use:   "perf",
	Short: "Performance Management",
	Long:  `Performance Management & Benchmarking`,
	Example: `
// Run performance test
mesheryctl perf apply test-3 --name "a quick stress test" --url http://192.168.1.15/productpage --qps 300 --concurrent-requests 2 --duration 30s
	
// List performance profiles
mesheryctl perf profile sam-test

// List performance results
mesheryctl perf result sam-test

// Display Perf profile in JSON or YAML
mesheryctl perf result -o json
mesheryctl perf result -o yaml
	`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			availableSubCmds := []string{"apply", "profile", "result"}

			suggestedCmd := utils.FindClosestArg(args[0], availableSubCmds)
			if suggestedCmd != "" && suggestedCmd[0] == args[0][0] {
				return errors.New(utils.PerfError(fmt.Sprintf("'%s' is a invalid command for '%s'. Did you mean this?\n\t%s\n", args[0], cmd.CalledAs(), suggestedCmd)))
			}
			return errors.New(utils.PerfError(fmt.Sprintf("'%s' is a invalid command for '%s'. Use 'mesheryctl perf --help' to display usage guide.\n", args[0], cmd.CalledAs())))
		}
		return nil
	},
}

func init() {
	PerfCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "(required) Path to meshery auth config")
	PerfCmd.PersistentFlags().StringVarP(&outputFormatFlag, "output-format", "o", "", "(optional) format to display in [json|yaml]")
	PerfCmd.PersistentFlags().BoolVarP(&utils.SilentFlag, "yes", "y", false, "(optional) assume yes for user interactive prompts.")

	availableSubcommands = []*cobra.Command{profileCmd, resultCmd, applyCmd}
	PerfCmd.AddCommand(availableSubcommands...)
}
