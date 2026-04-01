#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"
}

@test "given all requirements met when running  mesheryctl connection --count then total number of available connections is displayed" {
    run $MESHERYCTL_BIN connection --count
    assert_success

    assert_output --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}