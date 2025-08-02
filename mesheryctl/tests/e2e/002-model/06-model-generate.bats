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
    assert_output --partial "Usage: mesheryctl model generate [ file | filePath | URL ]"
    assert_output --partial "Run 'mesheryctl model generate --help' to see detailed help message"
}

@test "mesheryctl model generate fails with invalid URL" {
    run $MESHERYCTL_BIN model generate --file "invalid-url"
    assert_failure
    assert_output --partial "invalid URL"
}

@test "mesheryctl model generate succeeds with valid URL and template" {
    run $MESHERYCTL_BIN model generate --file "$FIXTURES_DIR/valid-model" --template "$FIXTURES_DIR/valid-template.json"
    assert_success
    assert_output --partial "Model can be accessed from $TESTDATA_DIR"
}

@test "mesheryctl model generate fails with missing template for URL" {
    run $MESHERYCTL_BIN model generate --file "$FIXTURES_DIR/valid-model"
    assert_failure
    assert_output --partial "Template file is not present"
}

@test "mesheryctl model generate succeeds with valid CSV directory" {
    run $MESHERYCTL_BIN model generate --file "$FIXTURES_DIR/valid-csv-dir"
    assert_success
    assert_output --partial "Model can be accessed from $TESTDATA_DIR"
    assert_output --partial "Logs for the csv generation can be accessed $TESTDATA_DIR/logs"
}

@test "mesheryctl model generate fails with invalid CSV directory" {
    run $MESHERYCTL_BIN model generate --file "invalid-dir"
    assert_failure
    assert_output --partial "error reading file"
}

@test "mesheryctl model generate skips registration with --register flag" {
    run $MESHERYCTL_BIN model generate --file "$FIXTURES_DIR/valid-csv-dir" --register 
    assert_success
    assert_output --partial "Model can be accessed from $FIXTURES_DIR"

}