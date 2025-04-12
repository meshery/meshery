#!/usr/bin/env bats
setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries
}

@test "mesheryctl system help is succeeded and display help" {
    run $MESHERYCTL_BIN system --help
    assert_success
    
    assert_line "Use \"mesheryctl system [command] --help\" for more information about a command."
}

