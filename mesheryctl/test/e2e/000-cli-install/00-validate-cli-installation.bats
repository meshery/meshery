#!/usr/bin/env bats

# Basic tests to validate cli has been build and behave properly at root

@test "cli is available" {
    run $MESHERYCTL_BIN
    [ "$status" -eq 0 ]
}

@test "mesheryctl version return Client and Server" {
  actual_output=$($MESHERYCTL_BIN version | grep -E "Client|Server" | wc -l)
  [[ "$tatus" -eq 0 ]]
  [[ "$actual_output" -eq 2 ]]
}