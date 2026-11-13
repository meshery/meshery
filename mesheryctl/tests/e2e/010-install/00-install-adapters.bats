#!/usr/bin/env bats

# Test Plan #6 (Installation / Installer / Bash Script):
#   curl -L https://meshery.io/install | ADAPTERS=<adapter> bash -
# installs mesheryctl, asks for a platform and starts Meshery with a new
# context that holds only the requested adapter.
#
# The published installer runs unmodified. Only its side effects outside the
# test are replaced: the release download is served from a local archive, the
# binary goes to a temporary install directory, and `mesheryctl system start`
# is recorded instead of deploying Meshery. Every other mesheryctl call goes to
# the CLI built from this checkout, so it validates the components the
# installer passes to `system context create`. Each run gets its own HOME and
# install directory, so no context leaks into the next run or the rest of the
# suite.

setup_file() {
    export INSTALL_SCRIPT="$BATS_FILE_TMPDIR/install.sh"
    curl -fsSL "${MESHERY_INSTALL_SCRIPT_URL:-https://meshery.io/install}" -o "$INSTALL_SCRIPT"

    REAL_MESHERYCTL_BIN="$(cd "$(dirname "$MESHERYCTL_BIN")" && pwd)/$(basename "$MESHERYCTL_BIN")"
    export REAL_MESHERYCTL_BIN

    # The "released" mesheryctl: logs each call, records `system start`
    # without running it and hands every other call to the real CLI.
    local archive_dir="$BATS_FILE_TMPDIR/archive"
    mkdir -p "$archive_dir"
    cat > "$archive_dir/mesheryctl" <<'EOF'
#!/usr/bin/env bash
echo "$*" >> "$MESHERYCTL_CALL_LOG"
if [ "${1:-} ${2:-}" = "system start" ]; then
    exit 0
fi
exec "$REAL_MESHERYCTL_BIN" "$@"
EOF
    chmod +x "$archive_dir/mesheryctl"
    export RELEASE_ARCHIVE="$BATS_FILE_TMPDIR/mesheryctl.tar.gz"
    tar -czf "$RELEASE_ARCHIVE" -C "$archive_dir" mesheryctl

    # curl stand-in for the installer: logs the URL, writes the archive to -o.
    export STUB_BIN_DIR="$BATS_FILE_TMPDIR/stub-bin"
    mkdir -p "$STUB_BIN_DIR"
    cat > "$STUB_BIN_DIR/curl" <<'EOF'
#!/usr/bin/env bash
out=
while [ $# -gt 0 ]; do
    case "$1" in
        -o) out="$2"; shift ;;
        http*) echo "$1" >> "$CURL_CALL_LOG" ;;
    esac
    shift
done
cp "$RELEASE_ARCHIVE" "$out"
EOF
    chmod +x "$STUB_BIN_DIR/curl"
}

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    PLATFORM_UNDER_TEST="${MESHERY_PLATFORM:-kubernetes}"
}

# run_installer <workdir> [VAR=value ...]
# Runs the installer from <workdir> with a fresh HOME and install directory,
# with ADAPTERS and PLATFORM taken only from the arguments. Set
# INSTALL_LAUNCHER to prefix the command (e.g. "setsid -w"), or
# INSTALL_TTY_INPUT to run it on a pseudo-terminal that this text is typed
# into, which is where the installer reads its prompt answers from.
run_installer() {
    local workdir="$1"
    shift
    INSTALL_HOME="$workdir/home"
    INSTALL_DIR="$workdir/bin"
    INSTALLED_CONFIG="$INSTALL_HOME/.meshery/config.yaml"
    MESHERYCTL_CALL_LOG="$workdir/mesheryctl.log"
    CURL_CALL_LOG="$workdir/curl.log"
    mkdir -p "$INSTALL_HOME" "$INSTALL_DIR"
    cd "$workdir"

    local cmd=(env -u ADAPTERS -u PLATFORM
        HOME="$INSTALL_HOME"
        PATH="$INSTALL_DIR:$STUB_BIN_DIR:$PATH"
        MESHERY_INSTALL_DIR="$INSTALL_DIR"
        MESHERY_VERSION=v0.0.0-e2e
        MESHERYCTL_CALL_LOG="$MESHERYCTL_CALL_LOG"
        CURL_CALL_LOG="$CURL_CALL_LOG"
        "$@"
        bash "$INSTALL_SCRIPT")

    if [ -n "${INSTALL_TTY_INPUT:-}" ]; then
        run script -qefc "$(printf '%q ' "${cmd[@]}")" /dev/null <<< "$INSTALL_TTY_INPUT"
    else
        run ${INSTALL_LAUNCHER:-} "${cmd[@]}" < /dev/null
    fi
}

