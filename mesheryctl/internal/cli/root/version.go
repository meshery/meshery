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

package root

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Mesheryctl config - holds config handler
	mctlCfg *config.MesheryCtlConfig
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of mesheryctl",
	Long:  `Version of Meshery command line client - mesheryctl.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			// get the currCtx
			logrus.Errorf("error processing config: %v", err)
			userResponse := false
			userResponse = utils.AskForConfirmation("Looks like you are using an outdated config file. Do you want to generate a new config file?")
			if userResponse {
				utils.BackupConfigFile(utils.DefaultConfigPath)
				// Create config file if not present in meshery folder
				err = utils.CreateConfigFile()
				if err != nil {
					logrus.Errorf("unable to create config file")
				}

				// Add Token to context file
				err = config.AddTokenToConfig(utils.TemplateToken, utils.DefaultConfigPath)
				if err != nil {
					logrus.Errorf("unable to add token to config")
				}

				// Add Context to context file
				err = config.AddContextToConfig("local", utils.TemplateContext, utils.DefaultConfigPath, true)
				if err != nil {
					logrus.Errorf("unable to add context to config")
				}

				logrus.Printf(
					fmt.Sprintf("Default config file created at %s",
						utils.DefaultConfigPath,
					))

				mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
				if err != nil {
					logrus.Errorf("error unmarshaling config file")
				}
				currCtx, err := mctlCfg.GetCurrentContext()
				if err != nil {
					return err
				}
				err = currCtx.ValidateVersion()
				if err != nil {
					return err
				}
				return nil
			}
			return handlers.ErrUnmarshal(errors.New("invalid config file"), "meshconfig")
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = currCtx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		url := mctlCfg.GetBaseMesheryURL()
		build := constants.GetMesheryctlVersion()
		commitsha := constants.GetMesheryctlCommitsha()

		version := config.Version{
			Build:          "unavailable",
			CommitSHA:      "unavailable",
			ReleaseChannel: "unavailable",
		}

		header := []string{"", "Version", "GitSHA"}
		rows := [][]string{{"Client", build, commitsha}, {"Server", version.GetBuild(), version.GetCommitSHA()}}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/system/version", url), nil)
		if err != nil {
			utils.PrintToTable(header, rows)
			logrus.Errorf("\nUnable to get request context: %v", err)
			return
		}

		defer checkMesheryctlClientVersion(build)
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			utils.PrintToTable(header, rows)
			logrus.Errorf("\n  Unable to communicate with Meshery: %v", err)
			logrus.Errorf("  See https://docs.meshery.io for help getting started with Meshery.")
			return
		}

		// needs multiple defer as Body.Close needs a valid response
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.PrintToTable(header, rows)
			logrus.Errorf("\n  Invalid response: %v", err)
			return
		}

		err = json.Unmarshal(data, &version)
		if err != nil {
			utils.PrintToTable(header, rows)
			logrus.Errorf("\n  Unable to unmarshal data: %v", err)
			return
		}
		rows[1][1] = version.GetBuild()
		rows[1][2] = version.GetCommitSHA()
		utils.PrintToTable(header, rows)
	},
}

func checkMesheryctlClientVersion(build string) {
	logrus.Infof("\nChecking for latest version of mesheryctl...")

	// Inform user of the latest release version
	res, err := utils.GetLatestStableReleaseTag()
	if err != nil {
		logrus.Warn("\n  Unable to check for latest version of mesheryctl. ", err)
		return
	}
	// If user is running an outdated release, let them know.
	if res != build {
		logrus.Info("\n  ", build, " is not the latest release. Update to ", res, ".")
	} else { // If user is running the latest release, let them know.
		logrus.Info("\n  ", res, " is the latest release.")
	}
}
