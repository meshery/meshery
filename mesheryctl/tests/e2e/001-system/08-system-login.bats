#!/usr/bin/env bats

# End-to-end tests for `mesheryctl system login`.
#
# Every test starts from a clean state: Meshery is stopped and redeployed
# (--reset) and the saved login token is removed, so tests cannot influence
# each other.
#
# Environment:
#   E2E_HELPERS_PATH  path to the shared bats helpers (required)
#   MESHERYCTL_BIN    mesheryctl binary to test (default: mesheryctl on PATH)
#   AUTH_FILE         saved login token (default: ~/.meshery/auth.json)

DASHBOARD_PORT=9081
setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    # Do not open real browser windows during tests.
    export BROWSER=true

    DASHBOARD_PID=""
    DASHBOARD_LOG=""
    LOGIN_PID=""
    LOGIN_LOG=""

    # Start every test logged out.
    rm -f "${AUTH_FILE:-$HOME/.meshery/auth.json}"

    # Stopping an already-stopped Meshery instance is expected.
    run $MESHERYCTL_BIN system stop --yes >/dev/null 2>&1 || true
}

teardown() {
    # Stop login and dashboard processes, if still running.
    if [[ -n "${LOGIN_PID:-}" ]]; then
        stop_process_tree "$LOGIN_PID"
        wait "$LOGIN_PID" 2>/dev/null || true
    fi

    if [[ -n "${DASHBOARD_PID:-}" ]]; then
        stop_process_tree "$DASHBOARD_PID"
        wait "$DASHBOARD_PID" 2>/dev/null || true
    fi

    # Remove temporary log files.
    [[ -z "${LOGIN_LOG:-}" ]] || rm -f "$LOGIN_LOG"
    [[ -z "${DASHBOARD_LOG:-}" ]] || rm -f "$DASHBOARD_LOG"

    rm -f "${AUTH_FILE:-$HOME/.meshery/auth.json}"

    # Stop Meshery regardless of the test outcome.
    "$MESHERYCTL_BIN" system stop --force -y >/dev/null 2>&1 || true
}

# Terminates a process and its descendants.
# Uses SIGTERM first, followed by SIGKILL if necessary.
stop_process_tree() {
    local pid="$1"
    local child
    local attempt

    # Terminate child processes before their parent.
    for child in $(pgrep -P "$pid" 2>/dev/null || true); do
        stop_process_tree "$child"
    done

    kill -TERM "$pid" 2>/dev/null || return 0

    # Allow the process time to terminate gracefully.
    for ((attempt = 0; attempt < 10; attempt++)); do
        kill -0 "$pid" 2>/dev/null || return 0
        sleep 0.1
    done

    # Force termination if the process is still running.
    kill -KILL "$pid" 2>/dev/null || true
}

# Starts Meshery on Kubernetes and opens the dashboard port-forward.
start_meshery() {
    run "$MESHERYCTL_BIN" system start \
        --skip-browser \
        --skip-update \
        --reset \
        -p kubernetes \
        --yes

    assert_success
    # Matches both "Meshery deployed." and "Meshery deployed on Kubernetes."
    assert_output --partial "Meshery deployed"

    start_dashboard
}

# Starts dashboard port forwarding in the background and waits until the
# Meshery API is reachable. Runs in the test shell (not under `run`) so that
# DASHBOARD_PID is visible to teardown.
start_dashboard() {
    local i

    DASHBOARD_LOG="$(mktemp)" || return 1

    # `3>&-` keeps the background process from holding bats' fd 3 open,
    # which would make bats hang when the run finishes.
    "$MESHERYCTL_BIN" system dashboard \
        --port-forward \
        -p "$DASHBOARD_PORT" \
        >"$DASHBOARD_LOG" 2>&1 3>&- &

    DASHBOARD_PID=$!

    # Wait up to approximately 30 seconds for the Meshery API.
    for ((i = 0; i < 30; i++)); do
        if ! kill -0 "$DASHBOARD_PID" 2>/dev/null; then
            echo "Dashboard port forwarding exited unexpectedly:" >&3
            cat "$DASHBOARD_LOG" >&3
            return 1
        fi

        if curl --silent --fail \
            --connect-timeout 1 \
            --max-time 1 \
            "http://localhost:${DASHBOARD_PORT}/api/system/version" \
            >/dev/null 2>&1; then
            return 0
        fi

        sleep 1
    done

    echo "Timed out waiting for the Meshery API on port ${DASHBOARD_PORT}:" >&3
    cat "$DASHBOARD_LOG" >&3
    return 1
}

