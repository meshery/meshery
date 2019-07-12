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
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull new Meshery images from Docker Hub.",
	Long:  `Poll Docker Hub for new Meshery container images and pulls if new image version(s) are available.`,
	Run: func(cmd *cobra.Command, args []string) {
		Init()
		log.Println(mesheryPreUpdateMessage)
		if err := exec.Command("docker-compose", "-f", dockerComposeFile, "pull").Run(); err != nil {
			log.Fatal(dockerComposeWarningMessage, err)
		}
		log.Println(mesheryPostUpdateMessage)

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
