#!/usr/bin/env bats

# Load required BATS libraries
load '../helpers/bats-support/load'
load '../helpers/bats-assert/load'

# Test 1: Verify successful deletion of a design
@test "mesheryctl design delete is succeeded" {
  # Check for design ID from previous test
  DESIGN_ID_FILE="${TEMP_TEST_DATA_DIR}/design/id"

  if [ ! -f "$DESIGN_ID_FILE" ]; then
    skip "No design ID available to delete"
  fi

  DESIGN_ID=$(cat "$DESIGN_ID_FILE")
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