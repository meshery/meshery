#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/tests_helpers"
	_tests_helper

}


@test "mesheryctl model --count return total numbers of models" {
  run $MESHERYCTL_BIN model --count
  assert_success

  assert_output --regexp "^Total number of models: [0-9]+$"
}
