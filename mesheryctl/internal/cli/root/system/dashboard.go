package system

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	mesheryServiceName = "meshery"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Open Meshery UI",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check if meshery is running or not
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		running, _ := utils.IsMesheryRunning(currCtx.GetPlatform())
		if !running {
			return errors.New(`meshery server is not running. run "mesheryctl system start" to start meshery.`)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}

		switch currCtx.GetPlatform() {
		case "docker":
			log.Info("Opening Meshery (" + currCtx.GetEndpoint() + ") in browser.")
			err = utils.NavigateToBrowser(currCtx.GetEndpoint())
			if err != nil {
				log.Warn("Failed to open Meshery in browser, please point your browser to " + currCtx.GetEndpoint() + " to access Meshery.")
			}
		case "kubernetes":
			return OpenTunnelToMeshery()
		}

		return nil
	},
}

func OpenTunnelToMeshery() error {
	client, err := meshkitkube.New([]byte(""))
	if err != nil {
		return err
	}

	// get v1 interface to our cluster
	v1ClientSet := client.KubeClient.CoreV1()

	svc, err := v1ClientSet.Services(utils.MesheryNamespace).Get(context.TODO(), mesheryServiceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Select pod/s given the service data available
	set := labels.Set(svc.Spec.Selector)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	pods, err := v1ClientSet.Pods(utils.MesheryNamespace).List(context.TODO(), listOptions)
	if err != nil {
		return fmt.Errorf("error listing pods: %s", err)
	}

	// Will select first running Pod available
	var mesheryPod *corev1.Pod

	for _, pod := range pods.Items {
		pod := pod // prevents aliasing address of loop variable which is the same in each iteration
		if pod.Status.Phase == "Running" {
			mesheryPod = &pod
			break
		}
	}
	if mesheryPod == nil {
		return fmt.Errorf("no running meshery pod available")
	}

	dialer, err := utils.DialerToPod(&client.RestConfig, client.KubeClient, mesheryPod.Name, mesheryPod.Namespace)
	if err != nil {
		return err
	}

	var opts meshkitkube.ServiceOptions
	opts.Name = "meshery"
	opts.Namespace = utils.MesheryNamespace
	opts.APIServerURL = client.RestConfig.Host
	endpoint, err := meshkitkube.GetServiceEndpoint(context.TODO(), client.KubeClient, &opts)
	if err != nil {
		return err
	}

	portForwarder, err := utils.NewPortForwarder(dialer, fmt.Sprintf("%d:%d", endpoint.Internal.Port, svc.Spec.Ports[0].TargetPort.IntValue()))
	if err != nil {
		return fmt.Errorf("error setting up port forwarding: %s", err)
	}

	err = portForwarder.Start(func(*utils.PortForwarder) error {
		url := fmt.Sprintf("http://localhost:%d", endpoint.Internal.Port)
		log.Info("Opening Meshery (" + url + ") in browser.")
		err = utils.NavigateToBrowser(url)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("port forwarding failed: %s", err)
	}

	// The command should only exit when a signal is received from the OS.
	// Exiting before will result in port forwarding to stop causing the browser
	// if open to not render the dashboard.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// portforwarder.Stop() triggered implicitly by SIGINT. Ensure it completes
	// before exiting.
	<-portForwarder.Done()

	return nil
}
