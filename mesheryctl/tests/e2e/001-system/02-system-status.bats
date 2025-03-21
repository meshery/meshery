#!/usr/bin/env bats

@test "mesheryctl system status is succeeded" {
    run $MESHERYCTL_BIN system status -y
    [ "$status" -eq 0 ]
}
