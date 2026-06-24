#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
}

@test "given help flag when running mesheryctl system config aks --help then usage information is displayed" {
    run $MESHERYCTL_BIN system config aks --help

    assert_success
    assert_output --partial "Configure Meshery to connect to AKS cluster"
    assert_output --partial "mesheryctl system config aks --token auth.json"
}

@test "given more than one argument when running mesheryctl system config aks then an error message is displayed" {
    run $MESHERYCTL_BIN system config aks extra-arg

    assert_failure
    assert_output --partial "more than one config name provided"
}

@test "given azure CLI is not installed when running mesheryctl system config aks then an error is displayed" {
    PATH=/dev/null run $MESHERYCTL_BIN system config aks

    assert_failure
}
