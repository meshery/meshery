setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"

    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/connection"
}

teardown_file() {
    rm -rf "$TESTDATA_DIR"
}

@test "[TC-1067][cut=Kubernetes Connection] given no connection-id provided as an argument when running mesheryctl connection delete then an error message is displayed" {
    run $MESHERYCTL_BIN connection delete
    assert_failure
    assert_output --partial "connection name or ID isn't specified"
    assert_output --partial "Error"
}

@test "[TC-1066][cut=Kubernetes Connection] given a non existing connection-id is provided as an argument when running mesheryctl connection delete non-existing-id then an error message is displayed" {
    NONEXISTENT_ID="00000000-0000-0000-0000-000000000000"

    run $MESHERYCTL_BIN connection delete "$NONEXISTENT_ID"
    assert_success
    assert_output --partial "No connection with ID"
}

@test "[TC-1067][cut=Kubernetes Connection] given an invalid connection-id is provided as an argument when running mesheryctl connection delete invalid-connection-id then an error message is displayed" {
    INVALID_ID="0000"

    run $MESHERYCTL_BIN connection delete "$INVALID_ID"
    assert_failure
    assert_output --partial "Invalid ID format"
    assert_output --partial "Error"
}

# Scenario X9 (documented gap): mesheryctl exposes only create/list/view/delete.
# The connect/disconnect/ignore lifecycle transitions are UI-only (PUT status),
# so there is no CLI subcommand for them. Assert their absence so the gap is
# tracked and a future parity feature is caught the moment it lands.
@test "[TC-1069][cut=Kubernetes Connection] given the connection command when inspecting subcommands then no CLI transition (connect/disconnect/ignore) command exists" {
    run $MESHERYCTL_BIN connection --help
    assert_success
    refute_output --partial "disconnect"
    refute_output --partial "ignore"

    run $MESHERYCTL_BIN connection disconnect
    assert_failure
    run $MESHERYCTL_BIN connection ignore
    assert_failure
}

# Scenario M9 (documented gap): there is no CLI command to view or reset MeshSync
# data for a connection (flushMeshsync is UI/REST-only). Assert absence.
@test "[TC-1049][cut=Kubernetes Connection] given the connection command when inspecting subcommands then no CLI meshsync data command exists" {
    run $MESHERYCTL_BIN connection meshsync
    assert_failure
    run $MESHERYCTL_BIN connection flush
    assert_failure
}

@test "[TC-1065][cut=Kubernetes Connection] given a valid connection-id is provided as an argument when running mesheryctl connection delete connection-id then the existing connection is deleted" {
    if [ ! -f "$TESTDATA_DIR/id" ]; then
        skip "No connection ID available to delete"
    fi

    CONNECTION_ID="$(cat "$TESTDATA_DIR/id")"
    [ -n "$CONNECTION_ID" ] || skip "Empty connection ID"

    run $MESHERYCTL_BIN connection delete "$CONNECTION_ID"
    assert_success
    assert_output --partial "deleted"
}
