#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
   _load_bats_libraries

   load "$E2E_HELPERS_PATH/constants"

}

@test "Help for components list command" {
  run $MESHERYCTL_BIN component list --help
  assert_success
  assert_output --partial "List all components registered in Meshery Server"
}

@test "Count components" {
  run $MESHERYCTL_BIN component list --count
  assert_success
  clean_output="$(echo "$output" | sed -r 's/\x1B\[[0-9;]*[A-Za-z]//g')"
  output="$clean_output"
  assert_output --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}

@test "component list shows header and pagination info" {
  run $MESHERYCTL_BIN component list --page 1
  assert_line --index 0 --partial "Total number of components:"
  assert_line --index 1 --partial "Page:"
  assert_line --index 2 --partial "MODEL"
  assert_line --index 3 --regexp '^[[:space:]]*[a-zA-Z0-9]'
}
