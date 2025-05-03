#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
    MESHERYCTL_DIR=$(dirname "$MESHERYCTL_BIN")
    export TEST_IMGS_OUTPUT="tests"
    export TEST_MODELS_OUTPUT="tests"
    export SHEET_ID="1kldsrMSwItkvlcazw8HEMmSrgsHPJsTsKY7f1hvJNAc"
    export FIXTURES_DIR="$MESHERYCTL_DIR/tests/e2e/006-registry/fixtures/registry-publish" 
    export CRED=$(base64 -w 0 $FIXTURES_DIR/cred.json)
    export INVALID_CRED="1234567890"
    export INVALID_SHEET="1kldsrMSwItkvlcazw8HEMmSrgsHPJsTsKY7f1hvJNA"
}

@test "mesheryctl registry publish displays usage instructions when no arguments provided" {
  run $MESHERYCTL_BIN registry publish
  assert_failure

  assert_output --partial "Error: [ system, google sheet credential, sheet-id, models output path, imgs output path] are required"
  assert_output --partial "Usage:"
  assert_output --partial "mesheryctl registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]"
  assert_output --partial "mesheryctl registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path] -o [output-format]"
  assert_output --partial "Run 'mesheryctl registry publish --help'"
  assert_output --partial "See https://docs.meshery.io/reference/mesheryctl/registry/publish for usage details"
}

@test "mesheryctl registry publish fails with invalid system" {
  run $MESHERYCTL_BIN registry publish invalid-system $CRED $SHEET_ID $TEST_MODELS_OUTPUT $TEST_IMGS_OUTPUT

  assert_output --partial "invalid system: invalid-system"
}

@test "mesheryctl registry publish fails with invalid google credential" {
  # This will fail because we're passing an invalid credential
  run $MESHERYCTL_BIN registry publish remote-provider $INVALID_CRED $SHEET_ID "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT"
  assert_failure

  assert_output --partial "Error: Invalid JWT Token: Ensure the provided token is a base64-encoded, valid Google Spreadsheets API token."
  assert_output --partial "See https://docs.meshery.io/reference/mesheryctl/registry/publish for usage details"
}

@test "mesheryctl registry publish fails with invalid sheet ID" {
  # This will fail because we're passing an invalid sheet ID
  run $MESHERYCTL_BIN registry publish remote-provider $CRED $INVALID_SHEET "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT"
  assert_failure

  assert_output --partial "Request to Google Spreadsheet did not succeed"
  assert_output --partial "Returned error: googleapi: Error 404: Requested entity was not found., notFound"
  assert_output --partial "See https://docs.meshery.io/reference/mesheryctl/registry/publish for usage details"
}

@test "mesheryctl registry publish remote-provider succeeds" {
  run $MESHERYCTL_BIN registry publish remote-provider "$CRED" "$SHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT"
  assert_success

  assert [ -f "../tests/accurate/model.json" ]
}

@test "mesheryctl registry publish website succeeds with md format" {
  run $MESHERYCTL_BIN registry publish website "$CRED" "$SHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT" -o md
  assert_success

  assert [ -f "../tests/accurate.md" ]
}

  @test "mesheryctl registry publish website succeeds with js format" {
  run $MESHERYCTL_BIN registry publish website "$CRED" "$SHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT" -o js
  assert_success

  assert [ -f "tests/data.js" ]
  }

@test "mesheryctl registry publish fails with invalid output format" {
  run $MESHERYCTL_BIN registry publish website "$CRED" "$SHEET_ID" "$TEST_MODELS_OUTPUT" "$TEST_IMGS_OUTPUT" -o invalid
  assert_failure

  assert_output --partial "Error: invalid output format: invalid"
  assert_output --partial "See https://docs.meshery.io/reference/mesheryctl/registry/publish for usage details"
}

  teardown(){
  rm -rf ../tests
  rm -rf tests
  }