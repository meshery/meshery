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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	dockerclient "github.com/docker/docker/client"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/constants"
	c "github.com/meshery/meshery/mesheryctl/pkg/constants"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/handlers"
	"github.com/meshery/meshery/server/models"
	meshkitutils "github.com/meshery/meshkit/utils"
	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	failure int
)

type cmdSystemCheckFlags struct {
	Preflight      bool   `json:"preflight" validate:"boolean"`
	Pre            bool   `json:"pre" validate:"boolean"`
	ComponentsFlag bool   `json:"componentsFlag" validate:"boolean"`
	AdaptersFlag   bool   `json:"adaptersFlag" validate:"boolean"`
	OperatorsFlag  bool   `json:"operatorsFlag" validate:"boolean"`
	Adapter        string `json:"adapter" validate:"omitempty"`
}

var systemCheckFlags cmdSystemCheckFlags

type HealthCheckOptions struct {
	// PrintLogs to keep incheck to print logs during a healthcheck
	PrintLogs bool
	// IsPreRunE to keep incheck if healthchecks are triggered as as a PreRunE
	IsPreRunE bool

	// the command under which the healthchecks are being triggered
	Subcommand string
	// if RunDockerChecks is true we run docker checks
	RunDockerChecks bool
	// if RunKubernetesChecks is true we run k8s checks
	RunKubernetesChecks bool
	// if RunVersionChecks is true we run version checks
	RunVersionChecks bool
	// if RunComponentChecks is true we run component checks
	RunComponentChecks bool
	// if RunOperatorChecks is true we run operator checks
	RunOperatorChecks bool
}

type HealthChecker struct {
	Options *HealthCheckOptions

	// Things that being used while running these checks
	context *config.Context
	mctlCfg *config.MesheryCtlConfig
}

func NewHealthChecker(options *HealthCheckOptions) (*HealthChecker, error) {
	hc := &HealthChecker{
		Options: options,
	}
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return nil, err
	}
	err = mctlCfg.SetCurrentContext(tempContext)
	if err != nil {
		return nil, ErrSetCurrentContext(err)
	}
	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return nil, ErrGetCurrentContext(err)
	}
	if utils.PlatformFlag != "" {
		if utils.PlatformFlag == platformDocker || utils.PlatformFlag == platformKubernetes {
			currCtx.SetPlatform(utils.PlatformFlag)
		} else {
			return nil, ErrUnsupportedPlatform(utils.PlatformFlag, utils.CfgFile)
		}
	}

	hc.context = currCtx
	hc.mctlCfg = mctlCfg

	return hc, nil
}

