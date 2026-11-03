package models

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/core"
	"helm.sh/helm/v3/pkg/chart"
	rspb "helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

// playgroundInstalledChart is the chart the Meshery Playground's operator
// release was found on (meshery/meshery#22555): installed months before the
// floor existed, crash-looping ever since, and never touched by a server that
// carried the floor.
const playgroundInstalledChart = "v1.0.0"

// fakeOperatorRelease is an operatorRelease that serves a fixed installed chart
// and records every upgrade it is asked for, so the floor reconcile can be
// driven without a cluster or a chart repository.
type fakeOperatorRelease struct {
	installed  string
	found      bool
	readErr    error
	upgradeErr error

	reads    int
	upgrades []controllers.OperatorDeploymentConfig
}

func (f *fakeOperatorRelease) DeployedChartVersion() (string, bool, error) {
	f.reads++
	return f.installed, f.found, f.readErr
}

func (f *fakeOperatorRelease) Upgrade(depConfig controllers.OperatorDeploymentConfig) error {
	f.upgrades = append(f.upgrades, depConfig)
	return f.upgradeErr
}

// floorTestConnection is the connection every floor test reconciles. It names
// a connection so the events it produces are the persistable kind.
var floorTestConnection = K8sContext{ID: testContextID, Name: "playground", ConnectionID: "6a1c7e7e-6c4d-4b8e-9d1f-2f7d0f1c9a11"}

// newFloorTestHelper builds a helper in the state the connect path leaves it in
// just before DeployUndeployedOperators: operator mode, handlers attached, the
// operator's status read as status, and the chart version resolved exactly as
// AddCtxControllerHandlers resolves it - so a test's target is whatever real
// resolution yields for its configuration, not a value the test asserted into
// place. The returned channel receives every event the helper emits.
func newFloorTestHelper(t *testing.T, status controllers.MesheryControllerStatus, release *fakeOperatorRelease) (*MesheryControllersHelper, *stubController, <-chan interface{}) {
	t.Helper()
	mch := newTestControllersHelper(t)
	mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)

	systemID := core.Uuid(uuid.Must(uuid.NewV4()))
	mch.systemID = &systemID
	mch.eventBroadcaster = NewBroadcaster("operator-chart-floor-test")
	eventsCh, unsubscribe := mch.eventBroadcaster.Subscribe(systemID)
	t.Cleanup(unsubscribe)

	stub := &stubController{status: status}
	mch.ctxControllerHandlers = map[MesheryController]controllers.IMesheryController{MesheryOperator: stub}
	mch.ctxOperatorStatus = status
	mch.operatorRelease = release
	return mch, stub, eventsCh
}

// resolveAsAttached runs the resolution AddCtxControllerHandlers runs and
// records its outcome the way that function does.
func resolveAsAttached(t *testing.T, mch *MesheryControllersHelper) {
	t.Helper()
	depConfig, _, err := mch.pinnedOperatorDeploymentConfig()
	if err != nil {
		mch.setOperatorChartError(err)
		mch.attachedOperatorChartVersion = ""
		return
	}
	mch.setOperatorChartError(nil)
	mch.attachedOperatorChartVersion = depConfig.MesheryReleaseVersion
}

func nextEvent(t *testing.T, ch <-chan interface{}) *events.Event {
	t.Helper()
	select {
	case got := <-ch:
		event, ok := got.(*events.Event)
		if !ok {
			t.Fatalf("broadcast carried %T, want *events.Event", got)
		}
		return event
	case <-time.After(2 * time.Second):
		t.Fatal("expected an event, none was emitted")
		return nil
	}
}

// requireNoEvent fails if anything is emitted. Publication is asynchronous, so
// absence can only be shown by waiting.
func requireNoEvent(t *testing.T, ch <-chan interface{}) {
	t.Helper()
	select {
	case got := <-ch:
		if event, ok := got.(*events.Event); ok {
			t.Fatalf("expected no event, got %q (%v)", event.Description, event.Metadata["error"])
		}
		t.Fatalf("expected no event, got %T", got)
	case <-time.After(150 * time.Millisecond):
	}
}

