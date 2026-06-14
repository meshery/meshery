#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/design"
  mkdir -p "$TESTDATA_DIR"

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
}

@test "mesheryctl design import -f nginx.yaml displays imported and output the design Id" {
  # Import the design directly from the testdata directory
  run $MESHERYCTL_BIN design import -f "$FIXTURES_DIR/design-import/nginx.yaml" --source-type "Kubernetes Manifest"
  assert_success
  assert_output --partial "imported" || assert_output --partial "Design ID" || assert_output --partial "saved"

  # Extract and store design ID for subsequent tests
  DESIGN_ID=$(echo "$output" | grep -o '[0-9a-f]\{8\}')
  echo "$DESIGN_ID" > "$TESTDATA_DIR/id"
}

@test "mesheryctl design import with an invalid path displays an error message" {
  # Use a non-existent file path
  run $MESHERYCTL_BIN design import -f "$TESTDATA_DIR/design-import/nonexistent.yaml" --source-type "Kubernetes Manifest"
  # TODO: Update command to assert is failing
  # assert_failure
  assert_output --partial "Error" || assert_output --partial "no such file"
}
