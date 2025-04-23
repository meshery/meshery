#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

}

# TODO: Investigate to handle 'Invalid adapter name provided'
# @test "mesheryctl system check is succeeded" {
#    run $MESHERYCTL_BIN system check
#    assert_success
# }
