#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

  export REAL_KUBECONFIG="$HOME/.kube/config"

  mkdir -p "$HOME/.kube"
  cp "$REAL_KUBECONFIG" "$HOME/.kube/bats-kubeconfig" 2>/dev/null || touch "$HOME/.kube/bats-kubeconfig"
  export KUBECONFIG="$HOME/.kube/bats-kubeconfig"

  
  export MESHERYCTL_CONFIG_PATH="$HOME/.meshery/config.yaml"
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
  MESHERY_URL=$(get_meshery_url)

  run $MESHERYCTL_BIN system dashboard --skip-browser
  assert_success
  assert_output --regexp "$MESHERY_URL|Opening Meshery"
}
