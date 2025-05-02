#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries
}

@test "Search for a known component (e.g., 'prometheus')" {
  run mesheryctl component search prometheus
  [ "$status" -eq 0 ]
  [[ "$output" == *"prometheus"* ]]
}

# Test 2: Error when no query is provided
@test "Error when no search query is provided" {
  run mesheryctl component search
  [ "$status" -ne 0 ]
  [[ "$output" == *"Please enter component name to search"* ]]
}

# Test 3: Search for a random non-existent component
@test "Search for non-existent component" {
  run mesheryctl component search nosuchcomponent123
  [ "$status" -eq 0 ]
  [[ "$output" == *"No components found"* || "$output" == *"0 components"* || "$output" == *"[]"* ]]
}