// requireEventCode asserts the event carries the MeshKit error with code as a
// structured error - the shape the notification center renders - and names the
// connection, which is what makes a background reconcile's event persistable.
// It returns the error's long description, where the versions involved travel.
func requireEventCode(t *testing.T, event *events.Event, severity events.EventSeverity, code string) string {
	t.Helper()
	if event.Severity != severity {
		t.Fatalf("event %q severity = %q, want %q", event.Description, event.Severity, severity)
	}
	meshkitErr, ok := event.Metadata["error"].(*meshkiterrors.Error)
	if !ok {
		t.Fatalf("event %q error metadata is %T, want a structured *errors.Error", event.Description, event.Metadata["error"])
	}
	if meshkitErr.Code != code {
		t.Fatalf("event %q error code = %q, want %q", event.Description, meshkitErr.Code, code)
	}
	if got := event.Metadata["connectionID"]; got != floorTestConnection.ConnectionID {
		t.Fatalf("event %q connectionID = %v, want %q so it is persisted", event.Description, got, floorTestConnection.ConnectionID)
	}
	if event.ActedUpon != uuid.FromStringOrNil(floorTestConnection.ConnectionID) {
		t.Fatalf("event %q acts upon %v, want the connection", event.Description, event.ActedUpon)
	}
	return strings.Join(meshkitErr.LongDescription, " ")
}

func TestDecideOperatorChartFloor(t *testing.T) {
	tests := []struct {
		name      string
		installed string
		target    string
		source    OperatorChartVersionSource
		want      operatorChartFloorAction
	}{
		{"below floor, derived target at the floor", playgroundInstalledChart, MinimumOperatorChartVersion, OperatorChartVersionDerived, operatorChartFloorUpgrade},
		{"below floor, derived target above the floor", playgroundInstalledChart, "v1.0.74", OperatorChartVersionDerived, operatorChartFloorUpgrade},
		{"just below floor", "v1.0.50", "v1.0.74", OperatorChartVersionDerived, operatorChartFloorUpgrade},
		{"at the floor is untouched", MinimumOperatorChartVersion, "v1.0.74", OperatorChartVersionDerived, operatorChartFloorSatisfied},
		{"above the floor is untouched", "v1.0.64", "v1.0.74", OperatorChartVersionDerived, operatorChartFloorSatisfied},
		{"above the target is never downgraded", "v1.0.80", "v1.0.74", OperatorChartVersionDerived, operatorChartFloorSatisfied},
		{"unprefixed spelling at the floor is untouched", "1.0.51", "v1.0.74", OperatorChartVersionDerived, operatorChartFloorSatisfied},
		{"pin below the floor is never overridden", playgroundInstalledChart, "v1.0.40", OperatorChartVersionRequested, operatorChartFloorRefusePinned},
		{"pin equal to the below-floor install is never overridden", "v1.0.40", "v1.0.40", OperatorChartVersionRequested, operatorChartFloorRefusePinned},
		{"pin at or above the floor is what the release is upgraded to", playgroundInstalledChart, "v1.0.62", OperatorChartVersionRequested, operatorChartFloorUpgrade},
		{"derived target still below the floor has nothing to upgrade to", playgroundInstalledChart, "v1.0.40", OperatorChartVersionDerived, operatorChartFloorRefuseUnreachable},
		{"no resolved target has nothing to upgrade to", playgroundInstalledChart, "", OperatorChartVersionDerived, operatorChartFloorRefuseUnreachable},
		{"an unparseable install names no chart verified to run", "stable-latest", "v1.0.74", OperatorChartVersionDerived, operatorChartFloorUpgrade},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decideOperatorChartFloor(tt.installed, tt.target, tt.source); got != tt.want {
				t.Fatalf("decideOperatorChartFloor(%q, %q, %v) = %v, want %v", tt.installed, tt.target, tt.source, got, tt.want)
			}
		})
	}
}

