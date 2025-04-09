#!/usr/bin/env bats

# Load required BATS libraries
load '../helpers/bats-support/load'
load '../helpers/bats-assert/load'

# Test 1: Verify successful import of a valid design
@test "mesheryctl design import with nginx YAML is succeeded" {
  # Copy the test fixture to TEMP_TEST_DATA_DIR
  cp "$(dirname "$BATS_TEST_FILENAME")/../testdata/design-import/nginx.yaml" "$TEMP_TEST_DATA_DIR/nginx.yaml"
  
  # Define the path to the test fixture in TEMP_TEST_DATA_DIR
  local TESTDATA_PATH="$TEMP_TEST_DATA_DIR/nginx.yaml"

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
