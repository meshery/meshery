package models

import (
	"errors"
	"fmt"

	"github.com/meshery/meshkit/models/controllers"
	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	// OperatorReleaseNamespace and OperatorReleaseName identify the Helm
	// release Meshery Server installs the Meshery Operator as. They must match
	// meshkit's controllers.applyOperatorHelmChart, which installs the release
	// that operator controller handlers deploy and remove; that function is
	// unexported, so the two literals are restated here rather than shared. A
	// drift between them would make the floor reconcile look for a release that
	// is never there, which it reports as "no release".
	OperatorReleaseNamespace = "meshery"
	OperatorReleaseName      = "meshery-operator"
)

// operatorRelease is the Meshery Operator's Helm release on one cluster, as
// reconcileOperatorChartFloor needs it: what chart is installed, and a way to
// upgrade it.
//
// It exists beside meshkit's operator controller handler rather than inside it
// because the handler offers neither. It cannot report the installed chart -
// GetVersion returns the manager *image* tag, which on the old charts this
// reconcile is for is a moving tag such as stable-latest - and its Deploy
// refuses while the Deployment's rollout is incomplete ("already a Meshery
// Operator is being deployed"), which is permanently the case for an operator
// that crash-loops. A crash-looping operator on a chart too old to run is the
// one this reconcile most needs to reach.
type operatorRelease interface {
	// DeployedChartVersion returns the chart version of the release's
	// currently deployed revision. found is false, with a nil error, when the
	// cluster has no deployed revision of the release at all. A non-nil error
	// is always a structured ErrReadOperatorRelease.
	DeployedChartVersion() (version string, found bool, err error)

	// Upgrade upgrades the existing release to depConfig's chart version from
	// depConfig's repository, with depConfig's Helm overrides. It never
	// installs: a release that vanished since DeployedChartVersion read it
	// fails the upgrade rather than being recreated.
	Upgrade(depConfig controllers.OperatorDeploymentConfig) error
}

// helmOperatorRelease is the operatorRelease backed by the cluster itself.
type helmOperatorRelease struct {
	// secrets is where Helm keeps the release's revisions. meshkit installs
	// with Helm's default Secret storage driver, so that is the driver read.
	secrets corev1.SecretInterface

	// applyHelmChart is mesherykube.Client.ApplyHelmChart: the same Helm entry
	// point the operator handler installs through, so an upgrade here
	// downloads, renders and applies the chart exactly as an install would.
	applyHelmChart func(mesherykube.ApplyHelmChartConfig) error
}

// newHelmOperatorRelease returns the operatorRelease for the cluster client
// talks to.
func newHelmOperatorRelease(client *mesherykube.Client) operatorRelease {
	return &helmOperatorRelease{
		secrets:        client.KubeClient.CoreV1().Secrets(OperatorReleaseNamespace),
		applyHelmChart: client.ApplyHelmChart,
	}
}

// DeployedChartVersion reads the deployed revision's chart version from Helm's
// own release record - the same value `helm -n meshery list` prints in its
// CHART column - rather than from a label on a rendered resource, which a chart
// is free to omit or spell differently across versions.
//
// Only the *deployed* revision counts. A failed upgrade leaves the previous
// revision deployed and the operator running it, so reading the newest revision
// regardless of status would report a chart that never took effect, and stop
// retrying an upgrade that did not happen.
func (r *helmOperatorRelease) DeployedChartVersion() (string, bool, error) {
	rel, err := storage.Init(driver.NewSecrets(r.secrets)).Deployed(OperatorReleaseName)
	if err != nil {
		// Helm reports both "no revision of this release exists" and "none is
		// deployed" as ErrNoDeployedReleases here. Every other error - a
		// forbidden list, an unreachable API server - is a failure to look,
		// not an absence, and is returned.
		if errors.Is(err, driver.ErrNoDeployedReleases) {
			return "", false, nil
		}
		return "", false, ErrReadOperatorRelease(err.Error())
	}
	if rel.Chart == nil || rel.Chart.Metadata == nil || rel.Chart.Metadata.Version == "" {
		return "", false, ErrReadOperatorRelease(fmt.Sprintf("the deployed revision %d of Helm release %s/%s records no chart version", rel.Version, OperatorReleaseNamespace, OperatorReleaseName))
	}
	return rel.Chart.Metadata.Version, true, nil
}

// Upgrade applies the chart as an explicit Helm UPGRADE of the existing
// release, with the same namespace, release name, chart and overrides
// meshkit's operator handler installs with.
//
// Helm runs the chart's pre-upgrade hooks before it touches any release
// resource. Every published chart at or above MinimumOperatorChartVersion
// carries a pre-install,pre-upgrade CRD update Job (crds.updateJob, enabled by
// default) that server-side-applies the current brokers.meshery.io and
// meshsyncs.meshery.io CRDs, which Helm itself never upgrades from a chart's
// crds/ directory. If that Job fails, Helm aborts the upgrade before the
// operator Deployment changes, so a failed upgrade never leaves a new operator
// running against old CRDs.
//
// The overrides are passed whole, as an install passes them, and are never
// empty (setOverrideValues always names the operator chart's subcharts). That
// matters: Helm reuses a release's previous values on an upgrade that supplies
// none, and the previous values here are those of a chart too old to run.
func (r *helmOperatorRelease) Upgrade(depConfig controllers.OperatorDeploymentConfig) error {
	return r.applyHelmChart(mesherykube.ApplyHelmChartConfig{
		Namespace:   OperatorReleaseNamespace,
		ReleaseName: OperatorReleaseName,
		ChartLocation: mesherykube.HelmChartLocation{
			Repository: depConfig.HelmChartRepo,
			Chart:      OperatorChartName,
			Version:    depConfig.MesheryReleaseVersion,
		},
		Action:         mesherykube.UPGRADE,
		OverrideValues: depConfig.GetHelmOverrides(false),
	})
}
