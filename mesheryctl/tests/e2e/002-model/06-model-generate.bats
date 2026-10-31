#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries
    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/model"
    mkdir -p "$TESTDATA_DIR"
    export FIXTURES_DIR="$BATS_TEST_DIRNAME/fixtures"
    export CSV_FIXTURES_DIR="$BATS_TEST_DIRNAME/../../../internal/cli/root/model/fixtures/templates/template-csvs"
}

@test "[cut=Model][tg=Model Generate] given no file, filePath or URL when running mesheryctl model generate then usage instructions are displayed" {
    run $MESHERYCTL_BIN model generate
    assert_failure
    assert_output --partial "either --file flag, a URL or a path as argument must be specified"
    assert_output --partial "Usage: mesheryctl model generate [ file | filePath | URL ]"
    assert_output --partial "Run 'mesheryctl model generate --help' to see detailed help message"
}

@test "[cut=Model][tg=Model Generate] given an invalid file path when running mesheryctl model generate then an error message is displayed" {
    run $MESHERYCTL_BIN model generate --file "invalid-url"
    assert_failure
    assert_output --partial "open invalid-url: no such file or directory"
}

@test "[cut=Model][tg=Model Generate] given a valid URL and template when running mesheryctl model generate then the model is generated" {
    run $MESHERYCTL_BIN model generate --file "https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.crds.yaml" --template "$FIXTURES_DIR/valid-template.json"
    assert_success
    assert_output --partial "Generating model from URL: https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.crds.yaml"
    assert_output --partial "Imported model cert-manager"
    assert_output --partial "Model can be accessed from"
}

@test "[cut=Model][tg=Model Generate] given a URL without a template when running mesheryctl model generate then an error message is displayed" {
    run $MESHERYCTL_BIN model generate --file "https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.crds.yaml"
    assert_failure
    assert_output --partial "no template file is provided while using url for importing a model"
}

@test "[cut=Model][tg=Model Generate] given a non-existent template when running mesheryctl model generate then an error message is displayed" {
    run $MESHERYCTL_BIN model generate --file "https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.crds.yaml" --template "nonexistent-template.json"
    assert_failure
    assert_output --partial "open nonexistent-template.json: no such file or directory"
}

@test "[cut=Model][tg=Model Generate] given a valid CSV directory when running mesheryctl model generate then the model is generated" {
    run $MESHERYCTL_BIN model generate --file "$CSV_FIXTURES_DIR"
    assert_success
    assert_output --partial "Generating model from CSV files"
    assert_output --partial "Model can be accessed from"
    assert_output --partial "Logs for the csv generation can be accessed"
}

@test "[cut=Model][tg=Model Generate] given an invalid CSV directory when running mesheryctl model generate then an error message is displayed" {
    run $MESHERYCTL_BIN model generate --file "invalid-dir"
    assert_failure
    assert_output --partial "open invalid-dir: no such file or directory"
}
