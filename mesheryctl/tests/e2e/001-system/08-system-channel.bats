#!/usr/bin/env bats

# Only the v1.0.x releases publish Helm charts, so this is the version used wherever a
# pinned release must actually be deployed. Override with CHANNEL_E2E_PINNED_VERSION.
PINNED_VERSION="${CHANNEL_E2E_PINNED_VERSION:-v1.0.71}"

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
}

# Record the original channel/version once, so every test can be reverted to it.
# Exported variables are visible to every test and to teardown.
setup_file() {
    run "$MESHERYCTL_BIN" system channel view
    if [ "$status" -ne 0 ]; then
        echo "# setup_file: 'system channel view' failed (status $status): $output" >&3
        return 1
    fi

    export ORIGINAL_CHANNEL=$(echo "$output" | sed -n 's/^Channel: //p' | head -n1)
    export ORIGINAL_VERSION=$(echo "$output" | sed -n 's/^Version: //p' | head -n1)

    # Without both fields there is nothing safe to restore; stop rather than guess a default.
    if [ -z "$ORIGINAL_CHANNEL" ] || [ -z "$ORIGINAL_VERSION" ]; then
        echo "# setup_file: could not read channel and version from: $output" >&3
        return 1
    fi

    _save_deployment_env || return 1
}

# Deployment settings that CI patches into the in-cluster meshery deployment (kubectl set env)
# and that later suites (connection, diagnostics) depend on. A channel switch reinstalls the
# Helm release and drops them, so they are recorded here and re-applied after each redeploy.
DEPLOYMENT_ENV_NAMES="PROVIDER_BASE_URLS MESHSYNC_DEFAULT_DEPLOYMENT_MODE MESHERY_MANAGED_BROKER_PORTFORWARD"

# Platform of the current mesheryctl context, read from the meshconfig.
# Requires mikefarah yq v4 (the Python yq cannot parse these expressions).
_current_platform() {
    local config="${MESHERY_CONFIG_FILE_PATH:-$HOME/.meshery/config.yaml}"

    if ! yq --version 2>&1 | grep -q mikefarah; then
        echo "# setup_file: mikefarah yq v4 is required (found: $(yq --version 2>&1))" >&3
        return 1
    fi

    # Read the context name first: indexing .contexts with a missing current-context does not
    # yield null, it silently returns another context's platform.
    local context
    context=$(yq '.current-context' "$config") || return 1
    if [ -z "$context" ] || [ "$context" = "null" ]; then
        echo "# setup_file: no current-context set in $config" >&3
        return 1
    fi

    CTX="$context" yq '.contexts[strenv(CTX)].platform' "$config"
}

_save_deployment_env() {
    local platform
    if ! platform=$(_current_platform) || [ -z "$platform" ] || [ "$platform" = "null" ]; then
        echo "# setup_file: could not determine the platform of the current context" >&3
        return 1
    fi

    # Only a Kubernetes deployment carries settings to preserve; Docker has no cluster to query.
    [ "$platform" = "kubernetes" ] || return 0

    if ! command -v kubectl >/dev/null 2>&1; then
        echo "# setup_file: platform is kubernetes but kubectl is not installed" >&3
        return 1
    fi

    # Only NotFound means the deployment is absent. Timeouts, authorization failures and
    # unreachable clusters must fail setup rather than silently capture nothing.
    local get_err
    if ! get_err=$(kubectl -n meshery get deployment/meshery 2>&1 >/dev/null); then
        if echo "$get_err" | grep -q "NotFound"; then
            return 0
        fi
        echo "# setup_file: could not look up deployment/meshery: $get_err" >&3
        return 1
    fi

    local env_list name
    if ! env_list=$(kubectl -n meshery set env deployment/meshery --list); then
        echo "# setup_file: could not capture deployment/meshery environment" >&3
        return 1
    fi
    for name in $DEPLOYMENT_ENV_NAMES; do
        echo "$env_list" | grep "^$name=" >> "$BATS_FILE_TMPDIR/deployment_env" || true
    done
}

# Re-apply the recorded settings and wait for the rollout. A no-op when none were recorded
# (the file is missing or empty).
_restore_deployment_env() {
    [ -s "$BATS_FILE_TMPDIR/deployment_env" ] || return 0

    local settings=() line
    while IFS= read -r line; do
        settings+=("$line")
    done < "$BATS_FILE_TMPDIR/deployment_env"
    kubectl -n meshery rollout status deployment/meshery --timeout=180s || return 1
    kubectl -n meshery set env deployment/meshery "${settings[@]}" || return 1
    kubectl -n meshery rollout status deployment/meshery --timeout=180s
}

# Redeploy Meshery on the original channel. Retries transient failures (e.g. TLS timeouts)
# and, when a failed switch left a half-removed deployment, force-stops and starts it again.
_redeploy_original() {
    local original="$1" attempt
    for attempt in 1 2 3; do
        if "$MESHERYCTL_BIN" system channel switch "$original" -y; then
            _restore_deployment_env
            return
        fi
        echo "# teardown: redeploy attempt $attempt failed, cleaning up" >&3
        "$MESHERYCTL_BIN" system stop --force -y || true
        if "$MESHERYCTL_BIN" system start --skip-browser -y; then
            _restore_deployment_env || return 1
            "$MESHERYCTL_BIN" system dashboard --skip-browser
            return
        fi
        # Wait before the next attempt only; nothing follows the last one.
        [ "$attempt" -lt 3 ] && sleep 10
    done
    return 1
}

