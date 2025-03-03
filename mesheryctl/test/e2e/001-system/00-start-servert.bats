#!/usr/bin/env bats

setup() {
    load "$SUPPORT_DESTDIR/load"
    load "$ASSERT_DESTDIR/load"
    load "$DETIK_DESTDIR/detik.bash"

    DETIK_CLIENT_NAMESPACE="meshery"
    DETIK_CLIENT_NAME="kubectl"
}

@test "mesheryctl system start on kubernetes is succeeded" {
    run $MESHERYCTL_BIN system start -p kubernetes
    [ "$status" -eq 0 ]

    # give time to pods and services to be deployed to prevent flaky test
    sleep 40
}

@test "meshery pods are deployed" {
	run verify "there is 1 pod named '^meshery-[a-z0-9]{10}-[a-z0-9]{5}$'"
	[ "$status" -eq 0 ]

	run verify "there is 1 pod named '^meshery-broker-[0-9]+$'"
	[ "$status" -eq 0 ]
    
    run verify "there is 1 pod named '^meshery-meshsync-[a-z0-9]{10}-[a-z0-9]{5}$'"
	[ "$status" -eq 0 ]

	run verify "there is 1 pod named '^meshery-operator-[a-z0-9]{10}-[a-z0-9]{5}$'"
	[ "$status" -eq 0 ]
}


@test "meshery services are deployed" {
	run verify "there is 1 service named '^meshery$'"
	[ "$status" -eq 0 ]

	run verify "there is 1 service named '^meshery-broker$'"
	[ "$status" -eq 0 ]


	run verify "there is 1 service named '^meshery-operator$'"
	[ "$status" -eq 0 ]
}