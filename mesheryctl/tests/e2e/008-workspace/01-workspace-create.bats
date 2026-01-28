#!/usr/bin/env bats
setup() {
    load "$E2E_HELPERS_PATH/bats_libraries"
    load "$E2E_HELPERS_PATH/constants"
    _load_bats_libraries

    export TESTDATA_DIR="$TEMP_DATA_DIR/testdata/workspace"
    mkdir -p "$TESTDATA_DIR"
}