var linkDocCheck = map[string]string{
	"link":    "![check-usage](/reference/images/check.png)",
	"caption": "Usage of mesheryctl system check",
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Pre-deployment and post-deployment healthchecks for Meshery",
	Long: `Verify environment pre/post-deployment of Meshery.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/check`,
	Args: cobra.MaximumNArgs(1),
	Example: `
// Run all system checks for both pre and post-deployment scenarios
mesheryctl system check

// Run pre-deployment checks (Docker and Kubernetes)
mesheryctl system check --preflight

// Run pre-deployment checks (Docker and Kubernetes)
mesheryctl system check --pre

// Run checks for all Meshery adapters
mesheryctl system check --adapters

// Run checks on a specific Meshery adapter
mesheryctl system check --adapter meshery-istio:10000
mesheryctl system check --adapter meshery-istio

// Verify the health of Meshery Operator's deployment with MeshSync and Broker
mesheryctl system check --operator
	`,
	Annotations: linkDocCheck,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &systemCheckFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		hco := &HealthCheckOptions{
			PrintLogs:  true,
			IsPreRunE:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hco)
		if err != nil {
			return ErrHealthCheckFailed(err)
		}

		runPreflightCheck := systemCheckFlags.Pre || systemCheckFlags.Preflight
		if runPreflightCheck {
			err := hc.RunPreflightHealthChecks()
			if err != nil {
				return err
			}

			prefix := "✓✓"
			suffix := "met"
			output := "\n--------------\n--------------\n%s Meshery prerequisites %s"
			if failure > 0 {
				prefix = "!!"
				suffix = "not met"
			}
			utils.Log.Info(fmt.Sprintf(output, prefix, suffix))
			return nil
		}

		if systemCheckFlags.ComponentsFlag {
			if err := hc.runComponentsHealthChecks(); err != nil {
				return err
			}
		}

		if systemCheckFlags.Adapter != "" {
			if err := hc.runAdapterHealthChecks(systemCheckFlags.Adapter); err != nil {
				return err
			}
		}
		if systemCheckFlags.AdaptersFlag {
			if err := hc.runAdapterHealthChecks(""); err != nil {
				return err
			}
		}

		currContext, err := hc.mctlCfg.GetCurrentContext()
		if err != nil {
			return ErrGetCurrentContext(err)
		}
		currPlatform := currContext.GetPlatform()

		if systemCheckFlags.OperatorsFlag {
			if platformDocker == currPlatform {
				return ErrOperatorUnsupportedPlatform(currPlatform)
			}
			if err := hc.runOperatorHealthChecks(); err != nil {
				return err
			}
		}

		skipFullCheck := systemCheckFlags.OperatorsFlag || systemCheckFlags.AdaptersFlag || systemCheckFlags.Adapter != "" || systemCheckFlags.ComponentsFlag
		if skipFullCheck {
			return nil
		}

		hc.Options.RunComponentChecks = true
		// if platform is docker only then run docker checks
		hc.Options.RunDockerChecks = currPlatform == platformDocker
		hc.Options.RunVersionChecks = true
		// if platform is kubernetes only then run kubernetes checks
		hc.Options.RunKubernetesChecks = currPlatform == platformKubernetes
		// if platform is kubernetes only then run operator checks as well since operator is only supported on kubernetes
		hc.Options.RunOperatorChecks = currPlatform == platformKubernetes
		return hc.Run()
	},
}

// Run triggers all the healthchecks according to the requirements defined from struct HealthChecks
func (hc *HealthChecker) Run() error {
	// Run meshery docker checks
	if hc.Options.RunDockerChecks {
		if err := hc.runDockerHealthChecks(); err != nil {
			return err
		}
	}
	// Run meshery kubernetes checks
	if hc.Options.RunKubernetesChecks {
		if err := hc.runKubernetesHealthChecks(); err != nil {
			return err
		}
	}
	// Run meshery component version checks
	if hc.Options.RunVersionChecks {
		if err := hc.runMesheryVersionHealthChecks(); err != nil {
			return err
		}
	}
	// Run meshery component health checks
	if hc.Options.RunComponentChecks {
		if err := hc.runComponentsHealthChecks(); err != nil {
			return err
		}
	}

	if hc.Options.RunOperatorChecks {
		if err := hc.runOperatorHealthChecks(); err != nil {
			return err
		}
	}

	return nil
}

// Run preflight healthchecks to verify environment health
func (hc *HealthChecker) RunPreflightHealthChecks() error {
	// Docker healthchecks are only invoked when it's not a PreRunExecution
	// or it's a PreRunExecution and current platform is docker
	if !hc.Options.IsPreRunE || (hc.Options.IsPreRunE && hc.context.Platform == platformDocker) {
		//Run docker healthchecks
		if err := hc.runDockerHealthChecks(); err != nil {
			return err
		}
	}
	// Kubernetes healthchecks are only
	// invoked when it's not a PreRunExecution
	// or it's a PreRunExecution and current platform is kubernetes
	if !hc.Options.IsPreRunE || (hc.Options.IsPreRunE && hc.context.Platform == platformKubernetes) {
		if err := hc.runKubernetesHealthChecks(); err != nil {
			return err
		}
	}

	return nil
}

func wrapDockerContextError(err error, hc *HealthChecker) error {
	return ErrDockerContext(errors.Wrapf(err, "Meshery is not running locally, please ensure that the appropriate Docker context is selected for Meshery endpoint: %s. To list all configured contexts use `docker context ls`", hc.context.GetEndpoint()))
}

func wrapDockerNotRunningError(err error, hc *HealthChecker) error {
	return ErrDockerNotRunning(errors.Wrapf(err, "Please start Docker. Run `mesheryctl system %s` once Docker is started ", hc.Options.Subcommand))
}

func wrapDockerStartError(err error) error {
	return ErrDockerStart(errors.Wrapf(err, "failed to start Docker "))
}

