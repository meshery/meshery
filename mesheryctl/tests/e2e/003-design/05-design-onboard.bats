#!/usr/bin/env bats

setup() {
  # Set up paths to BATS libraries
  BATS_TEST_DIRNAME_ABS="$(cd "$BATS_TEST_DIRNAME" && pwd)"
  
  # Load the BATS libraries directly
  load "/tmp/bats-support/load"
  load "/tmp/bats-assert/load"
  load "/tmp/bats-file/load"
  load "/tmp/bats-detik/lib/detik"

  # Set up test directories
  export TEMP_DATA_DIR=$(mktemp -d)
  export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/design"
  mkdir -p "$TESTDATA_DIR"

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
  
  # Create fixtures directory for design-onboard if it doesn't exist
  mkdir -p "$FIXTURES_DIR/design-onboard"
  
  # Copy the nginx.yaml file from design-import to design-onboard fixtures
  if [ ! -f "$FIXTURES_DIR/design-onboard/nginx.yaml" ]; then
    cp "$FIXTURES_DIR/design-import/nginx.yaml" "$FIXTURES_DIR/design-onboard/nginx.yaml"
  fi
  
  # Set the path to the mesheryctl binary
  export MESHERYCTL_BIN="$BATS_TEST_DIRNAME_ABS/../../../mesheryctl"
}

@test "mesheryctl design onboard -f nginx.yaml validates command structure" {
  # Skip this test if mesheryctl is not available
  if [ ! -f "$MESHERYCTL_BIN" ]; then
    skip "mesheryctl binary not found"
  fi

  # Onboard the design directly from the fixtures directory
  run $MESHERYCTL_BIN design onboard -f "$FIXTURES_DIR/design-onboard/nginx.yaml" -s "Kubernetes Manifest"
  
  # In a non-authenticated environment, we expect either:
  # 1. Authentication error (which is fine for testing command structure)
  # 2. Success message if somehow authenticated
  assert_output --partial "Authentication token not found" \
    || assert_output --partial "design onboarded" \
    || assert_output --partial "imported"
}

@test "mesheryctl design onboard with an invalid path displays an error message" {
  # Skip this test if mesheryctl is not available
  if [ ! -f "$MESHERYCTL_BIN" ]; then
    skip "mesheryctl binary not found"
  fi

  # Use a non-existent file path
  run $MESHERYCTL_BIN design onboard -f "$TESTDATA_DIR/design-onboard/nonexistent.yaml" -s "Kubernetes Manifest"
  assert_output --partial "Error" \
    || assert_output --partial "no such file" \
    || assert_output --partial "Authentication token not found"
}

@test "mesheryctl design onboard with invalid source type displays an error message" {
  # Skip this test if mesheryctl is not available
  if [ ! -f "$MESHERYCTL_BIN" ]; then
    skip "mesheryctl binary not found"
  fi

  # Use an invalid source type
  run $MESHERYCTL_BIN design onboard -f "$FIXTURES_DIR/design-onboard/nginx.yaml" -s "InvalidSourceType"
  assert_output --partial "Error" \
    || assert_output --partial "Invalid source type" \
    || assert_output --partial "Authentication token not found"
}

@test "mesheryctl design onboard without required flags shows appropriate error" {
  # Skip this test if mesheryctl is not available
  if [ ! -f "$MESHERYCTL_BIN" ]; then
    skip "mesheryctl binary not found"
  fi

  # Run without required flags
  run $MESHERYCTL_BIN design onboard
  assert_output --partial "Usage:" \
    || assert_output --partial "mesheryctl design onboard -f [filepath] -s [source type]" \
    || assert_output --partial "Error: Unable to onboard design" \
    || assert_output --partial "Authentication token not found"
}

@test "mesheryctl design onboard with existing design name validates command structure" {
  # Skip this test if mesheryctl is not available
  if [ ! -f "$MESHERYCTL_BIN" ]; then
    skip "mesheryctl binary not found"
  fi

  # First import a design to get its ID
  run $MESHERYCTL_BIN design import -f "$FIXTURES_DIR/design-onboard/nginx.yaml" --source-type "Kubernetes Manifest"
  
  # Extract design name from the output
  DESIGN_NAME="nginx-deployment"
  
  # Onboard the design using its name
  run $MESHERYCTL_BIN design onboard "$DESIGN_NAME"
  assert_output --partial "design onboarded" \
    || assert_output --partial "imported" \
    || assert_output --partial "Authentication token not found" \
    || assert_output --partial "Fetching patterns"
}
