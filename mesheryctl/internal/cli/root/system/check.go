package system

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sVersion "k8s.io/apimachinery/pkg/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	preflight bool
	pre       bool
	failure   int = 0
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Meshery environment check",
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
	//Kubernetes and KubernetesVersion healthchecks are invoked in
	//both the cases PreRunExecution and !PreRunExecution
	//Run k8s API healthchecks
	if err = runKubernetesAPIHealthCheck(isPreRunExecution); err != nil {
		return err
	}
	//Run k8s plus kubectl minimum version healthchecks
	if err := runKubernetesVersionHealthCheck(isPreRunExecution); err != nil {
		return err
	}
	if !isPreRunExecution {
		if failure == 0 {
			fmt.Printf("\n--------------\n--------------\n✓✓ Meshery prerequisites met")
		} else {
			fmt.Printf("\n--------------\n--------------\n!! Meshery prerequisites not met")
		}
	}
	return nil
}

//Run healthchecks to verify if docker is running and active
func runDockerHealthCheck(isPreRunExecution bool, subcommand string, platform string) error {
	if !isPreRunExecution {
		fmt.Printf("\nDocker \n--------------\n")
	}
	//Check whether docker daemon is running or not
	err := exec.Command("docker", "ps").Run()
	if err != nil && isPreRunExecution { // if this is PreRunExec we trigger self installation
		fmt.Printf("!! Docker is not running\n")
		//If preRunExecution and the current platform is docker then we trigger docker installation
		//No auto installation of docker for windows
		if runtime.GOOS == "windows" {
			return errors.Wrapf(err, "Please start Docker. Run `mesheryctl system %s` once Docker is started.", subcommand)
		}
		err = utils.Startdockerdaemon(subcommand)
		if err != nil {
			return errors.Wrapf(err, "failed to start Docker.")
		}
	} else if err != nil { // warn incase of preflight check
		fmt.Printf("!! Docker is not running\n")
	} else if !isPreRunExecution { // log incase of preflight check
		fmt.Printf("√ Docker is running\n")
	}

	//Check for installed docker-compose on client system
	err = exec.Command("docker-compose", "-v").Run()
	if err != nil && isPreRunExecution { // if PreRunExec we trigger self installation
		fmt.Printf("!! docker-compose is not available\n")
		//No auto installation of Docker-compose for windows
		if runtime.GOOS == "windows" {
			return errors.Wrapf(err, "please install docker-compose. Run `mesheryctl system %s` after docker-compose is installed.", subcommand)
		}
		err = utils.InstallprereqDocker()
		if err != nil {
			return errors.Wrapf(err, "failed to install prerequisites. Run `mesheryctl system %s` after docker-compose is installed.", subcommand)
		}
	} else if err != nil { // this is preflight check we just log a warning
		fmt.Printf("!! docker-compose is not available\n")
	} else if !isPreRunExecution { // log incase of preflight check
		fmt.Printf("√ docker-compose is available\n")
	}

	return nil
}

//Run healthchecks to verify if kubernetes client can be initialized and can be queried
func runKubernetesAPIHealthCheck(isPreRunExecution bool) error {
	if !isPreRunExecution {
		fmt.Printf("\nKubernetes API \n--------------\n")
	}
	//Check whether k8s client can be initialized
	client, err := meshkitkube.New([]byte(""))
	if err != nil && !isPreRunExecution { // this is preflight check
		fmt.Printf("!! cannot initialize Kubernetes client\n")
		fmt.Printf("!! cannot query the Kubernetes API\n")
		failure++ // increase failure count
		return nil
	} else if err != nil {
		return errors.New("ctlK8sClient1000: !! cannot initialize a Kubernetes client. See https://docs.meshery.io/reference/error-codes")
	} else if !isPreRunExecution { // log in case of preflight check
		fmt.Printf("√ can initialize Kubernetes client\n")
	}

	//Check whether kubernetes can be queried
	podInterface := client.KubeClient.CoreV1().Pods("")
	_, err = podInterface.List(context.TODO(), v1.ListOptions{})
	if err != nil && !isPreRunExecution { // this is preflight check
		fmt.Printf("!! cannot query the Kubernetes API\n")
		failure++ // increase failure count
	} else if err != nil {
		return errors.New("ctlK8sConnect1001: !! cannot query the Kubernetes API. See https://docs.meshery.io/reference/error-codes")
	} else if !isPreRunExecution { // log in case of preflight check
		fmt.Printf("√ can query the Kubernetes API\n")
	}

	return nil
}

//Run healthchecks to verify kubectl and kubenetes version with
// minimum compatible versions
func runKubernetesVersionHealthCheck(isPreRunExecution bool) error {
	if !isPreRunExecution {
		fmt.Printf("\nKubernetes Version \n--------------\n")
	}
	//Check whether system has minimum supported versions of kubernetes and kubectl
	var kubeVersion *k8sVersion.Info
	kubeVersion, err := utils.GetK8sVersionInfo()
	if err != nil {
		// probably kubernetes isn't running
		if !isPreRunExecution { // log incase of preflight checks
			fmt.Printf("!! cannot check Kubernetes version\n")
		} else { // return as error in case of PreRunE
			return err
		}
		failure++ // increase failure count
	} else {
		// kubernetes is running so check the version
		err = utils.CheckK8sVersion(kubeVersion)
		if err != nil && !isPreRunExecution { // warn incase of preflight check
			fmt.Printf("%s\n", err)
			failure++
		} else if err != nil { // incase of PreRunExec return the error
			return err
		} else if !isPreRunExecution { // log incase of perflight check
			fmt.Printf("√ is running the minimum Kubernetes version\n")
		}
	}

	err = utils.CheckKubectlVersion()
	if err != nil && !isPreRunExecution { // warn incase of preflight check
		fmt.Printf("%s\n", err)
		failure++
	} else if err != nil {
		return err
	} else if !isPreRunExecution { // log incase of preflight check
		fmt.Printf("√ is running the minimum kubectl version\n")
	}

	return nil
}

func init() {
	checkCmd.Flags().BoolVarP(&preflight, "preflight", "", false, "Verify environment readiness to deploy Meshery")
	checkCmd.Flags().BoolVarP(&pre, "pre", "", false, "Verify environment readiness to deploy Meshery")
}
