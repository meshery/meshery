#!/usr/bin/env bats

setup() {
  
 load "$E2E_HELPERS_PATH/tests_helpers"
	_tests_helper
  MESHERYCTL_DIR=$(dirname "$MESHERYCTL_BIN")
  export TESTDATA_DIR="$MESHERYCTL_DIR/tests/e2e/002-model/testdata/"

  [ -d "$TESTDATA_DIR" ] || mkdir -p "$TESTDATA_DIR"
}

teardown() {
 [ -d "$TESTDATA_DIR" ] && rm -rf "$TESTDATA_DIR"
}

@test "mesheryctl model export displays usage instructions when no model name provided" {
  run $MESHERYCTL_BIN model export
  [ "$status" -ne 0 ]
  
  assert_output --partial "Please provide a model name."
  assert_output --partial "Usage: mesheryctl model export [model-name ]"
}

@test "mesheryctl model export succeeds with default options" {
  run $MESHERYCTL_BIN model export accurate -l $TESTDATA_DIR
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
