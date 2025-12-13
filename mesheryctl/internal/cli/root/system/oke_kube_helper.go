package system

import (
    "fmt"
    "os"
    "path"
    "strings"

    "github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
    "github.com/meshery/meshery/mesheryctl/pkg/utils"
    meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
    "github.com/pkg/errors"
    "helm.sh/helm/v3/pkg/action"
    "helm.sh/helm/v3/pkg/chart/loader"
    helmcliconfig "helm.sh/helm/v3/pkg/cli"
    clientcmd "k8s.io/client-go/tools/clientcmd"
    clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
    "k8s.io/cli-runtime/pkg/genericclioptions"
    log "github.com/sirupsen/logrus"
)

// GetKubeconfigPath returns the default kubeconfig path used by mesheryctl.
func GetKubeconfigPath() string {
    return path.Join(utils.MesheryFolder, "kubeconfig.yaml")
}

// IsOKEFromKubeconfig loads kubeconfig and returns true if any user has an exec.Command referencing OCI/OKE.
func IsOKEFromKubeconfig(kubeconfigPath string) (bool, error) {
    cfg, err := clientcmd.LoadFromFile(kubeconfigPath)
    if err != nil {
        return false, errors.Wrapf(err, "failed to load kubeconfig %s", kubeconfigPath)
    }
    return isOKEFromClientcmdConfig(cfg), nil
}

func isOKEFromClientcmdConfig(cfg *clientcmdapi.Config) bool {
    for _, auth := range cfg.AuthInfos {
        if auth == nil || auth.Exec == nil {
            continue
        }
        cmd := strings.ToLower(auth.Exec.Command)
        if strings.Contains(cmd, "oci") || strings.Contains(cmd, "oracle") || strings.Contains(cmd, "oci-cli") {
            return true
        }
        for _, a := range auth.Exec.Args {
            if strings.Contains(strings.ToLower(a), "oci") || strings.Contains(strings.ToLower(a), "oracle") {
                return true
            }
        }
    }
    return false
}

// ApplyHelmChartsSmart chooses direct-kubeconfig or MeshKit ApplyHelmChart based on kubeconfig detection.
// It defaults to the kubeconfig path used by mesheryctl.
func ApplyHelmChartsSmart(kubeClient *meshkitkube.Client, currCtx *config.Context, mesheryImageVersion string, act meshkitkube.HelmChartAction, callbackURL, providerURL string) error {
    kubeconfigPath := GetKubeconfigPath()

    isOKE := false
    if ok, err := IsOKEFromKubeconfig(kubeconfigPath); err == nil && ok {
        isOKE = true
    } else if err != nil {
        // log and fall back to MeshKit path; do not fail detection here
        log.Debugf("Unable to detect OKE from kubeconfig: %v", err)
    }

    if isOKE {
        log.Debug("OKE kubeconfig detected; using direct kubeconfig Helm apply")
        return applyHelmChartsWithKubeconfig(kubeconfigPath, currCtx, mesheryImageVersion, act, callbackURL, providerURL)
    }

    log.Debug("Using MeshKit ApplyHelmChart for non-OKE cluster")
    return applyHelmCharts(kubeClient, currCtx, mesheryImageVersion, act, callbackURL, providerURL)
}

// Apply Meshery helm charts using kubeconfig file directly (required for OKE exec auth)
func applyHelmChartsWithKubeconfig(kubeconfigPath string, currCtx *config.Context, mesheryImageVersion string, act meshkitkube.HelmChartAction, callbackURL, providerURL string) error {
	// Verify kubeconfig exists and is readable
	if _, err := os.Stat(kubeconfigPath); err != nil {
		return errors.Wrapf(err, "kubeconfig file not found at %s", kubeconfigPath)
	}

	// Get value overrides
	overrideValues := utils.SetOverrideValues(currCtx, mesheryImageVersion, callbackURL, providerURL)
	chartVersion := ""
	if mesheryImageVersion != "latest" {
		chartVersion = mesheryImageVersion
	}

	action := "install"
	if act == meshkitkube.UNINSTALL {
		action = "uninstall"
	}

	// Apply Meshery Server chart
	if err := applyHelmChartDirect(kubeconfigPath, "meshery", utils.HelmChartURL, utils.HelmChartName, chartVersion, utils.MesheryNamespace, overrideValues, act); err != nil {
		return fmt.Errorf("could not %s Meshery Server: %w", action, err)
	}

	// Apply Meshery Operator chart
	if err := applyHelmChartDirect(kubeconfigPath, "meshery-operator", utils.HelmChartURL, utils.HelmChartOperatorName, chartVersion, utils.MesheryNamespace, nil, act); err != nil {
		return fmt.Errorf("could not %s meshery-operator: %w", action, err)
	}

	return nil
}

// applyHelmChartDirect uses Helm libraries directly with kubeconfig file
func applyHelmChartDirect(kubeconfigPath, releaseName, repoURL, chartName, version, namespace string, values map[string]interface{}, actionType meshkitkube.HelmChartAction) error {
	// Set KUBECONFIG environment variable for Helm
	oldKubeconfig := os.Getenv("KUBECONFIG")
	if err := os.Setenv("KUBECONFIG", kubeconfigPath); err != nil {
		return errors.Wrap(err, "failed to set KUBECONFIG")
	}
	defer func() {
		if oldKubeconfig != "" {
			os.Setenv("KUBECONFIG", oldKubeconfig)
		} else {
			os.Unsetenv("KUBECONFIG")
		}
	}()

	// Create Helm settings
	settings := helmcliconfig.New()
	chartPathOptions := action.ChartPathOptions{
		RepoURL: repoURL,
		Version: version,
	}

	// Locate chart
	chartPath, err := chartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		return errors.Wrapf(err, "failed to locate chart %s", chartName)
	}

	// Load chart
	ch, err := loader.Load(chartPath)
	if err != nil {
		return errors.Wrapf(err, "failed to load chart from %s", chartPath)
	}

	// Create Helm action configuration
	actionConfig := new(action.Configuration)
	flags := genericclioptions.NewConfigFlags(false)
	flags.KubeConfig = &kubeconfigPath
	ns := namespace
	if ns == "" {
		ns = "default"
	}

	if err := actionConfig.Init(flags, ns, "secret", log.Debugf); err != nil {
		return errors.Wrap(err, "failed to initialize helm action")
	}

	// Execute action
	switch actionType {
	case meshkitkube.INSTALL:
		inst := action.NewInstall(actionConfig)
		inst.Namespace = ns
		inst.ReleaseName = releaseName
		inst.CreateNamespace = true
		if _, err := inst.Run(ch, values); err != nil {
			return errors.Wrapf(err, "helm install failed for %s", releaseName)
		}
	case meshkitkube.UNINSTALL:
		uninstall := action.NewUninstall(actionConfig)
		if _, err := uninstall.Run(releaseName); err != nil {
			return errors.Wrapf(err, "helm uninstall failed for %s", releaseName)
		}
	case meshkitkube.UPGRADE:
		up := action.NewUpgrade(actionConfig)
		up.Namespace = ns
		up.Install = true
		if _, err := up.Run(releaseName, ch, values); err != nil {
			return errors.Wrapf(err, "helm upgrade failed for %s", releaseName)
		}
	default:
		return fmt.Errorf("unknown action type: %v", actionType)
	}

	return nil
}
