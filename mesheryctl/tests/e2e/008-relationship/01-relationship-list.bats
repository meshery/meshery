#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"
}

@test "Count components" {
    run $MESHERYCTL_BIN exp relationship list --count
    assert_success
    assert_output --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}

@test "component list shows header and pagination info" {
    run $MESHERYCTL_BIN exp relationship list --page 1
    assert_success
    assert_output --partial $LIST_COMMAND_OUTPUT_REGEX_PATTERN
    assert_output --partial "KIND"
    assert_output --partial "MODEL"
}