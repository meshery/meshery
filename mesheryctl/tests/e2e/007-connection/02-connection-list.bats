setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"
}

@test "given all requirements met when running mesheryctl connection list --status status then the conections of specified status is displayed" {
    run $MESHERYCTL_BIN connection list --status connected --page 1

    assert_success
    assert_output --partial "Total number of connection"
}

@test "given all requirements met when running mesheryctl connection list --kind kind-name then the connections of specified kind is displayed" {
    run $MESHERYCTL_BIN connection list --kind kubernetes --page 1

    assert_success
    assert_output --partial "Total number of connection"
}

@test "given all requirements met when running mesheryctl connection list --pagesize size then the conections of specified size is displayed" {
    run $MESHERYCTL_BIN connection list --pagesize 1 --page 1

    assert_success
    assert_output --partial "Total number of connection"
}

@test "given all requirements met when running mesheryctl connection list then header of total number of connections followed by a list are displayed" {
    run $MESHERYCTL_BIN connection list --page 1

    assert_success
    assert_output --partial "Total number of connection"
}

@test "given all requirements met when running mesheryctl connection list --count then only the total number of available connections is displayed" {
    run $MESHERYCTL_BIN connection list --count

    assert_success
    assert_output --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}