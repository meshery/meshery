setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"
}

@test "given no model-id provided as an argument when running mesheryctl model delete then an error message is displayed" {
    run $MESHERYCTL_BIN model delete
    assert_failure
    assert_output --partial "[ model-id ] is required"
    assert_output --partial "Error"
}

@test "given an invalid model-id is provided as an argument when running mesheryctl model delete invalid-id then an error message is displayed" {
    INVALID_ID="0000"

    run $MESHERYCTL_BIN model delete "$INVALID_ID"
    assert_failure
    assert_output --partial "Invalid ID format"
    assert_output --partial "Error"
}

@test "given a non existing model-id is provided as an argument when running mesheryctl model delete non-existing-id then an error message is displayed" {
    NONEXISTENT_ID="00000000-0000-0000-0000-000000000000"

    run $MESHERYCTL_BIN model delete "$NONEXISTENT_ID"
    assert_failure
    assert_output --partial "Failed to delete model"
    assert_output --partial "Error"
}
