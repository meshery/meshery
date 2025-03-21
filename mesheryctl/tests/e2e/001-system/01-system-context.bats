#!/usr/bin/env bats

@test "mesheryctl system context view --all is succeeded" {
   run $MESHERYCTL_BIN system context view --all
   [ "$status" -eq 0 ]
   [[ "$output" =~ "endpoint" ]]
   [[ "$output" =~ "token" ]]
   [[ "$output" =~ "platform" ]]
   [[ "$output" =~ "provider" ]]
}

@test "mesheryctl system context view is succedeed" {
   run $MESHERYCTL_BIN system context view
   [ "$status" -eq 0 ]
   [[ "$output" =~ "Current Context" ]]
   [[ "$output" =~ "endpoint" ]]
   [[ "$output" =~ "token" ]]
   [[ "$output" =~ "platform" ]]
   [[ "$output" =~ "provider" ]]
}