# assert_installed_with_components <json array>
assert_installed_with_components() {
    assert_file_executable "$INSTALL_DIR/mesheryctl"
    assert_file_contains "$CURL_CALL_LOG" "/releases/download/v0.0.0-e2e/mesheryctl_0.0.0-e2e_"

    run yq '.["current-context"]' "$INSTALLED_CONFIG"
    assert_output "custom"
    run yq -o=json -I=0 '.contexts.custom.components' "$INSTALLED_CONFIG"
    assert_output "$1"

    run cat "$MESHERYCTL_CALL_LOG"
    assert_line "system start --yes --platform $PLATFORM_UNDER_TEST"
}

@test "[TC-6][cut=Installer / Bash Script][tg=Installation] given each supported adapter in ADAPTERS when running the installer then mesheryctl is installed and started with a new context holding only that adapter" {
    local mapping=(
        consul:meshery-consul
        istio:meshery-istio
        linkerd:meshery-linkerd
        nsm:meshery-nsm
        kuma:meshery-kuma
        traefik:meshery-traefik-mesh
        cilium:meshery-cilium
    )
    local pair adapter component
    for pair in "${mapping[@]}"; do
        adapter="${pair%%:*}"
        component="${pair#*:}"
        echo "checking ADAPTERS=$adapter (expect $component)"

        run_installer "$BATS_TEST_TMPDIR/$adapter" ADAPTERS="$adapter" PLATFORM="$PLATFORM_UNDER_TEST"
        assert_success
        assert_installed_with_components "[\"$component\"]"
    done
}

@test "[TC-6][cut=Installer / Bash Script][tg=Installation] given an unsupported adapter next to a supported one in ADAPTERS when running the installer then the unsupported one is skipped with a warning" {
    run_installer "$BATS_TEST_TMPDIR/mixed" ADAPTERS="unknown consul" PLATFORM="$PLATFORM_UNDER_TEST"
    assert_success
    assert_output --partial "WARNING: Adapter 'unknown' is not supported by this installer and will be skipped."

    assert_installed_with_components '["meshery-consul"]'
}

@test "[TC-6][cut=Installer / Bash Script][tg=Installation] given only unsupported adapters in ADAPTERS when running the installer then it exits with an error before installing mesheryctl" {
    run_installer "$BATS_TEST_TMPDIR/unsupported" ADAPTERS="unknown" PLATFORM="$PLATFORM_UNDER_TEST"
    assert_failure
    assert_output --partial "ERROR: No supported adapters were specified. Cannot create a custom context."

    assert_file_not_exists "$CURL_CALL_LOG"
    assert_file_not_exists "$INSTALL_DIR/mesheryctl"
}

@test "[TC-6][cut=Installer / Bash Script][tg=Installation] given an empty ADAPTERS when running the installer then Meshery is started without creating a custom context" {
    run_installer "$BATS_TEST_TMPDIR/empty" ADAPTERS="" PLATFORM="$PLATFORM_UNDER_TEST"
    assert_success
    assert_file_executable "$INSTALL_DIR/mesheryctl"

    run cat "$MESHERYCTL_CALL_LOG"
    refute_line --partial "system context create"
    assert_line "system start --yes --platform $PLATFORM_UNDER_TEST"
}

@test "[TC-6][cut=Installer / Bash Script][tg=Installation] given ADAPTERS without PLATFORM when running the installer then it prompts for a platform before installing mesheryctl" {
    if ! command -v setsid > /dev/null; then
        skip "setsid is required to run the installer without a terminal"
    fi

    # Without a controlling terminal the prompt's read from /dev/tty fails, so
    # the installer gives up after its retries instead of waiting for input.
    INSTALL_LAUNCHER="setsid -w" run_installer "$BATS_TEST_TMPDIR/prompt" ADAPTERS="consul"
    assert_failure
    assert_output --partial "Enter a platform to deploy Meshery. Available platforms [docker, kubernetes]:"
    assert_output --partial "Too many invalid attempts. Please try again."

    assert_file_not_exists "$CURL_CALL_LOG"
    assert_file_not_exists "$INSTALL_DIR/mesheryctl"
}

@test "[TC-6][cut=Installer / Bash Script][tg=Installation] given ADAPTERS without PLATFORM when the platform prompt is answered then mesheryctl is installed and started with a new context holding only that adapter" {
    if ! command -v script > /dev/null; then
        skip "script (util-linux) is required to answer the prompt from a terminal"
    fi

    INSTALL_TTY_INPUT="$PLATFORM_UNDER_TEST" run_installer "$BATS_TEST_TMPDIR/answered" ADAPTERS="consul"
    assert_success
    assert_output --partial "Enter a platform to deploy Meshery. Available platforms [docker, kubernetes]:"
    refute_output --partial "Invalid platform"

    assert_installed_with_components '["meshery-consul"]'
}
