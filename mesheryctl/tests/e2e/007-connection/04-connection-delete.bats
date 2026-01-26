setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"

    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/connection"
}

teardown_file() {
    rm -rf "$TESTDATA_DIR"
}

@test "given non valid connection-id is provided as an argument when running mesheryctl connection delete connection-id then a message error is displayed" {
    NONEXISTENT_ID="00000000-0000-0000-0000-000000000000"

    run $MESHERYCTL_BIN connection delete "$NONEXISTENT_ID"
    assert_failure
    assert_output --partial "Invalid API call" 
    assert_output --partial "Failed"
}

@test "given a valid connection-id is provided as an argument when running mesheryctl connection delete connection-id then the existing connection is deleted" {
    if [ ! -f "$TESTDATA_DIR/id" ]; then
        skip "No connection ID available to delete"
    fi

    CONNECTION_ID="$(cat "$TESTDATA_DIR/id")"
    [ -n "$CONNECTION_ID" ] || skip "Empty connection ID"

    run $MESHERYCTL_BIN connection delete "$CONNECTION_ID"
    assert_success
    assert_output --partial "deleted"
}