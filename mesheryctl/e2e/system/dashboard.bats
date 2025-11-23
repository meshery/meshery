#!/usr/bin/env bats

# Test: mesheryctl system dashboard command
# This checks for:
#   - Non-zero exit status
#   - Expected kubeconfig missing/server not reachable message

@test "mesheryctl system dashboard shows error when kubeconfig is missing" {
  run mesheryctl/bin/mesheryctl system dashboard

  # Exit code must be non-zero (failure expected)
  [ "$status" -ne 0 ]

  # Expected kubeconfig missing message
  [[ "$output" =~ "no such file or directory" ]]
}
