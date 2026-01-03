// Copyright Meshery Authors
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
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitutils "github.com/meshery/meshkit/utils"
	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// runPortForward is used for port-forwarding Meshery UI via `system dashboard`
	runPortForward bool
	localPort      int
)

// dashboardOptions holds values for command line flags that apply to the dashboard
// command.
type dashboardOptions struct {
	host    string // Host on which server is running inside the pod
	port    int    // The default port on which Meshery service is listening
	podPort int    // Port on which server is running inside the pod
}

// newDashboardOptions initializes dashboard options with default
// values for host, port, and which dashboard to show. Also, set
// max wait time duration for 300 seconds for the dashboard to
// become available
//
// These options may be overridden on the CLI at run-time
func newDashboardOptions() *dashboardOptions {
	return &dashboardOptions{
		host:    utils.MesheryDefaultHost,
		port:    utils.MesheryDefaultPort,
		podPort: 8080,
	}
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Open Meshery UI in browser.",
	Args:  cobra.NoArgs,
	Example: `
// Open Meshery UI in browser
mesheryctl system dashboard

// Open Meshery UI in browser and use port-forwarding (if default port is taken already)
mesheryctl system dashboard --port-forward

// Open Meshery UI in browser and use port-forwarding, listen on port 9081 locally, forwarding traffic to meshery server in the pod
mesheryctl system dashboard --port-forward -p 9081

// (optional) skip opening of MesheryUI in browser.
mesheryctl system dashboard --skip-browser

Note: Meshery's web-based user interface is embedded in Meshery Server and is available as soon as Meshery starts. The location and port that Meshery UI is exposed varies depending upon your mode of deployment. See accessing \"Meshery UI\" for additional deployment-specific options: https://docs.meshery.io/installation/accessing-meshery-ui.`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check if meshery is running or not
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.LogError.Error(err)
			return nil
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.LogError.Error(ErrGetCurrentContext(err))
			return nil
		}
		running, err := utils.IsMesheryRunning(currCtx.GetPlatform())
		if err != nil {
			utils.LogError.Error(err)
			return nil
		}
		if !running {
			utils.LogError.Error(utils.ErrMesheryServerNotRunning(currCtx.GetPlatform()))
			return nil
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.LogError.Error(err)
			return nil
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.LogError.Error(ErrGetCurrentContext(err))
			return nil
		}
		utils.Log.Debug("Fetching Meshery-UI endpoint")
		switch currCtx.GetPlatform() {
		case "docker":
			if runPortForward {
				utils.Log.Warnf("--port-forward is not supported using Docker as Meshery's deployment platform.")
			}
		case "kubernetes":
			client, err := meshkitkube.New([]byte(""))
			if err != nil {
				return err
			}

			// Run port forwarding for accessing Meshery UI
			if runPortForward {
				options := newDashboardOptions()

				signals := make(chan os.Signal, 1)
				signal.Notify(signals, os.Interrupt)
				defer signal.Stop(signals)

				portforward, err := utils.NewPortForward(
					cmd.Context(),
					client,
					utils.MesheryNamespace,
					"meshery",
					options.host,
					localPort,
					options.podPort,
					false,
				)
				if err != nil {
					utils.LogError.Error(ErrInitPortForward(err))
					return nil

				}

				if err = portforward.Init(); err != nil {
					// TODO: consider falling back to an ephemeral port if defaultPort is taken
					return ErrRunPortForward(err)
				}
				utils.Log.Info("Port-forwarding for Meshery UI...")

				mesheryURL := portforward.URLFor("")

				// ticker for keeping connection alive with pod each 10 seconds
				ticker := time.NewTicker(10 * time.Second)
				go func() {
					for {
						select {
						case <-signals:
							portforward.Stop()
							ticker.Stop()
							return
						case <-ticker.C:
							keepConnectionAlive(mesheryURL)
						}
					}
				}()
				utils.Log.Info(fmt.Sprintf("Forwarding port %v -> %v", options.podPort, localPort))
				utils.Log.Info("Meshery UI available at: ", mesheryURL)
				utils.Log.Info("Opening Meshery UI in default browser...")
				err = utils.NavigateToBrowser(mesheryURL)
				if err != nil {
					utils.Log.Warn(errors.Wrapf(err, "Opening Meshery UI in browser at %v.", currCtx.GetEndpoint()))
				}

				<-portforward.GetStop()
				return nil
			}

			var mesheryEndpoint string
			var endpoint *meshkitutils.Endpoint
			clientset := client.KubeClient
			var opts meshkitkube.ServiceOptions
			opts.Name = "meshery"
			opts.Namespace = utils.MesheryNamespace
			opts.APIServerURL = client.RestConfig.Host

			endpoint, err = meshkitkube.GetServiceEndpoint(context.TODO(), clientset, &opts)
			if err != nil {
				utils.LogError.Error(err) //the func return a meshkit error
				return nil
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
				err = config.UpdateContextInConfig(currCtx, mctlCfg.GetCurrentContextName())
				if err != nil {
					utils.LogError.Error(err)
					return nil
				}
			}

		}

		if !skipBrowserFlag {
			utils.Log.Info("Opening Meshery UI in browser at " + currCtx.GetEndpoint() + ".")
			err = utils.NavigateToBrowser(currCtx.GetEndpoint())
			if err != nil {
				utils.Log.Warn(errors.Wrap(err, "Failed to open Meshery UI in your browser, please point your browser to " + currCtx.GetEndpoint() + " to access Meshery UI.\n\nOr run `mesheryctl system dashboard --port-forward` to access Meshery UI via port-forwarding."))
			}
		} else {
			utils.Log.Info("Meshery UI available at: ", currCtx.GetEndpoint())
		}

		return nil
	},
}

// keepConnectionAlive to stop being timed out with port forwarding
func keepConnectionAlive(url string) {
	_, err := http.Get(url)
	if err != nil {
		utils.Log.Debugf("connection request failed %v", err)
	}
	utils.Log.Debugf("connection request success")
}

func init() {
	dashboardCmd.Flags().BoolVarP(&runPortForward, "port-forward", "", false, "(optional) Use port forwarding to access Meshery UI")
	dashboardCmd.Flags().IntVarP(&localPort, "port", "p", 9081, "(optional) Local port that is not in use from which traffic is to be forwarded to the server running inside the Pod.")

	dashboardCmd.Flags().BoolVarP(&skipBrowserFlag, "skip-browser", "", false, "(optional) skip opening of MesheryUI in browser.")
}
