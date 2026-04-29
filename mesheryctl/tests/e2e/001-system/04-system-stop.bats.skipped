#!/usr/bin/env bats 

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    run $MESHERYCTL_BIN system start -y
    assert_success
    
    export STOP_START_MSG="Stopping Meshery resources"
    export STOP_SUCCESS_MSG="Meshery resources are stopped"
    export RESET_MSG="reset"
}

@test "mesheryctl system stop with no arguments succeeds when running" {
    run $MESHERYCTL_BIN system stop -y
    assert_success
    assert_output --partial "$STOP_START_MSG"
    assert_output --partial "$STOP_SUCCESS_MSG"
}

@test "mesheryctl system stop --reset stops and resets config" {
    run $MESHERYCTL_BIN system stop --reset -y
    assert_success 
    assert_output --partial "$STOP_SUCCESS_MSG"
    assert_output --partial "$RESET_MSG"
}

@test "mesheryctl system stop --force handles forceful stop" {
    run $MESHERYCTL_BIN system stop --force -y
    assert_success
    assert_output --partial "$STOP_START_MSG"
    assert_output --partial "$STOP_SUCCESS_MSG"
}

@test "mesheryctl system stop with invalid arguments error" {
    run $MESHERYCTL_BIN system stop invalid-arguments
    assert_failure
    assert_output --partial "this command takes no arguments" 
}
