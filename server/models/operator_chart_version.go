package models

import (
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
)

const (
	// OperatorChartName is the chart published in ChartRepo that carries the
	// Meshery Operator, plus the Broker and MeshSync custom resources it
	// reconciles as subcharts.
	OperatorChartName = "meshery-operator"

	// MinimumOperatorChartVersion is the oldest published meshery-operator
	// chart verified to deploy successfully onto a current cluster.
	//
	// The floor means exactly one thing: below it, the chart cannot deploy at
	// all. It is not a "newest is better" preference - raising a server past a
	// chart that works is what a server upgrade is for, not what this floor is
	// for - so it must sit at the oldest version that actually renders clean.
	//
	// What was verified, by rendering the published archives with `helm
	// template`:
	//
	//   - v1.0.51 renders with no `kube-rbac-proxy` container and pins
	//     ENABLE_WEBHOOKS="false" on the manager. So does every chart above it.
	//   - Every published chart checked below it - v1.0.40, and the contiguous
	//     run v1.0.41 through v1.0.50 (v1.0.47 was never published) - renders a
	//     `registry.k8s.io/kubebuilder/kube-rbac-proxy:v0.16.0` sidecar and sets
	//     no ENABLE_WEBHOOKS. (Much older charts, v0.8.180 among them, name that
	//     same sidecar under gcr.io instead; both spellings appear in user
	//     reports.) Affected clusters report the sidecar stuck in
	//     ImagePullBackOff so the Pod never becomes Ready. Why the pull fails is
	//     deliberately not claimed here: the gcr.io copy of v0.16.0 is gone (404,
	//     empty tag list) but the registry.k8s.io copy still resolves, so that
	//     half is an environment fact, not a chart fact. The chart facts are
	//     enough on their own - an unset ENABLE_WEBHOOKS means *enabled* in
	//     current operator images, so the manager crash-loops on
	//     `open /tmp/k8s-webhook-server/serving-certs/tls.crt: no such file or
	//     directory` regardless of the sidecar - and both the sidecar reference
	//     and the missing ENABLE_WEBHOOKS are permanent, being baked into an
	//     immutable published archive.
	//
	// Charts older than v1.0.40 were not rendered, so nothing here claims a
	// cause for them - only that they are below the oldest chart verified to
	// work. v1.0.51 carries operator image 1.0.1 rather than the current 1.0.4;
	// that image is functional, and pinning the floor higher to reach a newer
	// operator image would substitute charts that provably deploy.
	//
	// The floor matters because the chart version a server asks for is derived
	// from the *server's own release* (see NewOperatorDeploymentConfig): a
	// server released before that chart was published would otherwise keep
	// installing an unusable operator, on every cluster it is ever pointed at,
	// for the rest of its life. The floor therefore applies only to that
	// derived default - an `operator.version` someone set deliberately is still
	// honored, so pinning an older chart on purpose remains possible.
	MinimumOperatorChartVersion = "v1.0.51"
)

// OperatorChartVersionSource records who asked for a chart version, which is
// what decides the behavior when that version turns out not to be published.
type OperatorChartVersionSource int

const (
	// OperatorChartVersionDerived means nobody chose this version: it is the
	// Meshery Server release stamped in at boot, standing in for "the chart
	// that goes with this server". Chart publishing is decoupled from server
	// releases and trails them, so this value is a *guess*, and correcting it
	// to a published version beats failing the deployment over it.
	OperatorChartVersionDerived OperatorChartVersionSource = iota

	// OperatorChartVersionRequested means a layer of the controllers
	// configuration set `operator.version` explicitly. An explicit request is
	// never silently substituted: an unpublished pin fails loudly so the user
	// learns the pin is wrong instead of quietly running something else.
	OperatorChartVersionRequested
)