// TestBelowFloorDerivedReleaseIsUpgraded is meshery/meshery#22555: an operator
// release installed on a chart below the floor, which a current server must
// upgrade on connect rather than leave crash-looping. Deploying is what that
// crash loop reports (its rollout never completes) and meshkit's Deploy refuses
// on it, so the upgrade must not go through Deploy.
func TestBelowFloorDerivedReleaseIsUpgraded(t *testing.T) {
	for _, status := range []controllers.MesheryControllerStatus{controllers.Deploying, controllers.Deployed} {
		t.Run(status.String(), func(t *testing.T) {
			release := &fakeOperatorRelease{installed: playgroundInstalledChart, found: true}
			mch, stub, eventsCh := newFloorTestHelper(t, status, release)
			resolveAsAttached(t, mch)

			mch.DeployUndeployedOperators(NewOperatorTracker(false), floorTestConnection)

			if len(release.upgrades) != 1 {
				t.Fatalf("release upgraded %d times, want 1", len(release.upgrades))
			}
			got := release.upgrades[0]
			if got.MesheryReleaseVersion != bootChartVersion {
				t.Fatalf("upgraded to chart %q, want the derived, published chart %q", got.MesheryReleaseVersion, bootChartVersion)
			}
			if got.HelmChartRepo != testChartRepo {
				t.Fatalf("upgraded from repository %q, want %q", got.HelmChartRepo, testChartRepo)
			}
			if got.GetHelmOverrides == nil {
				t.Fatal("the upgrade must carry the same Helm overrides an install does")
			}
			if stub.deploys != 0 {
				t.Fatalf("Deploy was called %d times; it refuses a Deploying operator and must not be the upgrade path", stub.deploys)
			}
			if err := mch.GetOperatorError(); err != nil {
				t.Fatalf("operator error = %v after a successful upgrade, want nil", err)
			}

			event := nextEvent(t, eventsCh)
			details := requireEventCode(t, event, events.Warning, ErrOperatorChartBelowFloorUpgradedCode)
			for _, want := range []string{playgroundInstalledChart, bootChartVersion, MinimumOperatorChartVersion} {
				if !strings.Contains(details, want) {
					t.Fatalf("upgrade notice %q does not name %s", details, want)
				}
			}
			if event.Metadata["installedOperatorChartVersion"] != playgroundInstalledChart || event.Metadata["operatorChartVersion"] != bootChartVersion {
				t.Fatalf("event metadata = %v, want the installed and upgraded chart versions", event.Metadata)
			}
		})
	}
}

// TestBelowFloorReleaseIsUpgradedToTheFlooredDerivedChart: an old server whose
// own release names a chart below the floor resolves to the oldest published
// chart at or above it, and that - not the server's own below-floor release -
// is what an installed below-floor release is upgraded to.
func TestBelowFloorReleaseIsUpgradedToTheFlooredDerivedChart(t *testing.T) {
	release := &fakeOperatorRelease{installed: playgroundInstalledChart, found: true}
	mch, _, eventsCh := newFloorTestHelper(t, controllers.Deploying, release)
	mch.oprDepConfig.MesheryReleaseVersion = "v1.0.40"
	resolveAsAttached(t, mch)

	mch.DeployUndeployedOperators(NewOperatorTracker(false), floorTestConnection)

	if len(release.upgrades) != 1 || release.upgrades[0].MesheryReleaseVersion != MinimumOperatorChartVersion {
		t.Fatalf("upgrades = %v, want one to the floor %s", release.upgrades, MinimumOperatorChartVersion)
	}
	requireEventCode(t, nextEvent(t, eventsCh), events.Warning, ErrOperatorChartBelowFloorUpgradedCode)
}

// TestReleaseAtOrAboveFloorIsUntouched: the floor only lifts releases out of
// the range known not to run. Anything at or above it is left exactly as it is
// - no upgrade, no event, no recorded error - even when a newer chart is
// available, and even while installation is withheld for an unrelated reason.
func TestReleaseAtOrAboveFloorIsUntouched(t *testing.T) {
	for _, installed := range []string{MinimumOperatorChartVersion, "1.0.51", "v1.0.62", "v1.0.80"} {
		t.Run(installed, func(t *testing.T) {
			release := &fakeOperatorRelease{installed: installed, found: true}
			mch, stub, eventsCh := newFloorTestHelper(t, controllers.Deployed, release)
			resolveAsAttached(t, mch)
			// A latched resolution failure concerns installing, not
			// observing; a healthy release must not be blocked or reported
			// on account of it.
			mch.setOperatorChartError(ErrNoOperatorChartPublished(OperatorChartName, testChartRepo))

			mch.DeployUndeployedOperators(NewOperatorTracker(false), floorTestConnection)

			if release.reads != 1 {
				t.Fatalf("release read %d times, want 1", release.reads)
			}
			if len(release.upgrades) != 0 || stub.deploys != 0 {
				t.Fatalf("a release at %s was changed: %d upgrades, %d deploys", installed, len(release.upgrades), stub.deploys)
			}
			if err := mch.GetOperatorError(); err != nil {
				t.Fatalf("operator error = %v, want nil for a release at or above the floor", err)
			}
			requireNoEvent(t, eventsCh)
		})
	}
}

