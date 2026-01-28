#!/usr/bin/env bats

setup() {

    load "$E2E_HELPERS_PATH/bats_libraries"
    load "$E2E_HELPERS_PATH/constants"
	_load_bats_libraries

}

@test "given no orgId is provided as an argument when running mesheryctl exp workspace create --orgId then the error message is displayed" {
    run $MESHERYCTL_BIN exp workspace create --orgId

    assert_failure 
    assert_output --partial "Error"
    assert_output --partial "needs an argument"
}

@test "given invalid orgId is provided as an argument when running mesheryctl exp workspace create --orgId --name --description then the error message is displayed" {
    run $MESHERYCTL_BIN exp workspace create --orgId foo --name name --description description

    assert_failure 
    assert_output --partial "Error"
    assert_output --partial "provide a valid organization ID"
}
