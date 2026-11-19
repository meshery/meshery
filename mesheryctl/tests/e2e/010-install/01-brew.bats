#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
}

@test "[TC-9][tg=Installation] Install mesheryctl via Homebrew" {
    if ! command -v brew &> /dev/null; then
        skip "brew is not installed on this system"
    fi

    run brew install mesheryctl
    assert_success

    run brew --prefix mesheryctl
    assert_success
    local brew_prefix="$output"

    local keg_path
    keg_path=$(cd "$brew_prefix" && pwd -P)
    local brew_version="${keg_path##*/}"
    brew_version="${brew_version%_*}"

    run "$brew_prefix/bin/mesheryctl" version
    assert_success
    local client_version
    client_version=$(echo "$output" | awk '/^Client/ {print $2}')
    assert_equal "$client_version" "v$brew_version"
}
