#!/usr/bin/env bats

# Only the v1.0.x releases publish Helm charts, so this is the version used wherever a
# pinned release must actually be deployed. Override with CHANNEL_E2E_PINNED_VERSION.
PINNED_VERSION="${CHANNEL_E2E_PINNED_VERSION:-v1.0.71}"

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
}

# Record the original channel/version once, so every test can be reverted to it.
setup_file() {
    run "$MESHERYCTL_BIN" system channel view
    export ORIGINAL_CHANNEL=$(echo "$output" | sed -n 's/^Channel: //p' | head -n1)
    export ORIGINAL_VERSION=$(echo "$output" | sed -n 's/^Version: //p' | head -n1)

    # A crashed earlier run may have left a test value in the config; never treat it as the default.
    case "$ORIGINAL_VERSION" in
        v0.5.0|v0.7.999|"$PINNED_VERSION") export ORIGINAL_VERSION="latest" ;;
    esac
    [ -n "$ORIGINAL_CHANNEL" ] || export ORIGINAL_CHANNEL="stable"
    [ -n "$ORIGINAL_VERSION" ] || export ORIGINAL_VERSION="latest"

    echo "$ORIGINAL_CHANNEL-$ORIGINAL_VERSION" > "$BATS_FILE_TMPDIR/original_channel"
}

# Redeploy Meshery on the original channel. Retries transient failures (e.g. TLS timeouts)
# and, when a failed switch left a half-removed deployment, force-stops and starts it again.
_redeploy_original() {
    local original="$1" attempt
    for attempt in 1 2 3; do
        if "$MESHERYCTL_BIN" system channel switch "$original" -y; then
            return 0
        fi
        echo "# teardown: redeploy attempt $attempt failed, cleaning up" >&3
        "$MESHERYCTL_BIN" system stop --force -y || true
        if "$MESHERYCTL_BIN" system start --skip-browser -y; then
            return 0
        fi
        sleep 10
    done
    return 1
}

# Runs after every test and never fails it: restore the original channel/version in the config,
# and if the test switched channels, repair and redeploy Meshery so the next test starts clean.
teardown() {
    local original
    original=$(cat "$BATS_FILE_TMPDIR/original_channel" 2>/dev/null)
    [ -n "$original" ] || return 0

    "$MESHERYCTL_BIN" system channel set "$original" || \
        echo "# teardown: could not restore channel to $original" >&3

    if [ -f "$BATS_FILE_TMPDIR/switched" ]; then
        rm -f "$BATS_FILE_TMPDIR/switched"
        _redeploy_original "$original" || \
            echo "# teardown: Meshery could not be redeployed on $original; fix it manually" >&3
    fi
    return 0
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

# A channel switch replaces the running deployment; teardown_file restores the original one.
# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given a running deployment when switching to edge then Meshery is redeployed on edge" {
    touch "$BATS_FILE_TMPDIR/switched"
    run "$MESHERYCTL_BIN" system channel switch edge -y

    assert_success
    assert_output --partial "Channel set to edge-latest"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given a running deployment when switching to stable then Meshery is redeployed on stable" {
    touch "$BATS_FILE_TMPDIR/switched"
    run "$MESHERYCTL_BIN" system channel switch stable -y

    assert_success
    assert_output --partial "Channel set to stable-latest"
}

# bats test_tags=system:channel
@test "[cut=System][tg=System Channel] given a running deployment when switching to a pinned stable version then Meshery is redeployed at that release" {
    touch "$BATS_FILE_TMPDIR/switched"
    run "$MESHERYCTL_BIN" system channel switch stable-$PINNED_VERSION -y

    assert_success
    assert_output --partial "Channel set to stable-$PINNED_VERSION"
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
