#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries
}

@test "mesheryctl system dashboard fails when kubeconfig is missing" {
  run $MESHERYCTL_BIN system dashboard
  assert_failure
  assert_output --regexp "Meshery Server is not running"
}

@test "mesheryctl system dashboard fails when server is unreachable" {
  run $MESHERYCTL_BIN system dashboard
  assert_failure
  assert_output --regexp "Meshery Server is not running|connection refused|unreachable|timed out"
}
