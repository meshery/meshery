package perf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Meshery's performance component",
	Long:  `Stop Meshery's performance testing component (meshery-nighthawk).`,
	Args:  cobra.NoArgs,
	Example: `
// Stop meshery-nighthawk
mesheryctl perf stop
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		_, err = mctlCfg.GetCurrentContext()
		if err != nil {
			return errors.Wrap(err, "failed to get current context")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := stopPerf(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to stop meshery-nighthawk"))
		}
		return nil
	},
}

func stopPerf() error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return err
	}

	running, err := isPerfRunning()
	if err != nil {
		return err
	}
	if !running {
		log.Info("meshery-nighthawk is not running")
		return nil
	}

	switch currCtx.GetPlatform() {
	case "docker":
		log.Info("Stopping meshery-nighthawk...")

		stop := exec.Command("docker", "rm", "-f", "meshery-nighthawk")
		stop.Stdout = os.Stdout
		stop.Stderr = os.Stderr

		if err := stop.Run(); err != nil {
			return errors.Wrap(err, "failed to stop meshery-nighthawk container")
		}

		log.Info("meshery-nighthawk stopped successfully")

	case "kubernetes":
		kubeClient, err := meshkitkube.New([]byte(""))
		if err != nil {
			return err
		}

		log.Info("Stopping meshery-nighthawk...")

		spinner := utils.CreateDefaultSpinner("Removing meshery-nighthawk from Kubernetes", "\nmeshery-nighthawk removed from Kubernetes.")
		spinner.Start()

		deletePropagation := metav1.DeletePropagationForeground
		deleteOptions := metav1.DeleteOptions{
			PropagationPolicy: &deletePropagation,
		}

		err = kubeClient.KubeClient.AppsV1().Deployments(utils.MesheryNamespace).Delete(context.Background(), "meshery-nighthawk", deleteOptions)
		if err != nil && !k8serror.IsNotFound(err) {
			spinner.Stop()
			return errors.Wrap(err, "failed to delete meshery-nighthawk deployment")
		}

		err = kubeClient.KubeClient.CoreV1().Services(utils.MesheryNamespace).Delete(context.Background(), "meshery-nighthawk", deleteOptions)
		if err != nil && !k8serror.IsNotFound(err) {
			spinner.Stop()
			return errors.Wrap(err, "failed to delete meshery-nighthawk service")
		}

		gracePeriod := int64(0)
		forcePodDelete := metav1.DeleteOptions{
			GracePeriodSeconds: &gracePeriod,
		}
		err = kubeClient.KubeClient.CoreV1().Pods(utils.MesheryNamespace).DeleteCollection(context.Background(), forcePodDelete, metav1.ListOptions{
			LabelSelector: "app=meshery-nighthawk",
		})
		if err != nil && !k8serror.IsNotFound(err) {
			spinner.Stop()
			return errors.Wrap(err, "failed to delete meshery-nighthawk pods")
		}

		err = wait.PollImmediate(2*time.Second, 60*time.Second, func() (bool, error) {
			pods, err := kubeClient.KubeClient.CoreV1().Pods(utils.MesheryNamespace).List(context.Background(), metav1.ListOptions{
				LabelSelector: "app=meshery-nighthawk",
			})
			if err != nil {
				return false, err
			}
			return len(pods.Items) == 0, nil
		})
		if err != nil {
			pods, listErr := kubeClient.KubeClient.CoreV1().Pods(utils.MesheryNamespace).List(context.Background(), metav1.ListOptions{
				LabelSelector: "app=meshery-nighthawk",
			})
			if listErr != nil {
				spinner.Stop()
				return errors.Wrap(err, "failed to list meshery-nighthawk pods")
			}

			for _, pod := range pods.Items {
				_ = kubeClient.KubeClient.CoreV1().Pods(utils.MesheryNamespace).Delete(context.Background(), pod.Name, forcePodDelete)
			}

			spinner.Stop()
			return errors.Wrap(err, "some pods may still be terminating, please check with 'kubectl get pods -n meshery'")
		}

		spinner.Stop()
		log.Info("meshery-nighthawk stopped successfully")

	default:
		return fmt.Errorf("platform %s not supported", currCtx.GetPlatform())
	}

	return nil
}

func init() {
}
