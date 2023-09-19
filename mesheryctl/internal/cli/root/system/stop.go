// Copyright 2023 Layer5, Inc.
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
	"regexp"
	"strings"

	"github.com/docker/docker/client"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	controllerConfig "sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/layer5io/meshery-operator/api/v1alpha1"
)

var (
	// forceDelete used to clean-up meshery resources forcefully
	forceDelete bool
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Meshery",
	Long:  `Stop all Meshery containers / remove all Meshery resources.`,
	Example: `
// Stop Meshery
mesheryctl system stop

// Reset Meshery's configuration file to default settings.
mesheryctl system stop --reset

// (optional) keep the Meshery namespace during uninstallation
mesheryctl system stop --keep-namespace

// Stop Meshery forcefully (use it when system stop doesn't work)
mesheryctl system stop --force
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
			return errors.New(utils.SystemLifeCycleError(fmt.Sprintf("this command takes no arguments. See '%s --help' for more information.\n", cmd.CommandPath()), "stop"))
		}
		if err := stop(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to stop Meshery"))
		}
		return nil
	},
}

var userResponse bool

func stop() error {
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

	ok, err := utils.AreMesheryComponentsRunning(currCtx.GetPlatform())
	if err != nil {
		return err
	}
	// if --force passed possibly no deployments running but other stale resource present
	if !ok && !forceDelete {
		log.Info("Meshery resources are not running. Nothing to stop.")
		return nil
	}

	switch currCtx.GetPlatform() {
	case "docker":
		// if the platform is docker, then stop all the running containers
		if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
			if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
				return ErrCreateDir(err, utils.MesheryFolder)
			}
		}

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return errors.Wrap(err, "unable to create docker client")
		}

		spinner := utils.CreateDefaultSpinner("Removing meshery containers in docker\n", "")
		spinner.Start()

		options := types.ContainerListOptions{
			Filters: filters.NewArgs(
				filters.Arg("name", "meshery"),
				filters.Arg("network", "meshery_default"),
				filters.Arg("label", "com.centurylinklabs.watchtower.enable=true"),
			),
		}
		contaienrs, err := cli.ContainerList(context.Background(), options)

		headers := []string{"Container ID", "Image", "Names", "Status"}
		rows := [][]string{}

		for _, container := range contaienrs {
			// removes container forcefully
			err = cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{Force: true})
			if err != nil {
				log.Errorf("Error while removing %s: %s", container.Names[0], err.Error())
				continue
			}
			image := container.Image
			// containers those are deployed programatically have image sha256 key
			// in Container.Image field instead of image name
			if isSha256(container.Image) {
				image = getImageName(cli, container.Image)
			}
			names := strings.Join(container.Names, ",")
			rows = append(rows, []string{container.ID[0:12], image, names, "removed"})
		}

		spinner.FinalMSG = fmt.Sprintf("Removed %d containers from docker\n\n", len(rows))
		spinner.Stop()

		if len(rows) > 0 {
			utils.PrintToTable(headers, rows)
		}

	case "kubernetes":
		client, err := meshkitkube.New([]byte(""))
		if err != nil {
			return err
		}
		// if the platform is kubernetes, stop the deployment by uninstalling the helm charts
		userResponse = false
		if utils.SilentFlag {
			userResponse = true
		} else {
			// ask user for confirmation
			userResponse = utils.AskForConfirmation("Meshery deployments will be deleted from your cluster. Are you sure you want to continue")
		}

		if !userResponse {
			log.Info("Stop aborted.")
			return nil
		}

		log.Info("Stopping Meshery resources...")

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
			// DryRun helm release uninstallation with helm pkg
			// if err = client.ApplyHelmChart(meshkitkube.ApplyHelmChartConfig{
			// 	Namespace: utils.MesheryNamespace,
			// 	ChartLocation: meshkitkube.HelmChartLocation{
			// 		Repository: utils.HelmChartURL,
			// 		Chart:      utils.HelmChartName,
			// 	},
			// 	Action: meshkitkube.UNINSTALL,
			// 	DryRun: true,
			// }); err != nil {
			// 	// Dry run failed, in such case; fallback to force cleanup
			// 	if err = utils.ForceCleanupCluster(); err != nil {
			// 		return errors.Wrap(err, "cannot stop Meshery")
			// 	}
			// }

			// Dry run passed; now delete meshery components with the helm pkg
			err := applyHelmCharts(client, currCtx, currCtx.GetVersion(), false, meshkitkube.UNINSTALL)
			if err != nil {
				return errors.Wrap(err, "cannot stop Meshery")
			}
		}

		// Delete the CRDs for brokers and meshsyncs
		if err = invokeDeleteCRDs(); err != nil {
			return err
		}

		if !utils.KeepNamespace {
			log.Info("Deleting Meshery Namespace...")
			if err = deleteNs(utils.MesheryNamespace, client.KubeClient); err != nil {
				return err
			}
			// Wait for the namespace to be deleted
			deleted, err := utils.CheckMesheryNsDelete()
			if err != nil || !deleted {
				log.Info("Meshery is taking too long to stop.\nPlease check the status of the pods by executing “mesheryctl system status”.")
			} else {
				log.Info("Meshery resources are stopped.")
			}
		} else {
			log.Info("Meshery resources are stopped.")
		}
	}

	// Reset Meshery config file to default settings
	if utils.ResetFlag {
		err := resetMesheryConfig()
		if err != nil {
			return ErrResetMeshconfig(err)
		}
	}
	return nil
}

func isSha256(input string) bool {
	// hexadecimal regular expression pattern
	pattern := `^sha256:[0-9a-fA-F]{64}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(input)
}

