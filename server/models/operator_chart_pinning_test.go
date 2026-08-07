package models

import (
	"errors"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models/connections"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/spf13/viper"
)

// reachableK8sContext is structurally complete enough to yield a Kubernetes
// client without contacting the cluster, so a test reaches the chart-version
// resolution that follows client creation.
//
// A context this complete makes AddCtxControllerHandlers attach REAL meshkit
// handlers, so any test that goes on to reach Deploy or Undeploy performs live
// Helm and cluster I/O: meshkit downloads the multi-megabyte index and the chart
// archive from the configured repository before it touches the cluster at all,
// and 127.0.0.1:6443 is a live cluster on any machine running Docker Desktop
// Kubernetes or k3s - where a Deploy would really install the operator into it.
// Never let a test in this package reach an install or a removal on a handler
// this produced. Assert the guard through operatorInstallTarget, or swap a
// stubController into ctxControllerHandlers first.
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

// unresolvableK8sContext is the counterpart to reachableK8sContext. It carries
// no cluster and no auth, so mesherykube.New rejects it and
// AddCtxControllerHandlers returns before it attaches a handler or resolves a
// chart version - which is what keeps a test using it away from meshkit's Helm
// client entirely.
//
// That it fails, rather than quietly falling back to whatever kubeconfig the
// machine happens to have, rests on a meshkit detail: ProcessConfig runs
// clientcmd.Validate, which rejects the generated document, and DetectKubeConfig
// returns that error immediately instead of continuing to rest.InClusterConfig,
// $KUBECONFIG or ~/.kube/config. Should that ever change, mesherykube.New would
// hand back a live client aimed at the developer's current kubectl context and a
// test relying on this seam could install the operator into it for real. Every
// such test therefore calls requireUnresolvableK8sContext first rather than
// assuming the seam still holds.
func unresolvableK8sContext() K8sContext {
	return K8sContext{ID: testContextID, Name: "unresolvable-cluster"}
}

