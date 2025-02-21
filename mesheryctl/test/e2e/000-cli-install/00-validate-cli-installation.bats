#!/usr/bin/env bats

# Basic tests to validate cli has been build and behave properly at root

@test "assert cli is available" {
    run $MESHERYCTL_BIN
    [ "$status" -eq 0 ]
}

@test "test help flag header output" {
  expected_output="As a self-service engineering platform, Meshery enables collaborative design and operation of cloud native infrastructure."
  actual_output=$($MESHERYCTL_BIN --help | head -n 1)
  [[ "$tatus" -eq 0 ]]
  [[ "$actual_output" = "$expected_output" ]]
}