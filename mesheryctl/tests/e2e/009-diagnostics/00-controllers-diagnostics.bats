#!/usr/bin/env bats
#
# End-to-end: kubernetes connection controller diagnostics + broker recovery.
#
# Exercises the full manual flow as an automated test:
#   1. find the kubernetes connection (Meshery autoconnects the current context)
#   2. ensure MeshSync is in operator mode (embedded mode yields no broker diagnostics)
#   3. GET /api/system/controllers/diagnostics  -> broker_unreachable (baseline)
#   4. port-forward svc/meshery-nats so Meshery can reach + authenticate to NATS
#   5. GET diagnostics again -> healthy, broker_unreachable cleared (CONNECTED)
#
# Preconditions (skips cleanly if unmet):
#   - a running Meshery server on :9081 with ~/.meshery/auth.json populated
#   - started in operator mode: MESHSYNC_DEFAULT_DEPLOYMENT_MODE=operator
#   - meshery-operator + meshsync + statefulset/meshery-nats deployed in ns "meshery"
#   - jq, kubectl, curl on PATH
#
# For a DETERMINISTIC unreachable baseline (step 3), start the server with the
# managed port-forward disabled: MESHERY_MANAGED_BROKER_PORTFORWARD=false . When
# it is enabled (default), Meshery reaches NATS on its own; the baseline assertion
# is then skipped and only the recovered/healthy end state is asserted.
#
# Run just this test:
#   cd mesheryctl && make e2e-no-build BATS_FOLDER_PATTERN=009-diagnostics \
#     BATS_FILE_PATTERN='00-controllers-diagnostics'

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
    load "$E2E_HELPERS_PATH/constants"

    export MESHERY_URL="${MESHERY_SERVER_URL:-http://localhost:9081}"
    export MESHERY_NS="meshery"
    export MESHERY_AUTH_FILE="${MESHERY_AUTH_FILE:-$HOME/.meshery/auth.json}"
}

# _need skips the test unless the given command is available.
_need() {
    command -v "$1" >/dev/null 2>&1 || skip "$1 is required for this test"
}

# _diagnostics <connectionId> -> prints the diagnostics JSON, sending the provider
# session cookies mesheryctl uses (token + meshery-provider from auth.json).
_diagnostics() {
    local conn_id="$1"
    local token provider
    token=$(jq -r '.token // empty' "$MESHERY_AUTH_FILE")
    provider=$(jq -r '."meshery-provider" // "Meshery"' "$MESHERY_AUTH_FILE")
    curl -s --cookie "token=${token}; meshery-provider=${provider}" \
        "${MESHERY_URL}/api/system/controllers/diagnostics?connectionId=${conn_id}"
}

@test "kubernetes connection broker diagnostics recover once NATS is reachable" {
    _need jq
    _need kubectl
    _need curl

    # --- preconditions ---
    kubectl -n "$MESHERY_NS" get statefulset meshery-nats >/dev/null 2>&1 \
        || skip "statefulset/meshery-nats not deployed in ns ${MESHERY_NS}"
    # Quick readiness check (skip fast rather than block); this test assumes the
    # environment is already up. If NATS is crash-looping, note the operator bug:
    # it injects the NATS token unquoted into nats.conf; it must be token: "$NATS_TOKEN".
    local nats_ready
    nats_ready=$(kubectl -n "$MESHERY_NS" get statefulset meshery-nats -o jsonpath='{.status.readyReplicas}' 2>/dev/null)
    [ "${nats_ready:-0}" -ge 1 ] 2>/dev/null \
        || skip "meshery-nats not Ready (readyReplicas=${nats_ready:-0}); if crash-looping, the operator injects the NATS token unquoted into nats.conf"

    # --- 1. resolve the kubernetes connection id ---
    run "$MESHERYCTL_BIN" connection list --kind kubernetes
    assert_success
    local conn_id
    conn_id=$(echo "$output" | grep -oiE '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' | head -n1)
    [ -n "$conn_id" ] || skip "no kubernetes connection found; connect the current context first"

    # --- 2. sanity: diagnostics endpoint is reachable + authed, returns valid JSON ---
    local diag
    diag=$(_diagnostics "$conn_id")
    echo "diagnostics(initial)=$diag"
    echo "$diag" | jq -e '.connectionId' >/dev/null \
        || skip "diagnostics endpoint not reachable/authed (is the server up on ${MESHERY_URL} with a valid auth.json?)"

    # Embedded MeshSync has no in-cluster broker, so no broker diagnostics apply.
    if echo "$diag" | jq -e '.diagnostics[]? | select(.code=="connection_inactive")' >/dev/null; then
        skip "connection is not active yet (no live session); connect it and retry"
    fi

    # --- 3. baseline (best-effort, deterministic only when managed forward disabled) ---
    if echo "$diag" | jq -e '.healthy == false and (.diagnostics[]?.code=="broker_unreachable")' >/dev/null; then
        run bash -c "echo '$diag' | jq -e '.diagnostics[] | select(.code==\"broker_unreachable\") | .remediation | length > 0'"
        assert_success   # broker_unreachable carries remediation steps
    fi

    # --- 4. make NATS reachable via a port-forward ---
    kubectl -n "$MESHERY_NS" port-forward svc/meshery-nats 4222:4222 >/dev/null 2>&1 &
    local pf_pid=$!
    # shellcheck disable=SC2064
    trap "kill ${pf_pid} 2>/dev/null || true" RETURN

    # --- 5. poll until the broker recovers (self-healing NATS retry reconnects) ---
    local healthy="false" i
    for i in $(seq 1 30); do
        diag=$(_diagnostics "$conn_id")
        if echo "$diag" | jq -e '.healthy == true and ((.diagnostics // []) | map(select(.code=="broker_unreachable")) | length == 0)' >/dev/null; then
            healthy="true"
            break
        fi
        sleep 2
    done
    echo "diagnostics(after port-forward)=$diag"

    assert_equal "$healthy" "true"
    # broker_unreachable must be gone
    run bash -c "echo '$diag' | jq -e '(.diagnostics // []) | map(select(.code==\"broker_unreachable\")) | length == 0'"
    assert_success
}
