#!/usr/bin/env bats
setup() {
    export ORIGINAL_HOME="$HOME"
    HOME="$TEMP_DATA_DIR"

    mkdir -p "$HOME"


    load "$E2E_HELPERS_PATH/bats_libraries"
    _load_bats_libraries

    load "$E2E_HELPERS_PATH/constants"

    export FAKE_ORGID="11223344-1122-1122-1122-112233445566"
}

@test "mesheryctl exp workspace create --orgId --name --description" {
    run $MESHERYCTL_BIN exp workspace create --orgId $FAKE_ORGID --name ghj --description ibhv


}

