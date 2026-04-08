#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    load "$E2E_HELPERS_PATH/constants"
    _load_bats_libraries
}

@test "given no argument provided when running mesheryctl workspace view then an error message is displayed" {
    run $MESHERYCTL_BIN workspace view

    assert_failure
    assert_output --partial "Error"
    assert_output --partial "please provide exactly one workspace name or ID"
}

@test "given a name without --orgId when running mesheryctl workspace view foo then an error message is displayed" {
    run $MESHERYCTL_BIN workspace view foo

    assert_failure
    assert_output --partial "Error"
    assert_output --partial "--orgId is required when searching by name"
}

@test "given an invalid orgId when running mesheryctl workspace view somename --orgId foo then an error message is displayed" {
    run $MESHERYCTL_BIN workspace view somename --orgId foo

    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid value for --orgId 'foo': must be a valid UUID"
}

@test "given non-existent UUID when running mesheryctl workspace view then workspace not found error is displayed" {
    "TODO: workspace view by UUID needs rework - tracked in #18490"
    run $MESHERYCTL_BIN workspace view 12345678-1234-1234-1234-123456789012

    assert_failure
    assert_output --partial "Error"
    assert_output --partial "WorkspaceID or OrgID cannot be empty"
}
