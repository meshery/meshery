#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

  export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/design"
  mkdir -p "$TESTDATA_DIR"

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
  export SHARED_FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures/design-shared"
}

@test "mesheryctl design onboard with valid file and source type should succeed" {
  run $MESHERYCTL_BIN design onboard -f "$SHARED_FIXTURES_DIR/nginx.yaml" -s "Kubernetes Manifest"

  assert_success
  assert_output --partial "design onboarded"
}

@test "mesheryctl design onboard with valid file should show auth error when not authenticated" {
  # Temporarily move token file if it exists
  if [ -f "$HOME/.meshery/config" ]; then
    mv "$HOME/.meshery/config" "$HOME/.meshery/config.bak"
  fi

  run $MESHERYCTL_BIN design onboard -f "$SHARED_FIXTURES_DIR/nginx.yaml" -s "Kubernetes Manifest"

  # Restore token file if it was moved
  if [ -f "$HOME/.meshery/config.bak" ]; then
    mv "$HOME/.meshery/config.bak" "$HOME/.meshery/config"
  fi

  assert_output --partial "Authentication token not found"
}

@test "mesheryctl design onboard with an invalid file path should display an error message" {
  run $MESHERYCTL_BIN design onboard -f "$TESTDATA_DIR/nonexistent.yaml" -s "Kubernetes Manifest"

  assert_failure
  assert_output --partial "Error: unable to read file"
}

@test "mesheryctl design onboard with invalid source type should display an error message" {
  run $MESHERYCTL_BIN design onboard -f "$SHARED_FIXTURES_DIR/nginx.yaml" -s "InvalidSourceType"

  assert_failure
  assert_output --partial "Error: invalid source type"
}

@test "mesheryctl design onboard without required flags should show appropriate error" {
  run $MESHERYCTL_BIN design onboard

  assert_failure
  assert_output --partial "Error: Unable to onboard design"
}

@test "mesheryctl design onboard with existing design name should succeed" {
  run $MESHERYCTL_BIN design import -f "$SHARED_FIXTURES_DIR/nginx.yaml" --source-type "Kubernetes Manifest"

  DESIGN_NAME="nginx-deployment"
  run $MESHERYCTL_BIN design onboard "$DESIGN_NAME"

  assert_success
  assert_output --partial "design onboarded"
}