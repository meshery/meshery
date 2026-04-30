#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
}

@test "given no subcommand when running mesheryctl system config then an error message is displayed" {
    run $MESHERYCTL_BIN system config

    assert_failure
    assert_output --partial "name of kubernetes cluster to configure Meshery not provided"
}

@test "given an invalid subcommand when running mesheryctl system config then an error message is displayed" {
    run $MESHERYCTL_BIN system config invalidcloud

    assert_failure
    assert_output --partial "invalid command: \"invalidcloud\""
}

@test "given help flag when running mesheryctl system config --help then deprecation warning is displayed" {
    run $MESHERYCTL_BIN system config --help

    assert_success
    assert_output --partial "Command \"config\" is deprecated"
    assert_output --partial "Please use 'mesheryctl connection create --type <k8s-type>' instead"
}

@test "given help flag when running mesheryctl system config gke --help then usage information is displayed" {
    run $MESHERYCTL_BIN system config gke --help

    assert_success
    assert_output --partial "Configure Meshery to connect to GKE cluster"
    assert_output --partial "mesheryctl system config gke --token auth.json"
}

@test "given more than one argument when running mesheryctl system config gke then an error message is displayed" {
    run $MESHERYCTL_BIN system config gke extra-arg

    assert_failure
    assert_output --partial "more than one config name provided"
}

@test "given no GKE cluster credentials when running mesheryctl system config gke then an error is displayed" {
    PATH=/dev/null run $MESHERYCTL_BIN system config gke

    assert_failure
}
