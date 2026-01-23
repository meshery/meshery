#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/registry"
    mkdir -p "$TESTDATA_DIR"  
    export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
    
    # add CREDS & ID for integration sheets here
    export TEST_SPREADSHEET_ID=""
    export TEST_SPREADSHEET_CRED=""
}

common_outputs_on_success() {
    assert_output --partial "Generated"
    assert_output --partial "Summary:"
    assert_output --partial "generated."
}

@test "mesheryctl registry generate displays usage instructions when no arguments are provided" {
    run $MESHERYCTL_BIN registry generate 
    assert_failure
    assert_output --partial "[ Spreadsheet ID | Registrant Connection Definition Path | Local Directory | Individual CSV files ] isn't specified"
    assert_output --partial "Usage:"
    assert_output --partial "mesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred \$CRED"
}

@test "mesheryctl registry generate fails when spreadsheet-id is provided without spreadsheet-cred" {
    run $MESHERYCTL_BIN registry generate --spreadsheet-id "test-id"
    assert_failure
    assert_output --partial "Spreadsheet Credentials is required"
}

@test "mesheryctl registry generate fails when registrant-def is provided without registrant-cred" {
    run $MESHERYCTL_BIN registry generate --registrant-def "$FIXTURES_DIR/connection-def.json"
    assert_failure
    assert_output --partial "Registrant Credentials is required"
}

@test "mesheryctl registry generate fails with invalid directory" {
    run $MESHERYCTL_BIN registry generate --directory "invalid-directory"
    assert_failure
    assert_output --partial "error reading the directory"
}

@test "mesheryctl registry generate fails with directory missing required CSV files" {
    run $MESHERYCTL_BIN registry generate --directory "$FIXTURES_DIR/incomplete-csv-dir"
    assert_failure
    assert_output --partial "error reading the directory"
    assert_output --partial "either the model csv or component csv is missing"
}

@test "mesheryctl registry generate supports model-specific generation with spreadsheet" {
    run $MESHERYCTL_BIN registry generate --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --model "flyte"
    assert_success
    common_outputs_on_success
}

@test "mesheryctl registry generate succeeds with valid spreadsheet credentials" {
    run $MESHERYCTL_BIN registry generate --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED"
    assert_success
    common_outputs_on_success
}