// get image name by image SHA256 key
func getImageName(cli *client.Client, imageShaKey string) string {
	imageInspect, _, err := cli.ImageInspectWithRaw(context.Background(), imageShaKey)
	if err != nil {
		return "N/A"
	}
	return imageInspect.RepoTags[0]
}

// invokeDeleteCRs is a wrapper of deleteCR to delete CR instances (brokers and meshsyncs)
func invokeDeleteCRs(client *meshkitkube.Client) error {
	const (
		brokerResourceName   = "brokers"
		brokerInstanceName   = "meshery-broker"
		meshsyncResourceName = "meshsyncs"
		meshsyncInstanceName = "meshery-meshsync"
	)

	if err := deleteCR(brokerResourceName, brokerInstanceName, client); err != nil {
		err = ErrStopMeshery(errors.Wrap(err, "cannot delete CR "+brokerInstanceName))
		if !forceDelete {
			return err
		}

		log.Debug(err)
	}

	if err := deleteCR(meshsyncResourceName, meshsyncInstanceName, client); err != nil {
		err = ErrStopMeshery(errors.Wrap(err, "cannot delete CR "+meshsyncInstanceName))
		if !forceDelete {
			return err
		}

		log.Debug(err)
	}

	return nil
}

// deleteCRs delete the specified CR instance in the clusters
func deleteCR(resourceName, instanceName string, client *meshkitkube.Client) error {
	return client.DynamicKubeClient.Resource(schema.GroupVersionResource{
		Group:    v1alpha1.GroupVersion.Group,
		Version:  v1alpha1.GroupVersion.Version,
		Resource: resourceName,
	}).Namespace(utils.MesheryNamespace).Delete(context.TODO(), instanceName, metav1.DeleteOptions{})
}

// invokeDeleteCRs is a wrapper of deleteCRD to delete CRDs (brokers and meshsyncs)
func invokeDeleteCRDs() error {
	const (
		brokerCRDName   = "brokers.meshery.layer5.io"
		meshsyncCRDName = "meshsyncs.meshery.layer5.io"
	)

	cfg := controllerConfig.GetConfigOrDie()
	client, err := apiextension.NewForConfig(cfg)
	if err != nil {
		err = ErrStopMeshery(errors.Wrap(err, "cannot invoke delete CRDs"))
		if !forceDelete {
			return err
		}

		log.Debug(err)
	}

	if err = deleteCRD(brokerCRDName, client); err != nil {
		err = ErrStopMeshery(errors.Wrap(err, "cannot delete CRD "+brokerCRDName))
		if !forceDelete {
			return err
		}

		log.Debug(err)
	}

	if err = deleteCRD(meshsyncCRDName, client); err != nil {
		err = ErrStopMeshery(errors.Wrap(err, "cannot delete CRD "+meshsyncCRDName))
		if !forceDelete {
			return err
		}

		log.Debug(err)
	}

	return nil
}

// deleteCRs delete the specified CRD in the clusters
func deleteCRD(name string, client *apiextension.Clientset) error {
	return client.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func deleteNs(ns string, client *kubernetes.Clientset) error {
	return client.CoreV1().Namespaces().Delete(context.TODO(), ns, metav1.DeleteOptions{})
}

func init() {
	stopCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
	stopCmd.Flags().BoolVar(&utils.KeepNamespace, "keep-namespace", false, "(optional) keep the Meshery namespace during uninstallation")
	stopCmd.Flags().BoolVar(&forceDelete, "force", false, "(optional) uninstall Meshery resources forcefully")
}