// Run healthchecks to verify if docker is running and active
func (hc *HealthChecker) runDockerHealthChecks() error {
	if hc.Options.PrintLogs {
		utils.Log.Info("\nDocker \n--------------")
	}
	endpointParts := strings.Split(hc.context.GetEndpoint(), ":")

	// Check whether docker daemon is running using Docker client API
	dockerCli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		if endpointParts[1] != "//localhost" {
			return wrapDockerContextError(err, hc)
		}
		if hc.Options.IsPreRunE { // if this is PreRunExec we trigger self installation
			utils.Log.Warn(fmt.Errorf("!! Docker is not running"))
			//If preRunExecution and the current platform is docker then we trigger docker installation
			//No auto installation of docker for windows
			if runtime.GOOS == "windows" {
				return wrapDockerNotRunningError(err, hc)
			}
			err = utils.Startdockerdaemon(hc.Options.Subcommand)
			if err != nil {
				return wrapDockerStartError(err)
			}
		} else if hc.Options.PrintLogs { // warn incase of printing logs
			utils.Log.Warn(fmt.Errorf("!! Docker is not running"))
		} else { // else we're supposed to grab errors
			return ErrDockerUnknown(err)
		}
		if hc.context.Platform == platformDocker {
			failure++
		}
	} else {
		// Try to ping the Docker daemon to verify it's actually running
		_, err = dockerCli.Ping(context.Background())
		if err != nil {
			if endpointParts[1] != "//localhost" {
				return wrapDockerContextError(err, hc)
			}
			if hc.Options.IsPreRunE { // if this is PreRunExec we trigger self installation
				utils.Log.Warn(fmt.Errorf("!! Docker is not running"))
				if runtime.GOOS == "windows" {
					return wrapDockerNotRunningError(err, hc)
				}
				err = utils.Startdockerdaemon(hc.Options.Subcommand)
				if err != nil {
					return wrapDockerStartError(err)
				}
			} else if hc.Options.PrintLogs { // warn incase of printing logs
				utils.Log.Warn(fmt.Errorf("!! Docker is not running"))
			} else { // else we're supposed to grab errors
				return ErrDockerUnknown(err)
			}
			if hc.context.Platform == platformDocker {
				failure++
			}
		} else { // if not error we check if we are supposed to print logs
			// logging if we're supposed to
			if hc.Options.PrintLogs {
				utils.Log.Info("✓ Docker is running")
			}
		}
	}

	// Since we now use docker compose sdk, we don't need to check for docker-compose  binary
	if hc.Options.PrintLogs {
		utils.Log.Info("✓ docker-compose is available")
	}

	return nil
}

// Run healthchecks to verify if kubernetes client can be initialized and can be queried
func (hc *HealthChecker) runKubernetesAPIHealthCheck() error {
	if hc.Options.PrintLogs {
		utils.Log.Info("\nKubernetes API \n--------------")
	}
	//Check whether k8s client can be initialized
	client, err := meshkitkube.New([]byte(""))
	if err != nil {
		if hc.context.Platform == platformKubernetes { // increase failure count
			failure++
		}
		if hc.Options.PrintLogs { // print logs if we're supposed to
			utils.Log.Warn(fmt.Errorf("!! cannot initialize Kubernetes client"))
			utils.Log.Warn(fmt.Errorf("!! cannot query the Kubernetes API"))
			return nil
		}
		// else we're supposed to grab the error
		errMsg := fmt.Errorf("%s. Your %s context is configured to run Meshery on Kubernetes using the %s token",
			err.Error(), hc.mctlCfg.CurrentContext, hc.context.Token)
		return ErrK8sConfig(errMsg)
	}

	if hc.Options.PrintLogs { // print logs if we're supposed to
		utils.Log.Info("✓ can initialize Kubernetes client")
	}

	//Check whether kubernetes can be queried
	podInterface := client.KubeClient.CoreV1().Pods("")
	_, err = podInterface.List(context.TODO(), v1.ListOptions{})
	if err != nil {
		if hc.context.Platform == platformKubernetes { // increase failure count
			failure++
		}
		if hc.Options.PrintLogs { // log incase we're supposed to
			utils.Log.Warn(fmt.Errorf("!! cannot query the Kubernetes API"))
			return nil
		}
		return ErrK8SQuery(err)
	}

	if hc.Options.PrintLogs { // log incase we're supposed to
		utils.Log.Info("✓ can query the Kubernetes API")
	}

	return nil
}

