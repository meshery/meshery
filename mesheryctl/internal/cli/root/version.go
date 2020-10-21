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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/cfg"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-latest"
)

var (
	// Mesheryctl config - holds config handler
	mctlCfg    *cfg.MesheryCtl
	mesheryCfg *cfg.Version
)

// RequestErr is the error handler for the request
func requestErr(err error, url string) bool {
	if err != nil {
		logrus.Infof("Server Version: Unavailable \t  GitSHA: Unavailable")
		logrus.Errorf("\nError Occurred: %v", err)
		logrus.Errorf("\nCould not communicate with Meshery at %s", url+"/server/version")
		logrus.Errorf("Ensure that Meshery is available.\n See Meshery Documentation (https://docs.meshery.io) for help.\n")
		return true
	}
	return false
}

func checkLatestVersion(err error, serverVersion string) error {
	githubTag := &latest.GithubTag{
		Owner:      utils.GetMesheryGitHubOrg(),
		Repository: utils.GetMesheryGitHubRepo(),
	}
	if err != nil {
		return errors.Wrap(err, "could not reach Meshery GitHub repo")
	}
	// Compare current running Meshery server version to the latest available Meshery release on GitHub.
	res, err := latest.Check(githubTag, serverVersion)
	if err != nil {
		return errors.Wrap(err, "failed to compare latest and current version of Meshery")
	}
	// If user is running an outdated release, let them know.
	if res.Outdated {
		logrus.Info("\n", serverVersion, " is not the latest Meshery release. Update to v", res.Current, ". Run `mesheryctl system update`")
	}

	// If user is running the latest release, let them know.
	if res.Latest {
		logrus.Info("\n", serverVersion, " is the latest Meshery release.")
	}
	return nil
}

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

		logrus.Infof("Client Version: %v \t  GitSHA: %v", build, commitsha)

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/server/version", url), nil)
		if requestErr(err, url) {
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if requestErr(err, url) {
			return
		}

		data, err := ioutil.ReadAll(resp.Body)
		if requestErr(err, url) {
			return
		}

		version := cfg.Version{}
		err = json.Unmarshal(data, &version)
		if requestErr(err, url) {
			return
		}

		logrus.Infof("Server Version: %v \t  GitSHA: %v", version.GetBuild(), version.GetCommitSHA())

		// Inform user of the latest release version
		err = checkLatestVersion(err, version.GetBuild())
		if err != nil {
			logrus.Warn("\nfailed to check for latest version of Meshery")
		}
	},
}
