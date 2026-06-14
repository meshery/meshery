#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  load "$E2E_HELPERS_PATH/constants"

}

@test "mesheryctl model list --page 1 return total numbers of models" {
  run $MESHERYCTL_BIN model list --page 1
  assert_success

  assert_line --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}

@test "mesheryctl model list --count return total numbers of models" {
  run $MESHERYCTL_BIN model list --count
  assert_success

  assert_output --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}
