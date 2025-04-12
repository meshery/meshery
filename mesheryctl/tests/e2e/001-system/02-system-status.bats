#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

}

@test "mesheryctl system status display mesehry operator running pods" {
    run $MESHERYCTL_BIN system status -y
    assert_success

    assert_line --regexp "meshery(|-(broker|meshsync|operator))"
    assert_line --regexp "Running"
}
