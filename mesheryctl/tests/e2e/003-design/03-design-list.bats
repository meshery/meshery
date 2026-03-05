#!/usr/bin/env bats

# Setup function to load libraries
setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  load "$E2E_HELPERS_PATH/constants"

  export DESIGN_LIST_OUTPUT__HEADER_REGEX_PATTERN="^DESIGN[[:space:]]ID[[:space:]]+(USER[[:space:]]ID)?[[:space:]]+NAME[[:space:]]+CREATED[[:space:]]+UPDATED[[:space:]]+$"
}

@test "given all requirements met when running mesheryctl design list --page 1 then the total numbers of designs and a list of designs are displayed" {
  run $MESHERYCTL_BIN design list --page 1
  assert_success

  assert_line --regexp "$DESIGN_LIST_OUTPUT__HEADER_REGEX_PATTERN"
  assert_line --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}

@test "given all requirements met when running mesheryctl design list --count then the total numbers of designs is displayed" {
  run $MESHERYCTL_BIN design list --count
  assert_success
  
  assert_line --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}