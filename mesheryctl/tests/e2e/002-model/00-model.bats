#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  load "$E2E_HELPERS_PATH/constants"
}


@test "mesheryctl model --count return total numbers of models" {
  run $MESHERYCTL_BIN model --count
  assert_success

  assert_output --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}
