#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures/component-search"
  run $MESHERYCTL_BIN model import -f $FIXTURES_DIR/valid-model
  export TESTDATA_DIR="$BATS_TEST_DIRNAME/testdata/component-view"
}

@test "view command fails with no arguments" {
  run $MESHERYCTL_BIN component view
  assert_failure
  assert_output --partial "Error: [component name] is required but not specified"
  assert_output --partial "Usage: mesheryctl component view [component-name]"
  assert_output --partial "Run 'mesheryctl component view --help' to see detailed help message"
}

@test "view command fails with too many arguments" {
  run $MESHERYCTL_BIN component view comp1 comp2
  assert_failure
  assert_output --partial "Error: too many arguments specified"
  assert_output --partial "Usage: mesheryctl component view [component-name]"
  assert_output --partial "Run 'mesheryctl component view --help' to see detailed help message"
}

@test "view command fails with invalid output format" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN component view ACL --output-format xml"
  assert_failure
  assert_output --partial "Error: output-format choice is invalid or not provided, use [json|yaml]"
}

@test "view command displays YAML output for a known component" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN component view component_cli-e2e-test -o yaml | grep -Ev 'created_at|updated_at|deleted_at'"
  assert_success
  assert_output --partial "$(cat "$TESTDATA_DIR/exp_out_yaml.txt")"
}

@test "view command displays JSON output for a known component" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN component view component_cli-e2e-test -o json | grep -Ev 'created_at|updated_at|deleted_at'"
  assert_success
  assert_output --partial "$(cat "$TESTDATA_DIR/exp_out_json.txt")"
}

@test "view command saves YAML output file" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN component view component_cli-e2e-test -o yaml --save"
  assert_success
  assert_output --partial "Saving output as YAML file"  

  assert_file_exist "$HOME/.meshery/component_component_cli-e2e-test.yaml"
}

@test "view command saves JSON output file" {
  run bash -c "printf '\n' | $MESHERYCTL_BIN component view component_cli-e2e-test -o json --save"
  assert_success
  assert_output --partial "Saving output as JSON file"  

  assert_file_exist "$HOME/.meshery/component_component_cli-e2e-test.json"
}
