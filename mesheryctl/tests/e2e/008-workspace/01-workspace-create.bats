#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    load "$E2E_HELPERS_PATH/constants"
	_load_bats_libraries
}

@test "given an invalid orgId provided as an argument when running mesheryctl exp workspace create --orgId invalid-org-id then an error message is displayed" {
    run $MESHERYCTL_BIN exp workspace create --orgId foo

    assert_failure 
    assert_output --partial "Error"
    assert_output --partial "Invalid Argument"
    assert_output --partial "[ name | description ] not specified"
}

