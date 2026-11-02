#!/usr/bin/env bash
#
# Assert that every meshery-operator subchart in this repository advertises the
# same Meshery Operator release as the parent chart it ships under.
#
# WHAT THIS IS FOR
#
# The subcharts ship only the Broker and MeshSync custom resources; the operator
# release that reconciles them is what their appVersion names. Left at a moving
# tag - "appVersion: stable-latest", which is what both carried before they were
# pinned - an immutable published archive advertises whatever happened to be
# newest at some unrelated later moment. This check keeps the subcharts pinned to
# a real, published Meshery Operator release and keeps the in-repo values
# consistent with one another. `helm lint` does not compare a parent's appVersion
# with its subcharts', so nothing else catches the drift.
#
# WHAT THIS IS NOT FOR
#
# It does not assert - and must not be read as asserting - that the *published*
# archive has a parent and subcharts that agree. They do not, by design:
# .github/workflows/helm-chart-releaser.yml packages the chart through
# helm-gh-pages with `app_version: <Meshery Server tag>`, and that rewrites only
# the top-level Chart.yaml. So the published parent advertises a Meshery Server
# release while the subcharts advertise a Meshery Operator release - two separate
# numbering lines. Confirmed against the published meshery-operator-v1.0.65.tgz,
# whose parent reads v1.0.65 while its subcharts read what the repository said.
#
# WHAT IS INSPECTED
#
#   1. The standalone chart at install/kubernetes/helm/meshery-operator. Its
#      parent version/appVersion are stamped by hack/sync-downstream.sh in
#      meshery/meshery-operator on every operator release, and that mechanism
#      touches only the parent. This repository owns these files, so drift here
#      is a HARD FAILURE.
#
#   2. The packaged copy the meshery chart vendors at
#      install/kubernetes/helm/meshery/charts/meshery-operator-*.tgz - the
#      artifact `helm install meshery/meshery` actually deploys, and the
#      deployment shape the original operator-image bug report came from. It is
#      vendored from the published chart repository, so it cannot be corrected in
#      place: hand-repacking it would diverge from upstream and be reverted by
#      the next `helm dependency update`. Its subcharts can carry the pin only
#      once a meshery-operator chart built from the pinned sources in (1) has
#      been published and re-vendored here. Until then drift is reported LOUDLY
#      but NON-FATALLY, because a hard failure would block the very change that
#      makes the fix possible.
#
#      TO PROMOTE (2) TO A HARD FAILURE: set BUNDLED_CHART_DRIFT_IS_FATAL to
#      true, below. Do that once - and only once - a published meshery-operator
#      chart carrying pinned subchart appVersions has been re-vendored into
#      install/kubernetes/helm/meshery/charts/. Nothing else needs to change.

set -euo pipefail

# BUNDLED_CHART_DRIFT_IS_FATAL promotes the vendored-archive check from a warning
# to a build failure. See "WHAT IS INSPECTED" (2) for the condition that makes
# flipping this safe.
BUNDLED_CHART_DRIFT_IS_FATAL="${BUNDLED_CHART_DRIFT_IS_FATAL:-false}"

# BUNDLED_CHART_GLOB locates the vendored operator archive. The version is
# globbed rather than named so a `helm dependency update` that bumps it does not
# silently take the archive out of this check's reach.
BUNDLED_CHART_GLOB="${BUNDLED_CHART_GLOB:-install/kubernetes/helm/meshery/charts/meshery-operator-*.tgz}"

# read_app_version prints the top-level appVersion of a Chart.yaml as a bare
# version string.
#
# The trailing comment goes first, then trailing whitespace, then the quotes.
# Any other order defeats itself: stripping quotes from `"1.0.5" # stamped`
# finds no quote at the end of the line, so the closing one survives into the
# comparison and the check fails against a subchart that reads identically.
# --self-test pins that ordering.
#
# A malformed or missing chart RETURNS 2 rather than 1, so a caller can tell
# "this tree could not be read" from "this tree disagrees with itself" and report
# the one that actually happened. It must be return, not exit: every call site
# is a command substitution, so an exit would end only that nested subshell and
# leave the caller comparing an empty string - which reports drift, the wrong
# answer, and makes the read-failure branch unreachable. Callers must therefore
# check the status of each read rather than only its output.
read_app_version() {
  local file="$1" value
  if [ ! -f "$file" ]; then
    echo "check-operator-chart-appversions: $file does not exist" >&2
    return 2
  fi
  value="$(sed -n 's/^appVersion:[[:space:]]*//p' "$file" | head -n 1)"
  # YAML requires whitespace before an inline comment, so anchoring on that
  # leaves a '#' inside the value alone.
  value="$(printf '%s' "$value" | sed -e 's/[[:space:]]#.*$//' -e 's/[[:space:]]*$//')"
  value="${value%\"}"
  value="${value#\"}"
  value="${value%\'}"
  value="${value#\'}"
  if [ -z "$value" ]; then
    echo "check-operator-chart-appversions: $file declares no appVersion" >&2
    return 2
  fi
  printf '%s' "$value"
}

