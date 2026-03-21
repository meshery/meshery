#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    # allow CI/environment to provide integration spreadsheet credentials
    export TEST_SPREADSHEET_ID="${TEST_SPREADSHEET_ID:-}"
    export TEST_SPREADSHEET_CRED="${TEST_SPREADSHEET_CRED:-}"
}

require_spreadsheet_credentials() {
    if [[ -z "$TEST_SPREADSHEET_ID" || -z "$TEST_SPREADSHEET_CRED" ]]; then
        skip "Spreadsheet credentials not configured"
    fi
}

@test "given no arguments when running mesheryctl registry update then an error is displayed" {
    run $MESHERYCTL_BIN registry update
    assert_failure
    assert_output --partial "error updating registry"
}

@test "given spreadsheet-id without spreadsheet-cred when running mesheryctl registry update then an error about missing flag is displayed" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "test-id"
    assert_failure
    assert_output --partial "if any flags in the group [spreadsheet-id spreadsheet-cred] are set they must all be set"
    assert_output --partial "spreadsheet-cred"
}

@test "given spreadsheet-cred without spreadsheet-id when running mesheryctl registry update then an error about missing flag is displayed" {
    run $MESHERYCTL_BIN registry update --spreadsheet-cred "test-cred"
    assert_failure
    assert_output --partial "if any flags in the group [spreadsheet-id spreadsheet-cred] are set they must all be set"
    assert_output --partial "spreadsheet-id"
}

@test "given invalid spreadsheet credentials when running mesheryctl registry update then an error is displayed" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "invalid-id" --spreadsheet-cred "invalid-cred"
    assert_failure
    assert_output --partial "Invalid JWT credentials"
}

@test "given an invalid model name when running mesheryctl registry update then zero models are updated" {
    require_spreadsheet_credentials
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --model "nonexistent-model"
    assert_success
    assert_output --partial "Updated 0 models and 0 components"
}

@test "given valid spreadsheet credentials and a model name when running mesheryctl registry update then that model is updated" {
    require_spreadsheet_credentials
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --model "kubernetes"
    assert_success
    assert_output --regexp "Updated [1-9][0-9]* models? and [0-9]+ components?"
}

@test "given valid spreadsheet credentials when running mesheryctl registry update then models are updated successfully" {
    require_spreadsheet_credentials
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED"
    assert_success
    assert_output --regexp "Updated [1-9][0-9]* models? and [0-9]+ components?"
}
