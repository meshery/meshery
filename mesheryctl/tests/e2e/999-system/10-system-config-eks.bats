#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
}

@test "given help flag when running mesheryctl system config eks --help then usage information is displayed" {
    run $MESHERYCTL_BIN system config eks --help

    assert_success
    assert_output --partial "Configure Meshery to connect to EKS cluster"
    assert_output --partial "mesheryctl system config eks --token auth.json"
}

@test "given more than one argument when running mesheryctl system config eks then an error message is displayed" {
    run $MESHERYCTL_BIN system config eks extra-arg

    assert_failure
    assert_output --partial "more than one config name provided"
}

@test "given aws CLI is not installed when running mesheryctl system config eks then an error is displayed" {
    PATH=/dev/null run $MESHERYCTL_BIN system config eks

    assert_failure
}
