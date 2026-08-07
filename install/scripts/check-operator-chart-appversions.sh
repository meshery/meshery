#!/usr/bin/env bash
#
# Assert that the meshery-operator chart and its two subcharts advertise the
# same appVersion.
#
# The subcharts ship only the Broker and MeshSync custom resources; the operator
# release that reconciles them is what their appVersion names. That value is
# baked into an immutable published archive, so a subchart left behind
# advertises an operator release it was not published alongside - which is what
# "appVersion: stable-latest" did before it was pinned.
#
# The parent's version/appVersion are stamped by hack/sync-downstream.sh in
# meshery/meshery-operator on every operator release, and that mechanism touches
# only the parent. `helm lint` does not compare a parent's appVersion with its
# subcharts', so nothing else catches the drift. This check does.

set -euo pipefail

# read_app_version prints the top-level appVersion of a Chart.yaml as a bare
# version string.
#
# The trailing comment goes first, then trailing whitespace, then the quotes.
# Any other order defeats itself: stripping quotes from `"1.0.5" # stamped`
# finds no quote at the end of the line, so the closing one survives into the
# comparison and the check fails against a subchart that reads identically.
# --self-test pins that ordering.
read_app_version() {
  local file="$1" value
  if [ ! -f "$file" ]; then
    echo "check-operator-chart-appversions: $file does not exist" >&2
    exit 1
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
    exit 1
  fi
  printf '%s' "$value"
}

# self_test asserts read_app_version against every spelling a Chart.yaml may
# legitimately use, so the strip ordering above cannot silently regress.
self_test() {
  local dir failures=0 case_line expected got
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

  rm -rf "$dir"
  if [ "$failures" -ne 0 ]; then
    exit 1
  fi
  echo "check-operator-chart-appversions: self-test passed (${#cases[@]} cases)"
}

if [ "${1:-}" = "--self-test" ]; then
  self_test
  exit 0
fi

CHART_DIR="${1:-install/kubernetes/helm/meshery-operator}"

PARENT_CHART="${CHART_DIR}/Chart.yaml"
BROKER_CHART="${CHART_DIR}/charts/meshery-broker/Chart.yaml"
MESHSYNC_CHART="${CHART_DIR}/charts/meshery-meshsync/Chart.yaml"

parent="$(read_app_version "$PARENT_CHART")"
broker="$(read_app_version "$BROKER_CHART")"
meshsync="$(read_app_version "$MESHSYNC_CHART")"

if [ "$parent" != "$broker" ] || [ "$parent" != "$meshsync" ]; then
  cat >&2 <<EOF
check-operator-chart-appversions: the meshery-operator chart and its subcharts disagree on appVersion.

  ${PARENT_CHART}: ${parent}
  ${BROKER_CHART}: ${broker}
  ${MESHSYNC_CHART}: ${meshsync}

All three must name the same Meshery Operator release. The parent's appVersion is
stamped automatically on every operator release; update both subcharts to match it.
EOF
  exit 1
fi

echo "check-operator-chart-appversions: meshery-operator and both subcharts agree on appVersion ${parent}"
