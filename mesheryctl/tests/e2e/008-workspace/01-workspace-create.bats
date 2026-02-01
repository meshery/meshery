#!/usr/bin/env bats

setup() {

    load "$E2E_HELPERS_PATH/bats_libraries"
    load "$E2E_HELPERS_PATH/constants"
	_load_bats_libraries

}

@test "given --orgId without value is provided as an argument when running mesheryctl exp workspace create --orgId then an error message is displayed" {
    run $MESHERYCTL_BIN exp workspace create --orgId

    assert_failure 
    assert_output --partial "Error"
    assert_output --partial "Invalid Argument"
}

@test "given an invalid orgId provided as an argument when running mesheryctl exp workspace create --orgId invalid-org-id then an error message is displayed" {
    run $MESHERYCTL_BIN exp workspace create --orgId foo

    assert_failure 
    assert_output --partial "Error"
    assert_output --partial "Invalid Argument"
    assert_output --partial "[ name | description ] not specified"
}

@test "given an invalid orgId provided when running mesheryctl exp workspace create --orgId invalid-org-id --name --description then an error message is displayed" {
    run $MESHERYCTL_BIN exp workspace create --orgId foo --name name --description description

    assert_failure 
    assert_output --partial "Error"
    assert_output --partial "Invalid Argument"
    assert_output --partial "provide a valid organization ID"
}
@test "given an invalid orgId and missing of --description flag when running mesheryctl exp workspace create --orgId invalid-org-id --name then an error message is displayed" {
    run $MESHERYCTL_BIN exp workspace create --orgId foo --name name

    assert_failure 
    assert_output --partial "Error"
    assert_output --partial "Invalid Argument"
    assert_output --partial "[ description ] not specified"
}
@test "given an invalid orgId and missing of --name flag when running mesheryctl exp workspace create --orgId invalid-org-id --name --description then an error message is displayed" {
    run $MESHERYCTL_BIN exp workspace create --orgId foo --description description

    assert_failure 
    assert_output --partial "Error"
    assert_output --partial "Invalid Argument"
    assert_output --partial "[ name ] not specified"
}
