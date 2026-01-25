setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"

    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/connection"
}

@test "meshery connection view invalid connection-id shows error" {
    run $MESHERYCTL_BIN connection view
    assert_failure
    assert_output --partial "Error" || assert_output --partial "Invalid Argument"
}

@test "meshery connection view connection-id displays connection details" {
    if [ ! -f "$TESTDATA_DIR/id" ]; then
        skip "No connection ID available to view"
    fi

    CONNECTION_ID="$(cat "$TESTDATA_DIR/id")"
    [ -n "$CONNECTION_ID" ] || skip "Empty connection ID"

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID"
    assert_success
    assert_output --partial "id" || assert_output --partial "name" || assert_output --partial "metadata"
}