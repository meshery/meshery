package models

import (
	"errors"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models/connections"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/spf13/viper"
)

// reachableK8sContext is structurally complete enough to yield a Kubernetes
// client without contacting the cluster, so a test reaches the chart-version
// resolution that follows client creation.
func reachableK8sContext() K8sContext {
	return K8sContext{
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
	}
}

// publishedCharts is the catalogue every test in package models resolves
// against: three current charts, the oldest one verified to deploy
// (MinimumOperatorChartVersion), and two that predate it. It deliberately stops
// below the newest *server* releases, which is the real-world condition that
// broke operator deployment - chart publishing trails Meshery Server releases,
// so a current server routinely asks the repository for a chart that does not
// exist yet.
//
// It is one variable rather than one per test file on purpose: the helper that
// builds a MesheryControllersHelper resolves against it too, so two catalogues
// drifting apart would make a reconcile test assert against a chart set its own
// helper never saw.
var publishedCharts = []string{"v1.0.64", "v1.0.63", "v1.0.62", "v1.0.51", "v1.0.40", "v0.7.9"}

// testChartRepo stands in for the repository the catalogue was read from. It is
// deliberately not ChartRepo: resolution must report the repository it actually
// resolved against, so a mirror or an in-cluster repository is never described
// to the user as the default one.
const testChartRepo = "https://example.invalid/charts"

