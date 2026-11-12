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
# other than 0.0.0.0/0, and refuses the prod context, which is the cluster the
# playground runs on rather than the sandbox cluster it manages.
set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <sandbox-kube-context> <playground-egress-ip-or-cidr>" >&2
  exit 2
fi

context=$1
cidr=$2
[[ $cidr == */* ]] || cidr="$cidr/32"

# The playground Pod runs on the prod cluster, which also carries a
# meshery-operator release. Applying here would grant cluster-wide RBAC and expose
# the production Broker, so refuse it; the sandbox context is staging.
if [[ $context == prod ]]; then
  echo "error: context 'prod' is the cluster the playground runs on, not the sandbox cluster it manages; pass the sandbox context (staging)" >&2
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

"${apply[@]}" -f "$here/meshsync-discovery-interim-rbac.yaml"
# Single quotes are deliberate: envsubst substitutes only the variable named literally.
# shellcheck disable=SC2016
PLAYGROUND_EGRESS_CIDR=$cidr envsubst '${PLAYGROUND_EGRESS_CIDR}' \
  <"$here/broker-service-exposure.yaml" | "${apply[@]}" -f -

# MeshSync backs off between failed watches; restart it so read access takes effect now.
kubectl --context "$context" -n meshery delete pod -l app=meshery,component=meshsync --wait=false

echo "Waiting for the Broker to publish an external endpoint..."
for _ in $(seq 1 60); do
  # A transient kubectl failure must retry until the timeout, not exit under set -e.
  endpoint=$(kubectl --context "$context" -n meshery get broker meshery-broker \
    -o jsonpath='{.status.endpoint.external}' || true)
  # An external endpoint from before this apply does not prove the new source range
  # took effect; the operator updates the Service asynchronously. Only report
  # success once the Service carries the new CIDR.
  ranges=$(kubectl --context "$context" -n meshery get svc meshery-nats \
    -o jsonpath='{.spec.loadBalancerSourceRanges}' || true)
  if [[ -n $endpoint && $ranges == *"$cidr"* ]]; then
    echo "Broker external endpoint: $endpoint (reachable only from $cidr)"
    exit 0
  fi
  sleep 10
done
echo "error: no external endpoint after 10 minutes; check 'kubectl --context $context -n meshery describe svc meshery-nats'" >&2
exit 1
