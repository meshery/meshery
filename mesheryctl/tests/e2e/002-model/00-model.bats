#!/usr/bin/env bats

@test "mesheryctl model --count is succeeded and return total numbers of models" {
  run $MESHERYCTL_BIN model --count
  [ "$status" -eq 0 ]

  run grep -E "^Total number of  models : [0-9]+$" <<< $(echo "$output" | head -n 1)
  [ "$status" -eq 0 ]
}
