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
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset Meshery's configuration",
	Long:  `Reset Meshery to it's default configuration.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return resetMesheryConfig(false)
	},
}

// resets meshery config, skips conirmation if skipConfirmation is true
func resetMesheryConfig(skipConfirmation bool) error {
	// ask user for confirmation
	if !skipConfirmation {
		userResponse := utils.AskForConfirmation("Meshery config file will be reset to system defaults. Are you sure you want to continue")
		if !userResponse {
			log.Info("Reset aborted.")
			return nil
		}
	}

	currCtxName, currCtx, err := config.GetCurrentContext(tempContext)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve current-context")
	}

	log.Info("Meshery resetting...\n")
	log.Printf("Fetching default docker-compose file as per current-context: %s...\n", currCtxName)
	err = utils.DownloadDockerComposeFile(currCtx, true)
	if err != nil {
		return errors.Wrap(err, "failed to fetch docker-compose file")
	}

	log.Printf("Current Context: %s", currCtxName)
	log.Printf("Channel: %s", currCtx.Channel)
	log.Printf("Version: %s\n", currCtx.Version)

	log.Info("...Meshery config (" + utils.DockerComposeFile + ") now reset to default settings.")
	return nil
}
