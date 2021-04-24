package system

import (
	"context"
	"errors"
	"os/exec"
	"strings"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
)

var PreCheckCmd = &cobra.Command{
	Use:   "preflight",
	Short: "run basic checks to verify environment for deployment",
	Long:  `To verify environment readiness for a Meshery deloyment `,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPreDeploymentChecks()
	},
}

func runPreDeploymentChecks() error {
	log.Info("\nDocker \n--------------")
	//Check whether docker daemon is running or not
	if err := exec.Command("docker", "ps").Run(); err != nil {
		log.Error("Docker is not running")
	} else {
		log.Info("√ Docker is running")
	}
	//Check for installed docker-compose on client system
	if err := exec.Command("docker-compose", "-v").Run(); err != nil {
		log.Info("Docker-Compose is not available")
	} else {
		log.Info("√ docker-compose is available")
	}

	log.Info("\nKubernetes API \n--------------")
	//Check whether k8s client can be initialized
	client, err := meshkitkube.New([]byte(""))
	if err != nil {
		return errors.New("cannot initialize the client")
	} else {
		log.Info("√ can initialize the client")
	}
	//Check whether kubernetes can be queried
	podInterface := client.KubeClient.CoreV1().Pods("")
	_, err = podInterface.List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return errors.New("cannot query the Kubernetes API")
	} else {
		log.Info("√ can query the Kubernetes API")
	}

	log.Info("\nKubernetes Version \n--------------")
	//Check whether system has minimum supported versions of kubernetes and kubectl
	minK8sVersion := "1.12.0"
	minVersionArray := strings.Split(minK8sVersion, ".")

	//Check kubernetes version with minimum version
	serverOutput, err := exec.Command("bash", "-c", "kubectl version --short | grep -i server").Output()
	if err != nil {
		return errors.New("failed to fetch kubernetes version")
	}
	serverVersion := string(serverOutput)
	spliter := strings.Split(serverVersion, "v")
	spliter = strings.Split(spliter[1], ".")
	if spliter[0] > minVersionArray[0] {
		log.Info("√ is running the minimum Kubernetes version")
	} else if spliter[0] == minVersionArray[0] && spliter[1] >= minVersionArray[1] {
		log.Info("√ is running the minimum Kubernetes version")
	} else {
		return errors.New("system isn't running the minimum Kubernetes version")
	}

	//Check kubectl version with minimum version
	clientOutput, err := exec.Command("bash", "-c", "kubectl version --short | grep -i client").Output()
	if err != nil {
		return errors.New("failed to fetch kubectl version")
	}
	clientVersion := string(clientOutput)
	spliter = strings.Split(clientVersion, "v")
	spliter = strings.Split(spliter[1], ".")
	if spliter[0] > minVersionArray[0] {
		log.Info("√ is running the minimum kubectl version")
	} else if spliter[0] == minVersionArray[0] && spliter[1] >= minVersionArray[1] {
		log.Info("√ is running the minimum kubectl version")
	} else {
		return errors.New("system isn't running the minimum kubectl version")
	}

	return nil
}
