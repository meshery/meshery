#!/usr/bin/env bats
#
# E2E tests for the Meshery installer script.
#
# Notes on design:
# - Full-environment cleanup (Helm releases, CRDs, namespace, local binaries)
#   is expensive (namespace teardown alone can take up to a minute), so it
#   runs once defensively in setup_file() to clear any state left over from a
#   crashed previous run, and then once after every test in teardown(). It is
#   intentionally NOT re-run in the per-test setup(), to avoid paying the
#   cleanup cost twice per test.
# - Each `expect` invocation is wrapped in `timeout` so a hung installer
#   cannot stall the whole CI job past its own step timeout.

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
      docker rm -f $containers >/dev/null 2>&1 || true
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

setup_file() {
  # Defensive one-time cleanup in case a previous run crashed mid-test and
  # left state behind. This does NOT run between every test — teardown()
  # handles that — so it doesn't double the per-test cleanup cost.
  remove_all "setup_file"
}

setup() {
  # Each test runs in its own subshell, so the helper libraries must be
  # re-sourced here every time (loading them once in setup_file() would not
  # persist into each test's subshell). This is intentionally cheap —
  # the expensive environment cleanup lives in setup_file()/teardown(),
  # not here.
  load '../test_helper/bats-support/load'
  load '../test_helper/bats-assert/load'
}

teardown() {
  remove_all "teardown"
}

# Helper: run an `expect` script with a hard wall-clock timeout so a hang
# in the installer cannot stall the whole suite.
run_expect() {
  local script="$1"
  local wall_timeout="${2:-300}"
  run timeout "$wall_timeout" expect -c "$script"
}


@test "Given wrong platforms provided when installing then installer exits after 5 invalid attempts" {
  run_expect '
    set timeout 150
    log_user 1

    spawn bash -c {curl -L https://meshery.io/install | PLATFORM=bob bash -}

    expect {
      -re {(?i)select.*platform|platform.*:} { send "deomon\r" }
      timeout { exit 2 }
    }

    expect {
      -re {(?i)select.*platform|platform.*:} { send "youtube\r" }
      timeout { exit 2 }
    }

    expect {
      -re {(?i)select.*platform|platform.*:} { send "meshery\r" }
      timeout { exit 2 }
    }

    expect {
      -re {(?i)select.*platform|platform.*:} { send "\r" }
      timeout { exit 2 }
    }

    expect {
      -re {(?i)select.*platform|platform.*:} { send "bob\r" }
      timeout { exit 2 }
    }

    expect eof
  ' 180

  assert_success
  assert_output --partial "Too many invalid attempts. Please try again."
  # The script should reject each bad value before ultimately bailing out,
  # not just print the final message with a lucky earlier match.
  refute_output --partial "Starting Meshery..."
}


@test "Given docker as the platform when installing then docker is accepted" {
  run_expect '
    set timeout 300

    spawn bash -c {curl -L https://meshery.io/install | PLATFORM=docker bash -}

    expect eof
  ' 320

  assert_success
  assert_output --partial "Starting Meshery..."
  assert_output --partial "mesheryctl installed."
  assert_output --partial "Meshery endpoint is accessible."

  # Functional check, not just log-message matching: confirm mesheryctl was
  # actually placed on PATH and that Meshery believes itself to be running.
  run bash -lc "command -v mesheryctl"
  assert_success

  run bash -lc "mesheryctl system status"
  assert_success
  assert_output --partial "Running"
}


@test "Given kubernetes as the platform when installing then kubernetes is accepted" {
  run_expect '
    set timeout 300

    spawn bash -c {curl -L https://meshery.io/install | PLATFORM=kubernetes bash -}

    expect eof
  ' 320

  assert_success
  assert_output --partial "Meshery deployed on Kubernetes."

  # Functional check: verify the Meshery pods actually exist and reach
  # Ready state, rather than trusting only the printed success message.
  run bash -lc "kubectl wait --for=condition=Ready pod -l app.kubernetes.io/part-of=meshery -n meshery --timeout=120s"
  assert_success

  run bash -lc "kubectl get pods -n meshery --no-headers"
  assert_success
  refute_output --partial "CrashLoopBackOff"
  refute_output --partial "Error"
}


@test "Given DEPLOY_MESHERY=false when installing, Only install mesheryctl binary" {
  run_expect '
    set timeout 300

    spawn bash -c {curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -}

    expect eof
  ' 320

  assert_success
  assert_output --partial "mesheryctl installed."
  assert_output --partial "Run \"mesheryctl system start\" to start Meshery."

  # DEPLOY_MESHERY=false should mean nothing was actually deployed — verify
  # that no meshery namespace/resources were created as a side effect.
  run bash -lc "kubectl get namespace meshery"
  assert_failure
}


@test "Given Homebrew when installing Meshery CLI then mesheryctl is installed" {
  run brew install mesheryctl
  assert_success
  assert_output --partial "mesheryctl"

  run mesheryctl version
  assert_success
  assert_output --partial "Client"
  assert_output --partial "Server"
}
