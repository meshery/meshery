#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    # Use system temp directory instead of local temp
    export TESTDATA_DIR="$(mktemp -d)/testdata/registry"

    export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"

    # Add creds and ID for accessing integration sheets
    export TEST_SPREADSHEET_ID=""
    export TEST_SPREADSHEET_CRED=""

    export TEST_MODELS_OUTPUT="$TESTDATA_DIR/models-output"
    export TEST_IMGS_OUTPUT="$TESTDATA_DIR/imgs-output"
    mkdir -p "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT"
}

assert_missing_argument_message() {
    assert_output --partial "[ system, google sheet credential, sheet-id, models output path, imgs output path] are required"
}

@test "mesheryctl registry publish fails when insufficient arguments are provided" {
    run $MESHERYCTL_BIN registry publish website
    assert_failure
    assert_missing_argument_message
}

@test "mesheryctl registry publish fails when only 4 arguments are provided" {
    run $MESHERYCTL_BIN registry publish website cred id path
    assert_failure
    assert_missing_argument_message
}

@test "mesheryctl registry publish fails with invalid credentials" {
    run $MESHERYCTL_BIN registry publish website "invalid-cred" "invalid-id" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT"
    assert_failure
    assert_output --partial "Invalid JWT Token: Ensure the provided token is a base64-encoded, valid Google Spreadsheets API token"
}

@test "mesheryctl registry publish succeeds with meshery system" {
    run $MESHERYCTL_BIN registry publish meshery "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT"
    assert_success
    assert_output --partial "Downloaded CSV from:"
    assert_output --partial "Parsing Models..."
    assert_output --partial "Parsing Components..."
    assert_output --partial "Parsing Relationships..."
    assert_output --partial "Removing CSV file:"
}

@test "mesheryctl registry publish succeeds with remote-provider system" {
    run $MESHERYCTL_BIN registry publish remote-provider "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT"
    assert_success
    assert_output --partial "Downloaded CSV from:"
    assert_output --partial "Parsing Models..."
    assert_output --partial "Parsing Components..."
    assert_output --partial "Parsing Relationships..."
    assert_output --partial "Total model published: 1"
    assert_output --partial "Removing CSV file:"
}

@test "mesheryctl registry publish succeeds with website system and md output format" {
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT" -o md
    assert_success
    assert_success
    assert_output --partial "Downloaded CSV from:"
    assert_output --partial "Parsing Models..."
    assert_output --partial "Parsing Components..."
    assert_output --partial "Parsing Relationships..."
    assert_output --partial "no relationships found for contour-operator"
    assert_output --partial "Removing CSV file:"
}

@test "mesheryctl registry publish succeeds with website system and mdx output format" {
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT" -o mdx
    assert_success
    assert_output --partial "Downloaded CSV from:"
    assert_output --partial "Parsing Models..."
    assert_output --partial "Parsing Components..."
    assert_output --partial "Parsing Relationships..."
    assert_output --partial "no relationships found for contour-operator"
    assert_output --partial "Removing CSV file:"
}

@test "mesheryctl registry publish succeeds with website system and js output format" {
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT" -o js
    assert_success
    assert_output --partial "Downloaded CSV from:"
    assert_output --partial "Parsing Models..."
    assert_output --partial "Parsing Components..."
    assert_output --partial "Parsing Relationships..."
    assert_output --partial "no relationships found for contour-operator"
    assert_output --partial "Removing CSV file:"
}

@test "mesheryctl registry publish fails with website system and invalid output format" {
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT" -o invalid
    assert_failure
    assert_output --partial "Downloaded CSV from:"
    assert_output --partial "Parsing Models..."
    assert_output --partial "Parsing Components..."
    assert_output --partial "Parsing Relationships..."
    assert_output --partial "Error: invalid output format: invalid"
}

@test "mesheryctl registry publish creates output directories if they don't exist" {
    local new_models_dir="$TESTDATA_DIR/new-models"
    local new_imgs_dir="$TESTDATA_DIR/new-imgs"

    run $MESHERYCTL_BIN registry publish remote-provider "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$new_models_dir" "$new_imgs_dir"
    assert_success
    assert_output --partial "Downloaded CSV from:"
    assert_output --partial "Parsing Models..."
    assert_output --partial "Parsing Components..."
    assert_output --partial "Parsing Relationships..."
    assert_output --partial "Total model published: 1"
    assert_output --partial "Removing CSV file:"
}

@test "mesheryctl registry publish handles non-existent output paths gracefully" {
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TESTDATA_DIR/non/existent/models" "$TESTDATA_DIR/non/existent/imgs" -o md
    assert_success
    assert_output --partial "Downloaded CSV from:"
    assert_output --partial "Parsing Models..."
    assert_output --partial "Parsing Components..."
    assert_output --partial "Parsing Relationships..."
    assert_output --partial "no relationships found for contour-operator"
    assert_output --partial "Removing CSV file:"
}