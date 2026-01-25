setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"

    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/connection"
}

@test "mesheryctl connection delete fails for non-existent connection ID" {
    NONEXISTENT_ID="00000000-0000-0000-0000-000000000000"

    run $MESHERYCTL_BIN connection delete "$NONEXISTENT_ID"
    assert_failure
    assert_output --partial "Invalid API call" || assert_output --partial "Failed"
}

@test "mesheryctl connection delete connection-id removes an existing connection if available" {
    if [ ! -f "$TESTDATA_DIR/id" ]; then
    skip "No connection ID available to delete"
    fi

    CONNECTION_ID="$(cat "$TESTDATA_DIR/id")"
    [ -n "$CONNECTION_ID" ] || skip "Empty connection ID"

    run $MESHERYCTL_BIN connection delete "$CONNECTION_ID"
    assert_success
    assert_output --partial "deleted"
}