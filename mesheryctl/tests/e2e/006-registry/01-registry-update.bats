# bats file: tests/e2e/006-registry/03-registry-update.bats

#!/usr/bin/env bats

setup() {
	load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

	load "$E2E_HELPERS_PATH/constants"
	export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures/registry-update"
	export TEMP_DATA_DIR="$BATS_TEST_DIRNAME/fixtures/registry-update/temp"
	export LOG_PATH="$HOME/.meshery/logs/registry"
	# Dummy base64 credential and spreadsheet id for testing
	export INVALID_CRED="invalid-base64"
	export VALID_SHEET_ID="1kldsrMSwItkvlcazw8HEMmSrgsHPJsTsKY7f1hvJNAc"
	export INVALID_SHEET_ID="invalid-sheet-id"
	export VALID_CRED=$(cat $FIXTURES_DIR/cred.json | base64)
}

teardown() {
    # Remove any custom models directories created during tests
    rm -rf "$TEMP_DATA_DIR"
    # Remove log files created by the tests
    rm -f "$LOG_PATH/registry-update"
}

# It does not currently show usage instructions when no flags are provided.
@test "mesheryctl registry update with no arguments displays usage instructions with no flags" {
	run $MESHERYCTL_BIN registry update
	assert_failure
	assert_output --partial "unexpected end of JSON input"
	assert_output --partial "Error: unexpected end of JSON input"
}

@test "mesheryctl registry update with only spreadsheet-id fails" {
	run $MESHERYCTL_BIN registry update --spreadsheet-id "$SHEET_ID"
	assert_failure
	assert_output --partial "Error: if any flags in the group [spreadsheet-id spreadsheet-cred] are set they must all be set; missing [spreadsheet-cred]"
}

@test "mesheryctl registry update with only spreadsheet-cred fails" {
	run $MESHERYCTL_BIN registry update --spreadsheet-cred "$VALID_CRED"
	assert_failure
	assert_output --partial "Error: if any flags in the group [spreadsheet-id spreadsheet-cred] are set they must all be set; missing [spreadsheet-id]"
}

@test "mesheryctl registry update with invalid spreadsheet-id fails" {
	run $MESHERYCTL_BIN registry update --spreadsheet-id "INVALID_SHEET" --spreadsheet-cred "$VALID_CRED"
	assert_failure
	assert_output --partial "googleapi: Error 404: Requested entity was not found., notFound"
    assert_output --partial "Error: googleapi: Error 404: Requested entity was not found., notFound"
}

@test "mesheryctl registry update with invalid spreadsheet-cred fails" {
	run $MESHERYCTL_BIN registry update --spreadsheet-id "$SHEET_ID" --spreadsheet-cred "$INVALID_CRED"
	assert_failure
	assert_output --partial "invalid character '\u008a' looking for beginning of value"
	assert_output --partial "Error: invalid character '\u008a' looking for beginning of value"
}

@test "mesheryctl registry update with valid spreadsheet-id and cred succeeds" {
	run $MESHERYCTL_BIN registry update --spreadsheet-id $SHEET_ID --spreadsheet-cred $VALID_CRED
	assert_success
	assert_output --partial "Downloaded CSV from:" || assert_output --partial "Downloading CSV from:"
	assert_output --partial "Parsing Components..."
	assert_output --regexp "Updated [0-9]+ models and [0-9]+ components"
	assert_output --partial "refer $LOG_PATH for detailed registry update logs"
}

@test "mesheryctl registry update with --model flag for a specific model succeeds" {
	run $MESHERYCTL_BIN registry update --spreadsheet-id "$SHEET_ID" --spreadsheet-cred "$VALID_CRED" --model "test-model"
	assert_success
	assert_output --partial "Downloaded CSV from:" || assert_output --partial "Downloading CSV from:"
	assert_output --partial "Parsing Components..."
	assert_output --regexp "Updated [0-1]+ models and [0-9]+ components"
	assert_output --partial "refer $LOG_PATH for detailed registry update logs"
}

@test "mesheryctl registry update with invalid model name" {
	run $MESHERYCTL_BIN registry update --spreadsheet-id "$SHEET_ID" --spreadsheet-cred "$VALID_CRED" --model "nonexistent-model"
	assert_success
	assert_output --partial "Downloaded CSV from:" || assert_output --partial "Downloading CSV from:"
	assert_output --partial "Parsing Components..."
	assert_output --regexp "Updated 0 models and 0 components"
	assert_output --partial "refer $LOG_PATH for detailed registry update logs"
}

@test "mesheryctl registry update with --input flag uses custom models directory" {
    mkdir -p "$TEMP_DATA_DIR/custom-models"

    run $MESHERYCTL_BIN registry update --spreadsheet-id "$SHEET_ID" --spreadsheet-cred "$VALID_CRED" --input "$TEMP_DATA_DIR/custom-models"
	assert_success
    assert_output --partial "Downloaded CSV from:" || assert_output --partial "Downloading CSV from:"
	assert_output --partial "Parsing Components..."
	assert_output --regexp "Updated [0-9]+ models and [0-9]+ components"
	assert_output --partial "refer $LOG_PATH for detailed registry update logs"
}

@test "mesheryctl registry update creates log file in logs/registry" {
	# Ensure no previous log exists
	rm -f "$LOG_PATH/registry-update"  
	run $MESHERYCTL_BIN registry update --spreadsheet-id "$SHEET_ID" --spreadsheet-cred "$VALID_CRED"
	assert_file_exist $LOG_PATH/registry-update
}