// requireUnresolvableK8sContext fails the test unless unresolvableK8sContext
// still cannot yield a Kubernetes client. Call it before any code path that
// would otherwise install through a real handler, so a meshkit change that
// relaxes the validation above surfaces here instead of as a live helm install.
func requireUnresolvableK8sContext(t *testing.T) {
	t.Helper()
	probe := newTestControllersHelper(t)
	probe.AddCtxControllerHandlers(unresolvableK8sContext())
	if probe.GetOperatorError() == nil {
		t.Fatal("unresolvableK8sContext now yields a Kubernetes client: meshkit's kubeconfig validation changed, and tests built on this seam can reach a real Helm install")
	}
	if probe.attachedOperatorHandler() != nil {
		t.Fatal("unresolvableK8sContext attached an operator handler: the seam these tests rely on no longer holds")
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
// requests a chart that *is* published but cannot run (a kube-rbac-proxy
// sidecar that affected clusters cannot pull, no webhook certificate). It must
// be raised to the oldest *published*
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

// TestUnreadableIndexWithholdsInstallationNotObservation draws the line the
// chart version actually sits on. meshkit's operator handler reads its
// deployment config in Deploy and Undeploy only - GetStatus and GetVersion
// never touch it - so a version is what it takes to *install* the operator, not
// to *watch* one. Withholding the handler when the repository could not be read
// therefore cost a healthy, already-installed operator its status and its image
// tag, which is the exact value the troubleshooting guide tells users to read.
//
// All three handlers attach; only installation is withheld, and why is recorded
// for the diagnostics API.
func TestUnreadableIndexWithholdsInstallationNotObservation(t *testing.T) {
	mch := newTestControllersHelper(t)
	indexErr := errors.New("repository unreachable")
	mch.chartVersions = func(string, string) ([]string, error) { return nil, indexErr }
	mch.attachedOperatorChartVersion = "v1.0.64"

	mch.AddCtxControllerHandlers(reachableK8sContext())

	handlers := mch.GetControllerHandlersForEachContext()
	for controller, name := range map[MesheryController]string{
		MesheryBroker:   "broker",
		Meshsync:        "meshsync",
		MesheryOperator: "operator",
	} {
		if h, ok := handlers[controller]; !ok || h == nil {
			t.Fatalf("the %s handler observes the cluster and must be attached", name)
		}
	}
	if chartErr := mch.GetOperatorChartError(); !errors.Is(chartErr, indexErr) {
		t.Fatalf("operator chart error = %v, want the index failure that withholds installation", chartErr)
	}
	if mch.GetOperatorError() == nil {
		t.Fatal("the reason installation is withheld must be recorded for diagnostics")
	}
	if mch.attachedOperatorChartVersion != "" {
		t.Fatalf("attached chart version = %q, want it cleared so the next reconcile retries", mch.attachedOperatorChartVersion)
	}
}

// TestChartResolutionFailureIsClearedByALaterSuccess: the refusal is a live
// condition, not a latch. A repository that comes back must not leave
// installation withheld until the server restarts.
func TestChartResolutionFailureIsClearedByALaterSuccess(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.chartVersions = func(string, string) ([]string, error) { return nil, errors.New("repository unreachable") }
	mch.AddCtxControllerHandlers(reachableK8sContext())
	if mch.GetOperatorChartError() == nil {
		t.Fatal("expected the unreadable index to withhold installation")
	}

	mch.chartVersions = func(string, string) ([]string, error) { return publishedCharts, nil }
	mch.AddCtxControllerHandlers(reachableK8sContext())

	if err := mch.GetOperatorChartError(); err != nil {
		t.Fatalf("operator chart error = %v, want it cleared once a version resolved", err)
	}
	if mch.attachedOperatorChartVersion != bootChartVersion {
		t.Fatalf("attached chart version = %q, want %q", mch.attachedOperatorChartVersion, bootChartVersion)
	}
}

// testContextID is the Kubernetes context identifier the lifecycle call sites
// are given. It is non-empty on purpose: ErrOperatorHandlerNotAttached renders
// it into the cause the user reads, and an operator triaging several
// connections needs to be told which one lost its handler.
const testContextID = "ctx-1"

// stubController is an IMesheryController that records what it was asked to do,
// so the deploy/undeploy call sites can be tested without a cluster.
type stubController struct {
	status    controllers.MesheryControllerStatus
	deploys   int
	undeploys int
	deployErr error
}

func (s *stubController) GetName() string                                { return "stub" }
func (s *stubController) GetStatus() controllers.MesheryControllerStatus { return s.status }
func (s *stubController) Deploy(bool) error                              { s.deploys++; return s.deployErr }
func (s *stubController) Undeploy() error                                { s.undeploys++; return nil }
func (s *stubController) GetPublicEndpoint() (string, error)             { return "", nil }
func (s *stubController) GetVersion() (string, error)                    { return "1.0.5", nil }
func (s *stubController) GetEndpointForPort(string) (string, error)      { return "", nil }

// TestDeployIsRefusedWhileNoChartVersionResolves is the other half of attaching
// the handler unconditionally: the handler exists for observation, so the
// refusal has to live at the point of installation. Handing Helm a version the
// repository does not publish would surface as an opaque chart-not-found error
// instead of this connection's recorded resolution failure.
func TestDeployIsRefusedWhileNoChartVersionResolves(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
	stub := &stubController{status: controllers.NotDeployed}
	mch.ctxControllerHandlers = map[MesheryController]controllers.IMesheryController{MesheryOperator: stub}
	mch.ctxOperatorStatus = controllers.NotDeployed

	chartErr := ErrNoOperatorChartPublished(OperatorChartName, testChartRepo)
	mch.setOperatorChartError(chartErr)

	mch.DeployUndeployedOperators(NewOperatorTracker(false), testContextID)

	if stub.deploys != 0 {
		t.Fatalf("Deploy was called %d times with no installable chart version", stub.deploys)
	}
	if !errors.Is(mch.GetOperatorError(), chartErr) {
		t.Fatalf("operator error = %v, want the chart resolution failure", mch.GetOperatorError())
	}

	// Cleared, the same handler installs.
	mch.setOperatorChartError(nil)
	mch.DeployUndeployedOperators(NewOperatorTracker(false), testContextID)
	if stub.deploys != 1 {
		t.Fatalf("Deploy was called %d times once a chart version resolved, want 1", stub.deploys)
	}
}

// TestUserInitiatedDeployIsRefusedWhileNoChartVersionResolves closes the second
// way in. The GraphQL changeOperatorStatus mutation used to reach around the
// helper for a raw meshkit handler, so the guard the connect-time path applies
// did not exist on the path a user actually clicks. Both now share
// operatorInstallTarget and must refuse identically.
func TestUserInitiatedDeployIsRefusedWhileNoChartVersionResolves(t *testing.T) {
	mch := newTestControllersHelper(t)
	stub := &stubController{status: controllers.NotDeployed}
	mch.ctxControllerHandlers = map[MesheryController]controllers.IMesheryController{MesheryOperator: stub}

	// The repository is still unreadable, so the retry this path performs
	// resolves again and fails again. The refusal must then carry the *fresh*
	// failure rather than falling through to a Deploy.
	indexErr := errors.New("repository unreachable")
	mch.chartVersions = func(string, string) ([]string, error) { return nil, indexErr }
	mch.setOperatorChartError(ErrNoOperatorChartPublished(OperatorChartName, testChartRepo))

	err := mch.SetOperatorDeployment(reachableK8sContext(), true)
	if !errors.Is(err, indexErr) {
		t.Fatalf("err = %v, want the re-resolution's own failure returned to the caller", err)
	}
	if stub.deploys != 0 {
		t.Fatalf("Deploy was called %d times with no installable chart version", stub.deploys)
	}
	if !errors.Is(mch.GetOperatorError(), indexErr) {
		t.Fatalf("operator error = %v, want the chart resolution failure recorded", mch.GetOperatorError())
	}

	// With nothing latched there is no retry, and the same handler installs.
	mch = newTestControllersHelper(t)
	stub = &stubController{status: controllers.NotDeployed}
	mch.ctxControllerHandlers = map[MesheryController]controllers.IMesheryController{MesheryOperator: stub}
	if err := mch.SetOperatorDeployment(K8sContext{ID: testContextID}, true); err != nil {
		t.Fatalf("unexpected error with a resolvable chart version: %v", err)
	}
	if stub.deploys != 1 {
		t.Fatalf("Deploy was called %d times with a resolvable chart version, want 1", stub.deploys)
	}
}

// TestUserInitiatedDeployRetriesResolutionAfterTheRepositoryRecovers is what
// makes ErrHelmChartIndex's remediation - "confirm the repository is reachable,
// then retry" - a true instruction rather than a dead end.
//
// operatorChartError is written only where handlers are attached, so a connect
// that landed during a repository outage used to refuse every later install for
// the life of that connection however long ago the repository came back.
// Clearing the latch alone would not have been enough either: the handler
// attached on the failure path still carries the raw unresolved chart version,
// so re-resolving and re-attaching have to happen together.
// It is asserted through the two steps SetOperatorDeployment composes for a
// latched deploy - AddCtxControllerHandlers then operatorInstallTarget - rather
// than by calling SetOperatorDeployment itself, because a successful retry ends
// in a Deploy on a real meshkit handler, which downloads the chart archive and
// would install the operator for real on a machine where 127.0.0.1:6443 is a
// live cluster. That SetOperatorDeployment performs this retry at all is pinned
// by TestUserInitiatedDeployIsRefusedWhileNoChartVersionResolves, which only
// passes because the refusal it gets back is the re-resolution's own failure
// rather than the stale latch.
func TestUserInitiatedDeployRetriesResolutionAfterTheRepositoryRecovers(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.setOperatorChartError(ErrNoOperatorChartPublished(OperatorChartName, testChartRepo))
	mch.attachedOperatorChartVersion = ""

	// The catalogue is readable again, so re-attaching must resolve rather than
	// leave the refusal standing.
	mch.AddCtxControllerHandlers(reachableK8sContext())

	if chartErr := mch.GetOperatorChartError(); chartErr != nil {
		t.Fatalf("operator chart error = %v, want the latch cleared by a resolution that succeeded", chartErr)
	}
	if mch.attachedOperatorChartVersion != bootChartVersion {
		t.Fatalf("attached chart version = %q, want %q re-attached by the retry",
			mch.attachedOperatorChartVersion, bootChartVersion)
	}

	handler, attached, err := mch.operatorInstallTarget(testContextID)
	if err != nil {
		t.Fatalf("the install guard still refuses after a successful re-resolution: %v", err)
	}
	if !attached || handler == nil {
		t.Fatal("a resolved re-attach must yield a handler to install through")
	}
}

// TestUserInitiatedUndeployIsNeverSilent: a user who switches the operator off
// must be told when nothing could be done, whether or not Meshery ever observed
// an operator here. This is the case the teardown gate below deliberately does
// not cover.
func TestUserInitiatedUndeployIsNeverSilent(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.ctxControllerHandlers = nil

	err := mch.SetOperatorDeployment(K8sContext{ID: testContextID}, false)
	if err == nil {
		t.Fatal("a user-initiated undeploy with no handler must not report success")
	}
	if code := meshkiterrors.GetCode(err); code != ErrOperatorHandlerNotAttachedCode {
		t.Fatalf("error code = %q, want %q (from %v)", code, ErrOperatorHandlerNotAttachedCode, err)
	}
	if !strings.Contains(err.Error(), testContextID) {
		t.Fatalf("the error must name the context it is about, got %v", err)
	}
}

// TestUndeployIsNotBlockedByAnUnresolvableChartVersion: removal is the
// direction to attempt rather than refuse. Refusing would leave the operator
// running on a cluster the user asked to have it taken off - even though the
// attempt itself may still fail in meshkit, which downloads the chart archive
// for UNINSTALL from the very repository that could not be read.
func TestUndeployIsNotBlockedByAnUnresolvableChartVersion(t *testing.T) {
	for _, tt := range []struct {
		name string
		act  func(*MesheryControllersHelper)
	}{
		{
			name: "reconcile",
			act: func(mch *MesheryControllersHelper) {
				mch.UndeployDeployedOperators(NewOperatorTracker(false), testContextID)
			},
		},
		{
			name: "user initiated",
			act: func(mch *MesheryControllersHelper) {
				if err := mch.SetOperatorDeployment(K8sContext{ID: testContextID}, false); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			mch := newTestControllersHelper(t)
			stub := &stubController{status: controllers.Deployed}
			mch.ctxControllerHandlers = map[MesheryController]controllers.IMesheryController{MesheryOperator: stub}
			mch.ctxOperatorStatus = controllers.Deployed
			mch.setOperatorChartError(ErrNoOperatorChartPublished(OperatorChartName, testChartRepo))

			tt.act(mch)

			if stub.undeploys != 1 {
				t.Fatalf("Undeploy was called %d times, want 1: an unresolvable chart version must not strand the operator", stub.undeploys)
			}
		})
	}
}

// TestMissingOperatorHandlerIsNeverASilentNoOp pins the condition that used to
// be unreachable and is now routine: with no operator handler attached, deploy
// and undeploy did nothing and said nothing, so a user who switched the
// operator off saw success while it kept running on the cluster.
func TestMissingOperatorHandlerIsNeverASilentNoOp(t *testing.T) {
	tests := []struct {
		name   string
		status controllers.MesheryControllerStatus
		act    func(*MesheryControllersHelper)
	}{
		{
			name:   "deploy",
			status: controllers.NotDeployed,
			act: func(mch *MesheryControllersHelper) {
				mch.DeployUndeployedOperators(NewOperatorTracker(false), testContextID)
			},
		},
		{
			name:   "undeploy",
			status: controllers.Deployed,
			act: func(mch *MesheryControllersHelper) {
				mch.UndeployDeployedOperators(NewOperatorTracker(false), testContextID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, handlers := range []map[MesheryController]controllers.IMesheryController{
				nil,
				{MesheryBroker: &stubController{}},
			} {
				mch := newTestControllersHelper(t)
				mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
				mch.ctxControllerHandlers = handlers
				mch.ctxOperatorStatus = tt.status

				tt.act(mch)

				err := mch.GetOperatorError()
				if err == nil {
					t.Fatalf("%s with no operator handler reported nothing; it must not read as success", tt.name)
				}
				if code := meshkiterrors.GetCode(err); code != ErrOperatorHandlerNotAttachedCode {
					t.Fatalf("error code = %q, want %q (from %v)", code, ErrOperatorHandlerNotAttachedCode, err)
				}
				if !strings.Contains(err.Error(), testContextID) {
					t.Fatalf("%s: the error must name the context it is about, got %v", tt.name, err)
				}
			}
		})
	}
}

// TestTeardownOfANeverConnectedConnectionIsQuiet: disconnect and delete run
// UndeployDeployedOperators unconditionally, and a connection whose cluster was
// never reachable has no handler *and never had an operator*. Alerting that its
// removal failed would be a second alert for a non-event, on top of the
// Kubernetes-client failure the connect attempt already reported.
//
// The discriminator is whether the operator status was ever genuinely observed,
// which is exactly what distinguishes "we saw one and can no longer act on it"
// from "we never reached this cluster".
func TestTeardownOfANeverConnectedConnectionIsQuiet(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.ctxControllerHandlers = nil
	if mch.ctxOperatorStatus != controllers.Unknown {
		t.Fatalf("seeded operator status = %v, want the unobserved Unknown this test is about", mch.ctxOperatorStatus)
	}

	mch.UndeployDeployedOperators(NewOperatorTracker(false), testContextID)

	if err := mch.GetOperatorError(); err != nil {
		t.Fatalf("tearing down a never-connected connection reported %v; it has no operator to fail to remove", err)
	}
}

// TestMissingHandlerNeverClobbersTheSpecificDiagnostic: "no operator handler is
// attached" is the *consequence* of the kubeconfig or Kubernetes-client failure
// that AddCtxControllerHandlers already recorded by name. Letting the
// consequence overwrite the cause meant tearing a connection down degraded the
// very diagnostic this work exists to provide.
func TestMissingHandlerNeverClobbersTheSpecificDiagnostic(t *testing.T) {
	mch := newTestControllersHelper(t)
	clientErr := errors.New("failed to create Kubernetes client: no such host")
	mch.setOperatorError(clientErr)
	mch.ctxControllerHandlers = nil
	// Observed, so the teardown gate above does not apply and the report runs.
	mch.ctxOperatorStatus = controllers.Deployed

	mch.UndeployDeployedOperators(NewOperatorTracker(false), testContextID)

	if !errors.Is(mch.GetOperatorError(), clientErr) {
		t.Fatalf("operator error = %v, want the actionable client failure to survive the teardown", mch.GetOperatorError())
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

// TestReconcileLeavesTheAttachedVersionClearedWhenResolutionFails guards the
// retry path. When re-attaching cannot resolve a chart version it clears
// attachedOperatorChartVersion *deliberately*, so the next reconcile tries
// again instead of reading a stale value as "already at the desired version".
// Restoring the pre-reconcile value over that clearing stranded operator
// lifecycle: reverting operator.version would then match the stale attached
// version, short-circuit, and never resolve again.
//
// The handler itself stays attached - it is what keeps reporting the operator's
// status and image tag - so the restore must key on the cleared version, not on
// whether a handler is present.
func TestReconcileLeavesTheAttachedVersionClearedWhenResolutionFails(t *testing.T) {
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
		t.Fatal("no chart version resolved, so nothing can have been redeployed")
	}
	if err == nil {
		t.Fatal("expected the failed re-attach to be reported")
	}
	if h := mch.GetControllerHandlersForEachContext()[MesheryOperator]; h == nil {
		t.Fatal("the operator handler observes the cluster and must stay attached")
	}
	if mch.GetOperatorChartError() == nil {
		t.Fatal("installation must stay withheld until a chart version resolves")
	}
	if mch.attachedOperatorChartVersion != "" {
		t.Fatalf("attached chart version = %q, want it left cleared so the next reconcile retries",
			mch.attachedOperatorChartVersion)
	}
}

// TestInstallGuardRefusesOnALatchedChartError pins the guard itself, which all
// three install sites - the FSM deploy, the user-initiated deploy, and the
// chart-version reconcile - now route through. A handler is attached for
// observation even when no version could be resolved, so "a handler exists" is
// never on its own enough to install.
func TestInstallGuardRefusesOnALatchedChartError(t *testing.T) {
	mch := newTestControllersHelper(t)
	stub := &stubController{status: controllers.NotDeployed}
	mch.ctxControllerHandlers = map[MesheryController]controllers.IMesheryController{MesheryOperator: stub}

	latched := ErrNoOperatorChartPublished(OperatorChartName, testChartRepo)
	mch.setOperatorChartError(latched)

	handler, attached, err := mch.operatorInstallTarget(testContextID)
	if !errors.Is(err, latched) {
		t.Fatalf("err = %v, want the latched chart failure", err)
	}
	if !attached {
		t.Fatal("a handler was attached, so the refusal is about the chart version, not a missing handler")
	}
	if handler != nil {
		t.Fatal("a refused install must yield nothing to install through")
	}
}

// TestChartErrorLatchSurvivesAReattachThatNeverResolves is the sequence the
// guard could not see. The refusal used to be cleared at the top of
// AddCtxControllerHandlers, before the kubeconfig and Kubernetes-client steps
// that return early - and those returns deliberately leave the previously
// attached operator handler in place so observation survives. A connection that
// first attached during a chart-repository outage, and whose credentials then
// stopped working, therefore ended up with a handler still carrying the raw
// unresolved chart version and no refusal against it, and the very next
// DeployUndeployedOperators in the FSM chain would hand that version to Helm.
//
// The refusal now lifts only where a resolution actually succeeded.
func TestChartErrorLatchSurvivesAReattachThatNeverResolves(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)

	indexErr := errors.New("repository unreachable")
	mch.chartVersions = func(string, string) ([]string, error) { return nil, indexErr }
	mch.AddCtxControllerHandlers(reachableK8sContext())
	if mch.GetOperatorChartError() == nil {
		t.Fatal("an unreadable index must withhold installation")
	}
	if mch.ctxControllerHandlers[MesheryOperator] == nil {
		t.Fatal("the operator handler must stay attached for observation")
	}

	// The repository is readable again, but the stored kubeconfig no longer
	// yields a Kubernetes client, so this run returns before it ever resolves a
	// chart version - leaving the previous handler, and its unresolved version,
	// in place.
	requireUnresolvableK8sContext(t)
	mch.chartVersions = func(string, string) ([]string, error) { return publishedCharts, nil }
	mch.AddCtxControllerHandlers(unresolvableK8sContext())

	if mch.GetOperatorChartError() == nil {
		t.Fatal("a run that never reached resolution cleared the refusal, leaving a stale handler unguarded")
	}
	if mch.GetOperatorError() == nil {
		t.Fatal("the Kubernetes client failure must still be recorded")
	}

	// Observed through a stub so that a regression refuses into a counter rather
	// than into meshkit's Helm client.
	stub := &stubController{status: controllers.NotDeployed}
	mch.ctxControllerHandlers[MesheryOperator] = stub
	mch.ctxOperatorStatus = controllers.NotDeployed

	mch.DeployUndeployedOperators(NewOperatorTracker(false), testContextID)

	if stub.deploys != 0 {
		t.Fatalf("Deploy was called %d times with an unresolved chart version still attached", stub.deploys)
	}
}

// TestReconcileLeavesTheLatchStandingWhenReattachFails pins the reconcile
// counterpart of the latch fix, and only that: its re-attach fails at the
// Kubernetes client, so the reconcile returns at the operator-error branch
// before the install site is reached. What it proves is that a re-attach which
// never resolved a chart version neither installs nor lifts the refusal, and
// that the pre-reconcile chart version is restored so the next attempt retries.
//
// It deliberately does NOT cover the reconcile's wiring to the shared install
// guard - nothing here reaches that branch. TestReconcileInstallGoesThroughTheSharedGuard
// is what covers it.
func TestReconcileLeavesTheLatchStandingWhenReattachFails(t *testing.T) {
	requireUnresolvableK8sContext(t)

	mch := newTestControllersHelper(t)
	mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
	mch.SetControllersConfig(operatorVersionConfig("", connections.MeshsyncDeploymentModeOperator))
	mch.attachedOperatorChartVersion = "v1.0.51"

	stub := &stubController{status: controllers.NotDeployed}
	mch.ctxControllerHandlers = map[MesheryController]controllers.IMesheryController{MesheryOperator: stub}
	mch.setOperatorChartError(ErrNoOperatorChartPublished(OperatorChartName, testChartRepo))

	_, redeployed, err := mch.ReconcileOperatorChartVersion(unresolvableK8sContext(), NewOperatorTracker(false))
	if redeployed {
		t.Fatal("the re-attach failed, so nothing can have been redeployed")
	}
	if err == nil {
		t.Fatal("expected the failed re-attach to be reported")
	}
	if code := meshkiterrors.GetCode(err); code != ErrReconcileOperatorChartVersionCode {
		t.Fatalf("error code = %q, want %q (from %v)", code, ErrReconcileOperatorChartVersionCode, err)
	}
	if stub.deploys != 0 {
		t.Fatalf("Deploy was called %d times after a re-attach that never resolved a chart version", stub.deploys)
	}
	if mch.GetOperatorChartError() == nil {
		t.Fatal("a reconcile that never resolved must leave the refusal standing")
	}
	if mch.attachedOperatorChartVersion != "v1.0.51" {
		t.Fatalf("attached chart version = %q, want the pre-reconcile value restored so the refusal is retried",
			mch.attachedOperatorChartVersion)
	}
}

// TestReconcileInstallGoesThroughTheSharedGuard covers the reconcile's install
// site itself: that it asks operatorInstallTarget for the handler rather than
// reading the handler map, so a standing chart refusal stops it.
//
// The state that distinguishes the two - a re-attach that succeeded and yet left
// a refusal standing - is one AddCtxControllerHandlers cannot produce, because it
// clears the refusal exactly where it resolves. That is why the re-attach is
// injected here: without the seam this branch is unreachable, nothing would
// notice an install site that stopped consulting the guard, and that is precisely
// how the gap survived the first time.
//
// The injected re-attach installs a stub, so removing the guard call or
// replacing it with a direct map read makes stub.deploys 1 and err nil - this
// test fails by construction, and still never reaches meshkit's Helm client.
func TestReconcileInstallGoesThroughTheSharedGuard(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
	mch.SetControllersConfig(operatorVersionConfig("", connections.MeshsyncDeploymentModeOperator))
	mch.attachedOperatorChartVersion = "v1.0.51"

	latched := ErrNoOperatorChartPublished(OperatorChartName, testChartRepo)
	stub := &stubController{status: controllers.NotDeployed}
	reattached := 0
	mch.reattachControllerHandlers = func(K8sContext) {
		reattached++
		mch.setOperatorError(nil)
		mch.setOperatorChartError(latched)
		mch.ctxControllerHandlers = map[MesheryController]controllers.IMesheryController{MesheryOperator: stub}
		mch.attachedOperatorChartVersion = bootChartVersion
	}

	// The context is inert: the injected re-attach never builds a client from it.
	_, redeployed, err := mch.ReconcileOperatorChartVersion(
		K8sContext{ID: testContextID},
		NewOperatorTracker(false),
	)

	if reattached != 1 {
		t.Fatalf("the re-attach ran %d times, want 1: this test is not exercising the path it claims to", reattached)
	}
	if stub.deploys != 0 {
		t.Fatalf("Deploy was called %d times with a chart refusal standing: the install site is not going through operatorInstallTarget", stub.deploys)
	}
	if redeployed {
		t.Fatal("installation was refused, so nothing can have been redeployed")
	}
	if err == nil {
		t.Fatal("expected the refusal to be reported")
	}
	if code := meshkiterrors.GetCode(err); code != ErrReconcileOperatorChartVersionCode {
		t.Fatalf("error code = %q, want %q (from %v)", code, ErrReconcileOperatorChartVersionCode, err)
	}
	if !strings.Contains(err.Error(), "advertises no released") {
		t.Fatalf("the refusal must carry the latched chart failure, got %v", err)
	}
	if mch.attachedOperatorChartVersion != "v1.0.51" {
		t.Fatalf("attached chart version = %q, want the pre-reconcile value restored so the refusal is retried",
			mch.attachedOperatorChartVersion)
	}
}

// TestPrereleaseChartsAreNeverSelectedAutomatically covers the whole automatic
// path. meshery.io has carried prerelease operator charts (v0.6.0-rc.6,
// v0.6.0-rc.5l, v0.5.0-rc-5g), and by semver a release candidate outranks every
// stable chart below it - so an rc published ahead of its release would
// otherwise be what every server whose own chart is not yet published installs
// onto a production cluster.
func TestPrereleaseChartsAreNeverSelectedAutomatically(t *testing.T) {
	t.Run("the newest-published fallback skips a prerelease", func(t *testing.T) {
		published := []string{"v1.0.66-rc.1", "v1.0.65", "v1.0.51"}
		version, reason := mustResolve(t, published, "v1.0.70", OperatorChartVersionDerived)
		if version != "v1.0.65" {
			t.Fatalf("chart version = %q, want the newest *stable* chart v1.0.65", version)
		}
		if reason == "" {
			t.Fatal("expected the fallback to be explained")
		}
	})

	t.Run("an unpinned server release skips a prerelease", func(t *testing.T) {
		published := []string{"v1.0.66-rc.1", "v1.0.65", "v1.0.51"}
		version, _ := mustResolve(t, published, "", OperatorChartVersionDerived)
		if version != "v1.0.65" {
			t.Fatalf("chart version = %q, want the newest *stable* chart v1.0.65", version)
		}
	})

	t.Run("the floor raise skips a prerelease sitting at the boundary", func(t *testing.T) {
		// v1.0.52-rc.1 is the oldest published chart at or above the floor, and
		// selecting it would land a release candidate on a cluster whose only
		// fault was running an old server.
		published := []string{"v1.0.53", "v1.0.52-rc.1", "v1.0.40"}
		version, reason := mustResolve(t, published, "v1.0.40", OperatorChartVersionDerived)
		if version != "v1.0.53" {
			t.Fatalf("chart version = %q, want the oldest *stable* chart at or above the floor, v1.0.53", version)
		}
		if reason == "" {
			t.Fatal("expected the floor raise to be explained")
		}
	})

	t.Run("a repository publishing only prereleases fails rather than selecting one", func(t *testing.T) {
		_, _, err := ResolveOperatorChartVersion(testChartRepo, []string{"v1.0.66-rc.1", "v1.0.65-rc.2"}, "v1.0.70", OperatorChartVersionDerived)
		if err == nil {
			t.Fatal("expected resolution to fail rather than select a release candidate nobody asked for")
		}
		if code := meshkiterrors.GetCode(err); code != ErrNoOperatorChartPublishedCode {
			t.Fatalf("error code = %q, want %q (from %v)", code, ErrNoOperatorChartPublishedCode, err)
		}
	})
}

// TestPrereleaseChartsRemainReachableWhenNamed is the other side of the rule:
// exclusion applies to what Meshery *chooses*, never to what it is *told*. A
// prerelease pinned by name is honored exactly, and a prerelease server release
// that matches a published prerelease chart is not a choice either - it names
// one release, and it is the one the server was built alongside.
func TestPrereleaseChartsRemainReachableWhenNamed(t *testing.T) {
	published := []string{"v1.0.66-rc.1", "v1.0.65", "v1.0.51"}

	t.Run("an explicit prerelease pin is honored", func(t *testing.T) {
		version, reason := mustResolve(t, published, "v1.0.66-rc.1", OperatorChartVersionRequested)
		if version != "v1.0.66-rc.1" {
			t.Fatalf("chart version = %q, want the explicitly requested v1.0.66-rc.1", version)
		}
		if reason != "" {
			t.Fatalf("an honored explicit request is not a substitution, got reason %q", reason)
		}
	})

	t.Run("a prerelease server release matching a published chart is used verbatim", func(t *testing.T) {
		version, reason := mustResolve(t, published, "v1.0.66-rc.1", OperatorChartVersionDerived)
		if version != "v1.0.66-rc.1" {
			t.Fatalf("chart version = %q, want the matching published chart v1.0.66-rc.1", version)
		}
		if reason != "" {
			t.Fatalf("naming a published release exactly is not a substitution, got reason %q", reason)
		}
	})

	t.Run("an unusable explicit pin is pointed at a stable version", func(t *testing.T) {
		_, _, err := ResolveOperatorChartVersion(testChartRepo, published, "v9.9.9", OperatorChartVersionRequested)
		if err == nil {
			t.Fatal("expected an unpublished explicit pin to fail")
		}
		if !strings.Contains(err.Error(), "v1.0.65") {
			t.Fatalf("the remedy must name a stable chart, got %v", err)
		}
		if strings.Contains(err.Error(), "v1.0.66-rc.1") {
			t.Fatalf("the remedy must not recommend a release candidate, got %v", err)
		}
	})
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
