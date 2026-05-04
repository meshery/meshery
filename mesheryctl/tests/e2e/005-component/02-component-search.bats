#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures/component-search" 
  run $MESHERYCTL_BIN model import -f $FIXTURES_DIR/valid-model
  assert_success
}

@test "Search for a known component" {
  run $MESHERYCTL_BIN component search component_cli-e2e-test --page 1 --pagesize 1
  assert_success
  assert_output --partial "model_cli-e2e-test"
}

@test "given no search query when mesheryctl component search is run then an invalid argument error is displayed" {
  run $MESHERYCTL_BIN component search
  assert_failure

  assert_output --partial "Error: Only one argument must be provided"
  assert_output --partial "Invalid Argument"
  assert_output --partial "Usage: mesheryctl component search [query-text]"
  assert_output --partial "Run 'mesheryctl component search --help' to see detailed help message"
}

@test "given a non-existent component when mesheryctl component search component-name is run then no components are found" {
  run $MESHERYCTL_BIN component search nosuchcomponent123
  assert_success
  assert_output --partial "No components found"
}