# Runs after every test: restore the original channel/version in the config and, if the test
# switched channels, redeploy Meshery. Any failure fails the test, and the "switched" marker is
# only removed once the redeploy succeeds, so the next teardown retries the recovery.
teardown() {
    local original="$ORIGINAL_CHANNEL-$ORIGINAL_VERSION"

    if ! "$MESHERYCTL_BIN" system channel set "$original"; then
        echo "# teardown: could not restore channel to $original" >&3
        return 1
    fi

    if [ -f "$BATS_FILE_TMPDIR/switched" ]; then
        if ! _redeploy_original "$original"; then
            echo "# teardown: Meshery could not be redeployed on $original; fix it manually" >&3
            return 1
        fi
        rm -f "$BATS_FILE_TMPDIR/switched"
    fi
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given a current context when running mesheryctl system channel view then its channel and version are displayed" {
    run "$MESHERYCTL_BIN" system channel view

    assert_success
    assert_line "Channel: $ORIGINAL_CHANNEL"
    assert_line "Version: $ORIGINAL_VERSION"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given the edge channel when running mesheryctl system channel set edge then the current context is set to edge" {
    run "$MESHERYCTL_BIN" system channel set edge

    assert_success
    assert_output --partial "Channel set to edge-latest"

    run "$MESHERYCTL_BIN" system channel view
    assert_success
    assert_line --regexp "^Channel: edge$"
    assert_line --regexp "^Version: latest$"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given the stable channel when running mesheryctl system channel set stable then the current context is set to stable" {
    run "$MESHERYCTL_BIN" system channel set stable

    assert_success
    assert_output --partial "Channel set to stable-latest"

    run "$MESHERYCTL_BIN" system channel view
    assert_success
    assert_line --regexp "^Channel: stable$"
    assert_line --regexp "^Version: latest$"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given an invalid channel when running mesheryctl system channel set bob then an error is returned" {
    run "$MESHERYCTL_BIN" system channel set bob

    assert_failure
    assert_output --partial "Invalid release channel"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given a stable release version when running mesheryctl system channel set stable-<version> then the current context is pinned to that version" {
    run "$MESHERYCTL_BIN" system channel set stable-$PINNED_VERSION

    assert_success
    assert_output --partial "Channel set to stable-$PINNED_VERSION"

    run "$MESHERYCTL_BIN" system channel view
    assert_success
    assert_line --regexp "^Channel: stable$"
    assert_line --regexp "^Version: $PINNED_VERSION$"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given a non-latest edge version when running mesheryctl system channel set edge-v0.5 then an error explains edge only supports latest" {
    run "$MESHERYCTL_BIN" system channel set edge-v0.5

    assert_failure
    assert_output --partial "Invalid edge release version"
    assert_output --partial "only supports 'latest'"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given an invalid stable version when running mesheryctl system channel set stable-v0.7.999 then an error is returned" {
    run "$MESHERYCTL_BIN" system channel set stable-v0.7.999

    assert_failure
    assert_output --partial "not a valid Meshery release"
}

# A channel switch replaces the running deployment; teardown restores the original one.
# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given a running deployment when switching to edge then Meshery is redeployed on edge" {
    touch "$BATS_FILE_TMPDIR/switched"
    run "$MESHERYCTL_BIN" system channel switch edge -y

    assert_success
    assert_output --partial "Channel set to edge-latest"

    run "$MESHERYCTL_BIN" system channel view
    assert_success
    assert_line "Channel: edge"
    assert_line "Version: latest"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given a running deployment when switching to stable then Meshery is redeployed on stable" {
    touch "$BATS_FILE_TMPDIR/switched"
    run "$MESHERYCTL_BIN" system channel switch stable -y

    assert_success
    assert_output --partial "Channel set to stable-latest"

    run "$MESHERYCTL_BIN" system channel view
    assert_success
    assert_line "Channel: stable"
    assert_line "Version: latest"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given a running deployment when switching to a pinned stable version then Meshery is redeployed at that release" {
    touch "$BATS_FILE_TMPDIR/switched"
    run "$MESHERYCTL_BIN" system channel switch stable-$PINNED_VERSION -y

    assert_success
    assert_output --partial "Channel set to stable-$PINNED_VERSION"

    run "$MESHERYCTL_BIN" system channel view
    assert_success
    assert_line "Channel: stable"
    assert_line "Version: $PINNED_VERSION"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given a non-latest edge version when running mesheryctl system channel switch edge-v05 then an error explains edge only supports latest" {
    run "$MESHERYCTL_BIN" system channel switch edge-v05 -y

    assert_failure
    assert_output --partial "Invalid edge release version"
    assert_output --partial "only supports 'latest'"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given an invalid subcommand when running mesheryctl system channel bob then an error is returned" {
    run "$MESHERYCTL_BIN" system channel bob

    assert_failure
    assert_output --partial "invalid subcommand"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given an invalid stable version when switching channel then the version is rejected" {
    run "$MESHERYCTL_BIN" system channel switch stable-v0.7.999 -y

    assert_failure
    assert_output --partial "not a valid Meshery release"
}