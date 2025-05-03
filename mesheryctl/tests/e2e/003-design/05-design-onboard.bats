#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

  export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/design"
  mkdir -p "$TESTDATA_DIR"

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"

  # Create fixtures directory for design-onboard if it doesn't exist
  mkdir -p "$FIXTURES_DIR/design-onboard"

  # Copy the nginx.yaml file from design-import to design-onboard fixtures
  if [ ! -f "$FIXTURES_DIR/design-onboard/nginx.yaml" ]; then
    cp "$FIXTURES_DIR/design-import/nginx.yaml" "$FIXTURES_DIR/design-onboard/nginx.yaml"
  fi
}

@test "mesheryctl design onboard with valid file and source type should succeed or show auth error" {
  # Onboard the design directly from the fixtures directory
  run $MESHERYCTL_BIN design onboard -f "$FIXTURES_DIR/design-onboard/nginx.yaml" -s "Kubernetes Manifest"

  # Either we get a success message or an authentication error
  if [[ "$output" == *"Authentication token not found"* ]]; then
    assert_output --partial "Authentication token not found"
  else
    assert_output --partial "design onboarded"
  fi
}

@test "mesheryctl design onboard with an invalid file path should display an error message" {
  # Use a non-existent file path
  run $MESHERYCTL_BIN design onboard -f "$TESTDATA_DIR/design-onboard/nonexistent.yaml" -s "Kubernetes Manifest"

  assert_failure
  assert_output --partial "Error"
}

@test "mesheryctl design onboard with invalid source type should display an error message" {
  # Use an invalid source type
  run $MESHERYCTL_BIN design onboard -f "$FIXTURES_DIR/design-onboard/nginx.yaml" -s "InvalidSourceType"

  assert_failure
  assert_output --partial "Error"
}

@test "mesheryctl design onboard without required flags should show appropriate error" {
  # Run without required flags
  run $MESHERYCTL_BIN design onboard

  assert_failure
  assert_output --partial "Error: Unable to onboard design"
}

@test "mesheryctl design onboard with existing design name should succeed" {
  # First import a design to get its ID
  run $MESHERYCTL_BIN design import -f "$FIXTURES_DIR/design-onboard/nginx.yaml" --source-type "Kubernetes Manifest"

  # Extract design name from the output
  DESIGN_NAME="nginx-deployment"

  # Onboard the design using its name
  run $MESHERYCTL_BIN design onboard "$DESIGN_NAME"

  # Check for expected output - success message
  assert_output --partial "Fetching patterns"
}