// TestPinnedBelowFloorIsNotOverriddenAndReported: a deliberate operator.version
// is never overridden, even when both it and the installed release are below
// the floor - but leaving it is not silent.
func TestPinnedBelowFloorIsNotOverriddenAndReported(t *testing.T) {
	release := &fakeOperatorRelease{installed: playgroundInstalledChart, found: true}
	mch, stub, eventsCh := newFloorTestHelper(t, controllers.Deploying, release)
	mch.SetControllersConfig(operatorVersionConfig("v1.0.40", connections.MeshsyncDeploymentModeOperator))
	resolveAsAttached(t, mch)
	if mch.attachedOperatorChartVersion != "v1.0.40" {
		t.Fatalf("precondition: the pin resolved to %q, want it honored as v1.0.40", mch.attachedOperatorChartVersion)
	}

	mch.DeployUndeployedOperators(NewOperatorTracker(false), floorTestConnection)

	if len(release.upgrades) != 0 || stub.deploys != 0 {
		t.Fatalf("a pinned release was changed: %d upgrades, %d deploys", len(release.upgrades), stub.deploys)
	}
	details := requireEventCode(t, nextEvent(t, eventsCh), events.Error, ErrOperatorChartPinnedBelowFloorCode)
	for _, want := range []string{playgroundInstalledChart, "v1.0.40", MinimumOperatorChartVersion} {
		if !strings.Contains(details, want) {
			t.Fatalf("refusal %q does not name %s", details, want)
		}
	}
	if code := meshkiterrors.GetCode(mch.GetOperatorError()); code != ErrOperatorChartPinnedBelowFloorCode {
		t.Fatalf("operator error code = %q, want the refusal recorded for diagnostics", code)
	}
}

// TestPinnedAtOrAboveFloorIsWhatABelowFloorReleaseIsUpgradedTo: installing the
// chart a user pinned is honoring the pin, not overriding it.
func TestPinnedAtOrAboveFloorIsWhatABelowFloorReleaseIsUpgradedTo(t *testing.T) {
	release := &fakeOperatorRelease{installed: playgroundInstalledChart, found: true}
	mch, _, eventsCh := newFloorTestHelper(t, controllers.Deploying, release)
	mch.SetControllersConfig(operatorVersionConfig("v1.0.62", connections.MeshsyncDeploymentModeOperator))
	resolveAsAttached(t, mch)

	mch.DeployUndeployedOperators(NewOperatorTracker(false), floorTestConnection)

	if len(release.upgrades) != 1 || release.upgrades[0].MesheryReleaseVersion != "v1.0.62" {
		t.Fatalf("upgrades = %v, want one to the pinned v1.0.62", release.upgrades)
	}
	requireEventCode(t, nextEvent(t, eventsCh), events.Warning, ErrOperatorChartBelowFloorUpgradedCode)
}

// TestFloorUnreachableIsReported: when the repository publishes nothing at or
// above the floor, the derived version is below it too, and upgrading to it
// would trade one chart that cannot run for another.
func TestFloorUnreachableIsReported(t *testing.T) {
	release := &fakeOperatorRelease{installed: playgroundInstalledChart, found: true}
	mch, _, eventsCh := newFloorTestHelper(t, controllers.Deploying, release)
	mch.chartVersions = func(string, string) ([]string, error) { return []string{"v1.0.40", "v0.7.9"}, nil }
	resolveAsAttached(t, mch)

	mch.DeployUndeployedOperators(NewOperatorTracker(false), floorTestConnection)

	if len(release.upgrades) != 0 {
		t.Fatalf("release upgraded %d times to a chart below the floor", len(release.upgrades))
	}
	details := requireEventCode(t, nextEvent(t, eventsCh), events.Error, ErrOperatorChartFloorUnreachableCode)
	if !strings.Contains(details, testChartRepo) {
		t.Fatalf("refusal %q does not name the repository that was read", details)
	}
	if code := meshkiterrors.GetCode(mch.GetOperatorError()); code != ErrOperatorChartFloorUnreachableCode {
		t.Fatalf("operator error code = %q, want the refusal recorded", code)
	}
}

