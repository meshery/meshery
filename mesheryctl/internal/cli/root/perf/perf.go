// Copyright 2020 Layer5, Inc.
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

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
	outputFormatFlag     string
	tokenPath            string
)

// PerfCmd represents the Performance Management CLI command
var PerfCmd = &cobra.Command{
	Use:     "perf",
	Short:   "Performance Management",
	Long:    `Performance Management & Benchmarking using Meshery CLI.`,
	Example: "mesheryctl perf --name \"a quick stress test\" --url http://192.168.1.15/productpage --qps 300 --concurrent-requests 2 --duration 30s --token \"provider=Meshery\"",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		//Check prerequisite
		return system.RunPreflightHealthChecks(true, cmd.Use)
	},
}

func init() {
	PerfCmd.PersistentFlags().StringVarP(&tokenPath, "token", "t", "", "(required) Path to meshery auth config")
	PerfCmd.PersistentFlags().StringVarP(&outputFormatFlag, "output-format", "o", "", "(optional) format to display in [json|yaml]")

	availableSubcommands = []*cobra.Command{listCmd, viewCmd, applyCmd}
	PerfCmd.AddCommand(availableSubcommands...)
}
