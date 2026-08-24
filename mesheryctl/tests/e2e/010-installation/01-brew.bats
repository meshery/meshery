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

    run brew list --versions mesheryctl
    assert_success
    local brew_version
    brew_version=$(echo "$output" | awk '{print $2}')

    run "$brew_prefix/bin/mesheryctl" version
    assert_success
    assert_line --regexp "Client.*$brew_version"
}