// TestReleaseReadErrorIsSurfaced: failing to read the release is a failure to
// look, not evidence of a healthy release, and is reported as such.
func TestReleaseReadErrorIsSurfaced(t *testing.T) {
	readErr := ErrReadOperatorRelease(`secrets is forbidden: User "meshery" cannot list resource "secrets" in the namespace "meshery"`)
	release := &fakeOperatorRelease{readErr: readErr}
	mch, _, eventsCh := newFloorTestHelper(t, controllers.Deploying, release)
	resolveAsAttached(t, mch)

	mch.DeployUndeployedOperators(NewOperatorTracker(false), floorTestConnection)

	if len(release.upgrades) != 0 {
		t.Fatalf("release upgraded %d times without its chart version being read", len(release.upgrades))
	}
	requireEventCode(t, nextEvent(t, eventsCh), events.Error, ErrReadOperatorReleaseCode)
	if !stderrors.Is(mch.GetOperatorError(), readErr) {
		t.Fatalf("operator error = %v, want the read failure recorded", mch.GetOperatorError())
	}
}

// TestUpgradeErrorIsSurfacedAndRetried: a failed upgrade is reported with its
// cause, recorded for diagnostics, and attempted again on the next connect,
// because the deployed revision is still below the floor.
func TestUpgradeErrorIsSurfacedAndRetried(t *testing.T) {
	helmErr := stderrors.New(`pre-upgrade hooks failed: 1 error occurred: job meshery-operator-crd-update failed: BackoffLimitExceeded`)
	release := &fakeOperatorRelease{installed: playgroundInstalledChart, found: true, upgradeErr: helmErr}
	mch, _, eventsCh := newFloorTestHelper(t, controllers.Deploying, release)
	resolveAsAttached(t, mch)

	mch.DeployUndeployedOperators(NewOperatorTracker(false), floorTestConnection)

	details := requireEventCode(t, nextEvent(t, eventsCh), events.Error, ErrUpgradeOperatorReleaseCode)
	if !strings.Contains(details, helmErr.Error()) {
		t.Fatalf("failure %q does not carry the Helm cause", details)
	}
	if code := meshkiterrors.GetCode(mch.GetOperatorError()); code != ErrUpgradeOperatorReleaseCode {
		t.Fatalf("operator error code = %q, want the upgrade failure recorded", code)
	}

	mch.DeployUndeployedOperators(NewOperatorTracker(false), floorTestConnection)
	if len(release.upgrades) != 2 {
		t.Fatalf("release upgrade attempted %d times across two connects, want a retry", len(release.upgrades))
	}
	requireEventCode(t, nextEvent(t, eventsCh), events.Error, ErrUpgradeOperatorReleaseCode)
}

// TestLatchedChartErrorWithholdsTheUpgrade: a below-floor release with no
// resolvable chart version takes the same refusal every other install path
// does. The resolution failure was already emitted at attach, so it is
// recorded, not emitted again.
func TestLatchedChartErrorWithholdsTheUpgrade(t *testing.T) {
	release := &fakeOperatorRelease{installed: playgroundInstalledChart, found: true}
	mch, _, eventsCh := newFloorTestHelper(t, controllers.Deploying, release)
	chartErr := ErrNoOperatorChartPublished(OperatorChartName, testChartRepo)
	mch.setOperatorChartError(chartErr)

	mch.DeployUndeployedOperators(NewOperatorTracker(false), floorTestConnection)

	if len(release.upgrades) != 0 {
		t.Fatalf("release upgraded %d times with no resolvable chart version", len(release.upgrades))
	}
	if !stderrors.Is(mch.GetOperatorError(), chartErr) {
		t.Fatalf("operator error = %v, want the chart resolution failure", mch.GetOperatorError())
	}
	requireNoEvent(t, eventsCh)
}

