#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries
}

# tests config file and authenticatin file set properly

@test "meshery config.yaml file as been created" {
    assert_exists "$MESHERY_CONFIG_FILE_PATH"
}

@test "meshery config.yaml provider is Meshery" {
    run yq '.contexts.local.provider' "$MESHERY_CONFIG_FILE_PATH"
    assert_success

    assert_output  --partial "Layer5"
}

@test "mesehry auth.json file as been created" {
    assert_exists "$MESHERY_AUTH_FILE"
}

@test "meshery auth.json file meshery provider is Meshery" {
    run jq '."meshery-provider"' "$MESHERY_AUTH_FILE"
    assert_success

    assert_output "\"Layer5\""
}