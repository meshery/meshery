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
	"errors"

	"github.com/spf13/cobra"
)

const (
	url     = "http://localhost:9081"
	fileURL = "https://raw.githubusercontent.com/layer5io/meshery/master/docker-compose.yaml"
)

var systemDetails = `
System level functionality to configure and interact with the Meshery client and server.

Usage:
  mesheryctl system [command]

Available Commands:
  cleanup     Clean up Meshery
  help        Help about any command
  logs        Print logs
  start       Start Meshery
  status      Check Meshery status
  stop        Stop Meshery
  update      Pull new Meshery images from Docker Hub

Flags:
  -h, --help            help for mesheryctl system

Use "mesheryctl system [command] --help" for more information about a command.
`

var (
	availableSubcommands = []*cobra.Command{}
)

// SystemCmd represents the system level commands
var SystemCmd = &cobra.Command{
	Use:   "system",
	Short: "Meshery Lifecycle Management",
	Long:  `Manage the state and configuration of Meshery server, adapters, and client.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, subcommand := range availableSubcommands {
			if args[0] == subcommand.Name() {
				return nil
			}
		}

		return errors.New("sub-command not found : " + "\"" + args[0] + "\"")
	},
}

func init() {
	availableSubcommands = []*cobra.Command{
		cleanupCmd,
		logsCmd,
		startCmd,
		stopCmd,
		statusCmd,
		updateCmd,
	}
	SystemCmd.AddCommand(availableSubcommands...)
}
