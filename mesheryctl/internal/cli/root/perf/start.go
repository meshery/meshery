package perf

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	skipMeshConfigFlag bool
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Meshery's performance component",
	Long:  `Start Meshery's performance testing component (meshery-nighthawk).`,
	Args:  cobra.NoArgs,
	Example: `
// Start meshery-nighthawk
mesheryctl perf start

// Start meshery-nighthawk without adding entry to meshconfig
mesheryctl perf start --skip-meshconfig
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
		if err := startPerf(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to start meshery-nighthawk"))
		}
		return nil
	},
}

func updateMeshConfig() error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return err
	}

	components := currCtx.GetComponents()

	for _, comp := range components {
		if comp == "meshery-nighthawk" {
			return nil
		}
	}

	components = append(components, "meshery-nighthawk")

	newCtx := config.Context{
		Channel:    currCtx.GetChannel(),
		Components: components,
		Endpoint:   currCtx.GetEndpoint(),
		Platform:   currCtx.GetPlatform(),
		Token:      currCtx.GetToken(),
		Version:    currCtx.GetVersion(),
		Provider:   currCtx.GetProvider(),
	}

	if err := config.UpdateContextInConfig(&newCtx, mctlCfg.GetCurrentContextName()); err != nil {
		return errors.Wrap(err, "failed to update meshconfig")
	}

	return nil
}

func startPerf() error {
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
	if running {
		log.Info("meshery-nighthawk is already running")
		return nil
	}

	switch currCtx.GetPlatform() {
	case "docker":
		log.Info("Starting meshery-nighthawk...")

		remove := exec.Command("docker", "rm", "-f", "meshery-nighthawk")
		_ = remove.Run()

		start := exec.Command("docker", "run", "-d", "--name", "meshery-nighthawk",
			"--network", "meshery_default",
			"meshery/meshery-nighthawk:edge-latest")
		start.Stdout = os.Stdout
		start.Stderr = os.Stderr

		if err := start.Run(); err != nil {
			return errors.Wrap(err, "failed to start meshery-nighthawk container")
		}

		if !skipMeshConfigFlag {
			log.Info("Adding meshery-nighthawk entry to meshconfig...")
			if err := updateMeshConfig(); err != nil {
				log.Warn("Failed to update meshconfig: ", err)
			}
		}

		log.Info("meshery-nighthawk started successfully")

	case "kubernetes":
		kubeClient, err := meshkitkube.New([]byte(""))
		if err != nil {
			return err
		}

		log.Info("Starting meshery-nighthawk...")

		spinner := utils.CreateDefaultSpinner("Deploying meshery-nighthawk on Kubernetes", "\nmeshery-nighthawk deployed on Kubernetes.")
		spinner.Start()

		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "meshery-nighthawk",
				Namespace: utils.MesheryNamespace,
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "meshery-nighthawk",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "meshery-nighthawk",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "meshery-nighthawk",
								Image: "meshery/meshery-nighthawk:edge-latest",
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse("500m"),
										corev1.ResourceMemory: resource.MustParse("512Mi"),
									},
									Requests: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse("200m"),
										corev1.ResourceMemory: resource.MustParse("256Mi"),
									},
								},
							},
						},
					},
				},
			},
		}

		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "meshery-nighthawk",
				Namespace: utils.MesheryNamespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": "meshery-nighthawk",
				},
				Ports: []corev1.ServicePort{
					{
						Name:       "grpc",
						Port:       9999,
						TargetPort: intstr.FromInt(9999),
					},
				},
			},
		}

		_, err = kubeClient.KubeClient.AppsV1().Deployments(utils.MesheryNamespace).Create(context.Background(), deployment, metav1.CreateOptions{})
		if err != nil && !k8serror.IsAlreadyExists(err) {
			spinner.Stop()
			return errors.Wrap(err, "failed to create meshery-nighthawk deployment")
		}

		_, err = kubeClient.KubeClient.CoreV1().Services(utils.MesheryNamespace).Create(context.Background(), service, metav1.CreateOptions{})
		if err != nil && !k8serror.IsAlreadyExists(err) {
			spinner.Stop()
			return errors.Wrap(err, "failed to create meshery-nighthawk service")
		}

		spinner.Stop()
		log.Info("meshery-nighthawk started successfully")

		if !skipMeshConfigFlag {
			log.Info("Adding meshery-nighthawk entry to meshconfig...")
			if err := updateMeshConfig(); err != nil {
				log.Warn("Failed to update meshconfig: ", err)
			}
		}

	default:
		return fmt.Errorf("platform %s not supported", currCtx.GetPlatform())
	}

	return nil
}

func isPerfRunning() (bool, error) {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return false, err
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return false, err
	}

	switch currCtx.GetPlatform() {
	case "docker":
		cmd := exec.Command("docker", "ps", "--filter", "name=meshery-nighthawk", "--format", "{{.Names}}")
		output, err := cmd.Output()
		if err != nil {
			return false, err
		}
		return len(output) > 0, nil

	case "kubernetes":
		kubeClient, err := meshkitkube.New([]byte(""))
		if err != nil {
			return false, err
		}

		pods, err := kubeClient.KubeClient.CoreV1().Pods(utils.MesheryNamespace).List(context.Background(), metav1.ListOptions{
			LabelSelector: "app=meshery-nighthawk",
		})
		if err != nil {
			return false, err
		}

		for _, pod := range pods.Items {
			if pod.Status.Phase == "Running" {
				return true, nil
			}
		}

		return false, nil

	default:
		return false, fmt.Errorf("platform %s not supported", currCtx.GetPlatform())
	}
}

func init() {
	startCmd.Flags().BoolVar(&skipMeshConfigFlag, "skip-meshconfig", false, "(optional) skip adding entry to meshconfig")
}
