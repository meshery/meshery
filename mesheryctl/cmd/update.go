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

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull new Meshery images from Docker Hub.",
	Long:  `Poll Docker Hub for new Meshery container images and pulls if new image version(s) are available.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Updating Meshery now...")

		if _, err := os.Stat(dockerComposeFile); os.IsNotExist(err) {
			if err := downloadFile(dockerComposeFile, fileURL); err != nil {
				log.Fatal("update cmd: ", err)
			}
		}
		if err := exec.Command("docker-compose", "-f", dockerComposeFile, "pull").Run(); err != nil {
			log.Fatal("[ERROR] Please, install docker-compose. The error message: \n", err)
		}

		log.Info("Meshery is now up-to-date")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
