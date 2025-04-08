#!/usr/bin/env bats

load '../helpers/bats-support/load'
load '../helpers/bats-assert/load'

setup() {
  # Check if Meshery server is running
  if ! curl -s "http://localhost:9081/api/system/version" > /dev/null 2>&1; then
    echo "WARNING: Meshery server is not running, tests may fail"
  fi
}

@test "mesheryctl design delete is succeeded" {
  # Check if the design ID file exists
  if [ ! -f "${TEMP_DIR}/meshery_test_design_id" ]; then
    skip "No design ID available to delete"
  fi
  
  DESIGN_ID=$(cat "${TEMP_DIR}/meshery_test_design_id")
  echo "Using design ID from file: $DESIGN_ID"
  
  # Skip if design ID is empty
  [ -n "$DESIGN_ID" ] || skip "Empty design ID"
  
  # Get list of designs 
  run $MESHERYCTL_BIN design list
  
  # Skip if Meshery server is not running
  if [[ "$output" =~ "connection refused" ]]; then
    skip "Meshery server is not running, skipping test"
  fi
  
  assert_success
  BEFORE_DELETE="$output"
  echo "Designs before delete: $BEFORE_DELETE"
  
  # Delete the design
  run $MESHERYCTL_BIN design delete "$DESIGN_ID"
  echo "Delete command output: $output"
  
  # Skip if Meshery server is not running
  if [[ "$output" =~ "connection refused" ]]; then
    skip "Meshery server is not running, skipping test"
  fi
  
  # Check if command executed without error
  assert_success
  assert_output --partial "deleted" || assert_output --partial "removed"
  
  # Verify the design is deleted by trying to view it
  run $MESHERYCTL_BIN design view "$DESIGN_ID"
  
  # Skip if Meshery server is not running
  if [[ "$output" =~ "connection refused" ]]; then
    skip "Meshery server is not running, skipping test"
  fi
  
  # Should fail as design no longer exists
  assert_failure
  assert_output --partial "not found" || assert_output --partial "error"
}

@test "mesheryctl design delete for non-existent ID should give appropriate response" {
  # Use a non-existent design ID
  NONEXISTENT_ID="00000000-0000-0000-0000-000000000000"
  
  # Attempt to delete non-existent design
  run $MESHERYCTL_BIN design delete "$NONEXISTENT_ID"
  echo "Delete command output: $output"
  
  # Skip if Meshery server is not running
  if [[ "$output" =~ "connection refused" ]]; then
    skip "Meshery server is not running, skipping test"
  fi
  
  # Adjust expected response based on actual error message
  # The command might return various errors for non-existent designs
  assert_output --partial "Error" || assert_output --partial "not found" || assert_output --partial "failed"
  
  # Clean up the temporary files - but only attempt to remove files that exist
  if [ -f "${TEMP_DIR}/meshery_test_design_id" ]; then
    rm -f "${TEMP_DIR}/meshery_test_design_id"
  fi
}