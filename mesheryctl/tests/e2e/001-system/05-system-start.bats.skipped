#!/usr/bin/env bats

setup(){
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    run $MESHERYCTL_BIN system stop --yes

    export START_HEADER="Start Meshery and each of its cloud native components"
    export START_USAGE="mesheryctl system start [flags]"
    export START_BASIC="mesheryctl system start"
    export START_SKIP_BROWSER="mesheryctl system start --skip-browser"
    export START_SKIP_UPDATE="mesheryctl system start --skip-update"
    export START_RESET="mesheryctl system start --reset"
    export START_PLATFORM_DOCKER="mesheryctl system start -p docker"
    export START_PROVIDER="mesheryctl system start --provider Meshery"
    export START_MESSAGE_DOCKER="Starting Meshery..."
}

teardown() {
    run $MESHERYCTL_BIN system stop --force -y
}

common_start_assertions() {
    assert_success 
    assert_output --partial "$START_MESSAGE_DOCKER"
}

common_kube_assertions() {
    if ! kubectl cluster-info >/dev/null 2>&1; then
        skip "Kubernetes cluster not accessible"
    fi
}

@test "mesheryctl system start --help displays help" {
    run $MESHERYCTL_BIN system start --help
    assert_success
    assert_output --partial "$START_HEADER"
    assert_output --partial "$START_USAGE"
    assert_output --partial "$START_BASIC"
    assert_output --partial "$START_SKIP_BROWSER"
    assert_output --partial "$START_SKIP_UPDATE"
    assert_output --partial "$START_RESET"
    assert_output --partial "$START_PLATFORM_DOCKER"
    assert_output --partial "$START_PROVIDER"
}

@test "mesheryctl system start -p invalidplatform displays error for invalid platform" {
    run $MESHERYCTL_BIN system start -p invalidplatform
    assert_failure
    assert_output --partial "$INVALID_PLATFORM_REGEX"
}

@test "mesheryctl system start --skip-browser succeeds" {
    run $MESHERYCTL_BIN system start --skip-browser -y
    common_start_assertions
}

@test "mesheryctl system start --skip-update succeeds" {
    run $MESHERYCTL_BIN system start --skip-update -p docker
    assert_success
    assert_output --partial "$SKIP_UPDATE_REGEX"
    assert_success
    assert_output --partial "$DOCKER_PS_REGEX"
    run $MESHERYCTL_BIN system stop -y
}

@test "mesheryctl system start --reset succeeds" {
    run $MESHERYCTL_BIN system start --reset -y
    common_start_assertions
}

@test "mesheryctl system start -p docker succeeds" {
    run $MESHERYCTL_BIN system start -p docker -y
    common_start_assertions
}

@test "mesheryctl system start -p kubernetes succeeds" {
    common_kube_assertions
    run $MESHERYCTL_BIN system start -p kubernetes 
    common_start_assertions
}

@test "mesheryctl system start --provider succeeds" {
    run $MESHERYCTL_BIN system start --provider Meshery 
    common_start_assertions
}

@test "mesheryctl system start with all flags on docker succeeds" {
    run $MESHERYCTL_BIN system start --skip-browser --skip-update --reset -p docker -y
    common_start_assertions
}

@test "mesheryctl system start with all flags on kubernetes succeeds" {
    common_kube_assertions
    run $MESHERYCTL_BIN system start --skip-browser --skip-update --reset -p kubernetes 
    common_start_assertions
}