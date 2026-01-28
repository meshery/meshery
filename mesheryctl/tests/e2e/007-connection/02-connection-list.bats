setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"
}

@test "given all requirements met when mesheryctl connection list then header of total number of connections followed by a list are displayed" {
    run $MESHERYCTL_BIN connection list

    assert_success
    assert_output --partial "Total number of connection"
}

@test "given all requirements met when running mesheryctl connection list --count then only the total number of available connections is displayed" {
    run $MESHERYCTL_BIN connection list --count

    assert_success
    assert_output --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}