# Initiates Meshery provider authentication without waiting for browser
# interaction. Succeeds as soon as "Initiating login..." is printed.
verify_meshery_login_initiation() {
    local found=0
    local i

    LOGIN_LOG="$(mktemp)" || return 1

    # Start login independently so its output can be monitored.
    "$MESHERYCTL_BIN" system login -p Meshery \
        >"$LOGIN_LOG" 2>&1 3>&- &

    LOGIN_PID=$!

    # Wait up to 30 seconds for the authentication initiation message.
    for ((i = 0; i < 30; i++)); do
        if grep -Fq "Initiating login..." "$LOGIN_LOG"; then
            found=1
            break
        fi

        # Stop waiting if the login process exits prematurely.
        if ! kill -0 "$LOGIN_PID" 2>/dev/null; then
            # Re-check: it may have printed the message right before exiting.
            grep -Fq "Initiating login..." "$LOGIN_LOG" && found=1
            break
        fi

        sleep 1
    done

    # Terminate the CLI without waiting for browser-based authentication.
    stop_process_tree "$LOGIN_PID"
    wait "$LOGIN_PID" 2>/dev/null || true
    LOGIN_PID=""

    if [[ "$found" -ne 1 ]]; then
        echo "Meshery authentication was not initiated."
        cat "$LOGIN_LOG"
        return 1
    fi

    cat "$LOGIN_LOG"
}

@test "Given Meshery is deployed, When Local is selected interactively, Then authentication should succeed" {
    start_meshery

    run run_expect '
        set timeout 30
        log_user 1

        spawn $env(MESHERYCTL_BIN) system login

        # Wait for the provider selection menu.
        expect {
            -re {Select a Provider} {}
            timeout {
                puts "Provider selection menu did not appear"
                exit 2
            }
            eof {
                puts "Login process exited before displaying the provider menu"
                exit 2
            }
        }

        # Navigate until Local becomes the highlighted provider,
        # instead of relying on its position in the list.
        set found 0

        for {set i 0} {$i < 15 && !$found} {incr i} {
            expect {
                -re {▸[^\r\n]*\yLocal\y} {
                    set found 1
                }
                -re {▸[^\r\n]*} {
                    send "\033\[B"
                    after 300
                }
                timeout {
                    puts "Could not determine the highlighted provider"
                    exit 3
                }
                eof {
                    puts "Login process exited during provider selection"
                    exit 3
                }
            }
        }

        if {!$found} {
            puts "Local provider was not found in the selection menu"
            exit 4
        }

        # Select Local and wait for authentication to complete.
        send "\r"

        expect {
            eof {}
            timeout {
                puts "Login did not finish after selecting Local"
                exit 6
            }
        }

        # wait returns: pid spawn_id os_error_flag value
        # os_error_flag is 0 when value is the real exit status.
        set wait_result [wait]

        if {[lindex $wait_result 2] != 0} {
            puts "Failed to retrieve login process exit status"
            exit 5
        }

        exit [lindex $wait_result 3]
    ' 60

    assert_success
    assert_output --partial "authenticated"
}

@test "Given Meshery is deployed, When Local is selected explicitly, Then authentication should succeed" {
    start_meshery

    run "$MESHERYCTL_BIN" system login -p Local

    assert_success
    assert_output --partial "authenticated"
}

@test "Given Meshery is deployed, When the Meshery provider is selected, Then authentication should be initiated" {
    start_meshery

    run verify_meshery_login_initiation

    assert_success
    assert_output --partial "Initiating login..."
}

@test "Given Meshery is deployed, When an unknown provider is selected, Then authentication should fail" {
    start_meshery

    run "$MESHERYCTL_BIN" system login -p Unknown

    assert_failure
    assert_output --partial "the specified provider 'Unknownn' is not available. "
}
