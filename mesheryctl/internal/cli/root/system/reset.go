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

// resets meshery config
func resetMesheryConfig() error {
	log.Info("Meshery config file will be reset to system defaults. Are you sure you want to continue [y/n]?")
	res := ""
	fmt.Scanf("%s", &res)

	if res == "n" || res == "N" {
		log.Info("Aborting reset...")
	} else if res == "y" || res == "Y" {

		log.Info("Meshery resetting...\n")

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		currentContext := mctlCfg.CurrentContext
		currChannel := mctlCfg.Contexts[currentContext].Channel
		currVersion := mctlCfg.Contexts[currentContext].Version

		fileURL := ""

		if currChannel == "edge" {
			fileURL = "https://raw.githubusercontent.com/layer5io/meshery/master/docker-compose.yaml"
		} else if currChannel == "stable" {
			fileURL = "https://raw.githubusercontent.com/layer5io/meshery/" + currVersion + "/docker-compose.yaml"
		}

		log.Printf("Current Context: %s", currentContext)
		log.Printf("Channel: %s", currChannel)
		log.Printf("Version: %s\n", currVersion)
		log.Printf("Fetching default docker-compose file at version: %s...\n", currVersion)

		if err := utils.DownloadFile(utils.DockerComposeFile, fileURL); err != nil {
			return errors.Wrapf(err, utils.SystemError(fmt.Sprintf("failed to download %s file from %s", utils.DockerComposeFile, fileURL)))
		}
		log.Info("...Meshery config (" + utils.DockerComposeFile + ") now reset to default settings.")
	} else {
		log.Error("Invalid response! Please try again.")
	}
	return nil
}
