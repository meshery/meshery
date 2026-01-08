#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
}

# bats test_tags=system:token
@test "mesheryctl system token list displays available tokens including default" {
    run $MESHERYCTL_BIN system token list
    assert_success

    assert_line --partial "Available tokens:"
    assert_line --partial "default"
}

# bats test_tags=system:token
@test "mesheryctl system token list --help displays help" {
    run $MESHERYCTL_BIN system token list --help
    assert_success

    assert_line --partial "List all the tokens"
}

# bats test_tags=system:token
@test "mesheryctl system token list with one argument failed" {
    run $MESHERYCTL_BIN system token list invalid-arg
    assert_failure

    assert_output --partial "accepts 0 arg(s), received 1"
}

# bats test_tags=system:token
@test "mesheryctl system token create successfully creates a token" {
    run $MESHERYCTL_BIN system token create test-token
    assert_success

    assert_output --partial "Token test-token created"

    # Verify token appears in list
    run $MESHERYCTL_BIN system token list
    assert_success
    assert_output --partial "test-token"

    # Teardown: delete the token to ensure test isolation
    run $MESHERYCTL_BIN system token delete test-token
    assert_success
}

# bats test_tags=system:token
@test "mesheryctl system token create --help displays help" {
    run $MESHERYCTL_BIN system token create --help
    assert_success

    assert_output --partial "Create the token"
}

# bats test_tags=system:token
@test "mesheryctl system token create without token name fails" {
    run $MESHERYCTL_BIN system token create
    assert_failure

    assert_output --partial "token name is required"
}

# bats test_tags=system:token
@test "mesheryctl system token delete successfully removes a token" {
    # Setup: Create a token to delete
    run $MESHERYCTL_BIN system token create test-delete-token
    assert_success

    # Delete the token
    run $MESHERYCTL_BIN system token delete test-delete-token
    assert_success

    assert_output --partial "Token test-delete-token deleted"

    # Verify token no longer appears in list
    run $MESHERYCTL_BIN system token list
    assert_success
    refute_output --partial "test-delete-token"
}

# bats test_tags=system:token
@test "mesheryctl system token delete --help displays help" {
    run $MESHERYCTL_BIN system token delete --help
    assert_success

    assert_output --partial "Delete the token"
}

# bats test_tags=system:token
@test "mesheryctl system token delete without token name fails" {
    run $MESHERYCTL_BIN system token delete
    assert_failure

    assert_output --partial "token name is required"
}

# bats test_tags=system:token
@test "mesheryctl system token full lifecycle: create, verify, delete, verify" {
    # Create token
    run $MESHERYCTL_BIN system token create lifecycle-test-token
    assert_success
    assert_output --partial "Token lifecycle-test-token created"

    # Verify appears in list
    run $MESHERYCTL_BIN system token list
    assert_success
    assert_output --partial "lifecycle-test-token"

    # Delete token
    run $MESHERYCTL_BIN system token delete lifecycle-test-token
    assert_success
    assert_output --partial "Token lifecycle-test-token deleted"

    # Verify no longer in list
    run $MESHERYCTL_BIN system token list
    assert_success
    refute_output --partial "lifecycle-test-token"
}