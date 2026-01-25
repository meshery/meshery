setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"
}

@test "mesheryctl connection list runs successfully" {
    run $MESHERYCTL_BIN connection list

    assert_success
    assert_output --partial "Total number of connection"
}

@test "mesheryctl connection list --count returns only count" {
    run $MESHERYCTL_BIN connection list --count

    assert_success
    assert_output --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}