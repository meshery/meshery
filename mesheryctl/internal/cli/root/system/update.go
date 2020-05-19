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

package system

import (
	"os"
	"os/exec"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull new Meshery images from Docker Hub.",
	Long:  `Pull Docker Hub for new Meshery container images and pulls if new image version(s) are available.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if _, err := os.Stat(utils.DockerComposeFile); os.IsNotExist(err) {
			if err := utils.DownloadFile(utils.DockerComposeFile, fileURL); err != nil {
				log.Fatal("update cmd: ", err)
			}
		}

		updateMesheryContainers()

		log.Info("Meshery is now up-to-date")
	},
}

func updateMesheryContainers() {

	log.Info("Updating Meshery now...")
	start := exec.Command("docker-compose", "-f", utils.DockerComposeFile, "pull")
	start.Stdout = os.Stdout
	start.Stderr = os.Stderr
	if err := start.Run(); err != nil {
		log.Fatal(err)
	}

}
