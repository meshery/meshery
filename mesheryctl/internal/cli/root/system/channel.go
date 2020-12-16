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
	"errors"
	"io/ioutil"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var channelSet, channelSwitch, channelName string

// configCmd represents the config command
var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "Switch between release channels",
	Long:  `Subscribe to a release channel. Choose between either 'stable' or 'edge' channels.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if len(channelSet) < 0 && len(channelSwitch) < 0 {
			log.Fatal("Provide 'stable' or 'edge' as the release channel name.")
		}

		if channelSet == "" {
			channelName = channelSwitch
		} else {
			channelName = channelSet
		}

		switch channelName {
		case "stable":
			if err := searchAndReplace(utils.DockerComposeFile, "edge", "stable"); err != nil {
				log.Fatal("Error subscribing to release channel:", err)
				return
			} else {
				log.Info("Subscribed to stable channel.")
			}
		case "edge":
			if err := searchAndReplace(utils.DockerComposeFile, "stable", "edge"); err != nil {
				log.Fatal("Error subscribing to release channel:", err)
				return
			} else {
				log.Info("Subscribed to edge channel.")
			}
		default:
			log.Fatal("Please subscribe to either the 'stable' or 'edge' release channel.")
		}

		if channelSwitch != "" {
			if utils.IsMesheryRunning() {
				if err := stop(); err != nil {
					log.Fatal("Failed to stop Meshery:", err)
					return
				}
			}

			if err := start(); err != nil {
				log.Fatal("Failed to start Meshery:", err)
				return
			}
			log.Info("Successfully switched channel...")
		}

	},
}

func init() {
	channelCmd.Flags().StringVarP(&channelSet, "set", "s", "", "Release channel to be set for Meshery and its adapters.")
	channelCmd.Flags().StringVarP(&channelSwitch, "switch", "c", "", "Release channel to be switch for Meshery and its adapters.")
}

func searchAndReplace(mesheryConfig, currentChannel, newChannel string) error {
	fileData, err := ioutil.ReadFile(mesheryConfig)
	if err != nil {
		log.Errorf("unable to read meshery config file: %v", err)
		return err
	}

	if !strings.Contains(string(fileData), currentChannel) {
		log.Errorf("unable to switch channel with %s", newChannel)
		return errors.New("unable to set channel")
	}

	outputFile := strings.ReplaceAll(string(fileData), currentChannel, newChannel)

	err = ioutil.WriteFile(mesheryConfig, []byte(outputFile), 0644)
	if err != nil {
		log.Errorf("unable to change meshery config: %v", err)
		return errors.New("unable to write meshery config")
	}

	return nil
}
