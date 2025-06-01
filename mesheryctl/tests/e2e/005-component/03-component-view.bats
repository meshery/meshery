#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures/component-search"
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
  run $MESHERYCTL_BIN component view ACL --output-format yaml
  assert_success
  assert_output --partial "apiVersion"
}

@test "view command displays JSON output for a known component" {
  run mesh$MESHERYCTL_BINryctl component view ACL --output-format json
  assert_success
  assert_output --partial "\"apiVersion\"" 
}



