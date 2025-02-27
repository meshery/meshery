#!/usr/bin/env bats

setup() {
    load "$HELPERS_LIBS/bats-support/load"
    load "$HELPERS_LIBS/bats-assert/load"
    load "$HELPERS_LIBS/bats-detik/lib/util"
    load "$HELPERS_LIBS/bats-detik/lib/detik"

    DETIK_CLIENT_NAMESPACE=meshery
}

# starting server on kubernetes platform

@test "mesheryctl system start on kubernetes is succeeded" {
    run $MESHERYCTL_BIN system start -p kubernetes
    [ "$status" -eq 0 ]

    # wait for deployment
    sleep 40

    run verify "there are 1 pods named 'meshery'"
	[ "$status" -eq 0 ]
	
	run verify "there is 1 service named 'meshery'"
	[ "$status" -eq 0 ]
}

