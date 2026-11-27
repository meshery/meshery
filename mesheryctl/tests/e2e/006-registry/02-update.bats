#!/usr/bin/env bats

setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    # allow CI/environment to provide integration spreadsheet credentials
    export TEST_SPREADSHEET_ID="${TEST_SPREADSHEET_ID:-}"
    export TEST_SPREADSHEET_CRED="${TEST_SPREADSHEET_CRED:-}"

    # registry update rewrites every component JSON the spreadsheet names
    # beneath its input root, so it must never be pointed at the repository's
    # tracked models tree. Each test gets a disposable root that it owns.
    export MODELS_DIR="$TEMP_DATA_DIR/registry-update/$BATS_TEST_NAME"
    mkdir -p "$MODELS_DIR"
}

teardown() {
    rm -rf "$MODELS_DIR"
}

# Copies one version of a tracked model into the disposable input root. The copy
# is independent of the source, so the command under test cannot reach the
# tracked models tree.
seed_model_fixture() {
    local model="$1"
    local source_root="../../../models/$model"

    # A model root holds one directory per version, and may also hold loose
    # files such as model_template.json that registry update ignores.
    local version
    version="$(cd "$source_root" 2>/dev/null && ls -d */ 2>/dev/null | head -1 | tr -d '/')"
    if [[ -z "$version" ]]; then
        fail "no model version directory found under $source_root"
    fi

    mkdir -p "$MODELS_DIR/$model"
    cp -R "$source_root/$version" "$MODELS_DIR/$model/$version"
}

require_spreadsheet_credentials() {
    if [[ -z "$TEST_SPREADSHEET_ID" || -z "$TEST_SPREADSHEET_CRED" ]]; then
        skip "Spreadsheet credentials not configured"
    fi
}

@test "given no arguments when running mesheryctl registry update then an error is displayed" {
    run $MESHERYCTL_BIN registry update
    assert_failure
    assert_output --partial "error updating registry"
}

@test "given spreadsheet-id without spreadsheet-cred when running mesheryctl registry update then an error about missing flag is displayed" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "test-id"
    assert_failure
    assert_output --partial "if any flags in the group [spreadsheet-id spreadsheet-cred] are set they must all be set"
    assert_output --partial "spreadsheet-cred"
}

@test "given spreadsheet-cred without spreadsheet-id when running mesheryctl registry update then an error about missing flag is displayed" {
    run $MESHERYCTL_BIN registry update --spreadsheet-cred "test-cred"
    assert_failure
    assert_output --partial "if any flags in the group [spreadsheet-id spreadsheet-cred] are set they must all be set"
    assert_output --partial "spreadsheet-id"
}

@test "given invalid spreadsheet credentials when running mesheryctl registry update then an error is displayed" {
    run $MESHERYCTL_BIN registry update --spreadsheet-id "invalid-id" --spreadsheet-cred "invalid-cred"
    assert_failure
    assert_output --partial "Invalid JWT credentials"
}

@test "given an invalid model name when running mesheryctl registry update then zero models are updated" {
    require_spreadsheet_credentials
    run $MESHERYCTL_BIN registry update --input "$MODELS_DIR" --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --model "nonexistent-model"
    assert_success
    assert_output --partial "Updated 0 models and 0 components"
}

@test "given valid spreadsheet credentials and a model name when running mesheryctl registry update then that model is updated" {
    require_spreadsheet_credentials
    seed_model_fixture "kubernetes"
    run $MESHERYCTL_BIN registry update --input "$MODELS_DIR" --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED" --model "kubernetes"
    assert_success
    assert_output --regexp "Updated [1-9][0-9]* models? and [0-9]+ components?"
}

@test "given valid spreadsheet credentials when running mesheryctl registry update then models are updated successfully" {
    require_spreadsheet_credentials
    seed_model_fixture "kubernetes"
    run $MESHERYCTL_BIN registry update --input "$MODELS_DIR" --spreadsheet-id "$TEST_SPREADSHEET_ID" --spreadsheet-cred "$TEST_SPREADSHEET_CRED"
    assert_success
    assert_output --regexp "Updated [1-9][0-9]* models? and [0-9]+ components?"
}
