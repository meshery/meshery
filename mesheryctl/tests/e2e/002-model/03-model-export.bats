#!/usr/bin/env bats

setup() {
  E2ETEST_DIR=$(dirname "$BATS_TEST_DIRNAME")
  load "$E2ETEST_DIR/$SUPPORT_DESTDIR"
  load "$E2ETEST_DIR/$ASSERT_DESTDIR"  
  MESHERYCTL_DIR=$(dirname "$MESHERYCTL_BIN")
  export TESTDATA_DIR="$MESHERYCTL_DIR/tests/e2e/002-model/testdata/model-export/"
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
}

@test "mesheryctl model export succeeds with tar output type" {
  run $MESHERYCTL_BIN model export accurare -o tar
  [ "$status" -eq 0 ]

  assert_output --partial "file_type=tar"
}

@test "mesheryctl model export succeeds with json output format" {
  run $MESHERYCTL_BIN model export accurare -t json
  [ "$status" -eq 0 ]

  assert_output --partial "output_format=json"
}

@test "mesheryctl model export includes version when specified" {
  run $MESHERYCTL_BIN model export accurate --version v0.7.3
  [ "$status" -eq 0 ]

  assert_output --partial "&version=v0.7.3"
}

@test "mesheryctl model export handles discard flags correctly" {
  run $MESHERYCTL_BIN model export accurate --discard-components --discard-relationships
  [ "$status" -eq 0 ]

  assert_output --partial "components=false"
  assert_output --partial "relationships=false"
}
