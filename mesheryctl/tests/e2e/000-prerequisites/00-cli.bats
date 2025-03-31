#!/usr/bin/env bats

# Basic tests to validate cli has been build and behave properly at root

@test "cli is available" {
    if [[ -z "$MESHERYCTL_BIN" ]]; then
        echo "Error: MESHERYCTL_BIN is not defined. Set it before running tests."
        exit 1
    fi
    
    run $MESHERYCTL_BIN
    [ "$status" -eq 0 ]
}

@test "mesheryctl version return Client and Server" {
    run $MESHERYCTL_BIN version

    [[ "$status" -eq 0 ]] 
    actual_output=$(echo "$output" | grep -E "Client|Server" | wc -l)
    [[ "$actual_output" -eq 2 ]]
}