#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries
}

@test "mesheryctl system dashboard fails when kubeconfig is missing" {
  run $MESHERYCTL_BIN system dashboard
  assert_failure
  assert_output --regexp "kubeconfig|no such file|no.*directory|cluster|unreachable"
}
