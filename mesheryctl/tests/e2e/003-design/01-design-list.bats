#!/usr/bin/env bats

# Load required BATS libraries
load '../helpers/bats-support/load'
load '../helpers/bats-assert/load'

# Test 1: Verify basic design list command succeeds
@test "mesheryctl design list is succeeded" {
  run $MESHERYCTL_BIN design list
  assert_success
  assert_output --partial "DESIGN ID"
}

# Test 2: Verify design list with page parameter succeeds
@test "mesheryctl design list with page parameter is succeeded" {
  run $MESHERYCTL_BIN design list --page 1
  assert_success
  # Fix: Match the exact output format with two spaces between "of" and "patterns"
  assert_output --partial "Total number of  patterns :"
}

# Test 3: Verify design list with verbose flag succeeds
@test "mesheryctl design list with verbose flag is succeeded" {
  run $MESHERYCTL_BIN design list -v
  assert_success
  # Check for either designs or "No pattern(s) found" message
  if [[ "$output" == *"No pattern(s) found"* ]]; then
    assert_output --partial "No pattern(s) found"
  else
    assert_output --partial "DESIGN ID"
  fi
}