// Run healthchecks to verify kubectl and kubenetes version with
// minimum compatible versions
func (hc *HealthChecker) runKubernetesVersionHealthCheck() error {
	if hc.Options.PrintLogs {
		utils.Log.Info("\nKubernetes Version \n--------------")
	}
	//Check whether system has minimum supported versions of kubernetes and kubectl
	kubeVersion, err := utils.GetK8sVersionInfo()
	if err != nil {
		if hc.context.Platform == platformKubernetes { // increase failure count
			failure++
		}
		// probably kubernetes isn't running
		if hc.Options.PrintLogs { // log if we're supposed to
			utils.Log.Warn(fmt.Errorf("!! cannot check Kubernetes version"))
		} else { // else we're supposed to catch the error
			return err
		}
	} else {
		// kubernetes is running so check the version
		err = utils.CheckK8sVersion(kubeVersion)
		if err != nil {
			if hc.context.Platform == platformKubernetes { // increase failure count
				failure++
			}
			if hc.Options.PrintLogs { // log if we're supposed to
				utils.Log.Warnf("!! %s", err)
			} else { // else we gotta catch the error
				return err
			}
		} else { // if not error we check if we are supposed to print logs
			if hc.Options.PrintLogs { // log if we're supposed to
				utils.Log.Info("✓ running the minimum Kubernetes version")
			}
		}
	}

	err = utils.CheckKubectlVersion()
	if err != nil {
		if hc.context.Platform == platformKubernetes { // increase failure count
			failure++
		}
		if hc.Options.PrintLogs { // log if we're supposed to
			utils.Log.Warnf("!! %s", err)
		} else { // else we gotta catch the error
			return err
		}
	} else { // if not error we check if we are supposed to print logs
		if hc.Options.PrintLogs { // log if we're supposed to
			utils.Log.Info("✓ running the minimum kubectl version")
		}
	}

	return nil
}

// runKubernetesHealthChecks runs checks regarding k8s api and k8s plus kubectl version
func (hc *HealthChecker) runKubernetesHealthChecks() error {
	// Run k8s API healthchecks
	if err := hc.runKubernetesAPIHealthCheck(); err != nil {
		return err
	}
	// Run k8s plus kubectl minimum version healthchecks
	err = hc.runKubernetesVersionHealthCheck()
	return err
}

// runMesheryVersionHealthChecks runs checks regarding meshery version and mesheryctl version
func (hc *HealthChecker) runMesheryVersionHealthChecks() error {
	skipServerLogs := false

	if hc.Options.PrintLogs {
		utils.Log.Info("\nMeshery Version \n--------------")
	}

	url := hc.mctlCfg.GetBaseMesheryURL()
	var serverVersion *config.Version
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/system/version", url), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	// failed to fetch response for server version
	if err != nil || resp.StatusCode != 200 {
		if hc.Options.PrintLogs { // log if we're supposed to
			utils.Log.Info("!! failed to check Meshery Server version. try starting Meshery with `mesheryctl system start`")
			skipServerLogs = true
		} else {
			return err
		}
	}

	// skip this part as we failed to get a response from the api
	if !skipServerLogs {
		// needs multiple defer as Body.Close needs a valid response
		defer func() { _ = resp.Body.Close() }()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Errorf("\n  Invalid response: %v", err)
		}

		err = json.Unmarshal(data, &serverVersion)
		if err != nil {
			return errors.Errorf("\n  Unable to unmarshal data: %v", err)
		}

		isOutdated, _, err := handlers.CheckLatestVersion(serverVersion.GetBuild())
		if err != nil {
			return err
		}
		if hc.Options.PrintLogs { // log if we're supposed to
			if !*isOutdated {
				utils.Log.Infof("✓ Meshery Server is up-to-date (stable-%s)", serverVersion.GetBuild())
			} else {
				utils.Log.Info("!! Meshery Server is not up-to-date")
			}
		} else { // else we grab the error
			if !*isOutdated {
				return errors.New("!! Meshery Server is not up-to-date")
			}
		}
	}

	latestVersions, err := meshkitutils.GetLatestReleaseTagsSorted(c.GetMesheryGitHubOrg(), c.GetMesheryGitHubRepo())
	if err != nil {
		if hc.Options.PrintLogs { // log if we're supposed to
			utils.Log.Info("!! failed to fetch latest release tag of mesheryctl")
			// skip further for client as we failed to check latest tag on github
			return nil
		}
		return err
	}
	if len(latestVersions) == 0 {
		return fmt.Errorf("no versions found")
	}
	latest := latestVersions[len(latestVersions)-1]
	version := constants.GetMesheryctlVersion()
	if hc.Options.PrintLogs { // log if we're supposed to
		if latest == version {
			utils.Log.Infof("✓ CLI is up-to-date (stable-%s)", version)
		} else {
			utils.Log.Infof("!! CLI is not up-to-date (stable-%s)", version)
		}
	} else { // else we grab the error
		if latest != version {
			return ErrSystemCheckInvalidCliVersion(version)
		}
	}

	return nil
}

