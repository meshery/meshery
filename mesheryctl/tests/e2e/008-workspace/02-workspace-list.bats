#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    load "$E2E_HELPERS_PATH/constants"
    _load_bats_libraries
}

@test "given missing --orgId flag when running mesheryctl workspace list then an error message is displayed" {
    run $MESHERYCTL_BIN workspace list

    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid value for --orgId ''"
}

@test "given an invalid orgId is provided as an argument when running mesheryctl workspace list --orgId invalid-orgid then the error message is displayed" {
    ORGANIZATION_ID="foo"

    run $MESHERYCTL_BIN workspace list --orgId "$ORGANIZATION_ID"
    assert_failure

    assert_output --partial "Error"
    assert_output --partial "Invalid value for --orgId 'foo': must be a valid UUID"
}

@test "given non-existent orgId provided when running mesheryctl workspace list --orgId non-existent-orgId then an error message is displayed" {
    NON_EXISTENT_ORGANIZATION_ID="00000000-0000-0000-0000-000000000000"

    run $MESHERYCTL_BIN workspace list --orgId "$NON_EXISTENT_ORGANIZATION_ID"
    if [ "$status" -eq 0 ]; then
        assert_output --partial "No workspaces found"
    else
        assert_output --partial "Meshery server"
    fi
}

@test "given a valid orgId when running mesheryctl workspace list then it should pass validation" {
    run $MESHERYCTL_BIN workspace list --orgId "$ORGANIZATION_ID" --count
    
    refute_output --partial "must be a valid UUID"
    
    if [ "$status" -eq 0 ]; then
        assert_output --partial "Total number of workspaces"
    else
        assert_output --partial "Unable to connect"
    fi
}