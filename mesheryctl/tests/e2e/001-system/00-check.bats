#!/usr/bin/env bats

@test "mesheryctl system status is succeeded" {
    run $MESHERYCTL_BIN system status -y
    [ "$status" -eq 0 ]
}

# @test "mesheryctl system check is succeeded" {
#    run $MESHERYCTL_BIN system check
#    [ "$status" -eq 0 ]
# }