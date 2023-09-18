// Copyright 2023 Layer5, Inc.
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
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mctlCfg *config.MesheryCtlConfig
var err error

var showForAllContext bool

// PrintChannelAndVersionToStdout to return current release channel details
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
	Example: `
// View current release channel
mesheryctl system channel view
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemChannelSubError("this command takes no arguments.\n", "view"))
		}
		mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
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

		err = mctlCfg.SetCurrentContext(focusedContext)
		if err != nil {
			utils.Log.Error(ErrSetCurrentContext(err))
			return nil
		}

		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(ErrGetCurrentContext(err))
			return nil
		}
		log.Print(PrintChannelAndVersionToStdout(*currCtx, focusedContext))
		log.Println()
		return nil
	},
}

var setCmd = &cobra.Command{
	Use:   "set [stable|stable-version|edge|edge-version]",
	Short: "set release channel and version",
	Long:  `Set release channel and version of context in focus`,
	Example: `
// Subscribe to release channel or version
mesheryctl system channel set [stable|stable-version|edge|edge-version]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = `Please provide either 'stable' or 'edge' release channel. Usage: mesheryctl system channel set [stable|stable-version|edge|edge-version]`
		if len(args) == 0 {
			return fmt.Errorf("release channel not specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments.\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		focusedContext := mctlCfg.CurrentContext

		if len(tempContext) > 0 {
			focusedContext = tempContext
		}

		channelVersion := args[0]

		channelNameSeperated := strings.SplitN(channelVersion, "-", 2)

		if !IsBetaOrStable(channelNameSeperated[0]) {
			return errors.New("No release channel subscription found. " +
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
					currCtx := mctlCfg.Contexts[focusedContext]
					currCtx.Version = channelNameSeperated[1]
				}
			}
			version = channelNameSeperated[1]
		}

		ContextContent, ok := mctlCfg.Contexts[focusedContext]
		if !ok {
			utils.Log.Error(ErrContextContent())
			return nil
		}

		ContextContent.Version = version
		ContextContent.Channel = channelNameSeperated[0]

		err = ContextContent.ValidateVersion()
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		mctlCfg.Contexts[focusedContext] = ContextContent
		viper.Set("contexts", mctlCfg.Contexts)
		err = viper.WriteConfig()
		if err != nil {
			utils.Log.Error(ErrWriteConfig(err))
			return nil
		}
		log.Infof("Channel set to %s-%s", ContextContent.Channel, ContextContent.Version)
		return nil
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch [stable|stable-version|edge|edge-version]",
	Short: "switch release channel and version",
	Long:  `Switch release channel and version of context in focus`,
	Example: `
// Switch between release channels
mesheryctl system channel switch [stable|stable-version|edge|edge-version]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = `Please provide either 'stable' or 'edge' release channel. Usage: mesheryctl system channel switch [stable|stable-version|edge|edge-version]`
		if len(args) == 0 {
			return fmt.Errorf("release channel not specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments.\n\n%v", errMsg)
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		hcOptions := &HealthCheckOptions{
			IsPreRunE:  true,
			PrintLogs:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hcOptions)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		return hc.RunPreflightHealthChecks()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		userResponse := false

		mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
		}
		focusedContext := tempContext
		if focusedContext == "" {
			focusedContext = mctlCfg.CurrentContext
		}

		//skip asking confirmation if -y flag used
		if utils.SilentFlag {
			userResponse = true
		} else {
			// ask user for confirmation
			userResponse = utils.AskForConfirmation("The Meshery deployment in context '" + focusedContext + "' will be replaced with a new Meshery deployment and channel subscription. Are you sure you want to continue")
		}

		if !userResponse {
			utils.Log.Error(ErrSwitchChannelResponse())
			return nil
		}

		err = setCmd.RunE(cmd, args)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		err = restartCmd.RunE(cmd, args)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		return nil
	},
}

// channelCmd represents the channel command
var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "Switch between release channels",
	Long:  `Subscribe to a release channel. Choose between either 'stable' or 'edge' channels.`,
	Example: `
// Subscribe to release channel or version
mesheryctl system channel

// To set the channel
mesheryctl system channel set [stable|stable-version|edge|edge-version]

// To pin/set the channel to a specific version
mesheryctl system channel set stable-v0.6.0

// To view release channel and version
mesheryctl system channel view

// To switch release channel and version
mesheryctl system channel switch [stable|stable-version|edge|edge-version]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		// If no subcommands are provided, show usage
		if len(args) == 0 {
			return cmd.Help()
		}

		// If an invalid subcommand is provided, return error
		if !utils.IsValidSubcommand(availableSubcommands, args[0]) {
			return errors.New(utils.SystemChannelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [set/switch/view]. Use 'mesheryctl system channel --help' to display usage guide.\n", args[0]), "channel"))
		}

		return nil
	},
}

func init() {
	viewCmd.Flags().BoolVarP(&showForAllContext, "all", "a", false, "Show release channel for all contexts")
	channelCmd.AddCommand(viewCmd, setCmd, switchCmd)
}