// TestNoReleaseIsOrdinary: an operator that is not the meshery-operator Helm
// release - a subchart of the Meshery chart, or applied without Helm - has no
// chart version to judge, and is not an error.
func TestNoReleaseIsOrdinary(t *testing.T) {
	release := &fakeOperatorRelease{found: false}
	mch, _, eventsCh := newFloorTestHelper(t, controllers.Deployed, release)
	resolveAsAttached(t, mch)

	mch.DeployUndeployedOperators(NewOperatorTracker(false), floorTestConnection)

	if len(release.upgrades) != 0 {
		t.Fatalf("upgraded a release that does not exist")
	}
	if err := mch.GetOperatorError(); err != nil {
		t.Fatalf("operator error = %v, want nil when there is no release", err)
	}
	requireNoEvent(t, eventsCh)
}

// TestFloorIsNotCheckedWhereTheOperatorIsNotManaged: nothing reads the release,
// let alone upgrades it, while the operator is disabled server-wide, was
// explicitly undeployed, or the connection runs MeshSync embedded.
func TestFloorIsNotCheckedWhereTheOperatorIsNotManaged(t *testing.T) {
	cases := map[string]func(*MesheryControllersHelper) *OperatorTracker{
		"operator disabled": func(*MesheryControllersHelper) *OperatorTracker { return NewOperatorTracker(true) },
		"undeployed for this context": func(*MesheryControllersHelper) *OperatorTracker {
			ot := NewOperatorTracker(false)
			ot.Undeployed(testContextID, true)
			return ot
		},
		"embedded mode": func(mch *MesheryControllersHelper) *OperatorTracker {
			mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeEmbedded)
			return NewOperatorTracker(false)
		},
		"status unknown": func(mch *MesheryControllersHelper) *OperatorTracker {
			mch.ctxOperatorStatus = controllers.Unknown
			return NewOperatorTracker(false)
		},
	}
	for name, arrange := range cases {
		t.Run(name, func(t *testing.T) {
			release := &fakeOperatorRelease{installed: playgroundInstalledChart, found: true}
			mch, _, eventsCh := newFloorTestHelper(t, controllers.Deploying, release)
			resolveAsAttached(t, mch)
			ot := arrange(mch)

			mch.DeployUndeployedOperators(ot, floorTestConnection)

			if release.reads != 0 || len(release.upgrades) != 0 {
				t.Fatalf("release read %d times and upgraded %d times on an unmanaged operator", release.reads, len(release.upgrades))
			}
			requireNoEvent(t, eventsCh)
		})
	}
}

// TestHandlersAndReleaseAttachAndDetachTogether: operatorRelease is reachable
// exactly when the handlers are, so a torn-down connection cannot be upgraded
// through a stale release, and an attached one always has a release to read.
func TestHandlersAndReleaseAttachAndDetachTogether(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.AddCtxControllerHandlers(reachableK8sContext())
	if mch.attachedOperatorHandler() == nil || mch.operatorRelease == nil {
		t.Fatalf("handler attached = %v, release attached = %v; want both", mch.attachedOperatorHandler() != nil, mch.operatorRelease != nil)
	}
	if _, ok := mch.operatorRelease.(*helmOperatorRelease); !ok {
		t.Fatalf("attached release is %T, want the Helm-backed release", mch.operatorRelease)
	}

	mch.RemoveCtxControllerHandler(context.Background(), testContextID)
	if mch.operatorRelease != nil {
		t.Fatal("the release must be detached with the handlers")
	}
}

// helmReleaseRecord stores a release revision the way Helm's Secret driver does,
// so the reader is exercised against Helm's real storage format.
func helmReleaseRecord(t *testing.T, store *driver.Secrets, revision int, status rspb.Status, chartVersion string) {
	t.Helper()
	rel := &rspb.Release{
		Name:      OperatorReleaseName,
		Namespace: OperatorReleaseNamespace,
		Version:   revision,
		Info:      &rspb.Info{Status: status},
	}
	if chartVersion != "-" {
		rel.Chart = &chart.Chart{Metadata: &chart.Metadata{Name: OperatorChartName, Version: chartVersion, APIVersion: chart.APIVersionV2}}
	}
	key := "sh.helm.release.v1." + OperatorReleaseName + ".v" + strconv.Itoa(revision)
	if err := store.Create(key, rel); err != nil {
		t.Fatalf("storing release revision %d: %v", revision, err)
	}
}

