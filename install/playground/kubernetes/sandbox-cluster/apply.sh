#!/usr/bin/env bash
# Applies the playground's sandbox-cluster settings for Meshery Operator, MeshSync and
# Broker. See ../../readme.md, "Sandbox cluster: Meshery Operator, MeshSync and Broker".
#
# Usage: apply.sh <sandbox-kube-context> <playground-egress-ip-or-cidr>
#   e.g. apply.sh staging 203.0.113.10
#
# The Broker CRD does not validate loadBalancerSourceRanges, so an unsubstituted or
# malformed value would be accepted and only fail later inside the operator's reconcile.
# This script refuses to apply anything until the address is a valid IPv4 address/CIDR
# other than 0.0.0.0/0, and accepts only the sandbox cluster's context (staging, unless
# MESHERY_PLAYGROUND_SANDBOX_CONTEXT names another), so a mistyped context can never
# apply the grant and the Broker exposure to the prod cluster the playground runs on.
set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <sandbox-kube-context> <playground-egress-ip-or-cidr>" >&2
  exit 2
fi

context=$1
cidr=$2
[[ $cidr == */* ]] || cidr="$cidr/32"

# The sandbox cluster is the one the playground manages, and its kubeconfig
# context is staging. Accept only that context: the Broker existence check below
# proves a Broker exists, not which cluster this is, and the playground's own
# prod cluster also carries a meshery-operator release, so a mistyped context
# would grant cluster-wide RBAC and expose the production Broker. A sandbox
# whose context was renamed can opt in explicitly with
# MESHERY_PLAYGROUND_SANDBOX_CONTEXT=<context>.
sandbox_context=${MESHERY_PLAYGROUND_SANDBOX_CONTEXT:-staging}
if [[ $context != "$sandbox_context" ]]; then
  echo "error: context '$context' is not the sandbox cluster (expected '$sandbox_context'); refusing to apply" >&2
  exit 2
fi

ipv4_cidr='^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])/(3[0-2]|[12]?[0-9])$'
if [[ ! $cidr =~ $ipv4_cidr ]]; then
  echo "error: '$2' is not an IPv4 address or CIDR" >&2
  exit 2
fi

# 0.0.0.0/0 passes the format check but opens the Broker to every address,
# including the unauthenticated NATS monitoring port. Refuse it.
if [[ $cidr == 0.0.0.0/* ]]; then
  echo "error: '$2' opens the Broker to every address; pass the playground cluster's NAT egress address instead" >&2
  exit 2
fi

here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
apply=(kubectl --context "$context" apply --server-side --field-manager=meshery-playground-ops)

if ! kubectl --context "$context" -n meshery get broker meshery-broker >/dev/null; then
  echo "error: broker meshery-broker not found in namespace meshery on context '$context';" \
    "install the meshery-operator chart first" >&2
  exit 1
fi

# Single quotes are deliberate: envsubst substitutes only the variable named literally.
# shellcheck disable=SC2016
PLAYGROUND_EGRESS_CIDR=$cidr envsubst '${PLAYGROUND_EGRESS_CIDR}' \
  <"$here/broker-service-exposure.yaml" | "${apply[@]}" -f -

# Apply the discovery RBAC only after the Broker apply succeeds: with
# set -euo pipefail a failed Broker step exits before the grant is active.
"${apply[@]}" -f "$here/meshsync-discovery-interim-rbac.yaml"

# MeshSync backs off between failed watches; restart it so read access takes effect now.
kubectl --context "$context" -n meshery delete pod -l app=meshery,component=meshsync --wait=false

echo "Waiting for the Broker to publish an external endpoint..."
for _ in $(seq 1 60); do
  # A transient kubectl failure must retry until the timeout, not exit under set -e.
  endpoint=$(kubectl --context "$context" -n meshery get broker meshery-broker \
    -o jsonpath='{.status.endpoint.external}' || true)
  # An external endpoint from before this apply does not prove the new source range
  # took effect; the operator updates the Service asynchronously. Only report
  # success once the Service carries exactly the requested CIDR and nothing else:
  # a substring match would accept a Service that still permits every source
  # (for example [0.0.0.0/0, <cidr>]).
  ranges=$(kubectl --context "$context" -n meshery get svc meshery-nats \
    -o jsonpath='{.spec.loadBalancerSourceRanges[*]}' || true)
  if [[ -n $endpoint && $ranges == "$cidr" ]]; then
    echo "Broker external endpoint: $endpoint (reachable only from $cidr)"
    exit 0
  fi
  sleep 10
done
echo "error: no external endpoint after 10 minutes; check 'kubectl --context $context -n meshery describe svc meshery-nats'" >&2
exit 1
