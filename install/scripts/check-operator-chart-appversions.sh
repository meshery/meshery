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

CHART_DIR="${1:-install/kubernetes/helm/meshery-operator}"

PARENT_CHART="${CHART_DIR}/Chart.yaml"
BROKER_CHART="${CHART_DIR}/charts/meshery-broker/Chart.yaml"
MESHSYNC_CHART="${CHART_DIR}/charts/meshery-meshsync/Chart.yaml"

# read_app_version prints the top-level appVersion of a Chart.yaml, with any
# surrounding quotes stripped.
read_app_version() {
  local file="$1" value
  if [ ! -f "$file" ]; then
    echo "check-operator-chart-appversions: $file does not exist" >&2
    exit 1
  fi
  value="$(sed -n 's/^appVersion:[[:space:]]*//p' "$file" | head -n 1)"
  value="${value%\"}"
  value="${value#\"}"
  value="${value%\'}"
  value="${value#\'}"
  # Trim a trailing comment and any trailing whitespace.
  value="${value%%#*}"
  value="$(printf '%s' "$value" | sed 's/[[:space:]]*$//')"
  if [ -z "$value" ]; then
    echo "check-operator-chart-appversions: $file declares no appVersion" >&2
    exit 1
  fi
  printf '%s' "$value"
}

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
