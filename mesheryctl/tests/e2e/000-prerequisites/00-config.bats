#!/usr/bin/env bats

# tests config file and authenticatin file set properly

@test "meshery config.yaml file as been created" {
    [ -f "$MESHERY_CONFIG_FILE_PATH" ]
}

@test "meshery config.yaml provider is Meshery" {
    run yq '.contexts.local.provider' "$MESHERY_CONFIG_FILE_PATH"
    [ "$status" -eq 0 ]
    [ "$output" = "Meshery" ]       
}

@test "mesehry auth.json file as been created" {
    [ -f "$MESHERY_AUTH_FILE" ]
}

@test "meshery auth.json file meshery provider is Meshery" {
    run jq '."meshery-provider"' "$MESHERY_AUTH_FILE"
    [ "$status" -eq 0 ]
    [ "$output" = "\"Meshery\"" ]       
}