func TestHelmOperatorReleaseReadsTheDeployedRevision(t *testing.T) {
	tests := []struct {
		name      string
		revisions []struct {
			status  rspb.Status
			version string
		}
		wantVersion string
		wantFound   bool
	}{
		{name: "no release at all", wantFound: false},
		{
			name: "one deployed revision",
			revisions: []struct {
				status  rspb.Status
				version string
			}{{rspb.StatusDeployed, playgroundInstalledChart}},
			wantVersion: playgroundInstalledChart, wantFound: true,
		},
		{
			name: "a failed upgrade does not count as installed",
			revisions: []struct {
				status  rspb.Status
				version string
			}{{rspb.StatusDeployed, playgroundInstalledChart}, {rspb.StatusFailed, "v1.0.74"}},
			wantVersion: playgroundInstalledChart, wantFound: true,
		},
		{
			name: "a superseded revision does not count as installed",
			revisions: []struct {
				status  rspb.Status
				version string
			}{{rspb.StatusSuperseded, playgroundInstalledChart}, {rspb.StatusDeployed, "v1.0.74"}},
			wantVersion: "v1.0.74", wantFound: true,
		},
		{
			name: "only a failed install is no deployed release",
			revisions: []struct {
				status  rspb.Status
				version string
			}{{rspb.StatusFailed, "v1.0.74"}},
			wantFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secrets := fake.NewSimpleClientset().CoreV1().Secrets(OperatorReleaseNamespace)
			store := driver.NewSecrets(secrets)
			for i, rev := range tt.revisions {
				helmReleaseRecord(t, store, i+1, rev.status, rev.version)
			}
			release := &helmOperatorRelease{secrets: secrets}

			version, found, err := release.DeployedChartVersion()
			if err != nil {
				t.Fatalf("DeployedChartVersion() error = %v", err)
			}
			if found != tt.wantFound || version != tt.wantVersion {
				t.Fatalf("DeployedChartVersion() = (%q, %v), want (%q, %v)", version, found, tt.wantVersion, tt.wantFound)
			}
		})
	}
}

// TestHelmOperatorReleaseSurfacesReadFailures: a forbidden or failed list is
// not an absent release, and a record with no chart version is not a chart
// version. Both come back as the structured read error.
func TestHelmOperatorReleaseSurfacesReadFailures(t *testing.T) {
	t.Run("forbidden list", func(t *testing.T) {
		client := fake.NewSimpleClientset()
		client.PrependReactor("list", "secrets", func(k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, apierrors.NewForbidden(schema.GroupResource{Resource: "secrets"}, "", stderrors.New(`User "meshery" cannot list resource "secrets"`))
		})
		release := &helmOperatorRelease{secrets: client.CoreV1().Secrets(OperatorReleaseNamespace)}

		_, found, err := release.DeployedChartVersion()
		if found || meshkiterrors.GetCode(err) != ErrReadOperatorReleaseCode {
			t.Fatalf("DeployedChartVersion() = (found %v, %v), want the structured read error", found, err)
		}
		if !strings.Contains(meshkiterrors.GetLDescription(err), "cannot list") {
			t.Fatalf("read error %q does not carry the Kubernetes cause", meshkiterrors.GetLDescription(err))
		}
	})
	t.Run("record without a chart version", func(t *testing.T) {
		secrets := fake.NewSimpleClientset().CoreV1().Secrets(OperatorReleaseNamespace)
		helmReleaseRecord(t, driver.NewSecrets(secrets), 1, rspb.StatusDeployed, "-")
		release := &helmOperatorRelease{secrets: secrets}

		_, found, err := release.DeployedChartVersion()
		if found || meshkiterrors.GetCode(err) != ErrReadOperatorReleaseCode {
			t.Fatalf("DeployedChartVersion() = (found %v, %v), want the structured read error", found, err)
		}
	})
}

