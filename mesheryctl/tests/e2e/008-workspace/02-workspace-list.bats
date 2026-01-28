#!/usr/bin/env bats
setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    load "$E2E_HELPERS_PATH/constants"
    _load_bats_libraries

}


@test "given missing --orgId flag when running mesheryctl exp workspace list then it fails displaying error message" {
    run $MESHERYCTL_BIN exp workspace list

    assert_failure
    assert_output --partial "Error"
    assert_output --partial "isn't specified"
}

@test "given no orgId is provided as an argument when running mesheryctl exp workspace list --orgId then the error message is displayed" {
    run $MESHERYCTL_BIN exp workspace list --orgId
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "needs an argument"
}

# todo change output
@test "given an invalid orgId is provided as an argument when running mesheryctl exp workspace list --orgId then the error message is displayed" {
    CONNECTION_ID="foo"

    run $MESHERYCTL_BIN exp workspace list --orgId "$CONNECTION_ID"
    assert_failure
}

@test "given non-existent orgId is provided as an argument when running mesheryctl exp workspace list --orgId then the error message is displayed" {
    CONNECTION_ID="00000000-0000-0000-0000-000000000000"

    run $MESHERYCTL_BIN exp workspace list --orgId "$CONNECTION_ID"
    assert_success
    assert_output --partial "No workspaces found"
}