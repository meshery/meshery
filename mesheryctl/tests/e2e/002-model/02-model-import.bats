#!/usr/bin/env bats

function setup() {
    MESHERYCTL_DIR=$(dirname "$MESHERYCTL_BIN")
    export TESTDATA_DIR="$MESHERYCTL_DIR/tests/e2e/002-model/testdata/model-import" 
}

check_output(){
  local expected_message="$1"
  
  run grep -Fq "$expected_message" <<< "$actual_output"
  if [ "$status" -ne 0 ]; then
    echo "FAILED: Expected message not found: \"$expected_message\"" >&2
    echo "Actual Output: " >&2
    echo "$actual_output" >&2
  fi
  
  [ "$status" -eq 0 ]
}

@test "mesheryctl model import displays usage instructions" {
  run $MESHERYCTL_BIN model import
  [ "$status" -ne 0 ]

  actual_output=$(echo "$output")

  check_output "Error: [ file | filepath | URL ] isn't specified" 
  check_output "Usage: mesheryctl model import [ file | filePath | URL ]"
}

@test "mesheryctl model import for empty model returns no components found" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/empty-model
  echo "$TESTDATA_DIR/empty-model" >&2
  [ "$status" -eq 0 ]

  actual_output=$(echo "$output")

  check_output "No component found in model provided. Models must have at least one component."
}

@test "mesheryctl model import fails for invalid files" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/invalid-model
  [ "$status" -eq 0 ]

  actual_output=$(echo "$output")

  check_output "The file 'model.json' could not be parsed due to invalid JSON syntax. Error details: The \`schemaVersion\` key is either empty or has an incorrect value."
  check_output "The file 'test.DS_Store' could not be used for the operation 'import' because the extension '.DS_Store' is not supported."
  check_output "The file 'test.svg' could not be used for the operation 'import' because the extension '.svg' is not supported."

}

@test "mesheryctl model import fails for zip format" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/test-model.zip
  [ "$status" -eq 0 ]

  actual_output=$(echo "$output")

  check_output "The archive may contain unsupported compression formats or features."
}

@test "mesheryctl model import fails for invalid model artifact" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/non-oci-model-archive.tar
  [ "$status" -eq 0 ]

  actual_output=$(echo "$output")

  check_output "An error occurred while attempting to extract the TAR archive"
  check_output "The archive may be non OCI compliant."
}

@test "mesheryctl model import succeeds for valid model as OCI artifact" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/oci-model-archive.tar
  [ "$status" -eq 0 ]

  actual_output=$(echo "$output")

  check_output "Imported model kccna24 (1 component)"
}

@test "mesheryctl model import succeeds for valid model in uncompressed format" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/valid-model
  [ "$status" -eq 0 ]

  actual_output=$(echo "$output")

  check_output "Imported model kccna24 (1 component)"
}

@test "mesheryctl model import succeeds for URL of valid model as OCI artifact" {
  # TODO: replace URL
  # run $MESHERYCTL_BIN model import -f https://github.com/meshery/meshery/raw/refs/heads/master/mesheryctl/tests/e2e/002-model/testdata/model-import/oci-model-archive.tar
  run $MESHERYCTL_BIN model import -f https://github.com/riyaa14/meshery-testing/raw/refs/heads/main/valid-models%20copy/kccna24-test-oci.tar
  [ "$status" -eq 0 ]

  actual_output=$(echo "$output")

  check_output "Imported model kccna24 (1 component)"
}

function teardown() {
  # remove models and assiciated components from registry after testing
  sqlite3 $HOME/.meshery/config/mesherydb.sql "DELETE FROM component_definition_dbs WHERE model_id IN (SELECT id FROM model_dbs WHERE name='kccna24');"
  sqlite3 $HOME/.meshery/config/mesherydb.sql "DELETE FROM model_dbs WHERE name='kccna24';"
}



