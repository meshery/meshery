#!/usr/bin/env bash


# Descirption: Centralization of helper functions for BATS Core and libraries

_refresh_meshery_auth_file() {
    if [[ -z "$MESHERY_PROVIDER_TOKEN" || -z "$MESHERY_AUTH_FILE" ]]; then
        return
    fi

    mkdir -p "$(dirname "$MESHERY_AUTH_FILE")"
    jq -nc --arg token "$MESHERY_PROVIDER_TOKEN" '{"meshery-provider":"Meshery","token":$token}' > "$MESHERY_AUTH_FILE"
}

_wait_for_meshery_server() {
    local timeout_seconds=${1:-180}
    local elapsed=0
    local ready_marker="${TEMP_DATA_DIR:-/tmp}/.meshery-server-ready"

    if [[ -f "$ready_marker" ]]; then
        return
    fi

    until curl --silent --show-error --fail http://localhost:9081/api/system/version > /dev/null; do
        if (( elapsed >= timeout_seconds )); then
            echo "Timed out waiting for Meshery server to become ready" >&2
            return 1
        fi

        sleep 2
        elapsed=$((elapsed + 2))
    done

    touch "$ready_marker"
}

_wait_for_meshery_kubernetes_components() {
    local timeout_seconds=${1:-180}
    local elapsed=0
    local ready_marker="${TEMP_DATA_DIR:-/tmp}/.meshery-k8s-ready"

    if [[ "$MESHERY_PLATFORM" != "kubernetes" ]]; then
        return
    fi

    if ! command -v kubectl >/dev/null 2>&1; then
        return
    fi

    if [[ -f "$ready_marker" ]]; then
        return
    fi

    until pods=$(kubectl get pods -n meshery --no-headers 2>/dev/null) \
        && grep -Eq '^meshery-[a-z0-9-]+[[:space:]]' <<<"$pods" \
        && grep -Eq '^meshery-broker-[0-9]+[[:space:]]' <<<"$pods" \
        && grep -Eq '^meshery-meshsync-[a-z0-9-]+[[:space:]]' <<<"$pods" \
        && grep -Eq '^meshery-operator-[a-z0-9-]+[[:space:]]' <<<"$pods"; do
        if (( elapsed >= timeout_seconds )); then
            echo "Timed out waiting for Meshery Kubernetes components to become available" >&2
            kubectl get pods -n meshery >&2 || true
            return 1
        fi

        sleep 2
        elapsed=$((elapsed + 2))
    done

    touch "$ready_marker"
}

_ensure_meshery_test_environment() {
    _refresh_meshery_auth_file
    _wait_for_meshery_server
    _wait_for_meshery_kubernetes_components
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
