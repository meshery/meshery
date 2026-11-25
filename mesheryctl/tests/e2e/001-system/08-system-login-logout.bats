#!/usr/bin/env bats

# End-to-end tests for `mesheryctl system login` and `mesheryctl system logout`.
#
# Assumes Meshery is already deployed and running in Kubernetes.
#
# A dedicated temporary Meshery configuration is used so these tests
# do not modify the user's existing Meshery configuration.
#
# Port forwarding is started once for the whole test file.
#
# Background jobs (port-forward, login) run in their own process group
# (via setsid, available on Ubuntu) so cleanup can signal the whole group.

# Execute mesheryctl using the isolated test configuration.
mesheryctl_cmd() {
    "$MESHERYCTL_BIN" --config "$TEST_MESHERY_CONFIG" "$@"
}

# Start a command in its own process group (PGID == PID).
# Sets PROCESS_GROUP_PID. Call directly, not inside $(...).
start_process_group() {
    set -m
    "$@" </dev/null &
    PROCESS_GROUP_PID=$!
    set +m
}

# Terminate every process in a process group created by start_process_group.
# Uses SIGTERM first, followed by SIGKILL if necessary.
stop_process_group() {
    local pgid="$1"
    local own_pgid
    local attempt

    [[ "$pgid" =~ ^[0-9]+$ ]] || return 0

    # Never signal the Bats process group.
    own_pgid="$(ps -o pgid= -p "$$" | tr -d ' ')"
    [[ "$pgid" != "$own_pgid" ]] || return 0

    # Send SIGTERM to the entire process group.
    kill -TERM -- "-$pgid" 2>/dev/null || return 0

    # Wait for graceful termination.
    for ((attempt = 0; attempt < 20; attempt++)); do
        kill -0 -- "-$pgid" 2>/dev/null || return 0
        sleep 0.1
    done

    # Force termination if processes remain.
    kill -KILL -- "-$pgid" 2>/dev/null || true
}

setup_file() {
    local source_config
    local service_port
    local i

    export MESHERY_NAMESPACE="meshery"
    export MESHERY_SERVICE="meshery"
    export DASHBOARD_PORT="9081"

    source_config="${MESHERY_CONFIG:-$HOME/.meshery/config.yaml}"

    # Ensure an existing Meshery configuration is available.
    if [[ ! -f "$source_config" ]]; then
        echo "Meshery configuration not found: $source_config" >&3
        return 1
    fi

    # Create an isolated configuration for this test run.
    TEST_CONFIG_DIR="$(mktemp -d)" || return 1
    export TEST_CONFIG_DIR

    export TEST_MESHERY_CONFIG="$TEST_CONFIG_DIR/config.yaml"

    cp "$source_config" "$TEST_MESHERY_CONFIG" || return 1

    # Retrieve the Kubernetes Service port, not its NodePort.
    service_port="$(kubectl get svc "$MESHERY_SERVICE" \
        -n "$MESHERY_NAMESPACE" \
        -o jsonpath='{.spec.ports[0].port}')" || return 1

    if [[ ! "$service_port" =~ ^[0-9]+$ ]]; then
        echo "Could not determine Meshery Service port." >&3
        return 1
    fi

    # Create a dedicated context pointing to our predictable local port.
    export TEST_CONTEXT="bats-e2e-$$"

    export TEST_ENDPOINT="http://127.0.0.1:${DASHBOARD_PORT}"

    if ! "$MESHERYCTL_BIN" --config "$TEST_MESHERY_CONFIG" \
        system context create "$TEST_CONTEXT" \
        --platform kubernetes \
        --url "$TEST_ENDPOINT" \
        --set \
        --yes; then
        echo "Failed to create the temporary Meshery context." >&3
        return 1
    fi

    # Start Kubernetes port forwarding independently of mesheryctl,
    # in its own process group.
    PORT_FORWARD_LOG="$(mktemp)" || return 1
    export PORT_FORWARD_LOG

    start_process_group kubectl port-forward \
        -n "$MESHERY_NAMESPACE" \
        "svc/$MESHERY_SERVICE" \
        "${DASHBOARD_PORT}:${service_port}" \
        >"$PORT_FORWARD_LOG" 2>&1 3>&-

    export PORT_FORWARD_PID=$PROCESS_GROUP_PID

    # Wait until the Meshery API becomes reachable.
    for ((i = 0; i < 60; i++)); do
        if ! kill -0 "$PORT_FORWARD_PID" 2>/dev/null; then
            echo "Kubernetes port forwarding exited unexpectedly:" >&3
            cat "$PORT_FORWARD_LOG" >&3
            return 1
        fi

        if curl --silent --fail \
            --connect-timeout 1 \
            --max-time 2 \
            "$TEST_ENDPOINT/api/system/version" \
            >/dev/null 2>&1; then
            return 0
        fi

        sleep 1
    done

    echo "Meshery API not reachable at $TEST_ENDPOINT:" >&3
    cat "$PORT_FORWARD_LOG" >&3

    stop_process_group "$PORT_FORWARD_PID"
    return 1
}

