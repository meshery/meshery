// Copyright 2023 Layer5, Inc.
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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitutils "github.com/layer5io/meshkit/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// runPortForward is used for port-forwarding Meshery UI via `system dashboard`
	runPortForward bool
)

// dashboardOptions holds values for command line flags that apply to the dashboard
// command.
type dashboardOptions struct {
	host    string
	port    int
	podPort int
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

// (optional) skip opening of MesheryUI in browser.
mesheryctl system dashboard --skip-browser
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check if meshery is running or not
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(ErrGetCurrentContext(err))
			return nil
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
			utils.Log.Error(err)
			return nil
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(ErrGetCurrentContext(err))
			return nil
		}
		log.Debug("Fetching Meshery-UI endpoint")

		switch currCtx.GetPlatform() {
		case "docker":
			break
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

				// get a free port number to bind port-forwarding
				port, err := utils.GetEphemeralPort()
				if err != nil {
					utils.Log.Error(ErrFailedGetEphemeralPort(err))
					return nil
				}

				portforward, err := utils.NewPortForward(
					cmd.Context(),
					client,
					utils.MesheryNamespace,
					"meshery",
					options.host,
					port,
					options.podPort,
					false,
				)
				if err != nil {
					utils.Log.Error(ErrInitPortForward(err))
					return nil

				}

				if err = portforward.Init(); err != nil {
					// TODO: consider falling back to an ephemeral port if defaultPort is taken
					return ErrRunPortForward(err)
				}
				log.Info("Starting Port-forwarding for Meshery UI")

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
				log.Info(fmt.Sprintf("Forwarding ports %v -> %v", options.podPort, port))
				log.Info("Meshery UI available at: ", mesheryURL)
				log.Info("Opening Meshery UI in the default browser.")
				err = utils.NavigateToBrowser(mesheryURL)
				if err != nil {
					log.Warn("Failed to open Meshery in browser, please point your browser to " + currCtx.GetEndpoint() + " to access Meshery.")
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
				utils.Log.Error(err) //the func return a meshkit error
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
					utils.Log.Error(err)
					return nil
				}
			}

		}

		if !skipBrowserFlag {
			log.Info("Opening Meshery (" + currCtx.GetEndpoint() + ") in browser.")
			err = utils.NavigateToBrowser(currCtx.GetEndpoint())
			if err != nil {
				log.Warn("Failed to open Meshery in your browser, please point your browser to " + currCtx.GetEndpoint() + " to access Meshery.")
			}
		} else {
			log.Info("Meshery UI available at: ", currCtx.GetEndpoint())
		}

		return nil
	},
}

// keepConnectionAlive to stop being timed out with port forwarding
func keepConnectionAlive(url string) {
	_, err := http.Get(url)
	if err != nil {
		log.Debugf("connection request failed %v", err)
	}
	log.Debugf("connection request success")
}

func init() {
	dashboardCmd.Flags().BoolVarP(&runPortForward, "port-forward", "", false, "(optional) Use port forwarding to access Meshery UI")
	dashboardCmd.Flags().BoolVarP(&skipBrowserFlag, "skip-browser", "", false, "(optional) skip opening of MesheryUI in browser.")
}
