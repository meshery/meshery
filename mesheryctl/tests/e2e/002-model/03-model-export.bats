#!/usr/bin/env bats

setup() {
  
 load "$E2E_HELPERS_PATH/tests_helpers"
	_tests_helper
  MESHERYCTL_DIR=$(dirname "$MESHERYCTL_BIN")
  export TESTDATA_DIR="$MESHERYCTL_DIR/tests/e2e/002-model/testdata/"

  if [ ! -d "$TESTDATA_DIR" ]; then
    mkdir -p "$TESTDATA_DIR"
  fi
}

teardown() {
  # Remove the model-export directory created in setup, if it exists
  if [ -d "$TESTDATA_DIR" ]; then
    rm -rf "$TESTDATA_DIR"
  fi
}

@test "mesheryctl model export displays usage instructions when no model name provided" {
  run $MESHERYCTL_BIN model export
  [ "$status" -ne 0 ]
  
  assert_output --partial "Please provide a model name."
  assert_output --partial "Usage: mesheryctl model export [model-name ]"
}

@test "mesheryctl model export succeeds with default options" {
  run $MESHERYCTL_BIN model export accurate -l $TESTDATA_DIR
  [ "$status" -eq 0 ]

  assert_output --partial "Exported model to $TESTDATA_DIR"
  if [ ! -f "$TESTDATA_DIR/accurate.tar" ]; then
    echo "Expected file accurate.tar was not found in $TESTDATA_DIR"
    exit 1
  fi
}

@test "mesheryctl model export succeeds with tar output type" {
  run $MESHERYCTL_BIN model export accurate -l $TESTDATA_DIR -o tar
  [ "$status" -eq 0 ]

  assert_output --partial "Exported model to $TESTDATA_DIR"
  if [ ! -f "$TESTDATA_DIR/accurate.tar.gz" ]; then
    echo "Expected file accurate.tar was not found in $TESTDATA_DIR"
    exit 1
  fi

}

@test "mesheryctl model export succeeds with json output format" {
  run $MESHERYCTL_BIN model export accurate -t json -l $TESTDATA_DIR
  [ "$status" -eq 0 ]

  assert_output --partial "Exported model to $TESTDATA_DIR"
  if [ ! -f "$TESTDATA_DIR/accurate.tar" ]; then
    echo "Expected file accurate.tar was not found in $TESTDATA_DIR"
    exit 1
  fi
}

@test "mesheryctl model export includes version when specified" {
  run $MESHERYCTL_BIN model export accurate --version v1.7.0 -l $TESTDATA_DIR
  [ "$status" -eq 0 ]

  assert_output --partial "Exported model to $TESTDATA_DIR"
  if [ ! -f "$TESTDATA_DIR/accurate.tar" ]; then
    echo "Expected file accurate.tar was not found in $TESTDATA_DIR"
    exit 1
  fi
}

@test "mesheryctl model export handles discard flags correctly" {
  run $MESHERYCTL_BIN model export accurate $TESTDATA_DIR --discard-components --discard-relationships
  [ "$status" -eq 0 ]

 assert_output --partial "Exported model to $TESTDATA_DIR"
 if [ ! -f "$TESTDATA_DIR/accurate.tar" ]; then
    echo "Expected file accurate.tar was not found in $TESTDATA_DIR"
    exit 1
  fi

}
