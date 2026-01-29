#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/model"
    mkdir -p "$TESTDATA_DIR"  
    export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
}

@test "given no file, filePath or URL is provided when running mesheryctl model generate then usage instructions are displayed" {
    run $MESHERYCTL_BIN model generate 
    assert_failure
    assert_output --partial "Usage: mesheryctl model generate [ file | filePath | URL ]"
    assert_output --partial "Run 'mesheryctl model generate --help' to see detailed help message"
}

@test "given an invalid URL is provided when running mesheryctl model generate then an error message is displayed" {
    run $MESHERYCTL_BIN model generate --file "invalid-url"
    assert_failure
    assert_output --partial "no such file or directory"
}

@test "given a valid CSV directory is provided when running mesheryctl model generate --file valid-csv-directory then the model is generated" {
    run $MESHERYCTL_BIN model generate --file "$FIXTURES_DIR/valid-model" --template "$FIXTURES_DIR/valid-template.json"
    assert_success
    assert_output --partial "Model can be accessed from $TESTDATA_DIR"
}

@test "given a valid URL without a template when running mesheryctl model generate then an error is returned" {
    run $MESHERYCTL_BIN model generate --file "$FIXTURES_DIR/valid-model"
    assert_failure
    assert_output --partial "error reading csv file"
}

@test "given a valid CSV directory is provided when running mesheryctl model generate then the model is generated successfully" {
    run $MESHERYCTL_BIN model generate --file "$FIXTURES_DIR/valid-csv-dir"
    assert_success
    assert_output --partial "Model can be accessed from $TESTDATA_DIR"
    assert_output --partial "Logs for the csv generation can be accessed $TESTDATA_DIR/logs"
}

@test "given an invalid CSV directory is provided when running mesheryctl model generate then an error is returned" {
    run $MESHERYCTL_BIN model generate --file "invalid-dir"
    assert_failure
    assert_output --partial "no such file or directory"
}

@test "given the --register flag is provided when running mesheryctl model generate then the model is generated without registry registration" {
    run $MESHERYCTL_BIN model generate --file "$FIXTURES_DIR/valid-csv-dir" --register 
    assert_success
    assert_output --partial "Model can be accessed from $FIXTURES_DIR"

}
