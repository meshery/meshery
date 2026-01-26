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

@test "given no connection-id is provided as an argument when running mesheryctl connection view then a message error is displayed" {
    run $MESHERYCTL_BIN connection view
    assert_failure
    assert_output --partial "Error" 
    assert_output --partial "Invalid Argument"
}

@test "given a valid connection-id is provided as an argument when running meshery connection view connection-id then the connection details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID"
    assert_success
    assert_output --partial "id" 
    assert_output --partial "name" 
    assert_output --partial "metadata" 
}

@test "given no connection-id is provided as an argument when running mesheryctl connection view --save then a message error is displayed" {
    run $MESHERYCTL_BIN connection view --save
    assert_failure
    assert_output --partial "Error" 
    assert_output --partial "Invalid Argument"
    assert_output --partial "ID isn't specified"
}

@test "given a valid connection-id is provided as an argument when running mesheryctl connection view --save then a details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view --save "$CONNECTION_ID"
    assert_success
    assert_file_exists "$HOME/.meshery/connection_minikube.yaml"
}

@test "given an invalid connection-id is provided as an argument when running mesheryctl connection view --save then a message error is displayed" {
    NONEXISTENT_ID="00000000-0000-0000-0000-000000000000"
    
    run $MESHERYCTL_BIN connection view --save "$NONEXISTENT_ID"
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid connection ID"
}

@test "given no argument is provided when running mesheryctl connection view connection-id --output-format then a message error is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format
    assert_failure
    assert_output --partial "Error" 
    assert_output --partial "flag needs an argument"
}

@test "given invalid argument is provided as an argument when running mesheryctl connection view connection-id --output-format then a message error is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format foo
    assert_failure
    assert_output --partial "Error" 
    assert_output --partial "output-format choice is invalid" 
    assert_output --partial "use [json|yaml]" 
}

@test "given a valid argument is provided as an argument when running mesheryctl connection view connection-id --output-format yaml then a details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format yaml
    assert_success
    assert_output --partial "id: $CONNECTION_ID" 
    assert_output --partial "name"
    assert_output --partial "metadata"
    assert_output --partial "user_id"
}

@test "given a valid argument is provided as an argument when running mesheryctl connection view connection-id --output-format json then a details in default format is displayed" {
    require_connection_id

    run $MESHERYCTL_BIN connection view "$CONNECTION_ID" --output-format json
    assert_success
    assert_output --partial "\"id\": \"$CONNECTION_ID\""
    assert_output --partial "\"name\""
    assert_output --partial "\"metadata\""
    assert_output --partial "\"user_id\""
}

@test "given an invalid connection-id is provided as an argument when running mesheryctl connection view --output-format json/yaml then a message error is displayed" {
    NONEXISTENT_ID="00000000-0000-0000-0000-000000000000"
    
    run $MESHERYCTL_BIN connection view "$NONEXISTENT_ID" --output-format json
    assert_failure
    assert_output --partial "Error"
    assert_output --partial "Invalid connection ID"
}

@test "given no connection-id is provided as an argument when running mesheryctl connection view --output-format then a message error is displayed" {
    run $MESHERYCTL_BIN connection view --output-format yaml
    assert_failure
    assert_output --partial "Error" 
    assert_output --partial "ID isn't specified" 
}