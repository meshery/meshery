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

package filter

import (
	"fmt"
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands []*cobra.Command
)

// FilterCmd represents the root command for filter commands
var FilterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Cloud Native Filter Management",
	Long:  ``,
	Example: `
// Base command for WASM filters
mesheryctl filter [subcommands]	
	`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.HasSubCommands() {
			cmd.Help()
			os.Exit(0)
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		running, _ := utils.IsMesheryRunning(currCtx.GetPlatform())
		if !running {
			return errors.New(`meshery server is not running. run "mesheryctl system start" to start meshery`)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.FilterError(fmt.Sprintf("'%s' is a invalid command.  Use 'mesheryctl filter --help' to display usage guide.\n", args[0])))
		}
		return nil
	},
}

func init() {
	FilterCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token file default from current context")
	availableSubcommands = []*cobra.Command{viewCmd, deleteCmd, listCmd, importCmd}

	FilterCmd.AddCommand(availableSubcommands...)
}
