#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"
}

@test "Count components" {
    run $MESHERYCTL_BIN exp relationship search --kind edge --page 1
    assert_success
}

@test "Count " {
    run $MESHERYCTL_BIN exp relationship search --model meshery-shapes --page 1
    assert_success
}

@test "Count dd" {
    run $MESHERYCTL_BIN exp relationship search --subtype reference --page 1
    assert_success
}
