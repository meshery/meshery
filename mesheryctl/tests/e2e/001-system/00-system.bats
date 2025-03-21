#!/usr/bin/env bats

@test "mesheryctl system help is succeeded and display help" {
    run $MESHERYCTL_BIN system --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Use \"mesheryctl system [command] --help\" for more information about a command." ]]
}

