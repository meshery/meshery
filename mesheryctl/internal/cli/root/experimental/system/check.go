// Copyright 2020 Layer5, Inc.
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
	"os/exec"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	preflightChecks bool = false
)

// StopCmd represents the stop command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run checks on Meshery to look for potential problems",
	Long: `Run checks on Meshery to look for potential problems

The check command will perform a series of checks to validate that Meshery is configured correctly.
It checks the pre-requisites and also the status of the Meshery components.`,
	Args: cobra.NoArgs,
	// XXX (nitishm): Do we need to run this as a pre-run?
	//PreRunE: func(cmd *cobra.Command, args []string) error {
	//	return utils.PreReqCheck(cmd.Use)
	//},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("Checking Meshery status...")
		if err := check(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to run Meshery checks"))
		}
		return nil
	},
}

func check() error {
	// XXX (nitishm): Do we require meshery to by up and running. Possibly no since the check command
	// should have the capability to respond with meshery server not running.

	// First check if Meshery is up and running.
	//if !utils.IsMesheryRunning() {
	//	log.Info("Meshery is not running. Nothing to check.")
	//	return nil
	//}

	// Check docker status:
	if err := exec.Command("docker", "ps").Run(); err != nil {
		log.Info("[-] Docker daemon is not running")
	} else {
		log.Info("[√] Docker daemon is running")
	}

	// Check docker-compose status
	if err := exec.Command("docker-compose", "-v").Run(); err != nil {
		log.Info("[-] docker-compose is not available")
	} else {
		log.Info("[√] docker-compose is available")
	}

	// TODO:
	// Check kubernetes status:

	// Kubernetes API
	// --------------
	// √ can initialize the client
	// √ can query the Kubernetes API
	// √ has required level of privileges
	//
	// Kubernetes Version
	// ------------------
	// √ is running the minimum Kubernetes API version
	// √ is running the minimum kubectl version

	if preflightChecks {
		return nil
	}

	// TODO:
	// Check Meshery versions:

	// TODO:
	// Meshery Version
	// ---------------
	// √ server is up-to-date (stable-v0.x.x)
	// √ cli is up-to-date (stable-v0.x.x)

	// TODO:
	// Check meshery adapter status

	// TODO:
	// Meshery Adapters
	// ---------------
	// √ meshery-istio adapter is reachable and running
	// √ meshery-consul adapter is reachable and running
	// √ meshery-linkerd adapter is reachable and running
	// .
	// .
	// .
	// √ meshery-nginx-sm adapter is reachable and running
	// √ meshery-octarine adapter is reachable and running

	return nil
}

func init() {
	checkCmd.Flags().BoolVarP(&preflightChecks, "pre", "",
		false, "(optional) Checks status of docker and kubernetes only.")

}