// isPinnedChartVersion reports whether v names one immutable chart release.
//
// Moving channel tags (stable-latest, edge-latest, latest), the "Not Set"
// placeholder viper hands back for an unstamped build, and the empty string all
// fail this test. None of them identifies a specific published archive, so none
// may ever reach Helm as a chart version - that is how a cluster ends up
// running an operator nobody can name.
func isPinnedChartVersion(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" {
		return false
	}
	_, err := semver.NewVersion(v)
	return err == nil
}

// isPrereleaseChartVersion reports whether v carries a semver prerelease
// segment (v1.0.66-rc.1, v0.6.0-rc.6). Such charts are published to this
// repository ahead of the release they precede, so by semver they outrank every
// stable chart below them - which is why they are excluded from every selection
// Meshery makes on the user's behalf. An unparseable version is not a
// prerelease; it is not a version at all, and isPinnedChartVersion has already
// rejected it before this is reached.
func isPrereleaseChartVersion(v string) bool {
	sv, err := semver.NewVersion(strings.TrimSpace(v))
	if err != nil {
		return false
	}
	return sv.Prerelease() != ""
}

// compareChartVersions orders two chart versions by semver. Unparseable
// versions sort below every parseable one, so they can never be picked as
// "newest".
func compareChartVersions(a, b string) int {
	av, aerr := semver.NewVersion(strings.TrimSpace(a))
	bv, berr := semver.NewVersion(strings.TrimSpace(b))
	switch {
	case aerr != nil && berr != nil:
		return strings.Compare(a, b)
	case aerr != nil:
		return -1
	case berr != nil:
		return 1
	default:
		return av.Compare(bv)
	}
}

