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

package perf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands []*cobra.Command
	outputFormatFlag     string
	tokenPath            string
)

// PerfCmd represents the Performance Management CLI command
var PerfCmd = &cobra.Command{
	Use:     "perf",
	Short:   "Performance Management",
	Long:    `Performance Management & Benchmarking using Meshery CLI.`,
	Example: "mesheryctl perf --name \"a quick stress test\" --url http://192.168.1.15/productpage --qps 300 --concurrent-requests 2 --duration 30s --token \"provider=Meshery\"",
	Args:    cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		//Check prerequisite
		return system.RunPreflightHealthChecks(true, cmd.Use)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Importing SMP Configuration from the file
		var req *http.Request
		if filePath != "" {
			smpConfig, err := ioutil.ReadFile(filePath)

			if err != nil {
				return err
			}

			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return errors.Wrap(err, "error processing config")
			}

			req, err = http.NewRequest("POST", mctlCfg.GetBaseMesheryURL()+"/api/perf/load-test-smp", bytes.NewBuffer(smpConfig))
			if err != nil {
				return errors.Wrapf(err, utils.PerfError("Failed to invoke performance test"))
			}
		} else {
			var err error

			if testName == "" {
				log.Debug("Test Name not provided")
				testName = utils.StringWithCharset(8)
				log.Debug("Using random test name: ", testName)
			}

			if testURL == "" {
				return errors.New(utils.PerfError("please enter a test URL"))
			}

			// Method to check if the entered Test URL is valid or not
			var validURL = govalidator.IsURL(testURL)

			if !validURL {
				return errors.New(utils.PerfError("please enter a valid test URL"))
			}

			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return errors.Wrap(err, "error processing config")
			}

			req, err = http.NewRequest("POST", mctlCfg.GetBaseMesheryURL()+"/api/perf/load-test", nil)
			if err != nil {
				return errors.Wrapf(err, utils.PerfError("Failed to invoke performance test"))
			}

			q := req.URL.Query()

			q.Add("name", testName)
			q.Add("loadGenerator", loadGenerator)
			q.Add("c", concurrentRequests)
			q.Add("url", testURL)
			q.Add("qps", qps)

			durLen := len(testDuration)

			q.Add("dur", string(testDuration[durLen-1]))
			q.Add("t", string(testDuration[:durLen-1]))

			if testMesh != "" {
				q.Add("mesh", testMesh)
			}
			req.URL.RawQuery = q.Encode()

		}

		if err := utils.AddAuthDetails(req, tokenPath); err != nil {
			return errors.Wrap(err, utils.PerfError("failed to add auth details to request"))
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrapf(err, utils.PerfError(fmt.Sprintf("failed to make request to %s", testURL)))
		}
		if utils.ContentTypeIsHTML(resp) {
			return errors.New("failed to run test")
		}
		log.Debug("Initiating Performance test ...")
		log.Debug(resp.Status)

		defer utils.SafeClose(resp.Body)
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, utils.PerfError("failed to read response body"))
		}
		log.Debug(string(data))

		if err := utils.UpdateAuthDetails(tokenPath); err != nil {
			return errors.Wrap(err, utils.PerfError("failed to update auth details"))
		}

		log.Debug("Test Completed Successfully!")
		return nil
	},
}

func init() {
	PerfCmd.PersistentFlags().StringVarP(&tokenPath, "token", "t", "", "(required) Path to meshery auth config")
	PerfCmd.PersistentFlags().StringVarP(&outputFormatFlag, "output-format", "o", "", "(optional) format to display in [json|yaml]")

	availableSubcommands = []*cobra.Command{listCmd, viewCmd, applyCmd}
	PerfCmd.AddCommand(availableSubcommands...)
}
