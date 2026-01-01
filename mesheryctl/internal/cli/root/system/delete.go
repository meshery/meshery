// Copyright Meshery Authors
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
	"context"
	"fmt"
	"os"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Meshery containers",
	Long:  `Delete Meshery containers. This command removes all Meshery containers created by docker-compose.`,
	Example: `
// Delete Meshery containers
mesheryctl system delete

// Delete Meshery containers without confirmation
mesheryctl system delete -y
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check prerequisite
		hcOptions := &HealthCheckOptions{
			IsPreRunE:  true,
			PrintLogs:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hcOptions)
		if err != nil {
			return ErrHealthCheckFailed(err)
		}
		return hc.RunPreflightHealthChecks()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemLifeCycleError(fmt.Sprintf("this command takes no arguments. See '%s --help' for more information.\n", cmd.CommandPath()), "delete"))
		}
		if err := deleteContainers(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to delete Meshery containers"))
		}
		return nil
	},
}

func deleteContainers() error {
	// Get viper instance used for context
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	// if a temp context is set using the -c flag, use it as the current context
	err = mctlCfg.SetCurrentContext(tempContext)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve current-context")
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return err
	}

	switch currCtx.GetPlatform() {
	case "docker":
		// if the platform is docker, then remove all the containers
		if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
			if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
				return ErrCreateDir(err, utils.MesheryFolder)
			}
		}

		// Ask for confirmation
		userResponse := false
		if utils.SilentFlag {
			userResponse = true
		} else {
			userResponse = utils.AskForConfirmation("Meshery containers will be deleted. Are you sure you want to continue")
		}

		if !userResponse {
			utils.Log.Info("Delete aborted.")
			return nil
		}

		utils.Log.Info("Deleting Meshery containers...")

		// Use compose library to remove containers
		composeClient, err := utils.NewComposeClient()
		if err != nil {
			return errors.Wrap(err, utils.SystemError("failed to create compose client"))
		}

		// Remove all Docker containers (equivalent to docker compose rm)
		if err := composeClient.Remove(context.Background(), utils.DockerComposeFile); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to delete Meshery containers"))
		}
		utils.Log.Info("Meshery containers deleted.")
	case "kubernetes":
		utils.Log.Info("The 'delete' command is only supported for the Docker platform.")
		utils.Log.Info("For Kubernetes, use 'mesheryctl system stop' to remove Meshery resources.")
	default:
		return errors.New(utils.SystemError(fmt.Sprintf("unsupported platform: %s", currCtx.GetPlatform())))
	}

	return nil
}

func init() {
	// No additional flags for delete command
}
