#!/usr/bin/env bats

# Setup function to load libraries and prepare test fixtures

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
  export TESTDATA_PATH="$FIXTURES_DIR/design-import/nginx.yaml"
}

@test "mesheryctl design apply applies design file" {
  run $MESHERYCTL_BIN design apply -f "$TESTDATA_PATH"

  assert_success
}

@test "mesheryctl design apply -f invalid path shows error" {
  INVALID_PATH="./test/invalid/path/nginx.yaml"
  run $MESHERYCTL_BIN design apply -f "$INVALID_PATH"

  assert_failure
  assert_output --partial "invalid"
  assert_output --partial "Enter a valid path"
}


