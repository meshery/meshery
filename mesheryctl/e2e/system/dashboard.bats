#!/usr/bin/env bats

# Test: mesheryctl system dashboard when Meshery Server is NOT running
# This checks for:
#   - Non-zero exit status
#   - Expected kubeconfig missing/server not running message
#   - Prints full output for debugging

@test "mesheryctl system dashboard shows error when kubeconfig is missing" {

  # Run the command
  run mesheryctl/bin/mesheryctl system dashboard

  # Print full output for debugging (always printed)
  echo "===== FULL OUTPUT BELOW ====="
  echo "$output"
  echo "===== END OUTPUT ====="

  # Exit code must be non-zero (failure expected)
  [ "$status" -ne 0 ]

  # Match the real error message you saw:
  #   Error: open /home/user/.kube/config: no such file or directory
  [[ "$output" =~ "no such file or directory" ]]
}