// ResolveOperatorChartVersion pins requested to a Meshery Operator chart
// version that actually exists in published, which is the set of versions the
// chart repository at repo advertises (any order; unparseable entries are
// ignored). repo is used only to name the repository that was actually read in
// the error and reason text, so a mirror or an in-cluster repository is
// reported as itself rather than as the default.
//
// It returns the version to deploy and, when that names a different release
// than the one asked for, a one-sentence reason to log and surface to the user.
// A substitution is never silent: reason is empty if and only if the returned
// version is the requested one. The returned version is always the
// repository's own spelling of that release, which is what Helm is handed - a
// request for "1.0.64" against a repository publishing "v1.0.64" is the same
// release, not a substitution, and resolves without a reason.
//
// Prereleases (v1.0.66-rc.1) are published to this repository, and by semver
// they outrank every stable chart below them. They are therefore excluded from
// every version Meshery *chooses* - the newest-published fallback and the
// floor raise - so a release candidate published ahead of its release is never
// what lands on a production cluster by default. They remain fully reachable:
// an explicit `operator.version` naming one is honored exactly, as is a derived
// version that already matches a published prerelease, because a prerelease
// server asking for its own prerelease chart chose nothing.
//
// The rules, in order:
//
//   - Nothing published at all: fail. No version could work.
//   - An explicit request must be pinned and published, or it fails.
//     Substituting for it would deploy something the user did not ask for.
//   - A derived (server-release) request that is published and at or above
//     MinimumOperatorChartVersion is used verbatim. This is the pre-existing
//     behavior and the common case on a current server.
//   - A derived request below the floor is raised to the oldest published
//     stable chart at or above it: the smallest change that yields a working
//     operator on an old server.
//   - A derived request that is unpinned or unpublished falls back to the
//     newest published stable chart, because chart publishing trails server
//     releases and the newest chart is the closest thing to "the chart for
//     this server".
//   - Any of the last two with no stable chart published at all: fail, rather
//     than select a prerelease nobody asked for.
func ResolveOperatorChartVersion(repo string, published []string, requested string, source OperatorChartVersionSource) (version, reason string, err error) {
	requested = strings.TrimSpace(requested)

	pinned := make([]string, 0, len(published))
	for _, v := range published {
		if isPinnedChartVersion(v) {
			pinned = append(pinned, strings.TrimSpace(v))
		}
	}
	if len(pinned) == 0 {
		return "", "", ErrNoOperatorChartPublished(OperatorChartName, repo)
	}
	sort.SliceStable(pinned, func(i, j int) bool { return compareChartVersions(pinned[i], pinned[j]) > 0 })

	// stable is the set every automatic selection draws from, in the same
	// newest-first order.
	stable := make([]string, 0, len(pinned))
	for _, v := range pinned {
		if !isPrereleaseChartVersion(v) {
			stable = append(stable, v)
		}
	}

	// newestSelectable is the version an error tells the user to reach for
	// instead. It names a stable chart whenever one is published, because
	// recommending a release candidate is not advice worth giving.
	newestSelectable := pinned[0]
	if len(stable) > 0 {
		newestSelectable = stable[0]
	}

	// publishedSpelling returns the repository's own spelling of the release v
	// names, or "" when the repository publishes no such release. Membership is
	// decided by semver, not by string equality, because Helm treats "1.0.64"
	// and "v1.0.64" as the same version and rejecting one spelling would refuse
	// a chart that exists. It searches every published chart, prereleases
	// included: naming a release exactly is not choosing one.
	publishedSpelling := func(v string) string {
		if !isPinnedChartVersion(v) {
			return ""
		}
		for _, p := range pinned {
			if compareChartVersions(p, v) == 0 {
				return p
			}
		}
		return ""
	}

	if source == OperatorChartVersionRequested {
		if !isPinnedChartVersion(requested) {
			return "", "", ErrOperatorChartNotPinned(requested, newestSelectable)
		}
		match := publishedSpelling(requested)
		if match == "" {
			return "", "", ErrOperatorChartNotPublished(requested, newestSelectable)
		}
		return match, "", nil
	}

	match := publishedSpelling(requested)
	if match != "" && compareChartVersions(match, MinimumOperatorChartVersion) >= 0 {
		return match, "", nil
	}

	// Everything below picks a chart the request did not name, so it picks from
	// stable releases only. With none published there is nothing to pick.
	if len(stable) == 0 {
		return "", "", ErrNoOperatorChartPublished(OperatorChartName, repo)
	}
	newest := stable[0]

	if match == "" {
		if !isPinnedChartVersion(requested) {
			return newest, "This Meshery Server carries no pinned release version, so Meshery Operator was deployed from the newest published chart, " + newest + ".", nil
		}
		return newest, "Meshery Operator chart " + requested + " is not published in " + repo +
			" (chart publishing trails Meshery Server releases), so the newest published chart, " + newest + ", was deployed instead.", nil
	}

	floored := oldestPublishedAtLeast(stable, MinimumOperatorChartVersion)
	if floored == "" {
		// Nothing published reaches the floor, so newest is no more likely to
		// deploy than the requested chart. Saying it is "the oldest working
		// chart" would send the user looking for a fix that has already been
		// applied; the actionable fact is that no chart verified to deploy is
		// published yet.
		return newest, "Meshery Operator chart " + requested + " is older than " + MinimumOperatorChartVersion +
			", the oldest chart verified to deploy successfully, and " + repo + " publishes no chart at or above " +
			MinimumOperatorChartVersion + " yet, so the newest published chart, " + newest +
			", was deployed instead. It is older too; Meshery Operator may not become ready until a chart at or above " +
			MinimumOperatorChartVersion + " is published.", nil
	}
	return floored, "Meshery Operator chart " + requested + " is older than " + MinimumOperatorChartVersion +
		", the oldest chart verified to deploy successfully, so " + floored +
		" was deployed instead. Set operator.version on the connection to override.", nil
}

// oldestPublishedAtLeast returns the smallest version in versions that is at or
// above floor, or "" when every version is older than floor.
func oldestPublishedAtLeast(versions []string, floor string) string {
	best := ""
	for _, v := range versions {
		if compareChartVersions(v, floor) < 0 {
			continue
		}
		if best == "" || compareChartVersions(v, best) < 0 {
			best = v
		}
	}
	return best
}
