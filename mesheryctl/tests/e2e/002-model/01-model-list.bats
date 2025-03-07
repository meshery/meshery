#!/usr/bin/env bats

@test "meshery model list --page 1 succed and return total numbers of models" {
  run $MESHERYCTL_BIN model list --page 1
  [ "$status" -eq 0 ]

  run grep -E "^Total number of  models : [0-9]+$" <<< $(echo "$output" | head -n 1)
  [ "$status" -eq 0 ]
}

@test "mesheryctl model list --count is succeeded and return total numbers of models" {
  run $MESHERYCTL_BIN model list --count
  [ "$status" -eq 0 ]

  run grep -E "^Total number of  models : [0-9]+$" <<< $(echo "$output" | head -n 1)
  [ "$status" -eq 0 ]
}
