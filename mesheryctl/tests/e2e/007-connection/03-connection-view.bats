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

@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given no connection-id is provided as an argument when running mesheryctl connection view then a message error is displayed" {
    run $MESHERYCTL_BIN connection view
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid Argument"
}

@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given a valid connection-id is provided as an argument when running meshery connection view connection-id then the connection details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID"
    assert_success
    assert_output --partial "id"
    assert_output --partial "name"
    assert_output --partial "metadata"
}

@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given no connection-id is provided as an argument when running mesheryctl connection view --save then a message error is displayed" {
    run $MESHERYCTL_BIN connection view --save
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid Argument"
    assert_output --partial "ID isn't specified"
}

@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given a valid connection-id is provided as an argument when running mesheryctl connection view --save then a details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view --save "$CONNECTION_ID"
    assert_success

    SAVED_FILE="$(echo "$output" | sed -n 's/.*Data saved to file: //p')"
    assert_file_exists "$SAVED_FILE"
}

# The nil UUID is syntactically well-formed, so view takes the by-id path
# (isArgumentUUID) and the server rejects it in GetConnectionByID
# (server/handlers/connections_handlers.go), which writes a 400
# models.ErrInvalidUUID. That envelope carries "Invalid UUID" as its short
# description and "invalid connection ID" as its cause. Assert those, not
# "Invalid connection ID" (capitalised) — that string is emitted nowhere in
# mesheryctl or meshery server, so the assertion could never pass.
@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given an invalid connection-id is provided as an argument when running mesheryctl connection view --save then a message error is displayed" {
    INVALID_ID="00000000-0000-0000-0000-000000000000"

    run $MESHERYCTL_BIN connection view --save "$INVALID_ID"
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid UUID"
    assert_output --partial "invalid connection ID"
}

@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given no argument is provided when running mesheryctl connection view connection-id --output-format then a message error is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "flag needs an argument"
}

@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given invalid argument is provided as an argument when running mesheryctl connection view connection-id --output-format then a message error is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format foo
    assert_failure
    assert_output --partial "Error"
    # display.ErrInvalidOutputFormat: short description "Invalid Output Format",
    # cause `Provided output format "foo" is invalid`, remediation
    # "Ensure using [json|yaml] as the output format". The previously asserted
    # "output-format choice is invalid" / "use [json|yaml]" are not emitted by
    # any of those four parts.
    assert_output --partial "Invalid Output Format"
    assert_output --partial "is invalid"
    assert_output --partial "[json|yaml]"
}

@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given a valid argument is provided as an argument when running mesheryctl connection view connection-id --output-format yaml then a details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format yaml
    assert_success
    assert_output --partial "id: $CONNECTION_ID"
    assert_output --partial "name"
    assert_output --partial "metadata"
    # Wire fields are camelCase (schemas Connection: json/yaml tags); createdAt is
    # always present, whereas userId is omitempty. Asserting "user_id" would fail
    # against the real output.
    assert_output --partial "createdAt"
}

@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given a valid argument is provided as an argument when running mesheryctl connection view connection-id --output-format json then a details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format json
    assert_success
    assert_output --partial "\"id\": \"$CONNECTION_ID\""
    assert_output --partial "\"name\""
    assert_output --partial "\"metadata\""
    # camelCase wire field (see the yaml case above); "user_id" is not emitted.
    assert_output --partial "\"createdAt\""
}

# Same 400 models.ErrInvalidUUID envelope as the --save case above; the
# rejection happens before any output formatting, so --output-format json does
# not change the error text.
@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given an invalid connection-id is provided as an argument when running mesheryctl connection view --output-format json/yaml then a message error is displayed" {
    INVALID_ID="00000000-0000-0000-0000-000000000000"

    run $MESHERYCTL_BIN connection view "$INVALID_ID" --output-format json
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid UUID"
    assert_output --partial "invalid connection ID"
}

@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given no connection-id is provided as an argument when running mesheryctl connection view --output-format then a message error is displayed" {
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
@test "[TC-1080][cut=Kubernetes Connection][tg=Connection Lifecycle] given a valid connection name is provided when running mesheryctl connection view name then the connection details are displayed" {
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