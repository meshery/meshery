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
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var channelName string

// configCmd represents the config command
var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "Switch the release channel",
	Long:  `Switch the available release channels 'Stable' & 'Edge'`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if len(channelName) < 0 {
			log.Fatal("Give a valid channel name")

		}
		if channelName == "" {
			log.Fatal("Invalid channel name")
		}

		switch channelName {
		case "stable":
			if err := useStableChannel(); err != nil {
				log.Fatal("Error generating config:", err)
				return
			}
		case "edge":
			if err := useEdgeChannel(); err != nil {
				log.Fatal("Error generating config:", err)
				return
			}
		default:
			log.Fatal("Channel name has to be Stable | Edge.")
		}

		if utils.IsMesheryRunning() {
			if err := stop(); err != nil {
				log.Fatal("Failed to stop meshery:", err)
				return
			}
		}

		if err := start(); err != nil {
			log.Fatal("Failed to start meshery:", err)
			return
		}
	},
}

func init() {
	configCmd.Flags().StringVarP(&channelName, "channelName", "c", "stable", "Release channel to be used for Meshery and its adapters.")
	_ = configCmd.MarkFlagRequired("channelName")
}

func useEdgeChannel() error {
	configPath := "$HOME" + "/" + utils.MesheryFolder + "/" + utils.DockerComposeFile

	script := fmt.Sprintf(`
	set -e
	set -o pipefail

	File_Location="%s"

	replace_to_edge() {
		echo -n "replacing stable with edge"
		sed -i -e 's/stable/edge/g' "${File_Location}"
		printf "done\n"
	}

	replace_to_edge

    `, configPath)

	switchToEdge := exec.Command("bash", "-c", script)
	switchToEdge.Stdout = os.Stdout
	switchToEdge.Stderr = os.Stderr

	return switchToEdge.Run()

}

func useStableChannel() error {
	configPath := "$HOME" + "/" + utils.MesheryFolder + "/" + utils.DockerComposeFile

	script := fmt.Sprintf(`
	set -e
	set -o pipefail

	File_Location="%s"

	replace_to_edge() {
		echo -n "replacing stable with edge"
		sed -i -e 's/edge/stable/g' "${File_Location}"
		printf "done\n"
	}

	replace_to_edge

    `, configPath)

	switchToStable := exec.Command("bash", "-c", script)
	switchToStable.Stdout = os.Stdout
	switchToStable.Stderr = os.Stderr

	return switchToStable.Run()

}
