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

@test "given no arguments when running mesheryctl registry publish then an error about required args is displayed" {
    run $MESHERYCTL_BIN registry publish
    assert_failure
    assert_output --partial "system, google sheet credential, sheet-id, models output path, imgs output path"
    assert_output --partial "are required"
}

@test "given fewer than five arguments when running mesheryctl registry publish then an error about required args is displayed" {
    run $MESHERYCTL_BIN registry publish website
    assert_failure
    assert_output --partial "system, google sheet credential, sheet-id, models output path, imgs output path"
    assert_output --partial "are required"
}

@test "given invalid spreadsheet credentials when running mesheryctl registry publish then an error is displayed" {
    run $MESHERYCTL_BIN registry publish website "invalid-cred" "invalid-id" "$TESTDATA_DIR/models" "$TESTDATA_DIR/imgs"
    assert_failure
    assert_output --partial "Invalid JWT Token"
}

@test "given an invalid output format when running mesheryctl registry publish website then an error is displayed" {
    require_spreadsheet_credentials
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TESTDATA_DIR/models" "$TESTDATA_DIR/imgs" -o "invalid"
    assert_failure
    assert_output --partial "invalid output format"
}

@test "given an invalid system when running mesheryctl registry publish then an error is displayed" {
    require_spreadsheet_credentials
    run $MESHERYCTL_BIN registry publish invalid-system "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TESTDATA_DIR/models" "$TESTDATA_DIR/imgs"
    assert_failure
    assert_output --partial "invalid system"
}

@test "given valid credentials when running mesheryctl registry publish to website with md format then models are published" {
    require_spreadsheet_credentials
    mkdir -p "$TESTDATA_DIR/models-out" "$TESTDATA_DIR/imgs-out"
    run $MESHERYCTL_BIN registry publish website "$TEST_SPREADSHEET_CRED" "$TEST_SPREADSHEET_ID" "$TESTDATA_DIR/models-out" "$TESTDATA_DIR/imgs-out" -o md
    assert_success
}
