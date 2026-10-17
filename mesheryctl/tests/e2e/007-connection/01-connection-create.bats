setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"

    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/connection"
    mkdir -p "$TESTDATA_DIR"
}

@test "[TC-1013][cut=Kubernetes Connection][tg=Connection Lifecycle] given missing --type flag when running mesheryctl connection create then it fails displaying error message" {
    run $MESHERYCTL_BIN connection create

    assert_failure
    assert_output --partial "Invalid Argument"
    assert_output --partial "Use --type flag"
}

@test "[TC-1013][cut=Kubernetes Connection][tg=Connection Lifecycle] given non valid argument for --type flag when running mesheryctl connection create --type then it fails displaying error message" {
    run $MESHERYCTL_BIN connection create --type foo

    assert_failure
    assert_output --partial "Invalid connection type"
    assert_output --partial "provide a valid connection"
    assert_output --partial "Error"
}

@test "[TC-1013][cut=Kubernetes Connection][tg=Connection Lifecycle] given no argument for --type flag when running mesheryctl connection create --type then it fails displaying error message" {
    run $MESHERYCTL_BIN connection create --type

    assert_failure
    assert_output --partial "flag needs an argument"
    assert_output --partial "Error"
}

@test "[TC-1013][cut=Kubernetes Connection][tg=Connection Lifecycle] given valid type minikube is provided when running mesheryctl connection create --type minikube then a new connection is created" {
    if ! command -v minikube >/dev/null 2>&1; then
        skip "minikube not installed"
    fi

    run $MESHERYCTL_BIN connection create --type minikube
    assert_success
    assert_output --partial "Minikube connection created"
    assert_output --partial "Token set in context minikube"

    # Capture the connection id emitted by `connection create` ("connection_id: <uuid>")
    # so the view/delete suites can exercise the positive lifecycle against a real
    # connection. See create.go setToken().
    CONNECTION_ID=$(
        echo "$output" \
        | grep -oiE 'connection_id: [0-9a-f-]{36}' \
        | head -n1 \
        | awk '{print $2}'
    )
    [ -n "$CONNECTION_ID" ] || fail "connection create did not emit a connection_id line"

    echo "$CONNECTION_ID" > "$TESTDATA_DIR/id"
}
