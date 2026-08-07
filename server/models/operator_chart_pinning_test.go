package models

import (
	"errors"
	"testing"

	"github.com/meshery/meshery/server/models/connections"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/spf13/viper"
)

// publishedCharts is the catalogue these tests resolve against: the two current
// charts, the first one that works (MinimumOperatorChartVersion), and two that
// predate it. It deliberately stops below the newest *server* releases, which
// is the real-world condition that broke operator deployment - chart publishing
// trails Meshery Server releases, so a current server routinely asks the
// repository for a chart that does not exist yet.
var publishedCharts = []string{"v1.0.64", "v1.0.63", "v1.0.62", "v1.0.40", "v0.7.9"}

func mustResolve(t *testing.T, published []string, requested string, source OperatorChartVersionSource) (string, string) {
	t.Helper()
	version, reason, err := ResolveOperatorChartVersion(published, requested, source)
	if err != nil {
		t.Fatalf("ResolveOperatorChartVersion(%q) returned an unexpected error: %v", requested, err)
	}
	return version, reason
}

// TestDerivedChartVersionIsAlwaysPublished is the headline regression. Before
// this change the chart version was the Meshery Server release verbatim, so a
// server released after the last chart publish asked Helm for an archive that
// does not exist and the operator never deployed. Resolution must now land on a
// version the repository actually publishes, and must say that it did.
func TestDerivedChartVersionIsAlwaysPublished(t *testing.T) {
	version, _ := mustResolve(t, publishedCharts, "v1.0.66", OperatorChartVersionDerived)

	if version != "v1.0.64" {
		t.Fatalf("chart version = %q, want the newest published chart v1.0.64", version)
	}
	for _, v := range publishedCharts {
		if v == version {
			return
		}
	}
	t.Fatalf("resolved chart version %q is not published", version)
}

func TestDerivedChartVersionSubstitutionIsNeverSilent(t *testing.T) {
	_, reason := mustResolve(t, publishedCharts, "v1.0.66", OperatorChartVersionDerived)
	if reason == "" {
		t.Fatal("expected a user-readable reason explaining the substitution, got none")
	}
}

// TestDerivedChartVersionIsFlooredToTheOldestWorkingChart pins the fix for the
// reported failure: a server old enough to predate MinimumOperatorChartVersion
// requests a chart that *is* published but cannot run (retired kube-rbac-proxy
// image, no webhook certificate). It must be raised to the oldest working
// chart - not to the newest, which would be a larger change than the defect
// requires.
func TestDerivedChartVersionIsFlooredToTheOldestWorkingChart(t *testing.T) {
	version, reason := mustResolve(t, publishedCharts, "v1.0.40", OperatorChartVersionDerived)

	if version != MinimumOperatorChartVersion {
		t.Fatalf("chart version = %q, want the floor %q", version, MinimumOperatorChartVersion)
	}
	if reason == "" {
		t.Fatal("expected the floor to be explained to the user")
	}
}

func TestDerivedChartVersionAtOrAboveTheFloorIsUsedVerbatim(t *testing.T) {
	for _, requested := range []string{MinimumOperatorChartVersion, "v1.0.64"} {
		version, reason := mustResolve(t, publishedCharts, requested, OperatorChartVersionDerived)
		if version != requested {
			t.Fatalf("chart version = %q, want %q used verbatim", version, requested)
		}
		if reason != "" {
			t.Fatalf("expected no substitution for %q, got reason %q", requested, reason)
		}
	}
}

// TestMovingChartVersionNeverReachesHelm covers every way a server can fail to
// name a release: an unstamped source build, viper's "Not Set" placeholder, and
// the moving channel tags. None of them identifies a published archive, so none
// may be handed to Helm as a chart version.
func TestMovingChartVersionNeverReachesHelm(t *testing.T) {
	for _, requested := range []string{"", "   ", "Not Set", "edge-latest", "stable-latest", "latest"} {
		t.Run(requested, func(t *testing.T) {
			if isPinnedChartVersion(requested) {
				t.Fatalf("%q must not be treated as a pinned chart version", requested)
			}
			version, reason := mustResolve(t, publishedCharts, requested, OperatorChartVersionDerived)
			if version != "v1.0.64" {
				t.Fatalf("chart version = %q, want the newest published chart", version)
			}
			if reason == "" {
				t.Fatal("expected the fallback to be explained")
			}
		})
	}
}

