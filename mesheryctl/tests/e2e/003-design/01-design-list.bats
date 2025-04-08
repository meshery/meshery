#!/usr/bin/env bats

load '../helpers/bats-support/load'
load '../helpers/bats-assert/load'

setup() {
  # Check if Meshery server is running
  if ! curl -s "http://localhost:9081/api/system/version" > /dev/null 2>&1; then
    echo "WARNING: Meshery server is not running, tests may fail"
  fi
}

@test "mesheryctl design list shows designs" {
  # Run the list command
  run $MESHERYCTL_BIN design list
  echo "List command output: $output"
  
  # Skip if Meshery server is not running
  if [[ "$output" =~ "connection refused" ]]; then
    skip "Meshery server is not running, skipping test"
  fi
  
  assert_success
  # Verify output format - should contain headers for designs
  assert_output --partial "DESIGN ID"
}

@test "mesheryctl design list shows total count" {
  # Run list to check for total count display
  run $MESHERYCTL_BIN design list
  echo "List command output: $output"
  
  # Skip if Meshery server is not running
  if [[ "$output" =~ "connection refused" ]]; then
    skip "Meshery server is not running, skipping test"
  fi
  
  assert_success
  # Check for the total count line (adjust to match actual output format)
  assert_output --partial "Total number of  patterns :"
}

@test "mesheryctl design list with page parameter works" {
  # Test pagination with page flag
  run $MESHERYCTL_BIN design list --page 1
  echo "List command output with page=1: $output"
  
  # Skip if Meshery server is not running
  if [[ "$output" =~ "connection refused" ]]; then
    skip "Meshery server is not running, skipping test"
  fi
  
  assert_success
  # Adjust to match actual output format
  assert_output --partial "Total number of  patterns :"
  assert_output --partial "DESIGN ID"
}

@test "mesheryctl design list with verbose flag works" {
  # Test verbose output
  run $MESHERYCTL_BIN design list -v
  echo "Verbose list command output: $output"
  
  # Skip if Meshery server is not running
  if [[ "$output" =~ "connection refused" ]]; then
    skip "Meshery server is not running, skipping test"
  fi
  
  assert_success
  # Should either show designs or indicate no patterns found
  if [[ "$output" == *"No pattern(s) found"* ]]; then
    assert_output --partial "No pattern(s) found"
  else
    assert_output --partial "DESIGN ID"
  fi
}