# compare_chart_tree reports the parent and subchart appVersions of an unpacked
# chart directory on stdout, one "  <path>: <version>" line each. It returns 1
# when a subchart disagrees with the parent and 2 when the tree could not be
# read at all.
#
# Subcharts are discovered rather than named, so a third one added later is
# covered the day it appears instead of passing unexamined.
compare_chart_tree() {
  local chart_dir="$1"
  local parent subchart_dir subchart_yaml subchart_version drift=0 found=0

  if ! parent="$(read_app_version "${chart_dir}/Chart.yaml")"; then
    return 2
  fi
  printf '  %s: %s\n' "${chart_dir}/Chart.yaml" "$parent"

  for subchart_dir in "${chart_dir}"/charts/*/; do
    subchart_yaml="${subchart_dir}Chart.yaml"
    [ -f "$subchart_yaml" ] || continue
    found=1
    if ! subchart_version="$(read_app_version "$subchart_yaml")"; then
      return 2
    fi
    printf '  %s: %s\n' "$subchart_yaml" "$subchart_version"
    if [ "$subchart_version" != "$parent" ]; then
      drift=1
    fi
  done

  if [ "$found" -eq 0 ]; then
    echo "check-operator-chart-appversions: ${chart_dir} has no subcharts to check" >&2
    return 2
  fi
  return "$drift"
}

# self_test asserts read_app_version against every spelling a Chart.yaml may
# legitimately use, so the strip ordering above cannot silently regress.
self_test() {
  local dir failures=0 case_line expected got tree_status
  dir="$(mktemp -d)"

  # Each case is "<expected>|<appVersion line>".
  local cases=(
    '1.0.4|appVersion: 1.0.4'
    '1.0.4|appVersion: "1.0.4"'
    "1.0.4|appVersion: '1.0.4'"
    '1.0.4|appVersion: "1.0.4" # stamped by sync-downstream'
    '1.0.4|appVersion: 1.0.4 # stamped by sync-downstream'
    "1.0.4|appVersion: '1.0.4'   "
    '1.0.4|appVersion:    1.0.4'
    'stable-latest|appVersion: stable-latest'
  )

  for case_line in "${cases[@]}"; do
    expected="${case_line%%|*}"
    printf '%s\n' "${case_line#*|}" > "${dir}/Chart.yaml"
    got="$(read_app_version "${dir}/Chart.yaml")"
    if [ "$got" != "$expected" ]; then
      echo "check-operator-chart-appversions: self-test failed for '${case_line#*|}': parsed '${got}', want '${expected}'" >&2
      failures=$((failures + 1))
    fi
  done

  # compare_chart_tree must call agreement agreement and drift drift, on the
  # same shapes the two real chart trees take.
  mkdir -p "${dir}/tree/charts/meshery-broker" "${dir}/tree/charts/meshery-meshsync"
  printf 'appVersion: "1.0.5"\n' > "${dir}/tree/Chart.yaml"
  printf 'appVersion: "1.0.5"\n' > "${dir}/tree/charts/meshery-broker/Chart.yaml"
  printf 'appVersion: "1.0.5"\n' > "${dir}/tree/charts/meshery-meshsync/Chart.yaml"
  if ! compare_chart_tree "${dir}/tree" > /dev/null; then
    echo "check-operator-chart-appversions: self-test failed: agreeing charts reported as drift" >&2
    failures=$((failures + 1))
  fi
  printf 'appVersion: stable-latest\n' > "${dir}/tree/charts/meshery-meshsync/Chart.yaml"
  if compare_chart_tree "${dir}/tree" > /dev/null; then
    echo "check-operator-chart-appversions: self-test failed: a subchart on stable-latest reported as agreement" >&2
    failures=$((failures + 1))
  fi

  # An unreadable tree must report 2, not 1. What this pins is that
  # compare_chart_tree CHECKS the status of every read: drop either guard and an
  # unreadable chart is silently downgraded to ordinary drift (status 1), which
  # is what shipped before and what made the read-failure branch unreachable.
  # It does not distinguish read_app_version returning from exiting - with the
  # caller's guard in place both propagate - so do not read it as pinning that.
  # These two cases expect a NON-ZERO status, so each call must be written in a
  # form errexit tolerates. A bare `compare_chart_tree ...` followed by `case $?`
  # aborts the whole script under `set -e` before the status is ever inspected -
  # the assertion never runs and the self-test exits 2 instead of reporting.
  printf 'appVersion: "1.0.5"\n' > "${dir}/tree/charts/meshery-meshsync/Chart.yaml"
  printf 'name: no-appversion-here\n' > "${dir}/tree/Chart.yaml"
  tree_status=0
  compare_chart_tree "${dir}/tree" > /dev/null 2>&1 || tree_status=$?
  if [ "$tree_status" -ne 2 ]; then
    echo "check-operator-chart-appversions: self-test failed: a parent with no appVersion must report 2 (unreadable), got ${tree_status}" >&2
    failures=$((failures + 1))
  fi
  printf 'appVersion: "1.0.5"\n' > "${dir}/tree/Chart.yaml"
  printf 'name: broken\n' > "${dir}/tree/charts/meshery-broker/Chart.yaml"
  tree_status=0
  compare_chart_tree "${dir}/tree" > /dev/null 2>&1 || tree_status=$?
  if [ "$tree_status" -ne 2 ]; then
    echo "check-operator-chart-appversions: self-test failed: a subchart with no appVersion must report 2 (unreadable), got ${tree_status}" >&2
    failures=$((failures + 1))
  fi

  rm -rf "$dir"
  if [ "$failures" -ne 0 ]; then
    exit 1
  fi
  echo "check-operator-chart-appversions: self-test passed (${#cases[@]} parse cases, 4 tree cases)"
}

# check_bundled_chart inspects the vendored packaged operator chart. It never
# fails the build on its own; the caller decides, so that the promotion
# described in the header is a one-line change.
check_bundled_chart() {
  local archive="$1" work root report status=0

  if ! command -v tar > /dev/null 2>&1; then
    echo "check-operator-chart-appversions: tar is unavailable, skipping ${archive}" >&2
    return 0
  fi

  work="$(mktemp -d)"
  # shellcheck disable=SC2064
  trap "rm -rf '${work}'" RETURN

  if ! tar -xzf "$archive" -C "$work" 2> /dev/null; then
    echo "check-operator-chart-appversions: could not unpack ${archive}, skipping it" >&2
    return 0
  fi

  root="$(find "$work" -mindepth 2 -maxdepth 2 -name Chart.yaml -print 2> /dev/null | head -n 1)"
  if [ -z "$root" ]; then
    echo "check-operator-chart-appversions: ${archive} contains no chart, skipping it" >&2
    return 0
  fi
  root="$(dirname "$root")"

  report="$(compare_chart_tree "$root")" || status=$?
  # The report names the temporary extraction path, which means nothing to a
  # reader; rewrite it to the archive it came out of.
  report="${report//${work}\//${archive}!}"

  if [ "$status" -eq 0 ]; then
    echo "check-operator-chart-appversions: ${archive} and its subcharts agree on appVersion"
    return 0
  fi
  if [ "$status" -ne 1 ]; then
    echo "check-operator-chart-appversions: ${archive} could not be read, skipping it" >&2
    return 0
  fi

  cat >&2 <<EOF

================================================================================
check-operator-chart-appversions: WARNING - the vendored Meshery Operator chart
disagrees with its own subcharts on appVersion.

  ${archive}

${report}

This archive is vendored from the published chart repository, so it cannot be
corrected here: it will carry the pin only after a meshery-operator chart built
from install/kubernetes/helm/meshery-operator is published and re-vendored. This
is a warning, not a failure, for exactly that reason - see the header of
$0 for how to promote it once that has happened.
================================================================================

EOF
  return 1
}

if [ "${1:-}" = "--self-test" ]; then
  self_test
  exit 0
fi

CHART_DIR="${1:-install/kubernetes/helm/meshery-operator}"

standalone_report=""
standalone_status=0
standalone_report="$(compare_chart_tree "$CHART_DIR")" || standalone_status=$?

if [ "$standalone_status" -ne 0 ] && [ "$standalone_status" -ne 1 ]; then
  # read_app_version or compare_chart_tree already named the exact problem on
  # stderr; restating it as an appVersion disagreement would misdescribe it.
  exit "$standalone_status"
fi

if [ "$standalone_status" -ne 0 ]; then
  cat >&2 <<EOF
check-operator-chart-appversions: the meshery-operator chart and its subcharts disagree on appVersion.

${standalone_report}

All of them must name the same published Meshery Operator release. The parent's
appVersion is stamped automatically on every operator release and that stamping
touches only the parent, so update the subcharts to match it. A moving tag such
as stable-latest is never acceptable here: the archive is immutable, so the
application it advertises would change out from under it.
EOF
  exit 1
fi

echo "check-operator-chart-appversions: ${CHART_DIR} and its subcharts agree on appVersion"

bundled_drift=0
shopt -s nullglob
for bundled_archive in $BUNDLED_CHART_GLOB; do
  check_bundled_chart "$bundled_archive" || bundled_drift=1
done
shopt -u nullglob

if [ "$bundled_drift" -ne 0 ] && [ "$BUNDLED_CHART_DRIFT_IS_FATAL" = "true" ]; then
  echo "check-operator-chart-appversions: vendored chart drift is configured as fatal" >&2
  exit 1
fi

exit 0
