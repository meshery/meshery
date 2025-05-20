#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries
  MESHERYCTL_DIR=$(dirname "$MESHERYCTL_BIN")
  export TESTDATA_DIR="$MESHERYCTL_DIR/tests/e2e/002-model/testdata/model-view"
}

@test "mesheryctl model view displays usage instructions when no model name is provided" {
  run $MESHERYCTL_BIN model view
  [ "$status" -ne 0 ]
  actual_output=$(echo "$output")
  assert_output --partial "Error: model name isn't specified"
  assert_output --partial "Usage: mesheryctl model view [model-name]"
}

@test "mesheryctl model view displays an existing model" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN model view amd-gpu | grep -Ev 'created_at|updated_at|deleted_at'"
  [ "$status" -eq 0 ]
  assert_output --partial "$(cat "$TESTDATA_DIR/exp_out_existing_model.txt")"
}


@test "mesheryctl model view handles non-existent models gracefully" {
  run $MESHERYCTL_BIN model view non-existent-model
  [ "$status" -eq 0 ]
  assert_output --partial "No model(s) found for the given name  non-existent-model"
}

@test "mesheryctl model view supports JSON output" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN model view amd-gpu -o json | grep -Ev 'created_at|updated_at|deleted_at'"
  [ "$status" -eq 0 ]
  assert_output --partial "$(cat "$TESTDATA_DIR/exp_out_json.txt")"
}

@test "mesheryctl model view supports YAML output" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN model view amd-gpu -o yaml | grep -Ev 'created_at|updated_at|deleted_at'"
  [ "$status" -eq 0 ]
  assert_output --partial "$(cat "$TESTDATA_DIR/exp_out_yaml.txt")"
}