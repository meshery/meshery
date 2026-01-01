#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries
}

# bats test_tags=system:token
@test "mesheryctl system token list displays available tokens" {
    run $MESHERYCTL_BIN system token list
    assert_success
    
    assert_line --partial "Available tokens:"
}

# bats test_tags=system:token
@test "mesheryctl system token list shows default token" {
    run $MESHERYCTL_BIN system token list
    assert_success
    
    assert_line --partial "default"
}

# bats test_tags=system:token  
@test "mesheryctl system token list --help displays help" {
    run $MESHERYCTL_BIN system token list --help
    assert_success
    
    assert_line --partial "List all the tokens"
}

# bats test_tags=system:token
@test "mesheryctl system token list rejects invalid arguments" {
    run $MESHERYCTL_BIN system token list invalid-arg
    assert_failure
    
    assert_output --partial "accepts 0 arg(s), received 1"
}
