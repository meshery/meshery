setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"

    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/connection"
    mkdir -p "$TESTDATA_DIR"
}

@test "given missing --type flag when running mesheryctl connection create then it fails displaying error message" {
    run $MESHERYCTL_BIN connection create

    assert_failure
    assert_output --partial "Invalid Argument"
    assert_output --partial "Use --type flag"
}

@test "given non valid argument for --type flag when running mesheryctl connection create --type then it fails displaying error message" {
    run $MESHERYCTL_BIN connection create --type foo

    assert_failure
    assert_output --partial "Invalid connection type"
    assert_output --partial "provide a valid connection"
    assert_output --partial "Error"
}

@test "given no argument for --type flag when running mesheryctl connection create --type then it fails displaying error message" {
    run $MESHERYCTL_BIN connection create --type

    assert_failure
    assert_output --partial "flag needs an argument"
    assert_output --partial "Error"
}

@test "given valid type minikube is provided when running mesheryctl connection create --type minikube then a new connection is created" {
    if ! command -v minikube >/dev/null 2>&1; then
        skip "minikube not installed"
    fi

    run $MESHERYCTL_BIN connection create --type minikube
    assert_success
    assert_output --partial "Minikube connection created" 
    assert_output --partial "Token set in context minikube"

    #Extract connection ID if present and store it in temp dir
    CONNECTION_ID=$(
        echo "$output" |
        grep -oE '"connection_id"\s*:\s*"[^"]+"' |
        cut -d '"' -f 4
    )
    [ -n "$CONNECTION_ID" ] || skip "No ID found in output"

    echo "$CONNECTION_ID" > "$TESTDATA_DIR/id"

}