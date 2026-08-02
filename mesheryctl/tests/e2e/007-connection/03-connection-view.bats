setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"

    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/connection"

    if [ -f "$TESTDATA_DIR/id" ]; then
    CONNECTION_ID="$(cat "$TESTDATA_DIR/id")"
    export CONNECTION_ID
    fi
}

require_connection_id() {
    [ -n "$CONNECTION_ID" ] || skip "No connection ID available"
}

@test "[TC-1080][cut=Kubernetes Connection] given no connection-id is provided as an argument when running mesheryctl connection view then a message error is displayed" {
    run $MESHERYCTL_BIN connection view
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid Argument"
}

@test "[TC-1080][cut=Kubernetes Connection] given a valid connection-id is provided as an argument when running meshery connection view connection-id then the connection details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID"
    assert_success
    assert_output --partial "id"
    assert_output --partial "name"
    assert_output --partial "metadata"
}

@test "[TC-1080][cut=Kubernetes Connection] given no connection-id is provided as an argument when running mesheryctl connection view --save then a message error is displayed" {
    run $MESHERYCTL_BIN connection view --save
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid Argument"
    assert_output --partial "ID isn't specified"
}

@test "[TC-1080][cut=Kubernetes Connection] given a valid connection-id is provided as an argument when running mesheryctl connection view --save then a details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view --save "$CONNECTION_ID"
    assert_success

    SAVED_FILE="$(echo "$output" | sed -n 's/.*Data saved to file: //p')"
    assert_file_exists "$SAVED_FILE"
}

@test "[TC-1080][cut=Kubernetes Connection] given an invalid connection-id is provided as an argument when running mesheryctl connection view --save then a message error is displayed" {
    NONEXISTENT_ID="00000000-0000-0000-0000-000000000000"

    run $MESHERYCTL_BIN connection view --save "$NONEXISTENT_ID"
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid connection ID"
}

@test "[TC-1080][cut=Kubernetes Connection] given no argument is provided when running mesheryctl connection view connection-id --output-format then a message error is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "flag needs an argument"
}

@test "[TC-1080][cut=Kubernetes Connection] given invalid argument is provided as an argument when running mesheryctl connection view connection-id --output-format then a message error is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format foo
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "output-format choice is invalid"
    assert_output --partial "use [json|yaml]"
}

@test "[TC-1080][cut=Kubernetes Connection] given a valid argument is provided as an argument when running mesheryctl connection view connection-id --output-format yaml then a details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format yaml
    assert_success
    assert_output --partial "id: $CONNECTION_ID"
    assert_output --partial "name"
    assert_output --partial "metadata"
    assert_output --partial "user_id"
}

@test "[TC-1080][cut=Kubernetes Connection] given a valid argument is provided as an argument when running mesheryctl connection view connection-id --output-format json then a details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format json
    assert_success
    assert_output --partial "\"id\": \"$CONNECTION_ID\""
    assert_output --partial "\"name\""
    assert_output --partial "\"metadata\""
    assert_output --partial "\"user_id\""
}

@test "[TC-1080][cut=Kubernetes Connection] given an invalid connection-id is provided as an argument when running mesheryctl connection view --output-format json/yaml then a message error is displayed" {
    NONEXISTENT_ID="00000000-0000-0000-0000-000000000000"

    run $MESHERYCTL_BIN connection view "$NONEXISTENT_ID" --output-format json
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid connection ID"
}

@test "[TC-1080][cut=Kubernetes Connection] given no connection-id is provided as an argument when running mesheryctl connection view --output-format then a message error is displayed" {
    run $MESHERYCTL_BIN connection view --output-format yaml
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "ID isn't specified"
}

# Scenario L5 (view by name): the kubernetes connection's name is its context
# name (e.g. "minikube"). Resolving a connection by name exercises the
# fetchConnectionByName search path (view.go), distinct from view-by-id.
#
# fetchConnectionByName issues a substring `search` and, when more than one
# connection matches, falls to an interactive promptui.Select that cannot be
# answered under bats `run` (no TTY). To stay deterministic and non-blocking we
# only exercise the by-name path when the name resolves to a single kubernetes
# connection and skip otherwise, so the prompt can never be reached.
@test "[TC-1080][cut=Kubernetes Connection] given a valid connection name is provided when running mesheryctl connection view name then the connection details are displayed" {
    require_connection_id

    # Resolve the connection's name from its id.
    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format json
    assert_success
    CONNECTION_NAME="$(echo "$output" | sed -n 's/.*"name": "\([^"]*\)".*/\1/p' | head -n1)"
    [ -n "$CONNECTION_NAME" ] || skip "connection has no name to resolve by"

    # Guard: count kubernetes connections whose row carries this name. The
    # server's by-name search is a substring match, so grep -F (fixed-string,
    # substring) is the conservative mirror — any collision over-counts and
    # skips rather than risking the interactive prompt.
    run $MESHERYCTL_BIN connection list --kind kubernetes --pagesize 10000
    assert_success
    MATCH_COUNT="$(echo "$output" | grep -cF "$CONNECTION_NAME" || true)"
    [ "$MATCH_COUNT" -eq 1 ] || skip "connection name '$CONNECTION_NAME' is not unique (matched $MATCH_COUNT); by-name resolution would prompt"

    run $MESHERYCTL_BIN connection view "$CONNECTION_NAME" --output-format json
    assert_success
    assert_output --partial "\"id\": \"$CONNECTION_ID\""
    assert_output --partial "\"name\": \"$CONNECTION_NAME\""
}