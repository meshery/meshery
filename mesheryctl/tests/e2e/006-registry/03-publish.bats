#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/registry-publish"
    mkdir -p "$TESTDATA_DIR"

    # allow CI/environment to provide integration spreadsheet credentials
    export TEST_SPREADSHEET_ID="${TEST_SPREADSHEET_ID:-}"
    export TEST_SPREADSHEET_CRED="${TEST_SPREADSHEET_CRED:-}"
}

require_spreadsheet_credentials() {
    if [[ -z "$TEST_SPREADSHEET_ID" || -z "$TEST_SPREADSHEET_CRED" ]]; then
        skip "Spreadsheet credentials not configured"
    fi
}

@test "given no flags when running mesheryctl registry publish then an error about missing required flags is displayed" {
    run $MESHERYCTL_BIN registry publish
    assert_failure
    assert_output --partial "missing required flag(s)"
    assert_output --partial "--system"
}

@test "given legacy positional arguments when running mesheryctl registry publish then an error about positional arguments is displayed" {
    run $MESHERYCTL_BIN registry publish website "invalid-cred" "invalid-id" "$TESTDATA_DIR/models" "$TESTDATA_DIR/imgs"
    assert_failure
    assert_output --partial "unexpected positional arguments"
}

@test "given website system without imgs-output-path when running mesheryctl registry publish then an error is displayed" {
    run $MESHERYCTL_BIN registry publish --system website --spreadsheet-cred "invalid-cred" --spreadsheet-id "invalid-id" --models-output-path "$TESTDATA_DIR/models"
    assert_failure
    assert_output --partial "--imgs-output-path is required"
}

@test "given an invalid system when running mesheryctl registry publish then an error is displayed" {
    run $MESHERYCTL_BIN registry publish --system invalid-system --spreadsheet-cred "invalid-cred" --spreadsheet-id "invalid-id" --models-output-path "$TESTDATA_DIR/models" --imgs-output-path "$TESTDATA_DIR/imgs"
    assert_failure
    assert_output --partial "invalid system"
}

@test "given invalid spreadsheet credentials when running mesheryctl registry publish then an error is displayed" {
    run $MESHERYCTL_BIN registry publish --system website --spreadsheet-cred "invalid-cred" --spreadsheet-id "invalid-id" --models-output-path "$TESTDATA_DIR/models" --imgs-output-path "$TESTDATA_DIR/imgs"
    assert_failure
    assert_output --partial "Invalid JWT Token"
}

@test "given an invalid output format when running mesheryctl registry publish website then an error is displayed" {
    require_spreadsheet_credentials
    run $MESHERYCTL_BIN registry publish --system website --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --spreadsheet-id "$TEST_SPREADSHEET_ID" --models-output-path "$TESTDATA_DIR/models" --imgs-output-path "$TESTDATA_DIR/imgs" -o "invalid"
    assert_failure
    assert_output --partial "invalid output format"
}

@test "given valid credentials when running mesheryctl registry publish to website with md format then models are published" {
    require_spreadsheet_credentials
    mkdir -p "$TESTDATA_DIR/models-out" "$TESTDATA_DIR/imgs-out"
    run $MESHERYCTL_BIN registry publish --system website --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --spreadsheet-id "$TEST_SPREADSHEET_ID" --models-output-path "$TESTDATA_DIR/models-out" --imgs-output-path "$TESTDATA_DIR/imgs-out" -o md
    assert_success
}