teardown_file() {
    # Stop the port-forward process group.
    if [[ -n "${PORT_FORWARD_PID:-}" ]]; then
        stop_process_group "$PORT_FORWARD_PID"
        wait "$PORT_FORWARD_PID" 2>/dev/null || true
    fi

    # Remove temporary files and configuration.
    [[ -z "${PORT_FORWARD_LOG:-}" ]] || rm -f "$PORT_FORWARD_LOG"
    [[ -z "${TEST_CONFIG_DIR:-}" ]] || rm -rf "$TEST_CONFIG_DIR"
}

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    # Start each test logged out.
    # Logging out when already logged out is expected.
    run mesheryctl_cmd system logout --yes
}

teardown() {
    run mesheryctl_cmd system logout --yes
}

# Initiates Meshery provider authentication without waiting for browser
# interaction. Succeeds as soon as "Initiating login..." is printed.
verify_meshery_login_initiation() {
    local found=0
    local i
    local login_log
    local login_pgid

    login_log="$(mktemp)" || return 1

    # Start login in its own process group so that any helper it spawns
    # (for example a browser opener) is cleaned up with it.
    # mesheryctl_cmd is intentionally not used here: setsid must execute
    # the binary directly so the recorded PID is the CLI's own PID.
    start_process_group "$MESHERYCTL_BIN" --config "$TEST_MESHERY_CONFIG" \
        system login -p Meshery \
        >"$login_log" 2>&1 3>&-

    login_pgid=$PROCESS_GROUP_PID

    # Wait up to 30 seconds for the authentication initiation message.
    for ((i = 0; i < 30; i++)); do
        if grep -Fq "Initiating login..." "$login_log"; then
            found=1
            break
        fi

        # Stop waiting if the login process exits prematurely.
        if ! kill -0 "$login_pgid" 2>/dev/null; then
            grep -Fq "Initiating login..." "$login_log" && found=1
            break
        fi

        sleep 1
    done

    # Terminate the entire login process group.
    stop_process_group "$login_pgid"

    # Reap the background CLI process.
    wait "$login_pgid" 2>/dev/null || true

    if [[ "$found" -ne 1 ]]; then
        echo "Meshery authentication was not initiated."
        cat "$login_log"
        rm -f "$login_log"
        return 1
    fi

    cat "$login_log"
    rm -f "$login_log"
}

@test "[cut=System][tg=System Login] Given Meshery is deployed, When Local is selected explicitly, Then authentication should succeed" {
    run mesheryctl_cmd system login -p Local

    assert_success
    assert_output --partial "authenticated"
}

@test "[cut=System][tg=System Login] Given Meshery is deployed, When the Meshery provider is selected, Then authentication should be initiated" {
    run verify_meshery_login_initiation

    assert_success
    assert_output --partial "Initiating login..."
}

@test "[cut=System][tg=System Login] Given Meshery is deployed, When an unknown provider is selected, Then authentication should fail" {
    run mesheryctl_cmd system login -p Unknown

    assert_failure
    assert_output --partial 'Provider "Unknown" is not offered by this Meshery Server.'
}

@test "[cut=System][tg=System Logout] Given the user is logged in, When logout is run, Then the session should end" {
    # Log in first so there is a session to end.
    run mesheryctl_cmd system login -p Local
    assert_success

    run mesheryctl_cmd system logout

    assert_success
    assert_output --regexp 'logged out'

}
