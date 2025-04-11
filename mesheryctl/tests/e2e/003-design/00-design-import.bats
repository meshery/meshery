#!/usr/bin/env bats

# Setup function to load libraries and prepare test fixtures
setup() {
  # Load libraries
  load '../helpers/bats-support/load'
  load '../helpers/bats-assert/load'
  cp "$(dirname "$BATS_TEST_FILENAME")/../testdata/design-import/nginx.yaml" "$TEMP_TEST_DATA_DIR/nginx.yaml"
  local TESTDATA_PATH="$TEMP_TEST_DATA_DIR/nginx.yaml"
}

# Test 1: Verify successful import of a valid design
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
  mkdir -p "${TEMP_TEST_DATA_DIR}/design"
  echo "$DESIGN_ID" > "${TEMP_TEST_DATA_DIR}/design/id"
}

# Test 2: Verify failure when importing an invalid file
@test "mesheryctl design import for invalid file fails" {
  # Use a non-existent file path
  run $MESHERYCTL_BIN design import -f "${TEMP_TEST_DATA_DIR}/nonexistent.yaml" --source-type "Kubernetes Manifest"
  assert_failure
  assert_output --partial "Error" || assert_output --partial "no such file"
}
