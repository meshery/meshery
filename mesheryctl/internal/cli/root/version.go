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

package root

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/cfg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (

	// BuildClient - holds the build info of Client
	Build string

	// CommitSHA - holds the Git-SHA info
	CommitSHA string

	// Mesheryctl config - holds config handler
	mctlCfg *cfg.MesheryCtl
)

type Version struct {
	Build     string `json:"build,omitempty"`
	CommitSHA string `json:"commitsha,omitempty"`
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
		logrus.Infof("Client Version: %v \t  GitSHA: %v", Build, CommitSHA)
		version := &Version{
			Build:     "Unavailable",
			CommitSHA: "Unavailable",
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/server/version", mctlCfg.GetBaseMesheryURL()), nil)
		if err != nil {
			logrus.Infof("Server Version: %v \t  GitSHA: %v", version.Build, version.CommitSHA)
			logrus.Errorf("\nCould not communicate with Meshery at $BASE_URL\n")
			logrus.Errorf("Ensure that Meshery is available. See (https://docs.meshery.io)\n\n")
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logrus.Infof("Server Version: %v \t  GitSHA: %v", version.Build, version.CommitSHA)
			logrus.Errorf("\nCould not communicate with Meshery at $BASE_URL\n")
			logrus.Errorf("Ensure that Meshery is available. See (https://docs.meshery.io)\n\n")
			return
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Infof("Server Version: %v \t  GitSHA: %v", version.Build, version.CommitSHA)
			logrus.Errorf("\nCould not communicate with Meshery at $BASE_URL\n")
			logrus.Errorf("Ensure that Meshery is available. See (https://docs.meshery.io)\n\n")
			return
		}

		err = json.Unmarshal(data, version)
		if err != nil {
			logrus.Infof("Server Version: %v \t  GitSHA: %v", version.Build, version.CommitSHA)
			logrus.Errorf("\nCould not communicate with Meshery at $BASE_URL\n")
			logrus.Errorf("Ensure that Meshery is available. See (https://docs.meshery.io)\n\n")
			return
		}
		logrus.Infof("Server Version: %v \t  GitSHA: %v", version.Build, version.CommitSHA)
	},
}
