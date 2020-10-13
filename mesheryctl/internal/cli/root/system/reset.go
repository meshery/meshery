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

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset Meshery's configuration",
	Long:  `Reset Meshery to it's default configuration.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return resetMesheryConfig()
	},
}

// resets meshery config
func resetMesheryConfig() error {
	log.Info("Meshery resetting...")
	if err := utils.DownloadFile(utils.DockerComposeFile, fileURL); err != nil {
		return errors.Wrapf(err, utils.SystemError(fmt.Sprintf("failed to download %s file from %s", utils.DockerComposeFile, fileURL)))
	}
	log.Info("Meshery config (" + utils.DockerComposeFile + ") reset to default settings.")
	return nil
}
