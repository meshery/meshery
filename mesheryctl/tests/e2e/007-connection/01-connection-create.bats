setup(){
    load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"

    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/connection"
    mkdir -p "$TESTDATA_DIR"
}

@test "mesheryctl connection create fails without --type flag" {
    run $MESHERYCTL_BIN connection create

    assert_failure
    assert_output --partial "Invalid Argument"
}

@test "mesheryctl connection create --type minikube succeeds" {
    if ! command -v minikube >/dev/null 2>&1; then
        skip "minikube not installed"
    fi

    run $MESHERYCTL_BIN connection create --type minikube
    assert_success
    assert_output --partial "Minikube connection created" || assert_output --partial "Token set in context minikube"

    #Extract connection ID if present and store it in temp dir
    CONNECTION_ID=$(
        echo "$output" |
        grep -oE '"connection_id"\s*:\s*"[^"]+"' |
        head -n1 |
        sed -E 's/.*"connection_id"\s*:\s*"([^"]+)".*/\1/'
    )
    [ -n "$CONNECTION_ID" ] || skip "No ID found in output"

    echo "$CONNECTION_ID" > "$TESTDATA_DIR/id"

}