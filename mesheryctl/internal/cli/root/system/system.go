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

package system

import (
	"fmt"

	config "github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system/context"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// TODO: `context` support for `system start` change the docker URL to dynamic URL
	fileURL         = "https://raw.githubusercontent.com/layer5io/meshery/master/docker-compose.yaml"
	manifestsURL    = "https://api.github.com/repos/layer5io/meshery/git/trees/3ba314a0870d5291be6216b4d60d2bc9675a39b2"
	rawManifestsURL = "https://raw.githubusercontent.com/layer5io/meshery/master/install/deployment_yamls/k8s/"
	gitHubFolder    = "https://github.com/layer5io/meshery/tree/master/install/deployment_yamls/k8s"
)

var (
	availableSubcommands = []*cobra.Command{}
	url                  = ""
	// flag to change the current context to a temporary context
	tempContext = ""
)

// SystemCmd represents Meshery Lifecycle Management cli commands
var SystemCmd = &cobra.Command{
	Use:   "system",
	Short: "Meshery Lifecycle Management",
	Long:  `Manage the state and configuration of Meshery server, adapters, and client.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		url = mctlCfg.GetBaseMesheryURL()
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
		context.ContextCmd,
		completionCmd,
		channelCmd,
	}
	// --context flag to temporarily change context. This is global to all system commands
	SystemCmd.PersistentFlags().StringVarP(&tempContext, "context", "c", "", "(optional) temporarily change the current context.")
	SystemCmd.AddCommand(availableSubcommands...)
}
