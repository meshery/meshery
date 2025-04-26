#!/usr/bin/env bats

# Setup function to load libraries
setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries
  
  load "$E2E_HELPERS_PATH/constants"

  export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/design"
}

@test "mesheryctl design delete removes the associate design" {
  DESIGN_ID=$(cat "$TESTDATA_DIR/id")
  [ -n "$DESIGN_ID" ] || skip "Empty design ID"

  # Delete the design
  run $MESHERYCTL_BIN design delete "$DESIGN_ID"
  assert_success
  assert_output --partial "deleted" || assert_output --partial "Design"
}

# Test 2: Verify appropriate response for deleting a non-existent design
@test "mesheryctl design delete for non-existent ID gives appropriate response" {
  # Use a non-existent design ID
  NONEXISTENT_ID="00000000-0000-0000-0000-000000000000"

  # Attempt to delete non-existent design
  run $MESHERYCTL_BIN design delete "$NONEXISTENT_ID"
  assert_failure
  assert_output --partial "Error" || assert_output --partial "not found" || assert_output --partial "failed" || assert_output --partial "No Design"
}