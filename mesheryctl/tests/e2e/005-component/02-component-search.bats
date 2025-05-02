#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures/component-search" 
  run $MESHERYCTL_BIN model import -f $FIXTURES_DIR/valid-model
}

@test "Search for a known component" {
  run $MESHERYCTL_BIN component search component_cli-e2e-test
  assert_success
  assert_output --partial "model_cli-e2e-test"
}

@test "Error when no search query is provided" {
  run $MESHERYCTL_BIN component search
  assert_failure

  assert_output --partial "Error: [search term] isn't specified. Please enter component name to search"
  assert_output --partial "Usage: mesheryctl exp component search [query-text]"
  assert_output --partial "Run 'mesheryctl exp component search --help' to see detailed help message"
}

@test "Search for non-existent component" {
  run $MESHERYCTL_BIN component search nosuchcomponent123
  assert_success
  assert_output --partial "No components found"
}
