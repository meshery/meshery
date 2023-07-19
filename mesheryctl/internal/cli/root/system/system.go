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

package system

import (
	"fmt"

	config "github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands = []*cobra.Command{}
	// flag to change the current context to a temporary context
	tempContext = ""
	// token path
)

// SystemCmd represents Meshery Lifecycle Management cli commands
var SystemCmd = &cobra.Command{
	Use:   "system",
	Short: "Meshery Lifecycle Management",
	Long:  `Manage the state and configuration of Meshery server, components, and client.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("'%s' is a invalid command.  Use 'mesheryctl system --help' to display usage guide.\n", args[0])))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		mctlCfg.GetBaseMesheryURL()
		return nil
	},
}

func init() {
	availableSubcommands = []*cobra.Command{
		resetCmd,
		logsCmd,
		startCmd,
		stopCmd,
		restartCmd,
		statusCmd,
		updateCmd,
		configCmd,
		ContextCmd,
		channelCmd,
		providerCmd,
		checkCmd,
		loginCmd,
		logoutCmd,
		tokenCmd,
		dashboardCmd,
		modelCmd,
	}
	// --context flag to temporarily change context. This is global to all system commands
	SystemCmd.PersistentFlags().StringVarP(&tempContext, "context", "c", "", "(optional) temporarily change the current context.")
	SystemCmd.PersistentFlags().BoolVarP(&utils.SilentFlag, "yes", "y", false, "(optional) assume yes for user interactive prompts.")
	SystemCmd.AddCommand(availableSubcommands...)
}
