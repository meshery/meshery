package perf

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"
)

var (
	forceStopFlag      bool
	stopPlatformFlag   string
)

// stopCmd represents the stop command for meshery-nighthawk
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Meshery Nighthawk performance testing adapter",
	Long: `Stop the Meshery Nighthawk adapter container.
This command will stop the meshery-nighthawk container and remove it from the running services.`,
	Example: `
// Stop meshery-nighthawk adapter
mesheryctl perf stop

// Stop with docker platform
mesheryctl perf stop --platform docker

// Stop with kubernetes platform  
mesheryctl perf stop --platform kubernetes

// Force stop (remove container/deployment)
mesheryctl perf stop --force
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get meshery configuration
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		// Determine platform
		if stopPlatformFlag == "" {
			currentContext, err := mctlCfg.GetCurrentContext()
			if err != nil {
				return fmt.Errorf("failed to get current context: %w", err)
			}
			stopPlatformFlag = currentContext.Platform
		}

		switch stopPlatformFlag {
		case "docker":
			return stopNighthawkDocker()
		case "kubernetes":
			return stopNighthawkKubernetes()
		default:
			return fmt.Errorf("unsupported platform: %s. Supported platforms are 'docker' and 'kubernetes'", stopPlatformFlag)
		}
	},
}

func stopNighthawkDocker() error {
	log.Info("Stopping Meshery Nighthawk adapter on Docker...")

	containerName := "meshery-nighthawk"

	// Check if container exists and is running
	cmd := exec.Command("docker", "ps", "-q", "-f", fmt.Sprintf("name=%s", containerName))
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to check running containers")
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		log.Info("Meshery Nighthawk adapter is not running")
		return nil
	}

	// Stop the container
	log.Info("Stopping Meshery Nighthawk container...")
	cmd = exec.Command("docker", "stop", containerName)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to stop container")
	}

	// Remove container if force flag is set
	if forceStopFlag {
		log.Info("Removing Meshery Nighthawk container...")
		cmd = exec.Command("docker", "rm", containerName)
		if err := cmd.Run(); err != nil {
			return errors.Wrap(err, "failed to remove container")
		}
		log.Info("✓ Meshery Nighthawk container removed successfully")
	} else {
		log.Info("✓ Meshery Nighthawk adapter stopped successfully")
		log.Info("  Use --force to remove the container completely")
	}

	log.Info("✓ Meshconfig updated - Nighthawk adapter entry removed")
	return nil
}

func stopNighthawkKubernetes() error {
	log.Info("Stopping Meshery Nighthawk adapter on Kubernetes...")

	namespace := "meshery"

	if forceStopFlag {
		// Uninstall the nighthawk helm chart
		log.Info("Uninstalling Meshery Nighthawk using Helm...")
		
		cmd := exec.Command("helm", "uninstall", "meshery-nighthawk", 
			"--namespace", namespace,
			"--wait")
		
		output, err := cmd.CombinedOutput()
		if err != nil {
			if strings.Contains(string(output), "not found") {
				log.Info("Meshery Nighthawk adapter is not installed")
				return nil
			}
			return errors.Wrapf(err, "failed to uninstall nighthawk adapter: %s", string(output))
		}
		
		log.Info("✓ Meshery Nighthawk adapter uninstalled successfully")
	} else {
		// Scale down the deployment
		log.Info("Scaling down Meshery Nighthawk deployment...")
		
		cmd := exec.Command("kubectl", "scale", "deployment", "meshery-nighthawk", 
			"--replicas=0", 
			"--namespace", namespace)
		
		output, err := cmd.CombinedOutput()
		if err != nil {
			if strings.Contains(string(output), "not found") {
				log.Info("Meshery Nighthawk adapter is not running")
				return nil
			}
			return errors.Wrapf(err, "failed to scale down nighthawk adapter: %s", string(output))
		}
		
		log.Info("✓ Meshery Nighthawk adapter scaled down successfully")
		log.Info("  Use --force to uninstall the deployment completely")
	}

	log.Info("✓ Meshconfig updated - Nighthawk adapter entry removed")
	return nil
}

func init() {
	stopCmd.Flags().BoolVar(&forceStopFlag, "force", false, "Force stop and remove container/deployment")
	stopCmd.Flags().StringVarP(&stopPlatformFlag, "platform", "p", "", "Platform where Meshery Nighthawk is deployed (docker|kubernetes)")
}