// runComponentsHealthChecks runs health checks for Adapters, Operator, and all deployments in meshery ecosystem
func (hc *HealthChecker) runComponentsHealthChecks() error {
	if hc.Options.PrintLogs {
		utils.Log.Info("\nMeshery Components \n--------------")
	}
	return hc.runAdapterHealthChecks("")
}

// runOperatorHealthChecks executes health-checks for Operators
func (hc *HealthChecker) runOperatorHealthChecks() error {
	if hc.Options.PrintLogs {
		utils.Log.Info("\nMeshery Operators \n--------------")
	}
	clientMesh, err := meshkitkube.New([]byte(""))
	if err != nil {
		return err
	}

	// List the pods in the `meshery` Namespace
	podList, err := utils.GetPodList(clientMesh, utils.MesheryNamespace)
	if err != nil {
		return err
	}

	operatorCheck := false
	brokerCheck := false
	meshsyncCheck := false

	// Traverse the list of pods and check if the operator, broker and meshsync are running
	// This loop is made because if the clusterIP is not public, then it can't be accessed from outside the cluster
	for _, pod := range podList.Items {
		name := utils.GetCleanPodName(pod.GetName())

		if name == "meshery-operator" {
			operatorCheck = true
		}

		if name == "meshery-broker" {
			brokerCheck = true
		}

		if name == "meshery-meshsync" {
			meshsyncCheck = true
		}
	}

	if !operatorCheck {
		return errors.New("!! Meshery Operator is not running")
	} else {
		utils.Log.Info("✓ Meshery Operator is running")
	}

	if !meshsyncCheck {
		utils.Log.Info("!! Meshsync is not running")
	} else {
		utils.Log.Info("✓ Meshsync is running")
	}

	if !brokerCheck {
		utils.Log.Info("!! Meshery Broker is not running")
	} else {
		utils.Log.Info("✓ Meshery Broker is running")

		// Check if broker CR contains Status.Endpoint object with External and Internal parts
		const (
			brokerGroup    = "meshery.io"
			brokerVersion  = "v1alpha1"
			brokerResource = "brokers"
			brokerName     = "meshery-broker"
		)
		brokerGVR := schema.GroupVersionResource{
			Group:    brokerGroup,
			Version:  brokerVersion,
			Resource: brokerResource,
		}

		brokerCR, err := clientMesh.DynamicKubeClient.Resource(brokerGVR).Namespace(utils.MesheryNamespace).Get(context.Background(), brokerName, v1.GetOptions{})
		if err != nil {
			utils.Log.Info("!! Could not retrieve Meshery Broker CR")
		} else {
			// Check if status.endpoint exists with external and internal parts
			status, found, err := unstructured.NestedMap(brokerCR.Object, "status")
			if err != nil {
				utils.Log.Infof("!! Error parsing Meshery Broker CR status: %v", err)
			} else if !found {
				utils.Log.Info("!! Meshery Broker CR does not contain Status section")
			} else {
				endpoint, endpointFound, err := unstructured.NestedMap(status, "endpoint")
				if err != nil {
					utils.Log.Infof("!! Error parsing Meshery Broker CR status.endpoint: %v", err)
				} else if !endpointFound {
					utils.Log.Info("!! Meshery Broker CR does not contain Status.Endpoint")
				} else {
					external, externalFound, _ := unstructured.NestedString(endpoint, "external")
					internal, internalFound, _ := unstructured.NestedString(endpoint, "internal")

					if externalFound && internalFound && external != "" && internal != "" {
						utils.Log.Infof("✓ Meshery Broker CR contains Status.Endpoint (External: %s, Internal: %s)", external, internal)
					} else {
						missingParts := []string{}
						if !externalFound || external == "" {
							missingParts = append(missingParts, "External")
						}
						if !internalFound || internal == "" {
							missingParts = append(missingParts, "Internal")
						}
						utils.Log.Infof("!! Meshery Broker CR Status.Endpoint missing: %s", strings.Join(missingParts, ", "))
					}
				}
			}
		}
	}

	return nil
}

