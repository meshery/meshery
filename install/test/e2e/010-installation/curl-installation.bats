#!/usr/bin/env bats

remove_all() {
  # Stop Meshery if it is running.
  mesheryctl system stop --yes 2>/dev/null || true

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

  # Wait until the namespace is completely removed.
  for _ in {1..30}; do
    if ! kubectl get namespace meshery >/dev/null 2>&1; then
      break
    fi

    sleep 2
  done

  # Remove locally installed mesheryctl, if this test installed it.
  if [ -f "$HOME/.meshery" ]; then
    rm -rf "$HOME/.meshery"
  fi
  # Remove meshery.tar.gz if it exists in the test directory.
  if [ -f "$BATS_TEST_DIRNAME/meshery.tar.gz" ]; then
    rm -f "$BATS_TEST_DIRNAME/meshery.tar.gz"
  fi
  # Remove Homebrew-installed mesheryctl.
  if brew list mesheryctl >/dev/null 2>&1; then
    brew uninstall mesheryctl >/dev/null 2>&1 || true
  fi

  # Remove mesheryctl installed outside Homebrew.
  if [ -f "/usr/local/bin/mesheryctl" ]; then
    rm -f "/usr/local/bin/mesheryctl"
  fi

  # Remove Homebrew installation of mesheryctl.
  brew uninstall mesheryctl 2>/dev/null || true

}

setup() {
  load '../test_helper/bats-support/load'
  load '../test_helper/bats-assert/load'

  remove_all
}

teardown() {
  remove_all
}


@test "Given wrong platforms provided when installing then installer exits after 5 invalid attempts" {
  run expect -c "
    set timeout 30

    spawn bash -c {curl -L https://meshery.io/install | PLATFORM=bob bash -}

    expect \"platform\"
    send \"deomon\r\"

    expect \"platform\"
    send \"youtube\r\"

    expect \"platform\"
    send \"meshery\r\"

    expect \"platform\"
    send \"\r\"

    expect \"platform\"
    send \"bob\r\"

    expect eof
  "

  assert_success
  assert_output --partial "Too many invalid attempts. Please try again."
}


@test "Given docker as the platform when installing then docker is accepted" {
  run expect -c "
    set timeout 120

    spawn bash -c {curl -L https://meshery.io/install | PLATFORM=docker bash -}

    expect eof
  "

  assert_success
  assert_output --partial "Starting Meshery..."
  assert_output --partial "mesheryctl installed."
  assert_output --partial "Meshery endpoint is accessible."
}


@test "Given kubernetes as the platform when installing then kubernetes is accepted" {
  run expect -c "
    set timeout 120

    spawn bash -c {curl -L https://meshery.io/install | PLATFORM=kubernetes bash -}

    expect eof
  "

  assert_success
  assert_output --partial "Meshery deployed on Kubernetes."
}

@test "Given DEPLOY_MESHERY=false when installing, Only install mesheryctl binary" {
  run expect -c "
    set timeout 120

    spawn bash -c {curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -}

    expect eof
  "

  assert_success
  assert_output --partial "mesheryctl installed."
  assert_output --partial "Run \"mesheryctl system start\" to start Meshery."
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
