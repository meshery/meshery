#!/usr/bin/env bats
setup() {
    load "$E2E_HELPERS_PATH/tests_helpers"
	_tests_helper
}

@test "mesheryctl system help is succeeded and display help" {
    run $MESHERYCTL_BIN system --help
    [ "$status" -eq 0 ]
    assert_line "Use \"mesheryctl system [command] --help\" for more information about a command."
}

