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
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// BuildClient - holds the build info of Client
var BuildClient string

// BuildServer - holds the build info of server
var BuildServer string

// CommitSHA - holds the Git-SHA info
var CommitSHA string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of mesheryctl",
	Long:  `Version of Meshery command line client - mesheryctl.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Client Version: %v \t  GitSHA: %v", BuildClient, CommitSHA)
		log.Info("Server Version: ", BuildServer)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
