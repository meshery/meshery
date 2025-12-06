#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

  export REAL_KUBECONFIG="$HOME/.kube/config"

  mkdir -p "$HOME/.kube"
  cp "$REAL_KUBECONFIG" "$HOME/.kube/bats-kubeconfig" 2>/dev/null || touch "$HOME/.kube/bats-kubeconfig"
  export KUBECONFIG="$HOME/.kube/bats-kubeconfig"

  export MESHERYCTL_CONFIG_PATH="$HOME/.meshery/config.yaml"

  echo "DEBUG: BATS using isolated kubeconfig: $KUBECONFIG"
}

get_platform() {
  $MESHERYCTL_BIN system context view 2>/dev/null \
    | grep -i "platform:" \
    | sed 's/.platform:[[:space:]]//'
}

wait_for_meshery_k8s() {
  echo "DEBUG: Using kubeconfig: $KUBECONFIG"

  for i in {1..40}; do
    POD_COUNT=$(kubectl get pods -n meshery 2>/dev/null | grep -c "Running")
    if [ "$POD_COUNT" -gt 0 ]; then
      echo "DEBUG: Meshery pods detected"
      return 0
    fi
    echo "DEBUG: Waiting for Meshery pods..."
    sleep 2
  done

  echo "DEBUG: Kubernetes NOT ready — printing pods:"
  kubectl get pods -n meshery || true
  return 1
}

wait_for_meshery_docker() {
  for i in {1..40}; do
    STATUS=$($MESHERYCTL_BIN system status 2>/dev/null | grep -i "Running")
    if [ -n "$STATUS" ]; then
      return 0
    fi
    sleep 2
  done
  return 1
}

get_meshery_url() {
  $MESHERYCTL_BIN system context view \
    | grep endpoint \
    | sed 's/.endpoint:[[:space:]]//'
}

@test "mesheryctl system dashboard fails when kubeconfig is missing" {
  mv "$KUBECONFIG" "$KUBECONFIG.bak" || true
  > "$KUBECONFIG"

  run $MESHERYCTL_BIN system dashboard --skip-browser
  assert_output --regexp "Meshery Server is not running|kubeconfig|not found|current-context must exist"

  mv "$KUBECONFIG.bak" "$KUBECONFIG" || true
}

@test "mesheryctl system dashboard fails when meshery server is unreachable" {
  echo "" > "$KUBECONFIG"

  run $MESHERYCTL_BIN system dashboard --skip-browser
  assert_output --regexp "Meshery Server is not running|connection refused|current-context must exist"

  cp "$REAL_KUBECONFIG" "$KUBECONFIG" 2>/dev/null || true
}

@test "mesheryctl system dashboard succeeds when meshery server is running" {
  PLATFORM=$(get_platform)
  [ -z "$PLATFORM" ] && PLATFORM="kubernetes"

  echo "DEBUG: Platform detected: $PLATFORM"

  if [ "$PLATFORM" = "kubernetes" ]; then
      wait_for_meshery_k8s
  else
      $MESHERYCTL_BIN system start >/dev/null 2>&1 &
      wait_for_meshery_docker
  fi

  MESHERY_URL=$(get_meshery_url)
  echo "DEBUG: Using Meshery URL: $MESHERY_URL"

  run $MESHERYCTL_BIN system dashboard --skip-browser
  assert_success
  assert_output --regexp "$MESHERY_URL|Opening Meshery"
}
