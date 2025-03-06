#!/usr/bin/env bats
function setup() {
    MESHERY_CONFIG_FILE_PATH="$HOME/.meshery/config.yaml"
    MESHERY_AUTH_FILE="$HOME/.meshery/auth.json"
}


@test "mesheryctl system status is succeeded" {
    run $MESHERYCTL_BIN system status -y
    [ "$status" -eq 0 ]
}

@test "mesehry config file as been created" {
    [ -f "$MESHERY_CONFIG_FILE_PATH" ]
}

@test "meshery config provider is Meshery" {
    run yq '.contexts.local.provider' "$MESHERY_CONFIG_FILE_PATH"
    [ "$status" -eq 0 ]
    [ "$output" = "Meshery" ]       
}

@test "mesehry  auth.json file as been created" {
    [ -f "$MESHERY_AUTH_FILE" ]
}

@test "meshery auth provider is Meshery" {
    run jq '."meshery-provider"' "$MESHERY_AUTH_FILE"
    [ "$status" -eq 0 ]
    [ "$output" = "Meshery" ]       
}

#@test "mesheryctl system check is succeeded" {
#    run $MESHERYCTL_BIN system check
#    [ "$status" -eq 0 ]
#}