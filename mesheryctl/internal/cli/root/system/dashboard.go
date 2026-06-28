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
	"strings"
	"syscall"
	"time"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitutils "github.com/meshery/meshkit/utils"
	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cmdDashboardFlags holds the flag values for the dashboard command.
type cmdDashboardFlags struct {
	PortForward bool `json:"port-forward" validate:"boolean"`
	Port        int  `json:"port" validate:"omitempty"`
	SkipBrowser bool `json:"skip-browser" validate:"boolean"`
}

var dashboardCmdFlags cmdDashboardFlags

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
		// validate flags
		if err := mesheryctlflags.ValidateCmdFlags(cmd, &dashboardCmdFlags); err != nil {
			return err
		}
		// check if meshery is running or not
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return ErrGetCurrentContext(err)
		}
		running, err := utils.IsMesheryRunning(currCtx.GetPlatform())
		if err != nil {
			return err
		}
		if !running {
			return utils.ErrMesheryServerNotRunning(currCtx.GetPlatform())
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return ErrGetCurrentContext(err)
		}
		utils.Log.Debugf("Fetching Meshery-UI endpoint for platform: %s", currCtx.GetPlatform())
		switch currCtx.GetPlatform() {
		case platformDocker:
			if dashboardCmdFlags.PortForward {
				utils.Log.Warn(fmt.Errorf("--port-forward is not supported using Docker as Meshery's deployment platform"))
			}
		case platformKubernetes:
			client, err := meshkitkube.New([]byte(""))
			if err != nil {
				return err
			}

			// Run port forwarding for accessing Meshery UI
			if dashboardCmdFlags.PortForward {
				options := newDashboardOptions()

				signals := make(chan os.Signal, 1)
				signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
				defer signal.Stop(signals)

				// Check if user explicitly specified a port
				portExplicitlySet := cmd.Flags().Changed("port")

				// Attempt to create port forward connection
				// First try with the requested port, then fallback to ephemeral if not explicitly set
				var portforward *utils.PortForward
				var err error

				portforward, err = utils.NewPortForward(
					cmd.Context(),
					client,
					utils.MesheryNamespace,
					"meshery",
					options.host,
					dashboardCmdFlags.Port,
					options.podPort,
					false,
				)
				if err != nil {
					return ErrInitPortForward(err)
				}

				if err = portforward.Init(); err != nil {
					if !isPortAlreadyInUseError(err) {
						return ErrRunPortForward(err)
					}

					// If local listen port is in use and user didn't explicitly set it, try with ephemeral port.
					if shouldUseEphemeralPortFallback(portExplicitlySet) {
						utils.Log.Warnf("Port %d is already in use, attempting to use an ephemeral port...", dashboardCmdFlags.Port)

						portforward, err = utils.NewPortForward(
							cmd.Context(),
							client,
							utils.MesheryNamespace,
							"meshery",
							options.host,
							0, // Use ephemeral port
							options.podPort,
							false,
						)
						if err != nil {
							return ErrInitPortForward(err)
						}

						if err = portforward.Init(); err != nil {
							return ErrRunPortForward(err)
						}
						utils.Log.Info("Successfully bound to ephemeral port")
					} else if portExplicitlySet {
						// User explicitly set port and the port is in use.
						return fmt.Errorf("port %d is already in use. Please specify a different port with -p flag or omit the flag to use an ephemeral port", dashboardCmdFlags.Port)
					}
				}
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
				utils.Log.Infof("Port forwarding established - pod port %d forwarded to local port %d", options.podPort, portforward.GetLocalPort())
				err = utils.NavigateToBrowser(mesheryURL)
				if err != nil {
					utils.Log.Warnf("Opening Meshery UI in browser at %s.", currCtx.GetEndpoint())
				}

				<-portforward.GetStop()
				return nil
			}

			var mesheryEndpoint string
			endpoint, err := utils.GetMesheryEndpoint(context.TODO(), client)
			if err != nil {
				utils.Log.Debugf("Error while GetMesheryEndpoint\n- Endpoint: %v\n- Error: %v", endpoint, err)
				return err //the func return a meshkit error
			}

			mesheryEndpoint = fmt.Sprintf("%s://%s:%d", utils.EndpointProtocol, endpoint.Internal.Address, endpoint.Internal.Port)
			utils.Log.Debugf("Set Meshery Endpoint: %s", mesheryEndpoint)
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
					u, _ := url.Parse(client.RestConfig.Host)
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
					return err
				}
			}

		}

		if !dashboardCmdFlags.SkipBrowser {
			utils.Log.Infof("Opening Meshery UI in browser at %s", currCtx.GetEndpoint())
			err = utils.NavigateToBrowser(currCtx.GetEndpoint())
			if err != nil {
				utils.Log.Warnf("Failed to open Meshery UI in your browser, please point your browser to %s to access Meshery UI.\n\nOr run `mesheryctl system dashboard --port-forward` to access Meshery UI via port-forwarding.", currCtx.GetEndpoint())
			}
		} else {
			utils.Log.Infof("Meshery UI available at: %s", currCtx.GetEndpoint())
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

func shouldUseEphemeralPortFallback(portExplicitlySet bool) bool {
	return !portExplicitlySet
}

func isPortAlreadyInUseError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "address already in use") ||
		strings.Contains(errStr, "bind") ||
		strings.Contains(errStr, "unable to listen on any of the requested ports")
}

func init() {
	dashboardCmd.Flags().BoolVarP(&dashboardCmdFlags.PortForward, "port-forward", "", false, "(optional) Use port forwarding to access Meshery UI")
	dashboardCmd.Flags().IntVarP(&dashboardCmdFlags.Port, "port", "p", utils.MesheryDefaultPort, "(optional) Local port that is not in use from which traffic is to be forwarded to the server running inside the Pod. If not specified and default port is in use, an ephemeral port will be used.")
	dashboardCmd.Flags().BoolVarP(&dashboardCmdFlags.SkipBrowser, "skip-browser", "", false, "(optional) skip opening of MesheryUI in browser.")
}
