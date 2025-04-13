#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/tests_helpers"
	_tests_helper
}

# Basic tests to validate cli has been build and behave properly at root

@test "cli is available" {
    if [[ -z "$MESHERYCTL_BIN" ]]; then
        echo "Error: MESHERYCTL_BIN is not defined. Set it before running tests."
        exit 1
    fi
    
    run $MESHERYCTL_BIN
    assert_success
}

@test "mesheryctl version return Client and Server" {
    run $MESHERYCTL_BIN version
    assert_success
    
    assert_line --regexp "(Client|Server)"
}