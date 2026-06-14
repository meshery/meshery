#!/usr/bin/env bats

# tests to ensure the server is up and running

setup() {
	load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

    DETIK_CLIENT_NAMESPACE="meshery"
    DETIK_CLIENT_NAME="kubectl"
}

@test "meshery pod is deployed" {
	run verify "there are more than 0 pod named '^meshery-[a-z0-9]+-[a-z0-9]+$'"
	assert_success
}

@test "meshery-broker pod is deployed" {
	run verify "there are more than 0 pod named '^meshery-broker-[0-9]+$'"
	assert_success
}

@test "meshery-meshsync pod is deployed" {
    run verify "there are more than 0 pod named '^meshery-meshsync-[a-z0-9]+-[a-z0-9]+$'"
	assert_success
}

@test "meshery-operator pod is deployed" {
	run verify "there are more than 0 pod named '^meshery-operator-[a-z0-9]+-[a-z0-9]+$'"
	assert_success
}


@test "meshery service is deployed" {
	run verify "there is 1 service named '^meshery$'"
	assert_success
}

@test "meshery-broker service is deployed" {
	run verify "there is 1 service named '^meshery-broker$'"
	assert_success
}

@test "meshery-operator service is deployed" {
	run verify "there is 1 service named '^meshery-operator$'"
	assert_success
}
