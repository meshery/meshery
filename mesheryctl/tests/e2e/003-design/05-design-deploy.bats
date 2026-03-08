#!/usr/bin/env bats

setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
  _load_bats_libraries

  export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/design"
  mkdir -p "$TESTDATA_DIR"

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
  export SHARED_FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures/design-import"
}

@test "given --file and --source flags provided when mesheryctl design deploy --file file-path --source source-type then design is deployed" {
  run $MESHERYCTL_BIN design deploy -f "$SHARED_FIXTURES_DIR/nginx.yaml" -s "Kubernetes Manifest"

  assert_success
  assert_output --partial "design deployed"
}

@test "given an invalid file path when mesheryctl design deploy --file invalid-file-path --source source-type then an error message is displayed" {
  run $MESHERYCTL_BIN design deploy -f "$TESTDATA_DIR/nonexistent.yaml" -s "Kubernetes Manifest"

  assert_failure
  assert_output --partial "Invalid value for --file '$TESTDATA_DIR/nonexistent.yaml'"
}

@test "given an invalid source type when mesheryctl design deploy --file file-path --source invalid-source-type then an error message is displayed" {
  run $MESHERYCTL_BIN design deploy -f "$SHARED_FIXTURES_DIR/nginx.yaml" -s "InvalidSourceType"

  assert_failure
  assert_output --partial "Invalid value for --source-type 'InvalidSourceType'"
  assert_output --partial "Invalid flag(s) provided"
  assert_output --partial "Provide valid flag value and try again"
}

@test "given no required flags when mesheryctl design deploy then an appropriate error is displayed" {
  run $MESHERYCTL_BIN design deploy

  assert_failure
  assert_output --partial "Unable to deploy design"
}
