#!/usr/bin/env bats

setup() {
    E2ETEST_DIR=$(dirname "$BATS_TEST_DIRNAME")
    load "$E2ETEST_DIR/$SUPPORT_DESTDIR"
    load "$E2ETEST_DIR/$ASSERT_DESTDIR"

    MESHERYCTL_DIR=$(dirname "$MESHERYCTL_BIN")
    export TESTDATA_DIR="$MESHERYCTL_DIR/tests/e2e/002-model/testdata/model-import" 
}

@test "mesheryctl model import displays usage instructions" {
  run $MESHERYCTL_BIN model import
  [ "$status" -ne 0 ]

  assert_output --partial "Error: [ file | filepath | URL ] isn't specified" 
  assert_output --partial "Usage: mesheryctl model import [ file | filePath | URL ]"
}

@test "mesheryctl model import for empty model returns no components found" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/empty-model
  echo "$TESTDATA_DIR/empty-model" >&2
  [ "$status" -eq 0 ]

  assert_output --partial "No component found in model provided. Models must have at least one component."
}

@test "mesheryctl model import fails for invalid files" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/invalid-model
  [ "$status" -eq 0 ]

  assert_output --partial "The file 'model.json' could not be parsed due to invalid JSON syntax. Error details: The \`schemaVersion\` key is either empty or has an incorrect value."
  assert_output --partial "The file 'test.svg' could not be used for the operation 'import' because the extension '.svg' is not supported."

}

@test "mesheryctl model import fails for zip format" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/test-model.zip
  [ "$status" -eq 0 ]

  assert_output --partial "The archive may contain unsupported compression formats or features."
}

@test "mesheryctl model import fails for invalid model artifact" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/non-oci-model-archive.tar
  [ "$status" -eq 0 ]

  assert_output --partial "An error occurred while attempting to extract the TAR archive"
  assert_output --partial "The archive may be non OCI compliant."
}

@test "mesheryctl model import succeeds for valid model as OCI artifact" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/oci-model-archive.tar
  [ "$status" -eq 0 ]

  assert_output --partial "Imported model kccna24 (1 component)"
}

@test "mesheryctl model import succeeds for valid model in uncompressed format" {
  run $MESHERYCTL_BIN model import -f $TESTDATA_DIR/valid-model
  [ "$status" -eq 0 ]

  assert_output --partial "Imported model kccna24 (1 component)"
}

@test "mesheryctl model import succeeds for URL of valid model as OCI artifact" {
  # TODO: replace URL
  # run $MESHERYCTL_BIN model import -f https://github.com/meshery/meshery/raw/refs/heads/master/mesheryctl/tests/e2e/002-model/testdata/model-import/oci-model-archive.tar
  run $MESHERYCTL_BIN model import -f https://github.com/riyaa14/meshery-testing/raw/refs/heads/main/valid-models%20copy/kccna24-test-oci.tar
  [ "$status" -eq 0 ]

  assert_output --partial "Imported model kccna24 (1 component)"
}

teardown() {
  # remove models and assiciated components from registry after testing
  sqlite3 $HOME/.meshery/config/mesherydb.sql "DELETE FROM component_definition_dbs WHERE model_id IN (SELECT id FROM model_dbs WHERE name='kccna24');"
  sqlite3 $HOME/.meshery/config/mesherydb.sql "DELETE FROM model_dbs WHERE name='kccna24';"
}



