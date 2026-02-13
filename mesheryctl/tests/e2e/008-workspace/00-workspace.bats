#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    load "$E2E_HELPERS_PATH/constants"
    _load_bats_libraries
}

@test "given --help flag provided when running mesheryctl exp workspace --help displays command usage" {
    run $MESHERYCTL_BIN exp workspace --help

    assert_success
    assert_output --partial "Examples:"    
    assert_output --partial "Available Commands:"    
}
