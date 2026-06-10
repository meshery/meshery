#!/usr/bin/env bash


# Descirption: Centralization of helper functions for BATS Core and libraries

_start_port_forward() {
    if [[ "$MESHERY_PLATFORM" != "kubernetes" ]]; then
        return
    fi
    echo "start: Port forwarding"
    nohup kubectl -n meshery port-forward svc/meshery 9081:9081 > /dev/null 2>&1 &
    export MESHERY_SERVER_PORT_FORWARD_PID="$!"
    echo "done: Port forwarding"
}

_wait_for_meshery_server() {
    local timeout_seconds=${1:-180}
    local elapsed=0

    echo "Waiting for Meshery server to be ready..."
    until curl --silent --fail http://localhost:9081/api/system/version > /dev/null 2>&1; do
        if (( elapsed >= timeout_seconds )); then
            echo "Timed out waiting for Meshery server" >&2
            return 1
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done
    echo "Meshery server is ready"
}

_ensure_meshery_test_environment() {
    _start_port_forward
    _wait_for_meshery_server
}

_load_bats_libraries() {
    export BATS_LIB_PATH=${BATS_LIB_PATH:-"/usr/lib"}
    
    # Loading BATS libraries
    bats_load_library bats-support
    bats_load_library bats-assert
    bats_load_library bats-file

    # Required due to difference between local setup and ci
    if [ -f "$BATS_LIB_PATH/bats-detik/lib/detik.bash" ]
    then
        # For local setup
        bats_load_library bats-detik/lib/detik.bash
    else
        # for ci
        bats_load_library bats-detik/detik.bash
    fi

    _ensure_meshery_test_environment
}
