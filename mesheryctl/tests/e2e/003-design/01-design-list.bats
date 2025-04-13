#!/usr/bin/env bats

# Setup function to load libraries
setup() {
  load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  load "$E2E_HELPERS_PATH/constants"
}

@test "mesheryctl design list returns total numbers of designs" {
  run $MESHERYCTL_BIN design list
  assert_success
  assert_line --partial "DESIGN ID" && \
    assert_line --partial "USER ID" && \
    assert_line --partial "NAME" && \
    assert_line --partial "CREATED" && \
    assert_line --partial "UPDATED"
  assert_line --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}

@test "mesheryctl design list --page 1 returns total numbers of designs" {
  run $MESHERYCTL_BIN design list --page 1
  assert_success
  # Fix: Match the exact output format with two spaces between "of" and "patterns"
  assert_line --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}

@test "mesheryctl design list -v returns total numbers of designs" {
  run $MESHERYCTL_BIN design list -v
  assert_success
  # Check for either designs or "No pattern(s) found" message
  assert_line --regexp "(No pattern(s) found|.*DESIGN ID.*)" || assert_line --regexp "$LIST_COMMAND_OUTPUT_REGEX_PATTERN"
}