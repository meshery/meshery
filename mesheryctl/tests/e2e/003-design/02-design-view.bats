#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

  export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/design"
}

@test "mesheryctl design view using ID displays the design content" {
  if [[ ! -f "$TESTDATA_DIR/id" ]]; then
    skip "Design ID file not found (import test may have failed)"
  fi
  
  DESIGN_ID=$(cat "$TESTDATA_DIR/id")
  [ -n "$DESIGN_ID" ] || skip "Empty design ID"

  run $MESHERYCTL_BIN design view "$DESIGN_ID"
  assert_success
  assert_output --partial "nginx-deployment"
}

@test "mesheryctl design view using name displays the design content" {
  run $MESHERYCTL_BIN design view "nginx-deployment"
  assert_success
  assert_output --partial "nginx-deployment"
}

@test "mesheryctl design view with no argument throws an error" {
  run $MESHERYCTL_BIN design view
  assert_failure
  assert_output --partial "Provide a design name or ID"
}

@test "mesheryctl design view with nonexistent design throws an error" {
  run $MESHERYCTL_BIN design view "nonexistent-design"
  assert_failure
  assert_output --partial "not found"
}
