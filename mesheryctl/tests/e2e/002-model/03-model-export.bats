#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/model"
  mkdir -p "$TESTDATA_DIR"

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
}



@test "mesheryctl model export displays usage instructions when no model name provided" {
  run $MESHERYCTL_BIN model export
  assert_failure
  
  assert_output --partial "Please provide a model name."
  assert_output --partial "Usage: mesheryctl model export [model-name ]"
}

@test "mesheryctl model export succeeds with default options" {
  run $MESHERYCTL_BIN model export accurate -l "$TESTDATA_DIR"
  assert_success

  assert_output --partial "Exported model to $TESTDATA_DIR"
  assert_file_exists "$TESTDATA_DIR/accurate.tar"
}

@test "mesheryctl model export succeeds with tar output type" {
  run $MESHERYCTL_BIN model export accurate -l $TESTDATA_DIR -o tar
  assert_success

  assert_output --partial "Exported model to $TESTDATA_DIR"
  
  assert_file_exists "$TESTDATA_DIR/accurate.tar.gz"

}

@test "mesheryctl model export succeeds with json output format" {
  run $MESHERYCTL_BIN model export accurate -t json -l $TESTDATA_DIR
  assert_success

  assert_output --partial "Exported model to $TESTDATA_DIR"
  assert_file_exists "$TESTDATA_DIR/accurate.tar"
}

@test "mesheryctl model export includes version when specified" {
  run $MESHERYCTL_BIN model export accurate --version v1.7.0 -l $TESTDATA_DIR
  assert_success

  assert_output --partial "Exported model to $TESTDATA_DIR"
  assert_file_exists "$TESTDATA_DIR/accurate.tar"
}

@test "mesheryctl model export handles discard flags correctly" {
  run $MESHERYCTL_BIN model export accurate $TESTDATA_DIR --discard-components --discard-relationships -l $TESTDATA_DIR
  assert_success

 assert_output --partial "Exported model to $TESTDATA_DIR"
 assert_file_exists "$TESTDATA_DIR/accurate.tar"

}
