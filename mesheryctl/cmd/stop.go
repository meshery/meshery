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
		logs.Info("Stopping Meshery...")
		if _, err := os.Stat(mesheryLocalFolder); os.IsNotExist(err) {
			os.Mkdir(mesheryLocalFolder, 0777)
		}

		if err := downloadFile(dockerComposeFile, fileURL); err != nil {
			log.Fatal(err)
		}
		if err := exec.Command("docker-compose", "-f", dockerComposeFile, "stop").Run(); err != nil {
			log.Fatal("[ERROR] Please install docker-compose. The error message: \n", err)
		}

		if err := exec.Command("docker-compose", "-f", dockerComposeFile, "rm", "-f").Run(); err != nil {
			log.Fatal("[ERROR] Please install docker-compose. The error message: \n", err)
		}

		if err := exec.Command("docker", "volume", "prune", "-f").Run(); err != nil {
			log.Fatal("[ERROR] Please install docker-compose. The error message: \n", err)
		}

		logs.Info("Meshery is stopped.")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