// TestExplicitChartVersionFailsLoudlyWhenUnusable asserts the deliberate
// asymmetry: a version nobody chose may be corrected, but one an operator typed
// is never quietly replaced - the deployment fails so the mistake is visible.
func TestExplicitChartVersionFailsLoudlyWhenUnusable(t *testing.T) {
	tests := []struct {
		name      string
		requested string
		wantCode  string
	}{
		{
			name:      "a version the repository does not publish",
			requested: "v9.9.9",
			wantCode:  ErrOperatorChartNotPublishedCode,
		},
		{
			name:      "a moving tag, which names no specific archive",
			requested: "stable-latest",
			wantCode:  ErrOperatorChartNotPinnedCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, reason, err := ResolveOperatorChartVersion(publishedCharts, tt.requested, OperatorChartVersionRequested)
			if err == nil {
				t.Fatalf("expected %q to fail, got version %q", tt.requested, version)
			}
			if version != "" || reason != "" {
				t.Fatalf("a failed resolution must yield nothing to deploy, got version %q reason %q", version, reason)
			}
			if code := meshkiterrors.GetCode(err); code != tt.wantCode {
				t.Fatalf("error code = %q, want %q (from %v)", code, tt.wantCode, err)
			}
		})
	}
}

// TestExplicitChartVersionBelowTheFloorIsHonored keeps the escape hatch open:
// the floor exists to rescue servers that did not choose their chart version.
// Someone who deliberately pins an old chart is not overruled.
func TestExplicitChartVersionBelowTheFloorIsHonored(t *testing.T) {
	version, reason := mustResolve(t, publishedCharts, "v0.7.9", OperatorChartVersionRequested)
	if version != "v0.7.9" {
		t.Fatalf("chart version = %q, want the explicitly requested v0.7.9", version)
	}
	if reason != "" {
		t.Fatalf("an honored explicit request is not a substitution, got reason %q", reason)
	}
}

func TestResolutionFailsWhenNothingUsableIsPublished(t *testing.T) {
	for _, published := range [][]string{nil, {}, {"stable-latest", "edge-latest"}} {
		_, _, err := ResolveOperatorChartVersion(published, "v1.0.64", OperatorChartVersionDerived)
		if err == nil {
			t.Fatalf("expected a failure when the repository publishes %v", published)
		}
		if code := meshkiterrors.GetCode(err); code != ErrNoOperatorChartPublishedCode {
			t.Fatalf("error code = %q, want %q", code, ErrNoOperatorChartPublishedCode)
		}
	}
}

// TestFloorFallsBackToNewestWhenNothingReachesIt guards the degenerate case
// where the repository has not yet published anything at or above the floor:
// the newest published chart is still better than a version known to be broken,
// and it must not resolve to "".
func TestFloorFallsBackToNewestWhenNothingReachesIt(t *testing.T) {
	version, reason := mustResolve(t, []string{"v1.0.40", "v0.7.9"}, "v0.7.9", OperatorChartVersionDerived)
	if version != "v1.0.40" {
		t.Fatalf("chart version = %q, want the newest published chart v1.0.40", version)
	}
	if reason == "" {
		t.Fatal("expected the substitution to be explained")
	}
}

// TestPinnedDeploymentConfigSurfacesAnUnreadableIndex asserts the failure is
// loud rather than a guess: if the repository index cannot be read there is no
// way to know which versions exist, so no version may be deployed.
func TestPinnedDeploymentConfigSurfacesAnUnreadableIndex(t *testing.T) {
	mch := newTestControllersHelper(t)
	indexErr := errors.New("repository unreachable")
	mch.chartVersions = func(string, string) ([]string, error) { return nil, indexErr }

	_, _, err := mch.pinnedOperatorDeploymentConfig()
	if err == nil {
		t.Fatal("expected an unreadable chart index to fail the resolution")
	}
	if !errors.Is(err, indexErr) {
		t.Fatalf("expected the index failure to be reported, got %v", err)
	}
}

// TestUnreadableIndexStillAttachesBrokerAndMeshSync asserts the blast radius of
// an unreachable chart repository is confined to operator lifecycle. The Broker
// and MeshSync handlers only observe what is already in the cluster, so
// withholding them too would turn a momentary repository blip into lost
// discovery on a cluster whose operator is already installed and healthy.
func TestUnreadableIndexStillAttachesBrokerAndMeshSync(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.chartVersions = func(string, string) ([]string, error) { return nil, errors.New("repository unreachable") }
	mch.attachedOperatorChartVersion = "v1.0.64"

	// A structurally complete context yields a Kubernetes client without
	// contacting the cluster, so this reaches the chart-version resolution.
	mch.AddCtxControllerHandlers(K8sContext{
		ID:   "ctx-1",
		Name: "test-cluster",
		Cluster: map[string]interface{}{
			"name":    "test-cluster",
			"cluster": map[string]interface{}{"server": "https://127.0.0.1:6443"},
		},
		Auth: map[string]interface{}{
			"name": "test-user",
			"user": map[string]interface{}{"token": "test-token"},
		},
		Server: "https://127.0.0.1:6443",
	})

	handlers := mch.GetControllerHandlersForEachContext()
	if _, ok := handlers[MesheryOperator]; ok {
		t.Fatal("no chart version is safe to install, so no operator handler may be attached")
	}
	for controller, name := range map[MesheryController]string{MesheryBroker: "broker", Meshsync: "meshsync"} {
		if h, ok := handlers[controller]; !ok || h == nil {
			t.Fatalf("the %s handler needs no chart version and must still be attached", name)
		}
	}
	if mch.GetOperatorError() == nil {
		t.Fatal("the reason operator lifecycle is withheld must be recorded for diagnostics")
	}
	if mch.attachedOperatorChartVersion != "" {
		t.Fatalf("attached chart version = %q, want it cleared so the next reconcile retries", mch.attachedOperatorChartVersion)
	}
}

