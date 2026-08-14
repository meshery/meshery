#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries
    
    # Verify MESHERYCTL_BIN is set
    if [[ -z "$MESHERYCTL_BIN" ]]; then
        echo "Error: MESHERYCTL_BIN is not defined. Set it before running tests."
        exit 1
    fi
}

teardown() {
    # Clean up any running Meshery instances started during this test
    $MESHERYCTL_BIN system stop 2>/dev/null || true
    sleep 1
}

@test "[TC-0001][tg=Installation] install mesheryctl and start Meshery on an invalid platform" {
    # Test that installing mesheryctl with an invalid platform fails gracefully
    
    run bash -c 'curl -L https://meshery.io/install | PLATFORM=bob bash - < /dev/null'
    
    assert_output --partial "curl: Failed writing body"
}

@test "[TC-0002][tg=Installation] install mesheryctl and start Meshery on a valid platform" {
    # Test that installing mesheryctl with an valid platform succeed
    
    run bash -c 'curl -L https://meshery.io/install | PLATFORM=docker bash - < /dev/null'
    
    assert_success  
}

@test "[TC-0003][tg=Installation] mesheryctl system start with invalid platform flag should handle gracefully" {
    # Test that an invalid platform parameter is rejected
    run $MESHERYCTL_BIN system start -p invalidplatform

    assert_failure
}
