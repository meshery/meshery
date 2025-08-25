#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/registry"
    mkdir -p "$TESTDATA_DIR"  
    export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
    
    # Add creds and ID for accessing integration sheets
    export TEST_SPREADSHEET_ID=""
    export TEST_SPREADSHEET_CRED=""
}

@test "mesheryctl registry update displays usage instructions when no arguments are provided" {
    run $MESHERYCTL_BIN registry update 
    assert_failure
    assert_output --partial "unexpected end of JSON input"
}

@test "mesheryctl registry update fails when spreadsheet-id is provided without spreadsheet-cred" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "test-id"
    assert_failure
    assert_output --partial "missing [spreadsheet-cred]"
}

@test "mesheryctl registry update fails when spreadsheet-cred is provided without spreadsheet-id" {
    run $MESHERYCTL_BIN registry update --spreadsheet-cred "test-cred"
    assert_failure
    assert_output --partial "missing [spreadsheet-id]"
}

@test "mesheryctl registry update fails with invalid spreadsheet credentials" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "invalid-id" --spreadsheet-cred "invalid-cred"
    assert_failure
}

@test "mesheryctl registry update succeeds with non-existent input directory but updates nothing" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --input "/non/existent/path"
    assert_success
    assert_output --partial "Updated 0 models and 0 components"
    assert_output --partial "refer"
    assert_output --partial "logs/registry"
}

@test "mesheryctl registry update displays help message" {
    run $MESHERYCTL_BIN registry update --help
    assert_success
    assert_output --partial "Updates the component metadata"
    assert_output --partial "SVGs, shapes, styles and other"
    assert_output --partial "Google Spreadsheet"
    assert_output --partial "Usage:"
    assert_output --partial "mesheryctl registry update [flags]"
    assert_output --partial "Examples:"
    assert_output --partial "Update models from Meshery Integration Spreadsheet"
    assert_output --partial "--spreadsheet-id"
    assert_output --partial "--spreadsheet-cred"
    assert_output --partial "--input"
    assert_output --partial "--model"
}

@test "mesheryctl registry update succeeds with valid spreadsheet credentials and test fixtures, displays logs and summary of the updated components" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --input "$FIXTURES_DIR/test-models"
    assert_success
}

@test "mesheryctl registry update supports model-specific updates" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --input "$FIXTURES_DIR/test-models" --model "kubernetes"
    assert_success
}

@test "mesheryctl registry update handles empty models directory gracefully" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --input "$FIXTURES_DIR/empty-models"
    assert_success
}
