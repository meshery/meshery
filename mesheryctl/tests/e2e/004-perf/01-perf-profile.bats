#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries
}

@test "mesheryctl perf profile lists available profiles including test-profile" {
    run $MESHERYCTL_BIN perf profile

    assert_success

    # Check for table headers
    assert_line --partial "NAME"
    assert_line --partial "ID"
    assert_line --partial "RESULTS"
    assert_line --partial "LOAD-GENERATOR"
    assert_line --partial "LAST-RUN"

    # Check for the test-profile row by name only
    # Using the name test-profile since it was created during the previous test
    assert_line --partial "test-profile"
}

@test "mesheryctl perf profile returns message when profile is not found" {
    run $MESHERYCTL_BIN perf profile not-found-profile

    assert_success
    assert_line --partial "No Performance Profiles to display"
}