// runAdapterHealthChecks executes health-checks for Adapters
// If no adapter is specified all the adapters are checked
func (hc *HealthChecker) runAdapterHealthChecks(adapterName string) error {
	url := hc.mctlCfg.GetBaseMesheryURL()
	client := &http.Client{}
	var adapters []*models.Adapter
	prefs, err := utils.GetSessionData(hc.mctlCfg)
	if err != nil {
		return err
	}
	for _, adapter := range prefs.MeshAdapters {
		if adapterName != "" {
			name := strings.Split(adapter.Location, ":")[0]
			if adapterName == name || adapterName == adapter.Location {
				adapters = append(adapters, adapter)
				break
			}
		} else {
			adapters = append(adapters, adapter)
		}
	}
	if len(adapters) == 0 {
		utils.Log.Info("- No components configured in current context")
		return nil
	}
	for _, adapter := range adapters {
		name := adapter.Location
		skipAdapter := false
		req, err := utils.NewRequest("GET", fmt.Sprintf("%s/api/system/adapters?adapter=%s", url, name), nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			if hc.Options.PrintLogs { // incase we're printing logs
				utils.Log.Infof("!! failed to connect to Meshery Adapter for %s ", name)
				skipAdapter = true
			} else { // or we're supposed to grab the errors
				return utils.ErrFailedToConnectAdapter(name, err)
			}
			continue
		}
		if !skipAdapter {
			// needs multiple defer as Body.Close needs a valid response
			defer func() { _ = resp.Body.Close() }()
			if resp.StatusCode != 200 {
				if hc.Options.PrintLogs { // incase we're printing logs
					utils.Log.Infof("!! Meshery Adapter for %s is running but not reachable", name)
				} else { // or we're supposed to grab the errors
					return utils.ErrAdapterNotReachable(name)
				}
			} else { // if status == 200 we check if we are supposed to print logs
				if hc.Options.PrintLogs { // incase we're printing logs
					utils.Log.Infof("✓ %s adapter is running and reachable", name)
				}
			}
		}
	}
	return nil
}

// mesheryReadinessHealthCheck is waiting for Meshery to start, returns (ready, error)
func mesheryReadinessHealthCheck() (bool, error) {
	kubeClient, err := meshkitkube.New([]byte(""))
	if err != nil {
		return false, err
	}
	if err := utils.WaitForPodRunning(kubeClient, "meshery", utils.MesheryNamespace, 300); err != nil {
		return false, err
	}

	return true, nil
}

func init() {
	checkCmd.Flags().BoolVarP(&systemCheckFlags.Preflight, "preflight", "", false, "Verify environment readiness to deploy Meshery")
	checkCmd.Flags().BoolVarP(&systemCheckFlags.Pre, "pre", "", false, "Verify environment readiness to deploy Meshery")
	checkCmd.Flags().BoolVarP(&systemCheckFlags.ComponentsFlag, "components", "", false, "Check status of Meshery components")
	checkCmd.Flags().BoolVarP(&systemCheckFlags.AdaptersFlag, "adapters", "", false, "Check status of meshery adapters")
	checkCmd.Flags().StringVarP(&systemCheckFlags.Adapter, "adapter", "", "", "Check status of specified meshery adapter")
	checkCmd.Flags().BoolVarP(&systemCheckFlags.OperatorsFlag, "operator", "", false, "Verify the health of Meshery Operator's deployment with MeshSync and Broker")
}
