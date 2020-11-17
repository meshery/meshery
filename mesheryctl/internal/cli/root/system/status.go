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
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Meshery status",
	Args:  cobra.NoArgs,
	Long:  `Check status of Meshery and Meshery adapters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("Meshery status... \n")

		start := exec.Command("docker-compose", "-f", utils.DockerComposeFile, "ps")

		outputStd, err := start.Output()
		if err != nil {
			return errors.Wrap(err, utils.SystemError("failed to get Meshery status"))
		}

		outputString := string(outputStd)
		if strings.Contains(outputString, "meshery") {
			log.Info(outputString)
		} else {
			log.Info("Meshery is not running, run `mesheryctl system start` to start Meshery")
		}

		return nil
	},
}
