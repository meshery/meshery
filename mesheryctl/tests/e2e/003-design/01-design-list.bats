#!/usr/bin/env bats

@test "mesheryctl design list shows designs" {
  # Run the list command
  run $MESHERYCTL_BIN design list
  echo "List command output: $output"
  [ "$status" -eq 0 ]
  
  # Verify output format - should contain headers for designs
  echo "$output" | grep -q "DESIGN ID"
  [ "$?" -eq 0 ]
}

@test "mesheryctl design list shows total count" {
  # Run list to check for total count display
  run $MESHERYCTL_BIN design list
  echo "List command output: $output"
  [ "$status" -eq 0 ]
  
  # Check for the total count line
  echo "$output" | grep -q "Total number of patterns:"
  [ "$?" -eq 0 ]
}

@test "mesheryctl design list with page parameter works" {
  # Test pagination with page flag
  run $MESHERYCTL_BIN design list --page 1
  echo "List command output with page=1: $output"
  [ "$status" -eq 0 ]
  
 
  echo "$output" | grep -q "Total number of  patterns :"
  [ "$?" -eq 0 ]
  
  # Verify the table header exists
  echo "$output" | grep -q "DESIGN ID"
  [ "$?" -eq 0 ]
  
  # Check that data rows appear
  echo "$output" | grep -E "[0-9a-f]{8}"
  [ "$?" -eq 0 ]
}

@test "mesheryctl design list with verbose flag works" {
  # Test verbose output
  run $MESHERYCTL_BIN design list -v
  echo "Verbose list command output: $output"
  [ "$status" -eq 0 ]
  
  # Should either show designs or indicate no patterns found
  if echo "$output" | grep -q "No pattern(s) found"; then
    [ "$?" -eq 0 ]
  else
    echo "$output" | grep -q "DESIGN ID"
    [ "$?" -eq 0 ]
  fi
}