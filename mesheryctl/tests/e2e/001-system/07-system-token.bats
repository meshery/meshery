#!/usr/bin/env bats

# Shared token name to allow state to flow across ordered E2E tests.
TOKEN="e2e-test-token"

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
}

# Best-effort cleanup to keep the environment consistent between runs.
cleanup_token_best_effort() {
    run "$MESHERYCTL_BIN" system token delete "$TOKEN"
    true
}

# Ensure a clean starting state for the suite.
setup_file() {
    cleanup_token_best_effort
}

# Ensure no residual state is left after the suite completes.
teardown_file() {
    cleanup_token_best_effort
}

# bats test_tags=system:token
@test "mesheryctl system token create fails without token name" {
    run "$MESHERYCTL_BIN" system token create
    assert_failure
}

# bats test_tags=system:token
@test "mesheryctl system token create succeeds" {
    run "$MESHERYCTL_BIN" system token create "$TOKEN"
    assert_success
}

# bats test_tags=system:token
@test "mesheryctl system token list displays available tokens" {
    run "$MESHERYCTL_BIN" system token list
    assert_success
    assert_output --partial "Available tokens"
}

# bats test_tags=system:token
@test "mesheryctl system token list includes created token" {
    run "$MESHERYCTL_BIN" system token list
    assert_success
    assert_output --partial "$TOKEN"
}

# bats test_tags=system:token
@test "mesheryctl system token delete fails for non existing token" {
    run "$MESHERYCTL_BIN" system token delete "__non_existing_token__"
    assert_failure
}

# bats test_tags=system:token
@test "mesheryctl system token delete succeeds for created token" {
    run "$MESHERYCTL_BIN" system token delete "$TOKEN"
    assert_success
}

# bats test_tags=system:token
@test "mesheryctl system token list does not include deleted token" {
    run "$MESHERYCTL_BIN" system token list
    assert_success
    refute_output --partial "$TOKEN"
}
