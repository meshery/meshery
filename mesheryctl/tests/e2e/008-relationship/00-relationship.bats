#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"
}

@test "Help for components list command" {
    run $MESHERYCTL_BIN exp relationship --help
    assert_success
    assert_output --partial "Available Commands:"
    assert_output --partial "Usage:"
}

@test "Count components" {
    run $MESHERYCTL_BIN exp relationship --count
    assert_success
    assert_output --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}
