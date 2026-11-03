#!/usr/bin/env bash

# Notes on design:
# - Full-environment cleanup (Helm releases, CRDs, namespace, local binaries)
#   is expensive (namespace teardown alone can take up to a minute), so it
#   runs once defensively in setup_file() to clear any state left over from a
#   crashed previous run, and then once after every test in teardown(). It is
#   intentionally NOT re-run in the per-test setup(), to avoid paying the
#   cleanup cost twice per test.

remove_all() {
  local reason="${1:-cleanup}"
  echo "# --- remove_all (${reason}) ---" >&3

  # Stop Meshery if mesheryctl is available.
  if command -v mesheryctl >/dev/null 2>&1; then
    mesheryctl system stop --yes 2>/dev/null || true
  fi

  # Stop/remove any Meshery docker containers directly, in case
  # `mesheryctl system stop` above didn't fully tear down the docker
  # platform deployment (e.g. mesheryctl itself was never installed).
  if command -v docker >/dev/null 2>&1; then
    local containers
    containers="$(docker ps -aq --filter "name=meshery" 2>/dev/null || true)"
    if [ -n "$containers" ]; then
      # shellcheck disable=SC2086
      docker ps -aq --filter "name=meshery" | xargs -r docker rm -f >/dev/null 2>&1 || true
    fi
  fi

  # Remove Meshery Helm releases.
  helm uninstall meshery -n meshery 2>/dev/null || true
  helm uninstall meshery-operator -n meshery 2>/dev/null || true

  # Remove Meshery custom resources.
  kubectl delete brokers.meshery.io meshery-broker \
    -n meshery \
    --ignore-not-found \
    --wait=false \
    2>/dev/null || true

  kubectl delete meshsyncs.meshery.io meshery-meshsync \
    -n meshery \
    --ignore-not-found \
    --wait=false \
    2>/dev/null || true

  # Remove finalizers from Meshery custom resources if they remain.
  kubectl patch brokers.meshery.io meshery-broker \
    -n meshery \
    --type=json \
    -p='[{"op":"remove","path":"/metadata/finalizers"}]' \
    2>/dev/null || true

  kubectl patch meshsyncs.meshery.io meshery-meshsync \
    -n meshery \
    --type=json \
    -p='[{"op":"remove","path":"/metadata/finalizers"}]' \
    2>/dev/null || true

  # Delete the Meshery namespace.
  kubectl delete namespace meshery \
    --wait=false \
    2>/dev/null || true

  # Remove namespace finalizer if the namespace is stuck terminating.
  if kubectl get namespace meshery >/dev/null 2>&1; then
    kubectl patch namespace meshery \
      --type=json \
      -p='[{"op":"remove","path":"/spec/finalizers"}]' \
      2>/dev/null || true
  fi

  # Wait until the namespace is completely removed, logging clearly if we
  # give up rather than silently proceeding into a still-terminating state.
  local waited=0
  local max_wait=60
  local interval=2
  while kubectl get namespace meshery >/dev/null 2>&1; do
    if [ "$waited" -ge "$max_wait" ]; then
      echo "# WARNING: namespace 'meshery' still present after ${max_wait}s; proceeding anyway" >&3
      break
    fi
    sleep "$interval"
    waited=$((waited + interval))
  done

  # Remove locally installed mesheryctl config/state, if this test installed it.
  if [ -d "$HOME/.meshery" ]; then
    rm -rf "$HOME/.meshery"
  fi

  # Remove meshery.tar.gz if it exists in the test directory.
  if [ -f "$BATS_TEST_DIRNAME/meshery.tar.gz" ]; then
    rm -f "$BATS_TEST_DIRNAME/meshery.tar.gz"
  fi

  # Remove Homebrew-installed mesheryctl.
  if command -v brew >/dev/null 2>&1; then
    brew uninstall mesheryctl 2>/dev/null || true
  fi

  # Remove mesheryctl installed outside Homebrew. Try with sudo as a
  # fallback since /usr/local/bin and /usr/bin are often root-owned on CI
  # runners, and a plain `rm -f` there would silently no-op under `|| true`.
  for bin_dir in /usr/local/bin /usr/bin "$HOME/.local/bin"; do
    if [ -f "$bin_dir/mesheryctl" ]; then
      rm -f "$bin_dir/mesheryctl" 2>/dev/null \
        || sudo rm -f "$bin_dir/mesheryctl" 2>/dev/null \
        || echo "# WARNING: could not remove $bin_dir/mesheryctl" >&3
    fi
  done
}