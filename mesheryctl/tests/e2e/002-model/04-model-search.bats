#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries
}

@test "mesheryctl model search displays usage instructions when no query provided" {
  run $MESHERYCTL_BIN model search
  assert_failure

  assert_output --partial "Please provide a model name."
  assert_output --partial "Usage: mesheryctl model search [query-text]"
}

@test "mesheryctl model search succeeds with default options" {
  run $MESHERYCTL_BIN model search accurate
  assert_success

  assert_output --partial "accurate"
  assert_output --partial "MODEL"
  assert_output --partial "CATEGORY"
  assert_output --partial "VERSION"
  
}

@test "mesheryctl model search for non-existing model" {
  run $MESHERYCTL_BIN model search random-model
  assert_success

  assert_output --partial "No models found"
}

