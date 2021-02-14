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
	"os"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Meshery",
	Long:  `Stop all Meshery containers, remove their instances and prune their connected volumes.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if tempContext != "" {
			return utils.PreReqCheck(cmd.Use, tempContext)
		}
		return utils.PreReqCheck(cmd.Use, "")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("Stopping Meshery...")
		if err := stop(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to stop Meshery"))
		}
		return nil
	},
}

func stop() error {
	if !utils.IsMesheryRunning() {
		log.Info("Meshery is not running. Nothing to stop.")
		return nil
	}
	if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
		if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
			return errors.Wrapf(err, utils.SystemError(fmt.Sprintf("failed to mkdir %s", utils.MesheryFolder)))
		}
	}

	// Stop all Docker containers
	stop := exec.Command("docker-compose", "-f", utils.DockerComposeFile, "stop")
	stop.Stdout = os.Stdout
	stop.Stderr = os.Stderr

	if err := stop.Run(); err != nil {
		return errors.Wrap(err, utils.SystemError("failed to stop meshery - could not stop some containers."))
	}

	// Remove all Docker containers
	stop = exec.Command("docker-compose", "-f", utils.DockerComposeFile, "rm", "-f")
	stop.Stderr = os.Stderr

	if err := stop.Run(); err != nil {
		return errors.Wrap(err, utils.SystemError("failed to stop meshery"))
	}

	// Mesheryctl uses a docker volume for persistence. This volume should only be cleared when user wants
	// to start from scratch with a fresh install.
	// if err := exec.Command("docker", "volume", "prune", "-f").Run(); err != nil {
	// 	log.Fatal("[ERROR] Please install docker-compose. The error message: \n", err)
	// }

	log.Info("Meshery is stopped.")

	// Reset Meshery config file to default settings
	if utils.ResetFlag {
		err := resetMesheryConfig()
		if err != nil {
			return errors.Wrap(err, utils.SystemError("failed to reset meshery config"))
		}
	}
	return nil
}

func init() {
	stopCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
}