func mustResolve(t *testing.T, published []string, requested string, source OperatorChartVersionSource) (string, string) {
	t.Helper()
	version, reason, err := ResolveOperatorChartVersion(testChartRepo, published, requested, source)
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
// image, no webhook certificate). It must be raised to the oldest *published*
// chart at or above the floor - not to the newest, which would be a larger
// change than the defect requires, and not to the floor constant itself, which
// the repository is not obliged to publish.
func TestDerivedChartVersionIsFlooredToTheOldestWorkingChart(t *testing.T) {
	t.Run("the floor itself is published", func(t *testing.T) {
		version, reason := mustResolve(t, publishedCharts, "v1.0.40", OperatorChartVersionDerived)

		if version != MinimumOperatorChartVersion {
			t.Fatalf("chart version = %q, want the floor %q", version, MinimumOperatorChartVersion)
		}
		if reason == "" {
			t.Fatal("expected the floor to be explained to the user")
		}
	})

	t.Run("the floor itself is not published", func(t *testing.T) {
		// The floor is a policy boundary, not a promise that the repository
		// carries that exact release. Handing MinimumOperatorChartVersion to
		// Helm here would ask for an archive this repository does not publish,
		// and jumping to the newest would upgrade further than the defect
		// requires.
		version, reason := mustResolve(t, []string{"v1.0.64", "v1.0.62", "v1.0.40"}, "v1.0.40", OperatorChartVersionDerived)

		if version != "v1.0.62" {
			t.Fatalf("chart version = %q, want the oldest published chart at or above the floor, v1.0.62", version)
		}
		if reason == "" {
			t.Fatal("expected the floor to be explained to the user")
		}
	})
}

// TestDerivedChartVersionAtOrAboveTheFloorIsUsedVerbatim guards the floor from
// being raised past a chart that works. Each version listed here was rendered
// from its published archive and confirmed to carry no kube-rbac-proxy sidecar
// and ENABLE_WEBHOOKS="false", so substituting any of them would tell a user
// their working chart cannot deploy. They are spelled out as literals rather
// than derived from MinimumOperatorChartVersion precisely so that raising the
// constant fails here instead of silently redefining what the test asserts.
func TestDerivedChartVersionAtOrAboveTheFloorIsUsedVerbatim(t *testing.T) {
	for _, requested := range []string{"v1.0.51", "v1.0.62", "v1.0.63", "v1.0.64"} {
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
			version, reason, err := ResolveOperatorChartVersion(testChartRepo, publishedCharts, tt.requested, OperatorChartVersionRequested)
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
		_, _, err := ResolveOperatorChartVersion(testChartRepo, published, "v1.0.64", OperatorChartVersionDerived)
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
//
// It must also not be described as working. Calling a below-floor chart "the
// oldest working chart" sends the user looking for a fix that has already been
// applied, when the actionable fact is that no working chart is published yet.
func TestFloorFallsBackToNewestWhenNothingReachesIt(t *testing.T) {
	version, reason := mustResolve(t, []string{"v1.0.40", "v0.7.9"}, "v0.7.9", OperatorChartVersionDerived)
	if version != "v1.0.40" {
		t.Fatalf("chart version = %q, want the newest published chart v1.0.40", version)
	}
	if reason == "" {
		t.Fatal("expected the substitution to be explained")
	}
	if strings.Contains(reason, "working chart") {
		t.Fatalf("a below-floor chart must not be described as working: %q", reason)
	}
	if !strings.Contains(reason, MinimumOperatorChartVersion) {
		t.Fatalf("the reason must name the minimum that nothing published reaches: %q", reason)
	}
}

// TestChartVersionSpellingIsSemverNotStringEquality: Helm treats "1.0.64" and
// "v1.0.64" as the same version. Matching on raw string equality rejected an
// explicit pin that the repository does in fact publish, and reported a
// substitution on the derived path for a version that needed none. The version
// handed back is always the repository's own spelling, because that is what
// reaches Helm and what the attached-version comparison is made against.
func TestChartVersionSpellingIsSemverNotStringEquality(t *testing.T) {
	for _, source := range []OperatorChartVersionSource{OperatorChartVersionRequested, OperatorChartVersionDerived} {
		version, reason := mustResolve(t, publishedCharts, "1.0.64", source)
		if version != "v1.0.64" {
			t.Fatalf("chart version = %q, want the repository's own spelling v1.0.64", version)
		}
		if reason != "" {
			t.Fatalf("the same release in another spelling is not a substitution, got reason %q", reason)
		}
	}
}

// TestUnpublishedVersionsNameTheRepositoryThatWasRead: the catalogue comes from
// the deployment config's chart repository, which need not be the default one.
// Telling the user to check a repository Meshery did not read sends them to the
// wrong index.
func TestUnpublishedVersionsNameTheRepositoryThatWasRead(t *testing.T) {
	const mirror = "https://mirror.invalid/charts"

	_, _, err := ResolveOperatorChartVersion(mirror, nil, "v1.0.64", OperatorChartVersionDerived)
	if err == nil || !strings.Contains(err.Error(), mirror) {
		t.Fatalf("an empty catalogue must name the repository it was read from, got %v", err)
	}

	_, reason, err := ResolveOperatorChartVersion(mirror, publishedCharts, "v1.0.66", OperatorChartVersionDerived)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(reason, mirror) {
		t.Fatalf("the substitution reason must name the repository that was read, got %q", reason)
	}
	if strings.Contains(reason, ChartRepo) {
		t.Fatalf("the substitution reason must not name a repository that was not read, got %q", reason)
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

	mch.AddCtxControllerHandlers(reachableK8sContext())

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

// TestReconcileLeavesTheAttachedVersionClearedWhenNoHandlerAttaches guards the
// retry path. When re-attaching cannot resolve a chart version it attaches no
// operator handler and clears attachedOperatorChartVersion *deliberately*, so
// the next reconcile tries again instead of reading a stale value as "already
// at the desired version". Restoring the pre-reconcile value over that clearing
// stranded operator lifecycle: reverting operator.version would then match the
// stale attached version, short-circuit, and never re-attach a handler.
func TestReconcileLeavesTheAttachedVersionClearedWhenNoHandlerAttaches(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
	mch.SetControllersConfig(operatorVersionConfig("v1.0.63", connections.MeshsyncDeploymentModeOperator))
	mch.attachedOperatorChartVersion = "v1.0.64"

	// The index is readable when the reconcile computes the desired version and
	// unreadable by the time re-attaching resolves it again - a TTL expiry, or a
	// blip, between the two reads.
	reads := 0
	mch.chartVersions = func(string, string) ([]string, error) {
		reads++
		if reads == 1 {
			return publishedCharts, nil
		}
		return nil, errors.New("repository unreachable")
	}

	_, redeployed, err := mch.ReconcileOperatorChartVersion(reachableK8sContext(), NewOperatorTracker(false))
	if redeployed {
		t.Fatal("no operator handler attached, so nothing can have been redeployed")
	}
	if err == nil {
		t.Fatal("expected the failed re-attach to be reported")
	}
	if _, ok := mch.GetControllerHandlersForEachContext()[MesheryOperator]; ok {
		t.Fatal("no chart version resolved, so no operator handler may be attached")
	}
	if mch.attachedOperatorChartVersion != "" {
		t.Fatalf("attached chart version = %q, want it left cleared so the next reconcile retries",
			mch.attachedOperatorChartVersion)
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
	// Helm treats "1.0.64" and "v1.0.64" as one release, so they must compare
	// equal rather than by their differing spellings.
	if compareChartVersions("v1.0.64", "1.0.64") != 0 {
		t.Fatal("equal versions must compare equal")
	}
	if compareChartVersions("stable-latest", "v0.0.1") >= 0 {
		t.Fatal("an unparseable version must never outrank a real one")
	}
}
