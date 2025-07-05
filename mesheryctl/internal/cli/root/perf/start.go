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

package perf

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"
)

var (
	skipUpdateConfigFlag bool
	platformFlag         string
)

// startCmd represents the start command for meshery-nighthawk
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Meshery Nighthawk performance testing adapter",
	Long: `Start the Meshery Nighthawk adapter container for performance testing.
This command will start the meshery-nighthawk container and optionally add an entry to the meshconfig.`,
	Example: `
// Start meshery-nighthawk adapter
mesheryctl perf start

// Start with docker platform
mesheryctl perf start --platform docker

// Start with kubernetes platform  
mesheryctl perf start --platform kubernetes

// Start without updating meshconfig
mesheryctl perf start --skip-update-config
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get meshery configuration
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// Get current context
		currentContext, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return errors.Wrap(err, "failed to get current context")
		}

		// Determine platform
		if platformFlag == "" {
			platformFlag = currentContext.Platform
		}

		switch platformFlag {
		case "docker":
			return startNighthawkDocker(mctlCfg)
		case "kubernetes":
			return startNighthawkKubernetes(mctlCfg)
		default:
			return fmt.Errorf("unsupported platform: %s. Supported platforms are 'docker' and 'kubernetes'", platformFlag)
		}
	},
}

func startNighthawkDocker(mctlCfg *config.MesheryCtlConfig) error {
	log.Info("Starting Meshery Nighthawk adapter on Docker...")

	containerName := "meshery-nighthawk"
	imageName := "meshery/meshery-nighthawk:edge-latest"

	// Check if container already exists
	cmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=%s", containerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to check for existing container")
	}

	if strings.TrimSpace(string(output)) == containerName {
		// Container exists, check if it's running
		cmd = exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", containerName), "--format", "{{.Names}}")
		output, err = cmd.Output()
		if err != nil {
			return errors.Wrap(err, "failed to check container status")
		}

		if strings.TrimSpace(string(output)) == containerName {
			log.Info("Meshery Nighthawk adapter is already running")
			return updateMeshConfig(mctlCfg, containerName)
		}

		// Container exists but not running, start it
		log.Info("Starting existing Meshery Nighthawk container...")
		cmd = exec.Command("docker", "start", containerName)
		if output, err := cmd.CombinedOutput(); err != nil {
			return errors.Wrapf(err, "failed to start existing container: %s", string(output))
		}
	} else {
		// Pull the latest image
		log.Info("Pulling Meshery Nighthawk image...")
		cmd = exec.Command("docker", "pull", imageName)
		if output, err := cmd.CombinedOutput(); err != nil {
			return errors.Wrapf(err, "failed to pull image: %s", string(output))
		}

		// Create and start new container
		log.Info("Creating new Meshery Nighthawk container...")
		cmd = exec.Command("docker", "run", "-d",
			"--name", containerName,
			"-p", "10013:10013",
			"--restart", "unless-stopped",
			imageName)
		
		if output, err := cmd.CombinedOutput(); err != nil {
			return errors.Wrapf(err, "failed to create and start container: %s", string(output))
		}
	}

	// Wait for container to be ready
	log.Info("Waiting for Meshery Nighthawk adapter to be ready...")
	time.Sleep(5 * time.Second)

	log.Info("✓ Meshery Nighthawk adapter started successfully on Docker")
	log.Info("  Adapter is available at: http://localhost:10013")

	// Update mesh config if not skipped
	if !skipUpdateConfigFlag {
		return updateMeshConfig(mctlCfg, containerName)
	}

	return nil
}

func startNighthawkKubernetes(mctlCfg *config.MesheryCtlConfig) error {
	log.Info("Starting Meshery Nighthawk adapter on Kubernetes...")

	// Get current context
	_, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return errors.Wrap(err, "failed to get current context")
	}

	// Get kubeconfig path
	kubeconfig := ""
	if kubeconfig == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return errors.Wrap(err, "failed to get home directory")
		}
		kubeconfig = filepath.Join(homeDir, ".kube", "config")
	}

	namespace := "meshery"

	// Check if meshery-nighthawk is already deployed
	cmd := exec.Command("kubectl", "get", "deployment", "meshery-nighthawk", "-n", namespace)
	if err := cmd.Run(); err == nil {
		log.Info("Meshery Nighthawk adapter is already deployed")
		// Scale up if it's scaled down
		cmd = exec.Command("kubectl", "scale", "deployment", "meshery-nighthawk", "--replicas=1", "-n", namespace)
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Warnf("Failed to scale up deployment: %s", string(output))
		}
		return updateMeshConfig(mctlCfg, "meshery-nighthawk")
	}

	// Apply the nighthawk helm chart
	log.Info("Deploying Meshery Nighthawk using Helm...")
	
	// Use helm command to install
	cmd = exec.Command("helm", "install", "meshery-nighthawk", 
		"./install/kubernetes/helm/meshery/charts/meshery-nighthawk/",
		"--namespace", namespace,
		"--create-namespace",
		"--wait")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "already exists") {
			log.Info("Meshery Nighthawk adapter is already installed")
			// Try to upgrade instead
			cmd = exec.Command("helm", "upgrade", "meshery-nighthawk",
				"./install/kubernetes/helm/meshery/charts/meshery-nighthawk/",
				"--namespace", namespace,
				"--wait")
			output, err = cmd.CombinedOutput()
			if err != nil {
				return errors.Wrapf(err, "failed to upgrade nighthawk adapter: %s", string(output))
			}
		} else {
			return errors.Wrapf(err, "failed to install nighthawk adapter: %s", string(output))
		}
	}

	log.Info("✓ Meshery Nighthawk adapter started successfully on Kubernetes")
	log.Info("  Check status with: kubectl get pods -n meshery | grep nighthawk")

	// Update mesh config if not skipped
	if !skipUpdateConfigFlag {
		return updateMeshConfig(mctlCfg, "meshery-nighthawk")
	}

	return nil
}

func updateMeshConfig(mctlCfg *config.MesheryCtlConfig, adapterName string) error {
	if skipUpdateConfigFlag {
		return nil
	}

	log.Info("Updating meshconfig with Nighthawk adapter entry...")

	// Get current context
	ctx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return errors.Wrap(err, "failed to get current context")
	}

	// Add adapter to the context components
	// Check if adapter is already in the components list
	found := false
	for _, component := range ctx.Components {
		if component == adapterName {
			found = true
			break
		}
	}

	if !found {
		ctx.Components = append(ctx.Components, adapterName)
		log.Info("✓ Meshconfig updated with Nighthawk adapter entry")
	} else {
		log.Info("✓ Nighthawk adapter already exists in meshconfig")
	}
	
	return nil
}

func init() {
	startCmd.Flags().BoolVar(&skipUpdateConfigFlag, "skip-update-config", false, "Skip adding entry to meshconfig")
	startCmd.Flags().StringVarP(&platformFlag, "platform", "p", "", "Platform to deploy Meshery Nighthawk to (docker|kubernetes)")
}
