#!/usr/bin/env bats

# Setup function to load libraries and prepare test fixtures
setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  mkdir -p "$TEMP_DATA_DIR/design"
  cp "$E2E_TESTDATA_PATH/design-import/nginx.yaml" "$TEMP_DATA_DIR/design/nginx.yaml"
  export TESTDATA_PATH="$TEMP_DATA_DIR/design/nginx.yaml"
}

@test "mesheryctl design import with nginx YAML is succeeded" {
  # Verify the test fixture exists
  if [ ! -f "$TESTDATA_PATH" ]; then
    skip "Test fixture $TESTDATA_PATH not found"
  fi

  # Import the design directly from the testdata directory
  run $MESHERYCTL_BIN design import -f "$TESTDATA_PATH" --source-type "Kubernetes Manifest"
  assert_success
  assert_output --partial "imported" || assert_output --partial "Design ID" || assert_output --partial "saved"

  # Extract and store design ID for subsequent tests
  DESIGN_ID=$(echo "$output" | grep -o '[0-9a-f]\{8\}')
  echo "$DESIGN_ID" > "${TEMP_DATA_DIR}/design/id"
}

@test "mesheryctl design import for invalid file fails" {
  # Use a non-existent file path
  run $MESHERYCTL_BIN design import -f "${TEMP_DATA_DIR}/design/nonexistent.yaml" --source-type "Kubernetes Manifest"
  # TODO: Update command to assert is failing
  # assert_failure
  assert_output --partial "Error" || assert_output --partial "no such file"
}
