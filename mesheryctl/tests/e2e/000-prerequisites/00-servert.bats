#!/usr/bin/env bats

# tests to ensure the server is up and running

setup() {
    load "$SUPPORT_DESTDIR"
    load "$ASSERT_DESTDIR"
    load "$DETIK_DESTDIR"

    DETIK_CLIENT_NAMESPACE="meshery"
    DETIK_CLIENT_NAME="kubectl"
}

@test "meshery server is running pods are deployed" {
	run verify "there is 1 pod named '^meshery-[a-z0-9]{10}-[a-z0-9]{5}$'"
	[ "$status" -eq 0 ]

	run verify "there is 1 pod named '^meshery-broker-[0-9]+$'"
	[ "$status" -eq 0 ]
    
    run verify "there is 1 pod named '^meshery-meshsync-[a-z0-9]{10}-[a-z0-9]{5}$'"
	[ "$status" -eq 0 ]

	run verify "there is 1 pod named '^meshery-operator-[a-z0-9]{10}-[a-z0-9]{5}$'"
	[ "$status" -eq 0 ]
}


@test "meshery server is running services are deployed" {
	run verify "there is 1 service named '^meshery$'"
	[ "$status" -eq 0 ]

	run verify "there is 1 service named '^meshery-broker$'"
	[ "$status" -eq 0 ]


	run verify "there is 1 service named '^meshery-operator$'"
	[ "$status" -eq 0 ]
}