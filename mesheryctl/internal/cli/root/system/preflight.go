package system

import (
	"context"
	"errors"
	"os/exec"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
)

var preCheckCmd = &cobra.Command{
	Use:   "pre",
	Short: "run basic checks to verify environment for deployment",
	Long:  `To verify environment readiness for a Meshery deloyment `,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("Docker /n--------------")
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

		log.Info("Kubernetes API /n--------------")
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

		return nil
	},
}
