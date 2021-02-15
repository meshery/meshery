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
	"os"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull new Meshery images from Docker Hub.",
	Long:  `Pull Docker Hub for new Meshery container images and pulls if new image version(s) are available.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		if tempContext != "" {
			return utils.PreReqCheck(cmd.Use, tempContext)
		}
		return utils.PreReqCheck(cmd.Use, "")
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if _, err := os.Stat(utils.DockerComposeFile); os.IsNotExist(err) {
			if err := utils.DownloadFile(utils.DockerComposeFile, fileURL); err != nil {
				return errors.Wrapf(err, utils.SystemError(fmt.Sprintf("failed to download %s file from %s", utils.DockerComposeFile, fileURL)))
			}
		}
		err := resetMesheryConfig()
		if err != nil {
			return errors.Wrap(err, utils.SystemError("failed to update meshery containers"))
		}

		err = utils.UpdateMesheryContainers()
		if err != nil {
			return errors.Wrap(err, utils.SystemError("failed to update meshery containers"))
		}

		log.Info("Meshery is now up-to-date")
		return nil
	},
}
