// Copyright 2019 The Meshery Authors
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

package cmd

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Meshery",
	Long:  `Stop all Meshery containers, remove their instances and prune their connected volumes.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Stopping Meshery...")
		if !isMesheryRunning() {
			log.Info("Meshery is not running. Nothing to stop.")
			return 
		}
		if _, err := os.Stat(mesheryFolder); os.IsNotExist(err) {
			if err := os.Mkdir(mesheryFolder, 0777); err != nil {
				log.Fatal(err)
			}
		}

		if err := downloadFile(dockerComposeFile, fileURL); err != nil {
			log.Fatal(err)
		}

		// Stop all Docker containers
		stop := exec.Command("docker-compose", "-f", dockerComposeFile, "stop")
		stop.Stdout = os.Stdout
		stop.Stderr = os.Stderr

		if err := stop.Run(); err != nil {
			log.Fatal("[ERROR] Could not completely stop all containers. The error message: \n", err)
		}

		// Remove all Docker containers
		stop = exec.Command("docker-compose", "-f", dockerComposeFile, "rm", "-f")
		stop.Stderr = os.Stderr

		if err := stop.Run(); err != nil {
			log.Fatal("[ERROR] Could not completely remove all containers. The error message: \n", err)
		}

		// Mesheryctl uses a docker volume for persistence. This volume should only be cleared when user wants
		// to start from scratch with a fresh install.
		// if err := exec.Command("docker", "volume", "prune", "-f").Run(); err != nil {
		// 	log.Fatal("[ERROR] Please install docker-compose. The error message: \n", err)
		// }

		log.Info("Meshery is stopped.")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
