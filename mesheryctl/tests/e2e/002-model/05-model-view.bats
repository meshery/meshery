#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries
  MESHERYCTL_DIR=$(dirname "$MESHERYCTL_BIN")
  REQUIRED_FIELDS='
    .id != null and
    .schemaVersion != null and
    .displayName != null and
    .status != null and
    .subCategory != null and
    .model != null and
    .name != null and
    .version != null and
    .registrant.id != null and
    .category.id != null and
    .relationships_count != null and
    .components_count != null and
    has(\"components\") and
    has(\"relationships\")
    '
  export TESTDATA_DIR="$MESHERYCTL_DIR/tests/e2e/002-model/testdata/model-view"
}

@test "mesheryctl model view displays usage instructions when no model name is provided" {
  run $MESHERYCTL_BIN model view
  
  assert_failure
  
  assert_output --partial "Error: model name isn't specified"
  assert_output --partial "Usage: mesheryctl model view [model-name]"
}

@test "mesheryctl model view displays an existing model" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN model view amd-gpu -o json \
    | yq -e \"$REQUIRED_FIELDS\""

  assert_success
  assert_output "true"
}

@test "mesheryctl model view handles non-existent models gracefully" {
  run $MESHERYCTL_BIN model view non-existent-model

  assert_success
  assert_output --partial "No model(s) found for the given name  non-existent-model"
}

@test "mesheryctl model view supports JSON output" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN model view amd-gpu -o json \
    | jq -e \"$REQUIRED_FIELDS\""

  assert_success
  assert_output "true"
}

@test "mesheryctl model view supports YAML output" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN model view amd-gpu -o json \
    | yq -e \"$REQUIRED_FIELDS\""

  assert_success
  assert_output "true"
}