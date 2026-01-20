#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries
}

@test "mesheryctl perf apply fails when new profile and URL is missing" {
    run $MESHERYCTL_BIN perf apply dummy-test-profile -y

    assert_failure
    assert_line --partial "Unable to get URL for performing test"
}

@test "mesheryctl perf apply runs successfully with test-profile and valid URL" {
    run $MESHERYCTL_BIN perf apply test-profile --url "https://google.com" -y

    assert_success
    assert_line --partial "Initiating Performance test"
    assert_line --partial "Test Completed"
}

@test "mesheryctl perf apply uses specified load generator wrk2" {
    run $MESHERYCTL_BIN perf apply test-profile --load-generator wrk2

    assert_success
    assert_line --partial "Initiating Performance test"
    assert_line --partial "Test Completed"
}

@test "mesheryctl perf apply with URL and mesh istio runs successfully" {
    run $MESHERYCTL_BIN perf apply test-profile --url https://192.168.1.15/productpage --mesh istio

    assert_success
    assert_line --partial "Initiating Performance test"
    assert_line --partial "Test Completed"
}
