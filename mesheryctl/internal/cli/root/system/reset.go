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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset Meshery's configuration",
	Long:  `Reset Meshery to it's default configuration.`,
	Args:  cobra.NoArgs,
	Example: `
// Resets meshery.yaml file with a copy from Meshery repo
mesheryctl system reset

! Refer below image link for usage
* Usage of mesheryctl system reset
# ![reset-usage](/assets/img/mesheryctl/reset.png)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return resetMesheryConfig()
	},
}

// resets meshery config, skips conirmation if skipConfirmation is true
func resetMesheryConfig() error {
	userResponse := false
	if utils.SilentFlag {
		userResponse = true
	} else {
		// ask user for confirmation
		userResponse = utils.AskForConfirmation("Meshery config file will be reset to system defaults. Are you sure you want to continue")
	}
	if !userResponse {
		log.Info("Reset aborted.")
		return nil
	}

	// Get viper instance used for context
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return ErrProcessingMctlConfig(err)
	}
	// get the platform, channel and the version of the current context
	// if a temp context is set using the -c flag, use it as the current context
	if tempContext != "" {
		err = mctlCfg.SetCurrentContext(tempContext)
		if err != nil {
			return ErrSettingTemporaryContext(err)
		}
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return ErrRetrievingCurrentContext(err)
	}

	log.Info("Meshery resetting...\n")
	log.Printf("Current Context: %s", mctlCfg.GetCurrentContextName())
	log.Printf("Channel: %s", currCtx.GetChannel())
	log.Printf("Version: %s", currCtx.GetVersion())
	log.Printf("Platform: %s\n", currCtx.GetPlatform())

	// Reset the config file to the default context
	defaultContext := utils.TemplateContext
	defaultContext.Platform = currCtx.Platform
	err = config.AddContextToConfig(mctlCfg.GetCurrentContextName(), defaultContext, utils.DefaultConfigPath, true, true)
	if err != nil {
		return ErrSettingDefaultContextToConfig(err)
	}

	if err = fetchManifests(mctlCfg); err != nil {
		return err
	}

	return nil
}

// Fetches manifests for meshery components based on the current context
func fetchManifests(mctlCfg *config.MesheryCtlConfig) error {
	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return ErrRetrievingCurrentContext(err)
	}

	switch currCtx.GetPlatform() {
	case "docker":

		log.Printf("Fetching default docker-compose file as per current-context: %s...", mctlCfg.GetCurrentContextName())
		err = utils.DownloadDockerComposeFile(currCtx, true)
		if err != nil {
			return ErrDownloadFile(err, utils.DockerComposeFile)
		}

		err = utils.CreateManifestsFolder()

		if err != nil {
			return ErrCreateManifestsFolder(err)
		}

		log.Printf("...fetching Meshery Operator manifests for Kubernetes...")
		err = utils.DownloadOperatorManifest()

		if err != nil {
			return ErrDownloadFile(err, "operator manifest")
		}

		log.Info("...meshconfig (" + utils.DockerComposeFile + ") now reset to default settings.")

	case "kubernetes":

		log.Printf("Fetching Meshery Server and Meshery Operator manifests for  %s context...", mctlCfg.GetCurrentContextName())
		// fetch the manifest files corresponding to the version specified
		_, err := utils.FetchManifests(currCtx)

		if err != nil {
			return err
		}

		log.Info("...meshconfig has been reset to default settings.")

	default:
		return fmt.Errorf("the platform %s is not supported currently. The supported platforms are:\ndocker\nkubernetes\nPlease check %s/config.yaml file", currCtx.Platform, utils.MesheryFolder)
	}

	return nil
}
