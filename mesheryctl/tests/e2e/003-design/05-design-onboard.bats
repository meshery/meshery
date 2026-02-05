#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

  export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/design"
  mkdir -p "$TESTDATA_DIR"

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
  export SHARED_FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures/design-import"
}

@test "given --file and --source flags provided when mesheryctl design onboard --file file-path --source source-type then design is onboarded" {
  run $MESHERYCTL_BIN design onboard -f "$SHARED_FIXTURES_DIR/nginx.yaml" -s "Kubernetes Manifest"

  assert_success
  assert_output --partial "design onboarded"
}

@test "given an invalid file path when mesheryctl design onboard --file invalid-file-path --source source-type then an error message is displayed" {
  run $MESHERYCTL_BIN design onboard -f "$TESTDATA_DIR/nonexistent.yaml" -s "Kubernetes Manifest"

  assert_failure
  assert_output --partial "Error: unable to read file"
}

@test "given an invalid source type when mesheryctl design onboard --file file-path --source invalid-source-type then an error message is displayed" {
  run $MESHERYCTL_BIN design onboard -f "$SHARED_FIXTURES_DIR/nginx.yaml" -s "InvalidSourceType"

  assert_failure
  assert_output --partial "Error: invalid source type"
}

@test "given no required flags when mesheryctl design onboard then an appropriate error is displayed" {
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