// TestPinnedDeploymentConfigQueriesTheOperatorChart asserts the pin is computed
// against the operator chart in the server-wide repository, and that pinning
// replaces the chart version and nothing else.
func TestPinnedDeploymentConfigQueriesTheOperatorChart(t *testing.T) {
	mch := newTestControllersHelper(t)
	var gotRepo, gotChart string
	mch.chartVersions = func(repo, chart string) ([]string, error) {
		gotRepo, gotChart = repo, chart
		return publishedCharts, nil
	}

	depConfig, _, err := mch.pinnedOperatorDeploymentConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotRepo != "https://example.invalid/charts" {
		t.Fatalf("queried repo = %q, want the server-wide chart repository", gotRepo)
	}
	if gotChart != OperatorChartName {
		t.Fatalf("queried chart = %q, want %q", gotChart, OperatorChartName)
	}
	if depConfig.HelmChartRepo != "https://example.invalid/charts" || depConfig.GetHelmOverrides == nil {
		t.Fatal("pinning must replace the chart version only, leaving the server-wide fields intact")
	}
	if mch.oprDepConfig.MesheryReleaseVersion != bootChartVersion {
		t.Fatalf("server-wide deployment config was mutated: %q", mch.oprDepConfig.MesheryReleaseVersion)
	}
}

// TestReconcileComparesPinnedVersionsNotRequestedOnes pins a loop that the
// naive comparison would create. The attached handler always records a
// *published* version, so comparing it against the raw request would differ
// forever whenever the request is being substituted - re-running the Helm
// upgrade against the cluster on every reconcile.
func TestReconcileComparesPinnedVersionsNotRequestedOnes(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
	// A server release with no published chart: pinning lands on v1.0.64, which
	// is exactly what the attached handler already carries.
	mch.oprDepConfig.MesheryReleaseVersion = "v1.0.66"
	mch.SetControllersConfig(operatorVersionConfig("", connections.MeshsyncDeploymentModeOperator))
	mch.attachedOperatorChartVersion = "v1.0.64"

	chartVersion, redeployed, err := mch.ReconcileOperatorChartVersion(
		K8sContext{ID: "ctx-1", Name: "test-cluster"},
		NewOperatorTracker(false),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if redeployed {
		t.Fatal("the pinned version already matches the attached handler; redeploying is an upgrade loop")
	}
	if chartVersion != "v1.0.64" {
		t.Fatalf("chart version = %q, want the pinned v1.0.64", chartVersion)
	}
}

// TestNewOperatorDeploymentConfigNeverStampsAMovingVersion asserts the boot-time
// config carries either a real release or nothing at all. It previously reached
// for the newest *GitHub release*, which is not the newest *published chart* -
// spending a rate-limited API call to produce a version that frequently does not
// exist in the chart repository.
func TestNewOperatorDeploymentConfigNeverStampsAMovingVersion(t *testing.T) {
	tests := []struct {
		build string
		want  string
	}{
		{build: "", want: ""},
		{build: "Not Set", want: ""},
		{build: "edge-latest", want: ""},
		{build: "stable-latest", want: ""},
		{build: "v1.0.64", want: "v1.0.64"},
	}

	original := viper.GetString("BUILD")
	t.Cleanup(func() { viper.Set("BUILD", original) })

	for _, tt := range tests {
		t.Run(tt.build, func(t *testing.T) {
			viper.Set("BUILD", tt.build)
			got := NewOperatorDeploymentConfig(nil).MesheryReleaseVersion
			if got != tt.want {
				t.Fatalf("MesheryReleaseVersion = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCompareChartVersionsOrdersBySemverNotLexically(t *testing.T) {
	// "v1.0.9" sorts after "v1.0.64" lexically, which would pick a broken chart
	// as "newest".
	if compareChartVersions("v1.0.64", "v1.0.9") <= 0 {
		t.Fatal("v1.0.64 must compare greater than v1.0.9")
	}
	if compareChartVersions("v1.0.63", MinimumOperatorChartVersion) != 0 {
		t.Fatal("equal versions must compare equal")
	}
	if compareChartVersions("stable-latest", "v0.0.1") >= 0 {
		t.Fatal("an unparseable version must never outrank a real one")
	}
}
