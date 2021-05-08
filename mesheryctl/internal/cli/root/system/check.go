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

var (
	preflight bool
	pre       bool
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Meshery enviroment check",
	Long:  `Verify environment pre/post-deployment of Meshery.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		//incase --preflight or --pre has been passed as a flag
		if preflight || pre {
			return RunPreflightHealthChecks(false, cmd.Use)
		}
		return nil
	},
}

//Run preflight healthchecks to verify environment health
func RunPreflightHealthChecks(isPreRunExecution bool, subcommand string) error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}
	currCtx, err := mctlCfg.SetCurrentContext(tempContext)
	if err != nil {
		return err
	}
	platform := currCtx.Platform
	//Docker healthchecks are only invoked when it's not a PreRunExecution
	// or it's a PreRunExecution and current platform is docker
	if !isPreRunExecution || (isPreRunExecution && platform == "docker") {
		//Run docker healthchecks
		if err := runDockerHealthCheck(isPreRunExecution, subcommand, platform); err != nil {
			return err
		}
	}
	//Kubernetes and KubernetesVersion healthchecks are only invoked when it's not a PreRunExecution
	// or it's a PreRunExecution and current platform is kubernetes
	if !isPreRunExecution || (isPreRunExecution && platform == "kubernetes") {
		//Run k8s API healthchecks
		if err = runKubernetesAPIHealthCheck(isPreRunExecution); err != nil {
			return err
		}
		//Run k8s plus kubectl minimum version healthchecks
		if err := runKubernetesVersionHealthCheck(isPreRunExecution); err != nil {
			return err
		}
	}
	if !isPreRunExecution {
		log.Info("\n--------------\n--------------\n✓✓ Meshery prerequisites met")
	}
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
		if isPreRunExecution {
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
		if isPreRunExecution {
			//No auto installation of Docker-compose for windows
			if runtime.GOOS == "windows" {
				return errors.Wrapf(err, "please install docker-compose. Run `mesheryctl system %s` after docker-compose is installed.", subcommand)
			}
			err = utils.InstallprereqDocker()
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
func runKubernetesAPIHealthCheck(isPreRunExecution bool) error {
	if !isPreRunExecution {
		log.Info("\nKubernetes API \n--------------")
	}
	//Check whether k8s client can be initialized
	client, err := meshkitkube.New([]byte(""))
	if err != nil {
		return errors.New("!! cannot initialize the client")
	}
	if !isPreRunExecution {
		log.Info("√ can initialize the client")
	}

	//Check whether kubernetes can be queried
	podInterface := client.KubeClient.CoreV1().Pods("")
	_, err = podInterface.List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return errors.New("!! cannot query the Kubernetes API")
	}
	if !isPreRunExecution {
		log.Info("√ can query the Kubernetes API")
	}

	return nil
}

//Run healthchecks to verify kubectl and kubenetes version with
// minimum compatible versions
func runKubernetesVersionHealthCheck(isPreRunExecution bool) error {
	if !isPreRunExecution {
		log.Info("\nKubernetes Version \n--------------")
	}
	//Check whether system has minimum supported versions of kubernetes and kubectl
	var kubeVersion *k8sVersion.Info
	kubeVersion, err := utils.GetK8sVersionInfo()
	if err != nil {
		return err
	}
	if err = utils.CheckK8sVersion(kubeVersion); err != nil {
		return err
	}
	if !isPreRunExecution {
		log.Info("√ is running the minimum Kubernetes version")
	}

	if err = utils.CheckKubectlVersion(); err != nil {
		return err
	}
	if !isPreRunExecution {
		log.Info("√ is running the minimum kubectl version")
	}

	return nil
}

func init() {
	checkCmd.Flags().BoolVarP(&preflight, "preflight", "", false, "Verify environment readiness to deploy Meshery")
	checkCmd.Flags().BoolVarP(&pre, "pre", "", false, "Verify environment readiness to deploy Meshery")
}
