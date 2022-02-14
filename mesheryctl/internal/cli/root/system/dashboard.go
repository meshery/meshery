package system

import (
	"context"
	"fmt"
	"net/url"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitutils "github.com/layer5io/meshkit/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Open Meshery UI in browser.",
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
			return errors.New(`meshery server is not running. run "mesheryctl system start" to start meshery`)
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
		if !skipBrowserFlag {
			log.Debug("Fetching Meshery-UI endpoint")
		}

		switch currCtx.GetPlatform() {
		case "docker":
			break
		case "kubernetes":
			var mesheryEndpoint string
			var endpoint *meshkitutils.Endpoint
			kubeClient, err := meshkitkube.New([]byte(""))
			if err != nil {
				return err
			}
			clientset := kubeClient.KubeClient
			var opts meshkitkube.ServiceOptions
			opts.Name = "meshery"
			opts.Namespace = utils.MesheryNamespace
			opts.APIServerURL = kubeClient.RestConfig.Host

			endpoint, err = meshkitkube.GetServiceEndpoint(context.TODO(), clientset, &opts)
			if err != nil {
				return err
			}

			mesheryEndpoint = fmt.Sprintf("%s://%s:%d", utils.EndpointProtocol, endpoint.Internal.Address, endpoint.Internal.Port)
			currCtx.SetEndpoint(mesheryEndpoint)
			if !meshkitutils.TcpCheck(&meshkitutils.HostPort{
				Address: endpoint.Internal.Address,
				Port:    endpoint.Internal.Port,
			}, nil) {
				currCtx.SetEndpoint(fmt.Sprintf("%s://%s:%d", utils.EndpointProtocol, endpoint.External.Address, endpoint.External.Port))
				if !meshkitutils.TcpCheck(&meshkitutils.HostPort{
					Address: endpoint.External.Address,
					Port:    endpoint.External.Port,
				}, nil) {
					u, _ := url.Parse(opts.APIServerURL)
					if meshkitutils.TcpCheck(&meshkitutils.HostPort{
						Address: u.Hostname(),
						Port:    endpoint.External.Port,
					}, nil) {
						mesheryEndpoint = fmt.Sprintf("%s://%s:%d", utils.EndpointProtocol, u.Hostname(), endpoint.External.Port)
						currCtx.SetEndpoint(mesheryEndpoint)
					}
				}
			}

			if err == nil {
				err = config.UpdateContextInConfig(viper.GetViper(), currCtx, mctlCfg.GetCurrentContextName())
				if err != nil {
					return err
				}
			}

		}

		if !skipBrowserFlag {
			log.Info("Opening Meshery (" + currCtx.GetEndpoint() + ") in browser.")
			err = utils.NavigateToBrowser(currCtx.GetEndpoint())
			if err != nil {
				log.Warn("Failed to open Meshery in browser, please point your browser to " + currCtx.GetEndpoint() + " to access Meshery.")
			}
		}

		return nil
	},
}

func init() {
	dashboardCmd.Flags().BoolVarP(&skipBrowserFlag, "skip-browser", "", false, "(optional) skip opening of MesheryUI in browser.")
}
