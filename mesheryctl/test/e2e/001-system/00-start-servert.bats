#!/usr/bin/env bats

# starting server on kubernetes platform

@test "mesheryctl system start on kubernetes is succeeded" {
    run $MESHERYCTL_BIN system start -p kubernetes
    [ "$status" -eq 0 ]
}

