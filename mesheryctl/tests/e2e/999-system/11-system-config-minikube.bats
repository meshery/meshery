#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
}

@test "given help flag when running mesheryctl system config minikube --help then usage information is displayed" {
    run $MESHERYCTL_BIN system config minikube --help

    assert_success
    assert_output --partial "Configure Meshery to connect to minikube cluster"
    assert_output --partial "mesheryctl system config minikube --token auth.json"
}

@test "given more than one argument when running mesheryctl system config minikube then an error message is displayed" {
    run $MESHERYCTL_BIN system config minikube extra-arg

    assert_failure
    assert_output --partial "more than one config name provided"
}
