#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/tests_helpers"
	_tests_helper

}

@test "meshery model list --page 1 succed and return total numbers of models" {
  run $MESHERYCTL_BIN model list --page 1
  assert_success

  assert_line --regexp "^Total number of models: [0-9]+$"
}

@test "mesheryctl model list --count is succeeded and return total numbers of models" {
  run $MESHERYCTL_BIN model list --count
  assert_success

  assert_output --regexp "^Total number of models: [0-9]+$"
}
