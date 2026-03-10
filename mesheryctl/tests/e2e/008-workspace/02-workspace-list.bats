#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    load "$E2E_HELPERS_PATH/constants"
    _load_bats_libraries
}


@test "given missing --orgId flag when running mesheryctl exp workspace list then an error message is displayed" {
    run $MESHERYCTL_BIN exp workspace list

    assert_failure
    assert_output --partial "Error"
    assert_output --partial "[ Organization ID ] isn't specified"
}

# todo change output
@test "given an invalid orgId is provided as an argument when running mesheryctl exp workspace list --orgId invalid-orgid then the error message is displayed" {
    ORGANIZATION_ID="foo"

    run $MESHERYCTL_BIN exp workspace list --orgId "$ORGANIZATION_ID"
    assert_failure

    assert_output --partial "Error"
    assert_output --partial "Server emitted an error"
}

@test "given non-existent orgId provided when running mesheryctl exp workspace list --orgId non-existent-orgId then an error message is displayed" {
    NON_EXISTENT_ORGANIZATION_ID="00000000-0000-0000-0000-000000000000"

    run $MESHERYCTL_BIN exp workspace list --orgId "$NON_EXISTENT_ORGANIZATION_ID"
    assert_success
    assert_output --partial "No workspaces found"
}