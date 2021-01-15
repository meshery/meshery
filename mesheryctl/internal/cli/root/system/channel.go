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
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var channelSet, channelSwitch, channelName, overrideContext string

func PrintChannelAndVersionToStdout(ctx config.Context) {
	log.Printf("Channel: %v", ctx.Channel)
	log.Printf("Version: %v", ctx.Version)
}

// channelCmd represents the config command
var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "Switch between release channels",
	Long:  `Subscribe to a release channel. Choose between either 'stable' or 'edge' channels.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		focusedContext := mctlCfg.CurrentContext

		if len(overrideContext) > 0 {
			focusedContext = overrideContext
		}

		if len(channelSet) == 0 && len(channelSwitch) == 0 {
			PrintChannelAndVersionToStdout(mctlCfg.Contexts[focusedContext])
			return nil
		}

		if channelSet == "" {
			channelName = channelSwitch
		} else {
			channelName = channelSet
		}

		channelNameSeperated := strings.Split(channelName, ":")

		ContextContent, ok := mctlCfg.Contexts[focusedContext]
		if !ok {
			return errors.New("error while trying to fetch context content")
		}

		switch channelNameSeperated[0] {
		case "stable":
			ContextContent.Channel = channelNameSeperated[0]
			if len(channelNameSeperated) == 1 {
				ContextContent.Version = "latest"
			} else {
				ContextContent.Version = channelNameSeperated[1]
			}

		case "edge":
			ContextContent.Channel = channelNameSeperated[0]
			ContextContent.Version = "latest"

		default:
			currentChannel := mctlCfg.Contexts[focusedContext].Channel
			if currentChannel == "" {
				return errors.New("No release channel subscription found." +
					"Please subscribe to either the 'stable' or 'edge' release channel")
			}
		}

		mctlCfg.Contexts[focusedContext] = ContextContent
		viper.Set("contexts", mctlCfg.Contexts)
		err = viper.WriteConfig()
		if err != nil {
			return err
		}
		PrintChannelAndVersionToStdout(mctlCfg.Contexts[focusedContext])

		if channelSwitch != "" {
			err = startCmd.RunE(cmd, nil)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	channelCmd.Flags().StringVarP(&channelSet, "set", "s", "", "Release channel to be set for Meshery and its adapters.")
	channelCmd.Flags().StringVarP(&channelSwitch, "switch", "w", "", "Release channel to be switch for Meshery and its adapters.")
	channelCmd.Flags().StringVarP(&overrideContext, "context", "c", "", "Override specified context with current context.")
}
