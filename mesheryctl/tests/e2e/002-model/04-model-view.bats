#!/usr/bin/env bats

setup() {
  # Set the directory where mesheryctl is located if needed.
  MESHERYCTL_DIR=$(dirname "$MESHERYCTL_BIN")
  export TESTDATA_DIR="$MESHERYCTL_DIR/tests/e2e/002-model/testdata/model-view"
}

# Utility function to check if a particular string appears in the actual output.
check_output() {
  local expected="$1"
  run grep -Fq "$expected" <<< "$actual_output"
  if [ "$status" -ne 0 ]; then
    echo "FAILED: Expected message not found: '$expected'" >&2
    echo "Actual output:" >&2
    echo "$actual_output" >&2
  fi
  [ "$status" -eq 0 ]
}

@test "mesheryctl model view displays usage instructions when no model name is provided" {
  run $MESHERYCTL_BIN model view
  [ "$status" -ne 0 ]
  actual_output=$(echo "$output")
  check_output "model name isn't specified"
  check_output "Usage: mesheryctl model view [model-name]"
}

@test "mesheryctl model view displays an existing model" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN model view amd-gpu"
  [ "$status" -eq 0 ]
  echo "$output" > actual_output.txt
  diff actual_output.txt "$TESTDATA_DIR/exp_out_existing_model.txt"
}


@test "mesheryctl model view handles non-existent models gracefully" {
  run $MESHERYCTL_BIN model view non-existent-model
  [ "$status" -eq 0 ]
  actual_output=$(echo "$output")
  check_output "No model(s) found for the given name  non-existent-model"
}

@test "mesheryctl model view supports JSON output" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN model view amd-gpu -o json"
  [ "$status" -eq 0 ]
  echo "$output" > actual_output.txt
  diff actual_output.txt "$TESTDATA_DIR/exp_out_json.txt"
}

@test "mesheryctl model view supports YAML output" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN model view amd-gpu -o yaml"
  [ "$status" -eq 0 ]
  echo "$output" > actual_output.txt
  diff actual_output.txt "$TESTDATA_DIR/exp_out_yaml.txt"
}

teardown() {
  rm -f actual_output.txt
}