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

	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Meshery",
	Long: `Uninstall Meshery and completely remove the Meshery installation.

Stop and remove all Meshery containers (Docker) or delete all Meshery resources
from the cluster, including custom resource instances, custom resource
definitions, the Helm release and the Meshery namespace (Kubernetes).
	Find more information at: https://docs.meshery.io/reference/references/mesheryctl/system/uninstall`,
	Example: `
// Uninstall Meshery
mesheryctl system uninstall

// (optional) keep the Meshery namespace during uninstallation
mesheryctl system uninstall --keep-namespace

// Uninstall Meshery forcefully (use it when system uninstall doesn't work)
mesheryctl system uninstall --force
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
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
			return errors.New(utils.SystemLifeCycleError(fmt.Sprintf("this command takes no arguments. See '%s --help' for more information.\n", cmd.CommandPath()), "uninstall"))
		}
		if err := uninstallMeshery(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to uninstall Meshery"))
		}
		return nil
	},
}

func uninstallMeshery() error {
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
	case platformDocker:
		// if the platform is docker, then stop and remove all the containers
		if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
			if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
				return ErrCreateDir(err, utils.MesheryFolder)
			}
		}

		// without the compose file there is no Docker installation to remove
		if _, err := os.Stat(utils.DockerComposeFile); os.IsNotExist(err) {
			utils.Log.Info("Meshery is not installed. Nothing to uninstall.")
			return nil
		}

		// Ask for confirmation
		userResponse := false
		if utils.SilentFlag {
			userResponse = true
		} else {
			userResponse = utils.AskForConfirmation("Meshery containers will be stopped and removed. Are you sure you want to continue")
		}

		if !userResponse {
			utils.Log.Info("Uninstall aborted.")
			return nil
		}

		utils.Log.Info("Uninstalling Meshery resources...")

		// Use compose library instead of exec.Command
		composeClient, err := utils.NewComposeClient()
		if err != nil {
			return utils.ErrDockerComposeClient(err)
		}

		// Stop all Docker containers
		if err := composeClient.Stop(context.Background(), utils.DockerComposeFile); err != nil {
			return utils.ErrDockerComposeStop(err)
		}

		// Remove all Docker containers
		if err := composeClient.Remove(context.Background(), utils.DockerComposeFile); err != nil {
			return ErrUninstallMeshery(err)
		}
		utils.Log.Info("Meshery uninstalled.")
	case platformKubernetes:
		client, err := meshkitkube.New([]byte(""))
		if err != nil {
			return err
		}

		// ask user for confirmation
		userResponse := false
		if utils.SilentFlag {
			userResponse = true
		} else {
			userResponse = utils.AskForConfirmation("Meshery deployments will be deleted from your cluster. Are you sure you want to continue")
		}

		if !userResponse {
			utils.Log.Info("Uninstall aborted.")
			return nil
		}

		utils.Log.Info("Uninstalling Meshery resources...")

		// Delete the CR instances for brokers and meshsyncs
		// this needs to be executed before deleting the helm release, or the CR instances cannot be found for some reason
		if err = invokeDeleteCRs(client); err != nil {
			return err
		}

		if forceDelete {
			if err = utils.ForceCleanupCluster(); err != nil {
				return err
			}
		} else {
			// Delete meshery components by uninstalling the helm release
			err := applyHelmCharts(client, currCtx, currCtx.GetVersion(), false, meshkitkube.UNINSTALL, "", "")
			if err != nil {
				return errors.Wrap(err, "cannot uninstall Meshery")
			}
		}

		// Delete the CRDs for brokers and meshsyncs
		if err = invokeDeleteCRDs(); err != nil {
			return err
		}

		if !utils.KeepNamespace {
			utils.Log.Info("Deleting Meshery Namespace...")
			if err = deleteNs(utils.MesheryNamespace, client.KubeClient); err != nil {
				return err
			}
			// Wait for the namespace to be deleted
			deleted, err := utils.CheckMesheryNsDelete()
			if err != nil || !deleted {
				utils.Log.Info("Meshery is taking too long to uninstall.\nPlease check the status of the resources by executing “mesheryctl system status”.")
			} else {
				utils.Log.Info("Meshery uninstalled.")
			}
		} else {
			utils.Log.Info("Meshery uninstalled. The Meshery namespace was kept.")
		}
	default:
		return ErrUnsupportedPlatform(currCtx.GetPlatform(), utils.CfgFile)
	}

	return nil
}

func init() {
	uninstallCmd.Flags().BoolVar(&forceDelete, "force", false, "(optional) uninstall Meshery resources forcefully")
	uninstallCmd.Flags().BoolVar(&utils.KeepNamespace, "keep-namespace", false, "(optional) keep the Meshery namespace during uninstallation")
}
