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
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// cleanupCmd represents the cleanup command
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleanup Meshery's configuration",
	Long:  `Reset Meshery to it's default configuration.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		resetMesheryConfig()
	},
}

// resets meshery config
func resetMesheryConfig() {
	log.Info("Meshery resetting...")
	if err := utils.DownloadFile(utils.DockerComposeFile, fileURL); err != nil {
		log.Fatal("Error while resetting:", err)
	}
	log.Info("Meshery config (" + utils.DockerComposeFile + ") settings reset to defaults.")
}
