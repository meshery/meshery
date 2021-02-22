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
	"regexp"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mctlCfg *config.MesheryCtlConfig
var err error

var showForAllContext bool

// PrintChannelAndVersionToStdout to return curren release channel details
func PrintChannelAndVersionToStdout(ctx config.Context, contextName string) string {
	return fmt.Sprintf("Context: %v\nChannel: %v\nVersion: %v", contextName, ctx.Channel, ctx.Version)
}

// IsBetaOrStable to determine which release channel is being used
func IsBetaOrStable(str string) bool {
	return str == "edge" || str == "stable"
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view release channel and version",
	Long:  `View release channel and version of context in focus`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}
		focusedContext := tempContext
		if focusedContext == "" {
			focusedContext = mctlCfg.CurrentContext
		}

		if showForAllContext {
			for k, v := range mctlCfg.Contexts {
				log.Println(PrintChannelAndVersionToStdout(v, k))
				log.Println()
			}
			log.Printf("Current Context: %v", focusedContext)
			return nil
		}
		log.Print(PrintChannelAndVersionToStdout(mctlCfg.Contexts[focusedContext], focusedContext))
		log.Println()
		return nil
	},
}

var setCmd = &cobra.Command{
	Use:   "set [stable|stable-version|edge|edge-version]",
	Short: "set release channel and version",
	Long:  `Set release channel and version of context in focus`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		focusedContext := mctlCfg.CurrentContext

		if len(tempContext) > 0 {
			focusedContext = tempContext
		}

		channelVersion := args[0]

		channelNameSeperated := strings.SplitN(channelVersion, "-", 2)

		if !IsBetaOrStable(channelNameSeperated[0]) {
			return errors.New("No release channel subscription found." +
				"Please subscribe to either the 'stable' or 'edge' release channel")
		}

		version := "latest"

		if len(channelNameSeperated) > 1 {
			if channelNameSeperated[0] == "edge" {
				if channelNameSeperated[1] != "latest" {
					return errors.New("edge channel only supports latest as version argument")
				}
			} else if channelNameSeperated[0] == "stable" {
				if channelNameSeperated[1] != "latest" {
					matched, err := regexp.Match("v\\d+\\.\\d+.\\d+\\-(?:alpha|beta|rc)-\\d+", []byte(channelNameSeperated[1]))
					if err != nil || !matched && channelNameSeperated[1] != "latest" {
						return errors.New(fmt.Sprintf("%v is not a valid version tag", channelNameSeperated[1]))
					}
				}
			}
			version = channelNameSeperated[1]
		}

		ContextContent, ok := mctlCfg.Contexts[focusedContext]
		if !ok {
			return errors.New("error while trying to fetch context content")
		}

		ContextContent.Version = version
		ContextContent.Channel = channelNameSeperated[0]

		mctlCfg.Contexts[focusedContext] = ContextContent
		viper.Set("contexts", mctlCfg.Contexts)
		err = viper.WriteConfig()
		if err != nil {
			return err
		}

		return nil
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch [stable|stable-version|edge|edge-version]",
	Short: "switch release channel and version",
	Long:  `Switch release channel and version of context in focus`,
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if tempContext != "" {
			return utils.PreReqCheck(cmd.Use, tempContext)
		}
		return utils.PreReqCheck(cmd.Use, "")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err = setCmd.RunE(cmd, args)
		if err != nil {
			return err
		}
		err = startCmd.RunE(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

// channelCmd represents the config command
var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "Switch between release channels",
	Long:  `Subscribe to a release channel. Choose between either 'stable' or 'edge' channels.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}
		err = viewCmd.RunE(cmd, args)
		if err != nil {
			return err
		}
		err = cmd.Usage()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	viewCmd.Flags().BoolVarP(&showForAllContext, "all", "a", false, "Show release channel for all contexts")
	channelCmd.AddCommand(viewCmd, setCmd, switchCmd)
}
