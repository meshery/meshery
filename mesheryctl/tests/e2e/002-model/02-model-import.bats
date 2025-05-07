#!/usr/bin/env bats

setup() {
   load "$E2E_HELPERS_PATH/bats_libraries"
	_load_bats_libraries

  export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures/model-import" 
}

@test "mesheryctl model import displays usage instructions" {
  run $MESHERYCTL_BIN model import
  assert_failure

  assert_output --partial "Error: [ file | filepath | URL ] isn't specified" 
  assert_output --partial "Usage: mesheryctl model import [ file | filePath | URL ]"
}

@test "mesheryctl model import for empty model returns no components found" {
  run $MESHERYCTL_BIN model import -f $FIXTURES_DIR/empty-model
  assert_failure

  assert_output --partial "No component found in model provided. Models must have at least one component."
}

@test "mesheryctl model import fails for invalid files" {
  run $MESHERYCTL_BIN model import -f $FIXTURES_DIR/invalid-model
  assert_failure

  assert_output --partial "The file 'model.json' could not be parsed due to invalid JSON syntax. Error details: The \`schemaVersion\` key is either empty or has an incorrect value."
  assert_output --partial "The file 'test.svg' could not be used for the operation 'import' because the extension '.svg' is not supported."

}

@test "mesheryctl model import fails for zip format" {
  run $MESHERYCTL_BIN model import -f $FIXTURES_DIR/test-model.zip
  assert_failure

  assert_output --partial "The archive may contain unsupported compression formats or features."
}

@test "mesheryctl model import fails for invalid model artifact" {
  run $MESHERYCTL_BIN model import -f $FIXTURES_DIR/non-oci-model-archive.tar
  assert_failure

  assert_output --partial "An error occurred while attempting to extract the TAR archive"
  assert_output --partial "The archive may be non OCI compliant."
}

@test "mesheryctl model import succeeds for valid model as OCI artifact" {
  run $MESHERYCTL_BIN model import -f $FIXTURES_DIR/valid-oci-model-archive.tar
  assert_success

  assert_output --partial "Imported model model-import_cli-e2e-test (1 component)"
  assert_output --partial "model-import-component_cli-e2e-test"
}

@test "mesheryctl model import succeeds for valid model in uncompressed format" {
  run $MESHERYCTL_BIN model import -f $FIXTURES_DIR/valid-model
  assert_success

  assert_output --partial "Imported model model-import_cli-e2e-test (1 component)"
  assert_output --partial "model-import-component_cli-e2e-test"
}

@test "mesheryctl model import succeeds for URL of valid model as OCI artifact" {
  skip
  run $MESHERYCTL_BIN model import -f https://github.com/meshery/meshery/blob/master/mesheryctl/tests/e2e/002-model/fixtures/model-import/valid-oci-model-archive.tar
  assert_success

  assert_output --partial "Imported model model-import_cli-e2e-test (1 component)"
  assert_output --partial "model-import-component_cli-e2e-test"
}



