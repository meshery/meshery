#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/model"
    mkdir -p "$TESTDATA_DIR"  
    export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
}

@test "mesheryctl model generate displays usage instructions when no file, filePath or URL is provided" {
    run $MESHERYCTL_BIN model generate
    assert_failure
    assert_output --partial "[ file | filepath | URL ] isn't specified"
    assert_output --partial "Usage: mesheryctl model generate [ file | filePath | URL ]"
    assert_output --partial "Run 'mesheryctl model generate --help' to see detailed help message"
}

@test "mesheryctl model generate --help displays help information" {
    run $MESHERYCTL_BIN model generate --help
    assert_success
    assert_output --partial "Generate models by specifying the directory, file, or URL"
    assert_output --partial "Usage:"
    assert_output --partial "Flags:"
    assert_output --partial "-f, --file"
    assert_output --partial "-t, --template"
    assert_output --partial "-r, --register"
}

@test "mesheryctl model generate fails with missing template for URL" {
    run $MESHERYCTL_BIN model generate -f "https://example.com/model"
    assert_failure
    assert_output --partial "no template file is provided while using url for importing a model"
}

@test "mesheryctl model generate fails with invalid URL and template" {
    run $MESHERYCTL_BIN model generate -f "https://example.com/model" -t "nonexistent-template.json"
    assert_failure
    assert_output --partial "open nonexistent-template.json: no such file or directory"
}

@test "mesheryctl model generate can successfully process URL with template" {
    run $MESHERYCTL_BIN model generate -f "https://github.com/Azure/azure-service-operator/releases/download/v2.14.0/azureserviceoperator_customresourcedefinitions_v2.14.0.yaml" -t "$FIXTURES_DIR/valid-template.json" --register
    # NOTE: This test only validates successful URL processing without server registration
    assert_output --partial "Generating model from URL: https://github.com/Azure/azure-service-operator/releases/download/v2.14.0/azureserviceoperator_customresourcedefinitions_v2.14.0.yaml"
}

@test "mesheryctl model generate fails with invalid CSV directory" {
    run $MESHERYCTL_BIN model generate -f "invalid-dir"
    assert_failure
    assert_output --partial "open invalid-dir: no such file or directory"
}

@test "mesheryctl model generate fails with non-existent file" {
    run $MESHERYCTL_BIN model generate -f "non-existent-file.csv"
    assert_failure
    assert_output --partial "open non-existent-file.csv: no such file or directory"
}

@test "mesheryctl model generate can successfully process CSV directory" {
    run $MESHERYCTL_BIN model generate -f "$FIXTURES_DIR/real-valid-csv" --register
    # NOTE: This test only validates successful CSV processing without server registration
    assert_output --partial "Generating model from CSV files"
}