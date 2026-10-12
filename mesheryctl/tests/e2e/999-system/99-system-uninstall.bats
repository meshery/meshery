#!/usr/bin/env bats

# Runs in the 999-system suite so it executes after every suite that needs a
# running Meshery: `system uninstall` removes the whole installation, leaving
# nothing behind for the suites that follow.

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    export UNINSTALL_START_MSG="Uninstalling Meshery resources"
    export UNINSTALL_SUCCESS_MSG="Meshery uninstalled"
}

@test "mesheryctl system uninstall with invalid arguments error" {
    run $MESHERYCTL_BIN system uninstall invalid-arguments
    assert_failure
    assert_output --partial "this command takes no arguments"
}

@test "mesheryctl system uninstall completely removes Meshery" {
    run $MESHERYCTL_BIN system uninstall -y
    assert_success
    assert_output --partial "$UNINSTALL_START_MSG"
    assert_output --partial "$UNINSTALL_SUCCESS_MSG"

    # Full teardown: no Meshery component is left running afterwards
    run $MESHERYCTL_BIN system status -y
    assert_success
    assert_output --partial "Meshery is not running"
}
