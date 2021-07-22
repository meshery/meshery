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
	"github.com/pkg/errors"

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
		return errors.Wrap(err, "error processing config")
	}
	// get the platform, channel and the version of the current context
	// if a temp context is set using the -c flag, use it as the current context
	if tempContext != "" {
		err = mctlCfg.SetCurrentContext(tempContext)
		if err != nil {
			return errors.Wrap(err, "failed to set temporary context")
		}
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve current-context")
	}

	log.Info("Meshery resetting...\n")
	log.Printf("Current Context: %s", mctlCfg.GetCurrentContextName())
	log.Printf("Channel: %s", currCtx.GetChannel())
	log.Printf("Version: %s", currCtx.GetVersion())
	log.Printf("Platform: %s\n", currCtx.GetPlatform())

	switch currCtx.GetPlatform() {
	case "docker":

		log.Printf("Fetching default docker-compose file as per current-context: %s...\n", mctlCfg.GetCurrentContextName())
		err = utils.DownloadDockerComposeFile(currCtx, true)
		if err != nil {
			return errors.Wrap(err, "failed to fetch docker-compose file")
		}

		err = utils.CreateManifestsFolder()

		if err != nil {
			return err
		}

		log.Printf("Fetching Meshery Operator manifests...\n")
		err = utils.DownloadOperatorManifest()

		if err != nil {
			return err
		}

		log.Info("...Meshery config (" + utils.DockerComposeFile + ") now reset to default settings.")

	case "kubernetes":

		version := currCtx.GetVersion()
		if version == "latest" {
			if currCtx.GetChannel() == "edge" {
				version = "master"
			} else {
				version, err = utils.GetLatestStableReleaseTag()
				if err != nil {
					return err
				}
			}
		}

		log.Printf("Fetching Meshery Server and Meshery Operator manifests...")
		// fetch the manifest files corresponding to the version specified
		_, err := utils.FetchManifests(version)

		if err != nil {
			return err
		}

		log.Info("...meshconfig has been reset to default settings.")

	default:
		return fmt.Errorf("the platform %s is not supported currently. The supported platforms are:\ndocker\nkubernetes\nPlease check %s/config.yaml file", currCtx.Platform, utils.MesheryFolder)
	}
	return nil
}
