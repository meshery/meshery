#!/usr/bin/env bats

# tests to ensure the server is up and running

setup() {
	load "$E2E_HELPERS_PATH/tests_helpers"
	_tests_helper

    DETIK_CLIENT_NAMESPACE="meshery"
    DETIK_CLIENT_NAME="kubectl"
}

@test "meshery pod is deployed" {
	run verify "there is 1 pod named '^meshery-[a-z0-9]{10}-[a-z0-9]{5}$'"
	assert_success
}

@test "meshery-broker pod is deployed" {
	run verify "there is 1 pod named '^meshery-broker-[0-9]+$'"
	assert_success
}

@test "meshery-meshsync pod is deployed" {
    run verify "there is 1 pod named '^meshery-meshsync-[a-z0-9]{10}-[a-z0-9]{5}$'"
	assert_success
}

@test "meshery-operator pod is deployed" {
	run verify "there is 1 pod named '^meshery-operator-[a-z0-9]{10}-[a-z0-9]{5}$'"
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