// TestHelmOperatorReleaseUpgradesTheInstalledRelease pins what reaches Helm: an
// explicit UPGRADE (never an install) of the release meshkit installs, from the
// connection's repository at the target chart, with the full override set an
// install uses - non-empty, so Helm does not carry the old chart's values over.
func TestHelmOperatorReleaseUpgradesTheInstalledRelease(t *testing.T) {
	var applied []mesherykube.ApplyHelmChartConfig
	release := &helmOperatorRelease{applyHelmChart: func(cfg mesherykube.ApplyHelmChartConfig) error {
		applied = append(applied, cfg)
		return nil
	}}
	overrides := map[string]interface{}{"fullnameOverride": "meshery-operator"}
	var deleteRequested []bool
	err := release.Upgrade(controllers.OperatorDeploymentConfig{
		HelmChartRepo:         testChartRepo,
		MesheryReleaseVersion: "v1.0.74",
		GetHelmOverrides: func(del bool) map[string]interface{} {
			deleteRequested = append(deleteRequested, del)
			return overrides
		},
	})
	if err != nil {
		t.Fatalf("Upgrade() error = %v", err)
	}
	if len(applied) != 1 {
		t.Fatalf("ApplyHelmChart called %d times, want 1", len(applied))
	}
	got := applied[0]
	if got.Action != mesherykube.UPGRADE {
		t.Fatalf("Helm action = %v, want UPGRADE", got.Action)
	}
	if got.Namespace != "meshery" || got.ReleaseName != "meshery-operator" {
		t.Fatalf("release = %s/%s, want meshery/meshery-operator", got.Namespace, got.ReleaseName)
	}
	if got.ChartLocation != (mesherykube.HelmChartLocation{Repository: testChartRepo, Chart: OperatorChartName, Version: "v1.0.74"}) {
		t.Fatalf("chart location = %+v", got.ChartLocation)
	}
	if len(deleteRequested) != 1 || deleteRequested[0] {
		t.Fatalf("overrides requested with delete = %v, want the install set (delete=false)", deleteRequested)
	}
	if len(got.OverrideValues) == 0 || got.OverrideValues["fullnameOverride"] != "meshery-operator" {
		t.Fatalf("override values = %v, want the install overrides", got.OverrideValues)
	}

	helmErr := stderrors.New("UPGRADE FAILED")
	release.applyHelmChart = func(mesherykube.ApplyHelmChartConfig) error { return helmErr }
	if err := release.Upgrade(controllers.OperatorDeploymentConfig{GetHelmOverrides: func(bool) map[string]interface{} { return overrides }}); !stderrors.Is(err, helmErr) {
		t.Fatalf("Upgrade() error = %v, want the Helm failure returned", err)
	}
}

// TestControllerEventsCarryStructuredMeshKitErrors: the notification center
// renders metadata.error as a MeshKit error (LongDescription, ProbableCause,
// SuggestedRemediation). Flattening it to a string dropped all three, and the
// error code with them.
func TestControllerEventsCarryStructuredMeshKitErrors(t *testing.T) {
	meshkitErr := ErrOperatorChartPinnedBelowFloor(playgroundInstalledChart, "v1.0.40")
	if got, ok := eventErrorMetadata(meshkitErr).(*meshkiterrors.Error); !ok || got.Code != ErrOperatorChartPinnedBelowFloorCode {
		t.Fatalf("eventErrorMetadata(meshkit error) = %#v, want the error itself", eventErrorMetadata(meshkitErr))
	}
	wrapped := ErrReconcileOperatorChartVersion(meshkitErr)
	if got, ok := eventErrorMetadata(wrapped).(*meshkiterrors.Error); !ok || got.Code != ErrReconcileOperatorChartVersionCode {
		t.Fatalf("eventErrorMetadata(outer meshkit error) = %#v, want the outer error", eventErrorMetadata(wrapped))
	}
	if got := eventErrorMetadata(stderrors.New("plain")); got != "plain" {
		t.Fatalf("eventErrorMetadata(plain error) = %#v, want its text", got)
	}

	// What is persisted is the JSON of the metadata, and the notification
	// center reads these keys from it.
	encoded, err := json.Marshal(map[string]any{"error": eventErrorMetadata(meshkitErr)})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, key := range []string{`"Code":"` + ErrOperatorChartPinnedBelowFloorCode + `"`, `"LongDescription"`, `"ProbableCause"`, `"SuggestedRemediation"`} {
		if !strings.Contains(string(encoded), key) {
			t.Fatalf("persisted error metadata %s lacks %s", encoded, key)
		}
	}
}
