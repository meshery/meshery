#!/usr/bin/env bats

setup(){
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    export START_HEADER="Start Meshery and each of its cloud native components"
    export START_USAGE="mesheryctl system start [flags]"
    export START_BASIC="mesheryctl system start"
    export START_SKIP_BROWSER="mesheryctl system start --skip-browser"
    export START_SKIP_UPDATE="mesheryctl system start --skip-update"
    export START_RESET="mesheryctl system start --reset"
    export START_PLATFORM_DOCKER="mesheryctl system start -p docker"
    export START_PROVIDER="mesheryctl system start --provider Meshery"
    export START_MESSAGE_DOCKER="Starting Meshery..."
    export START_MESSAGE_KUBE="Meshery is starting..."
    export SKIP_UPDATE_MESSAGE="Skipping Meshery update..."
}

common_start_assertions() {
    assert_success 
    assert_output --partial "$START_MESSAGE_DOCKER"
}

common_kube_assertions() {
    assert_success
    assert_output --partial "$START_MESSAGE_KUBE"
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
    assert_output --partial "the platform invalidplatform is not supported currently"
    assert_output --partial "the supported platforms are:"
    assert_output --partial "docker"
    assert_output --partial "kubernetes"
}

@test "mesheryctl system start --skip-browser succeeds" {
    run $MESHERYCTL_BIN system start --skip-browser
    common_start_assertions
}

@test "mesheryctl system start --skip-update succeeds" {
    run $MESHERYCTL_BIN system start --skip-update -p docker
    assert_success
    assert_output --partial "$SKIP_UPDATE_MESSAGE"
    assert_output --partial "$START_MESSAGE_DOCKER"  
}

@test "mesheryctl system start --reset succeeds" {
    run $MESHERYCTL_BIN system start --reset
    common_start_assertions
}

@test "mesheryctl system start -p docker succeeds" {
    run $MESHERYCTL_BIN system start -p docker -y
    common_start_assertions
    run docker ps --filter name=meshery_meshery_1 --format "{{.Status}}"
    assert_output --partial "Up"
}

@test "mesheryctl system start -p kubernetes succeeds" {
    run $MESHERYCTL_BIN system start -p kubernetes 
    common_kube_assertions
}

@test "mesheryctl system start --provider succeeds" {
    run $MESHERYCTL_BIN system start --provider Meshery 
    common_start_assertions
}

@test "mesheryctl system start with all flags on docker succeeds" {
    run $MESHERYCTL_BIN system start --skip-browser --skip-update --reset -p docker -y
    common_start_assertions
    run docker ps --filter name=meshery_meshery_1 --format "{{.Status}}"
    assert_output --partial "Up"
}

@test "mesheryctl system start with all flags on kubernetes succeeds" {
    run $MESHERYCTL_BIN system start --skip-browser --skip-update --reset -p kubernetes 
    common_kube_assertions
}