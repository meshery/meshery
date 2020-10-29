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
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/cfg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Mesheryctl config - holds config handler
	mctlCfg    *cfg.MesheryCtl
	mesheryCfg *cfg.Version
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of mesheryctl",
	Long:  `Version of Meshery command line client - mesheryctl.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		mctlCfg, err = cfg.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		url := mctlCfg.GetBaseMesheryURL()
		build := mctlCfg.GetVersion().GetBuild()
		commitsha := mctlCfg.GetVersion().GetCommitSHA()

		version := cfg.Version{
			Build:     "unavailable",
			CommitSHA: "unavailable",
		}

		logrus.Infof("Client Version: %v \t  GitSHA: %v", build, commitsha)

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/server/version", url), nil)
		if err != nil {
			logrus.Infof("Server Version: %v \t  GitSHA: %v", version.Build, version.CommitSHA)
			logrus.Errorf("\nUnable to get request context: %v", err)
			return
		}

		defer checkMesheryctlClientVersion(build)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logrus.Infof("Server Version: %v \t  GitSHA: %v", version.Build, version.CommitSHA)
			logrus.Errorf("\nUnable to communicate with Meshery: %v", err)
			logrus.Errorf("Ensure that Meshery is available, see Meshery Documentation (https://docs.meshery.io) for more help.\n")
			return
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Infof("Server Version: %v \t  GitSHA: %v", version.Build, version.CommitSHA)
			logrus.Errorf("\nInvalid response: %v", err)
			return
		}

		err = json.Unmarshal(data, &version)
		if err != nil {
			logrus.Infof("Server Version: %v \t  GitSHA: %v", version.Build, version.CommitSHA)
			logrus.Errorf("\nUnable to unmarshal data: %v", err)
			return
		}

		logrus.Infof("Server Version: %v \t  GitSHA: %v", version.GetBuild(), version.GetCommitSHA())
	},
}

func checkMesheryctlClientVersion(build string) {
	logrus.Infof("Checking for latest version of Meshery....")

	// Inform user of the latest release version
	_, err := handlers.CheckLatestVersion(build)
	if err != nil {
		logrus.Warn("\nfailed to check for latest version of Meshery")
	}
}
