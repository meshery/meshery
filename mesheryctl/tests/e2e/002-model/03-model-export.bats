#!/usr/bin/env bats

load helpers/test_helpers.bats

# -----------------------------------------------------------------------------
# This file tests the `mesheryctl model export` command for:
#   1) Missing model name (usage instructions)
#   2) Default export options (oci output, yaml format)
#   3) Overriding file type with "-o tar"
#   4) Overriding output format with "-t json"
#   5) Including a version with "--version"
#   6) Handling discard flags for components and relationships
# -----------------------------------------------------------------------------

setup() {
  MESHERYCTL_DIR=$(dirname "$MESHERYCTL_BIN")
  export TESTDATA_DIR="$MESHERYCTL_DIR/tests/e2e/002-model/testdata/model-export"
}

@test "mesheryctl model export displays usage instructions when no model name provided" {
  run $MESHERYCTL_BIN model export
  [ "$status" -ne 0 ]
  actual_output=$(echo "$output")

  check_output "Please provide a model name."
  check_output "Usage: mesheryctl model export [model-name]"
}

@test "mesheryctl model export succeeds with default options" {
  run $MESHERYCTL_BIN model export $TESTDATA_DIR
  [ "$status" -eq 0 ]
  actual_output=$(echo "$output")

  check_output "api/meshmodels/export?name=$TESTDATA_DIR"
  check_output "output_format=yaml"
  check_output "file_type=oci"
  check_output "components=true"
  check_output "relationships=true"
}

@test "mesheryctl model export succeeds with tar output type" {
  run $MESHERYCTL_BIN model export  $TESTDATA_DIR -o tar
  [ "$status" -eq 0 ]
  actual_output=$(echo "$output")

  check_output "file_type=tar"
}

@test "mesheryctl model export succeeds with json output format" {
  run $MESHERYCTL_BIN model export $TESTDATA_DIR -t json
  [ "$status" -eq 0 ]
  actual_output=$(echo "$output")

  check_output "output_format=json"
}

@test "mesheryctl model export includes version when specified" {
  run $MESHERYCTL_BIN model export $TESTDATA_DIR --version v0.7.3
  [ "$status" -eq 0 ]
  actual_output=$(echo "$output")

  check_output "&version=v0.7.3"
}

@test "mesheryctl model export handles discard flags correctly" {
  run $MESHERYCTL_BIN model export $TESTDATA_DIR --discard-components --discard-relationships
  [ "$status" -eq 0 ]
  actual_output=$(echo "$output")

  check_output "components=false"
  check_output "relationships=false"
}
