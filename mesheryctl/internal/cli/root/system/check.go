package system

import (
	"context"
	"os/exec"
	"runtime"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sVersion "k8s.io/apimachinery/pkg/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PreCheckCmd = &cobra.Command{
	Use:   "preflight",
	Short: "Meshery pre-flight check",
	Long:  `Verify environment readiness to deploy Meshery.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		currCtx, err := mctlCfg.SetCurrentContext(tempContext)
		if err != nil {
			return err
		}
		return runPreflightHealthChecks(false, cmd.Use, currCtx.Platform)

	},
}

//Run preflight healthchecks to verify environment health
func runPreflightHealthChecks(isPreRunExecution bool, subcommand string, platform string) error {
	//Run docker healthchecks
	if err := runDockerHealthCheck(isPreRunExecution, subcommand, platform); err != nil {
		return err
	}
	//Run k8s API healthchecks
	if err = runKubernetesAPIHealthCheck(); err != nil {
		return err
	}
	//Run k8s plus kubectl minimum version healthchecks
	if err := runKubernetesVersionHealthCheck(); err != nil {
		return err
	}
	log.Info("\n--------------\n--------------\n✓✓ Meshery prerequisites met")
	return nil
}

//Run healthchecks to verify if docker is running and active
func runDockerHealthCheck(isPreRunExecution bool, subcommand string, platform string) error {
	if !isPreRunExecution {
		log.Info("\nDocker \n--------------")
	}
	//Check whether docker daemon is running or not
	if err := exec.Command("docker", "ps").Run(); err != nil {
		log.Warn("!! Docker is not running")
		//If preRunExecution and the current platform is docker then we trigger docker installation
		if isPreRunExecution && platform == "docker" {
			//No auto installation of docker for windows
			if runtime.GOOS == "windows" {
				return errors.Wrapf(err, "Please start Docker. Run `mesheryctl system %s` once Docker is started.", subcommand)
			}
			err = utils.Startdockerdaemon(subcommand)
			if err != nil {
				return errors.Wrapf(err, "failed to start Docker.")
			}
		}
	}
	if !isPreRunExecution {
		log.Info("√ Docker is running")
	}

	//Check for installed docker-compose on client system
	if err := exec.Command("docker-compose", "-v").Run(); err != nil {
		log.Warn("!! docker-compose is not available")
		if isPreRunExecution && platform == "docker" {
			//No auto installation of Docker-compose for windows
			if runtime.GOOS == "windows" {
				return errors.Wrapf(err, "please install docker-compose. Run `mesheryctl system %s` after docker-compose is installed.", subcommand)
			}
			err = utils.Installprereq()
			if err != nil {
				return errors.Wrapf(err, "failed to install prerequisites. Run `mesheryctl system %s` after docker-compose is installed.", subcommand)
			}
		}
	}
	if !isPreRunExecution {
		log.Info("√ docker-compose is available")
	}

	return nil
}

//Run healthchecks to verify if kubernetes client can be initialized and can be queried
func runKubernetesAPIHealthCheck() error {
	log.Info("\nKubernetes API \n--------------")
	//Check whether k8s client can be initialized
	client, err := meshkitkube.New([]byte(""))
	if err != nil {
		return errors.New("!! cannot initialize the client")
	}
	log.Info("√ can initialize the client")

	//Check whether kubernetes can be queried
	podInterface := client.KubeClient.CoreV1().Pods("")
	_, err = podInterface.List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return errors.New("!! cannot query the Kubernetes API")
	}
	log.Info("√ can query the Kubernetes API")

	return nil
}

//Run healthchecks to verify kubectl and kubenetes version with
// minimum compatible versions
func runKubernetesVersionHealthCheck() error {
	log.Info("\nKubernetes Version \n--------------")
	//Check whether system has minimum supported versions of kubernetes and kubectl
	var kubeVersion *k8sVersion.Info
	kubeVersion, err := utils.GetK8sVersionInfo()
	if err != nil {
		return err
	}
	if err = utils.CheckK8sVersion(kubeVersion); err != nil {
		return err
	}
	log.Info("√ is running the minimum Kubernetes version")

	if err = utils.CheckKubectlVersion(); err != nil {
		return err
	}
	log.Info("√ is running the minimum kubectl version")

	return nil
}
