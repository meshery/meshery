package system

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/models"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sVersion "k8s.io/apimachinery/pkg/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	preflight      bool
	pre            bool
	componentsFlag bool
	failure        int
)

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
		return nil, errors.Wrap(err, "error processing config")
	}
	err = mctlCfg.SetCurrentContext(tempContext)
	if err != nil {
		return nil, err
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return nil, err
	}

	if utils.PlatformFlag != "" {
		if utils.PlatformFlag == "docker" || utils.PlatformFlag == "kubernetes" {
			currCtx.SetPlatform(utils.PlatformFlag)
		} else {
			return nil, ErrUnsupportedPlatform(utils.PlatformFlag, utils.CfgFile)
		}
	}

	hc.context = currCtx
	hc.mctlCfg = mctlCfg

	return hc, nil
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Meshery environment check",
	Long:  `Verify environment pre/post-deployment of Meshery.`,
	Args:  cobra.MaximumNArgs(1),
	Example: `
// Run system checks for both pre and post mesh deployment scenarios on Meshery
mesheryctl system check

// Run Pre-mesh deployment checks (Docker and Kubernetes)
mesheryctl system check --preflight

// Run checks on specific mesh adapter
mesheryctl system check --adapter meshery-istio:10000

// Verify the health of Meshery Operator's deployment with MeshSync and Broker
mesheryctl system check --operator

// Runs diagnostic checks and bundles up to open an issue if present
mesheryctl system check --report

! Refer below image link for usage
* Usage of mesheryctl system check
# ![check-usage](/assets/img/mesheryctl/check.jpeg)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		hco := &HealthCheckOptions{
			PrintLogs:  true,
			IsPreRunE:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hco)
		if err != nil {
			return errors.New("failed to initialize a healthchecker")
		}

		// if --pre or --preflight has been passed we run preflight checks
		if pre || preflight {
			// Run preflight checks
			err := hc.RunPreflightHealthChecks()
			if err != nil {
				return err
			}
			// Print End
			if failure == 0 {
				log.Info("\n--------------\n--------------\n✓✓ Meshery prerequisites met")
			} else {
				log.Info("\n--------------\n--------------\n!! Meshery prerequisites not met")
			}
			return nil
		} else if componentsFlag { // if --components has been passed we run checks related to components
			return hc.runComponentsHealthChecks()
		}

		// if no flags passed we run complete system check
		hc.Options.RunComponentChecks = true
		hc.Options.RunDockerChecks = true
		hc.Options.RunKubernetesChecks = true
		hc.Options.RunVersionChecks = true
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

	return nil
}

// Run preflight healthchecks to verify environment health
func (hc *HealthChecker) RunPreflightHealthChecks() error {
	// Docker healthchecks are only invoked when it's not a PreRunExecution
	// or it's a PreRunExecution and current platform is docker
	if !hc.Options.IsPreRunE || (hc.Options.IsPreRunE && hc.context.Platform == "docker") {
		//Run docker healthchecks
		if err := hc.runDockerHealthChecks(); err != nil {
			return err
		}
	}
	// Kubernetes healthchecks are only
	// invoked when it's not a PreRunExecution
	// or it's a PreRunExecution and current platform is kubernetes
	if !hc.Options.IsPreRunE || (hc.Options.IsPreRunE && hc.context.Platform == "kubernetes") {
		if err := hc.runKubernetesHealthChecks(); err != nil {
			return err
		}
	}

	return nil
}

// Run healthchecks to verify if docker is running and active
func (hc *HealthChecker) runDockerHealthChecks() error {
	if hc.Options.PrintLogs {
		log.Info("\nDocker \n--------------")
	}
	//Check whether docker daemon is running or not
	err := exec.Command("docker", "ps").Run()
	if err != nil {
		if hc.Options.IsPreRunE { // if this is PreRunExec we trigger self installation
			log.Warn("!! Docker is not running")
			//If preRunExecution and the current platform is docker then we trigger docker installation
			//No auto installation of docker for windows
			if runtime.GOOS == "windows" {
				return errors.Wrapf(err, "Please start Docker. Run `mesheryctl system %s` once Docker is started ", hc.Options.Subcommand)
			}
			err = utils.Startdockerdaemon(hc.Options.Subcommand)
			if err != nil {
				return errors.Wrapf(err, "failed to start Docker ")
			}
		} else if hc.Options.PrintLogs { // warn incase of printing logs
			log.Warn("!! Docker is not running")
		} else { // else we're supposed to grab errors
			return err
		}

		if hc.context.Platform == "docker" {
			failure++
		}
	} else { // if not error we check if we are supposed to print logs
		// logging if we're supposed to
		if hc.Options.PrintLogs {
			log.Info("✓ Docker is running")
		}
	}

	//Check for installed docker-compose on client system
	err = exec.Command("docker-compose", "-v").Run()
	if err != nil {
		if hc.Options.IsPreRunE { // if PreRunExec we trigger self installation
			log.Warn("!! docker-compose is not available")
			//No auto installation of Docker-compose for windows
			if runtime.GOOS == "windows" {
				return errors.Wrapf(err, "please install docker-compose. Run `mesheryctl system %s` after docker-compose is installed ", hc.Options.Subcommand)
			}
			err = utils.InstallprereqDocker()
			if err != nil {
				return errors.Wrapf(err, "failed to install prerequisites. Run `mesheryctl system %s` after docker-compose is installed ", hc.Options.Subcommand)
			}
		} else if hc.Options.PrintLogs { // warn incase of printing logs
			log.Warn("!! docker-compose is not available")
		} else { // else we're supposed to grab the error
			return err
		}

		if hc.context.Platform == "docker" {
			failure++
		}
	} else { // if not error we check if we are supposed to print logs
		// logging if we're supposed to
		if hc.Options.PrintLogs {
			log.Info("✓ docker-compose is available")
		}
	}

	return nil
}

// Run healthchecks to verify if kubernetes client can be initialized and can be queried
func (hc *HealthChecker) runKubernetesAPIHealthCheck() error {
	if hc.Options.PrintLogs {
		log.Info("\nKubernetes API \n--------------")
	}
	//Check whether k8s client can be initialized
	client, err := meshkitkube.New([]byte(""))
	if err != nil {
		if hc.context.Platform == "kubernetes" { // increase failure count
			failure++
		}
		if hc.Options.PrintLogs { // print logs if we're supposed to
			log.Warn("!! cannot initialize Kubernetes client")
			log.Warn("!! cannot query the Kubernetes API")
			return nil
		}
		// else we're supposed to grab the error
		errMsg := fmt.Errorf("%s. Your %s context is configured to run Meshery on Kubernetes using the %s token",
			err.Error(), hc.mctlCfg.CurrentContext, hc.context.Token)
		return ErrK8sConfig(errMsg)
	}

	if hc.Options.PrintLogs { // print logs if we're supposed to
		log.Info("✓ can initialize Kubernetes client")
	}

	//Check whether kubernetes can be queried
	podInterface := client.KubeClient.CoreV1().Pods("")
	_, err = podInterface.List(context.TODO(), v1.ListOptions{})
	if err != nil {
		if hc.context.Platform == "kubernetes" { // increase failure count
			failure++
		}
		if hc.Options.PrintLogs { // log incase we're supposed to
			log.Warn("!! cannot query the Kubernetes API")
			return nil
		}
		return ErrK8SQuery(err)
	}

	if hc.Options.PrintLogs { // log incase we're supposed to
		log.Info("✓ can query the Kubernetes API")
	}

	return nil
}

// Run healthchecks to verify kubectl and kubenetes version with
// minimum compatible versions
func (hc *HealthChecker) runKubernetesVersionHealthCheck() error {
	if hc.Options.PrintLogs {
		log.Info("\nKubernetes Version \n--------------")
	}
	//Check whether system has minimum supported versions of kubernetes and kubectl
	var kubeVersion *k8sVersion.Info
	kubeVersion, err := utils.GetK8sVersionInfo()
	if err != nil {
		if hc.context.Platform == "kubernetes" { // increase failure count
			failure++
		}
		// probably kubernetes isn't running
		if hc.Options.PrintLogs { // log if we're supposed to
			log.Warn("!! cannot check Kubernetes version")
		} else { // else we're supposed to catch the error
			return err
		}
	} else {
		// kubernetes is running so check the version
		err = utils.CheckK8sVersion(kubeVersion)
		if err != nil {
			if hc.context.Platform == "kubernetes" { // increase failure count
				failure++
			}
			if hc.Options.PrintLogs { // log if we're supposed to
				log.Warnf("!! %s", err)
			} else { // else we gotta catch the error
				return err
			}
		} else { // if not error we check if we are supposed to print logs
			if hc.Options.PrintLogs { // log if we're supposed to
				log.Info("✓ running the minimum Kubernetes version")
			}
		}
	}

	err = utils.CheckKubectlVersion()
	if err != nil {
		if hc.context.Platform == "kubernetes" { // increase failure count
			failure++
		}
		if hc.Options.PrintLogs { // log if we're supposed to
			log.Warnf("!! %s", err)
		} else { // else we gotta catch the error
			return err
		}
	} else { // if not error we check if we are supposed to print logs
		if hc.Options.PrintLogs { // log if we're supposed to
			log.Info("✓ running the minimum kubectl version")
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
		log.Info("\nMeshery Version \n--------------")
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
			log.Info("!! failed to check Meshery Server version. try starting Meshery with `mesheryctl system start`")
			skipServerLogs = true
		} else {
			return err
		}
	}

	// skip this part as we failed to get a response from the api
	if !skipServerLogs {
		// needs multiple defer as Body.Close needs a valid response
		defer resp.Body.Close()
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
				log.Infof("✓ Meshery Server is up-to-date (stable-%s)", serverVersion.GetBuild())
			} else {
				log.Info("!! Meshery Server is not up-to-date")
			}
		} else { // else we grab the error
			if !*isOutdated {
				return errors.New("!! Meshery Server is not up-to-date")
			}
		}
	}

	latest, err := utils.GetLatestStableReleaseTag()
	if err != nil {
		if hc.Options.PrintLogs { // log if we're supposed to
			log.Info("!! failed to fetch latest release tag of mesheryctl")
			// skip further for client as we failed to check latest tag on github
			return nil
		}
		return err
	}

	version := constants.GetMesheryctlVersion()
	if hc.Options.PrintLogs { // log if we're supposed to
		if latest == version {
			log.Infof("✓ CLI is up-to-date (stable-%s)", version)
		} else {
			log.Infof("!! CLI is not up-to-date (stable-%s)", version)
		}
	} else { // else we grab the error
		if latest != version {
			return errors.New("!! CLI is not up-to-date")
		}
	}

	return nil
}

// runComponentsHealthChecks runs health checks for Adapters, Operator, and all deployments in meshery ecosystem
func (hc *HealthChecker) runComponentsHealthChecks() error {
	if hc.Options.PrintLogs {
		log.Info("\nMeshery Components \n--------------")
	}
	return hc.runAdapterHealthChecks()
}

// runAdapterHealthChecks executes health-checks for Adapters
func (hc *HealthChecker) runAdapterHealthChecks() error {
	url := hc.mctlCfg.GetBaseMesheryURL()
	client := &http.Client{}

	// Request to grab running adapters and ports
	req, err := utils.NewRequest("GET", fmt.Sprintf("%s/api/system/adapters", url), nil)
	if err != nil {
		if hc.Options.PrintLogs {
			log.Info("!! Authentication token not found. Please supply a valid user token. Login with `mesheryctl system login`")
			// skip further as we failed to attach token
			return nil
		}
		return errors.New("Authentication token not found. Please supply a valid user token. Login with `mesheryctl system login`")
	}

	var adapters []*models.Adapter
	resp, err := client.Do(req)
	// failed to grab response from the request
	if err != nil || resp.StatusCode != 200 {
		if hc.Options.PrintLogs {
			log.Info("!! Failed to connect to Meshery Adapters")
			// skip further as we failed to grab the response from server
			return nil
		}
		return err
	}

	// needs multiple defer as Body.Close needs a valid response
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Errorf("\n  Invalid response: %v", err)
	}

	err = json.Unmarshal(data, &adapters)
	if err != nil {
		return errors.Errorf("\n  Unable to unmarshal data: %v", err)
	}

	// check for each adapter
	for _, adapter := range adapters {
		skipAdapter := false

		name := strings.Split(adapter.Location, ":")[0]
		if adapter.Ops == nil {
			if hc.Options.PrintLogs { // incase we're printing logs
				log.Infof("!! %s adapter is not running", name)
			}
			continue
		}
		req, err := utils.NewRequest("GET", fmt.Sprintf("%s/api/system/adapters?adapter=%s", url, adapter.Name), nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err != nil {
			if hc.Options.PrintLogs { // incase we're printing logs
				log.Infof("!! failed to connect to Meshery Adapter for %s ", name)
				skipAdapter = true
			} else { // or we're supposed to grab the errors
				return fmt.Errorf("!! failed to connect to Meshery Adapter for%s adapter: %s", name, err)
			}
			continue
		}

		// skip the adapter as we failed to receive response for adapter
		if !skipAdapter {
			// needs multiple defer as Body.Close needs a valid response
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				if hc.Options.PrintLogs { // incase we're printing logs
					log.Infof("!! Meshery Adapter for %s is running but not reachable", name)
				} else { // or we're supposed to grab the errors
					return fmt.Errorf("!! Meshery Adapter for %s is running, but not reachable", name)
				}
			} else { // if status == 200 we check if we are supposed to print logs
				if hc.Options.PrintLogs { // incase we're printing logs
					log.Infof("✓ %s adapter is running and reachable", name)
				}
			}
		}
	}

	return nil
}

func (hc *HealthChecker) runOperatorHealthChecks() error {
	//TODO
	return nil
}

func (hc *HealthChecker) runMesheryReadinessHealthChecks() error {
	ready, err := mesheryReadinessHealthCheck()
	if err != nil || !ready {
		if hc.Options.PrintLogs { // incase we're printing logs
			log.Infof("!! Meshery failed to reach Running state")
		} else { // or we're supposed to grab the errors
			return fmt.Errorf("!! Meshery failed to reach Running state. %s", err)
		}
	}
	if hc.Options.PrintLogs { // incase we're printing logs
		log.Infof("✓ Meshery is in Running state")
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
	checkCmd.Flags().BoolVarP(&preflight, "preflight", "", false, "Verify environment readiness to deploy Meshery")
	checkCmd.Flags().BoolVarP(&pre, "pre", "", false, "Verify environment readiness to deploy Meshery")
	checkCmd.Flags().BoolVarP(&componentsFlag, "components", "", false, "Check status of Meshery components")
}
