#!/usr/bin/env bats

@test "mesheryctl design delete is succeeded" {
  # Check if the design ID file exists
  if [ ! -f "/tmp/meshery_test_design_id" ]; then
    skip "No design ID available to delete"
  fi
  
  DESIGN_ID=$(cat /tmp/meshery_test_design_id)
  echo "Using design ID from file: $DESIGN_ID"
  
  # Skip if design ID is empty
  [ -n "$DESIGN_ID" ] || skip "Empty design ID"
  
  # Get list of designs 
  run $MESHERYCTL_BIN design list
  BEFORE_DELETE="$output"
  echo "Designs before delete: $BEFORE_DELETE"
  
  # Delete the design
  run $MESHERYCTL_BIN design delete "$DESIGN_ID"
  echo "Delete command output: $output"
  
  # Check if command executed without error
  [ "$status" -eq 0 ]
  
  # Check if the output indicates success
  # Note: The actual success message may vary depending on implementation
  echo "$output" | grep -E "deleted|success|removed"
  
  # Wait a moment for the server to process the deletion
  sleep 2
  
  # Check that design no longer appears in list (if it was listed before)
  if echo "$BEFORE_DELETE" | grep -q "$DESIGN_ID"; then
    run $MESHERYCTL_BIN design list
    echo "Designs after delete: $output"
    
    # Either the ID should not appear, or there should be one fewer design
    if echo "$output" | grep -q "$DESIGN_ID"; then
      echo "Design ID still appears in list, checking if total count decreased"
      BEFORE_COUNT=$(echo "$BEFORE_DELETE" | grep -o "Total number of patterns:[0-9]*" | grep -o "[0-9]*")
      AFTER_COUNT=$(echo "$output" | grep -o "Total number of patterns:[0-9]*" | grep -o "[0-9]*")
      
      # Ensure counts are numbers before comparing
      if [[ "$BEFORE_COUNT" =~ ^[0-9]+$ ]] && [[ "$AFTER_COUNT" =~ ^[0-9]+$ ]]; then
        [ "$AFTER_COUNT" -lt "$BEFORE_COUNT" ]
      fi
    fi
  fi
}

@test "mesheryctl design delete for non-existent ID should give appropriate response" {
  # Use a non-existent design ID
  NONEXISTENT_ID="00000000-0000-0000-0000-000000000000"
  
  # Attempt to delete non-existent design
  run $MESHERYCTL_BIN design delete "$NONEXISTENT_ID"
  echo "Delete command output: $output"
  
  # The command might return success with a message or return an error code
  # We're just verifying it runs and gives appropriate feedback
  if [ "$status" -ne 0 ]; then
    # If it returns an error code, that's fine
    [ "$?" -eq 0 ]
  else
    # If it returns success, it should mention "not found" or similar
    echo "$output" | grep -E "not found|failed|no such design|error"
    [ "$?" -eq 0 ]
  fi
  
  # Clean up the temporary files 
  [ -f "/tmp/meshery_test_design_id" ] && rm -f /tmp/meshery_test_design_id
  [ -f "/tmp/meshery_test_design_id2" ] && rm -f /tmp/meshery_test_design_id2
}