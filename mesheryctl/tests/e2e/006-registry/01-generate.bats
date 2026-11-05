#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/registry"
    mkdir -p "$TESTDATA_DIR"  
    export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
    
    # allow CI/environment to provide integration spreadsheet credentials
    export TEST_SPREADSHEET_ID="${TEST_SPREADSHEET_ID:-}"
    export TEST_SPREADSHEET_CRED="${TEST_SPREADSHEET_CRED:-}"
}

common_outputs_on_success() {
    assert_output --partial "Generated"
    assert_output --partial "Summary:"
    assert_output --partial "generated."
}

require_spreadsheet_credentials() {
    if [[ -z "$TEST_SPREADSHEET_ID" || -z "$TEST_SPREADSHEET_CRED" ]]; then
        skip "Spreadsheet credentials not configured"
    fi
}

@test "given no arguments when running mesheryctl registry generate then usage instructions are displayed" {
    run $MESHERYCTL_BIN registry generate 
    assert_failure
    assert_output --partial "[ Spreadsheet ID | Registrant Connection Definition Path | Local Directory | Individual CSV files ] isn't specified"
    assert_output --partial "Usage:"
    assert_output --partial "mesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred \$CRED"
}

@test "given spreadsheet-id without spreadsheet-cred when running mesheryctl registry generate then an error message is displayed" {
    run $MESHERYCTL_BIN registry generate --spreadsheet-id "test-id"
    assert_failure
    assert_output --partial "Spreadsheet Credentials is required"
}

@test "given registrant-def without registrant-cred when running mesheryctl registry generate then an error message is displayed" {
    run $MESHERYCTL_BIN registry generate --registrant-def "$FIXTURES_DIR/connection-def.json"
    assert_failure
    assert_output --partial "Registrant Credentials is required"
}

@test "given an invalid directory when running mesheryctl registry generate then an error message is displayed" {
    run $MESHERYCTL_BIN registry generate --directory "invalid-directory"
    assert_failure
    assert_output --partial "error reading the directory"
}

@test "given a directory missing required csv files when running mesheryctl registry generate then an error message is displayed" {
    run $MESHERYCTL_BIN registry generate --directory "$FIXTURES_DIR/incomplete-csv-dir"
    assert_failure
    assert_output --partial "error reading the directory"
    assert_output --partial "either the model csv or component csv is missing"
}

@test "given valid spreadsheet credentials and a model name when running mesheryctl registry generate then only that model is generated" {
    require_spreadsheet_credentials
    run $MESHERYCTL_BIN registry generate --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --model "flyte"
    assert_success
    common_outputs_on_success
}

@test "given valid spreadsheet credentials when running mesheryctl registry generate then models are generated successfully" {
    require_spreadsheet_credentials
    run $MESHERYCTL_BIN registry generate --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED"
    assert_success
    common_outputs_on_success
}
