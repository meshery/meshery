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

# Meshery Operator >= 1.0.0 renders the broker from the official NATS chart, so
# the workload is meshery-nats (pod meshery-nats-0); earlier operators used
# meshery-broker. Accept either, because the operator version is a property of
# the cluster under test, not of this suite. Only the workload was renamed - the
# Broker custom resource is still meshery-broker.
@test "meshery-broker pod is deployed" {
	run verify "there are more than 0 pod named '^meshery-(nats|broker)-[0-9]+$'"
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

# See the pod test above for why both names are accepted. The NATS chart also
# publishes a meshery-nats-headless service; the anchored pattern matches only
# the addressable one, so the count stays 1.
@test "meshery-broker service is deployed" {
	run verify "there is 1 service named '^meshery-(nats|broker)$'"
	assert_success
}

@test "meshery-operator service is deployed" {
	run verify "there is 1 service named '^meshery-operator$'"
	assert_success
}
